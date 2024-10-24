// Code generated from Pkl module `org.kdeps.pkl.Parameters`. DO NOT EDIT.
package parameters

// Class representing the definition of a parameter argument.
type Args struct {
	// A flag indicating whether this parameter is required. Defaults to `true`.
	Required bool `pkl:"required"`

	// A description of the parameter, providing additional context or usage information.
	Description *string `pkl:"description"`
}
