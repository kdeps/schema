// Code generated from Pkl module `org.kdeps.pkl.Exec`. DO NOT EDIT.
package exec

import "github.com/apple/pkl-go/pkl"

// Class representing an executable resource, which includes the command to be executed,
// its environment variables, and various output/error properties.
type ResourceExec struct {
	// A mapping of environment variable names to their values.
	Env *map[string]string `pkl:"env"`

	// The command to be executed.
	Command string `pkl:"command"`

	// The standard error output from the execution.
	Stderr *string `pkl:"stderr"`

	// The standard output from the execution.
	Stdout *string `pkl:"stdout"`

	// The exit code of the executed command. Defaults to 0.
	ExitCode *int `pkl:"exitCode"`

	// The file path where the stdout value of this resource is saved
	File *string `pkl:"file"`

	// The listing of the item iteration results
	ItemValues *[]string `pkl:"itemValues"`

	// A timestamp of when the command was executed, represented as an unsigned 64-bit integer.
	Timestamp *pkl.Duration `pkl:"timestamp"`

	// The timeout duration (in seconds) for the command execution. Defaults to 60 seconds.
	TimeoutDuration *pkl.Duration `pkl:"timeoutDuration"`
}
