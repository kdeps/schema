// Code generated from Pkl module `org.kdeps.pkl.Tool`. DO NOT EDIT.
package tool

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Tool interface {
}

var _ Tool = ToolImpl{}

// Abstractions for Tool records
type ToolImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Tool
func LoadFromPath(ctx context.Context, path string) (ret Tool, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return ret, err
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Tool
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Tool, error) {
	var ret ToolImpl
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
