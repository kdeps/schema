// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type LLM interface {
	GetResources() *map[string]*ResourceChat
}

var _ LLM = (*LLMImpl)(nil)

// Abstractions for Kdeps LLM Resource
//
// This module provides an abstraction layer for managing resources related to
// large language model (LLM) interactions within the Kdeps system.
//
// It defines the [ResourceChat] class, which encapsulates the metadata and responses
// related to LLM model interactions. The class allows for managing prompts, responses,
// file generations, image generation flags, and the handling of JSON responses.
//
// Key functionalities include:
// - Managing a collection of resources that represent LLM interactions through a mapping of unique
// resource IDs to [ResourceChat] objects.
// - Providing methods to retrieve various pieces of information related to the LLM interaction,
// such as the prompt text, response text, file paths, JSON keys, and whether image generation was
// involved.
type LLMImpl struct {
	// A mapping of resource IDs to their associated [ResourceChat] objects.
	Resources *map[string]*ResourceChat `pkl:"resources"`
}

// A mapping of resource IDs to their associated [ResourceChat] objects.
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
