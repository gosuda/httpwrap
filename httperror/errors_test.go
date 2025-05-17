package httperror

import (
	"testing"
)

func TestHttpError_Error(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
		want    string
	}{
		{
			name:    "400 Bad Request",
			code:    400,
			message: "Invalid input",
			want:    "400: Invalid input",
		},
		{
			name:    "404 Not Found",
			code:    404,
			message: "Resource not found",
			want:    "404: Resource not found",
		},
		{
			name:    "500 Internal Server Error",
			code:    500,
			message: "Internal server error",
			want:    "500: Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.code, tt.message)
			if got := e.Error(); got != tt.want {
				t.Errorf("HttpError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpError_StatusCode(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "400 Bad Request",
			code: 400,
		},
		{
			name: "404 Not Found",
			code: 404,
		},
		{
			name: "500 Internal Server Error",
			code: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.code, "Test message")
			if got := e.StatusCode(); got != tt.code {
				t.Errorf("HttpError.StatusCode() = %v, want %v", got, tt.code)
			}
		})
	}
}

func TestHttpError_ErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Simple message",
			message: "Simple error message",
		},
		{
			name:    "Empty message",
			message: "",
		},
		{
			name:    "Message with special characters",
			message: "Error: $#@! special chars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(400, tt.message)
			if got := e.ErrorMessage(); got != tt.message {
				t.Errorf("HttpError.ErrorMessage() = %v, want %v", got, tt.message)
			}
		})
	}
}

func TestErrorHelperFunctions(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(string) *HttpError
		code    int
		message string
	}{
		{
			name:    "BadRequest",
			fn:      BadRequest,
			code:    400,
			message: "Bad request error",
		},
		{
			name:    "Unauthorized",
			fn:      Unauthorized,
			code:    401,
			message: "Unauthorized error",
		},
		{
			name:    "Forbidden",
			fn:      Forbidden,
			code:    403,
			message: "Forbidden error",
		},
		{
			name:    "NotFound",
			fn:      NotFound,
			code:    404,
			message: "Not found error",
		},
		{
			name:    "MethodNotAllowed",
			fn:      MethodNotAllowed,
			code:    405,
			message: "Method not allowed error",
		},
		{
			name:    "InternalServerError",
			fn:      InternalServerError,
			code:    500,
			message: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.fn(tt.message)
			if e.Code != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, e.Code)
			}
			if e.Message != tt.message {
				t.Errorf("Expected message %s, got %s", tt.message, e.Message)
			}
		})
	}
}
