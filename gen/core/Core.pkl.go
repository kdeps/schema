// Code generated from Pkl module `org.kdeps.pkl.Core`. DO NOT EDIT.
package core

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Core interface {
}

var _ Core = (*CoreImpl)(nil)

// Core abstractions for Kdeps operations
// This module provides unified functions for agent resolution and generic pklres operations
type CoreImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Core
func LoadFromPath(ctx context.Context, path string) (ret Core, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Core
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Core, error) {
	var ret CoreImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
