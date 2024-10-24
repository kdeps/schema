// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

// Class representing the configuration settings for the API server.
type APIServerSettings struct {
	// The IP address the API server will bind to. Defaults to "127.0.0.1".
	HostIP string `pkl:"hostIP"`

	// The port number the API server will listen on. Defaults to 3000.
	PortNum uint16 `pkl:"portNum"`

	// A listing of routes configured for the API server.
	//
	// Each route defines a path and the allowed HTTP methods for that path.
	Routes []*APIServerRoutes `pkl:"routes"`
}
