// Code generated from Pkl module `org.kdeps.pkl.Utils`. DO NOT EDIT.
package utils

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Utils interface {
}

var _ Utils = (*UtilsImpl)(nil)

// Tools for Kdeps Resources
//
// This module includes tools for interacting with Kdeps
type UtilsImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Utils
func LoadFromPath(ctx context.Context, path string) (ret Utils, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Utils
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Utils, error) {
	var ret UtilsImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
