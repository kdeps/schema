// Code generated from Pkl module `org.kdeps.pkl.Docker`. DO NOT EDIT.
package docker

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Docker interface {
}

var _ Docker = (*DockerImpl)(nil)

// This module defines the settings and configurations for Docker-related
// resources within the KDEPS framework. It allows for the specification
// of package management, including additional package repositories (PPAs)
// and models to be used within Docker containers.
type DockerImpl struct {
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Docker
func LoadFromPath(ctx context.Context, path string) (ret Docker, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Docker
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Docker, error) {
	var ret DockerImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
