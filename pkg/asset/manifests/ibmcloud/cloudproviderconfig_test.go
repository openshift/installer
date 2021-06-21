package ibmcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfig(t *testing.T) {
	expectedConfig := `[global]
version = 1.1.0
[kubernetes]
config-file = /mnt/etc/kubernetes/controller-manager-kubeconfig
[load-balancer-deployment]
image = [REGISTRY]/[NAMESPACE]/keepalived:[TAG]
application = keepalived
vlan-ip-config-map = ibm-cloud-provider-vlan-ip-config
[provider]
accountID = 1e1f75646aef447814a6d907cc83fb3c
clusterID = ocp4-8pxks

`

	actualConfig, err := CloudProviderConfig("ocp4-8pxks", "1e1f75646aef447814a6d907cc83fb3c")
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expectedConfig, actualConfig, "unexpected cloud provider config")
}
