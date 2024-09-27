// Code generated from Pkl module `org.kdeps.pkl.Http`. DO NOT EDIT.
package http

type ResourceHTTPClient struct {
	Method string `pkl:"method"`

	Url string `pkl:"url"`

	Data *string `pkl:"data"`

	Headers *map[string]string `pkl:"headers"`

	Response *string `pkl:"response"`

	ResponseData *string `pkl:"responseData"`

	ResponseBody *string `pkl:"responseBody"`

	ResponseHeaders *map[string]string `pkl:"responseHeaders"`
}
