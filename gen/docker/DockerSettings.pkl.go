// Code generated from Pkl module `org.kdeps.pkl.Docker`. DO NOT EDIT.
package docker

// Class representing the settings for Docker configurations.
// It includes options for specifying packages, PPAs, and models.
type DockerSettings struct {
	// Sets the timezone (see the TZ Identifier here: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
	Timezone string `pkl:"timezone"`

	// Sets if Anaconda3 will be pre-installed in the Image
	InstallAnaconda bool `pkl:"installAnaconda"`

	// Conda packages to install when `installAnaconda` is set to true.
	//
	// Example:
	// condaPackages {
	//   ["base"] { // The name of the Anaconda environment
	//     ["main"] = "diffuser"  // Package "diffuser" from the "main" channel
	//   }
	// }
	CondaPackages *map[string]map[string]string `pkl:"condaPackages"`

	// Python packages that will be pre-installed.
	PythonPackages *[]string `pkl:"pythonPackages"`

	// A list of packages to be installed in the Docker container.
	Packages *[]string `pkl:"packages"`

	// A list of APT or PPA repos to be added.
	Repositories *[]string `pkl:"repositories"`

	// A mandatory list of models to be used in the Docker environment.
	Models []string `pkl:"models"`

	// Sets the Ollama Docker version to be use as the base image
	OllamaImageTag string `pkl:"ollamaImageTag"`

	// A mapping of build arguments variable name
	Args *map[string]string `pkl:"args"`

	// A mapping of build env variable names that persist in the image and container
	Env *map[string]string `pkl:"env"`
}
