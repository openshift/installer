package gcp

import (
	"context"
	"fmt"
	"sort"
	"time"

	gcpclient "github.com/openshift/installer/pkg/client/gcp"
	"github.com/pkg/errors"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// AvailabilityZones retrieves a list of availability zones for the given project and region.
func AvailabilityZones(project, region string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	ssn, err := gcpclient.GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	svc, err := compute.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}

	regionURL := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s",
		project, region)
	req := svc.Zones.List(project).Filter(fmt.Sprintf("(region eq %s) (status eq UP)", regionURL))

	var zones []string
	if err := req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, z := range page.Items {
			zones = append(zones, z.Name)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to list zones")
	}
	if len(zones) == 0 {
		return nil, errors.New("no zone was found")
	}

	sort.Strings(zones)
	return zones, nil
}
