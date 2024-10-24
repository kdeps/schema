// Code generated from Pkl module `org.kdeps.pkl.Docker`. DO NOT EDIT.
package docker

// Class representing the settings for Docker configurations.
// It includes options for specifying packages, PPAs, and models.
type DockerSettings struct {
	// A list of packages to be installed in the Docker container.
	Packages *[]string `pkl:"packages"`

	// A list of Personal Package Archives (PPAs) to be added.
	Ppa *[]string `pkl:"ppa"`

	// A mandatory list of models to be used in the Docker environment.
	Models []string `pkl:"models"`
}
