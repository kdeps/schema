// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

// Class representing the configuration settings for the API server.
type APIServerSettings struct {
	// The IP address the API server will bind to. Defaults to "127.0.0.1".
	HostIP string `pkl:"HostIP"`

	// The port number the API server will listen on. Defaults to 3000.
	PortNum uint16 `pkl:"PortNum"`

	// A list of trusted proxies (IPv4, IPv6, or CIDR ranges).
	// If set, only requests passing through these proxies will have their `X-Forwarded-For`
	// header trusted.
	// If unset, all proxies—including potentially malicious ones—are considered trusted,
	// which may expose the server to IP spoofing and other attacks.
	TrustedProxies *[]string `pkl:"TrustedProxies"`

	// A listing of routes configured for the API server.
	//
	// Each route defines a path and the allowed HTTP methods for that path.
	Routes []APIServerRoutes `pkl:"Routes"`

	// Cross-Origin Resource Sharing (CORS) configuration
	CORS CORSConfig `pkl:"CORS"`
}
