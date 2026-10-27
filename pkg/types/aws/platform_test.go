package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

// Programs that vendor these types, such as Hive, unmarshal a user
// install-config and marshal it again. Installers older than 4.19 do not know
// the vpc field, so an empty VPC must not appear in the output.
func TestPlatformMarshalOmitsEmptyVPC(t *testing.T) {
	out, err := yaml.Marshal(Platform{Region: "us-east-1"})
	require.NoError(t, err)
	assert.Equal(t, "region: us-east-1\n", string(out))
}

func TestPlatformMarshalKeepsConfiguredVPC(t *testing.T) {
	p := Platform{
		Region: "us-east-1",
		VPC: VPC{
			Subnets: []Subnet{{ID: "subnet-0123456789abcdef0"}},
		},
	}
	out, err := yaml.Marshal(p)
	require.NoError(t, err)
	assert.Contains(t, string(out), "vpc:")
	assert.Contains(t, string(out), "id: subnet-0123456789abcdef0")
}
