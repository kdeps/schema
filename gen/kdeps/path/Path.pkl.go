// Code generated from Pkl module `org.kdeps.pkl.Kdeps`. DO NOT EDIT.
package path

import (
	"encoding"
	"fmt"
)

// Defines the paths where Kdeps configurations can be stored.
type Path string

const (
	User    Path = "user"
	Project Path = "project"
	Xdg     Path = "xdg"
)

// String returns the string representation of Path
func (rcv Path) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(Path)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Path.
func (rcv *Path) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "user":
		*rcv = User
	case "project":
		*rcv = Project
	case "xdg":
		*rcv = Xdg
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid Path`, str)
	}
	return nil
}
