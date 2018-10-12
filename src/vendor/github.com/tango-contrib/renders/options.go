// Copyright 2017 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renders

import (
	"html/template"
	"net/http"
	"path"
)

// Options is a struct for specifying configuration options for the render.Renderer middleware
type Options struct {
	// if reload templates
	Reload bool
	// Directory to load templates. Default is "templates"
	Directory string
	// Extensions to parse template files from. Defaults to [".tmpl"]
	Extensions []string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	Funcs template.FuncMap
	// Vars is a data map for global
	Vars T
	// Appends the given charset to the Content-Type header. Default is "UTF-8".
	Charset string
	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
	HTMLContentType string
	// default Delims
	DelimsLeft, DelimsRight string
	// where the file stored
	FileSystem http.FileSystem
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.Extensions) == 0 {
		opt.Extensions = []string{".html"}
	}
	if len(opt.HTMLContentType) == 0 {
		opt.HTMLContentType = ContentHTML
	}
	if len(opt.DelimsLeft) == 0 {
		opt.DelimsLeft = "{{"
	}
	if len(opt.DelimsRight) == 0 {
		opt.DelimsRight = "}}"
	}
	if opt.FileSystem == nil {
		opt.FileSystem = http.Dir(opt.Directory)
	}

	return opt
}

// IsExtMatch is a file name match the ext
func (o Options) IsExtMatch(fileName string) bool {
	ext := path.Ext(fileName)
	for _, extension := range o.Extensions {
		if ext == extension {
			return true
		}
	}
	return false
}
