package aws

import (
	"context"
	"fmt"

	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/quota"
	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// Load load the quota information for a region. It provides information
// about the usage and limit for each resource quota.
func Load(ctx context.Context, region string, endpoints []typesaws.ServiceEndpoint, services ...string) ([]quota.Quota, error) {
	client, err := awsconfig.NewServiceQuotasClient(ctx, awsconfig.EndpointOptions{
		Region:    region,
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create service quotas client: %w", err)
	}

	records, err := loadLimits(ctx, client, services...)
	if err != nil {
		return nil, fmt.Errorf("failed to load limits for servicequotas: %w", err)
	}
	return newQuota(region, records), nil
}

func newQuota(region string, limits []record) []quota.Quota {
	var ret []quota.Quota
	for _, limit := range limits {
		q := quota.Quota{
			Service: limit.Service,
			Name:    fmt.Sprintf("%s/%s", limit.Service, limit.Name),
			Region:  region,
			InUse:   0,
			Limit:   limit.Value,
		}
		if limit.global {
			q.Region = "global"
		}
		ret = append(ret, q)
	}
	return ret
}
