// Code generated from Pkl module `org.kdeps.pkl.HTTP`. DO NOT EDIT.
package http

import "github.com/apple/pkl-go/pkl"

// Class representing an HTTP client resource, which includes details
// about the HTTP method, URL, request data, headers, and response.
type ResourceHTTPClient struct {
	// The HTTP method to be used for the request.
	Method string `pkl:"Method"`

	// The URL to which the request will be sent.
	Url string `pkl:"Url"`

	// Optional data to be sent with the request.
	Data *[]string `pkl:"Data"`

	// A mapping of headers to be included in the request.
	Headers *map[string]string `pkl:"Headers"`

	// A mapping of parameters to be included in the request.
	Params *map[string]string `pkl:"Params"`

	// The response received from the HTTP request.
	Response *ResponseBlock `pkl:"Response"`

	// The file path where the response body value of this resource is saved
	File *string `pkl:"File"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"ItemValues"`

	// A timestamp of when the request was made, represented as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"Timestamp"`

	// The timeout duration (in seconds) for the HTTP request. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`
}
