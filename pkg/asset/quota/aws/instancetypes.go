package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// InstanceTypeInfo describes the instance type
type InstanceTypeInfo struct {
	Name string
	vCPU int64
}

// InstanceTypes returns information of the all the instance types available for a region.
// It returns a map of instance type name to it's information.
func InstanceTypes(ctx context.Context, region string, serviceEndpoints []typesaws.ServiceEndpoint) (map[string]InstanceTypeInfo, error) {
	ret := map[string]InstanceTypeInfo{}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return ret, fmt.Errorf("failed to create AWS config: %w", err)
	}

	client := ec2.NewFromConfig(cfg,
		func(o *ec2.Options) {
			for _, endpoint := range serviceEndpoints {
				if strings.EqualFold(endpoint.Name, ec2.ServiceID) {
					o.BaseEndpoint = aws.String(endpoint.URL)
				}
			}
		},
	)

	paginator := ec2.NewDescribeInstanceTypesPaginator(client, &ec2.DescribeInstanceTypesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return ret, fmt.Errorf("failed to get ec2 instance types: %w", err)
		}
		for _, info := range page.InstanceTypes {
			ti := InstanceTypeInfo{Name: string(info.InstanceType)}
			if info.VCpuInfo == nil {
				continue
			}
			ti.vCPU = int64(aws.ToInt32(info.VCpuInfo.DefaultVCpus))
			ret[ti.Name] = ti
		}
	}

	return ret, nil
}
