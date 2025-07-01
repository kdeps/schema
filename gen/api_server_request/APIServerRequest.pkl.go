// Code generated from Pkl module `org.kdeps.pkl.APIServerRequest`. DO NOT EDIT.
package apiserverrequest

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerRequest interface {
	GetPath() string

	GetIP() string

	GetID() string

	GetMethod() string

	GetData() *string

	GetParams() *map[string]string

	GetHeaders() *map[string]string

	GetFiles() *map[string]*APIServerRequestUploads
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
// - [Path]: The URI path of the incoming request.
// - [Method]: The HTTP method used for the request.
// - [Data]: The request body data.
// - [Files]: A mapping of uploaded files and their metadata.
type APIServerRequestImpl struct {
	// The URI path of the incoming request.
	Path string `pkl:"Path"`

	// The Client IP Address
	IP string `pkl:"IP"`

	// The Request ID
	ID string `pkl:"ID"`

	// The HTTP method used for the request. Must be a valid method, as determined by [isValidHTTPMethod].
	Method string `pkl:"Method"`

	// The request body, if provided.
	Data *string `pkl:"Data"`

	// A mapping of query parameters included in the request.
	Params *map[string]string `pkl:"Params"`

	// A mapping of HTTP headers included in the request.
	Headers *map[string]string `pkl:"Headers"`

	// Files uploaded with the request, represented as a mapping of file keys to upload metadata.
	Files *map[string]*APIServerRequestUploads `pkl:"Files"`
}

// The URI path of the incoming request.
func (rcv *APIServerRequestImpl) GetPath() string {
	return rcv.Path
}

// The Client IP Address
func (rcv *APIServerRequestImpl) GetIP() string {
	return rcv.IP
}

// The Request ID
func (rcv *APIServerRequestImpl) GetID() string {
	return rcv.ID
}

// The HTTP method used for the request. Must be a valid method, as determined by [isValidHTTPMethod].
func (rcv *APIServerRequestImpl) GetMethod() string {
	return rcv.Method
}

// The request body, if provided.
func (rcv *APIServerRequestImpl) GetData() *string {
	return rcv.Data
}

// A mapping of query parameters included in the request.
func (rcv *APIServerRequestImpl) GetParams() *map[string]string {
	return rcv.Params
}

// A mapping of HTTP headers included in the request.
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
