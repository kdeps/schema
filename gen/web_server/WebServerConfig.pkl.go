// Code generated from Pkl module `org.kdeps.pkl.WebServer`. DO NOT EDIT.
package webserver

// Configuration settings for the web server
type WebServerConfig struct {
	// The IP address the server binds to (default: "127.0.0.1")
	Host string `pkl:"host"`

	// The port the server listens on (default: 8080)
	Port uint16 `pkl:"port"`

	// List of routes configured for the server
	//
	// Each route specifies a path and its server behavior
	Routes []*WebServerRouteConfig `pkl:"routes"`
}
