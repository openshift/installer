package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	validType            = "valid-type"
	validZones           = []string{"zone-a", "zone-b"}
	validEncryptionKey   = "crn:v1:bluemix:public:kms:global:a/accountid:service:key:keyid"
	invalidEncryptionKey = "v1:bluemix:kms:global:a/accountid:service:key:keyid"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name        string
		machinepool *ibmcloud.MachinePool
		valid       bool
	}{
		{
			name:        "minimal",
			machinepool: &ibmcloud.MachinePool{},
			valid:       true,
		},
		{
			name: "valid type",
			machinepool: &ibmcloud.MachinePool{
				InstanceType: validType,
			},
			valid: true,
		},
		{
			name: "valid zones",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
			},
			valid: true,
		},
		{
			name: "valid bootVolume",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					EncryptionKey: validEncryptionKey,
				},
			},
			valid: true,
		},
		{
			name: "valid bootVolume",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					EncryptionKey: invalidEncryptionKey,
				},
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.machinepool, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
