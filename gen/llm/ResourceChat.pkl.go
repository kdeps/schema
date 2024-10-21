// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

type ResourceChat struct {
	Model string `pkl:"model"`

	Prompt string `pkl:"prompt"`

	Files *[]string `pkl:"files"`

	ImageGeneration *bool `pkl:"imageGeneration"`

	GeneratedFile *string `pkl:"generatedFile"`

	JsonResponse *bool `pkl:"jsonResponse"`

	JsonResponseKeys *[]string `pkl:"jsonResponseKeys"`

	Response *string `pkl:"response"`

	Timestamp *uint32 `pkl:"timestamp"`

	TimeoutSeconds *int `pkl:"timeoutSeconds"`
}
