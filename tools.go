// +build tools

// Official workaround to track tool dependencies with go modules:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	// dependency of hack/go-lint.sh
	_ "golang.org/x/lint"

	// dependency of generating CRD for install-config
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)
