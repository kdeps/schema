// Code generated from Pkl module `org.kdeps.pkl.Kdeps`. DO NOT EDIT.
package mode

import (
	"encoding"
	"fmt"
)

// Defines the mode of execution for Kdeps.
type Mode string

const (
	Docker Mode = "docker"
	Local  Mode = "local"
)

// String returns the string representation of Mode
func (rcv Mode) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(Mode)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Mode.
func (rcv *Mode) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "docker":
		*rcv = Docker
	case "local":
		*rcv = Local
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid Mode`, str)
	}
	return nil
}
