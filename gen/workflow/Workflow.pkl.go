// Code generated from Pkl module `org.kdeps.pkl.Workflow`. DO NOT EDIT.
package workflow

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/project"
)

type Workflow interface {
	GetAgentID() string

	GetDescription() string

	GetWebsite() *string

	GetAuthors() *[]string

	GetDocumentation() *string

	GetRepository() *string

	GetHeroImage() *string

	GetAgentIcon() *string

	GetVersion() string

	GetTargetActionID() string

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
	AgentID string `pkl:"AgentID"`

	// A description of the workflow, providing details about its purpose and behavior.
	Description string `pkl:"Description"`

	// A URI pointing to the website or landing page for the workflow, if available.
	Website *string `pkl:"Website"`

	// A listing of the authors or contributors to the workflow.
	Authors *[]string `pkl:"Authors"`

	// A URI pointing to the documentation for the workflow, if available.
	Documentation *string `pkl:"Documentation"`

	// A URI pointing to the repository where the workflow's code or configuration can be found.
	Repository *string `pkl:"Repository"`

	// Hero image to be used on this AI Agent.
	HeroImage *string `pkl:"HeroImage"`

	// The icon to be used on this AI agent.
	AgentIcon *string `pkl:"AgentIcon"`

	// The version of the workflow, following semantic versioning rules (e.g., 1.0.0).
	Version string `pkl:"Version"`

	// The default action to be performed by the workflow, validated to ensure proper formatting.
	TargetActionID string `pkl:"TargetActionID"`

	// A listing of external workflows referenced by this workflow, validated by format.
	Workflows []string `pkl:"Workflows"`

	// The project settings that this workflow depends on.
	Settings *project.Settings `pkl:"Settings"`
}

// The name of the workflow, validated to contain only alphanumeric characters.
func (rcv *WorkflowImpl) GetAgentID() string {
	return rcv.AgentID
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
func (rcv *WorkflowImpl) GetTargetActionID() string {
	return rcv.TargetActionID
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
