// Code generated from Pkl module `org.kdeps.pkl.Python`. DO NOT EDIT.
package python

import "github.com/apple/pkl-go/pkl"

// Class representing a Python execution resource, which includes the script to be executed,
// environment variables, and execution details such as outputs and exit codes.
type ResourcePython struct {
	// A mapping of environment variable names to their values.
	Env *map[string]string `pkl:"Env"`

	// The Python script to be executed.
	Command string `pkl:"Command"`

	// The standard error output of the script, if any.
	Stderr *string `pkl:"Stderr"`

	// The standard output of the script, if any.
	Stdout *string `pkl:"Stdout"`

	// The exit code of the script. Defaults to 0 (success).
	ExitCode *int `pkl:"ExitCode"`

	// The file path where the script output value of this resource is saved
	File *string `pkl:"File"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"ItemValues"`

	// A timestamp indicating when the command was executed, as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"Timestamp"`

	// The timeout duration (in seconds) for the script execution. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`
}
