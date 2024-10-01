// Code generated from Pkl module `org.kdeps.pkl.Env`. DO NOT EDIT.
package env

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Env struct {
	Resource *map[string]*ResourceEnv `pkl:"resource"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Env
func LoadFromPath(ctx context.Context, path string) (ret *Env, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Env
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Env, error) {
	var ret Env
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
