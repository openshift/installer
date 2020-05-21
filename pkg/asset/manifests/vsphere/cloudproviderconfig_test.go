package vsphere

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func TestCloudProviderConfig(t *testing.T) {
	platform := &vspheretypes.Platform{
		VCenter:          "test-name",
		Username:         "test-username",
		Password:         "test-password",
		Datacenter:       "test-datacenter",
		DefaultDatastore: "test-datastore",
	}
	expectedConfig := `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[Workspace]
server = "test-name"
datacenter = "test-datacenter"
default-datastore = "test-datastore"
folder = "/test-datacenter/vm/clusterID"

[VirtualCenter "test-name"]
datacenters = "test-datacenter"
`
	folderPath := fmt.Sprintf("/%s/vm/%s", "test-datacenter", "clusterID")
	actualConfig, err := CloudProviderConfig(folderPath, platform)
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}
