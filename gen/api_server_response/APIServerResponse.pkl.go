// Code generated from Pkl module `org.kdeps.pkl.APIServerResponse`. DO NOT EDIT.
package apiserverresponse

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerResponse struct {
	Success bool `pkl:"success"`

	Response *APIServerResponseBlock `pkl:"response"`

	File *string `pkl:"file"`

	Errors *APIServerErrorsBlock `pkl:"errors"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServerResponse
func LoadFromPath(ctx context.Context, path string) (ret *APIServerResponse, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a APIServerResponse
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*APIServerResponse, error) {
	var ret APIServerResponse
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
