// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing a multi-turn chat conversation.
type MultiChat struct {
	// The role or persona for this turn of the conversation.
	Role *string `pkl:"Role"`

	// The prompt text to be sent to the LLM model.
	Prompt *string `pkl:"Prompt"`

	// The content or message for this turn of the conversation.
	Content *string `pkl:"Content"`

	// A description of this turn of the conversation.
	Description *string `pkl:"Description"`
}
