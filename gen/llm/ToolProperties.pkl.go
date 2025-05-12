// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing a single parameter's properties in a tool definition
type ToolProperties struct {
	// Whether this parameter is required; defaults to true
	Required *bool `pkl:"required"`

	// Data type of the parameter (e.g., 'string', 'integer')
	Type *string `pkl:"type"`

	// description of the parameter for clarity
	Description *string `pkl:"description"`
}
