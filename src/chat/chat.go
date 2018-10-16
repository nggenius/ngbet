package chat

import (
	"config"
	"os/exec"
	"time"

	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
)

type msginfo struct {
	msg   string
	group string
}

var (
	session  *wxweb.Session
	msgqueue = make(chan *msginfo, 128)
)

func SendToRecommend(msg string) {
	for _, g := range config.Setting.Recommend {
		SendQQMessage(msg, g)
	}
}

func SendToBroadcast(msg string) {
	for _, g := range config.Setting.Broadcast {
		SendQQMessage(msg, g)
	}
}

// 发送QQ消息
func SendQQMessage(msg string, group string) {
	m := &msginfo{msg, group}
	select {
	case msgqueue <- m:
	default:
	}
}

func MessageLoop() {
	for {
		select {
		case m := <-msgqueue:
			cmd := exec.Command("qq", "send", "group", m.group, m.msg)
			cmd.Run()
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// 发送微信消息
func SendWeChatMessage(msg string) {
	friend := session.Cm.GetContactsByName("bet")
	//logs.Info(friend)
	if len(friend) == 0 {
		return
	}
	session.SendText(msg, session.Bot.UserName, friend[0].UserName)
}

func Run(cb func() error) {
	// 创建session, 一个session对应一个机器人
	// 二维码显示在终端上
	var err error
	session, err = wxweb.CreateSession(nil, nil, wxweb.WEB_MODE)
	if err != nil {
		logs.Error(err)
		return
	}

	session.SetAfterLogin(cb)
	//Register(session)
	// 登录并接收消息
	if err := session.LoginAndServe(false); err != nil {
		logs.Error("session exit, %s", err)
	}
}
