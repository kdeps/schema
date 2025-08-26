// Code generated from Pkl module `org.kdeps.pkl.Python`. DO NOT EDIT.
package python

import "github.com/apple/pkl-go/pkl"

// Represents an executable Python resource, including its associated script,
// environment variables, and execution details such as outputs and exit codes.
type ResourcePython struct {
	// A mapping of environment variable names to their values.
	Env *map[string]string `pkl:"Env"`

	// Specifies the conda environment in which this Python script will execute, if Anaconda is
	// installed.
	CondaEnvironment *string `pkl:"CondaEnvironment"`

	// The Python script to be executed.
	Script string `pkl:"Script"`

	// Captures the standard error output from the execution.
	Stderr *string `pkl:"Stderr"`

	// Captures the standard output from the execution.
	Stdout *string `pkl:"Stdout"`

	// The exit code of the executed command. Defaults to 0.
	ExitCode *int `pkl:"ExitCode"`

	// The file path where the Python stdout of this resource is saved
	File *string `pkl:"File"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"ItemValues"`

	// A timestamp indicating when the command was executed, as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"Timestamp"`

	// The maximum duration (in seconds) allowed for the command execution. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`
}
