package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name            string
		platform        *gcp.Platform
		credentialsMode types.CredentialsMode
		publishStrategy types.PublishingStrategy
		valid           bool
	}{
		{
			name: "minimal",
			platform: &gcp.Platform{
				Region: "us-east1",
			},
			valid: true,
		},
		{
			name: "invalid region",
			platform: &gcp.Platform{
				Region: "",
			},
			valid: false,
		},
		{
			name: "valid machine pool",
			platform: &gcp.Platform{
				Region:                 "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{},
			},
			valid: true,
		},
		{
			name: "valid subnets & network",
			platform: &gcp.Platform{
				Region:             "us-east1",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			valid: true,
		},
		{
			name: "missing subnets",
			platform: &gcp.Platform{
				Region:  "us-east1",
				Network: "valid-vpc",
			},
			valid: false,
		},
		{
			name: "subnets missing network",
			platform: &gcp.Platform{
				Region:        "us-east1",
				ComputeSubnet: "valid-compute-subnet",
			},
			valid: false,
		},
		{
			name: "unsupported GCP disk type",
			platform: &gcp.Platform{
				Region: "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						DiskType: "pd-standard",
					},
				},
			},
			valid: false,
		},

		{
			name: "supported GCP disk type",
			platform: &gcp.Platform{
				Region: "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						DiskType: "pd-ssd",
					},
				},
			},
			valid: true,
		},
		{
			name: "GCP valid network project data",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           true,
		},
		{
			name: "GCP invalid network project missing network",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP invalid network project missing compute subnet",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP invalid network project missing control plane subnet",
			platform: &gcp.Platform{
				Region:           "us-east1",
				NetworkProjectID: "valid-network-project",
				ProjectID:        "valid-project",
				Network:          "valid-vpc",
				ComputeSubnet:    "valid-compute-subnet",
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP invalid network project bad credentials mode",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			credentialsMode: types.MintCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP missing network project with private zone",
			platform: &gcp.Platform{
				Region:             "us-east1",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
				DNS: &gcp.DNS{
					PrivateZone: &gcp.DNSZone{
						Name: "test-private-zone-name",
					},
				},
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP missing Zone with private zone",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
				DNS: &gcp.DNS{
					PrivateZone: &gcp.DNSZone{
						ProjectID: "valid-project",
					},
				},
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           false,
		},
		{
			name: "GCP valid private zone",
			platform: &gcp.Platform{
				Region:             "us-east1",
				NetworkProjectID:   "valid-network-project",
				ProjectID:          "valid-project",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
				DNS: &gcp.DNS{
					PrivateZone: &gcp.DNSZone{
						ProjectID: "valid-project",
						Name:      "test-private-zone-name",
					},
				},
			},
			credentialsMode: types.PassthroughCredentialsMode,
			valid:           true,
		},
		{
			name: "invalid gcp endpoint no network",
			platform: &gcp.Platform{
				Region:    "us-east1",
				ProjectID: "valid-project",
				Endpoint: &gcp.PSCEndpoint{
					Name: "test-endpoint",
				},
			},
			valid: false,
		},
		{
			name: "invalid mutual exclusivity for endpoint and custom dns",
			platform: &gcp.Platform{
				UserProvisionedDNS: dns.UserProvisionedDNSEnabled,
				Region:             "us-east1",
				ProjectID:          "valid-project",
				Endpoint: &gcp.PSCEndpoint{
					Name: "test-endpoint",
				},
			},
			valid: false,
		},
		{
			name: "invalid firewall management configuration",
			platform: &gcp.Platform{
				UserProvisionedDNS:      dns.UserProvisionedDNSEnabled,
				FirewallRulesManagement: gcp.UnmanagedFirewallRules,
				Region:                  "us-east1",
				ProjectID:               "valid-project",
			},
			valid: false,
		},
		{
			name: "invalid firewall management",
			platform: &gcp.Platform{
				UserProvisionedDNS:      dns.UserProvisionedDNSEnabled,
				FirewallRulesManagement: "random-test",
				Region:                  "us-east1",
				ProjectID:               "valid-project",
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			credentialsMode := tc.credentialsMode
			if credentialsMode == "" {
				credentialsMode = types.MintCredentialsMode
			}

			publishStrategy := types.ExternalPublishingStrategy
			if tc.publishStrategy != "" {
				publishStrategy = tc.publishStrategy
			}

			// the only item currently used is the credentialsMode
			ic := types.InstallConfig{
				CredentialsMode: credentialsMode,
				Publish:         publishStrategy,
			}

			err := ValidatePlatform(tc.platform, field.NewPath("test-path"), &ic).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateUserLabels(t *testing.T) {
	fieldPath := "spec.platform.gcp.UserLabels"
	cases := []struct {
		name        string
		userLabels  []gcp.UserLabel
		expectedErr string
	}{
		{
			name:        "userLabels not configured",
			userLabels:  nil,
			expectedErr: "[]",
		},
		{
			name: "userLabels configured",
			userLabels: []gcp.UserLabel{
				{Key: "key1", Value: "value1"},
				{Key: "key_2", Value: "value_2"},
				{Key: "key-3", Value: "value-3"},
				{Key: "key4_", Value: "value4_"},
				{Key: "key5-", Value: "value5-"},
			},
			expectedErr: "[]",
		},
		{
			name: "userLabels configured is more than max limit",
			userLabels: []gcp.UserLabel{
				{Key: "key11", Value: "value11"}, {Key: "key18", Value: "value18"},
				{Key: "key19", Value: "value19"}, {Key: "key21", Value: "value21"},
				{Key: "key14", Value: "value14"}, {Key: "key22", Value: "value22"},
				{Key: "key25", Value: "value25"}, {Key: "key27", Value: "value27"},
				{Key: "key31", Value: "value31"}, {Key: "key9", Value: "value9"},
				{Key: "key10", Value: "value10"}, {Key: "key15", Value: "value15"},
				{Key: "key28", Value: "value28"}, {Key: "key29", Value: "value29"},
				{Key: "key32", Value: "value32"}, {Key: "key3", Value: "value3"},
				{Key: "key7", Value: "value7"}, {Key: "key17", Value: "value17"},
				{Key: "key20", Value: "value20"}, {Key: "key4", Value: "value4"},
				{Key: "key23", Value: "value23"}, {Key: "key26", Value: "value26"},
				{Key: "key12", Value: "value12"}, {Key: "key33", Value: "value33"},
				{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"},
				{Key: "key5", Value: "value5"}, {Key: "key8", Value: "value8"},
				{Key: "key30", Value: "value30"}, {Key: "key6", Value: "value6"},
				{Key: "key13", Value: "value13"}, {Key: "key16", Value: "value16"},
				{Key: "key24", Value: "value24"},
			},
			expectedErr: "[spec.platform.gcp.UserLabels: Too many: 33: must have at most 32 items]",
		},
		{
			name:        "userLabels contains key starting a number",
			userLabels:  []gcp.UserLabel{{Key: "1key", Value: "1value"}},
			expectedErr: "[spec.platform.gcp.UserLabels[1key]: Invalid value: \"1value\": label key is invalid or contains invalid characters. Label key can have a maximum of 63 characters and cannot be empty. Label key must begin with a lowercase letter, and must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name:        "userLabels contains key starting a uppercase letter",
			userLabels:  []gcp.UserLabel{{Key: "Key", Value: "1value"}},
			expectedErr: "[spec.platform.gcp.UserLabels[Key]: Invalid value: \"1value\": label key is invalid or contains invalid characters. Label key can have a maximum of 63 characters and cannot be empty. Label key must begin with a lowercase letter, and must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name:        "userLabels contains empty key",
			userLabels:  []gcp.UserLabel{{Key: "", Value: "value"}},
			expectedErr: "[spec.platform.gcp.UserLabels[]: Invalid value: \"value\": label key is invalid or contains invalid characters. Label key can have a maximum of 63 characters and cannot be empty. Label key must begin with a lowercase letter, and must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name: "userLabels contains key length greater than 63",
			userLabels: []gcp.UserLabel{
				{
					Key:   "thisisaverylongkeywithmorethan63characterswhichisnotallowedforgcpresourcelabelkey",
					Value: "value",
				},
			},
			expectedErr: "[spec.platform.gcp.UserLabels[thisisaverylongkeywithmorethan63characterswhichisnotallowedforgcpresourcelabelkey]: Invalid value: \"value\": label key is invalid or contains invalid characters. Label key can have a maximum of 63 characters and cannot be empty. Label key must begin with a lowercase letter, and must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name:        "userLabels contains key with invalid character",
			userLabels:  []gcp.UserLabel{{Key: "key/test", Value: "value"}},
			expectedErr: "[spec.platform.gcp.UserLabels[key/test]: Invalid value: \"value\": label key is invalid or contains invalid characters. Label key can have a maximum of 63 characters and cannot be empty. Label key must begin with a lowercase letter, and must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name: "userLabels contains value length greater than 63",
			userLabels: []gcp.UserLabel{
				{
					Key:   "key",
					Value: "thisisaverylongvaluewithmorethan63characterswhichisnotallowedforgcpresourcelabelvalue",
				},
			},
			expectedErr: "[spec.platform.gcp.UserLabels[key]: Invalid value: \"thisisaverylongvaluewithmorethan63characterswhichisnotallowedforgcpresourcelabelvalue\": label value is invalid or contains invalid characters. Label value can have a maximum of 63 characters and cannot be empty. Value must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name:        "userLabels contains empty value",
			userLabels:  []gcp.UserLabel{{Key: "key", Value: ""}},
			expectedErr: "[spec.platform.gcp.UserLabels[key]: Invalid value: \"\": label value is invalid or contains invalid characters. Label value can have a maximum of 63 characters and cannot be empty. Value must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name:        "userLabels contains value with invalid character",
			userLabels:  []gcp.UserLabel{{Key: "key", Value: "value*^%"}},
			expectedErr: "[spec.platform.gcp.UserLabels[key]: Invalid value: \"value*^%\": label value is invalid or contains invalid characters. Label value can have a maximum of 63 characters and cannot be empty. Value must contain only lowercase letters, numeric characters, and the following special characters `_-`]",
		},
		{
			name:        "userLabels contains key with prefix kubernetes-io",
			userLabels:  []gcp.UserLabel{{Key: "kubernetes-io_cluster", Value: "value"}},
			expectedErr: "[spec.platform.gcp.UserLabels[kubernetes-io_cluster]: Invalid value: \"value\": label key contains restricted prefix. Label key cannot have `kubernetes-io`, `openshift-io` prefixes]",
		},
		{
			name:        "userLabels contains allowed key prefix for_openshift-io",
			userLabels:  []gcp.UserLabel{{Key: "for_openshift-io", Value: "gcp"}},
			expectedErr: "[]",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserLabels(tt.userLabels, field.NewPath(fieldPath))
			if fmt.Sprintf("%v", err) != tt.expectedErr {
				t.Errorf("Got: %+v Want: %+v", err, tt.expectedErr)
			}
		})
	}
}
