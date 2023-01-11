package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"openshift-install": func() int {
			main()
			return 0
		},
	}))
}

func TestAgentIntegration(t *testing.T) {
	runAllIntegrationTests(t, "testdata/agent")
}
