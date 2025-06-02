// Package httperror provides HTTP error handling utilities with support for both RFC7807 and RFC9457 problem details formats.
// It offers a standard HttpError type for representing HTTP errors with status codes, messages, and optional content types,
// along with utilities for creating and handling structured error responses.
package httperror

import (
	"strconv"
)

// HttpError represents an HTTP error with a status code, message, and optional content type.
// It implements the error interface and provides methods to retrieve the status code and error message.
// The ContentType field allows customization of the response content type for different error formats.
type HttpError struct {
	Code        int    `json:"code"`                   // HTTP status code
	Message     string `json:"message"`                // Human-readable error message
	ContentType string `json:"content_type,omitempty"` // Optional content type for the error response
}

// New creates a new HttpError with the specified status code, message, and optional content type.
// The contentType parameter is optional; if provided, it sets the ContentType field for custom response formatting.
func New(code int, message string, contentType ...string) *HttpError {
	ct := ""
	if len(contentType) > 0 {
		ct = contentType[0]
	}
	return &HttpError{
		Code:        code,
		Message:     message,
		ContentType: ct,
	}
}

// Error returns a string representation of the HttpError in the format "code: message".
// This method implements the error interface.
func (e *HttpError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}

// StatusCode returns the HTTP status code associated with this error.
func (e *HttpError) StatusCode() int {
	return e.Code
}

// ErrorMessage returns the human-readable error message associated with this HttpError.
func (e *HttpError) ErrorMessage() string {
	return e.Message
}

// BadRequest creates a new HttpError with status code 400 (Bad Request).
// This indicates that the server cannot or will not process the request due to a client error.
func BadRequest(message string) *HttpError {
	return New(400, message)
}

// Unauthorized creates a new HttpError with status code 401 (Unauthorized).
// This indicates that the request requires user authentication.
func Unauthorized(message string) *HttpError {
	return New(401, message)
}

// PaymentRequired creates a new HttpError with status code 402 (Payment Required).
// This indicates that the request requires payment.
func PaymentRequired(message string) *HttpError {
	return New(402, message)
}

// Forbidden creates a new HttpError with status code 403 (Forbidden).
// This indicates that the server understood the request but refuses to authorize it.
func Forbidden(message string) *HttpError {
	return New(403, message)
}

// NotFound creates a new HttpError with status code 404 (Not Found).
// This indicates that the server cannot find the requested resource.
func NotFound(message string) *HttpError {
	return New(404, message)
}

// MethodNotAllowed creates a new HttpError with status code 405 (Method Not Allowed).
// This indicates that the request method is not supported by the server for the requested resource.
func MethodNotAllowed(message string) *HttpError {
	return New(405, message)
}

// NotAcceptable creates a new HttpError with status code 406 (Not Acceptable).
// This indicates that the resource is not capable of generating content acceptable to the client.
func NotAcceptable(message string) *HttpError {
	return New(406, message)
}

// ProxyAuthRequired creates a new HttpError with status code 407 (Proxy Authentication Required).
// This indicates that the client must first authenticate itself with the proxy.
func ProxyAuthRequired(message string) *HttpError {
	return New(407, message)
}

// RequestTimeout creates a new HttpError with status code 408 (Request Timeout).
// This indicates that the server did not receive a complete request message within the expected time.
func RequestTimeout(message string) *HttpError {
	return New(408, message)
}

// Conflict creates a new HttpError with status code 409 (Conflict).
// This indicates that the request could not be completed due to a conflict with the current state of the target resource.
func Conflict(message string) *HttpError {
	return New(409, message)
}

// Gone creates a new HttpError with status code 410 (Gone).
// This indicates that the resource is no longer available at the server and no forwarding address is known.
func Gone(message string) *HttpError {
	return New(410, message)
}

// LengthRequired creates a new HttpError with status code 411 (Length Required).
// This indicates that the server refuses to accept the request without a defined Content-Length.
func LengthRequired(message string) *HttpError {
	return New(411, message)
}

// PreconditionFailed creates a new HttpError with status code 412 (Precondition Failed).
// This indicates that the server does not meet one of the preconditions in the request header fields.
func PreconditionFailed(message string) *HttpError {
	return New(412, message)
}

// PayloadTooLarge creates a new HttpError with status code 413 (Payload Too Large).
// This indicates that the request is larger than the server is willing or able to process.
func PayloadTooLarge(message string) *HttpError {
	return New(413, message)
}

// URITooLong creates a new HttpError with status code 414 (URI Too Long).
// This indicates that the request URI is longer than the server is willing to interpret.
func URITooLong(message string) *HttpError {
	return New(414, message)
}

// UnsupportedMediaType creates a new HttpError with status code 415 (Unsupported Media Type).
// This indicates that the server refuses to accept the request because the payload format is unsupported.
func UnsupportedMediaType(message string) *HttpError {
	return New(415, message)
}

