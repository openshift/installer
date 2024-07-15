package configimage

import (
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/imagebased"
)

const (
	rawNMStateConfig = `
interfaces:
  - name: eth0
    type: ethernet
    state: up
    mac-address: 00:00:00:00:00:00
    ipv4:
      enabled: true
      address:
        - ip: 192.168.122.2
          prefix-length: 23
      dhcp: false`

	testCert = `-----BEGIN CERTIFICATE-----
MIICYTCCAcqgAwIBAgIJAI2kA+uXAbhOMA0GCSqGSIb3DQEBCwUAMEgxCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJDQTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzEUMBIG
A1UECgwLUmVkIEhhdCBJbmMwHhcNMTkwMjEyMTkzMjUzWhcNMTkwMjEzMTkzMjUz
WjBIMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFu
Y2lzY28xFDASBgNVBAoMC1JlZCBIYXQgSW5jMIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQC+HOC0mKig/oINAKPo88LqxDJ4l7lozdLtp5oGeqWrLUXSfkvXAkQY
2QYdvPAjpRfH7Ii7G0Asx+HTKdvula7B5fXDjc6NYKuEpTJZRV1ugntI97bozF/E
C2BBmxxEnJN3+Xe8RYXMjz5Q4aqPw9vZhlWN+0hrREl1Ea/zHuWFIQIDAQABo1Mw
UTAdBgNVHQ4EFgQUvTS1XjlvOdsufSyWxukyQu3LriEwHwYDVR0jBBgwFoAUvTS1
XjlvOdsufSyWxukyQu3LriEwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsF
AAOBgQB9gFcOXnzJrM65QqxeCB9Z5l5JMjp45UFC9Bj2cgwDHP80Zvi4omlaacC6
aavmnLd67zm9PbYDWRaOIWAMeB916Iwaw/v6I0jwhAk/VxX5Fl6cGlZu9jZ3zbFE
2sDqkwzIuSjCG2A23s6d4M1S3IXCCydoCSLMu+WhLkbboK6jEg==
-----END CERTIFICATE-----
`
)

func defaultInstallConfig() *InstallConfig {
	_, newCidr, err := net.ParseCIDR("192.168.111.0/24")
	if err != nil {
		return nil
	}
	_, machineNetCidr, err := net.ParseCIDR("10.10.11.0/24")
	if err != nil {
		return nil
	}

	return &InstallConfig{
		AssetBase: installconfig.AssetBase{
			Config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ocp-ibi",
				},
				BaseDomain: "testing.com",
				PullSecret: testSecret,
				SSHKey:     testSSHKey,
				ControlPlane: &types.MachinePool{
					Name:     "controlplane",
					Replicas: ptr.To[int64](1),
					Platform: types.MachinePoolPlatform{},
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{
							CIDR: ipnet.IPNet{IPNet: *machineNetCidr},
						},
					},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       ipnet.IPNet{IPNet: *newCidr},
							HostPrefix: 23,
						},
					},
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
					NetworkType: "OVNKubernetes",
				},
			},
		},
	}
}

func imageBasedConfig() *ImageBasedConfig {
	ibConfig := &ImageBasedConfig{
		Config: &imagebased.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "image-based-config",
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: imagebased.ImageBasedConfigVersion,
			},
			Hostname:        "somehostname",
			ReleaseRegistry: "quay.io",
			NetworkConfig: aiv1beta1.NetConfig{
				Raw: unmarshalJSON([]byte(rawNMStateConfig)),
			},
		},
	}
	return ibConfig
}

func unmarshalJSON(b []byte) []byte {
	output, err := yaml.JSONToYAML(b)
	if err != nil {
		return nil
	}
	return output
}
