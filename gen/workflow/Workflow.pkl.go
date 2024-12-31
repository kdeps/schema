// Code generated from Pkl module `org.kdeps.pkl.Workflow`. DO NOT EDIT.
package workflow

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/project"
)

type Workflow interface {
	GetName() string

	GetDescription() string

	GetWebsite() *string

	GetAuthors() *[]string

	GetDocumentation() *string

	GetRepository() *string

	GetHeroImage() *string

	GetAgentIcon() *string

	GetVersion() string

	GetAction() string

	GetWorkflows() []string

	GetSettings() *project.Settings
}

var _ Workflow = (*WorkflowImpl)(nil)

// Abstractions for Kdeps Workflow Management
//
// This module provides functionality for defining and managing workflows within the Kdeps system.
// It handles workflow validation, versioning, and linking to external actions, repositories, and
// documentation. Workflows are defined by a name, description, version, actions, and can reference
// external workflows and settings.
//
// This module also ensures the proper structure of workflows using validation checks for names,
// workflow references, action formats, and versioning patterns.
type WorkflowImpl struct {
	// The name of the workflow, validated to contain only alphanumeric characters.
	Name string `pkl:"name"`

	// A description of the workflow, providing details about its purpose and behavior.
	Description string `pkl:"description"`

	// A URI pointing to the website or landing page for the workflow, if available.
	Website *string `pkl:"website"`

	// A listing of the authors or contributors to the workflow.
	Authors *[]string `pkl:"authors"`

	// A URI pointing to the documentation for the workflow, if available.
	Documentation *string `pkl:"documentation"`

	// A URI pointing to the repository where the workflow's code or configuration can be found.
	Repository *string `pkl:"repository"`

	// Hero image to be used on this AI Agent.
	HeroImage *string `pkl:"heroImage"`

	// The icon to be used on this AI agent.
	AgentIcon *string `pkl:"agentIcon"`

	// The version of the workflow, following semantic versioning rules (e.g., 1.0.0).
	Version string `pkl:"version"`

	// The default action to be performed by the workflow, validated to ensure proper formatting.
	Action string `pkl:"action"`

	// A listing of external workflows referenced by this workflow, validated by format.
	Workflows []string `pkl:"workflows"`

	// The project settings that this workflow depends on.
	Settings *project.Settings `pkl:"settings"`
}

// The name of the workflow, validated to contain only alphanumeric characters.
func (rcv *WorkflowImpl) GetName() string {
	return rcv.Name
}

// A description of the workflow, providing details about its purpose and behavior.
func (rcv *WorkflowImpl) GetDescription() string {
	return rcv.Description
}

// A URI pointing to the website or landing page for the workflow, if available.
func (rcv *WorkflowImpl) GetWebsite() *string {
	return rcv.Website
}

// A listing of the authors or contributors to the workflow.
func (rcv *WorkflowImpl) GetAuthors() *[]string {
	return rcv.Authors
}

// A URI pointing to the documentation for the workflow, if available.
func (rcv *WorkflowImpl) GetDocumentation() *string {
	return rcv.Documentation
}

// A URI pointing to the repository where the workflow's code or configuration can be found.
func (rcv *WorkflowImpl) GetRepository() *string {
	return rcv.Repository
}

// Hero image to be used on this AI Agent.
func (rcv *WorkflowImpl) GetHeroImage() *string {
	return rcv.HeroImage
}

// The icon to be used on this AI agent.
func (rcv *WorkflowImpl) GetAgentIcon() *string {
	return rcv.AgentIcon
}

// The version of the workflow, following semantic versioning rules (e.g., 1.0.0).
func (rcv *WorkflowImpl) GetVersion() string {
	return rcv.Version
}

// The default action to be performed by the workflow, validated to ensure proper formatting.
func (rcv *WorkflowImpl) GetAction() string {
	return rcv.Action
}

// A listing of external workflows referenced by this workflow, validated by format.
func (rcv *WorkflowImpl) GetWorkflows() []string {
	return rcv.Workflows
}

// The project settings that this workflow depends on.
func (rcv *WorkflowImpl) GetSettings() *project.Settings {
	return rcv.Settings
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Workflow
func LoadFromPath(ctx context.Context, path string) (ret Workflow, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Workflow
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Workflow, error) {
	var ret WorkflowImpl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
