package ibmcloud

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// Subnet represents an IBM Cloud VPC Subnet
type Subnet struct {
	CIDR string
	CRN  string
	ID   string
	Name string
	VPC  string
	Zone string
}

// CreateSubnetName will build a subnet name based on the clusterID, subnet role, and subnet zone.
func CreateSubnetName(clusterID string, role string, zone string) (string, error) {
	switch role {
	case "master":
		return fmt.Sprintf("%s-subnet-control-plane-%s", clusterID, zone), nil
	case "worker":
		return fmt.Sprintf("%s-subnet-compute-%s", clusterID, zone), nil
	default:
		return "", fmt.Errorf("invalid role: %s", role)
	}
}

func getSubnets(ctx context.Context, client API, region string, subnetNames []string) (map[string]Subnet, error) {
	subnets := map[string]Subnet{}

	for _, name := range subnetNames {
		results, err := client.GetSubnetByName(ctx, name, region)
		if err != nil {
			return nil, errors.Wrapf(err, "getting subnet %s", name)
		}

		if results.ID == nil {
			return nil, errors.Errorf("%s has no ID", name)
		}

		if results.Ipv4CIDRBlock == nil {
			return nil, errors.Errorf("%s has no Ipv4CIDRBlock", *results.ID)
		}

		if results.CRN == nil {
			return nil, errors.Errorf("%s has no CRN", *results.ID)
		}

		if results.Name == nil {
			return nil, errors.Errorf("%s has no Name", *results.ID)
		}

		if results.VPC == nil {
			return nil, errors.Errorf("%s has no VPC", *results.ID)
		}

		if results.Zone == nil {
			return nil, errors.Errorf("%s has no Zone", *results.ID)
		}

		subnets[*results.ID] = Subnet{
			CIDR: *results.Ipv4CIDRBlock,
			CRN:  *results.CRN,
			ID:   *results.ID,
			Name: *results.Name,
			VPC:  *results.VPC.Name,
			Zone: *results.Zone.Name,
		}
	}

	return subnets, nil
}
