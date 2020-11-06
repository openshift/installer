// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

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

	arns := make([]*string, 0, len(privateSubnets)+len(publicSubnets))
	for _, subnet := range privateSubnets {
		arns = append(arns, aws.String(subnet.ARN))
	}
	for _, subnet := range publicSubnets {
		arns = append(arns, aws.String(subnet.ARN))
	}

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return err
	}

	tags := []*ec2.Tag{
		{
			Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", clusterID)),
			Value: aws.String("shared"),
		},
	}
	client := ec2.New(session, aws.NewConfig().WithRegion(installConfig.Config.Platform.AWS.Region))

	if _, err = client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
		Resources: arns,
		Tags:      tags,
	}); err != nil {
		return err
	}

	return nil
}
