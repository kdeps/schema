// Code generated from Pkl module `org.kdeps.pkl.Http`. DO NOT EDIT.
package http

type ResourceHTTPClient struct {
	Method string `pkl:"method"`

	Url string `pkl:"url"`

	Data *[]string `pkl:"data"`

	Headers *map[string]string `pkl:"headers"`

	ResponseData *[]string `pkl:"responseData"`

	ResponseHeaders *map[string]string `pkl:"responseHeaders"`

	Timestamp *uint32 `pkl:"timestamp"`

	TimeoutSeconds *int `pkl:"timeoutSeconds"`
}
