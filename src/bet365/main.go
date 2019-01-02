package main

import (
	"bet365/bet365"
	"chat"
	"config"
	"controller"
	"fmt"
	"html/template"
	"os"
	"ssq"

	"github.com/lunny/log"
	"github.com/lunny/tango"
	"github.com/mysll/toolkit"
	"github.com/robfig/cron"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
)

func Serv() {
	l := log.New(os.Stdout, "[tango] ", log.Ldefault())
	l.SetOutputLevel(log.Lfatal)
	t := tango.Classic(l)
	t.Use(
		events.Events(),
		tango.Static(tango.StaticOptions{
			RootPath: "./views/static",
			Prefix:   "static",
		}),
		renders.New(renders.Options{
			Reload:    true,
			Directory: "./views/templates",
			Funcs:     template.FuncMap{},
			Charset:   "UTF-8", // Appends the given charset to the Content-Type header. Default is UTF-8
			// Allows changing of output to XHTML instead of HTML. Default is "text/html"
			HTMLContentType: "text/html",
			DelimsLeft:      "<<<",
			DelimsRight:     ">>>", // default Delims is {{}}, if it conflicts with your javascript template such as angluar, you can change it.
		}),
	)

	t.Get("/", new(controller.Index))
	t.Get("/lucky", new(controller.Lucky))
	t.Get("/update", new(controller.Update))
	t.Get("/odd/:odd", new(controller.Odd))
	t.Post("/qqbot", new(controller.QQBridge))
	t.Run(8888)
}

var C = cron.New()

func luncky() string {
	balls := ssq.Millionaire()
	str := fmt.Sprintf("上期%d:\n %02d %02d %02d %02d %02d %02d | %02d\n", balls.Last.Expect,
		balls.Last.Red[0],
		balls.Last.Red[1],
		balls.Last.Red[2],
		balls.Last.Red[3],
		balls.Last.Red[4],
		balls.Last.Red[5],
		balls.Last.Blue)
	str += "本期推荐:\n"
	for _, v := range balls.Lucky {
		str += fmt.Sprintf("  %02d %02d %02d %02d %02d %02d | %02d\n",
			v.Red[0],
			v.Red[1],
			v.Red[2],
			v.Red[3],
			v.Red[4],
			v.Red[5],
			v.Blue)
	}
	return str
}

func main() {
	config.LoadConfig()
	C.AddFunc("0 0 12 * * *", func() {
		chat.SendToRecommend(bet365.Stat())
	})
	C.AddFunc("0 0 18 * * 0,2,4", func() {
		chat.SendToRecommend(luncky())
	})
	C.AddFunc("0 0 22 * * 0,2,4", func() {
		h := ssq.Histroy(true)
		if len(h) > 0 {
			chat.SendToRecommend(fmt.Sprintf("update ok, %v", h[0]))
		}
	})
	C.Start()
	go Serv()
	go bet365.Run(config.Setting.Bet365.WSURL, config.Setting.Bet365.Host, config.Setting.Bet365.Host)
	go chat.MessageLoop()
	toolkit.WaitForQuit()
}
