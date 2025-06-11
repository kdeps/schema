// Package assets makes the PKL schema available to downstream code/tests.
package assets

import "embed"

//go:embed ../deps/pkl/*.pkl
var PKLFS embed.FS
