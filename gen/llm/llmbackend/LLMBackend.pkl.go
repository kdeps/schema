// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llmbackend

import (
	"encoding"
	"fmt"
)

type LLMBackend string

const (
	Ollama      LLMBackend = "ollama"
	Openai      LLMBackend = "openai"
	Mistral     LLMBackend = "mistral"
	Huggingface LLMBackend = "huggingface"
	Groq        LLMBackend = "groq"
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
	case "openai":
		*rcv = Openai
	case "mistral":
		*rcv = Mistral
	case "huggingface":
		*rcv = Huggingface
	case "groq":
		*rcv = Groq
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid LLMBackend`, str)
	}
	return nil
}
