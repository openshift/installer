package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/common"
	"github.com/openshift/installer/pkg/types/none"
)

func TestValidateProvisioning(t *testing.T) {
	cases := []struct {
		name     string
		config   *types.InstallConfig
		platform *none.Platform
		expected string
	}{
		{
			name: "valid_empty_credentials",
			platform: platform().
				FencingCredentials().build(),
			expected: "",
		},
		{
			name: "valid_two_credentials",
			platform: platform().
				FencingCredentials(fc1(),
					fc2()).build(),
			expected: "",
		},
		{
			name: "invalid_number_of_credentials",
			platform: platform().
				FencingCredentials(fc1(),
					fc2(),
					fc3()).build(),
			expected: "Forbidden: there should be exactly two fencingCredentials to support the two node cluster, instead 3 fencingCredentials were found",
		},
		{
			name: "duplicate_bmc_address",
			platform: platform().
				FencingCredentials(
					fc1().BMCAddress("ipmi://192.168.111.1"),
					fc2().BMCAddress("ipmi://192.168.111.1")).build(),
			expected: "none.fencingCredentials\\[1\\].Address: Duplicate value: \"ipmi://192.168.111.1\"",
		},
		{
			name: "bmc_address_required",
			platform: platform().
				FencingCredentials(fc1().BMCAddress("")).build(),
			expected: "none.fencingCredentials\\[0\\].Address: Required value: missing Address",
		},
		{
			name: "bmc_username_required",
			platform: platform().
				FencingCredentials(fc1().BMCUsername("")).build(),
			expected: "none.fencingCredentials\\[0\\].Username: Required value: missing Username",
		},
		{
			name: "bmc_password_required",
			platform: platform().
				FencingCredentials(fc1().BMCPassword("")).build(),
			expected: "none.fencingCredentials\\[0\\].Password: Required value: missing Password",
		},
		{
			name: "host_name_required",
			platform: platform().
				FencingCredentials(fc1().HostName("")).build(),
			expected: "none.fencingCredentials\\[0\\].HostName: Required value: missing HostName",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
			}
			tc.config.None = tc.platform

			err := ValidateFencingCredentials(tc.platform.FencingCredentials, field.NewPath("none")).ToAggregate()

			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

type fcBuilder struct {
	common.FencingCredential
}

func fc1() *fcBuilder {
	return &fcBuilder{
		common.FencingCredential{
			HostName: "host1",
			Username: "root",
			Password: "password",
			Address:  "ipmi://192.168.111.1",
		},
	}
}

func fc2() *fcBuilder {
	return &fcBuilder{
		common.FencingCredential{
			HostName: "host2",
			Username: "root",
			Password: "password",
			Address:  "ipmi://192.168.111.2",
		},
	}
}

func fc3() *fcBuilder {
	return &fcBuilder{
		common.FencingCredential{
			HostName: "host3",
			Username: "root",
			Password: "password",
			Address:  "ipmi://192.168.111.3",
		},
	}
}

func (hb *fcBuilder) build() *common.FencingCredential {
	return &hb.FencingCredential
}

func (hb *fcBuilder) HostName(value string) *fcBuilder {
	hb.FencingCredential.HostName = value
	return hb
}

func (hb *fcBuilder) BMCAddress(value string) *fcBuilder {
	hb.FencingCredential.Address = value
	return hb
}

func (hb *fcBuilder) BMCUsername(value string) *fcBuilder {
	hb.FencingCredential.Username = value
	return hb
}

func (hb *fcBuilder) BMCPassword(value string) *fcBuilder {
	hb.FencingCredential.Password = value
	return hb
}

type platformBuilder struct {
	none.Platform
}

func platform() *platformBuilder {
	return &platformBuilder{
		none.Platform{}}
}

func (pb *platformBuilder) build() *none.Platform {
	return &pb.Platform
}

func (pb *platformBuilder) FencingCredentials(builders ...*fcBuilder) *platformBuilder {
	pb.Platform.FencingCredentials = nil
	for _, builder := range builders {
		pb.Platform.FencingCredentials = append(pb.Platform.FencingCredentials, builder.build())
	}
	return pb
}

type installConfigBuilder struct {
	types.InstallConfig
}

func installConfig() *installConfigBuilder {
	return &installConfigBuilder{
		InstallConfig: types.InstallConfig{},
	}
}

func (icb *installConfigBuilder) build() *types.InstallConfig {
	return &icb.InstallConfig
}
