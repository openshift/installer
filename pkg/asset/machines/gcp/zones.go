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

	regionURL := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s",
		project, region)
	filter := fmt.Sprintf("(region eq %s) (status eq UP)", regionURL)
	zones, err := gcpconfig.GetZones(ctx, svc, project, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list zones")
	}
	if len(zones) == 0 {
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ssn, err := gcpconfig.GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	svc, err := compute.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}

	machines, err := gcpconfig.GetMachineTypeList(ctx, svc, project, region, instanceType, "items/*/machineTypes(zone),nextPageToken")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get zones for instance type")
	}

	found := sets.New[string]()
	for _, machine := range machines {
		found.Insert(machine.Zone)
	}

	return sets.List(found), nil
}
