// Code generated from Pkl module `org.kdeps.pkl.Project`. DO NOT EDIT.
package project

import (
	"github.com/kdeps/schema/gen/api_server"
	"github.com/kdeps/schema/gen/docker"
)

type Settings struct {
	ApiServerMode *bool `pkl:"apiServerMode"`

	ApiServerSettings *apiserver.APIServerSettings `pkl:"apiServerSettings"`

	AgentSettings *docker.DockerSettings `pkl:"agentSettings"`
}
