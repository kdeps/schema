// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package roletype

import (
	"encoding"
	"fmt"
)

// Type of LLM roles
type RoleType string

const (
	User      RoleType = "user"
	System    RoleType = "system"
	Assistant RoleType = "assistant"
)

// String returns the string representation of RoleType
func (rcv RoleType) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(RoleType)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for RoleType.
func (rcv *RoleType) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "user":
		*rcv = User
	case "system":
		*rcv = System
	case "assistant":
		*rcv = Assistant
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid RoleType`, str)
	}
	return nil
}
