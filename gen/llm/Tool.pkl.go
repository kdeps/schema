// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing the details of a tool interaction with an LLM model
type Tool struct {
	// name of the function
	Name *string `pkl:"Name"`

	// path of the script or inline
	Script *string `pkl:"Script"`

	// description of what the tool does
	Description *string `pkl:"Description"`

	// mapping of parameter names to their properties
	Parameters *map[string]ToolProperties `pkl:"Parameters"`
}
