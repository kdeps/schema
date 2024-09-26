// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"github.com/kdeps/schema/gen/api_server"
	"github.com/kdeps/schema/gen/api_server_response"
	"github.com/kdeps/schema/gen/env"
	"github.com/kdeps/schema/gen/llm"
	"github.com/kdeps/schema/gen/tag"
)

type ResourceAction struct {
	Exec *string `pkl:"exec"`

	Env *[]*env.ResourceEnv `pkl:"env"`

	Tags *[]*tag.ResourceTag `pkl:"tags"`

	Chat *[]*llm.ResourceChat `pkl:"chat"`

	SkipCondition *[]bool `pkl:"skipCondition"`

	PreflightCheck *[]bool `pkl:"preflightCheck"`

	PostflightCheck *[]bool `pkl:"postflightCheck"`

	HttpClient *[]*apiserver.ResourceHTTPClient `pkl:"httpClient"`

	ApiResponse *apiserverresponse.APIServerResponse `pkl:"apiResponse"`
}
