// Code generated from Pkl module `org.kdeps.pkl.WebServer`. DO NOT EDIT.
package webserver

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type WebServer interface {
}

var _ WebServer = (*WebServerImpl)(nil)

// Configuration for the Kdeps Web Server
//
// This module defines settings and routes for the Kdeps Web Server, including
// server binding details (host and port) and route configurations. The server
// handles HTTP requests, routing them to appropriate handlers based on defined
// paths and server types.
type WebServerImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a WebServer
func LoadFromPath(ctx context.Context, path string) (ret WebServer, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a WebServer
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (WebServer, error) {
	var ret WebServerImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
