// Code generated from Pkl module `org.kdeps.pkl.Python`. DO NOT EDIT.
package python

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/utils"
)

type Python interface {
	utils.Utils

	GetResources() map[string]*ResourcePython
}

var _ Python = (*PythonImpl)(nil)

// Abstractions for Python script execution within KDEPS
//
// This module defines the structure for Python execution resources that can be used within the Kdeps framework.
// It handles Python script execution, environment variable management, capturing outputs,
// variables as well as exit codes. The module provides utilities for retrieving
// and managing Python execution resources based on their identifiers.
type PythonImpl struct {
	*utils.UtilsImpl

	// A mapping of resource actionIDs to their associated [ResourcePython] objects.
	// This mapping is populated from pklres storage.
	Resources map[string]*ResourcePython `pkl:"Resources"`
}

// A mapping of resource actionIDs to their associated [ResourcePython] objects.
// This mapping is populated from pklres storage.
func (rcv *PythonImpl) GetResources() map[string]*ResourcePython {
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
