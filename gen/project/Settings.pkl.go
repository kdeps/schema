// Code generated from Pkl module `org.kdeps.pkl.Project`. DO NOT EDIT.
package project

import (
	"github.com/kdeps/schema/gen/api_server"
	"github.com/kdeps/schema/gen/docker"
	"github.com/kdeps/schema/gen/security"
)

type Settings struct {
	ApiServerMode bool `pkl:"apiServerMode"`

	ApiServer *apiserver.APIServerSettings `pkl:"apiServer"`

	AgentSettings *docker.DockerSettings `pkl:"agentSettings"`

	SecuritySettings *security.Settings `pkl:"securitySettings"`
}
