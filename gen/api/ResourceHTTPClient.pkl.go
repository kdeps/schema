// Code generated from Pkl module `org.kdeps.pkl.API`. DO NOT EDIT.
package api

type ResourceHTTPClient struct {
	Method string `pkl:"method"`

	Url string `pkl:"url"`

	Data string `pkl:"data"`

	Output *string `pkl:"output"`

	Headers []*ResourceHTTPClientHeaders `pkl:"headers"`
}
