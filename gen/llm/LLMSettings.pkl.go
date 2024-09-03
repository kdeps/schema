// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import "github.com/kdeps/schema/gen/llm/llmbackend"

type LLMSettings struct {
	LlmAPIKeys *LLMAPIKeys `pkl:"llmAPIKeys"`

	LlmBackend llmbackend.LLMBackend `pkl:"llmBackend"`

	LlmModel string `pkl:"llmModel"`

	ModelFile *ModelFile `pkl:"modelFile"`
}
