// Code generated from Pkl module `org.kdeps.pkl.WebServer`. DO NOT EDIT.
package webserver

import "github.com/kdeps/schema/gen/web_server/webservertype"

// Configuration for a server route
type WebServerRoutes struct {
	// The URL path for the route
	Path string `pkl:"Path"`

	// Optional port for the application to be proxied.
	// Only applicable if serverType is "app". (default: 8052)
	AppPort *uint16 `pkl:"AppPort"`

	// Type of web server for this route, can either be "app" or "static". (default: "static")
	ServerType *webservertype.WebServerType `pkl:"ServerType"`

	// Public path relative to the "/data/" directory (default: "/web")
	PublicPath *string `pkl:"PublicPath"`

	// Optional command to execute for the route. Only applicable if serverType is "app".
	Command *string `pkl:"Command"`
}
