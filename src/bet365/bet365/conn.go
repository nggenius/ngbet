package bet365

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/levigross/grequests"
)

const (
	_ENCODINGS_NONE                  = "\x00"
	_TYPES_SUBSCRIBE                 = "\x16"
	_DELIMITERS_HANDSHAKE            = "\x03"
	_TYPES_TOPIC_STATUS_NOTIFICATION = "\x23"

	_DELIMITERS_RECORD  = "\x01"
	_DELIMITERS_MESSAGE = "\x08"

	_DELIMITERS_SPLIT = "\x7c"
)

var messagesSessionId = fmt.Sprintf("%s%sP%s__time,S_%%s%s",
	_TYPES_TOPIC_STATUS_NOTIFICATION,
	_DELIMITERS_HANDSHAKE,
	_DELIMITERS_RECORD,
	_ENCODINGS_NONE,
)

func genSessionId(sessionId string) string {
	return fmt.Sprintf(messagesSessionId, sessionId)
}

func genSubscription(item string) string {
	messagesSubscription := fmt.Sprintf("%s%s%%s%s", _TYPES_SUBSCRIBE, _ENCODINGS_NONE, _DELIMITERS_RECORD)
	s := fmt.Sprintf(messagesSubscription, item)
	return s
}

type bet365conn struct {
	c *websocket.Conn
}

func NewConn() *bet365conn {
	conn := new(bet365conn)
	return conn
}

func (b *bet365conn) Connect(addr string, origin string, getcookieurl string) error {
	var dialer = &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		EnableCompression: true,
		Subprotocols:      []string{"zap-protocol-v1"},
	}

	u := url.URL{Scheme: "wss", Host: addr, Path: "/zap/"}
	q := u.Query()
	u.RawQuery = q.Encode()

	log.Printf("connecting to %s", u.String())

	var header = http.Header{
		"Origin":                 []string{origin},
		"Sec-WebSocket-Protocol": []string{"zap-protocol-v1"},
		"User-Agent":             []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/5.0 (KHTML, like Gecko) Chrome/5.0 Safari/5.0"},
	}

	c, resp, err := dialer.Dial(u.String(), header)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("connect failed, code:%d", resp.StatusCode)
	}

	b.c = c

	r, err := grequests.Get(getcookieurl, nil)
	if err != nil {
		return err
	}
	res := r.RawResponse.Cookies()
	log.Println("Sessionid=" + res[1].Value)
	err = b.sendMessage([]byte(genSessionId(res[1].Value)))
	if err != nil {
		return err
	}

	return nil
}

func (b *bet365conn) sendMessage(data []byte) error {
	if b.c == nil {
		return fmt.Errorf("wss not connected")
	}

	log.Printf("send:%s\n", data)
	return b.c.WriteMessage(websocket.TextMessage, data)
}

func (b *bet365conn) subscibe(topic string) error {
	err := b.sendMessage([]byte(genSubscription(topic)))
	return err
}

func (b *bet365conn) ReadMessage() (messageType int, p []byte, err error) {
	return b.c.ReadMessage()
}

func (b *bet365conn) close() {
	if b.c != nil {
		b.c.Close()
	}
}
