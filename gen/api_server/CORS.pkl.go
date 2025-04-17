// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

import "github.com/apple/pkl-go/pkl"

// Cross-Origin Resource Sharing (CORS) configuration
type CORS struct {
	// Enables or disables CORS support (default: false)
	EnableCORS bool `pkl:"enableCORS"`

	// List of allowed origin domains for CORS requests (e.g., "https://example.com")
	//
	// If unset, no origins are allowed unless CORS is disabled
	AllowOrigins *[]string `pkl:"allowOrigins"`

	// List of HTTP methods allowed for CORS requests, validated by regex
	//
	// If unset, defaults to methods specified in the route configuration
	AllowMethods *[]string `pkl:"allowMethods"`

	// List of request headers allowed in CORS requests (e.g., "Content-Type")
	//
	// If unset, no additional headers are allowed
	AllowHeaders *[]string `pkl:"allowHeaders"`

	// List of response headers exposed to clients in CORS requests
	//
	// If unset, no headers are exposed beyond defaults
	ExposeHeaders *[]string `pkl:"exposeHeaders"`

	// Allows credentials (e.g., cookies, HTTP authentication) in CORS requests (default: true)
	AllowCredentials bool `pkl:"allowCredentials"`

	// Maximum duration (in hours) for which CORS preflight responses can be cached (default: 12 hours)
	MaxAge *pkl.Duration `pkl:"maxAge"`
}
