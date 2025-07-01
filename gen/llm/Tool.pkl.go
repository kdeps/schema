// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing the details of a tool interaction with an LLM model
type Tool struct {
	// The name of the tool.
	Name *string `pkl:"Name"`

	// The script content to execute for the tool.
	Script *string `pkl:"Script"`

	// A description of what the tool does.
	Description *string `pkl:"Description"`

	// A mapping of parameter names to their properties for tool configuration.
	Parameters *map[string]*ToolProperties `pkl:"Parameters"`
}
