// Code generated from Pkl module `org.kdeps.pkl.APIServerDeferredResponse`. DO NOT EDIT.
package apiserverdeferredresponse

// Class representing an error block in deferred API server responses.
type APIServerDeferredErrorsBlock struct {
	// The error code associated with the response.
	Code int `pkl:"code"`

	// The error message describing the issue.
	Message string `pkl:"message"`
}
