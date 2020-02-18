package azure

import (
	"bytes"
	"encoding/json"
)

//CloudProviderConfig is the azure cloud provider config
type CloudProviderConfig struct {
	TenantID                    string
	SubscriptionID              string
	GroupLocation               string
	ResourcePrefix              string
	NetworkResourceGroupName    string
	NetworkSecurityGroupName    string
	VirtualNetworkName          string
	SubnetName                  string
	ExcludeMasterFromStandardLB bool `json:"excludeMasterFromStandardLB"`
}

// JSON generates the cloud provider json config for the azure platform.
// managed resource names are matching the convention defined by capz
func (params CloudProviderConfig) JSON() (string, error) {
	resourceGroupName := params.ResourcePrefix + "-rg"
	config := config{
		authConfig: authConfig{
			Cloud:                       "AzurePublicCloud",
			TenantID:                    params.TenantID,
			SubscriptionID:              params.SubscriptionID,
			UseManagedIdentityExtension: true,
			// The cloud provider needs the clientID which is only known after terraform has run.
			// When left empty, the existing managed identity on the VM will be used.
			// By leaving it empty, we don't have to create the identity before running the installer.
			// We only need to know that there will be one assigned to the VM, and we control this.
			// ref: https://github.com/kubernetes/kubernetes/blob/4b7c607ba47928a7be77fadef1550d6498397a4c/staging/src/k8s.io/legacy-cloud-providers/azure/auth/azure_auth.go#L69
			UserAssignedIdentityID: "",
		},
		ResourceGroup:          resourceGroupName,
		Location:               params.GroupLocation,
		SubnetName:             params.SubnetName,
		SecurityGroupName:      params.NetworkSecurityGroupName,
		VnetName:               params.VirtualNetworkName,
		VnetResourceGroup:      params.NetworkResourceGroupName,
		RouteTableName:         params.ResourcePrefix + "-node-routetable",
		CloudProviderBackoff:   true,
		CloudProviderRateLimit: true,

		// The default rate limits for Azure cloud provider are https://github.com/kubernetes/kubernetes/blob/f8d2b6b982bb06fc64979ac53ae668284d9c003c/staging/src/k8s.io/legacy-cloud-providers/azure/azure.go#L51-L56
		// While the AKS recommends following rate limits for large clusters https://github.com/Azure/aks-engine/blob/0f6aa91fa1870d5be657c62374d11f7d6009121d/examples/largeclusters/kubernetes.json#L9-L15
		// 									default		AKS (large)	Change
		// cloudProviderBackoffRetries		6			6					NO
		// cloudProviderBackoffJitter		1.0			1					NO
		// cloudProviderBackoffExponent		1.5			1.5					NO
		// cloudProviderBackoffDuration		5			6					YES to 6
		// cloudProviderRateLimitQPS		3			3					YES to 6
		// cloudProviderRateLimitBucket		5			10					YES to 10
		CloudProviderBackoffDuration:      6,
		CloudProviderRateLimitQPS:         6,
		CloudProviderRateLimitQPSWrite:    6,
		CloudProviderRateLimitBucket:      10,
		CloudProviderRateLimitBucketWrite: 10,

		UseInstanceMetadata: true,
		//default to standard load balancer, supports tcp resets on idle
		//https://docs.microsoft.com/en-us/azure/load-balancer/load-balancer-tcp-reset
		LoadBalancerSku: "standard",
	}
	buff := &bytes.Buffer{}
	encoder := json.NewEncoder(buff)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(config); err != nil {
		return "", err
	}
	return buff.String(), nil
}
