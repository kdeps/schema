// Code generated from Pkl module `org.kdeps.pkl.APIServerDeferredResponse`. DO NOT EDIT.
package apiserverdeferredresponse

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerDeferredResponse interface {
	GetResources() map[uint32]*APIServerDeferred
}

var _ APIServerDeferredResponse = (*APIServerDeferredResponseImpl)(nil)

// Abstractions for Kdeps API Server Deferred Responses
type APIServerDeferredResponseImpl struct {
	// A mapping of resource IDs to their associated [APIServerDeferred] objects.
	Resources map[uint32]*APIServerDeferred `pkl:"resources"`
}

// A mapping of resource IDs to their associated [APIServerDeferred] objects.
func (rcv *APIServerDeferredResponseImpl) GetResources() map[uint32]*APIServerDeferred {
	return rcv.Resources
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServerDeferredResponse
func LoadFromPath(ctx context.Context, path string) (ret APIServerDeferredResponse, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a APIServerDeferredResponse
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (APIServerDeferredResponse, error) {
	var ret APIServerDeferredResponseImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
