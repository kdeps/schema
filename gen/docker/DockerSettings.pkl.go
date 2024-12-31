// Code generated from Pkl module `org.kdeps.pkl.Docker`. DO NOT EDIT.
package docker

// Class representing the settings for Docker configurations.
// It includes options for specifying packages, PPAs, and models.
type DockerSettings struct {
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

	// A listing of params variable names
	Parameters *[]string `pkl:"parameters"`
}
