// Code generated from Pkl module `org.kdeps.pkl.Python`. DO NOT EDIT.
package python

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/utils"
)

type Python interface {
	utils.Utils

	GetResources() *map[string]*ResourcePython
}

var _ Python = (*PythonImpl)(nil)

// This module defines the execution resources for the KDEPS framework.
// It facilitates the management and execution of Python-based commands,
// capturing their standard output, standard error, and handling environment
// variables as well as exit codes. The module provides utilities for retrieving
// and managing executable resources identified by unique resource IDs.
type PythonImpl struct {
	*utils.UtilsImpl

	// A mapping of resource IDs to their corresponding [ResourcePython] objects.
	Resources *map[string]*ResourcePython `pkl:"resources"`
}

// A mapping of resource IDs to their corresponding [ResourcePython] objects.
func (rcv *PythonImpl) GetResources() *map[string]*ResourcePython {
	return rcv.Resources
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Python
func LoadFromPath(ctx context.Context, path string) (ret Python, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Python
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Python, error) {
	var ret PythonImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
