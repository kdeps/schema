// Code generated from Pkl module `org.kdeps.pkl.Http`. DO NOT EDIT.
package http

type ResourceHTTPClient struct {
	Method string `pkl:"method"`

	Url string `pkl:"url"`

	Data *[]string `pkl:"data"`

	Headers *map[string]string `pkl:"headers"`

	Response *ResponseBlock `pkl:"response"`

	Timestamp *uint32 `pkl:"timestamp"`

	TimeoutSeconds *int `pkl:"timeoutSeconds"`
}
