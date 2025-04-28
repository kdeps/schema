// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing the details of a multi-prompt interaction with an LLM model
type MultiChat struct {
	// The role used to instruct the LLM model.
	Role *string `pkl:"role"`

	// The prompt text sent to the LLM model.
	Prompt *string `pkl:"prompt"`
}
