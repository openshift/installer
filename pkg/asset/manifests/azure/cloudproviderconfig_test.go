package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfig(t *testing.T) {

	config := CloudProviderConfig{
		GroupLocation:            "westeurope",
		ResourcePrefix:           "clusterid",
		SubscriptionID:           "subID",
		TenantID:                 "tenantID",
		NetworkResourceGroupName: "clusterid-rg",
		NetworkSecurityGroupName: "clusterid-node-nsg",
		VirtualNetworkName:       "clusterid-vnet",
		SubnetName:               "clusterid-node-subnet",
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
	"subnetName": "clusterid-node-subnet",
	"securityGroupName": "clusterid-node-nsg",
	"routeTableName": "clusterid-node-routetable",
	"primaryAvailabilitySetName": "",
	"vmType": "",
	"primaryScaleSetName": "",
	"cloudProviderBackoff": true,
	"cloudProviderBackoffRetries": 0,
	"cloudProviderBackoffExponent": 0,
	"cloudProviderBackoffDuration": 6,
	"cloudProviderBackoffJitter": 0,
	"cloudProviderRateLimit": true,
	"cloudProviderRateLimitQPS": 6,
	"cloudProviderRateLimitBucket": 10,
	"cloudProviderRateLimitQPSWrite": 6,
	"cloudProviderRateLimitBucketWrite": 10,
	"useInstanceMetadata": true,
	"loadBalancerSku": "standard",
	"excludeMasterFromStandardLB": null,
	"disableOutboundSNAT": null,
	"maximumLoadBalancerRuleCount": 0
}
`

	json, err := config.JSON()
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expected, json, "unexpected cloud provider config")
}
