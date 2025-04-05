package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfig(t *testing.T) {
	expectedConfig := `[global]
project-id      = test-project-id
regional        = true
multizone       = true
node-tags       = uid-master
node-tags       = uid-control-plane
node-tags       = uid-worker
node-instance-prefix = uid
external-instance-groups-prefix = uid
subnetwork-name = uid-worker-subnet

`
	actualConfig, err := CloudProviderConfig("uid", "test-project-id", "uid-worker-subnet", "", "", "")
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}

func TestCloudProviderConfigWithNPID(t *testing.T) {
	expectedConfig := `[global]
project-id      = test-project-id
regional        = true
multizone       = true
node-tags       = uid-master
node-tags       = uid-control-plane
node-tags       = uid-worker
node-instance-prefix = uid
external-instance-groups-prefix = uid
subnetwork-name = uid-worker-subnet
network-project-id = test-network-project-id

`
	actualConfig, err := CloudProviderConfig("uid", "test-project-id", "uid-worker-subnet", "test-network-project-id", "", "")
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}

func TestCloudProviderConfigWithEndpoints(t *testing.T) {
	expectedConfig := `[global]
project-id      = test-project-id
regional        = true
multizone       = true
node-tags       = uid-master
node-tags       = uid-control-plane
node-tags       = uid-worker
node-instance-prefix = uid
external-instance-groups-prefix = uid
subnetwork-name = uid-worker-subnet
api-endpoint = compute-testendpoint.p.googleapis.com
container-api-endpoint = container-testendpoint.p.googleapis.com

`
	actualConfig, err := CloudProviderConfig("uid", "test-project-id", "uid-worker-subnet", "", "compute-testendpoint.p.googleapis.com", "container-testendpoint.p.googleapis.com")
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}
