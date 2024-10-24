// Code generated from Pkl module `org.kdeps.pkl.Security`. DO NOT EDIT.
package security

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Security interface {
}

var _ Security = (*SecurityImpl)(nil)

// Abstractions for Kdeps Security Settings
//
// This module defines the security settings for the Kdeps framework, allowing for the enforcement
// of security rules and the configuration of allowed headers, URLs, executable commands, external workflows,
// models, and PPA (Personal Package Archive). These settings are crucial for ensuring secure
// operations within the Kdeps environment.
type SecurityImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Security
func LoadFromPath(ctx context.Context, path string) (ret Security, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Security
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Security, error) {
	var ret SecurityImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
