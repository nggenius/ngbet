package controller

type Index struct {
	RenderBase
}

func (i *Index) Get() error {
	return i.Render("index.html")
}
