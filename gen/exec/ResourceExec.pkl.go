// Code generated from Pkl module `org.kdeps.pkl.Exec`. DO NOT EDIT.
package exec

type ResourceExec struct {
	Env *map[string]string `pkl:"env"`

	Command string `pkl:"command"`

	Stderr *string `pkl:"stderr"`

	Stdout *string `pkl:"stdout"`

	ExitCode *int `pkl:"exitCode"`

	Timestamp *uint32 `pkl:"timestamp"`

	TimeoutSeconds *int `pkl:"timeoutSeconds"`
}
