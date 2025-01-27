package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

// knownPublicRegions is the subset of public AWS regions where RHEL CoreOS images are published.
// This subset does not include supported regions which are found in other partitions, such as us-gov-east-1.
// Returns: a list of region names.
func knownPublicRegions(architecture types.Architecture) []string {
	required := rhcos.AMIRegions(architecture)

	regions := []string{}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithDefaultRegion("us-east-1"))
	if err != nil {
		logrus.Debugf("failed to create aws config: %v", err)
		return regions
	}
	ec2Config := ec2.NewFromConfig(cfg)

	output, err := ec2Config.DescribeRegions(context.Background(), &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false),
	})
	if err != nil {
		logrus.Debugf("failed to describe regions: %v", err)
		return regions
	}
	for _, region := range output.Regions {
		if required.Has(*region.RegionName) {
			regions = append(regions, *region.RegionName)
		}
	}
	return regions
}

// IsKnownPublicRegion returns true if a specified region is Known to the installer.
// A known region is the subset of public AWS regions where RHEL CoreOS images are published.
func IsKnownPublicRegion(region string, architecture types.Architecture) bool {
	publicRegions := sets.New(knownPublicRegions(architecture)...)
	return publicRegions.Has(region)
}
