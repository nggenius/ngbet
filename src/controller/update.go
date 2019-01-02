package controller

import (
	"ssq"

	"github.com/lunny/tango"
)

type Update struct {
	RenderBase
	tango.Json
}

func (u *Update) Get() interface{} {
	u.Ctx.Header().Add("Access-Control-Allow-Origin", "*") //允许访问所有域
	lh := ssq.Histroy(true)
	return map[string]interface{}{
		"Status": 200,
		"Last":   lh[0],
	}
}
