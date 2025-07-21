// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing a single parameter's properties in a tool definition.
type ToolProperties struct {
	// Indicates if the parameter is required for the tool to function.
	Required *bool `pkl:"Required"`

	// The data type of the parameter (e.g., "string", "integer").
	Type *string `pkl:"Type"`

	// A description of the parameter's purpose.
	Description *string `pkl:"Description"`
}
