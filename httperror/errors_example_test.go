package httperror_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gosuda/httpwrap/httperror"
	"github.com/gosuda/httpwrap/wrapper/httpwrap"
)

// ExampleHttpError_contentType demonstrates how to use HttpError with custom Content-Type headers.
// This example shows various content types including JSON, RFC9457 Problem Details, HTML, XML, and plain text.
func ExampleHttpError_contentType() {
	// Create a new mux with error callback
	mux := httpwrap.NewMux(func(err error) {
		log.Printf("Error occurred: %v", err)
	})

	// Handler that returns JSON error with Content-Type
	mux.Handle("/json-error", func(w http.ResponseWriter, r *http.Request) error {
		jsonError := `{"error": "Invalid request", "code": "BAD_REQUEST", "timestamp": "2025-06-03T10:00:00Z"}`
		return httperror.New(400, jsonError, "application/json")
	})

	// Handler that returns RFC9457 Problem Details
	mux.Handle("/problem-details", func(w http.ResponseWriter, r *http.Request) error {
		problem := httperror.BadRequestProblem9457("Validation failed for user input")
		problem.WithType("https://example.com/errors/validation-failed")
		problem.WithInstance("/api/users/123")
		problem.WithExtension("fields", []string{"email", "username"})
		problem.WithTraceID("trace-abc123def456")
		
		return problem.ToHttpError() // This automatically sets Content-Type to application/problem+json
	})

	// Handler that returns HTML error with Content-Type
	mux.Handle("/html-error", func(w http.ResponseWriter, r *http.Request) error {
		htmlError := `<!DOCTYPE html>
<html>
<head><title>Error</title></head>
<body>
	<h1>404 - Page Not Found</h1>
	<p>The requested resource could not be found.</p>
</body>
</html>`
		return httperror.New(404, htmlError, "text/html")
	})

	// Handler that returns XML error with Content-Type
	mux.Handle("/xml-error", func(w http.ResponseWriter, r *http.Request) error {
		xmlError := `<?xml version="1.0" encoding="UTF-8"?>
<error>
	<code>500</code>
	<message>Internal Server Error</message>
	<details>Database connection failed</details>
</error>`
		return httperror.New(500, xmlError, "application/xml")
	})

	// Handler that returns plain text error (default behavior)
	mux.Handle("/text-error", func(w http.ResponseWriter, r *http.Request) error {
		return httperror.New(403, "Access denied")
	})

	// Success handler for comparison
	mux.Handle("/success", func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"message": "Success!", "status": "ok"}`))
		return nil
	})

	fmt.Println("Content-Type demo server configured with endpoints:")
	fmt.Println("  /json-error       - JSON error with application/json")
	fmt.Println("  /problem-details  - RFC9457 with application/problem+json")
	fmt.Println("  /html-error       - HTML error with text/html")
	fmt.Println("  /xml-error        - XML error with application/xml")
	fmt.Println("  /text-error       - Plain text error (default)")
	fmt.Println("  /success          - Success response for comparison")

	// Output:
	// Content-Type demo server configured with endpoints:
	//   /json-error       - JSON error with application/json
	//   /problem-details  - RFC9457 with application/problem+json
	//   /html-error       - HTML error with text/html
	//   /xml-error        - XML error with application/xml
	//   /text-error       - Plain text error (default)
	//   /success          - Success response for comparison
}

// ExampleNew_withContentType demonstrates creating HttpError with custom Content-Type.
func ExampleNew_withContentType() {
	// Create JSON error with application/json content type
	jsonErr := httperror.New(400, `{"error":"bad request"}`, "application/json")
	fmt.Printf("Status: %d, Message: %s, Content-Type: %s\n", 
		jsonErr.StatusCode(), jsonErr.ErrorMessage(), jsonErr.ContentType)

	// Create XML error with application/xml content type
	xmlErr := httperror.New(500, `<error><message>Server Error</message></error>`, "application/xml")
	fmt.Printf("Status: %d, Content-Type: %s\n", 
		xmlErr.StatusCode(), xmlErr.ContentType)

	// Create error without content type (uses default)
	textErr := httperror.New(404, "Not found")
	fmt.Printf("Status: %d, Content-Type: %s\n", 
		textErr.StatusCode(), textErr.ContentType)

	// Output:
	// Status: 400, Message: {"error":"bad request"}, Content-Type: application/json
	// Status: 500, Content-Type: application/xml
	// Status: 404, Content-Type: 
}

// ExampleRFC9457Error_ToHttpError demonstrates RFC9457 Problem Details with automatic Content-Type.
func ExampleRFC9457Error_ToHttpError() {
	// Create RFC9457 Problem Details error
	problem := httperror.BadRequestProblem9457("Invalid input data")
	problem.WithType("https://example.com/errors/validation")
	problem.WithInstance("/api/users/create")
	problem.WithExtension("field", "email")
	problem.WithTraceID("trace-123")

	// Convert to HttpError (automatically sets Content-Type to application/problem+json)
	httpErr := problem.ToHttpError()
	
	fmt.Printf("Status: %d\n", httpErr.StatusCode())
	fmt.Printf("Content-Type: %s\n", httpErr.ContentType)
	fmt.Printf("Contains JSON: %t\n", httpErr.ContentType == "application/problem+json")

	// Output:
	// Status: 400
	// Content-Type: application/problem+json
	// Contains JSON: true
}

// ExampleHttpError_differentContentTypes shows various content type examples.
func ExampleHttpError_differentContentTypes() {
	// JSON API error
	jsonErr := httperror.New(422, 
		`{"error":"validation_failed","fields":["email","password"]}`, 
		"application/json")
	
	// HTML error page
	htmlErr := httperror.New(404, 
		`<html><body><h1>Page Not Found</h1></body></html>`, 
		"text/html")
	
	// XML SOAP fault
	xmlErr := httperror.New(500, 
		`<?xml version="1.0"?><soap:Fault><faultstring>Server Error</faultstring></soap:Fault>`, 
		"application/soap+xml")
	
	// CSV error report
	csvErr := httperror.New(400, 
		"line,error\n1,missing required field\n2,invalid date format", 
		"text/csv")

	// Plain text (default - no content type specified)
	textErr := httperror.New(403, "Access denied")

	fmt.Printf("JSON Error - Type: %s, Status: %d\n", jsonErr.ContentType, jsonErr.StatusCode())
	fmt.Printf("HTML Error - Type: %s, Status: %d\n", htmlErr.ContentType, htmlErr.StatusCode())
	fmt.Printf("XML Error - Type: %s, Status: %d\n", xmlErr.ContentType, xmlErr.StatusCode())
	fmt.Printf("CSV Error - Type: %s, Status: %d\n", csvErr.ContentType, csvErr.StatusCode())
	fmt.Printf("Text Error - Type: '%s', Status: %d\n", textErr.ContentType, textErr.StatusCode())

	// Output:
	// JSON Error - Type: application/json, Status: 422
	// HTML Error - Type: text/html, Status: 404
	// XML Error - Type: application/soap+xml, Status: 500
	// CSV Error - Type: text/csv, Status: 400
	// Text Error - Type: '', Status: 403
}
