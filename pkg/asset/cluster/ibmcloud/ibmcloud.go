// Package ibmcloud extracts IBM Cloud metadata from install configurations.
package ibmcloud

import (
	"context"

	icibmcloud "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

// Metadata converts an install configuration to IBM Cloud metadata.
func Metadata(infraID string, config *types.InstallConfig, meta *icibmcloud.Metadata) *ibmcloud.Metadata {
	accountID, _ := meta.AccountID(context.TODO())
	cisCrn, _ := meta.CISInstanceCRN(context.TODO())
	subnets := []string{}
	controlPlaneSubnets, _ := meta.ControlPlaneSubnets(context.TODO())
	for id := range controlPlaneSubnets {
		subnets = append(subnets, id)
	}
	computeSubnets, _ := meta.ComputeSubnets(context.TODO())
	for id := range computeSubnets {
		subnets = append(subnets, id)
	}

	// TODO: For now we don't care about any duplicates in 'subnets', but might need to remove any if we need to
	// process the subnets data. Currently, if there is one or more subnet, we skip destroying all subnets (user-provided)

	return &ibmcloud.Metadata{
		AccountID:         accountID,
		BaseDomain:        config.BaseDomain,
		CISInstanceCRN:    cisCrn,
		Region:            config.Platform.IBMCloud.Region,
		ResourceGroupName: config.Platform.IBMCloud.ClusterResourceGroupName(infraID),
		Subnets:           subnets,
		VPC:               config.Platform.IBMCloud.GetVPCName(),
	}
}
