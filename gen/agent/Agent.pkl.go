// Code generated from Pkl module `org.kdeps.pkl.Agent`. DO NOT EDIT.
package agent

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Agent interface {
}

var _ Agent = (*AgentImpl)(nil)

// Abstractions for Agent ID resolution
type AgentImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Agent
func LoadFromPath(ctx context.Context, path string) (ret Agent, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Agent
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Agent, error) {
	var ret AgentImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
