// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import "github.com/kdeps/schema/gen/llm/llmbackend"

type LLMSettings struct {
	LlmAPIKeys *LLMAPIKeys `pkl:"llmAPIKeys"`

	LlmFallbackBackend llmbackend.LLMBackend `pkl:"llmFallbackBackend"`

	LlmFallbackModel string `pkl:"llmFallbackModel"`
}
