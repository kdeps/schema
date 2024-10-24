// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

// Class representing validation checks that can be performed on actions.
type ValidationCheck struct {
	// A listing of validation conditions.
	Validations *[]bool `pkl:"validations"`

	// An error associated with the validation check, if any.
	Error *APIError `pkl:"error"`
}
