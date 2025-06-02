package httperror

import (
	"strconv"
)

// HttpError represents an HTTP error with a status code and message.
// It implements the error interface and provides methods to retrieve the status code and error message.
// The error message is typically a human-readable description of the error.
// The status code is an integer representing the HTTP status code associated with the error.
// The error message and status code can be used to provide more context about the error
// and to help diagnose issues in the application.
type HttpError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	ContentType string `json:"content_type,omitempty"` // Optional content type for the error response
}

// New creates a new HttpError with the given status code and message.
// The status code is an integer representing the HTTP status code associated with the error.
// The message is a string representing the error message.
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

func (e *HttpError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}

func (e *HttpError) StatusCode() int {
	return e.Code
}

func (e *HttpError) ErrorMessage() string {
	return e.Message
}

// BadRequest indicates that the server cannot or will not process the request due to something that is perceived to be a client error.
// This code is typically used in situations where the client has sent a malformed request or has provided invalid data.
// For example, if a client tries to send a request with an invalid JSON payload, the server may respond with a 400 Bad Request status code.
func BadRequest(message string) *HttpError {
	return New(400, message)
}

// Unauthorized indicates that the request requires user authentication.
// The 401 status code is typically used in situations where the client must authenticate itself to get the requested response.
// For example, if a client tries to access a resource that requires authentication and does not provide valid credentials,
func Unauthorized(message string) *HttpError {
	return New(401, message)
}

// PaymentRequired indicates that the request requires payment.
// This code is typically used in situations where the client must pay for a service or resource before it can be accessed.
// For example, if a client tries to access a resource that requires a subscription or payment, the server may respond with a 402 Payment Required status code.
func PaymentRequired(message string) *HttpError {
	return New(402, message)
}

// Forbidden indicates that the server understood the request but refuses to authorize it.
// This code is typically used in situations where the client does not have permission to access the requested resource.
// For example, if a client tries to access a resource that is restricted to certain users or roles, the server may respond with a 403 Forbidden status code.
func Forbidden(message string) *HttpError {
	return New(403, message)
}

// NotFound indicates that the server cannot find the requested resource.
// This code is typically used in situations where the requested resource does not exist or is not available.
// For example, if a client tries to access a resource that has been deleted or moved, the server may respond with a 404 Not Found status code.
func NotFound(message string) *HttpError {
	return New(404, message)
}

// MethodNotAllowed indicates that the request method is not supported by the server for the requested resource.
// This code is typically used in situations where the client has made a request using an HTTP method that is not supported by the server.
// For example, if a client tries to use the PUT method to update a resource that only supports the GET method, the server may respond with a 405 Method Not Allowed status code.
func MethodNotAllowed(message string) *HttpError {
	return New(405, message)
}

// NotAcceptable indicates that the resource is not capable of generating content that is acceptable to the client.
// The server should generate a 406 Not Acceptable response if the client has indicated that it only accepts certain content types
// and the server cannot generate a response that meets those criteria.
func NotAcceptable(message string) *HttpError {
	return New(406, message)
}

// ProxyAuthRequired indicates that the client must first authenticate itself with the proxy.
// The 407 status code is similar to 401 Unauthorized, but indicates that the client must first authenticate itself with the proxy.
// This code is typically used in situations where the client must authenticate with a proxy server before accessing the target resource.
func ProxyAuthRequired(message string) *HttpError {
	return New(407, message)
}

// RequestTimeout indicates that the server did not receive a complete request message within the time that it was prepared to wait.
// This code is typically used in situations where the server is unable to process the request due to a timeout.
// For example, if a client takes too long to send the request body, the server may respond with a 408 Request Timeout status code.
func RequestTimeout(message string) *HttpError {
	return New(408, message)
}

// Conflict indicates that the request could not be completed due to a conflict with the current state of the target resource.
// This code is typically used in situations where the request could not be completed due to a conflict with the current state of the target resource.
// For example, if a user tries to create a resource that already exists, the server may respond with a 409 Conflict status code.
func Conflict(message string) *HttpError {
	return New(409, message)
}

// Gone indicates that the resource is no longer available at the server and no forwarding address is known.
// The server does not wish to make this information available to the client by including a payload
// and the client should not expect to be able to reuse the same resource in the future.
func Gone(message string) *HttpError {
	return New(410, message)
}

// LengthRequired indicates that the server refuses to accept the request without a defined Content-Length.
// The server should respond with a 411 Length Required status code if the client has not provided a Content-Length header
// and the server requires it to process the request.
func LengthRequired(message string) *HttpError {
	return New(411, message)
}

// PreconditionFailed indicates that the server does not meet one of the preconditions that the requester put on the request header fields.
// This code is typically used in situations where the client has specified certain conditions that must be met for the request to be processed.
// For example, if a client specifies an If-Match header and the resource has been modified since the request was made, the server may respond with a 412 Precondition Failed status code.
func PreconditionFailed(message string) *HttpError {
	return New(412, message)
}

