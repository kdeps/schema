// Code generated from Pkl module `org.kdeps.pkl.LLM`. DO NOT EDIT.
package llm

// Class representing the details of a chat interaction with an LLM model, including prompts, responses,
// file generation, and additional metadata.
type ResourceChat struct {
	// The name of the LLM model used for the chat.
	Model string `pkl:"model"`

	// The prompt text sent to the LLM model.
	Prompt string `pkl:"prompt"`

	// A listing of file paths or identifiers associated with the chat.
	Files *[]string `pkl:"files"`

	// Whether the chat involves image generation. Defaults to `false`.
	ImageGeneration *bool `pkl:"imageGeneration"`

	// The file path of a generated file from the LLM interaction.
	GeneratedFile *string `pkl:"generatedFile"`

	// Whether the LLM's response is in JSON format. Defaults to `false`.
	JsonResponse *bool `pkl:"jsonResponse"`

	// A listing of keys expected in the JSON response from the LLM model.
	JsonResponseKeys *[]string `pkl:"jsonResponseKeys"`

	// The actual response returned from the LLM model.
	Response *string `pkl:"response"`

	// A timestamp of when the response was generated, represented as an unsigned 32-bit integer.
	Timestamp *uint32 `pkl:"timestamp"`

	// The timeout duration (in seconds) for the LLM interaction. Defaults to 60 seconds.
	TimeoutSeconds *int `pkl:"timeoutSeconds"`
}
