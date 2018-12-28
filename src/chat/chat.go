package chat

import (
	"bytes"
	"config"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
)

type msginfo struct {
	msg   string
	group string
}

const (
	FREQUENCE = time.Second * 5
)

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

type Content struct {
	Content string `json:"content"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingTalkMsg struct {
	MsgType string  `json:"msgtype"`
	Text    Content `json:"text"`
	At      At      `json:"at"`
}

func makeDingTalkMsg(msg string) []byte {
	dmsg := DingTalkMsg{}
	dmsg.MsgType = "text"
	dmsg.Text.Content = msg
	dmsg.At.IsAtAll = false
	data, err := json.Marshal(dmsg)
	if err != nil {
		return nil
	}
	return data
}

func SendDingTalk(msg, webhook string) {
	data := makeDingTalkMsg(msg)
	if data == nil {
		return
	}

	for {
		resp, err := http.Post(webhook, "application/json", bytes.NewReader(data))
		if err != nil {
			time.Sleep(time.Second)
			log.Println(err)
			continue
		}

		resp.Body.Close()
		break
	}

}

func MessageLoop() {
	for {
		select {
		case m := <-msgqueue:
			SendDingTalk(m.msg, m.group)
			time.Sleep(FREQUENCE)
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
