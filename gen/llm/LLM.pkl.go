// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type LLM interface {
}

var _ LLM = (*LLMImpl)(nil)

// Abstractions for Kdeps LLM operations
//
// This module provides the structure for LLM (Large Language Model) operations within the Kdeps framework,
// including chat interactions, response handling, and model configuration. It defines classes and functions
// for managing LLM resources, processing prompts, and handling responses from various LLM models.
//
// This module is part of the `kdeps` schema and provides a unified interface for LLM operations across
// different models and providers.
//
// The module defines:
// - [ResourceChat]: For managing chat interactions with LLM models.
// - [MultiChat]: For managing multi-turn chat conversations.
// - [Tool]: For managing tool interactions with LLM models.
// - Functions for retrieving and processing LLM responses.
type LLMImpl struct {
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
