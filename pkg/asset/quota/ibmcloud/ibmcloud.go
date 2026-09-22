package ibmcloud

import (
	"context"
	"fmt"

	ibmcloudic "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
)

// defaultLimits are IBM Cloud VPC default resource limits.
// Actual limits may vary by account. The check warns rather than
// blocks when it cannot determine the actual limit.
var defaultLimits = map[string]int64{
	"is/floating-ip":    20,
	"is/security-group": 25,
	"is/load-balancer":  50,
	"is/instance":       200,
	"is/vpc":            10,
}

// Load fetches current IBM Cloud VPC resource usage for quota comparison.
// The security group count is scoped to the target VPC when deploying
// into an existing VPC, since the limit is per-VPC.
func Load(ctx context.Context, client *ibmcloudic.Client, config *types.InstallConfig) ([]quota.Quota, error) {
	region := config.Platform.IBMCloud.Region

	// Resolve VPC ID for scoped security group counting.
	var vpcID string
	if config.Platform.IBMCloud.VPCName != "" {
		vpc, err := client.GetVPCByName(ctx, config.Platform.IBMCloud.VPCName)
		if err != nil {
			return nil, fmt.Errorf("looking up VPC %q: %w", config.Platform.IBMCloud.VPCName, err)
		}
		if vpc != nil && vpc.ID != nil {
			vpcID = *vpc.ID
		}
	}

	type counter struct {
		name string
		fn   func() (int, error)
	}

	counters := []counter{
		{"is/floating-ip", func() (int, error) {
			fips, err := client.ListFloatingIPs(ctx, region)
			return len(fips), err
		}},
		{"is/security-group", func() (int, error) {
			sgs, err := client.ListSecurityGroups(ctx, region, vpcID)
			return len(sgs), err
		}},
		{"is/load-balancer", func() (int, error) {
			lbs, err := client.ListLoadBalancers(ctx, region)
			return len(lbs), err
		}},
		{"is/instance", func() (int, error) {
			instances, err := client.ListInstances(ctx, region)
			return len(instances), err
		}},
		{"is/vpc", func() (int, error) {
			vpcs, err := client.GetVPCs(ctx, region)
			return len(vpcs), err
		}},
	}

	quotas := make([]quota.Quota, 0, len(counters))
	for _, c := range counters {
		count, err := c.fn()
		if err != nil {
			return nil, fmt.Errorf("counting %s: %w", c.name, err)
		}
		quotas = append(quotas, quota.Quota{
			Service: "is",
			Name:    c.name,
			Region:  region,
			InUse:   int64(count),
			Limit:   defaultLimits[c.name],
		})
	}

	return quotas, nil
}
