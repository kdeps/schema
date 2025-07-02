// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

import "github.com/apple/pkl-go/pkl"

// Class representing the configuration settings for the API server.
type APIServerSettings struct {
	// The IP address the server binds to (default: "127.0.0.1")
	HostIP string `pkl:"HostIP"`

	// The port the server listens on (default: 8080)
	PortNum uint16 `pkl:"PortNum"`

	// The timeout duration (in seconds) for API requests. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`

	// A listing of trusted proxies (IPv4, IPv6, or CIDR ranges).
	// If set, only requests passing through these proxies will have their `X-Forwarded-For`
	// header trusted.
	// If unset, all proxies—including potentially malicious ones—are considered trusted,
	// which may expose the server to IP spoofing and other attacks.
	TrustedProxies *[]string `pkl:"TrustedProxies"`

	// List of routes configured for the server
	Routes []*APIServerRoutes `pkl:"Routes"`

	// CORS settings for the API server
	CORS *CORS `pkl:"CORS"`
}
