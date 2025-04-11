package conversion

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/find"

	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	validCIDR = "10.0.0.0/16"
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				APIVIPs:     []string{"192.168.111.0"},
				IngressVIPs: []string{"192.168.111.1"},
			},
		},
	}
}

func validLegacyUpiInstallConfig() *types.InstallConfig {
	installConfig := validInstallConfig()

	// The fields below are the original UPI required
	installConfig.VSphere.DeprecatedVCenter = "test-server"
	installConfig.VSphere.DeprecatedUsername = "test-username"
	installConfig.VSphere.DeprecatedPassword = "test-password"
	installConfig.VSphere.DeprecatedDefaultDatastore = "LocalDS_0"
	installConfig.VSphere.DeprecatedDatacenter = "DC0"

	return installConfig
}

func validLegacyIpiInstallConfig() *types.InstallConfig {
	installConfig := validLegacyUpiInstallConfig()

	installConfig.VSphere.DeprecatedCluster = "DC0_C0"
	installConfig.VSphere.DeprecatedNetwork = "DC0_DVPG0"
	installConfig.VSphere.DeprecatedFolder = ""
	installConfig.VSphere.DeprecatedResourcePool = ""

	return installConfig
}

func convertedLegacyUpiInstallConfig() *types.InstallConfig {
	installConfig := validInstallConfig()
	installConfig.VSphere = &vsphere.Platform{
		DeprecatedVCenter:          "test-server",
		DeprecatedUsername:         "test-username",
		DeprecatedPassword:         "test-password",
		DeprecatedDatacenter:       "DC0",
		DeprecatedDefaultDatastore: "LocalDS_0",
		DeprecatedFolder:           "",
		DeprecatedCluster:          "",
		DeprecatedResourcePool:     "",
		ClusterOSImage:             "",
		DeprecatedAPIVIP:           "",
		APIVIPs:                    []string{"192.168.111.0"},
		DeprecatedIngressVIP:       "",
		IngressVIPs:                []string{"192.168.111.1"},
		DefaultMachinePlatform:     nil,
		DeprecatedNetwork:          "",
		DiskType:                   "",
		VCenters: []vsphere.VCenter{{
			Server:      "test-server",
			Port:        443,
			Username:    "test-username",
			Password:    "test-password",
			Datacenters: []string{"DC0"},
		}},
		FailureDomains: []vsphere.FailureDomain{{
			Name:   "generated-failure-domain",
			Region: "generated-region",
			Zone:   "generated-zone",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "DC0",
				ComputeCluster: "",
				Networks:       []string{""},
				Datastore:      "/DC0/datastore/LocalDS_0",
				ResourcePool:   "",
				Folder:         "",
			},
		}},
	}

	return installConfig
}
func convertedLegacyIpiInstallConfig() *types.InstallConfig {
	installConfig := validInstallConfig()
	installConfig.VSphere = &vsphere.Platform{
		DeprecatedVCenter:          "test-server",
		DeprecatedUsername:         "test-username",
		DeprecatedPassword:         "test-password",
		DeprecatedDatacenter:       "DC0",
		DeprecatedDefaultDatastore: "LocalDS_0",
		DeprecatedFolder:           "",
		DeprecatedCluster:          "DC0_C0",
		DeprecatedResourcePool:     "",
		ClusterOSImage:             "",
		DeprecatedAPIVIP:           "",
		APIVIPs:                    []string{"192.168.111.0"},
		DeprecatedIngressVIP:       "",
		IngressVIPs:                []string{"192.168.111.1"},
		DefaultMachinePlatform:     nil,
		DeprecatedNetwork:          "DC0_DVPG0",
		DiskType:                   "",
		VCenters: []vsphere.VCenter{{
			Server:      "test-server",
			Port:        443,
			Username:    "test-username",
			Password:    "test-password",
			Datacenters: []string{"DC0"},
		}},
		FailureDomains: []vsphere.FailureDomain{{
			Name:   "generated-failure-domain",
			Region: "generated-region",
			Zone:   "generated-zone",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "DC0",
				ComputeCluster: "/DC0/host/DC0_C0",
				Networks:       []string{"DC0_DVPG0"},
				Datastore:      "/DC0/datastore/LocalDS_0",
				ResourcePool:   "",
				Folder:         "",
			},
		}},
	}

	return installConfig
}

