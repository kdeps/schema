// Code generated from Pkl module `org.kdeps.pkl.APIServerResponse`. DO NOT EDIT.
package apiserverresponse

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerResponse interface {
	GetSuccess() *bool

	GetMeta() *APIServerResponseMetaBlock

	GetResponse() *APIServerResponseBlock

	GetErrors() *[]*APIServerErrorsBlock
}

var _ APIServerResponse = (*APIServerResponseImpl)(nil)

// Abstractions for Kdeps API Server Responses
//
// This module provides the structure for handling API server responses in the Kdeps system.
// It includes classes and variables for managing both successful and error responses, as well as
// any files returned by the server. It also defines how data blocks and error blocks are structured
// in the API responses.
//
// This module is part of the `kdeps` schema and interacts with the API server to process responses.
//
// The module defines:
// - [APIServerResponseBlock]: For handling data returned in a successful response.
// - [APIServerErrorsBlock]: For managing error information in a failed API request.
// - [Success]: A flag indicating the success or failure of the API request.
// - [file]: A URI pointing to any file returned by the server in the response.
// - [Errors]: The error block containing details of the error if the request was unsuccessful.
type APIServerResponseImpl struct {
	// A Boolean flag indicating whether the API request was successful.
	//
	// - `true`: The request was successful.
	// - `false`: The request encountered an error.
	Success *bool `pkl:"Success"`

	// Additional metadata related to the API request.
	//
	// Provides request-specific details such as headers, properties, and tracking information.
	Meta *APIServerResponseMetaBlock `pkl:"Meta"`

	// The response block containing data returned by the API server in a successful request, if any.
	//
	// If the request was successful, this block contains the data associated with the response.
	// [APIServerResponseBlock]: Contains a listing of the returned data items.
	Response *APIServerResponseBlock `pkl:"Response"`

	// The error block containing details of any error encountered during the API request.
	//
	// If the request was unsuccessful, this block contains the error code and error message
	// returned by the server.
	// [APIServerErrorsBlock]: Contains the error code and message explaining the issue.
	Errors *[]*APIServerErrorsBlock `pkl:"Errors"`
}

// A Boolean flag indicating whether the API request was successful.
//
// - `true`: The request was successful.
// - `false`: The request encountered an error.
func (rcv *APIServerResponseImpl) GetSuccess() *bool {
	return rcv.Success
}

// Additional metadata related to the API request.
//
// Provides request-specific details such as headers, properties, and tracking information.
func (rcv *APIServerResponseImpl) GetMeta() *APIServerResponseMetaBlock {
	return rcv.Meta
}

// The response block containing data returned by the API server in a successful request, if any.
//
// If the request was successful, this block contains the data associated with the response.
// [APIServerResponseBlock]: Contains a listing of the returned data items.
func (rcv *APIServerResponseImpl) GetResponse() *APIServerResponseBlock {
	return rcv.Response
}

// The error block containing details of any error encountered during the API request.
//
// If the request was unsuccessful, this block contains the error code and error message
// returned by the server.
// [APIServerErrorsBlock]: Contains the error code and message explaining the issue.
func (rcv *APIServerResponseImpl) GetErrors() *[]*APIServerErrorsBlock {
	return rcv.Errors
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServerResponse
func LoadFromPath(ctx context.Context, path string) (ret APIServerResponse, err error) {
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
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (APIServerResponse, error) {
	var ret APIServerResponseImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
