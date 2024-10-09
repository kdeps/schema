// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

type ResourceChat struct {
	Model string `pkl:"model"`

	Prompt string `pkl:"prompt"`

	JsonResponse *bool `pkl:"jsonResponse"`

	Response *string `pkl:"response"`

	Timestamp *uint32 `pkl:"timestamp"`

	TimeoutSeconds *int `pkl:"timeoutSeconds"`
}
