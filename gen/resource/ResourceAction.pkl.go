// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"github.com/kdeps/schema/gen/api"
	"github.com/kdeps/schema/gen/env"
	"github.com/kdeps/schema/gen/llm"
	"github.com/kdeps/schema/gen/project"
)

type ResourceAction struct {
	Name string `pkl:"name"`

	Exec string `pkl:"exec"`

	Settings *project.Settings `pkl:"settings"`

	Skip *[]string `pkl:"skip"`

	Preflight *[]string `pkl:"preflight"`

	Env *[]*env.ResourceEnv `pkl:"env"`

	Chat *[]*llm.ResourceChat `pkl:"chat"`

	HttpClient *[]*api.ResourceHTTPClient `pkl:"httpClient"`
}
