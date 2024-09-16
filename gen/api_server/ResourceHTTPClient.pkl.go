// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

type ResourceHTTPClient struct {
	Method string `pkl:"method"`

	Url string `pkl:"url"`

	Data string `pkl:"data"`

	Headers []*ResourceHTTPClientHeaders `pkl:"headers"`
}
