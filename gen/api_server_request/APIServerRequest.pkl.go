// Code generated from Pkl module `org.kdeps.pkl.APIServerRequest`. DO NOT EDIT.
package apiserverrequest

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type APIServerRequest interface {
	GetPath() string

	GetMethod() string

	GetData() *string

	GetParams() *map[string]string

	GetHeaders() *map[string]string

	GetFiles() *map[string]*APIServerRequestUploads
}

var _ APIServerRequest = (*APIServerRequestImpl)(nil)

type APIServerRequestImpl struct {
	Path string `pkl:"path"`

	Method string `pkl:"method"`

	Data *string `pkl:"data"`

	Params *map[string]string `pkl:"params"`

	Headers *map[string]string `pkl:"headers"`

	Files *map[string]*APIServerRequestUploads `pkl:"files"`
}

func (rcv *APIServerRequestImpl) GetPath() string {
	return rcv.Path
}

func (rcv *APIServerRequestImpl) GetMethod() string {
	return rcv.Method
}

func (rcv *APIServerRequestImpl) GetData() *string {
	return rcv.Data
}

func (rcv *APIServerRequestImpl) GetParams() *map[string]string {
	return rcv.Params
}

func (rcv *APIServerRequestImpl) GetHeaders() *map[string]string {
	return rcv.Headers
}

func (rcv *APIServerRequestImpl) GetFiles() *map[string]*APIServerRequestUploads {
	return rcv.Files
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a APIServerRequest
func LoadFromPath(ctx context.Context, path string) (ret APIServerRequest, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a APIServerRequest
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (APIServerRequest, error) {
	var ret APIServerRequestImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
