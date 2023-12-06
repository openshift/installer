package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/azure"
)

func TestCloudProviderConfig(t *testing.T) {

	config := CloudProviderConfig{
		CloudName:                azure.PublicCloud,
		ResourceGroupName:        "clusterid-rg",
		GroupLocation:            "westeurope",
		ResourcePrefix:           "clusterid",
		SubscriptionID:           "subID",
		TenantID:                 "tenantID",
		NetworkResourceGroupName: "clusterid-rg",
		NetworkSecurityGroupName: "clusterid-node-nsg",
		VirtualNetworkName:       "clusterid-vnet",
		SubnetName:               "clusterid-worker-subnet",
	}
	expected := `{
	"cloud": "AzurePublicCloud",
	"tenantId": "tenantID",
	"aadClientId": "",
	"aadClientSecret": "",
	"aadClientCertPath": "",
	"aadClientCertPassword": "",
	"useManagedIdentityExtension": true,
	"userAssignedIdentityID": "",
	"subscriptionId": "subID",
	"resourceGroup": "clusterid-rg",
	"location": "westeurope",
	"vnetName": "clusterid-vnet",
	"vnetResourceGroup": "clusterid-rg",
	"subnetName": "clusterid-worker-subnet",
	"securityGroupName": "clusterid-node-nsg",
	"routeTableName": "clusterid-node-routetable",
	"vmType": "standard",
	"loadBalancerSku": "standard",
	"cloudProviderBackoff": true,
	"useInstanceMetadata": true,
	"excludeMasterFromStandardLB": false,
	"cloudProviderBackoffDuration": 6,
	"putVMSSVMBatchSize": 0,
	"enableMigrateToIPBasedBackendPoolAPI": false
}
`

	json, err := config.JSON()
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expected, json, "unexpected cloud provider config")
}
