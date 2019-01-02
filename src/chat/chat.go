package chat

import (
	"bytes"
	"config"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type msginfo struct {
	msg   string
	group string
}

const (
	FREQUENCE = time.Second * 5
)

var (
	msgqueue = make(map[string]chan *msginfo, 10)
)

func SendToRecommend(msg string) {
	for _, g := range config.Setting.Recommend {
		SendMessage(msg, g)
	}
}

func SendToBroadcast(msg string) {
	for _, g := range config.Setting.Broadcast {
		SendMessage(msg, g)
	}
}

// 发送消息
func SendMessage(msg string, group string) {
	m := &msginfo{msg, group}
	select {
	case msgqueue[group] <- m:
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
	for _, g := range config.Setting.Broadcast {
		if _, ok := msgqueue[g]; ok {
			continue
		}
		msgqueue[g] = make(chan *msginfo, 128)
		go loop(msgqueue[g])
	}
	for _, g := range config.Setting.Recommend {
		if _, ok := msgqueue[g]; ok {
			continue
		}
		msgqueue[g] = make(chan *msginfo, 128)
		go loop(msgqueue[g])
	}
}

func loop(ch chan *msginfo) {
	for {
		select {
		case m := <-ch:
			SendDingTalk(m.msg, m.group)
			time.Sleep(FREQUENCE)
		}
	}
}
