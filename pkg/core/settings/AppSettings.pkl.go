// Code generated from Pkl module `org.kdeps.pkl.Settings`. DO NOT EDIT.
package settings

import (
	"github.com/apple/pkl-go/pkl"
	"github.com/kdeps/schema/pkg/core/api"
	"github.com/kdeps/schema/pkg/core/llm"
)

type AppSettings struct {
	RunTimeout *pkl.Duration `pkl:"runTimeout"`

	InteractiveOnMissingValues bool `pkl:"interactiveOnMissingValues"`

	LlmSettings *llm.LLMSettings `pkl:"llmSettings"`

	ApiServerMode bool `pkl:"apiServerMode"`

	ApiServerSettings *api.APIServerSettings `pkl:"apiServerSettings"`
}