// PayloadTooLarge indicates that the request is larger than the server is willing or able to process.
// This code is typically used in situations where the client has sent a request that is too large for the server to handle.
// For example, if a client tries to upload a file that is too large, the server may respond with a 413 Payload Too Large status code.
func PayloadTooLarge(message string) *HttpError {
	return New(413, message)
}

// URITooLong indicates that the request URI is longer than the server is willing to interpret.
// This code is typically used in situations where the client has sent a request with a URI that is too long for the server to handle.
// For example, if a client tries to send a request with a very long query string, the server may respond with a 414 URI Too Long status code.
func URITooLong(message string) *HttpError {
	return New(414, message)
}

// UnsupportedMediaType indicates that the server refuses to accept the request because the payload format is in an unsupported format.
// This code is typically used in situations where the client has sent a request with a payload format that the server does not support.
// For example, if a client tries to send a request with a JSON payload but the server only supports XML, the server may respond with a 415 Unsupported Media Type status code.
func UnsupportedMediaType(message string) *HttpError {
	return New(415, message)
}

// RangeNotSatisfiable indicates that the server cannot provide the requested range for the resource.
// This code is typically used in situations where the client has requested a range of bytes from a resource
// and the server cannot provide that range.
// For example, if a client tries to request a range of bytes from a file that is smaller than the requested range,
// the server may respond with a 416 Range Not Satisfiable status code.
func RangeNotSatisfiable(message string) *HttpError {
	return New(416, message)
}

// ExpectationFailed indicates that the server cannot meet the requirements of the Expect request-header field.
// This code is typically used in situations where the client has specified certain expectations in the request header
// and the server cannot meet those expectations.
// For example, if a client specifies an Expect header with a value of 100-continue and the server cannot process the request,
// the server may respond with a 417 Expectation Failed status code.
func ExpectationFailed(message string) *HttpError {
	return New(417, message)
}

// MisdirectedRequest indicates that the request was directed at a server that is not able to produce a response.
// This code is typically used in situations where the client has sent a request to a server that is not able to process it.
// For example, if a client tries to send a request to a server that is not configured to handle that type of request,
// the server may respond with a 421 Misdirected Request status code.
func MisdirectedRequest(message string) *HttpError {
	return New(421, message)
}

// UnprocessableEntity indicates that the server understands the content type of the request entity,
// and the syntax of the request entity is correct, but it was unable to process the contained instructions.
// This code is typically used in situations where the client has sent a request with a valid syntax
// but the server cannot process the request due to semantic errors.
// For example, if a client tries to send a request with a valid JSON payload but the payload contains invalid data,
// the server may respond with a 422 Unprocessable Entity status code.
func UnprocessableEntity(message string) *HttpError {
	return New(422, message)
}

// Locked indicates that the resource that is being accessed is locked.
// This code is typically used in situations where the client has tried to access a resource that is locked
// and cannot be modified until the lock is released.
// For example, if a client tries to update a resource that is currently being edited by another user,
// the server may respond with a 423 Locked status code.
func Locked(message string) *HttpError {
	return New(423, message)
}

// FailedDependency indicates that the request failed due to failure of a previous request.
// This code is typically used in situations where the client has sent a request that depends on the success of a previous request
// and the previous request has failed.
// For example, if a client tries to send a request to update a resource that depends on the success of a previous request
// and the previous request has failed, the server may respond with a 424 Failed Dependency status code.
func FailedDependency(message string) *HttpError {
	return New(424, message)
}

// UpgradeRequired indicates that the client should switch to a different protocol.
// This code is typically used in situations where the server requires the client to switch to a different protocol
// in order to process the request.
// For example, if a client tries to send a request using HTTP/1.1 but the server requires HTTP/2,
// the server may respond with a 426 Upgrade Required status code.
func UpgradeRequired(message string) *HttpError {
	return New(426, message)
}

// PreconditionRequired indicates that the server requires the request to be conditional.
// This code is typically used in situations where the server requires the client to provide certain conditions
// in order to process the request.
// For example, if a client tries to send a request without providing an If-Match header
// and the server requires it to process the request, the server may respond with a 428 Precondition Required status code.
func PreconditionRequired(message string) *HttpError {
	return New(428, message)
}

// TooManyRequests indicates that the user has sent too many requests in a given amount of time.
// This code is typically used in situations where the client has exceeded the rate limit for requests
// and the server is unable to process the request.
// For example, if a client tries to send too many requests in a short period of time,
// the server may respond with a 429 Too Many Requests status code.
func TooManyRequests(message string) *HttpError {
	return New(429, message)
}

// RequestHeaderFieldsTooLarge indicates that the server is unwilling to process the request
// because its header fields are too large.
// This code is typically used in situations where the client has sent a request with header fields that are too large for the server to handle.
// For example, if a client tries to send a request with a very large cookie header,
// the server may respond with a 431 Request Header Fields Too Large status code.
func RequestHeaderFieldsTooLarge(message string) *HttpError {
	return New(431, message)
}

