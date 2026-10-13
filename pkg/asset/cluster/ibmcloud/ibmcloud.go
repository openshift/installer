// Package ibmcloud extracts IBM Cloud metadata from install configurations.
package ibmcloud

import (
	"context"

	icibmcloud "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

// Metadata converts an install configuration to IBM Cloud metadata.
func Metadata(infraID string, config *types.InstallConfig) *ibmcloud.Metadata {
	meta := icibmcloud.NewMetadata(config)
	accountID, _ := meta.AccountID(context.TODO())
	cisCrn, _ := meta.CISInstanceCRN(context.TODO())
	dnsInstance, _ := meta.DNSInstance(context.TODO())

	var dnsInstanceID string
	if dnsInstance != nil {
		dnsInstanceID = dnsInstance.ID
	}

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
		DNSInstanceID:     dnsInstanceID,
		Region:            config.Platform.IBMCloud.Region,
		ResourceGroupName: config.Platform.IBMCloud.ClusterResourceGroupName(infraID),
		ServiceEndpoints:  config.Platform.IBMCloud.ServiceEndpoints,
		Subnets:           subnets,
		VPC:               config.Platform.IBMCloud.GetVPCName(),
	}
}
