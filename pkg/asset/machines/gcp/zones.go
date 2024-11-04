package gcp

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/util/sets"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
)

// AvailabilityZones retrieves a list of availability zones for the given project and region.
func AvailabilityZones(project, region string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	ssn, err := gcpconfig.GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	svc, err := compute.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}

	zones, err := gcpconfig.GetZones(ctx, svc, project, region)
	if err != nil {
		return nil, errors.New("no zone was found")
	}

	zoneNames := make([]string, 0, len(zones))
	for _, z := range zones {
		zoneNames = append(zoneNames, z.Name)
	}

	sort.Strings(zoneNames)
	return zoneNames, nil
}

// ZonesForInstanceType retrieves a filtered list of availability zones where
// the particular instance type is available. This is mainly necessary for
// arm64, since the instance t2a-standard-* is not available in all
// availability zones.
func ZonesForInstanceType(project, region, instanceType string) ([]string, error) {
	ssn, err := gcpconfig.GetSession(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	svc, err := compute.NewService(context.Background(), option.WithCredentials(ssn.Credentials))
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
