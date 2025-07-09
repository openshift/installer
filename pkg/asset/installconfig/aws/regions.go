package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

const (
	isoPartition  = "aws-iso"
	isobPartition = "aws-iso-b"
)

// knownPublicRegions is the subset of public AWS regions where RHEL CoreOS images are published.
// This subset does not include supported regions which are found in other partitions, such as us-gov-east-1.
// Returns: a list of region names.
func knownPublicRegions(architecture types.Architecture) ([]string, error) {
	required := rhcos.AMIRegions(architecture)

	regions, err := GetRegions(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get aws regions: %w", err)
	}

	foundRegions := []string{}
	for _, region := range regions {
		if required.Has(region) {
			foundRegions = append(foundRegions, region)
		}
	}
	return foundRegions, nil
}

// IsKnownPublicRegion returns true if a specified region is Known to the installer.
// A known region is the subset of public AWS regions where RHEL CoreOS images are published.
func IsKnownPublicRegion(region string, architecture types.Architecture) (bool, error) {
	publicRegions, err := knownPublicRegions(architecture)
	if err != nil {
		return false, err
	}
	return sets.New(publicRegions...).Has(region), nil
}

// IsSecretRegion determines if the region is part of a secret partition.
func IsSecretRegion(region string) (bool, error) {
	endpoint, err := ec2.NewDefaultEndpointResolver().ResolveEndpoint(region, ec2.EndpointResolverOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to resolve AWS ec2 endpoint: %w", err)
	}

	switch endpoint.PartitionID {
	case isoPartition, isobPartition:
		return true, nil
	}

	return false, nil
}
