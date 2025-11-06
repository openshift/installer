package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{ //nolint:staticcheck //TODO OCPBUGS-64696
		"openshift-install": func() int {
			main()
			return 0
		},
	}))
}

func TestAgentIntegration(t *testing.T) {
	runAllIntegrationTests(t, "testdata/agent")
}

func TestImageBasedIntegration(t *testing.T) {
	runAllIntegrationTests(t, "testdata/imagebased")
}
