// Code generated from Pkl module `org.kdeps.pkl.APIServerDeferredResponse`. DO NOT EDIT.
package apiserverdeferredresponse

// Class representing the state of a deferred API server response.
type APIServerDeferred struct {
	// Indicates if the request was successful. Defaults to `false`.
	Success bool `pkl:"success"`

	// Indicates if the response has been resolved. Defaults to `false`.
	Resolved bool `pkl:"resolved"`

	// The response block containing data if the request was successful.
	Response *APIServerDeferredResponseBlock `pkl:"response"`

	// The error block if the request encountered issues.
	Errors *APIServerDeferredErrorsBlock `pkl:"errors"`
}
