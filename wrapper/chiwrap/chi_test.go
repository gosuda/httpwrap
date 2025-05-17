package chiwrap_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gosuda/httpwrap/httperror"
	"github.com/gosuda/httpwrap/wrapper/chiwrap"
)

func TestNewRouter(t *testing.T) {
	r := chiwrap.NewRouter(func(err error) {
		log.Printf("Router log test: Error occured: %v", err)
	})
	r.Get("/echo", func(writer http.ResponseWriter, request *http.Request) error {
		name := request.URL.Query().Get("name")
		if name == "" {
			return httperror.BadRequest("name is required")
		}

		writer.Write([]byte("Hello " + name))
		return nil
	})

	svr := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	time.Sleep(250 * time.Millisecond) // Wait for server to start

	func() {
		hc := http.Client{}
		resp, err := hc.Get("http://localhost:8080/echo?name=John")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(data) != "Hello John" {
			t.Fatalf("Unexpected response body: %s", string(data))
		}
	}()

	func() {
		hc := http.Client{}
		resp, err := hc.Get("http://localhost:8080/echo")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(bytes.TrimSpace(data)) != "name is required" {
			t.Fatalf("Unexpected response body: %s", string(data))
		}
	}()

	svr.Shutdown(context.Background())
}

func TestWithRouter(t *testing.T) {
	r := chiwrap.NewRouter(func(err error) {
		log.Printf("Router log test: Error occured: %v", err)
	})
	r.Route("/echo", func(r *chiwrap.Router) {
		r.Get("/name", func(writer http.ResponseWriter, request *http.Request) error {
			name := request.URL.Query().Get("name")
			if name == "" {
				return httperror.BadRequest("name is required")
			}

			writer.Write([]byte("Hello " + name))
			return nil
		})
	})

	svr := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	time.Sleep(250 * time.Millisecond) // Wait for server to start

	func() {
		hc := http.Client{}
		resp, err := hc.Get("http://localhost:8080/echo/name?name=John")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(data) != "Hello John" {
			t.Fatalf("Unexpected response body: %s", string(data))
		}
	}()

	func() {
		hc := http.Client{}
		resp, err := hc.Get("http://localhost:8080/echo/name")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(bytes.TrimSpace(data)) != "name is required" {
			t.Fatalf("Unexpected response body: %s", string(data))
		}
	}()

	svr.Shutdown(context.Background())
}
