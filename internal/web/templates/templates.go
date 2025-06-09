package templates

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) Render(w io.Writer, name string, binding interface{}, layouts ...string) error {
	component, ok := binding.(templ.Component)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid view component")
	}
	return component.Render(context.Background(), w)
}
