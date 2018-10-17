package controller

import (
	"bet365/bet365"
	"bet365/odd"
	"config"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type QQBridge struct {
	RenderBase
}

func (u *QQBridge) Post() string {
	u.Ctx.Header().Add("Access-Control-Allow-Origin", "*") //允许访问所有域

	body, err := u.Ctx.Body()
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	values, err := url.ParseQuery(string(body))

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	cmd, args := parseCommand(values.Get("content"))
	log.Println("[bridge]", cmd)
	result := do(values, cmd, args)
	//chat.SendQQMessage(result, values.Get("group"))
	return result
}

func parseCommand(content string) (cmd string, args []string) {
	str := strings.Replace(content, "[@ME]", "", 1)
	str = strings.TrimRight(strings.TrimLeft(str, " "), " ")
	cmds := strings.Split(str, " ")
	return cmds[0], cmds[1:]
}

func do(values url.Values, cmd string, args []string) string {
	switch cmd {
	case "help", "帮助":
		return "[所有命令]:\n  [odd|赔率]:返回欧亚转换，odd 1.5\n  [ssq|双色球]:双色球五注\n  [time|定时]:比赛定时提醒，time id 35\n  [size|大球]:大小球提醒，size id 0.5 2.0"
	case "odd", "赔率":
		if len(args) != 1 {
			return "[error] 参数错误"
		}
		odds := odd.GetOddStr(args[0])
		return strings.Join(odds, "\n")
	case "ssq", "双色球":
		balls := Millionaire()
		str := fmt.Sprintf("上期%d:\n %v %d\n", balls.Last.Expect, balls.Last.Red, balls.Last.Blue)
		str += "本期推荐:\n"
		for _, v := range balls.Lucky {
			str += fmt.Sprintf("  %v %d\n", v.Red, v.Blue)
		}
		return str
	case "reload":
		config.LoadConfig()
		return "reload ok"
	case "update":
		lh := histroy(true)
		return fmt.Sprintf("update ok, %v", lh[0])
	case "time", "定时":
		if len(args) != 2 {
			return "[error] 参数错误， example:time id 35"
		}

		return bet365.AddTimeNotify(values.Get("group"), values.Get("from"), args[0], args[1])
	case "size", "大球":
		if len(args) != 3 {
			return "[error] 参数错误， example:size id 0.5 2.5"
		}
		return bet365.AddSizeNotify(values.Get("group"), values.Get("from"), args[0], args[1], args[2])
	case "stat", "统计":
		return bet365.Stat()
	case "订阅":
		g := values.Get("group")
		for _, v := range config.Setting.Recommend {
			if v == g {
				return "[error] 群名称冲突，请修改群名称"
			}
		}
		config.Setting.Recommend = append(config.Setting.Recommend, g)
		return "增加订阅成功，将收到比赛推荐信息"
	default:
		return "[error]你说什么我听不懂,输入help或者帮助查看所有支持的命令"
	}
}
