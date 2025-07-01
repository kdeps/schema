// Code generated from Pkl module `org.kdeps.pkl.Project`. DO NOT EDIT.
package buildenv

import (
	"encoding"
	"fmt"
)

// Defines the environment type.
type BuildEnv string

const (
	Dev  BuildEnv = "dev"
	Prod BuildEnv = "prod"
)

// String returns the string representation of BuildEnv
func (rcv BuildEnv) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(BuildEnv)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for BuildEnv.
func (rcv *BuildEnv) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "dev":
		*rcv = Dev
	case "prod":
		*rcv = Prod
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid BuildEnv`, str)
	}
	return nil
}
