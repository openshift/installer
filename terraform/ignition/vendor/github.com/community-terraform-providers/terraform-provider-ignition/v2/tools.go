// +build tools

// tools is a dummy package that will be ignored for builds, but included for dependencies.
package tools

import (
	_ "github.com/kisielk/errcheck"
	_ "github.com/mitchellh/gox"
)