func convertedLegacyIpiZonalInstallConfig() *types.InstallConfig {
	installConfig := convertedLegacyIpiInstallConfig()

	installConfig.VSphere.VCenters = []vsphere.VCenter{
		{
			Server:   "test-server",
			Port:     443,
			Username: "test-username",
			Password: "test-password",
			Datacenters: []string{
				"DC0",
				"DC1",
				"DC2",
				"DC3",
			},
		},
	}

	installConfig.VSphere.FailureDomains = []vsphere.FailureDomain{
		{
			Name:   "region-1-zone-1a",
			Region: "region-1",
			Zone:   "zone-1a",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "DC0",
				ComputeCluster: "/DC0/host/DC0_C0",
				Networks:       []string{"DC0_DVPG0"},
				Datastore:      "/DC0/datastore/LocalDS_0",
			},
		},
		{
			Name:   "region-2-zone-2a",
			Region: "region-2",
			Zone:   "zone-2a",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "DC1",
				ComputeCluster: "/DC1/host/DC1_C0",
				Networks:       []string{"DC1_DVPG0"},
				Datastore:      "/DC1/datastore/LocalDS_0",
			},
		},
		{
			Name:   "region-3-zone-3a",
			Region: "region-3",
			Zone:   "zone-3a",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "DC2",
				ComputeCluster: "/DC2/host/DC2_C0",
				Networks:       []string{"DC2_DVPG0"},
				Datastore:      "/DC2/datastore/LocalDS_0",
			},
		},
		{
			Name:   "region-4-zone-4a",
			Region: "region-4",
			Zone:   "zone-4a",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "DC3",
				ComputeCluster: "/DC3/host/DC3_C0",
				Networks:       []string{"DC3_DVPG0"},
				Datastore:      "/DC3/datastore/LocalDS_0",
			},
		},
	}

	return installConfig
}

// validLegacyIpiZonalInstallConfig defines an install-config
// that was valid during 4.12 tech preview of zonal.
func validLegacyIpiZonalInstallConfig() *types.InstallConfig {
	installConfig := validLegacyIpiInstallConfig()

	installConfig.VSphere.FailureDomains = []vsphere.FailureDomain{
		{
			Name:   "region-1-zone-1a",
			Region: "region-1",
			Zone:   "zone-1a",
			Server: "",
			Topology: vsphere.Topology{
				Datacenter:     "DC0",
				ComputeCluster: "/DC0/host/DC0_C0",
				Networks:       []string{"DC0_DVPG0"},
				Datastore:      "LocalDS_0",
				ResourcePool:   "",
				Folder:         "",
			},
		},
		{
			Name:   "region-2-zone-2a",
			Region: "region-2",
			Zone:   "zone-2a",
			Server: "",
			Topology: vsphere.Topology{
				Datacenter:     "DC1",
				ComputeCluster: "/DC1/host/DC1_C0",
				Networks:       []string{"DC1_DVPG0"},
				Datastore:      "LocalDS_0",
				ResourcePool:   "",
				Folder:         "",
			},
		},
		{
			Name:   "region-3-zone-3a",
			Region: "region-3",
			Zone:   "zone-3a",
			Server: "",
			Topology: vsphere.Topology{
				Datacenter:     "DC2",
				ComputeCluster: "/DC2/host/DC2_C0",
				Networks:       []string{"DC2_DVPG0"},
				Datastore:      "LocalDS_0",
				ResourcePool:   "",
				Folder:         "",
			},
		},
		{
			Name:   "region-4-zone-4a",
			Region: "region-4",
			Zone:   "zone-4a",
			Server: "",
			Topology: vsphere.Topology{
				Datacenter:     "DC3",
				ComputeCluster: "/DC3/host/DC3_C0",
				Networks:       []string{"DC3_DVPG0"},
				Datastore:      "LocalDS_0",
				ResourcePool:   "",
				Folder:         "",
			},
		},
	}

	return installConfig
}

