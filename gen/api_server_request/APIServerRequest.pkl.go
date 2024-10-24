// Code generated from Pkl module `org.kdeps.pkl.APIServerRequest`. DO NOT EDIT.
package apiserverrequest

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerRequest interface {
	GetPath() string

	GetMethod() string

	GetData() *string

	GetParams() *map[string]string

	GetHeaders() *map[string]string

	GetFiles() *map[string]*APIServerRequestUploads
}

var _ APIServerRequest = (*APIServerRequestImpl)(nil)

// Abstractions for Kdeps API Server Requests
//
// This module provides functionality to handle and validate HTTP requests to the Kdeps API Server, including methods
// for parsing HTTP methods, request data, parameters, headers, and file uploads.
//
// Supported features:
// - Validation of HTTP methods.
// - Handling request body data, parameters, headers, and file uploads.
// - Functions to decode Base64 encoded request data.
// - File management utilities like retrieving file types and paths.
// - Filtering files by MIME type.
type APIServerRequestImpl struct {
	// Represents the request URI path.
	Path string `pkl:"path"`

	// The HTTP method used for the request. Must be a valid method, as determined by [isValidHTTPMethod].
	Method string `pkl:"method"`

	// The body data of the request, which is optional.
	Data *string `pkl:"data"`

	// Query parameters sent with the request.
	Params *map[string]string `pkl:"params"`

	// Headers sent with the request.
	Headers *map[string]string `pkl:"headers"`

	// Files uploaded with the request, represented as a mapping of file keys to upload metadata.
	Files *map[string]*APIServerRequestUploads `pkl:"files"`
}

// Represents the request URI path.
func (rcv *APIServerRequestImpl) GetPath() string {
	return rcv.Path
}

// The HTTP method used for the request. Must be a valid method, as determined by [isValidHTTPMethod].
func (rcv *APIServerRequestImpl) GetMethod() string {
	return rcv.Method
}

// The body data of the request, which is optional.
func (rcv *APIServerRequestImpl) GetData() *string {
	return rcv.Data
}

// Query parameters sent with the request.
func (rcv *APIServerRequestImpl) GetParams() *map[string]string {
	return rcv.Params
}

// Headers sent with the request.
func (rcv *APIServerRequestImpl) GetHeaders() *map[string]string {
	return rcv.Headers
}

// Files uploaded with the request, represented as a mapping of file keys to upload metadata.
func (rcv *APIServerRequestImpl) GetFiles() *map[string]*APIServerRequestUploads {
	return rcv.Files
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
