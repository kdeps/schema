package test

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

// NewTestEvaluator returns a PKL evaluator with all modules/resources allowed and custom readers.
func NewTestEvaluator(readers ...pkl.ResourceReader) (pkl.Evaluator, error) {
	opts := func(options *pkl.EvaluatorOptions) {
		pkl.WithDefaultAllowedResources(options)
		pkl.WithOsEnv(options)
		pkl.WithDefaultAllowedModules(options)
		pkl.WithDefaultCacheDir(options)
		options.Logger = pkl.NoopLogger
		options.ResourceReaders = readers
		options.AllowedModules = []string{".*"}
		options.AllowedResources = []string{".*"}
		options.ModulePaths = []string{"."}
	}
	return pkl.NewEvaluator(context.Background(), opts)
}
