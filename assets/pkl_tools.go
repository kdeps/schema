//go:build tools
// +build tools

package assets

// Blank-import so the assets package (and its embedded files) stay
// in the build graph even if nothing uses it directly.
import _ "github.com/kdeps/schema/assets"
