// Code generated from Pkl module `org.kdeps.pkl.Python`. DO NOT EDIT.
package python

import "github.com/apple/pkl-go/pkl"

// Represents an executable Python resource, including its associated script,
// environment variables, and execution details such as outputs and exit codes.
type ResourcePython struct {
	// A mapping of environment variable names to their values.
	Env *map[string]string `pkl:"env"`

	// Specifies the conda environment in which this Python script will execute, if Anaconda is
	// installed.
	CondaEnvironment *string `pkl:"condaEnvironment"`

	// The Python script to be executed.
	Script string `pkl:"script"`

	// Captures the standard error output from the execution.
	Stderr *string `pkl:"stderr"`

	// Captures the standard output from the execution.
	Stdout *string `pkl:"stdout"`

	// The exit code of the executed command. Defaults to 0.
	ExitCode *int `pkl:"exitCode"`

	// The file path where the Python stdout of this resource is saved
	File *string `pkl:"file"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"itemValues"`

	// A timestamp indicating when the command was executed, as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"timestamp"`

	// The maximum duration (in seconds) allowed for the command execution. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"timeoutDuration"`
}
