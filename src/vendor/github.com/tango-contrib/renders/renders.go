// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renders

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/lunny/tango"
	"github.com/oxtoacart/bpool"
)

// Version return the middleware's version
func Version() string {
	return "0.4.0329"
}

const (
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	ContentHTML    = "text/html"
	ContentXHTML   = "application/xhtml+xml"
	defaultCharset = "UTF-8"
)

// Provides a common buffer to execute templates.
type T map[string]interface{}

func (t T) Merge(at T) T {
	if len(at) <= 0 {
		return t
	}

	for k, v := range at {
		t[k] = v
	}
	return t
}

type Renders struct {
	Options
	cs        string
	pool      *bpool.BufferPool
	templates map[string]*template.Template
}

func New(options ...Options) *Renders {
	opt := prepareOptions(options)
	t, err := compile(opt)
	if err != nil {
		panic(err)
	}
	return &Renders{
		Options:   opt,
		cs:        prepareCharset(opt.Charset),
		pool:      bpool.NewBufferPool(64),
		templates: t,
	}
}

type IRenderer interface {
	SetRenderer(*Renders, *tango.Context, func(string), func(string), func(string, io.Reader))
}

// confirm Renderer implements IRenderer
var _ IRenderer = &Renderer{}

type Renderer struct {
	ctx                     *tango.Context
	renders                 *Renders
	before, after           func(string)
	afterBuf                func(string, io.Reader)
	compiledCharset         string
	Charset                 string
	HTMLContentType         string
	delimsLeft, delimsRight string
}

func (r *Renderer) SetRenderer(renders *Renders, ctx *tango.Context,
	before, after func(string), afterBuf func(string, io.Reader)) {
	r.renders = renders
	r.ctx = ctx
	r.before = before
	r.after = after
	r.afterBuf = afterBuf
	r.HTMLContentType = renders.Options.HTMLContentType
	r.compiledCharset = renders.cs
	r.delimsLeft = renders.Options.DelimsLeft
	r.delimsRight = renders.Options.DelimsRight
}

type Before interface {
	BeforeRender(string)
}

type After interface {
	AfterRender(string)
}

type AfterBuf interface {
	AfterRender(string, io.Reader)
}

func (r *Renders) RenderBytes(name string, bindings ...interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := r.Render(buf, name, bindings...)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *Renders) Render(w io.Writer, name string, bindings ...interface{}) error {
	var binding interface{}
	if len(bindings) > 0 {
		binding = bindings[0]
	}
	if t, ok := binding.(T); ok {
		binding = t.Merge(r.Options.Vars)
	}

	if r.Reload {
		var err error
		// recompile for easy development
		r.templates, err = compile(r.Options)
		if err != nil {
			return err
		}
	}

	buf, err := r.execute(name, binding)
	if err != nil {
		r.pool.Put(buf)
		return err
	}

	// template rendered fine, write out the result
	_, err = io.Copy(w, buf)
	r.pool.Put(buf)
	return err
}

func (r *Renders) execute(name string, binding interface{}) (*bytes.Buffer, error) {
	buf := r.pool.Get()
	name = alignTmplName(name)

	if rt, ok := r.templates[name]; ok {
		return buf, rt.ExecuteTemplate(buf, name, binding)
	}
	return buf, errors.New("template is not exist")
}

func (r *Renders) Handle(ctx *tango.Context) {
	if action := ctx.Action(); action != nil {
		if rd, ok := action.(IRenderer); ok {
			var before, after func(string)
			var afterBuf func(string, io.Reader)
			if b, ok := action.(Before); ok {
				before = b.BeforeRender
			}
			if a, ok := action.(After); ok {
				after = a.AfterRender
			}
			if a2, ok := action.(AfterBuf); ok {
				afterBuf = a2.AfterRender
			}

			rd.SetRenderer(r, ctx, before, after, afterBuf)
		}
	}

	ctx.Next()
}

func compile(options Options) (map[string]*template.Template, error) {
	return Load(options)
}

func prepareCharset(charset string) string {
	if len(charset) != 0 {
		return "; charset=" + charset
	}

	return "; charset=" + defaultCharset
}

// Render a template
//     r.Render("index.html")
//     r.Render("index.html", renders.T{
//                "name": value,
//           })
func (r *Renderer) Render(name string, bindings ...interface{}) error {
	return r.StatusRender(http.StatusOK, name, bindings...)
}

// RenderBytes Will not called before & after method.
func (r *Renderer) RenderBytes(name string, binding ...interface{}) ([]byte, error) {
	return r.renders.RenderBytes(name, binding...)
}

func (r *Renderer) StatusRender(status int, name string, bindings ...interface{}) error {
	var binding interface{}
	if len(bindings) > 0 {
		binding = bindings[0]
	}
	if t, ok := binding.(T); ok {
		binding = t.Merge(r.renders.Options.Vars)
	}

	if r.renders.Reload {
		var err error
		// recompile for easy development
		r.renders.templates, err = compile(r.renders.Options)
		if err != nil {
			return err
		}
	}

	buf, err := r.execute(name, binding)
	if err != nil {
		r.renders.pool.Put(buf)
		return err
	}

	var cs string
	if len(r.Charset) > 0 {
		cs = prepareCharset(r.Charset)
	} else {
		cs = r.compiledCharset
	}
	// template rendered fine, write out the result
	r.ctx.Header().Set(ContentType, r.HTMLContentType+cs)
	r.ctx.WriteHeader(status)
	_, err = io.Copy(r.ctx.ResponseWriter, buf)
	r.renders.pool.Put(buf)
	return err
}

func (r *Renderer) Template(name string) *template.Template {
	return r.renders.templates[alignTmplName(name)]
}

func (r *Renderer) execute(name string, binding interface{}) (*bytes.Buffer, error) {
	buf := r.renders.pool.Get()
	if r.before != nil {
		r.before(name)
	}
	if r.after != nil {
		defer r.after(name)
	}

	name = alignTmplName(name)

	if rt, ok := r.renders.templates[name]; ok {
		err := rt.Delims(r.delimsLeft, r.delimsRight).ExecuteTemplate(buf, name, binding)
		if err == nil && r.afterBuf != nil {
			var tmpBuf = bytes.NewBuffer(buf.Bytes())
			r.afterBuf(name, tmpBuf)
		}
		return buf, err
	}
	if r.afterBuf != nil {
		r.afterBuf(name, nil)
	}
	return buf, fmt.Errorf("template %s is not exist", name)
}
