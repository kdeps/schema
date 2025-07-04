// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

// Class representing an error returned from an API validation check.
type APIError struct {
	// The error code associated with the API error.
	Code *int `pkl:"Code"`

	// A message providing details about the error.
	Message *string `pkl:"Message"`
}
