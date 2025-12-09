package components

type Component interface {
	Render(RenderContext) string
}

type RenderContext struct {
	Width  int
	Height int
}
