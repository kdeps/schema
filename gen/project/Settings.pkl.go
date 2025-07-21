// Code generated from Pkl module `org.kdeps.pkl.Project`. DO NOT EDIT.
package project

import (
	"github.com/kdeps/schema/gen/api_server"
	"github.com/kdeps/schema/gen/docker"
	"github.com/kdeps/schema/gen/project/buildenv"
	"github.com/kdeps/schema/gen/web_server"
)

// Class representing the settings and configurations for a project.
type Settings struct {
	// Boolean flag to enable or disable API server mode for the project.
	//
	// - `true`: The project runs in API server mode.
	// - `false`: The project does not run in API server mode. Default is `false`.
	APIServerMode *bool `pkl:"APIServerMode"`

	// Settings for configuring the API server, which is optional.
	//
	// If API server mode is enabled, these settings provide additional configuration for the API server.
	// [APIServer.APIServerSettings]: Defines the structure and properties for API server settings.
	APIServer *apiserver.APIServerSettings `pkl:"APIServer"`

	// Boolean flag to enable or disable Web server mode for the project.
	//
	// - `true`: The project runs in Web server mode.
	// - `false`: The project does not run in Web server mode. Default is `false`.
	WebServerMode *bool `pkl:"WebServerMode"`

	// Settings for configuring the Web server, which is optional.
	//
	// If Web server mode is enabled, these settings provide additional configuration for the Web server.
	// [WebServer.WebServerConfig]: Defines the structure and properties for Web server settings.
	WebServer *webserver.WebServerSettings `pkl:"WebServer"`

	// Docker-related settings for the project's agent.
	//
	// These settings define how the Docker agent should be configured for the project.
	// [Docker.DockerSettings]: Includes properties such as docker image, container settings, and other
	// Docker-specific configurations.
	AgentSettings *docker.DockerSettings `pkl:"AgentSettings"`

	// Maximum number of concurrent requests allowed in the workflow.
	//
	// This setting controls the rate limiting behavior for workflow execution.
	// Default value is 5 concurrent requests.
	RateLimitMax *int `pkl:"RateLimitMax"`

	// Environment setting for the workflow execution.
	//
	// Specifies whether the workflow runs in development or production mode.
	// Valid values are "dev", "development", "prod", or "production".
	//
	// In production mode ("prod" or "production"):
	// - Gin framework runs in release mode (no debug output)
	// - Log level is set to WARN (less verbose)
	// - DEBUG environment variable is set to 0
	// - Debug logs are suppressed
	//
	// In development mode ("dev" or "development"):
	// - Gin framework runs in debug mode (verbose output)
	// - Log level is set to INFO or DEBUG (based on DEBUG env var)
	// - DEBUG environment variable is set to 1
	// - Full logging and debug information available
	Environment *buildenv.BuildEnv `pkl:"Environment"`
}
