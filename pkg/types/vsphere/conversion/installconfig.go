package conversion

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/soap"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var localLogger = logrus.New()

const (
	// GeneratedFailureDomainName is a placeholder name when one wasn't provided.
	GeneratedFailureDomainName string = "generated-failure-domain"
	// GeneratedFailureDomainRegion is a placeholder region when one wasn't provided.
	GeneratedFailureDomainRegion string = "generated-region"
	// GeneratedFailureDomainZone is a placeholder zone when one wasn't provided.
	GeneratedFailureDomainZone string = "generated-zone"
)

// GetFinder connects to vCenter via SOAP and returns the Finder object if the SOAP
// connection is successful. If the connection fails it returns nil.
// Errors are mostly ignored to support AI and agent installers.
func GetFinder(server, username, password string) (*find.Finder, error) {
	var finder *find.Finder

	if server != "" && password != "" && username != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		u, err := soap.ParseURL(server)
		if err != nil {
			return nil, err
		}
		u.User = url.UserPassword(username, password)

		client, err := govmomi.NewClient(ctx, u, false)
		if err != nil {
			// If bogus authentication is provided in the scenario of AI or assisted
			// just provide warning message. If this is IPI or UPI validation will
			// catch and halt on incorrect authentication.

			localLogger.Debugf("this can be safely ignored if non-deprecated platform spec fields are used"+
				"or installing via UPI, Assisted or Agent Installer. "+
				"Conversion of deprecated platform spec fields cannot continue without vCenter %s access, error: %v",
				server, err)

			return nil, nil
		}
		finder = find.NewFinder(client.Client, true)
	}

	return finder, nil
}

func findViaPathOrName(finder *find.Finder, objectPath, objectFindPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	elements, err := finder.ManagedObjectListChildren(ctx, objectFindPath)
	if err != nil {
		return "", err
	}

	for _, e := range elements {
		if e.Path == objectPath {
			return objectPath, nil
		}

		if path.Base(e.Path) == path.Base(objectPath) {
			return e.Path, nil
		}
	}
	return "", errors.New("unable to find object")
}

// fixNoVCentersScenario this function creates the VCenters slice
// with existing legacy vcenter authentication and configuration.
func fixNoVCentersScenario(platform *vsphere.Platform) {
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
	} else if platform.DeprecatedVCenter != "" {
		localLogger.Warn("vcenter field is deprecated please avoid using both vcenter and vcenters fields together")
	}
}

func fixTechPreviewZonalFailureDomainsScenario(platform *vsphere.Platform, finders map[string]*find.Finder) error {
	if len(platform.FailureDomains) > 0 {
		var err error

		for i := range platform.FailureDomains {
			computeCluster := platform.FailureDomains[i].Topology.ComputeCluster
			datastore := platform.FailureDomains[i].Topology.Datastore
			folder := platform.FailureDomains[i].Topology.Folder
			datacenter := platform.FailureDomains[i].Topology.Datacenter
			vCenter := platform.FailureDomains[i].Server
			fdName := platform.FailureDomains[i].Name

			finder, ok := finders[vCenter]
			if !ok {
				// This is when invalid config happens.  There is a check later in cycle that will print it out.  For now,
				// lets just log warning and return.
				localLogger.Warnf("unable to find finder for vCenter %v in order to do upconvert", vCenter)
				return nil
			}
			platform.FailureDomains[i].Topology.ComputeCluster, err = SetObjectPath(finder, "host", computeCluster, datacenter)
			if err != nil {
				return fmt.Errorf("unable to SetObjectPath for compute cluster of failure domain %s: %w", fdName, err)
			}

			platform.FailureDomains[i].Topology.Datastore, err = SetObjectPath(finder, "datastore", datastore, datacenter)
			if err != nil {
				return fmt.Errorf("unable to SetObjectPath for datastore of failure domain %s: %w", fdName, err)
			}

			platform.FailureDomains[i].Topology.Folder, err = SetObjectPath(finder, "vm", folder, datacenter)
			if err != nil {
				return fmt.Errorf("unable to SetObjectPath for folder of failure domain %s: %w", fdName, err)
			}
		}
	}
	return nil
}

