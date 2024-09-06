// Code generated from Pkl module `org.kdeps.pkl.Kdeps`. DO NOT EDIT.
package gpu

import (
	"encoding"
	"fmt"
)

type GPU string

const (
	Nvidia GPU = "nvidia"
	Amd    GPU = "amd"
	Cpu    GPU = "cpu"
)

// String returns the string representation of GPU
func (rcv GPU) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(GPU)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for GPU.
func (rcv *GPU) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "nvidia":
		*rcv = Nvidia
	case "amd":
		*rcv = Amd
	case "cpu":
		*rcv = Cpu
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid GPU`, str)
	}
	return nil
}