// RangeNotSatisfiable creates a new HttpError with status code 416 (Range Not Satisfiable).
// This indicates that the server cannot provide the requested range for the resource.
func RangeNotSatisfiable(message string) *HttpError {
	return New(416, message)
}

// ExpectationFailed creates a new HttpError with status code 417 (Expectation Failed).
// This indicates that the server cannot meet the requirements of the Expect request-header field.
func ExpectationFailed(message string) *HttpError {
	return New(417, message)
}

// MisdirectedRequest creates a new HttpError with status code 421 (Misdirected Request).
// This indicates that the request was directed at a server that is not able to produce a response.
func MisdirectedRequest(message string) *HttpError {
	return New(421, message)
}

// UnprocessableEntity creates a new HttpError with status code 422 (Unprocessable Entity).
// This indicates that the server understands the content type but was unable to process the contained instructions.
func UnprocessableEntity(message string) *HttpError {
	return New(422, message)
}

// Locked creates a new HttpError with status code 423 (Locked).
// This indicates that the resource that is being accessed is locked.
func Locked(message string) *HttpError {
	return New(423, message)
}

// FailedDependency creates a new HttpError with status code 424 (Failed Dependency).
// This indicates that the request failed due to failure of a previous request.
func FailedDependency(message string) *HttpError {
	return New(424, message)
}

// UpgradeRequired creates a new HttpError with status code 426 (Upgrade Required).
// This indicates that the client should switch to a different protocol.
func UpgradeRequired(message string) *HttpError {
	return New(426, message)
}

// PreconditionRequired creates a new HttpError with status code 428 (Precondition Required).
// This indicates that the server requires the request to be conditional.
func PreconditionRequired(message string) *HttpError {
	return New(428, message)
}

// TooManyRequests creates a new HttpError with status code 429 (Too Many Requests).
// This indicates that the user has sent too many requests in a given amount of time.
func TooManyRequests(message string) *HttpError {
	return New(429, message)
}

// RequestHeaderFieldsTooLarge creates a new HttpError with status code 431 (Request Header Fields Too Large).
// This indicates that the server is unwilling to process the request because its header fields are too large.
func RequestHeaderFieldsTooLarge(message string) *HttpError {
	return New(431, message)
}

// UnavailableForLegalReasons creates a new HttpError with status code 451 (Unavailable For Legal Reasons).
// This indicates that the server is denying access to the resource in response to a legal demand.
func UnavailableForLegalReasons(message string) *HttpError {
	return New(451, message)
}

// InternalServerError creates a new HttpError with status code 500 (Internal Server Error).
// This indicates that the server encountered an unexpected condition that prevented it from fulfilling the request.
func InternalServerError(message string) *HttpError {
	return New(500, message)
}

// NotImplemented creates a new HttpError with status code 501 (Not Implemented).
// This indicates that the server does not support the functionality required to fulfill the request.
func NotImplemented(message string) *HttpError {
	return New(501, message)
}

// BadGateway creates a new HttpError with status code 502 (Bad Gateway).
// This indicates that the server, while acting as a gateway or proxy, received an invalid response from the upstream server.
func BadGateway(message string) *HttpError {
	return New(502, message)
}

// ServiceUnavailable creates a new HttpError with status code 503 (Service Unavailable).
// This indicates that the server is currently unable to handle the request due to temporary overloading or maintenance.
func ServiceUnavailable(message string) *HttpError {
	return New(503, message)
}

// GatewayTimeout creates a new HttpError with status code 504 (Gateway Timeout).
// This indicates that the server, while acting as a gateway or proxy, did not receive a timely response from the upstream server.
func GatewayTimeout(message string) *HttpError {
	return New(504, message)
}

// HTTPVersionNotSupported creates a new HttpError with status code 505 (HTTP Version Not Supported).
// This indicates that the server does not support the HTTP protocol version used in the request.
func HTTPVersionNotSupported(message string) *HttpError {
	return New(505, message)
}

// VariantAlsoNegotiates creates a new HttpError with status code 506 (Variant Also Negotiates).
// This indicates that the server has an internal configuration error related to transparent content negotiation.
func VariantAlsoNegotiates(message string) *HttpError {
	return New(506, message)
}

// InsufficientStorage creates a new HttpError with status code 507 (Insufficient Storage).
// This indicates that the server is unable to store the representation needed to complete the request.
func InsufficientStorage(message string) *HttpError {
	return New(507, message)
}

// LoopDetected creates a new HttpError with status code 508 (Loop Detected).
// This indicates that the server detected an infinite loop while processing the request.
func LoopDetected(message string) *HttpError {
	return New(508, message)
}

// NotExtended creates a new HttpError with status code 510 (Not Extended).
// This indicates that further extensions to the request are required for the server to fulfill it.
func NotExtended(message string) *HttpError {
	return New(510, message)
}

// NetworkAuthenticationRequired creates a new HttpError with status code 511 (Network Authentication Required).
// This indicates that the client needs to authenticate to gain network access.
func NetworkAuthenticationRequired(message string) *HttpError {
	return New(511, message)
}
