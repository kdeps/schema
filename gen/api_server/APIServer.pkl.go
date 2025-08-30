// Code generated from Pkl module `org.kdeps.pkl.APIServer`. DO NOT EDIT.
package apiserver

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServer interface {
}

var _ APIServer = APIServerImpl{}

// Abstractions for Kdeps API Server Configuration
//
// This module defines the settings and routes for configuring the Kdeps API Server. It includes
// server settings such as host IP and port number, as well as route definitions. The API server
// is designed to handle incoming requests and route them to the appropriate handlers, ensuring
// proper management of HTTP methods and deferred processing.
type APIServerImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServer
func LoadFromPath(ctx context.Context, path string) (ret APIServer, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return ret, err
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
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
