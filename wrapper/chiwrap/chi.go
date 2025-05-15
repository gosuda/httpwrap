package chiwrap

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gosuda.org/httpwrap/httperror"
)

type Router struct {
	router      chi.Router
	errCallback func(err error)
}

func NewRouter(errCallback func(err error)) *Router {
	if errCallback == nil {
		errCallback = func(err error) {}
	}
	return &Router{
		router:      chi.NewRouter(),
		errCallback: errCallback,
	}
}

type HandlerFunc func(writer http.ResponseWriter, request *http.Request) error

func (r *Router) Handle(pattern string, handler HandlerFunc) {
	r.router.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Get(pattern string, handler HandlerFunc) {
	r.router.Get(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Post(pattern string, handler HandlerFunc) {
	r.router.Post(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Put(pattern string, handler HandlerFunc) {
	r.router.Put(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Delete(pattern string, handler HandlerFunc) {
	r.router.Delete(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Patch(pattern string, handler HandlerFunc) {
	r.router.Patch(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Options(pattern string, handler HandlerFunc) {
	r.router.Options(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Head(pattern string, handler HandlerFunc) {
	r.router.Head(pattern, func(writer http.ResponseWriter, request *http.Request) {
		if err := handler(writer, request); err != nil {
			he := &httperror.HttpError{}
			switch errors.As(err, &he) {
			case true:
				http.Error(writer, he.Message, he.Code)
			case false:
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			r.errCallback(err)
		}
	})
}

func (r *Router) Route(pattern string, callback func(r *Router)) {
	r.router.Route(pattern, func(router chi.Router) {
		callback(&Router{
			router:      router,
			errCallback: r.errCallback,
		})
	})
}

func (r *Router) Mount(pattern string, subRouter http.Handler) {
	r.router.Mount(pattern, subRouter)
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	r.router.ServeHTTP(writer, reader)
}
