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

	return &ibmcloud.Metadata{
		AccountID:         accountID,
		BaseDomain:        config.BaseDomain,
		CISInstanceCRN:    cisCrn,
		Region:            config.Platform.IBMCloud.Region,
		ResourceGroupName: config.Platform.IBMCloud.ClusterResourceGroupName(infraID),
		Subnets:           config.Platform.IBMCloud.Subnets,
		VPC:               config.Platform.IBMCloud.VPC,
	}
}
