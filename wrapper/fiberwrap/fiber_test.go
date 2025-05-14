package fiberwrap_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/snowmerak/httpwrap/httperror"
	"github.com/snowmerak/httpwrap/wrapper/fiberwrap"
)

func TestNewWrapper(t *testing.T) {
	w := fiberwrap.NewWrapper()
	w.Get("/echo", func(c *fiber.Ctx) error {
		name := c.Query("name")
		if name == "" {
			return httperror.BadRequest("name is required")
		}

		return c.SendString("Hello " + name)
	})

	go func() {
		if err := w.App().Listen(":8080"); err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()

	time.Sleep(250 * time.Millisecond) // Wait for server to start

	func() {
		resp, err := http.Get("http://localhost:8080/echo?name=John")
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
			t.Fatalf("Expected 'Hello John', got '%s'", data)
		}
	}()

	func() {
		resp, err := http.Get("http://localhost:8080/echo")
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
			t.Fatalf("Expected 'name is required', got '%s'", data)
		}
	}()
}
