// Code generated from Pkl module `org.kdeps.pkl.Exec`. DO NOT EDIT.
package exec

type ResourceExec struct {
	Env *map[string]string `pkl:"env"`

	Command string `pkl:"command"`

	Stderr *string `pkl:"stderr"`

	Stdout *string `pkl:"stdout"`

	Timestamp *uint32 `pkl:"timestamp"`
}