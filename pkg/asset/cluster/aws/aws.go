// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// Metadata converts an install configuration to AWS metadata.
func Metadata(clusterID, infraID string, config *types.InstallConfig) *awstypes.Metadata {
	return &awstypes.Metadata{
		Region: config.Platform.AWS.Region,
		Identifier: []map[string]string{{
			fmt.Sprintf("kubernetes.io/cluster/%s", infraID): "owned",
		}, {
			"openshiftClusterID": clusterID,
		}},
		ServiceEndpoints: config.AWS.ServiceEndpoints,
	}
}

// PreTerraform performs any infrastructure initialization which must
// happen before Terraform creates the remaining infrastructure.
func PreTerraform(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	if len(installConfig.Config.Platform.AWS.Subnets) == 0 {
		return nil
	}

	privateSubnets, err := installConfig.AWS.PrivateSubnets(ctx)
	if err != nil {
		return err
	}

	publicSubnets, err := installConfig.AWS.PublicSubnets(ctx)

	arns := make([]string, 0, len(privateSubnets)+len(publicSubnets))
	for _, subnet := range privateSubnets {
		arns = append(arns, subnet.ARN)
	}
	for _, subnet := range publicSubnets {
		arns = append(arns, subnet.ARN)
	}

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return err
	}

	request := &resourcegroupstaggingapi.TagResourcesInput{
		Tags: map[string]*string{
			fmt.Sprintf("kubernetes.io/cluster/%s", clusterID): aws.String("shared"),
		},
	}

	tagClient := resourcegroupstaggingapi.New(session, aws.NewConfig().WithRegion(installConfig.Config.Platform.AWS.Region))
	for i := 0; i < len(arns); i += 20 {
		request.ResourceARNList = make([]*string, 0, 20)
		for j := 0; i+j < len(arns) && j < 20; j++ {
			logrus.Debugf("Tagging %s with kubernetes.io/cluster/%s: shared", arns[i+j], clusterID)
			request.ResourceARNList = append(request.ResourceARNList, aws.String(arns[i+j]))
		}
		_, err = tagClient.TagResourcesWithContext(ctx, request)
		if err != nil {
			return err
		}
	}

	return nil
}
