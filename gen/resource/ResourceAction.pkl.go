// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"github.com/kdeps/schema/gen/api_server_response"
	"github.com/kdeps/schema/gen/exec"
	"github.com/kdeps/schema/gen/http"
	"github.com/kdeps/schema/gen/llm"
)

// Class representing an action that can be executed on a resource.
type ResourceAction struct {
	// Configuration for executing commands.
	Exec *exec.ResourceExec `pkl:"exec"`

	// Configuration for chat interactions with an LLM.
	Chat *llm.ResourceChat `pkl:"chat"`

	// A listing of conditions that determine if the action should be skipped.
	SkipCondition *[]bool `pkl:"skipCondition"`

	// A pre-flight validation check to be performed before executing the action.
	PreflightCheck *ValidationCheck `pkl:"preflightCheck"`

	// A post-flight validation check to be performed after executing the action.
	PostflightCheck *ValidationCheck `pkl:"postflightCheck"`

	// Configuration for HTTP client interactions.
	HttpClient *http.ResourceHTTPClient `pkl:"httpClient"`

	// Configuration for handling API responses.
	ApiResponse *apiserverresponse.APIServerResponse `pkl:"apiResponse"`
}
