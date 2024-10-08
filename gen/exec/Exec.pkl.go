// Code generated from Pkl module `org.kdeps.pkl.Exec`. DO NOT EDIT.
package exec

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Exec interface {
	GetResources() *map[string]*ResourceExec
}

var _ Exec = (*ExecImpl)(nil)

type ExecImpl struct {
	Resources *map[string]*ResourceExec `pkl:"resources"`
}

func (rcv *ExecImpl) GetResources() *map[string]*ResourceExec {
	return rcv.Resources
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Exec
func LoadFromPath(ctx context.Context, path string) (ret Exec, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Exec
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Exec, error) {
	var ret ExecImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
