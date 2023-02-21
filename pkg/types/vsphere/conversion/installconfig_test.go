package conversion

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

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
	installConfig.VSphere.DeprecatedDefaultDatastore = "test-datastore"
	installConfig.VSphere.DeprecatedDatacenter = "test-datacenter"

	return installConfig
}

func validLegacyIpiInstallConfig() *types.InstallConfig {
	installConfig := validLegacyUpiInstallConfig()

	installConfig.VSphere.DeprecatedCluster = "test-cluster"
	installConfig.VSphere.DeprecatedNetwork = "test-network"
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
		DeprecatedDatacenter:       "test-datacenter",
		DeprecatedDefaultDatastore: "test-datastore",
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
			Datacenters: []string{"test-datacenter"},
		}},
		FailureDomains: []vsphere.FailureDomain{{
			Name:   "generated-failure-domain",
			Region: "generated-region",
			Zone:   "generated-zone",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "test-datacenter",
				ComputeCluster: "",
				Networks:       []string{""},
				Datastore:      "/test-datacenter/datastore/test-datastore",
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
		DeprecatedDatacenter:       "test-datacenter",
		DeprecatedDefaultDatastore: "test-datastore",
		DeprecatedFolder:           "",
		DeprecatedCluster:          "test-cluster",
		DeprecatedResourcePool:     "",
		ClusterOSImage:             "",
		DeprecatedAPIVIP:           "",
		APIVIPs:                    []string{"192.168.111.0"},
		DeprecatedIngressVIP:       "",
		IngressVIPs:                []string{"192.168.111.1"},
		DefaultMachinePlatform:     nil,
		DeprecatedNetwork:          "test-network",
		DiskType:                   "",
		VCenters: []vsphere.VCenter{{
			Server:      "test-server",
			Port:        443,
			Username:    "test-username",
			Password:    "test-password",
			Datacenters: []string{"test-datacenter"},
		}},
		FailureDomains: []vsphere.FailureDomain{{
			Name:   "generated-failure-domain",
			Region: "generated-region",
			Zone:   "generated-zone",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "test-datacenter",
				ComputeCluster: "/test-datacenter/host/test-cluster",
				Networks:       []string{"test-network"},
				Datastore:      "/test-datacenter/datastore/test-datastore",
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
				"test-datacenter",
				"test-datacenter4",
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
				Datacenter:     "test-datacenter",
				ComputeCluster: "/test-datacenter1/host/test-computecluster1",
				Networks:       []string{"network1"},
				Datastore:      "/test-datacenter/datastore/datastore1",
			},
		},
		{
			Name:   "region-2-zone-2a",
			Region: "region-2",
			Zone:   "zone-2a",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "test-datacenter",
				ComputeCluster: "/test-datacenter1/host/test-computecluster2",
				Networks:       []string{"network2"},
				Datastore:      "/test-datacenter/datastore/datastore2",
			},
		},
		{
			Name:   "region-3-zone-3a",
			Region: "region-3",
			Zone:   "zone-3a",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "test-datacenter",
				ComputeCluster: "/test-datacenter1/host/test-computecluster3",
				Networks:       []string{"network3"},
				Datastore:      "/test-datacenter/datastore/datastore3",
			},
		},
		{
			Name:   "region-4-zone-4a",
			Region: "region-4",
			Zone:   "zone-4a",
			Server: "test-server",
			Topology: vsphere.Topology{
				Datacenter:     "test-datacenter4",
				ComputeCluster: "/test-datacenter4/host/test-computecluster4",
				Networks:       []string{"network4"},
				Datastore:      "/test-datacenter4/datastore/datastore4",
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
				Datacenter:     "",
				ComputeCluster: "/test-datacenter1/host/test-computecluster1",
				Networks:       []string{"network1"},
				Datastore:      "datastore1",
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
				Datacenter:     "",
				ComputeCluster: "/test-datacenter1/host/test-computecluster2",
				Networks:       []string{"network2"},
				Datastore:      "datastore2",
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
				Datacenter:     "",
				ComputeCluster: "/test-datacenter1/host/test-computecluster3",
				Networks:       []string{"network3"},
				Datastore:      "datastore3",
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
				Datacenter:     "test-datacenter4",
				ComputeCluster: "/test-datacenter4/host/test-computecluster4",
				Networks:       []string{"network4"},
				Datastore:      "datastore4",
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
			Datacenters: []string{"test-datacenter1", "test-datacenter4"},
		}},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "region-1-zone-1a",
				Region: "region-1",
				Zone:   "zone-1a",
				Server: "test-server",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter1",
					ComputeCluster: "/test-datacenter1/host/test-computecluster1",
					Networks:       []string{"network1"},
					Datastore:      "/test-datacenter1/datstore/test-datastore1",
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
					Datacenter:     "test-datacenter1",
					ComputeCluster: "/test-datacenter1/host/test-computecluster2",
					Networks:       []string{"network2"},
					Datastore:      "/test-datacenter1/datastore/test-datastore2",
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
					Datacenter:     "test-datacenter1",
					ComputeCluster: "/test-datacenter1/host/test-computecluster3",
					Networks:       []string{"network3"},
					Datastore:      "/test-datacenter1/datastore/test-datastore3",
					ResourcePool:   "/test-datacenter1/host/test-computecluster3/Resources/test-resourcepool4",
					Folder:         "",
				},
			},
			{
				Name:   "region-4-zone-4a",
				Region: "region-4",
				Zone:   "zone-4a",
				Server: "",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter4",
					ComputeCluster: "/test-datacenter4/host/test-computecluster4",
					Networks:       []string{"network4"},
					Datastore:      "/test-datacenter4/datastore/datastore4",
					ResourcePool:   "",
					Folder:         "/test-datacenter4/vm/test-folder4",
				},
			},
		},
	}

	return installConfig
}

