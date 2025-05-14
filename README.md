# httpwrap

`httpwrap` is a Go library that provides convenient wrappers for popular HTTP routers (`net/http`, `chi`, `fiber`). It simplifies error handling in HTTP handlers by allowing them to return an `error` type. The wrappers automatically handle these errors, converting them into appropriate HTTP error responses.

## Features

*   Simplified error handling in HTTP handlers.
*   Automatic conversion of `httperror.HttpError` to corresponding HTTP status codes and messages.
*   Fallback to HTTP 500 Internal Server Error for other error types.
*   Wrappers for:
    *   Standard `net/http` (`httpwrap`)
    *   `go-chi/chi/v5` (`chiwrap`)
    *   `gofiber/fiber/v2` (`fiberwrap`)
*   Customizable error callback for logging or other purposes (for `httpwrap` and `chiwrap`).

## Installation

To install `httpwrap` and its subpackages, use `go get`:

```bash
go get -u github.com/snowmerak/httpwrap
```

## httperror

The `httperror` package provides a convenient way to create Go errors that also carry specific HTTP status codes. This allows the `httpwrap` family of wrappers (`httpwrap`, `chiwrap`, `fiberwrap`) to automatically translate these errors into appropriate HTTP responses.

The package offers several functions to create these specialized errors:

*   `httperror.New(statusCode int, message string) *HttpError`: A generic function to create an error with any given HTTP `statusCode` and a custom `message`.
*   Specific helper functions for common HTTP status codes, such as:
    *   `httperror.BadRequest(message string) *HttpError` (for 400 Bad Request)
    *   `httperror.Unauthorized(message string) *HttpError` (for 401 Unauthorized)
    *   `httperror.Forbidden(message string) *HttpError` (for 403 Forbidden)
    *   `httperror.NotFound(message string) *HttpError` (for 404 Not Found)
    *   `httperror.InternalServerError(message string) *HttpError` (for 500 Internal Server Error)
    *   And many more, corresponding to standard HTTP status codes (e.g., `MethodNotAllowed`, `Conflict`, `ServiceUnavailable`, etc.).

When a handler in your application returns an error created by these functions (e.g., `httperror.New()` or `httperror.BadRequest()`), the respective wrapper will use the `statusCode` and `message` from this error to formulate the HTTP response. If a handler returns any other standard Go error, the wrappers will default to sending a 500 Internal Server Error.

## Usage

Below are examples of how to use each wrapper.

### 1. `httpwrap` (for standard `net/http`)

This wrapper is for the standard library's `http.ServeMux`.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/snowmerak/httpwrap/httperror"
	"github.com/snowmerak/httpwrap/wrapper/httpwrap"
)

func myHandler(w http.ResponseWriter, r *http.Request) error {
	// Simulate some logic
	shouldFail := r.URL.Query().Get("fail")
	if shouldFail == "true" {
		return httperror.BadRequest("Invalid request due to 'fail' parameter")
	}
	if shouldFail == "panic" {
		return fmt.Errorf("something unexpected happened")
	}
	_, err := w.Write([]byte("Hello from httpwrap!"))
	return err
}

func main() {
	mux := httpwrap.NewMux(func(err error) {
		// Optional: Custom error logging or handling
		fmt.Printf("An error occurred: %v\n", err)
	})

	mux.Handle("/hello", myHandler)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
```

### 2. `chiwrap` (for `go-chi/chi`)

This wrapper is for `go-chi/chi/v5`.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snowmerak/httpwrap/httperror"
	"github.com/snowmerak/httpwrap/wrapper/chiwrap"
)

func myChiHandler(w http.ResponseWriter, r *http.Request) error {
	param := chi.URLParam(r, "name")
	if param == "" {
		return httperror.BadRequest("Name parameter is missing")
	}
	if param == "error" {
		return fmt.Errorf("internal server error triggered")
	}
	_, err := w.Write([]byte(fmt.Sprintf("Hello, %s, from chiwrap!", param)))
	return err
}

func main() {
	router := chiwrap.NewRouter(func(err error) {
		// Optional: Custom error logging
		fmt.Printf("Chiwrap encountered an error: %v\n", err)
	})

	router.Get("/greet/{name}", myChiHandler)

	// You can also use other methods like Post, Put, Delete, etc.
	// router.Post("/submit", func(w http.ResponseWriter, r *http.Request) error { ... })

	// Mount the chiwrap router
	chiRouter := chi.NewRouter() // Or your existing chi router
	chiRouter.Mount("/chi", router) // Mount chiwrap.Router which implements http.Handler

	fmt.Println("Chi server starting on :8081...")
	if err := http.ListenAndServe(":8081", chiRouter); err != nil {
		fmt.Printf("Chi server failed to start: %v\n", err)
	}
}
```

### 3. `fiberwrap` (for `gofiber/fiber`)

This wrapper is for `gofiber/fiber/v2`.

```go
package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/snowmerak/httpwrap/httperror"
	"github.com/snowmerak/httpwrap/wrapper/fiberwrap"
)

func myFiberHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "0" {
		return httperror.NotFound("User with ID 0 not found")
	}
	if userID == "panic" {
		return fmt.Errorf("a critical error occurred")
	}
	return c.SendString(fmt.Sprintf("Hello from Fiber user %s!", userID))
}

func main() {
	// Using a new Fiber app
	fw := fiberwrap.NewWrapper()
	app := fw.App()

	// Or using an existing Fiber app
	// existingApp := fiber.New()
	// fw := fiberwrap.WithApp(existingApp)
	// app := fw.App()

	fw.Get("/fiber/user/:id", myFiberHandler)
	// fw.Post("/fiber/data", func(c *fiber.Ctx) error { ... })

	fmt.Println("Fiber server starting on :8082...")
	if err := app.Listen(":8082"); err != nil {
		fmt.Printf("Fiber server failed to start: %v\n", err)
	}
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

