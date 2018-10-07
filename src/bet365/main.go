package main

import (
	"bet365/bet365"
	"net/http"
	_ "net/http/pprof"

	"github.com/mysll/toolkit"
)

func main() {
	// go chat.Run(func() error {
	// 	chat.SendMessage("微信接入成功")
	// 	go ybf.Run()
	// 	return nil
	// })

	go http.ListenAndServe(":8888", nil)
	go bet365.Run("premws-pt3.365pushodds.com", "https://www.348365365.com", "https://www.348365365.com")
	toolkit.WaitForQuit()
}
