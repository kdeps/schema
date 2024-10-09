// Code generated from Pkl module `org.kdeps.pkl.Security`. DO NOT EDIT.
package security

type Settings struct {
	EnforceSecurityRules *bool `pkl:"enforceSecurityRules"`

	AllowedHttpServerHeaders *[]string `pkl:"allowedHttpServerHeaders"`

	AllowedHttpClientHeaders *[]string `pkl:"allowedHttpClientHeaders"`

	AllowedHttpClientUrl *[]string `pkl:"allowedHttpClientUrl"`

	AllowedCmdExecutable *[]string `pkl:"allowedCmdExecutable"`

	AllowedExternalWorkflow *[]string `pkl:"allowedExternalWorkflow"`

	AllowedModels *[]string `pkl:"allowedModels"`
}
