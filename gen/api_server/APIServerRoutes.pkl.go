// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

// Class representing a route in the API server configuration.
type APIServerRoutes struct {
	// The URL path for the route
	Path string `pkl:"Path"`

	// The HTTP methods for the route (GET, POST, etc.)
	Methods []string `pkl:"Methods"`
}
