package conversion

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var localLogger = logrus.New()

// ConvertInstallConfig modifies a given platform spec for the new requirements.
func ConvertInstallConfig(config *types.InstallConfig) error {
	platform := config.Platform.VSphere

	// Scenario: IPI or 4.12 Zonal IPI w/o vcenters defined
	if len(platform.VCenters) == 0 {
		createVCenters(platform)

		// Scenario: 4.12 Zonal IPI
		if len(platform.FailureDomains) > 0 {
			for i := range platform.FailureDomains {
				if platform.FailureDomains[i].Topology.Datacenter == "" {
					platform.FailureDomains[i].Topology.Datacenter = platform.DeprecatedDatacenter
				}
				if platform.FailureDomains[i].Server == "" {
					// Assumption: by the time it is possible to use multiple vcenters
					// it will be past 4.15
					// so this conversion can be removed.
					platform.FailureDomains[i].Server = platform.VCenters[0].Server
				}
			}
		}
	}

	// Scenario: Fields are not paths
	if len(platform.FailureDomains) > 0 {
		for i := range platform.FailureDomains {
			platform.FailureDomains[i].Topology.ComputeCluster = setComputeClusterPath(platform.FailureDomains[i].Topology.ComputeCluster,
				platform.FailureDomains[i].Topology.Datacenter)

			platform.FailureDomains[i].Topology.Datastore = setDatastorePath(platform.FailureDomains[i].Topology.Datastore,
				platform.FailureDomains[i].Topology.Datacenter)

			platform.FailureDomains[i].Topology.Folder = setFolderPath(platform.FailureDomains[i].Topology.Folder,
				platform.FailureDomains[i].Topology.Datacenter)
		}
	}

	// Scenario: legacy UPI or IPI
	if len(platform.FailureDomains) == 0 {
		localLogger.Warn("vsphere topology fields are now deprecated; please use failureDomains")

		platform.FailureDomains = make([]vsphere.FailureDomain, 1)
		platform.FailureDomains[0].Name = "generated-failure-domain"
		platform.FailureDomains[0].Server = platform.VCenters[0].Server
		platform.FailureDomains[0].Region = "generated-region"
		platform.FailureDomains[0].Zone = "generated-zone"

		platform.FailureDomains[0].Topology.Datacenter = platform.DeprecatedDatacenter
		platform.FailureDomains[0].Topology.ResourcePool = platform.DeprecatedResourcePool
		platform.FailureDomains[0].Topology.ComputeCluster = setComputeClusterPath(platform.DeprecatedCluster, platform.DeprecatedDatacenter)
		platform.FailureDomains[0].Topology.Networks = make([]string, 1)
		platform.FailureDomains[0].Topology.Networks[0] = platform.DeprecatedNetwork
		platform.FailureDomains[0].Topology.Datastore = setDatastorePath(platform.DeprecatedDefaultDatastore, platform.DeprecatedDatacenter)
		platform.FailureDomains[0].Topology.Folder = setFolderPath(platform.DeprecatedFolder, platform.DeprecatedDatacenter)
	}

	return nil
}

func setComputeClusterPath(cluster, datacenter string) string {
	if cluster != "" && !strings.HasPrefix(cluster, "/") {
		localLogger.Warnf("computeCluster as a non-path is now deprecated; please use the form: /%s/host/%s", datacenter, cluster)
		return fmt.Sprintf("/%s/host/%s", datacenter, cluster)
	}
	return cluster
}

func setDatastorePath(datastore, datacenter string) string {
	if datastore != "" && !strings.HasPrefix(datastore, "/") {
		localLogger.Warnf("datastore as a non-path is now deprecated; please use the form: /%s/datastore/%s", datacenter, datastore)
		return fmt.Sprintf("/%s/datastore/%s", datacenter, datastore)
	}
	return datastore
}

func setFolderPath(folder, datacenter string) string {
	if folder != "" && !strings.HasPrefix(folder, "/") {
		localLogger.Warnf("folder as a non-path is now deprecated; please use the form: /%s/vm/%s", datacenter, folder)
		return fmt.Sprintf("/%s/vm/%s", datacenter, folder)
	}
	return folder
}

func createVCenters(platform *vsphere.Platform) {
	localLogger.Warn("vsphere authentication fields are now deprecated; please use vcenters")

	platform.VCenters = make([]vsphere.VCenter, 1)
	platform.VCenters[0].Server = platform.DeprecatedVCenter
	platform.VCenters[0].Username = platform.DeprecatedUsername
	platform.VCenters[0].Password = platform.DeprecatedPassword
	platform.VCenters[0].Port = 443

	if platform.DeprecatedDatacenter != "" {
		platform.VCenters[0].Datacenters = append(platform.VCenters[0].Datacenters, platform.DeprecatedDatacenter)
	}

	// Scenario: Zonal IPI w/o vcenters defined
	// Confirms the list of datacenters from FailureDomains are updated
	// in vcenters[0].datacenters
	for _, failureDomain := range platform.FailureDomains {
		found := false
		if failureDomain.Topology.Datacenter != "" {
			for _, dc := range platform.VCenters[0].Datacenters {
				if dc == failureDomain.Topology.Datacenter {
					found = true
				}
			}

			if !found {
				platform.VCenters[0].Datacenters = append(platform.VCenters[0].Datacenters, failureDomain.Topology.Datacenter)
			}
		}
	}
}
