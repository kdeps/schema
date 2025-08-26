// Code generated from Pkl module `org.kdeps.pkl.APIServerResponse`. DO NOT EDIT.
package apiserverresponse

// Class representing error details returned in an API response when an error occurs.
type APIServerErrorsBlock struct {
	// The error code returned by the API server, typically an HTTP status code.
	Code int `pkl:"Code"`

	// A descriptive message explaining the error.
	Message string `pkl:"Message"`
}
