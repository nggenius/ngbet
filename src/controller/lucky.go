package controller

import (
	"ssq"

	"github.com/lunny/tango"
)

type Lucky struct {
	RenderBase
	tango.Json
}

func (l *Lucky) Get() interface{} {
	l.Ctx.Header().Add("Access-Control-Allow-Origin", "*") //允许访问所有域
	res := ssq.Millionaire()
	return map[string]interface{}{
		"Status": res.Status,
		"Last":   res.Last,
		"Lucky":  res.Lucky,
	}
}
