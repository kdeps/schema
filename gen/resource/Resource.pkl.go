// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

// Abstractions for Kdeps Resources
//
// This module defines the structure for resources used within the Kdeps framework,
// including actions that can be performed on these resources, validation checks,
// and error handling mechanisms. Each resource can define its actionID, name, description,
// category, dependencies, and how it runs.
type Resource struct {
	// The unique identifier for the resource, validated against [isValidActionID].
	ActionID string `pkl:"actionID"`

	// The name of the resource.
	Name string `pkl:"name"`

	// A description of the resource, providing additional context.
	Description string `pkl:"description"`

	// The category to which the resource belongs.
	Category string `pkl:"category"`

	// A listing of dependencies required by the resource, validated against [isValidDependency].
	Requires *[]string `pkl:"requires"`

	// Defines the action items to be processed individually in a loop.
	Items *[]string `pkl:"items"`

	// Defines the action to be taken for the resource.
	Run *ResourceAction `pkl:"run"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Resource
func LoadFromPath(ctx context.Context, path string) (ret *Resource, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Resource
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Resource, error) {
	var ret Resource
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
