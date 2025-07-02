// Code generated from Pkl module `org.kdeps.pkl.Exec`. DO NOT EDIT.
package exec

import "github.com/apple/pkl-go/pkl"

// Class representing an executable resource, which includes the command to be executed,
// environment variables, and execution details such as outputs and exit codes.
type ResourceExec struct {
	// A mapping of environment variable names to their values.
	Env *map[string]string `pkl:"Env"`

	// The command to be executed.
	Command string `pkl:"Command"`

	// The standard error output of the command, if any.
	Stderr *string `pkl:"Stderr"`

	// The standard output of the command, if any.
	Stdout *string `pkl:"Stdout"`

	// The exit code of the command. Defaults to 0 (success).
	ExitCode *int `pkl:"ExitCode"`

	// The file path where the command output value of this resource is saved
	File *string `pkl:"File"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"ItemValues"`

	// A timestamp of when the command was executed, represented as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"Timestamp"`

	// The timeout duration (in seconds) for the command execution. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"TimeoutDuration"`
}
