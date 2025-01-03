// Code generated from Pkl module `org.kdeps.pkl.Http`. DO NOT EDIT.
package http

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/utils"
)

type Http interface {
	utils.Utils

	GetResources() *map[string]*ResourceHTTPClient
}

var _ Http = (*HttpImpl)(nil)

// This module defines the settings and configurations for HTTP client
// resources within the KDEPS framework. It enables the management of
// HTTP requests, including method specifications, request data, headers,
// and handling of responses. This module provides functionalities to
// retrieve and manage HTTP client resources based on their identifiers.
type HttpImpl struct {
	*utils.UtilsImpl

	// A mapping of resource IDs to their associated [ResourceHTTPClient] objects.
	Resources *map[string]*ResourceHTTPClient `pkl:"resources"`
}

// A mapping of resource IDs to their associated [ResourceHTTPClient] objects.
func (rcv *HttpImpl) GetResources() *map[string]*ResourceHTTPClient {
	return rcv.Resources
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Http
func LoadFromPath(ctx context.Context, path string) (ret Http, err error) {
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
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Http, error) {
	var ret HttpImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
