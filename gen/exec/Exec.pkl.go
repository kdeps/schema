// Code generated from Pkl module `org.kdeps.pkl.Exec`. DO NOT EDIT.
package exec

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Exec interface {
	GetResources() *map[string]*ResourceExec
}

var _ Exec = (*ExecImpl)(nil)

// This module defines the execution resources for the KDEPS framework.
// It allows for the management and execution of commands, capturing their
// standard output and error, as well as handling environment variables and
// exit codes. The module provides functionalities to retrieve and manage
// executable resources based on their identifiers.
type ExecImpl struct {
	// A mapping of resource IDs to their associated [ResourceExec] objects.
	Resources *map[string]*ResourceExec `pkl:"resources"`
}

// A mapping of resource IDs to their associated [ResourceExec] objects.
func (rcv *ExecImpl) GetResources() *map[string]*ResourceExec {
	return rcv.Resources
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Exec
func LoadFromPath(ctx context.Context, path string) (ret Exec, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Exec
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Exec, error) {
	var ret ExecImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
