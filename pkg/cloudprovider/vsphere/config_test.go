package vsphere

import (
	"fmt"
	"strings"
	"testing"

	gmg "github.com/onsi/gomega"
)

const basicConfigINI = `
[Global]
server = 0.0.0.0
port = 443
user = user
password = password
insecure-flag = true
datacenters = us-west
ca-file = /some/path/to/a/ca.pem
`

const basicConfigYaml = `
global:
  caFile: /some/path/to/a/ca.pem
  datacenters:
  - us-west
  insecureFlag: true
  password: password
  port: 443
  server: 0.0.0.0
  user: user
`

const basicConfigVcenterSectionINI = `
[Global]
secret-name = "global-secret"
secret-namespace = "global-secret-ns"

[VirtualCenter "vc.rh.com"]
datacenters = "DC0,DC1"

[Labels]
region = "k8s-region"
zone = "k8s-zone"
`

const basicConfigVcenterSectionYAML = `
global:
  secretName: global-secret
  secretNamespace: global-secret-ns
labels:
  region: k8s-region
  zone: k8s-zone
vcenter:
  vc.rh.com:
    datacenters:
    - DC0
    - DC1
    server: vc.rh.com
`

const multiVCDCsConfigINI = `
[Global]
port = 443
insecure-flag = true
secret-name = "global-secret"
secret-namespace = "global-secret-ns"

[VirtualCenter "t1"]
server = "10.0.0.1"
datacenters = "DC0,DC1,DC2"
secret-name = "tenant1-secret"
secret-namespace = "kube-system"

[VirtualCenter "10.0.0.2"]
datacenters = "DC3"

[VirtualCenter "10.0.0.3"]
datacenters = "DC5,DC6"
ip-family = "ipv6"
`

const invalidConfig = `boom[]{}`

const invalidConfigWrongGlobalPort = `
[Global]
port = -443
insecure-flag = true
`

const invalidConfigWrongVCPort = `
[Global]
port = 443
insecure-flag = true

[VirtualCenter "10.0.0.3"]
datacenters = "DC5,DC6"
port = -1
ip-family = "ipv6"
`

func TestINIConfigConversion(t *testing.T) {
	t.Run("basic config yaml conversion", func(t *testing.T) {
		g := gmg.NewWithT(t)
		configStruct, err := ReadConfig([]byte(basicConfigINI))
		g.Expect(err).ToNot(gmg.HaveOccurred())

		convertedConfig, err := MarshalConfig(configStruct)
		g.Expect(err).ToNot(gmg.HaveOccurred())
		// Trim left emptyline just for keep constant more readable
		g.Expect(convertedConfig).To(gmg.BeEquivalentTo(strings.TrimLeft(basicConfigYaml, "\n")))
	})

	t.Run("basic config yaml conversion with vc section", func(t *testing.T) {
		g := gmg.NewWithT(t)
		configStruct, err := ReadConfig([]byte(basicConfigVcenterSectionINI))
		g.Expect(err).ToNot(gmg.HaveOccurred())

		convertedConfig, err := MarshalConfig(configStruct)
		g.Expect(err).ToNot(gmg.HaveOccurred())
		// Trim left emptyline just for keep constant more readable
		g.Expect(convertedConfig).To(gmg.BeEquivalentTo(strings.TrimLeft(basicConfigVcenterSectionYAML, "\n")))
	})

	t.Run("Test multi DCs config yaml conversion", func(t *testing.T) {
		g := gmg.NewWithT(t)
		configStruct, err := ReadConfig([]byte(multiVCDCsConfigINI))
		g.Expect(err).ToNot(gmg.HaveOccurred())

		g.Expect(len(configStruct.Vcenter)).To(gmg.Equal(3))

		VC1 := configStruct.Vcenter["t1"]
		g.Expect(len(VC1.Datacenters)).To(gmg.Equal(3))
		g.Expect(VC1.Datacenters).To(gmg.BeComparableTo([]string{"DC0", "DC1", "DC2"}))
		g.Expect(VC1.VCenterIP).To(gmg.Equal("10.0.0.1"))
		g.Expect(VC1.SecretNamespace).To(gmg.Equal("kube-system"))
		g.Expect(VC1.SecretName).To(gmg.Equal("tenant1-secret"))

		VC2 := configStruct.Vcenter["10.0.0.2"]
		g.Expect(len(VC2.Datacenters)).To(gmg.Equal(1))
		g.Expect(VC2.Datacenters).To(gmg.BeComparableTo([]string{"DC3"}))
		g.Expect(VC2.VCenterIP).To(gmg.Equal("10.0.0.2"))
		g.Expect(VC2.IPFamilyPriority).To(gmg.Equal([]string{}))

		VC3 := configStruct.Vcenter["10.0.0.3"]
		g.Expect(len(VC3.Datacenters)).To(gmg.Equal(2))
		g.Expect(VC3.Datacenters).To(gmg.BeComparableTo([]string{"DC5", "DC6"}))
		g.Expect(VC3.IPFamilyPriority).To(gmg.Equal([]string{"ipv6"}))

		_, err = MarshalConfig(configStruct)
		g.Expect(err).ToNot(gmg.HaveOccurred())
	})

	invalidConfigTestCases := []struct {
		name         string
		input        string
		errSubstring string
	}{
		{
			"rubbish", invalidConfig, "expected section header",
		},
		{
			"bad port", invalidConfigWrongGlobalPort, "invalid global port parameter: parsed int bigger than zero",
		},
		{
			"bad vc port", invalidConfigWrongVCPort, "invalid port parameter for vc 10.0.0.3: parsed int bigger than zero",
		},
		{
			"empty", "", "vSphere config is empty",
		},
	}

	for _, tc := range invalidConfigTestCases {
		t.Run(fmt.Sprintf("test invalid config: %s", tc.name), func(t *testing.T) {
			g := gmg.NewWithT(t)
			_, err := ReadConfig([]byte(tc.input))
			g.Expect(err).To(gmg.MatchError(gmg.ContainSubstring(tc.errSubstring)))
		})
	}
}

func TestEmptyNodesOmittedInYAML(t *testing.T) {
	g := gmg.NewWithT(t)

	// Create a config with an empty Nodes struct
	config := &CPIConfig{
		CommonConfig: CommonConfig{
			Global: Global{
				User:     "testuser",
				Password: "testpass",
			},
		},
		Nodes: Nodes{}, // Empty struct, should be omitted
	}

	yamlOutput, err := MarshalConfig(config)
	g.Expect(err).ToNot(gmg.HaveOccurred())

	// Verify that "nodes:" does not appear in the output
	g.Expect(yamlOutput).ToNot(gmg.ContainSubstring("nodes:"))
}

func TestNodesWithValuesIncludedInYAML(t *testing.T) {
	g := gmg.NewWithT(t)

	// Create a config with a populated Nodes struct
	config := &CPIConfig{
		CommonConfig: CommonConfig{
			Global: Global{
				User:     "testuser",
				Password: "testpass",
			},
		},
		Nodes: Nodes{
			InternalNetworkSubnetCIDR: "192.168.1.0/24",
		},
	}

	yamlOutput, err := MarshalConfig(config)
	g.Expect(err).ToNot(gmg.HaveOccurred())

	// Verify that "nodes:" DOES appear when the struct has values
	g.Expect(yamlOutput).To(gmg.ContainSubstring("nodes:"))
	g.Expect(yamlOutput).To(gmg.ContainSubstring("internalNetworkSubnetCidr: 192.168.1.0/24"))
}
