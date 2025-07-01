// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing the details of a multi-prompt interaction with an LLM model
type MultiChat struct {
	// The role of the speaker in the conversation (e.g., "user", "assistant").
	Role *string `pkl:"Role"`

	// The prompt text to be sent to the LLM model.
	Prompt *string `pkl:"Prompt"`
}
