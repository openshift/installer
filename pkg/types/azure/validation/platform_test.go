package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

func validPlatform() *azure.Platform {
	return &azure.Platform{
		Region:                      "eastus",
		BaseDomainResourceGroupName: "group",
		OutboundType:                azure.LoadbalancerOutboundType,
		CloudName:                   azure.PublicCloud,
	}
}

func validNetworkPlatform() *azure.Platform {
	p := validPlatform()
	p.NetworkResourceGroupName = "networkresourcegroup"
	p.VirtualNetwork = "virtualnetwork"
	p.Subnets = []azure.SubnetSpec{
		{
			Name: "controlplanesubnet",
			Role: v1beta1.SubnetControlPlane,
		},
		{
			Name: "computesubnet",
			Role: v1beta1.SubnetNode,
		},
	}

	return p
}

func validCustomerManagedKeys() *azure.Platform {
	p := validPlatform()
	p.CustomerManagedKey = &azure.CustomerManagedKey{
		KeyVault: azure.KeyVault{
			KeyName:       "test",
			Name:          "test",
			ResourceGroup: "test123",
		},
		UserAssignedIdentityKey: "12345678-1234-1234-1234-123456789123"}
	return p
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *azure.Platform
		wantSkip func(p *azure.Platform) bool
		expected string
	}{
		{
			name: "invalid region",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.Region = ""
				return p
			}(),
			expected: `^test-path\.region: Required value: region should be set to one of the supported Azure regions$`,
		},
		{
			name: "invalid baseDomainResourceGroupName",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.BaseDomainResourceGroupName = ""
				return p
			}(),
			expected: `^test-path\.baseDomainResourceGroupName: Required value: baseDomainResourceGroupName is the resource group name where the azure dns zone is deployed$`,
		},
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "valid machine pool",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &azure.MachinePool{}
				return p
			}(),
		},
		{
			name:     "valid subnets & virtual network",
			platform: validNetworkPlatform(),
		},
		{
			name: "missing subnets",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.Subnets = p.Subnets[1:]
				return p
			}(),
			expected: `^test-path\.controlPlaneSubnet: Required value: must provide a control plane subnet when a virtual network is specified$`,
		},
		{
			name: "subnets missing virtual network",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.Subnets = p.Subnets[0:1]
				p.VirtualNetwork = ""
				return p
			}(),
			expected: `^test-path\.virtualNetwork: Required value: must provide a virtual network when supplying subnets$`,
		},
		{
			name: "missing network resource group",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.NetworkResourceGroupName = ""
				return p
			}(),
			expected: `^\[test-path\.networkResourceGroupName: Required value: must provide a network resource group when a virtual network is specified, test-path\.networkResourceGroupName: Required value: must provide a network resource group when supplying subnets\]$`,
		},
		{
			name: "missing cloud name",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.CloudName = ""
				return p
			}(),
			expected: `^test-path\.cloudName: Unsupported value: "": supported values:`,
		},
		{
			name: "invalid cloud name",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.CloudName = azure.CloudEnvironment("AzureOtherCloud")
				return p
			}(),
			expected: `^test-path\.cloudName: Unsupported value: "AzureOtherCloud": supported values:`,
		},
		{
			name: "invalid outbound type",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.OutboundType = "random-egress"
				return p
			}(),
			expected: `^test-path\.outboundType: Unsupported value: "random-egress": supported values: "Loadbalancer", "NATGatewayMultiZone", "NATGatewaySingleZone", "UserDefinedRouting"$`,
		},
		{
			name: "invalid user defined type",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.OutboundType = azure.UserDefinedRoutingOutboundType
				return p
			}(),
			expected: `^test-path\.outboundType: Invalid value: "UserDefinedRouting": UserDefinedRouting is only allowed when installing to pre-existing network$`,
		},
		{
			name: "missing key vault name",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.KeyVault.Name = ""
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Required value: name of the key vault is required for storage account encryption$`,
		},
		{
			name: "invalid key vault name",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.KeyVault.Name = "1invalid"
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Invalid value: "1invalid": invalid name for key vault for encryption$`,
		},
		{
			name: "missing key vault key name",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.KeyVault.KeyName = ""
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Required value: key vault key name is required for storage account encryption$`,
		},
		{
			name: "invalid key vault key name",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.KeyVault.KeyName = "."
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Invalid value: ".": invalid key name for encryption$`,
		},
		{
			name: "missing resource group",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.KeyVault.ResourceGroup = ""
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Required value: resource group of the key vault is required for storage account encryption$`,
		},
		{
			name: "invalid resource group",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.KeyVault.ResourceGroup = "invalid."
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Invalid value: "invalid.": invalid resource group for encryption$`,
		},
		{
			name: "missing user assigned identity",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.UserAssignedIdentityKey = ""
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Required value: user assigned identity key is required for storage account encryption$`,
		},
		{
			name: "invalid user assigned identity",
			platform: func() *azure.Platform {
				p := validCustomerManagedKeys()
				p.CustomerManagedKey.UserAssignedIdentityKey = "-"
				return p
			}(),
			expected: `^test-path\.customerManagedKey: Invalid value: "-": invalid user assigned identity key for encryption$`,
		},
	}
	ic := types.InstallConfig{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantSkip != nil && tc.wantSkip(tc.platform) {
				t.Skip()
			}

			err := ValidatePlatform(tc.platform, types.ExternalPublishingStrategy, field.NewPath("test-path"), &ic).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

func TestValidateUserTags(t *testing.T) {
	fieldPath := "spec.platform.azure.userTags"
	cases := []struct {
		name     string
		userTags map[string]string
		wantErr  bool
	}{
		{
			name:     "userTags not configured",
			userTags: map[string]string{},
			wantErr:  false,
		},
		{
			name: "userTags configured",
			userTags: map[string]string{
				"key1": "value1", "key_2": "value_2", "key.3": "value.3", "key-4._": "value=4", "key.5_A": "value+5",
				"key-6": "value-6", "Key.-_7": "value@7", "key8_": "value8-", "key9A": "value9+", "key10a": "value10@"},
			wantErr: false,
		},
		{
			name: "userTags configured is more than max limit",
			userTags: map[string]string{
				"key1": "value1", "key2": "value2", "key3": "value3", "key4": "value4", "key5": "value5",
				"key6": "value6", "key7": "value7", "key8": "value8", "key9": "value9", "key10": "value10",
				"key11": "value11"},
			wantErr: true,
		},
		{
			name:     "userTags contains key starting with a number",
			userTags: map[string]string{"1key": "1value"},
			wantErr:  true,
		},
		{
			name:     "userTags contains key starting with a special character",
			userTags: map[string]string{"_key": "1value"},
			wantErr:  true,
		},
		{
			name:     "userTags contains key ending with a special character",
			userTags: map[string]string{"key@": "1value"},
			wantErr:  true,
		},
		{
			name:     "userTags contains empty key",
			userTags: map[string]string{"": "value"},
			wantErr:  true,
		},
		{
			name: "userTags contains key length greater than 128",
			userTags: map[string]string{
				"thisisaverylongkeywithmorethan128characterswhichisnotallowedforazureresourcetagkeysandthetagkeyvalidationshouldfailwithinvalidfieldvalueerror": "value"},
			wantErr: true,
		},
		{
			name:     "userTags contains key with invalid character",
			userTags: map[string]string{"key/test": "value"},
			wantErr:  true,
		},
		{
			name:     "userTags contains value length greater than 256",
			userTags: map[string]string{"key": "thisisaverylongvaluewithmorethan256characterswhichisnotallowedforazureresourcetagvaluesandthetagvaluevalidationshouldfailwithinvalidfieldvalueerrorrepeatthisisaverylongvaluewithmorethan256characterswhichisnotallowedforazureresourcetagvaluesandthetagvaluevalidationshouldfailwithinvalidfieldvalueerror"},
			wantErr:  true,
		},
		{
			name:     "userTags contains empty value",
			userTags: map[string]string{"key": ""},
			wantErr:  true,
		},
		{
			name:     "userTags contains value with invalid character",
			userTags: map[string]string{"key": "value*^%"},
			wantErr:  true,
		},
		{
			name:     "userTags contains key as name",
			userTags: map[string]string{"name": "value"},
			wantErr:  true,
		},
		{
			name:     "userTags contains allowed key name123",
			userTags: map[string]string{"name123": "value"},
			wantErr:  false,
		},
		{
			name:     "userTags contains key with prefix kubernetes.io",
			userTags: map[string]string{"kubernetes.io_cluster": "value"},
			wantErr:  true,
		},
		{
			name:     "userTags contains allowed key prefix for_openshift.io",
			userTags: map[string]string{"for_openshift.io": "azure"},
			wantErr:  false,
		},
		{
			name:     "userTags contains key with prefix azure",
			userTags: map[string]string{"azure": "microsoft"},
			wantErr:  true,
		},
		{
			name:     "userTags contains allowed key resourcename",
			userTags: map[string]string{"resourcename": "value"},
			wantErr:  false,
		},
		{
			name: "userTags contain duplicate keys",
			userTags: map[string]string{
				"environment": "test",
				"Environment": "lab",
				"key":         "value",
				"createdFor":  "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserTags(tt.userTags, field.NewPath(fieldPath))
			if (len(err) > 0) != tt.wantErr {
				t.Errorf("unexpected error, err: %v", err)
			}
		})
	}
}
