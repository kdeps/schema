// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import "github.com/apple/pkl-go/pkl"

// Class representing a chat interaction with an LLM model.
type ResourceChat struct {
	// The name of the LLM model to use for the chat interaction.
	Model *string `pkl:"Model"`

	// The role or persona for the chat interaction.
	Role *string `pkl:"Role"`

	// The prompt or message to send to the LLM model.
	Prompt *string `pkl:"Prompt"`

	// The response received from the LLM model.
	Response *string `pkl:"Response"`

	// The file path where the response is stored.
	File *string `pkl:"File"`

	// Whether the response should be in JSON format.
	JSONResponse *bool `pkl:"JSONResponse"`

	// A listing of specific keys to extract from the JSON response.
	JSONResponseKeys *[]string `pkl:"JSONResponseKeys"`

	// The timeout duration for the LLM request.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`

	// The timestamp when the request was made.
	Timestamp *pkl.Duration `pkl:"Timestamp"`

	// The scenario or context for the chat interaction.
	Scenario *[]*MultiChat `pkl:"Scenario"`

	// The tools available for the LLM to use.
	Tools *[]*Tool `pkl:"Tools"`

	// The files associated with the chat interaction.
	Files *[]string `pkl:"Files"`

	// A description of the chat interaction.
	Description *string `pkl:"Description"`

	// The listing of the item iteration results.
	ItemValues *[]string `pkl:"ItemValues"`
}