func fixLegacyPlatformScenario(platform *vsphere.Platform, finders map[string]*find.Finder) error {
	if len(platform.FailureDomains) == 0 {
		var err error
		localLogger.Warn("vsphere topology fields are now deprecated; please use failureDomains")

		platform.FailureDomains = make([]vsphere.FailureDomain, 1)
		platform.FailureDomains[0].Name = GeneratedFailureDomainName
		platform.FailureDomains[0].Server = platform.VCenters[0].Server
		platform.FailureDomains[0].Region = GeneratedFailureDomainRegion
		platform.FailureDomains[0].Zone = GeneratedFailureDomainZone

		platform.FailureDomains[0].Topology.Datacenter = platform.DeprecatedDatacenter
		platform.FailureDomains[0].Topology.ResourcePool = platform.DeprecatedResourcePool
		platform.FailureDomains[0].Topology.Networks = make([]string, 1)
		platform.FailureDomains[0].Topology.Networks[0] = platform.DeprecatedNetwork

		finder, ok := finders[platform.FailureDomains[0].Server]
		if !ok {
			return fmt.Errorf("unable to find finder for vCenter %v", platform.FailureDomains[0].Server)
		}
		platform.FailureDomains[0].Topology.ComputeCluster, err = SetObjectPath(finder, "host", platform.DeprecatedCluster, platform.DeprecatedDatacenter)
		if err != nil {
			return err
		}

		platform.FailureDomains[0].Topology.Datastore, err = SetObjectPath(finder, "datastore", platform.DeprecatedDefaultDatastore, platform.DeprecatedDatacenter)
		if err != nil {
			return err
		}

		platform.FailureDomains[0].Topology.Folder, err = SetObjectPath(finder, "vm", platform.DeprecatedFolder, platform.DeprecatedDatacenter)
		if err != nil {
			return err
		}
	} else if platform.DeprecatedDatacenter != "" || platform.DeprecatedFolder != "" || platform.DeprecatedCluster != "" || platform.DeprecatedDefaultDatastore != "" || platform.DeprecatedResourcePool != "" || platform.DeprecatedNetwork != "" {
		localLogger.Warn("vsphere topology fields are now deprecated; please avoid using failureDomains and the vsphere topology fields together")
	}

	return nil
}

// ConvertInstallConfig modifies a given platform spec for the new requirements.
func ConvertInstallConfig(config *types.InstallConfig) error {
	var err error
	platform := config.Platform.VSphere

	fixNoVCentersScenario(platform)
	finders := make(map[string]*find.Finder)
	for _, vcenter := range platform.VCenters {
		finder, err := GetFinder(vcenter.Server, vcenter.Username, vcenter.Password)
		if err != nil {
			return err
		}
		finders[vcenter.Server] = finder
	}
	err = fixTechPreviewZonalFailureDomainsScenario(platform, finders)
	if err != nil {
		return err
	}
	err = fixLegacyPlatformScenario(platform, finders)
	if err != nil {
		return err
	}

	return nil
}

// SetObjectPath based on the pathType will either determine the path for the type via
// a simple join of the datacenter, pathType and objectPath if finder is nil
// or via a connection to vCenter find of all child objects under the
// datacenter and pathType.
// pathType must only be "host", "vm", or "datastore".
func SetObjectPath(finder *find.Finder, pathType, objectPath, datacenter string) (string, error) {
	if objectPath != "" && !path.IsAbs(objectPath) {
		var joinedObjectPath string
		var joinedObjectFindPath string
		var paramName string

		switch pathType {
		case "host":
			paramName = "computeCluster"
		case "vm":
			paramName = "folder"
		case "datastore":
			paramName = "datastore"
		default:
			return "", errors.New("pathType can only be host, datastore or vm")
		}

		joinedObjectFindPath = path.Join("/", datacenter, pathType, "...")
		joinedObjectPath = path.Join("/", datacenter, pathType, objectPath)

		if finder == nil {
			localLogger.Warnf("%s as a non-path is now deprecated; please use the joined form: %s", paramName, joinedObjectPath)
			return joinedObjectPath, nil
		}

		newObjectPath, err := findViaPathOrName(finder, joinedObjectPath, joinedObjectFindPath)
		if err != nil {
			return "", err
		}

		if objectPath != newObjectPath {
			localLogger.Debugf("%s path changed from %s to %s", paramName, objectPath, newObjectPath)
		}
		localLogger.Warnf("%s as a non-path is now deprecated; please use the discovered form: %s", paramName, newObjectPath)

		return newObjectPath, nil
	}

	return objectPath, nil
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
