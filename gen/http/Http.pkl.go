// Code generated from Pkl module `org.kdeps.pkl.HTTP`. DO NOT EDIT.
package http

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/utils"
)

type HTTP interface {
	utils.Utils

	GetResources() *map[string]*ResourceHTTPClient
}

var _ HTTP = (*HTTPImpl)(nil)

// This module defines the settings and configurations for HTTP client
// resources within the KDEPS framework. It enables the management of
// HTTP requests, including method specifications, request data, headers,
// and handling of responses. This module provides functionalities to
// retrieve and manage HTTP client resources based on their identifiers.
type HTTPImpl struct {
	*utils.UtilsImpl

	// A mapping of resource IDs to their associated [ResourceHTTPClient] objects.
	Resources *map[string]*ResourceHTTPClient `pkl:"resources"`
}

// A mapping of resource IDs to their associated [ResourceHTTPClient] objects.
func (rcv *HTTPImpl) GetResources() *map[string]*ResourceHTTPClient {
	return rcv.Resources
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a HTTP
func LoadFromPath(ctx context.Context, path string) (ret HTTP, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a HTTP
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (HTTP, error) {
	var ret HTTPImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
