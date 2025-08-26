// Code generated from Pkl module `org.kdeps.pkl.WebServer`. DO NOT EDIT.
package webservertype

import (
	"encoding"
	"fmt"
)

// Type of web server
type WebServerType string

const (
	Static WebServerType = "static"
	App    WebServerType = "app"
)

// String returns the string representation of WebServerType
func (rcv WebServerType) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(WebServerType)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for WebServerType.
func (rcv *WebServerType) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "static":
		*rcv = Static
	case "app":
		*rcv = App
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid WebServerType`, str)
	}
	return nil
}
