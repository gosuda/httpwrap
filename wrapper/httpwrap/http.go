package httpwrap

import (
	"errors"
	"net/http"

	"github.com/gosuda/httpwrap/httperror"
)

type Mux struct {
	mux           *http.ServeMux
	errorCallback func(err error)
}

func NewMux(errorCallback func(err error)) *Mux {
	if errorCallback == nil {
		errorCallback = func(err error) {}
	}
	return &Mux{
		mux:           http.NewServeMux(),
		errorCallback: errorCallback,
	}
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

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

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
