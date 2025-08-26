// Code generated from Pkl module `org.kdeps.pkl.Skip`. DO NOT EDIT.
package skip

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Skip interface {
}

var _ Skip = (*SkipImpl)(nil)

// Skip condition functions used across all resources.
//
// Tools for creating skip logic validations
type SkipImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Skip
func LoadFromPath(ctx context.Context, path string) (ret Skip, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Skip
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Skip, error) {
	var ret SkipImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
