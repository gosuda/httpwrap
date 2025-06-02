// Package httpwrap provides HTTP error handling utilities for the standard net/http package.
// It wraps the standard HTTP multiplexer with enhanced error handling capabilities,
// supporting both simple error types and structured HttpError types with custom content types.
package httpwrap

import (
	"errors"
	"net/http"

	"github.com/gosuda/httpwrap/httperror"
)

// Mux wraps http.ServeMux with enhanced error handling capabilities.
// It provides automatic error handling and supports custom error callbacks.
type Mux struct {
	mux           *http.ServeMux
	errorCallback func(err error)
}

// NewMux creates a new Mux with the specified error callback function.
// If errorCallback is nil, a no-op function is used.
func NewMux(errorCallback func(err error)) *Mux {
	if errorCallback == nil {
		errorCallback = func(err error) {}
	}
	return &Mux{
		mux:           http.NewServeMux(),
		errorCallback: errorCallback,
	}
}

// HandlerFunc defines a function signature for HTTP handlers that can return errors.
// This allows for cleaner error handling in HTTP handlers.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Handle registers a new handler for the given pattern with automatic error handling.
// If the handler returns an error, it will be automatically converted to an appropriate HTTP response.
func (m *Mux) Handle(pattern string, handler HandlerFunc) {
	m.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
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
			m.errorCallback(err)
		}
	})
}

// ServeHTTP implements the http.Handler interface, delegating to the underlying ServeMux.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
