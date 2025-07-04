// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

import "github.com/apple/pkl-go/pkl"

// Cross-Origin Resource Sharing (CORS) configuration
type CORS struct {
	// Enable Cross-Origin Resource Sharing (CORS) for the API server
	EnableCORS *bool `pkl:"EnableCORS"`

	// List of allowed origins for CORS
	AllowOrigins *[]string `pkl:"AllowOrigins"`

	// List of allowed HTTP methods for CORS
	AllowMethods *[]string `pkl:"AllowMethods"`

	// List of allowed headers for CORS
	AllowHeaders *[]string `pkl:"AllowHeaders"`

	// List of exposed headers for CORS
	ExposeHeaders *[]string `pkl:"ExposeHeaders"`

	// Maximum age for CORS preflight requests (in seconds)
	MaxAge *pkl.Duration `pkl:"MaxAge"`

	// Allow credentials in CORS requests
	AllowCredentials *bool `pkl:"AllowCredentials"`
}
