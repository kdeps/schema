// Code generated from Pkl module `org.kdeps.pkl.WebServer`. DO NOT EDIT.
package webserver

// Configuration settings for the web server
type WebServerSettings struct {
	// The IP address the server binds to (default: "127.0.0.1")
	HostIP *string `pkl:"HostIP"`

	// The port the server listens on (default: 8080)
	PortNum *uint16 `pkl:"PortNum"`

	// A list of trusted proxies (IPv4, IPv6, or CIDR ranges).
	// If set, only requests passing through these proxies will have their `X-Forwarded-For`
	// header trusted.
	// If unset, all proxies—including potentially malicious ones—are considered trusted,
	// which may expose the server to IP spoofing and other attacks.
	TrustedProxies *[]string `pkl:"TrustedProxies"`

	// List of routes configured for the server
	//
	// Each route specifies a path and its server behavior
	Routes *[]*WebServerRoutes `pkl:"Routes"`
}
