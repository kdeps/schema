// Code generated from Pkl module `org.kdeps.pkl.Resource`. DO NOT EDIT.
package resource

// Class representing validation checks that can be performed on actions.
type ValidationCheck struct {
	// A listing of validation conditions.
	Validations *[]any `pkl:"Validations"`

	// An error associated with the validation check, if any.
	Error *APIError `pkl:"Error"`

	// Boolean flag to enable or disable retry functionality for the validation check.
	//
	// - `true`: The validation check will be retried if it fails.
	// - `false`: The validation check will not be retried. Default is `false`.
	Retry *bool `pkl:"Retry"`

	// The number of times to retry the validation check before considering it a failure.
	//
	// This property is only used when [Retry] is set to `true`.
	// Default value is 3 retry attempts.
	RetryTimes *int `pkl:"RetryTimes"`
}
