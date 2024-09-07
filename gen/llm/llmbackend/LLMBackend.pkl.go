// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llmbackend

import (
	"encoding"
	"fmt"
)

type LLMBackend string

const (
	Local          LLMBackend = "local"
	OpenaiApi      LLMBackend = "openai-api"
	MistralApi     LLMBackend = "mistral-api"
	HuggingfaceApi LLMBackend = "huggingface-api"
	GroqApi        LLMBackend = "groq-api"
)

// String returns the string representation of LLMBackend
func (rcv LLMBackend) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(LLMBackend)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for LLMBackend.
func (rcv *LLMBackend) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "local":
		*rcv = Local
	case "openai-api":
		*rcv = OpenaiApi
	case "mistral-api":
		*rcv = MistralApi
	case "huggingface-api":
		*rcv = HuggingfaceApi
	case "groq-api":
		*rcv = GroqApi
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid LLMBackend`, str)
	}
	return nil
}
