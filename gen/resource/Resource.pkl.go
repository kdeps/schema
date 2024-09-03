// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Resource struct {
	Resources []*AppResource `pkl:"resources"`

	Rag []*RAGResource `pkl:"rag"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Resource
func LoadFromPath(ctx context.Context, path string) (ret *Resource, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Resource
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Resource, error) {
	var ret Resource
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
