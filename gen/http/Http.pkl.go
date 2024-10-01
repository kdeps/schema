// Code generated from Pkl module `org.kdeps.pkl.Http`. DO NOT EDIT.
package http

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Http struct {
	Resource *map[string]*ResourceHTTPClient `pkl:"resource"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Http
func LoadFromPath(ctx context.Context, path string) (ret *Http, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Http
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Http, error) {
	var ret Http
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
