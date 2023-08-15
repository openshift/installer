package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name            string
		platform        *gcp.Platform
		credentialsMode types.CredentialsMode
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			credentialsMode := tc.credentialsMode
			if credentialsMode == "" {
				credentialsMode = types.MintCredentialsMode
			}

			// the only item currently used is the credentialsMode
			ic := types.InstallConfig{
				CredentialsMode: credentialsMode,
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

func TestValidateUserTags(t *testing.T) {
	fieldPath := "spec.platform.gcp.userTags"
	cases := []struct {
		name        string
		userTags    []gcp.UserTag
		expectedErr string
	}{
		{
			name:        "userTags not configured",
			userTags:    []gcp.UserTag{},
			expectedErr: "[]",
		},
		{
			name: "userTags configured",
			userTags: []gcp.UserTag{
				{ParentID: "1234567890", Key: "key_2", Value: "value_2"},
				{ParentID: "test-project-123", Key: "key.gcp", Value: "value.3"},
				{ParentID: "1234567890", Key: "keY", Value: "value"},
				{ParentID: "test-project-123", Key: "thisisalongkeywithinlimitof63_characters-whichisallowedfortags", Value: "value"},
				{ParentID: "1234567890", Key: "KEY4", Value: "hisisavaluewithin-63characters_{[(.@%=+: ,*#&)]}forgcptagvalue"},
				{ParentID: "test-project-123", Key: "key1", Value: "value1"},
			},
			expectedErr: "[]",
		},
		{
			name: "userTags configured is more than max limit",
			userTags: []gcp.UserTag{
				{ParentID: "1234567890", Key: "key29", Value: "value29"},
				{ParentID: "test-project-123", Key: "key33", Value: "value33"},
				{ParentID: "1234567890", Key: "key39", Value: "value39"},
				{ParentID: "test-project-123", Key: "key43", Value: "value43"},
				{ParentID: "1234567890", Key: "key5", Value: "value5"},
				{ParentID: "test-project-123", Key: "key6", Value: "value6"},
				{ParentID: "1234567890", Key: "key14", Value: "value14"},
				{ParentID: "test-project-123", Key: "key25", Value: "value25"},
				{ParentID: "1234567890", Key: "key20", Value: "value20"},
				{ParentID: "test-project-123", Key: "key24", Value: "value24"},
				{ParentID: "1234567890", Key: "key40", Value: "value40"},
				{ParentID: "test-project-123", Key: "key46", Value: "value46"},
				{ParentID: "1234567890", Key: "key1", Value: "value1"},
				{ParentID: "test-project-123", Key: "key2", Value: "value2"},
				{ParentID: "1234567890", Key: "key4", Value: "value4"},
				{ParentID: "test-project-123", Key: "key10", Value: "value10"},
				{ParentID: "1234567890", Key: "key51", Value: "value51"},
				{ParentID: "test-project-123", Key: "key8", Value: "value8"},
				{ParentID: "1234567890", Key: "key13", Value: "value13"},
				{ParentID: "test-project-123", Key: "key44", Value: "value44"},
				{ParentID: "1234567890", Key: "key48", Value: "value48"},
				{ParentID: "test-project-123", Key: "key9", Value: "value9"},
				{ParentID: "1234567890", Key: "key17", Value: "value17"},
				{ParentID: "test-project-123", Key: "key18", Value: "value18"},
				{ParentID: "1234567890", Key: "key30", Value: "value30"},
				{ParentID: "test-project-123", Key: "key36", Value: "value36"},
				{ParentID: "1234567890", Key: "key49", Value: "value49"},
				{ParentID: "test-project-123", Key: "key7", Value: "value7"},
				{ParentID: "1234567890", Key: "key15", Value: "value15"},
				{ParentID: "test-project-123", Key: "key22", Value: "value22"},
				{ParentID: "1234567890", Key: "key34", Value: "value34"},
				{ParentID: "test-project-123", Key: "key37", Value: "value37"},
				{ParentID: "1234567890", Key: "key38", Value: "value38"},
				{ParentID: "test-project-123", Key: "key47", Value: "value47"},
				{ParentID: "1234567890", Key: "key12", Value: "value12"},
				{ParentID: "test-project-123", Key: "key16", Value: "value16"},
				{ParentID: "1234567890", Key: "key23", Value: "value23"},
				{ParentID: "test-project-123", Key: "key28", Value: "value28"},
				{ParentID: "1234567890", Key: "key50", Value: "value50"},
				{ParentID: "test-project-123", Key: "key21", Value: "value21"},
				{ParentID: "1234567890", Key: "key26", Value: "value26"},
				{ParentID: "test-project-123", Key: "key35", Value: "value35"},
				{ParentID: "1234567890", Key: "key42", Value: "value42"},
				{ParentID: "test-project-123", Key: "key31", Value: "value31"},
				{ParentID: "1234567890", Key: "key32", Value: "value32"},
				{ParentID: "test-project-123", Key: "key41", Value: "value41"},
				{ParentID: "1234567890", Key: "key45", Value: "value45"},
				{ParentID: "test-project-123", Key: "key3", Value: "value3"},
				{ParentID: "1234567890", Key: "key11", Value: "value11"},
				{ParentID: "test-project-123", Key: "key19", Value: "value19"},
				{ParentID: "1234567890", Key: "key27", Value: "value27"},
			},
			expectedErr: "[spec.platform.gcp.userTags: Too many: 51: must have at most 50 items]",
		},
		{
			name:        "userTags contains key starting with a special character",
			userTags:    []gcp.UserTag{{ParentID: "1234567890", Key: "_key", Value: "1value"}},
			expectedErr: "[spec.platform.gcp.userTags[_key]: Invalid value: \"1value\": tag key is invalid or contains invalid characters. Tag key can have a maximum of 63 characters and cannot be empty. Tag key must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `._-`]",
		},
		{
			name:        "userTags contains key ending with a special character",
			userTags:    []gcp.UserTag{{ParentID: "1234567890", Key: "key@", Value: "1value"}},
			expectedErr: "[spec.platform.gcp.userTags[key@]: Invalid value: \"1value\": tag key is invalid or contains invalid characters. Tag key can have a maximum of 63 characters and cannot be empty. Tag key must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `._-`]",
		},
		{
			name:        "userTags contains empty key",
			userTags:    []gcp.UserTag{{ParentID: "1234567890", Key: "", Value: "value"}},
			expectedErr: "[spec.platform.gcp.userTags[]: Invalid value: \"value\": tag key is invalid or contains invalid characters. Tag key can have a maximum of 63 characters and cannot be empty. Tag key must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `._-`]",
		},
		{
			name: "userTags contains key length greater than 63",
			userTags: []gcp.UserTag{
				{
					ParentID: "1234567890",
					Key:      "thisisalongkeyforlimitof63_characters-whichisnotallowedfortagkey",
					Value:    "value",
				},
			},
			expectedErr: "[spec.platform.gcp.userTags[thisisalongkeyforlimitof63_characters-whichisnotallowedfortagkey]: Invalid value: \"value\": tag key is invalid or contains invalid characters. Tag key can have a maximum of 63 characters and cannot be empty. Tag key must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `._-`]",
		},
		{
			name:        "userTags contains key with invalid character",
			userTags:    []gcp.UserTag{{ParentID: "1234567890", Key: "key/test", Value: "value"}},
			expectedErr: "[spec.platform.gcp.userTags[key/test]: Invalid value: \"value\": tag key is invalid or contains invalid characters. Tag key can have a maximum of 63 characters and cannot be empty. Tag key must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `._-`]",
		},
		{
			name: "userTags contains value length greater than 63",
			userTags: []gcp.UserTag{
				{
					ParentID: "1234567890",
					Key:      "key",
					Value:    "hisisavaluewith-63characters_{[(.@%=+: ,*#&)]}allowedforgcptagvalue",
				},
			},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"hisisavaluewith-63characters_{[(.@%=+: ,*#&)]}allowedforgcptagvalue\": tag value is invalid or contains invalid characters. Tag value can have a maximum of 63 characters and cannot be empty. Tag value must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `_-.@%=+:,*#&(){}[]` and spaces]",
		},
		{
			name:        "userTags contains empty value",
			userTags:    []gcp.UserTag{{ParentID: "1234567890", Key: "key", Value: ""}},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"\": tag value is invalid or contains invalid characters. Tag value can have a maximum of 63 characters and cannot be empty. Tag value must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `_-.@%=+:,*#&(){}[]` and spaces]",
		},
		{
			name:        "userTags contains value with invalid character",
			userTags:    []gcp.UserTag{{ParentID: "1234567890", Key: "key", Value: "value*^%"}},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"value*^%\": tag value is invalid or contains invalid characters. Tag value can have a maximum of 63 characters and cannot be empty. Tag value must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `_-.@%=+:,*#&(){}[]` and spaces]",
		},
		{
			name:        "userTags contains empty ParentID",
			userTags:    []gcp.UserTag{{Key: "key", Value: "value*^%"}},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"value*^%\": tag parentID is invalid or contains invalid characters. ParentID can have a maximum of 32 characters and cannot be empty. ParentID can be either OrganizationID or ProjectID. OrganizationID must consist of decimal numbers, and cannot have leading zeroes and ProjectID must be 6 to 30 characters in length, can only contain lowercase letters, numbers, and hyphens, and must start with a letter, and cannot end with a hyphen]",
		},
		{
			name:        "userTags contains ParentID configured with invalid OrganizationID",
			userTags:    []gcp.UserTag{{ParentID: "00001234567890", Key: "key", Value: "value"}},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"value\": tag parentID is invalid or contains invalid characters. ParentID can have a maximum of 32 characters and cannot be empty. ParentID can be either OrganizationID or ProjectID. OrganizationID must consist of decimal numbers, and cannot have leading zeroes and ProjectID must be 6 to 30 characters in length, can only contain lowercase letters, numbers, and hyphens, and must start with a letter, and cannot end with a hyphen]",
		},
		{
			name:        "userTags contains ParentID configured with invalid ProjectID",
			userTags:    []gcp.UserTag{{ParentID: "test-project-123-", Key: "key", Value: "value"}},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"value\": tag parentID is invalid or contains invalid characters. ParentID can have a maximum of 32 characters and cannot be empty. ParentID can be either OrganizationID or ProjectID. OrganizationID must consist of decimal numbers, and cannot have leading zeroes and ProjectID must be 6 to 30 characters in length, can only contain lowercase letters, numbers, and hyphens, and must start with a letter, and cannot end with a hyphen]",
		},
		{
			name: "userTags contains ParentID configured with invalid OrganizationID length",
			userTags: []gcp.UserTag{
				{
					ParentID: "123456789012345678901234567890123",
					Key:      "key",
					Value:    "value",
				},
			},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"value\": tag parentID is invalid or contains invalid characters. ParentID can have a maximum of 32 characters and cannot be empty. ParentID can be either OrganizationID or ProjectID. OrganizationID must consist of decimal numbers, and cannot have leading zeroes and ProjectID must be 6 to 30 characters in length, can only contain lowercase letters, numbers, and hyphens, and must start with a letter, and cannot end with a hyphen]",
		},
		{
			name: "userTags contains ParentID configured with invalid ProjectID length",
			userTags: []gcp.UserTag{
				{
					ParentID: "test-project-123-test-project-123-test-project-123-test-project-123",
					Key:      "key",
					Value:    "value",
				},
			},
			expectedErr: "[spec.platform.gcp.userTags[key]: Invalid value: \"value\": tag parentID is invalid or contains invalid characters. ParentID can have a maximum of 32 characters and cannot be empty. ParentID can be either OrganizationID or ProjectID. OrganizationID must consist of decimal numbers, and cannot have leading zeroes and ProjectID must be 6 to 30 characters in length, can only contain lowercase letters, numbers, and hyphens, and must start with a letter, and cannot end with a hyphen]",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserTags(tt.userTags, field.NewPath(fieldPath))
			if fmt.Sprintf("%v", err) != tt.expectedErr {
				t.Errorf("Got: %+v Want: %+v", err, tt.expectedErr)
			}
		})
	}
}
