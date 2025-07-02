// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import "github.com/apple/pkl-go/pkl"

// Class representing the details of a chat interaction with an LLM model, including prompts, responses,
// and configuration options such as tools, JSON output, and timeout settings.
type ResourceChat struct {
	// The model to use for the conversation (e.g., "llama3.2").
	Model string `pkl:"Model"`

	// The role of the speaker in the conversation (e.g., "user", "assistant").
	Role *string `pkl:"Role"`

	// The prompt text to be sent to the LLM model.
	Prompt *string `pkl:"Prompt"`

	// A listing of multi-prompt scenarios to be executed in sequence.
	Scenario *[]*MultiChat `pkl:"Scenario"`

	// A listing of tools available for the LLM to use during the conversation.
	Tools *[]*Tool `pkl:"Tools"`

	// A listing of file paths that the LLM can reference or access.
	Files *[]string `pkl:"Files"`

	// A flag indicating whether the response should be in JSON format.
	JSONResponse *bool `pkl:"JSONResponse"`

	// A listing of specific keys to extract from the JSON response.
	JSONResponseKeys *[]string `pkl:"JSONResponseKeys"`

	// The response text returned by the LLM model.
	Response *string `pkl:"Response"`

	// The file path where the response is saved.
	File *string `pkl:"File"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"ItemValues"`

	// A timestamp of when the response was generated, represented as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"Timestamp"`

	// The timeout duration (in seconds) for the LLM request. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`
}
