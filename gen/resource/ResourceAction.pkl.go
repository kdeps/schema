// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

import (
	"github.com/kdeps/schema/gen/api_server_response"
	"github.com/kdeps/schema/gen/exec"
	"github.com/kdeps/schema/gen/http"
	"github.com/kdeps/schema/gen/llm"
)

type ResourceAction struct {
	Exec *exec.ResourceExec `pkl:"exec"`

	Chat *llm.ResourceChat `pkl:"chat"`

	SkipCondition *[]bool `pkl:"skipCondition"`

	PreflightCheck *ValidationCheck `pkl:"preflightCheck"`

	PostflightCheck *ValidationCheck `pkl:"postflightCheck"`

	HttpClient *http.ResourceHTTPClient `pkl:"httpClient"`

	ApiResponse *apiserverresponse.APIServerResponse `pkl:"apiResponse"`
}
