package controller

import (
	"bet365/odd"

	"github.com/lunny/tango"
)

type Odd struct {
	RenderBase
	tango.Json
}

func (u *Odd) Get() interface{} {
	u.Ctx.Header().Add("Access-Control-Allow-Origin", "*") //允许访问所有域
	o := u.Ctx.ParamFloat64(":odd", 0)
	result := odd.GetOdd(o)
	return map[string]interface{}{
		"Status": 200,
		"Odd":    result,
	}
}
