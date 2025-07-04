// Code generated from Pkl module `org.kdeps.pkl.Docker`. DO NOT EDIT.
package docker

// Class representing the settings for Docker configurations.
// It includes options for specifying packages, PPAs, and models.
type DockerSettings struct {
	// A mapping of environment variable names to their values.
	Params *map[string]string `pkl:"Params"`

	// Sets the tag version to be use as the base image
	OllamaTagVersion *string `pkl:"OllamaTagVersion"`

	// Install Anaconda Python on the Docker container
	InstallAnaconda *bool `pkl:"InstallAnaconda"`

	// Sets the timezone (see the TZ Identifier here: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
	Timezone *string `pkl:"Timezone"`

	// Conda packages to install when `InstallAnaconda` is set to true.
	//
	// Example:
	// CondaPackages {
	//   ["base"] { // The name of the Anaconda environment
	//     ["main"] = "diffuser"  // Package "diffuser" from the "main" channel
	//   }
	// }
	CondaPackages *map[string]map[string]string `pkl:"CondaPackages"`

	// Python packages that will be pre-installed.
	PythonPackages *[]string `pkl:"PythonPackages"`

	// A list of packages to be installed in the Docker container.
	Packages *[]string `pkl:"Packages"`

	// A list of APT or PPA repos to be added.
	Repositories *[]string `pkl:"Repositories"`

	// A mandatory list of LLM models to be used in the Docker environment.
	Models []string `pkl:"Models"`

	// A mapping of build arguments variable name
	Args *map[string]string `pkl:"Args"`

	// A mapping of build env variable names that persist in the image and container
	Env *map[string]string `pkl:"Env"`

	// A list of ports to be exposed in the Docker container.
	ExposedPorts *[]string `pkl:"ExposedPorts"`
}
