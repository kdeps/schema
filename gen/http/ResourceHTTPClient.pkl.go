// Code generated from Pkl module `org.kdeps.pkl.Http`. DO NOT EDIT.
package http

// Class representing an HTTP client resource, which includes details
// about the HTTP method, URL, request data, headers, and response.
type ResourceHTTPClient struct {
	// The HTTP method to be used for the request.
	Method string `pkl:"method"`

	// The URL to which the request will be sent.
	Url string `pkl:"url"`

	// Optional data to be sent with the request.
	Data *[]string `pkl:"data"`

	// A mapping of headers to be included in the request.
	Headers *map[string]string `pkl:"headers"`

	// The response received from the HTTP request.
	Response *ResponseBlock `pkl:"response"`

	// A timestamp of when the request was made, represented as an unsigned 32-bit integer.
	Timestamp *uint32 `pkl:"timestamp"`

	// The timeout duration (in seconds) for the HTTP request. Defaults to 60 seconds.
	TimeoutSeconds *int `pkl:"timeoutSeconds"`
}
