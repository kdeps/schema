// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

// Class representing a route in the API server configuration.
type APIServerRoutes struct {
	// The path for the route in the API server.
	Path string `pkl:"path"`

	// A listing of allowed HTTP methods for this route, validated by the HTTP method regex.
	Methods []string `pkl:"methods"`

	// A Boolean flag indicating whether the route should be processed in a deferred manner.
	//
	// - `true`: The route will be processed after the initial request handling.
	// - `false`: The route will be processed immediately. Default is `false`.
	DeferredApi *bool `pkl:"deferredApi"`
}