func validIpiInstallConfig() *types.InstallConfig {
	installConfig := validInstallConfig()

	installConfig.VSphere = &vsphere.Platform{
		DeprecatedVCenter:          "",
		DeprecatedUsername:         "",
		DeprecatedPassword:         "",
		DeprecatedDatacenter:       "",
		DeprecatedDefaultDatastore: "",
		DeprecatedFolder:           "",
		DeprecatedCluster:          "",
		DeprecatedResourcePool:     "",
		ClusterOSImage:             "",
		DeprecatedAPIVIP:           "",
		APIVIPs:                    []string{"192.168.111.0"},
		DeprecatedIngressVIP:       "",
		IngressVIPs:                []string{"192.168.111.1"},
		DefaultMachinePlatform:     nil,
		DeprecatedNetwork:          "",
		DiskType:                   "",
		VCenters: []vsphere.VCenter{{
			Server:      "test-server",
			Port:        443,
			Username:    "test-username",
			Password:    "test-password",
			Datacenters: []string{"DC0", "DC1", "DC2", "DC3"},
		}},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "region-1-zone-1a",
				Region: "region-1",
				Zone:   "zone-1a",
				Server: "test-server",
				Topology: vsphere.Topology{
					Datacenter:     "DC0",
					ComputeCluster: "/DC0/host/DC0_C0",
					Networks:       []string{"network1"},
					Datastore:      "/DC0/datstore/test-datastore1",
					ResourcePool:   "",
					Folder:         "",
				},
			},
			{
				Name:   "region-2-zone-2a",
				Region: "region-2",
				Zone:   "zone-2a",
				Server: "test-server",
				Topology: vsphere.Topology{
					Datacenter:     "DC1",
					ComputeCluster: "/DC1/host/DC1_C0",
					Networks:       []string{"network2"},
					Datastore:      "/DC0/datastore/test-datastore2",
					ResourcePool:   "",
					Folder:         "",
				},
			},
			{
				Name:   "region-3-zone-3a",
				Region: "region-3",
				Zone:   "zone-3a",
				Server: "test-server",
				Topology: vsphere.Topology{
					Datacenter:     "DC2",
					ComputeCluster: "/DC2/host/test-computecluster3",
					Networks:       []string{"network3"},
					Datastore:      "/DC2/datastore/test-datastore3",
					ResourcePool:   "/DC2/host/test-computecluster3/Resources/test-resourcepool4",
					Folder:         "",
				},
			},
			{
				Name:   "region-4-zone-4a",
				Region: "region-4",
				Zone:   "zone-4a",
				Server: "test-server",
				Topology: vsphere.Topology{
					Datacenter:     "DC3",
					ComputeCluster: "/DC3/host/DC3_C0",
					Networks:       []string{"network4"},
					Datastore:      "/DC3/datastore/LocalDS_0",
					ResourcePool:   "",
					Folder:         "/DC3/vm/",
				},
			},
		},
	}

	return installConfig
}

