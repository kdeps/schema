// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llmbackend

import (
	"encoding"
	"fmt"
)

type LLMBackend string

const (
	Ollama LLMBackend = "ollama"
	Api    LLMBackend = "api"
)

// String returns the string representation of LLMBackend
func (rcv LLMBackend) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(LLMBackend)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for LLMBackend.
func (rcv *LLMBackend) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "ollama":
		*rcv = Ollama
	case "api":
		*rcv = Api
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid LLMBackend`, str)
	}
	return nil
}
