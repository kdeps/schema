// Code generated from Pkl module `org.kdeps.pkl.Security`. DO NOT EDIT.
package security

// Class representing the security settings for the Kdeps framework.
type Settings struct {
	// A flag indicating whether to enforce security rules. Defaults to `false`.
	EnforceSecurityRules *bool `pkl:"enforceSecurityRules"`

	// A listing of allowed HTTP server headers.
	AllowedHttpServerHeaders *[]string `pkl:"allowedHttpServerHeaders"`

	// A listing of allowed HTTP client headers.
	AllowedHttpClientHeaders *[]string `pkl:"allowedHttpClientHeaders"`

	// A listing of allowed HTTP client URLs.
	AllowedHttpClientUrl *[]string `pkl:"allowedHttpClientUrl"`

	// A listing of allowed executable commands.
	AllowedCmdExecutable *[]string `pkl:"allowedCmdExecutable"`

	// A listing of allowed external workflows.
	AllowedExternalWorkflow *[]string `pkl:"allowedExternalWorkflow"`

	// A listing of allowed models.
	AllowedModels *[]string `pkl:"allowedModels"`

	// A listing of allowed Personal Package Archives (PPA).
	AllowedPPA *[]string `pkl:"allowedPPA"`
}
