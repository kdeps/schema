// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"github.com/kdeps/schema/pkg/core/api"
	"github.com/kdeps/schema/pkg/core/env"
	"github.com/kdeps/schema/pkg/core/llm"
	"github.com/kdeps/schema/pkg/core/settings"
)

type ResourceAction struct {
	Name string `pkl:"name"`

	Exec string `pkl:"exec"`

	Settings *settings.AppSettings `pkl:"settings"`

	Skip *[]string `pkl:"skip"`

	Check *[]string `pkl:"check"`

	Expect *[]string `pkl:"expect"`

	Env *[]*env.ResourceEnv `pkl:"env"`

	Chat *[]*llm.ResourceChat `pkl:"chat"`

	Schat *[]*llm.ResourceChatSchema `pkl:"schat"`

	Api *[]*api.ResourceAPI `pkl:"api"`
}
