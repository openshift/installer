package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func validPlatform() *ovirt.Platform {
	return &ovirt.Platform{
		ClusterID:              "0953a3fb-6682-49a6-8767-43f561045a59",
		StorageDomainID:        "57e42205-02ac-46e1-a6d1-07459d94bc51",
		VNICProfileID:          "57e42205-02ac-46e1-a6d1-07459d94bc52",
		NetworkName:            "ocp-blue",
		APIVIP:                 "10.0.0.1",
		DNSVIP:                 "10.0.0.2",
		IngressVIP:             "10.0.0.3",
		DefaultMachinePlatform: nil,
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *ovirt.Platform
		valid    bool
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
			valid:    true,
		},
		{
			name: "invalid when empty cluster ID",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.ClusterID = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid when empty storage ID",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.StorageDomainID = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid when API VIP is invalid",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.APIVIP = "1."
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid when DNS VIP is invalid",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.APIVIP = "1."
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid when INGRESS VIP is invalid",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.APIVIP = "1."
				return p
			}(),
			valid: false,
		},
		{
			name: "malformed vnic profile id",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.VNICProfileID = "abcd-sdf"
				return p
			}(),
			valid: false,
		},
		{
			name: "valid empty vnic profile id",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.VNICProfileID = ""
				return p
			}(),
			valid: true,
		},
		{
			name: "valid machine pool",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &ovirt.MachinePool{}
				return p
			}(),
			valid: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
