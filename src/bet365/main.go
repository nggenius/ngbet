package main

import (
	"bet365/bet365"
	"controller"
	"html/template"

	"github.com/lunny/tango"
	"github.com/mysll/toolkit"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
)

func Serv() {
	t := tango.Classic()
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
	t.Run(8888)
}

func main() {
	// go chat.Run(func() error {
	// 	chat.SendMessage("微信接入成功")
	// 	go ybf.Run()
	// 	return nil
	// })

	go Serv()
	go bet365.Run("premws-pt3.365pushodds.com", "https://www.348365365.com", "https://www.348365365.com")

	toolkit.WaitForQuit()
}
