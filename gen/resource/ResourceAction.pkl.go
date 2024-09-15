// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"github.com/kdeps/schema/gen/api"
	"github.com/kdeps/schema/gen/env"
	"github.com/kdeps/schema/gen/llm"
	"github.com/kdeps/schema/gen/project"
	"github.com/kdeps/schema/gen/tag"
)

type ResourceAction struct {
	Name *string `pkl:"name"`

	Exec *string `pkl:"exec"`

	Settings *project.Settings `pkl:"settings"`

	Skip *[]string `pkl:"skip"`

	PreflightCheck *[]string `pkl:"preflightCheck"`

	Env *[]*env.ResourceEnv `pkl:"env"`

	Tags *[]*tag.ResourceTag `pkl:"tags"`

	Chat *[]*llm.ResourceChat `pkl:"chat"`

	HttpClient *[]*api.ResourceHTTPClient `pkl:"httpClient"`
}
