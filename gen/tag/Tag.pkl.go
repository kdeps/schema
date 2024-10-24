// Code generated from Pkl module `org.kdeps.pkl.Tag`. DO NOT EDIT.
package tag

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Tag interface {
	GetTags() *map[string]string
}

var _ Tag = (*TagImpl)(nil)

// Abstractions for Kdeps Resource Tagging
//
// This module provides definitions and validations for tagging resources within the Kdeps framework.
// Tags are used to categorize and identify resources using alphanumeric names. Each tag is associated
// with a value and an optional timestamp, allowing for organized resource management and retrieval.
type TagImpl struct {
	// A mapping of tag names to their corresponding [str] objects.
	Tags *map[string]string `pkl:"tags"`
}

// A mapping of tag names to their corresponding [str] objects.
func (rcv *TagImpl) GetTags() *map[string]string {
	return rcv.Tags
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Tag
func LoadFromPath(ctx context.Context, path string) (ret Tag, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Tag
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Tag, error) {
	var ret TagImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
