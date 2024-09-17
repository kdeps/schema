// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

import "github.com/kdeps/schema/gen/api_server/apiserverresponsetype"

type APIServerRoutes struct {
	Path string `pkl:"path"`

	Methods []string `pkl:"methods"`

	ResponseType apiserverresponsetype.APIServerResponseType `pkl:"responseType"`
}