func TestConvertInstallConfig(t *testing.T) {
	logger, hook := test.NewNullLogger()

	sim, err := mock.StartSimulator(true)

	if err != nil {
		assert.NoError(t, err)
	}

	localLogger = logger
	tests := []struct {
		name                string
		actualInstallConfig func() *types.InstallConfig
		expectInstallConfig func() *types.InstallConfig
		expectWarn          []string
		expectLevel         logrus.Level
		simulated           bool
	}{
		{
			name:                "legacy upi conversion with simulator",
			actualInstallConfig: validLegacyUpiInstallConfig,
			expectInstallConfig: convertedLegacyUpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now deprecated; please use vcenters",
				"vsphere topology fields are now deprecated; please use failureDomains",
				"datastore as a non-path is now deprecated; please use the discovered form: /DC0/datastore/LocalDS_0",
			},
			simulated: true,
		},
		{
			name:                "legacy ipi conversion with simulator",
			actualInstallConfig: validLegacyIpiInstallConfig,
			expectInstallConfig: convertedLegacyIpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now deprecated; please use vcenters",
				"vsphere topology fields are now deprecated; please use failureDomains",
				"computeCluster as a non-path is now deprecated; please use the discovered form: /DC0/host/DC0_C0",
				"datastore as a non-path is now deprecated; please use the discovered form: /DC0/datastore/LocalDS_0",
			},
			simulated: true,
		},
		{
			name:                "zonal ipi with simulator",
			actualInstallConfig: validIpiInstallConfig,
			expectInstallConfig: validIpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn:          []string{""},
			simulated:           true,
		},
		{
			name: "legacy zonal ipi conversion with simulator",
			actualInstallConfig: func() *types.InstallConfig {
				installConfig := validLegacyIpiZonalInstallConfig()
				return installConfig
			},
			expectInstallConfig: convertedLegacyIpiZonalInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now deprecated; please use vcenters",
				"datastore as a non-path is now deprecated; please use the discovered form: /DC0/datastore/LocalDS_0",
				"datastore as a non-path is now deprecated; please use the discovered form: /DC1/datastore/LocalDS_0",
				"datastore as a non-path is now deprecated; please use the discovered form: /DC2/datastore/LocalDS_0",
				"datastore as a non-path is now deprecated; please use the discovered form: /DC3/datastore/LocalDS_0",
				"vsphere topology fields are now deprecated; please avoid using failureDomains and the vsphere topology fields together",
			},
			simulated: true,
		},
		{
			name:                "legacy upi conversion without simulator",
			actualInstallConfig: validLegacyUpiInstallConfig,
			expectInstallConfig: convertedLegacyUpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now deprecated; please use vcenters",
				"vsphere topology fields are now deprecated; please use failureDomains",
				"datastore as a non-path is now deprecated; please use the joined form: /DC0/datastore/LocalDS_0",
			},
			simulated: false,
		},
		{
			name:                "legacy ipi conversion without simulator",
			actualInstallConfig: validLegacyIpiInstallConfig,
			expectInstallConfig: convertedLegacyIpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now deprecated; please use vcenters",
				"vsphere topology fields are now deprecated; please use failureDomains",
				"computeCluster as a non-path is now deprecated; please use the joined form: /DC0/host/DC0_C0",
				"datastore as a non-path is now deprecated; please use the joined form: /DC0/datastore/LocalDS_0",
			},
			simulated: false,
		},
		{
			name:                "zonal ipi without simulator",
			actualInstallConfig: validIpiInstallConfig,
			expectInstallConfig: validIpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn:          []string{""},
			simulated:           false,
		},
		{
			name: "legacy zonal ipi conversion without simulator",
			actualInstallConfig: func() *types.InstallConfig {
				installConfig := validLegacyIpiZonalInstallConfig()
				return installConfig
			},
			expectInstallConfig: convertedLegacyIpiZonalInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now deprecated; please use vcenters",
				"datastore as a non-path is now deprecated; please use the joined form: /DC0/datastore/LocalDS_0",
				"datastore as a non-path is now deprecated; please use the joined form: /DC1/datastore/LocalDS_0",
				"datastore as a non-path is now deprecated; please use the joined form: /DC2/datastore/LocalDS_0",
				"datastore as a non-path is now deprecated; please use the joined form: /DC3/datastore/LocalDS_0",
				"vsphere topology fields are now deprecated; please avoid using failureDomains and the vsphere topology fields together",
			},
			simulated: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var finder *find.Finder
			g := gomega.NewWithT(t)

			actual := tt.actualInstallConfig()
			expect := tt.expectInstallConfig()

			platform := actual.Platform.VSphere

			if tt.simulated {
				finder, err = mock.GetFinder(sim)
				if err != nil {
					assert.NoError(t, err)
				}
			} else {
				finder = nil
			}

			// This section is duplication from ConvertInstallConfig()
			// which makes it easier to deal with the simulator and finder object.
			finders := make(map[string]*find.Finder)

			fixNoVCentersScenario(platform)
			finders[platform.VCenters[0].Server] = finder
			if err != nil {
				assert.NoError(t, err)
			}
			err = fixTechPreviewZonalFailureDomainsScenario(platform, finders)
			if err != nil {
				assert.NoError(t, err)
			}
			err = fixLegacyPlatformScenario(platform, finders)
			if err != nil {
				assert.NoError(t, err)
			}

			if tt.expectWarn[0] != "" {
				entries := hook.AllEntries()
				assert.NotEmpty(t, entries)
				for i, entry := range entries {
					g.Expect(entry.Level).To(gomega.Equal(tt.expectLevel))
					g.Expect(entry.Message).To(gomega.Equal(tt.expectWarn[i]))
				}
			}
			hook.Reset()

			g.Expect(actual).Should(gomega.BeComparableTo(expect))
		})
	}
}

