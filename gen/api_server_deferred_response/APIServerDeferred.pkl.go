// Code generated from Pkl module `org.kdeps.pkl.APIServerDeferredResponse`. DO NOT EDIT.
package apiserverdeferredresponse

type APIServerDeferred struct {
	Success bool `pkl:"success"`

	Resolved bool `pkl:"resolved"`

	Response *APIServerDeferredResponseBlock `pkl:"response"`

	Errors *APIServerDeferredErrorsBlock `pkl:"errors"`
}
