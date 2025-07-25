// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

// Class representing a route in the API server configuration.
type APIServerRoutes struct {
	// The path for the route in the API server.
	Path string `pkl:"Path"`

	// A listing of allowed HTTP methods for this route, validated by the HTTP method regex.
	Methods []string `pkl:"Methods"`
}
