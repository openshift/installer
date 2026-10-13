package logforwarding

import (
	"testing"

	//nolint:staticcheck
	. "github.com/onsi/ginkgo/v2"
	//nolint:staticcheck
	. "github.com/onsi/gomega"
)

func TestLogforwarding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logforwarding Suite")
}
