package controller

import (
	"fmt"
	"time"

	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type RenderBase struct {
	renders.Renderer
	tango.Ctx

	startTime time.Time
}

// Before
func (b *RenderBase) Before() {
	b.startTime = time.Now()

}

// After
func (b *RenderBase) After() {

}

// Render 渲染模板
func (b *RenderBase) Render(tmpl string, t ...renders.T) error {
	var ts = renders.T{}
	if len(t) > 0 {
		ts = t[0].Merge(renders.T{})
	}

	ts["costTime"] = fmt.Sprintf("%dms", time.Now().Sub(b.startTime)/1000000)

	return b.Renderer.Render(tmpl, ts)
}
