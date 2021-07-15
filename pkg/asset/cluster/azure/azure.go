// Package azure extracts AZURE metadata from install configurations.
package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// Metadata converts an install configuration to Azure metadata.
func Metadata(config *types.InstallConfig) *azure.Metadata {
	return &azure.Metadata{
		ARMEndpoint:                 config.Platform.Azure.ARMEndpoint,
		CloudName:                   config.Platform.Azure.CloudName,
		Region:                      config.Platform.Azure.Region,
		ResourceGroupName:           config.Azure.ResourceGroupName,
		BaseDomainResourceGroupName: config.Azure.BaseDomainResourceGroupName,
		ClusterName:                 config.ClusterName,
	}
}

// PreTerraform performs any infrastructure initialization which must
// happen before Terraform creates the remaining infrastructure.
func PreTerraform(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	if len(installConfig.Config.Azure.ResourceGroupName) == 0 {
		return nil
	}

	session, err := installConfig.Azure.Session()
	if err != nil {
		return errors.Wrap(err, "failed to get session")
	}

	client := resources.NewGroupsClientWithBaseURI(session.Environment.ResourceManagerEndpoint, session.Credentials.SubscriptionID)
	client.Authorizer = session.Authorizer
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	group, err := client.Get(ctx, installConfig.Config.Azure.ResourceGroupName)
	if err != nil {
		return errors.Wrap(err, "failed to get the resource group")
	}

	if group.Tags == nil {
		group.Tags = map[string]*string{}
	}
	group.Tags[fmt.Sprintf("kubernetes.io_cluster.%s", clusterID)] = to.StringPtr("owned")
	logrus.Debugf("Tagging %s with kubernetes.io/cluster/%s: shared", installConfig.Config.Azure.ResourceGroupName, clusterID)
	_, err = client.Update(ctx, installConfig.Config.Azure.ResourceGroupName, resources.GroupPatchable{
		Tags: group.Tags,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to tag the resource group %s", installConfig.Config.Azure.ResourceGroupName)
	}
	return nil
}
