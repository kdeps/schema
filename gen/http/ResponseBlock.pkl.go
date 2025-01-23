// Code generated from Pkl module `org.kdeps.pkl.HTTP`. DO NOT EDIT.
package http

// Class representing the response block of an HTTP request.
// It contains the body and headers of the response.
type ResponseBlock struct {
	// The body of the response.
	Body *string `pkl:"body"`

	// A mapping of response headers.
	Headers *map[string]string `pkl:"headers"`
}
