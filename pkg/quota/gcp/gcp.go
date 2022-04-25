package gcp

import (
	"context"
	"net/http"
	"sort"
	"strings"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	serviceusage "google.golang.org/api/serviceusage/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpclient "github.com/openshift/installer/pkg/client/gcp"
	"github.com/openshift/installer/pkg/quota"
)

// Load load the quota information for a project and provided services. It provides information
// about the usage and limit for each resource quota.
// roles/servicemanagement.quotaViewer role allows users to fetch the required details.
func Load(ctx context.Context, project string, services ...string) ([]quota.Quota, error) {
	ssn, err := gcpclient.GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	options := []option.ClientOption{
		option.WithCredentials(ssn.Credentials),
	}
	servicesSvc, err := serviceusage.NewService(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create services svc")
	}
	metricsSvc, err := monitoring.NewMetricClient(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create metrics svc")
	}

	limits, err := loadLimits(ctx, servicesSvc.Services, project, services...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load quota limits")
	}
	usages, err := loadUsage(ctx, metricsSvc, project)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load quota usages")
	}
	return newQuotas(usages, limits), nil
}

// record stores the data from quota limits and usages.
type record struct {
	Service string
	Name    string

	// this can either be a region or zones or const "global"
	Location string

	Value int64
}

// newQuotas combines the usage and quota limit to create a list of Quotas.
// newQuotas matches the limits to the corrsponding usage to create summary of the quota.
// Usages are usually reported for location as global or either per region or per zone, while the limits
// are usually set for location as global or per region and therfore the quota consolidate the zone's usages to
// it's region as sum.
// When there is no matching usage found for a limit, the usage is treated as zero.
func newQuotas(usages []record, limits []record) []quota.Quota {
	sort.Slice(usages, func(i, j int) bool {
		return usages[i].Service < usages[j].Service && usages[i].Name < usages[j].Name && usages[i].Location < usages[j].Location
	})
	sort.Slice(limits, func(i, j int) bool {
		return limits[i].Service < limits[j].Service && limits[i].Name < limits[j].Name && limits[i].Location < limits[j].Location
	})

	findUsage := func(l record) (record, bool) {
		for _, u := range usages {
			if !strings.EqualFold(l.Service, u.Service) {
				continue
			}
			if !strings.EqualFold(l.Name, u.Name) {
				continue
			}
			if !strings.EqualFold(l.Location, u.Location) {
				continue
			}
			return u, true
		}
		return record{}, false
	}

	var quotas []quota.Quota
	for _, l := range limits {
		q := quota.Quota{
			Service: l.Service,
			Name:    l.Name,
			Region:  l.Location,

			Limit: l.Value,
		}
		u, ok := findUsage(l)
		if !ok {
			q.InUse = int64(0)
		} else {
			q.InUse = u.Value
		}
		quotas = append(quotas, q)
	}
	return quotas
}

// IsUnauthorized checks if the error is un authorized.
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	var gErr *googleapi.Error
	if errors.As(err, &gErr) {
		return gErr.Code == http.StatusUnauthorized || gErr.Code == http.StatusForbidden
	}

	if grpcCode := status.Code(errors.Cause(err)); grpcCode != codes.OK {
		return grpcCode == codes.PermissionDenied
	}
	return false
}
