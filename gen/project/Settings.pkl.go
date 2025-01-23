// Code generated from Pkl module `org.kdeps.pkl.Project`. DO NOT EDIT.
package project

import (
	"github.com/kdeps/schema/gen/api_server"
	"github.com/kdeps/schema/gen/docker"
)

// Class representing the settings and configurations for a project.
type Settings struct {
	// Boolean flag to enable or disable API server mode for the project.
	//
	// - `true`: The project runs in API server mode.
	// - `false`: The project does not run in API server mode. Default is `false`.
	APIServerMode bool `pkl:"APIServerMode"`

	// Settings for configuring the API server, which is optional.
	//
	// If API server mode is enabled, these settings provide additional configuration for the API server.
	// [APIServer.APIServerSettings]: Defines the structure and properties for API server settings.
	APIServer *apiserver.APIServerSettings `pkl:"APIServer"`

	// Docker-related settings for the project's agent.
	//
	// These settings define how the Docker agent should be configured for the project.
	// [Docker.DockerSettings]: Includes properties such as Docker image, container settings, and other
	// Docker-specific configurations.
	AgentSettings *docker.DockerSettings `pkl:"agentSettings"`
}