// UnavailableForLegalReasons indicates that the server is denying access to the resource
// in response to a legal demand.
// This code is typically used in situations where the server is required to deny access to a resource
// due to legal reasons, such as a court order or government request.
func UnavailableForLegalReasons(message string) *HttpError {
	return New(451, message)
}

// InternalServerError indicates that the server encountered an unexpected condition
// that prevented it from fulfilling the request.
// This code is typically used in situations where the server has encountered an error
// that is not related to the client's request.
// For example, if a server encounters a database error while processing a request,
// the server may respond with a 500 Internal Server Error status code.
func InternalServerError(message string) *HttpError {
	return New(500, message)
}

// NotImplemented indicates that the server does not support the functionality required to fulfill the request.
// This code is typically used in situations where the server does not support the requested method
// or the requested resource.
// For example, if a client tries to send a request using a method that is not supported by the server,
// the server may respond with a 501 Not Implemented status code.
func NotImplemented(message string) *HttpError {
	return New(501, message)
}

// BadGateway indicates that the server, while acting as a gateway or proxy,
// received an invalid response from the upstream server.
// This code is typically used in situations where the server is acting as a gateway or proxy
// and receives an invalid response from the upstream server.
// For example, if a server is acting as a reverse proxy and receives an invalid response from the backend server,
// the server may respond with a 502 Bad Gateway status code.
func BadGateway(message string) *HttpError {
	return New(502, message)
}

// ServiceUnavailable indicates that the server is currently unable to handle the request
// due to temporary overloading or maintenance of the server.
// This code is typically used in situations where the server is temporarily unable to handle the request
// due to high load or maintenance.
// For example, if a server is undergoing maintenance and cannot process requests,
// the server may respond with a 503 Service Unavailable status code.
func ServiceUnavailable(message string) *HttpError {
	return New(503, message)
}

// GatewayTimeout indicates that the server, while acting as a gateway or proxy,
// did not receive a timely response from the upstream server.
// This code is typically used in situations where the server is acting as a gateway or proxy
// and does not receive a timely response from the upstream server.
// For example, if a server is acting as a reverse proxy and does not receive a response from the backend server within a certain time limit,
// the server may respond with a 504 Gateway Timeout status code.
func GatewayTimeout(message string) *HttpError {
	return New(504, message)
}

// HTTPVersionNotSupported indicates that the server does not support the HTTP protocol version
// that was used in the request message.
// This code is typically used in situations where the client has sent a request using an HTTP version
// that is not supported by the server.
// For example, if a client tries to send a request using HTTP/2 but the server only supports HTTP/1.1,
// the server may respond with a 505 HTTP Version Not Supported status code.
func HTTPVersionNotSupported(message string) *HttpError {
	return New(505, message)
}

// VariantAlsoNegotiates indicates that the server has an internal configuration error
// and the selected variant resource is configured to engage in transparent content negotiation
// itself and is therefore not a proper end point in the negotiation process.
// This code is typically used in situations where the server has an internal configuration error
// and the selected variant resource is not a proper end point in the negotiation process.
// For example, if a server is configured to use transparent content negotiation
// and the selected variant resource is not a proper end point in the negotiation process,
// the server may respond with a 506 Variant Also Negotiates status code.
func VariantAlsoNegotiates(message string) *HttpError {
	return New(506, message)
}

// InsufficientStorage indicates that the server is unable to store the representation needed to complete the request.
// This code is typically used in situations where the server is unable to store the representation needed to complete the request.
// For example, if a server is unable to store the representation needed to complete the request due to insufficient storage,
// the server may respond with a 507 Insufficient Storage status code.
func InsufficientStorage(message string) *HttpError {
	return New(507, message)
}

// LoopDetected indicates that the server detected an infinite loop while processing the request.
// This code is typically used in situations where the server detects an infinite loop while processing the request.
// For example, if a server is processing a request that causes an infinite loop,
// the server may respond with a 508 Loop Detected status code.
func LoopDetected(message string) *HttpError {
	return New(508, message)
}

// NotExtended indicates that further extensions to the request are required for the server to fulfill it.
// This code is typically used in situations where the server requires further extensions to the request
// in order to fulfill it.
// For example, if a server requires further extensions to the request in order to fulfill it,
// the server may respond with a 510 Not Extended status code.
func NotExtended(message string) *HttpError {
	return New(510, message)
}

// NetworkAuthenticationRequired indicates that the client needs to authenticate to gain network access.
// This code is typically used in situations where the client needs to authenticate to gain network access.
// For example, if a client tries to access a resource that requires authentication
// and the client is not authenticated, the server may respond with a 511 Network Authentication Required status code.
func NetworkAuthenticationRequired(message string) *HttpError {
	return New(511, message)
}
