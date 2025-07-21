// Code generated from Pkl module `org.kdeps.pkl.APIServerRequest`. DO NOT EDIT.
package apiserverrequest

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerRequest interface {
}

var _ APIServerRequest = (*APIServerRequestImpl)(nil)

// Abstractions for KDEPS API Server Request handling
//
// This module provides the structure for handling API server requests in the Kdeps system.
// It includes classes and variables for managing request data such as paths, methods, headers,
// query parameters, and uploaded files. It also provides functions for retrieving and processing
// request information, including file uploads and metadata extraction.
//
// This module is part of the `kdeps` schema and interacts with the API server to process incoming
// requests.
//
// The module defines:
// - [APIServerRequestUploads]: For managing metadata of uploaded files.
// - Functions for retrieving request data from the key-value store.
type APIServerRequestImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServerRequest
func LoadFromPath(ctx context.Context, path string) (ret APIServerRequest, err error) {
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
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (APIServerRequest, error) {
	var ret APIServerRequestImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
