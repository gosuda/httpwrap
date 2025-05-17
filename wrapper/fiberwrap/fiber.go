package fiberwrap

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/gosuda/httpwrap/httperror"
)

type Wrapper struct {
	app *fiber.App
}

func NewWrapper() *Wrapper {
	return &Wrapper{
		app: fiber.New(),
	}
}

func WithApp(app *fiber.App) *Wrapper {
	return &Wrapper{
		app: app,
	}
}

type HandlerFunc func(c *fiber.Ctx) error

func (a *Wrapper) Handle(method, path string, handler HandlerFunc) {
	a.app.Add(method, path, func(c *fiber.Ctx) error {
		if err := handler(c); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				return c.Status(he.Code).SendString(he.ErrorMessage())
			case false:
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}
		return nil
	})
}

func (a *Wrapper) Get(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodGet, path, handler)
}

func (a *Wrapper) Post(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodPost, path, handler)
}

func (a *Wrapper) Put(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodPut, path, handler)
}

func (a *Wrapper) Delete(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodDelete, path, handler)
}

func (a *Wrapper) Patch(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodPatch, path, handler)
}

func (a *Wrapper) Options(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodOptions, path, handler)
}

func (a *Wrapper) Head(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodHead, path, handler)
}

func (a *Wrapper) Connect(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodConnect, path, handler)
}

func (a *Wrapper) Trace(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodTrace, path, handler)
}

func (a *Wrapper) App() *fiber.App {
	return a.app
}
