package controller

import (
	"bet365/odd"
	"chat"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type QQBridge struct {
	RenderBase
}

func (u *QQBridge) Post() {
	u.Ctx.Header().Add("Access-Control-Allow-Origin", "*") //允许访问所有域

	body, err := u.Ctx.Body()
	if err != nil {
		log.Println(err)
		return
	}

	values, err := url.ParseQuery(string(body))

	if err != nil {
		log.Println(err)
		return
	}

	cmd, args := parseCommand(values.Get("content"))
	log.Println("[bridge]", cmd)
	result := do(cmd, args)
	chat.SendQQMessage(result, values.Get("group"))
	return
}

func parseCommand(content string) (cmd string, args []string) {
	str := strings.Replace(content, "[@ME]", "", 1)
	str = strings.TrimRight(strings.TrimLeft(str, " "), " ")
	cmds := strings.Split(str, " ")
	return cmds[0], cmds[1:]
}

func do(cmd string, args []string) string {
	switch cmd {
	case "help", "帮助":
		return "[所有命令]:\n\todd|赔率:返回欧亚转换，odd 1.5\n\tssq|双色球:双色球五注\n更多功能增加中"
	case "odd", "赔率":
		if len(args) != 1 {
			return "[error] 参数错误"
		}
		odds := odd.GetOddStr(args[0])
		return strings.Join(odds, "\n")
	case "ssq", "双色球":
		balls := Millionaire()
		str := fmt.Sprintf("上期%d:\n 红:%v  蓝:%d\n", balls.Last.Expect, balls.Last.Red, balls.Last.Blue)
		str += "本期推荐:\n"
		for _, v := range balls.Lucky {
			str += fmt.Sprintf("\t红:%v 蓝:%d\n", v.Red, v.Blue)
		}
		return str
	default:
		return "[error]你说什么我听不懂,输入help或者帮助查看所有支持的命令"
	}
}
