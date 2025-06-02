// Package fiberwrap provides HTTP error handling utilities for the Fiber web framework.
// It integrates the httperror package with Fiber's routing and middleware system,
// enabling automatic handling of structured HTTP errors with proper status codes and content types.
package fiberwrap

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/gosuda/httpwrap/httperror"
)

// Wrapper wraps a Fiber application with error handling capabilities.
type Wrapper struct {
	app *fiber.App
}

// NewWrapper creates a new Wrapper with a default Fiber application.
func NewWrapper() *Wrapper {
	return &Wrapper{
		app: fiber.New(),
	}
}

// WithApp creates a new Wrapper with an existing Fiber application.
func WithApp(app *fiber.App) *Wrapper {
	return &Wrapper{
		app: app,
	}
}

// HandlerFunc defines a handler function that can return an error.
type HandlerFunc func(c *fiber.Ctx) error

// Handle registers a handler function for the given HTTP method and path with error handling.
// It automatically converts httperror.HttpError instances to appropriate HTTP responses
// with proper status codes and content types.
func (a *Wrapper) Handle(method, path string, handler HandlerFunc) {
	a.app.Add(method, path, func(c *fiber.Ctx) error {
		if err := handler(c); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				// Set Content-Type if specified in HttpError
				if he.ContentType != "" {
					c.Set("Content-Type", he.ContentType)
				}
				return c.Status(he.Code).SendString(he.ErrorMessage())
			case false:
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}
		return nil
	})
}

// Get registers a GET handler for the given path.
func (a *Wrapper) Get(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodGet, path, handler)
}

// Post registers a POST handler for the given path.
func (a *Wrapper) Post(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodPost, path, handler)
}

// Put registers a PUT handler for the given path.
func (a *Wrapper) Put(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodPut, path, handler)
}

// Delete registers a DELETE handler for the given path.
func (a *Wrapper) Delete(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodDelete, path, handler)
}

// Patch registers a PATCH handler for the given path.
func (a *Wrapper) Patch(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodPatch, path, handler)
}

// Options registers an OPTIONS handler for the given path.
func (a *Wrapper) Options(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodOptions, path, handler)
}

// Head registers a HEAD handler for the given path.
func (a *Wrapper) Head(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodHead, path, handler)
}

// Connect registers a CONNECT handler for the given path.
func (a *Wrapper) Connect(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodConnect, path, handler)
}

// Trace registers a TRACE handler for the given path.
func (a *Wrapper) Trace(path string, handler HandlerFunc) {
	a.Handle(fiber.MethodTrace, path, handler)
}

// App returns the underlying Fiber application instance.
func (a *Wrapper) App() *fiber.App {
	return a.app
}
