package manifests

import (
	"testing"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/stretchr/testify/assert"
)

func TestGetNodeZeroIP(t *testing.T) {
	tests := []struct {
		name       string
		nmConfig   []aiv1beta1.NMStateConfig
		expectedIP string
	}{
		{
			name: "Valid NMStateConfig with only one IPV4 IP address",
			nmConfig: func() []aiv1beta1.NMStateConfig {
				config := GetDefaultConfig()
				config[0].Spec.NetConfig.Raw = []byte(`
interfaces:
    - ipv4:
        address:
            - ip: 192.168.122.2
`)
				return config
			}(),
			expectedIP: "192.168.122.2",
		},
		{
			name: "Valid NMStateConfig with only one IPV6 IP address",
			nmConfig: func() []aiv1beta1.NMStateConfig {
				config := GetDefaultConfig()
				config[0].Spec.NetConfig.Raw = []byte(`
interfaces:
    - ipv6:
        address:
            - ip: 2600:1700:c1a0:fe10:96e6:f7ff:fe91:2590
`)
				return config
			}(),
			expectedIP: "2600:1700:c1a0:fe10:96e6:f7ff:fe91:2590",
		},
		{
			name: "Valid NMStateConfig with multiple interfaces",
			nmConfig: func() []aiv1beta1.NMStateConfig {
				config := GetDefaultConfig()
				config[0].Spec.NetConfig.Raw = []byte(`
interfaces:
    - ipv4:
        address:
            - ip: 192.168.122.2
    - ipv4:
        address:
            - ip: 192.168.122.3
`)
				return config
			}(),
			expectedIP: "192.168.122.2",
		},
		{
			name: "Valid NMStateConfig with multiple IPV4 IP addresses in a single interface",
			nmConfig: func() []aiv1beta1.NMStateConfig {
				config := GetDefaultConfig()
				config[0].Spec.NetConfig.Raw = []byte(`
interfaces:
    - ipv4:
        address:
            - ip: 192.168.122.2
            - ip: 192.168.122.3
`)
				return config
			}(),
			expectedIP: "192.168.122.2",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n := NMConfig{
				nmConfig: func() ([]aiv1beta1.NMStateConfig, error) {
					return test.nmConfig, nil
				},
			}
			nodeZeroIP := n.GetNodeZeroIP()
			assert.Equal(t, nodeZeroIP, test.expectedIP)

		})
	}
}

func GetDefaultConfig() []aiv1beta1.NMStateConfig {
	var nmStateConfigList []aiv1beta1.NMStateConfig
	data := aiv1beta1.NMStateConfig{
		Spec: aiv1beta1.NMStateConfigSpec{
			NetConfig: aiv1beta1.NetConfig{
				Raw: []byte(""),
			},
		},
	}
	nmStateConfigList = append(nmStateConfigList, data)
	return nmStateConfigList
}
