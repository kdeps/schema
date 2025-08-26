// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/gen/api_server_response"
	"github.com/kdeps/schema/gen/exec"
	"github.com/kdeps/schema/gen/http"
	"github.com/kdeps/schema/gen/llm"
	"github.com/kdeps/schema/gen/python"
)

// Class representing an action that can be executed on a resource.
type ResourceAction struct {
	// Block for performing PKL expressions.
	Expr *pkl.Object `pkl:"Expr"`

	// Configuration for executing commands.
	Exec *exec.ResourceExec `pkl:"Exec"`

	// Configuration for python scripts.
	Python *python.ResourcePython `pkl:"Python"`

	// Configuration for chat interactions with an LLM.
	Chat *llm.ResourceChat `pkl:"Chat"`

	// A listing of conditions that determine if the action should be skipped.
	SkipCondition *[]any `pkl:"SkipCondition"`

	// A pre-flight validation check to be performed before executing the action.
	PreflightCheck *ValidationCheck `pkl:"PreflightCheck"`

	// A listing of allowed HTTP headers
	AllowedHeaders *[]string `pkl:"AllowedHeaders"`

	// A listing of allowed HTTP params
	AllowedParams *[]string `pkl:"AllowedParams"`

	// A listing of targeted HTTP methods
	RestrictToHTTPMethods *[]string `pkl:"RestrictToHTTPMethods"`

	// A listing of targeted HTTP routes
	RestrictToRoutes *[]string `pkl:"RestrictToRoutes"`

	// Configuration for HTTP client interactions.
	HTTPClient *http.ResourceHTTPClient `pkl:"HTTPClient"`

	// Configuration for handling API responses.
	APIResponse *apiserverresponse.APIServerResponse `pkl:"APIResponse"`
}
