package gcp

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/util/sets"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// ZonesForInstanceType retrieves a filtered list of availability zones where
// the particular instance type is available. This is mainly necessary for
// arm64, since the instance t2a-standard-* is not available in all
// availability zones.
func ZonesForInstanceType(project, region, instanceType string, endpoint *gcptypes.PSCEndpoint) ([]string, error) {
	opts := []option.ClientOption{}
	if endpoint != nil && !endpoint.ClusterUseOnly {
		opts = append(opts, gcpconfig.CreateEndpointOption(endpoint.Name, gcpconfig.ServiceNameGCPCompute))
	}
	svc, err := gcpconfig.GetComputeService(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pZones, err := gcpconfig.GetZones(ctx, svc, project, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get zones for project: %w", err)
	}
	pZoneNames := sets.New[string]()
	for _, z := range pZones {
		pZoneNames.Insert(z.Name)
	}

	machines, err := gcpconfig.GetMachineTypeList(ctx, svc, project, region, instanceType, "items/*/machineTypes(zone),nextPageToken")
	if err != nil {
		return nil, fmt.Errorf("failed to get zones for instance type: %w", err)
	}
	// Custom machine types do not show up in the list. Let's fallback to the project zones
	if len(machines) == 0 {
		return sets.List(pZoneNames), nil
	}

	zones := sets.New[string]()
	for _, machine := range machines {
		zones.Insert(machine.Zone)
	}

	// Not all instance zones might be available in the project
	return sets.List(zones.Intersection(pZoneNames)), nil
}