func Test_setComputeClusterPath(t *testing.T) {
	logger, hook := test.NewNullLogger()
	localLogger = logger

	tests := []struct {
		name           string
		computeCluster string
		datacenter     string
		expectLevel    string
		expectWarn     string
		expectCluster  string
	}{
		{
			name:           "correct cluster path",
			computeCluster: "/DC1/host/C1",
			datacenter:     "DC1",
			expectLevel:    ``,
			expectWarn:     ``,
			expectCluster:  "/DC1/host/C1",
		},
		{
			name:           "cluster name, not path",
			computeCluster: "C1",
			datacenter:     "DC1",
			expectLevel:    `warning`,
			expectWarn:     `computeCluster as a non-path is now deprecated; please use the joined form: /DC1/host/C1`,
			expectCluster:  "/DC1/host/C1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetObjectPath(nil, "host", tt.computeCluster, tt.datacenter)

			assert.NoError(t, err)

			assert.Equal(t, tt.expectCluster, got)
			if tt.expectWarn != "" {
				entries := hook.AllEntries()
				assert.NotEmpty(t, entries)
				for _, e := range entries {
					assert.Equal(t, tt.expectLevel, e.Level.String())
					assert.Regexp(t, tt.expectWarn, e.Message)
				}
			}
			hook.Reset()
		})
	}
}

func Test_setDatastorePath(t *testing.T) {
	logger, hook := test.NewNullLogger()
	localLogger = logger

	tests := []struct {
		name           string
		datastore      string
		datacenter     string
		expectLevel    string
		expectWarn     string
		expectDatstore string
	}{
		{
			name:           "correct datastore path",
			datastore:      "/DC1/datastore/DS1",
			datacenter:     "DC1",
			expectLevel:    ``,
			expectWarn:     ``,
			expectDatstore: "/DC1/datastore/DS1",
		},
		{
			name:           "cluster name, not path",
			datastore:      "DS1",
			datacenter:     "DC1",
			expectLevel:    `warning`,
			expectWarn:     `datastore as a non-path is now deprecated; please use the joined form: /DC1/datastore/DS1`,
			expectDatstore: "/DC1/datastore/DS1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetObjectPath(nil, "datastore", tt.datastore, tt.datacenter)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectDatstore, got)
			if tt.expectWarn != "" {
				entries := hook.AllEntries()
				assert.NotEmpty(t, entries)
				for _, e := range entries {
					assert.Equal(t, tt.expectLevel, e.Level.String())
					assert.Regexp(t, tt.expectWarn, e.Message)
				}
			}
			hook.Reset()
		})
	}
}

func Test_setFolderPath(t *testing.T) {
	logger, hook := test.NewNullLogger()
	localLogger = logger

	tests := []struct {
		name         string
		folder       string
		datacenter   string
		expectLevel  string
		expectWarn   string
		expectFolder string
	}{
		{
			name:         "correct folder path",
			folder:       "/DC1/vm/Folder1",
			datacenter:   "DC1",
			expectLevel:  ``,
			expectWarn:   ``,
			expectFolder: "/DC1/vm/Folder1",
		},
		{
			name:         "folder as name not path",
			folder:       "Folder1",
			datacenter:   "DC1",
			expectLevel:  `warning`,
			expectWarn:   `folder as a non-path is now deprecated; please use the joined form: /DC1/vm/Folder1`,
			expectFolder: "/DC1/vm/Folder1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetObjectPath(nil, "vm", tt.folder, tt.datacenter)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectFolder, got)
			if tt.expectWarn != "" {
				entries := hook.AllEntries()
				assert.NotEmpty(t, entries)
				for _, e := range entries {
					assert.Equal(t, tt.expectLevel, e.Level.String())
					assert.Regexp(t, tt.expectWarn, e.Message)
				}
			}
			hook.Reset()
		})
	}
}
