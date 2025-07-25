// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import "github.com/apple/pkl-go/pkl"

// Class representing the details of a chat interaction with an LLM model, including prompts, responses,
// file generation, and additional metadata.
type ResourceChat struct {
	// The name of the LLM model used for the chat.
	Model string `pkl:"Model"`

	// The role used to instruct the LLM model.
	Role *string `pkl:"Role"`

	// The prompt text sent to the LLM model.
	Prompt *string `pkl:"Prompt"`

	// A scenario is where a series of conditions to be sent for this chat.
	Scenario *[]*MultiChat `pkl:"Scenario"`

	// Tools that interact with the LLM model
	Tools *[]*Tool `pkl:"Tools"`

	// A listing of file paths or identifiers associated with the chat.
	Files *[]string `pkl:"Files"`

	// Whether the LLM's response is in JSON format. Defaults to `false`.
	JSONResponse *bool `pkl:"JSONResponse"`

	// A listing of keys expected in the JSON response from the LLM model.
	JSONResponseKeys *[]string `pkl:"JSONResponseKeys"`

	// The actual response returned from the LLM model.
	Response *string `pkl:"Response"`

	// The file path where the LLM response of this resource is saved
	File *string `pkl:"File"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"ItemValues"`

	// A timestamp of when the response was generated, represented as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"Timestamp"`

	// The timeout duration (in seconds) for the LLM interaction. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`
}