func TestConvertInstallConfig(t *testing.T) {
	logger, hook := test.NewNullLogger()
	localLogger = logger
	tests := []struct {
		name                string
		actualInstallConfig func() *types.InstallConfig
		expectInstallConfig func() *types.InstallConfig
		expectWarn          []string
		expectLevel         logrus.Level
	}{
		{
			name:                "legacy upi conversion",
			actualInstallConfig: validLegacyUpiInstallConfig,
			expectInstallConfig: convertedLegacyUpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now depreciated please use vcenters",
				"vsphere topology fields are now depreciated please use failureDomains",
				"datastore as a non-path is now depreciated please use the form: /test-datacenter/datastore/test-datastore",
			},
		},
		{
			name:                "legacy ipi conversion",
			actualInstallConfig: validLegacyIpiInstallConfig,
			expectInstallConfig: convertedLegacyIpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now depreciated please use vcenters",
				"vsphere topology fields are now depreciated please use failureDomains",
				"computeCluster as a non-path is now depreciated please use the form: /test-datacenter/host/test-cluster",
				"datastore as a non-path is now depreciated please use the form: /test-datacenter/datastore/test-datastore",
			},
		},
		{
			name:                "legacy zonal ipi conversion",
			actualInstallConfig: validLegacyIpiZonalInstallConfig,
			expectInstallConfig: convertedLegacyIpiZonalInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn: []string{
				"vsphere authentication fields are now depreciated please use vcenters",
				"datastore as a non-path is now depreciated please use the form: /test-datacenter/datastore/datastore1",
				"datastore as a non-path is now depreciated please use the form: /test-datacenter/datastore/datastore2",
				"datastore as a non-path is now depreciated please use the form: /test-datacenter/datastore/datastore3",
				"datastore as a non-path is now depreciated please use the form: /test-datacenter4/datastore/datastore4",
			},
		},
		{
			name:                "zonal ipi",
			actualInstallConfig: validIpiInstallConfig,
			expectInstallConfig: validIpiInstallConfig,
			expectLevel:         logrus.WarnLevel,
			expectWarn:          []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)

			actual := tt.actualInstallConfig()
			expect := tt.expectInstallConfig()

			if err := ConvertInstallConfig(actual); err != nil {
				g.Expect(err).Should(gomega.Equal(nil))
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
			expectWarn:     `computeCluster as a non-path is now depreciated please use the form: /DC1/host/C1`,
			expectCluster:  "/DC1/host/C1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setComputeClusterPath(tt.computeCluster, tt.datacenter)
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
			expectWarn:     `datastore as a non-path is now depreciated please use the form: /DC1/datastore/DS1`,
			expectDatstore: "/DC1/datastore/DS1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setDatastorePath(tt.datastore, tt.datacenter)
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
			expectWarn:   `folder as a non-path is now depreciated please use the form: /DC1/vm/Folder1`,
			expectFolder: "/DC1/vm/Folder1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setFolderPath(tt.folder, tt.datacenter)
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
