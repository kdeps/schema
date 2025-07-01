// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

// Class representing a route in the API server configuration.
type APIServerRoutes struct {
	// The URL path for the route
	Path string `pkl:"Path"`

	// The HTTP method for the route (GET, POST, etc.)
	Method string `pkl:"Method"`

	// The action ID that this route maps to
	ActionID string `pkl:"ActionID"`
}
