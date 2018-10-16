package main

import (
	"bet365/bet365"
	"chat"
	"config"
	"controller"
	"html/template"
	"os"

	"github.com/lunny/log"
	"github.com/lunny/tango"
	"github.com/mysll/toolkit"
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

func main() {
	config.LoadConfig()
	go Serv()
	go bet365.Run(config.Setting.Bet365.WSURL, config.Setting.Bet365.Host, config.Setting.Bet365.Host)
	go chat.MessageLoop()
	toolkit.WaitForQuit()
}
