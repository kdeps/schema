// Code generated from Pkl module `org.kdeps.pkl.Kdeps`. DO NOT EDIT.
package runmode

import (
	"encoding"
	"fmt"
)

type RunMode string

const (
	Docker RunMode = "docker"
	Local  RunMode = "local"
)

// String returns the string representation of RunMode
func (rcv RunMode) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(RunMode)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for RunMode.
func (rcv *RunMode) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "docker":
		*rcv = Docker
	case "local":
		*rcv = Local
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid RunMode`, str)
	}
	return nil
}
