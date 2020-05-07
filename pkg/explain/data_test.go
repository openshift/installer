package explain

import (
	"io/ioutil"
	"testing"
)

func loadCRD(t *testing.T) []byte {
	crd, err := ioutil.ReadFile("../../data/data/install.openshift.io_installconfigs.yaml")
	if err != nil {
		t.Fatalf("failed to load CRD: %v", err)
	}
	return crd
}
