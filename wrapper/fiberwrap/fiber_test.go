package fiberwrap_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gosuda/httpwrap/httperror"
	"github.com/gosuda/httpwrap/wrapper/fiberwrap"
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
			t.Errorf("Failed to start server: %v", err)
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

func TestWrapper_HandleWithContentType(t *testing.T) {
	w := fiberwrap.NewWrapper()

	// Test with JSON content type
	w.Get("/json", func(c *fiber.Ctx) error {
		return httperror.New(400, `{"error":"bad request"}`, "application/json")
	})

	// Test with RFC9457 problem details
	w.Get("/problem", func(c *fiber.Ctx) error {
		prob := httperror.BadRequestProblem9457("Invalid input data")
		return prob.ToHttpError()
	})

	// Test without content type
	w.Get("/text", func(c *fiber.Ctx) error {
		return httperror.New(404, "Not found")
	})

	go func() {
		if err := w.App().Listen(":8081"); err != nil {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	time.Sleep(250 * time.Millisecond) // Wait for server to start

	// Test JSON content type
	func() {
		resp, err := http.Get("http://localhost:8081/json")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 400 {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Fatalf("Expected content type application/json, got %s", contentType)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(data) != `{"error":"bad request"}` {
			t.Fatalf("Unexpected response body: %s", string(data))
		}
	}()

	// Test RFC9457 problem details content type
	func() {
		resp, err := http.Get("http://localhost:8081/problem")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 400 {
			t.Fatalf("Expected status code 400, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/problem+json" {
			t.Fatalf("Expected content type application/problem+json, got %s", contentType)
		}
	}()

	// Test default content type (no ContentType specified)
	func() {
		resp, err := http.Get("http://localhost:8081/text")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 404 {
			t.Fatalf("Expected status code 404, got %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(data) != "Not found" {
			t.Fatalf("Unexpected response body: %s", string(data))
		}
	}()
}
