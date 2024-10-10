// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserverresponsetype

import (
	"encoding"
	"fmt"
)

type APIServerResponseType string

const (
	Json    APIServerResponseType = "json"
	Yaml    APIServerResponseType = "yaml"
	Jsonnet APIServerResponseType = "jsonnet"
	Plist   APIServerResponseType = "plist"
	Xml     APIServerResponseType = "xml"
	Pcf     APIServerResponseType = "pcf"
)

// String returns the string representation of APIServerResponseType
func (rcv APIServerResponseType) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(APIServerResponseType)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for APIServerResponseType.
func (rcv *APIServerResponseType) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "json":
		*rcv = Json
	case "yaml":
		*rcv = Yaml
	case "jsonnet":
		*rcv = Jsonnet
	case "plist":
		*rcv = Plist
	case "xml":
		*rcv = Xml
	case "pcf":
		*rcv = Pcf
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid APIServerResponseType`, str)
	}
	return nil
}
