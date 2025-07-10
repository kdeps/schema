// Code generated from Pkl module `org.kdeps.pkl.PklResource`. DO NOT EDIT.
package pklresource

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type PklResource interface {
}

var _ PklResource = (*PklResourceImpl)(nil)

// Abstractions for PKL records stored in SQLite3
type PklResourceImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a PklResource
func LoadFromPath(ctx context.Context, path string) (ret PklResource, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a PklResource
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (PklResource, error) {
	var ret PklResourceImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
