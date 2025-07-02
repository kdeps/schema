// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServer interface {
}

var _ APIServer = (*APIServerImpl)(nil)

// Configuration for the Kdeps API Server
//
// This module defines the settings and routes for the Kdeps API Server, including server binding details
// (host and port) and route configurations. The server handles HTTP requests, routing them to appropriate
// handlers based on defined paths and HTTP methods.
type APIServerImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServer
func LoadFromPath(ctx context.Context, path string) (ret APIServer, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a APIServer
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (APIServer, error) {
	var ret APIServerImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
