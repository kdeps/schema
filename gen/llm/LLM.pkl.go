// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/utils"
)

type LLM interface {
	utils.Utils

	GetRequestID() *string

	GetResources() *map[string]*ResourceChat
}

var _ LLM = (*LLMImpl)(nil)

// Abstractions for KDEPS LLM Resources
//
// This module defines the structure for managing large language model (LLM) resources
// related to LLM model interactions. The class allows for managing prompts, responses,
// and additional configurations such as tools, scenarios, and output files. It also
// provides utilities for retrieving and managing LLM resources based on their identifiers.
//
// The module includes:
// - [ResourceChat]: A class for handling individual chat interactions with LLM models.
// - [Resource]: Mapping of resource actionIDs to [ResourceChat] objects.
type LLMImpl struct {
	*utils.UtilsImpl

	// The current request ID for pklres operations (injected by Go code)
	RequestID *string `pkl:"requestID"`

	// A mapping of resource actionIDs to their associated [ResourceChat] objects.
	Resources *map[string]*ResourceChat `pkl:"Resources"`
}

// The current request ID for pklres operations (injected by Go code)
func (rcv *LLMImpl) GetRequestID() *string {
	return rcv.RequestID
}

// A mapping of resource actionIDs to their associated [ResourceChat] objects.
func (rcv *LLMImpl) GetResources() *map[string]*ResourceChat {
	return rcv.Resources
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a LLM
func LoadFromPath(ctx context.Context, path string) (ret LLM, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a LLM
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (LLM, error) {
	var ret LLMImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
