// Code generated from Pkl module `org.kdeps.pkl.Parameters`. DO NOT EDIT.
package parameters

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Parameters interface {
	GetTags() *map[string]*Args
}

var _ Parameters = (*ParametersImpl)(nil)

// Abstractions for Kdeps Parameters
//
// This module defines the structure and validation for parameters used in Kdeps workflows.
// It includes a mapping of parameter IDs to their corresponding argument definitions,
// allowing for clear documentation and validation of input arguments in the system.
type ParametersImpl struct {
	// A mapping of parameter IDs to their corresponding [Args] objects.
	Tags *map[string]*Args `pkl:"tags"`
}

// A mapping of parameter IDs to their corresponding [Args] objects.
func (rcv *ParametersImpl) GetTags() *map[string]*Args {
	return rcv.Tags
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Parameters
func LoadFromPath(ctx context.Context, path string) (ret Parameters, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Parameters
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Parameters, error) {
	var ret ParametersImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
