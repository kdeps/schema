// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import "github.com/kdeps/schema/gen/llm/llmbackend"

type ResourceChatSchema struct {
	Backend llmbackend.LLMBackend `pkl:"backend"`

	Model string `pkl:"model"`

	Prompt string `pkl:"prompt"`

	Schema string `pkl:"schema"`

	Input *string `pkl:"input"`

	Output string `pkl:"output"`
}
