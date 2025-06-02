// Package chiwrap provides HTTP error handling utilities for the go-chi/chi router.
// It wraps the chi router with enhanced error handling capabilities,
// supporting both simple error types and structured HttpError types with custom content types.
package chiwrap

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gosuda/httpwrap/httperror"
)

// Router wraps chi.Router with enhanced error handling capabilities.
// It provides automatic error handling and supports custom error callbacks.
type Router struct {
	router      chi.Router
	errCallback func(err error)
}

// NewRouter creates a new Router with the specified error callback function.
// If errCallback is nil, a no-op function is used.
func NewRouter(errCallback func(err error)) *Router {
	if errCallback == nil {
		errCallback = func(err error) {}
	}
	return &Router{
		router:      chi.NewRouter(),
		errCallback: errCallback,
	}
}

// HandlerFunc defines a function signature for HTTP handlers that can return errors.
// This allows for cleaner error handling in HTTP handlers with chi router.
type HandlerFunc func(writer http.ResponseWriter, request *http.Request) error

// handleError processes an error and sends appropriate HTTP response with content type support.
func (r *Router) handleError(writer http.ResponseWriter, err error) {
	he := &httperror.HttpError{}
	switch errors.As(err, &he) {
	case true:
		// Set Content-Type if specified in HttpError
		if he.ContentType != "" {
			writer.Header().Set("Content-Type", he.ContentType)
			writer.WriteHeader(he.Code)
			writer.Write([]byte(he.Message))
		} else {
			http.Error(writer, he.Message, he.Code)
		}
	case false:
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	r.errCallback(err)
}

// Handle registers a new handler for the given pattern with automatic error handling.
// If the handler returns an error, it will be automatically converted to an appropriate HTTP response.
func (r *Router) Handle(pattern string, handler HandlerFunc) {
	r.router.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Get registers a new GET handler for the given pattern with automatic error handling.
func (r *Router) Get(pattern string, handler HandlerFunc) {
	r.router.Get(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Post registers a new POST handler for the given pattern with automatic error handling.
func (r *Router) Post(pattern string, handler HandlerFunc) {
	r.router.Post(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Put registers a new PUT handler for the given pattern with automatic error handling.
func (r *Router) Put(pattern string, handler HandlerFunc) {
	r.router.Put(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Delete registers a new DELETE handler for the given pattern with automatic error handling.
func (r *Router) Delete(pattern string, handler HandlerFunc) {
	r.router.Delete(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Patch registers a new PATCH handler for the given pattern with automatic error handling.
func (r *Router) Patch(pattern string, handler HandlerFunc) {
	r.router.Patch(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Options registers a new OPTIONS handler for the given pattern with automatic error handling.
func (r *Router) Options(pattern string, handler HandlerFunc) {
	r.router.Options(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Head registers a new HEAD handler for the given pattern with automatic error handling.
func (r *Router) Head(pattern string, handler HandlerFunc) {
	r.router.Head(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			r.handleError(writer, err)
		}
	})
}

// Route creates a new sub-router for the given pattern.
// The callback function receives a new Router instance that inherits the error callback.
func (r *Router) Route(pattern string, callback func(r *Router)) {
	r.router.Route(pattern, func(router chi.Router) {
		callback(&Router{
			router:      router,
			errCallback: r.errCallback,
		})
	})
}

// Mount attaches a sub-router or http.Handler to the routing pattern.
func (r *Router) Mount(pattern string, subRouter http.Handler) {
	r.router.Mount(pattern, subRouter)
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	r.router.ServeHTTP(writer, reader)
}
