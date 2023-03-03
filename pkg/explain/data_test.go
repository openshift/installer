package explain

import (
	"os"
	"testing"
)

func loadCRD(t *testing.T) []byte {
	crd, err := os.ReadFile("../../data/data/installconfig/install.openshift.io_installconfigs.yaml")
	if err != nil {
		t.Fatalf("failed to load CRD: %v", err)
	}
	return crd
}
