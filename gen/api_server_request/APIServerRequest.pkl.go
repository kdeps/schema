// Code generated from Pkl module `org.kdeps.pkl.APIServerRequest`. DO NOT EDIT.
package apiserverrequest

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerRequest struct {
	Url *string `pkl:"url"`

	Data []*APIServerRequestData `pkl:"data"`

	Method *string `pkl:"method"`

	Params []*APIServerRequestParams `pkl:"params"`

	Headers []*APIServerRequestHeader `pkl:"headers"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServerRequest
func LoadFromPath(ctx context.Context, path string) (ret *APIServerRequest, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a APIServerRequest
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*APIServerRequest, error) {
	var ret APIServerRequest
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}