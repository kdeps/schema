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

type LLMImpl struct {
	Resources *map[string]*ResourceChat `pkl:"resources"`
}

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
