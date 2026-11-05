package clusterapi

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

const componentsDir = "../../data/data/cluster-api"

// controllerSets lists the infrastructure providers whose manifests are installed
// into the same local control plane, in the order runSystem installs them. An
// object defined by more than one manifest in a set is installed twice, with the
// last one winning.
var controllerSets = map[string][]string{
	"azure":      {"azure", "azureaso"},
	"azurestack": {"azurestack", "azureaso"},
}

// manifestObject carries just enough of an API object to identify it.
type manifestObject struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Group string `json:"group"`
	} `json:"spec"`
}

// readComponents parses the objects out of an infrastructure components manifest.
func readComponents(t *testing.T, provider string) []manifestObject {
	t.Helper()

	path := filepath.Join(componentsDir, fmt.Sprintf("%s-infrastructure-components.yaml", provider))
	f, err := os.Open(path)
	require.NoError(t, err, "failed to open %s", path)
	defer f.Close()

	objects := []manifestObject{}
	reader := utilyaml.NewYAMLReader(bufio.NewReaderSize(f, 1024*1024))
	for {
		doc, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err, "failed to read a document from %s", path)

		obj := manifestObject{}
		require.NoError(t, yaml.Unmarshal(doc, &obj), "failed to unmarshal a document from %s", path)
		if obj.Kind == "" {
			continue
		}
		objects = append(objects, obj)
	}
	return objects
}

// TestNoDuplicateCRDs verifies that no CRD is defined by two manifests installed
// into the same local control plane. envtest updates a CRD that already exists,
// and the apiserver rejects an update that drops a version the first manifest
// recorded in status.storedVersions.
func TestNoDuplicateCRDs(t *testing.T) {
	for platform, providers := range controllerSets {
		t.Run(platform, func(t *testing.T) {
			definedBy := map[string]string{}
			for _, provider := range providers {
				for _, obj := range readComponents(t, provider) {
					if obj.Kind != "CustomResourceDefinition" {
						continue
					}
					if first, ok := definedBy[obj.Metadata.Name]; ok {
						assert.Failf(t, "duplicate CRD",
							"CRD %q is defined by both the %s and %s manifests", obj.Metadata.Name, first, provider)
						continue
					}
					definedBy[obj.Metadata.Name] = provider
				}
			}
		})
	}
}

// TestASOWebhooksAreOverwritten verifies that every Azure Service Operator webhook
// configuration shipped by a CAPZ manifest is also shipped by the azureaso
// manifest. envtest rewrites each manifest's webhooks to point at the controller
// it belongs to, so an ASO webhook the azureaso manifest does not overwrite is
// left pointing at the CAPZ webhook server, which does not serve ASO paths.
func TestASOWebhooksAreOverwritten(t *testing.T) {
	asoWebhooks := func(provider string) map[string]bool {
		names := map[string]bool{}
		for _, obj := range readComponents(t, provider) {
			switch obj.Kind {
			case "MutatingWebhookConfiguration", "ValidatingWebhookConfiguration":
			default:
				continue
			}
			if strings.HasPrefix(obj.Metadata.Name, "azureserviceoperator-") {
				names[obj.Metadata.Name] = true
			}
		}
		return names
	}

	owned := asoWebhooks("azureaso")
	require.NotEmpty(t, owned, "the azureaso manifest defines no ASO webhook configurations")

	for platform, providers := range controllerSets {
		t.Run(platform, func(t *testing.T) {
			for _, provider := range providers {
				if provider == "azureaso" {
					continue
				}
				for name := range asoWebhooks(provider) {
					assert.Truef(t, owned[name],
						"the %s manifest defines ASO webhook configuration %q, which the azureaso manifest does not overwrite", provider, name)
				}
			}
		})
	}
}
