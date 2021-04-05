// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"

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

	if err := tagSubnetEC2Instances(ctx, clusterID, installConfig); err != nil {
		return err
	}

	if err := tagIamRoles(ctx, clusterID, installConfig); err != nil {
		return err
	}

	return nil
}

func tagSubnetEC2Instances(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	if len(installConfig.Config.Platform.AWS.Subnets) == 0 {
		return nil
	}

	privateSubnets, err := installConfig.AWS.PrivateSubnets(ctx)
	if err != nil {
		return err
	}

	publicSubnets, err := installConfig.AWS.PublicSubnets(ctx)
	if err != nil {
		return err
	}

	//arns := make([]*string, 0, len(privateSubnets)+len(publicSubnets))
	ids := make([]*string, 0, len(privateSubnets)+len(publicSubnets))
	for id := range privateSubnets {
		ids = append(ids, aws.String(id))
	}
	for id := range publicSubnets {
		ids = append(ids, aws.String(id))
	}

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return err
	}
	key, value := sharedTag(clusterID)
	ec2Tags := []*ec2.Tag{{Key: &key, Value: &value}}
	ec2Client := ec2.New(session, aws.NewConfig().WithRegion(installConfig.Config.Platform.AWS.Region))

	if _, err = ec2Client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
		Resources: ids,
		Tags:      ec2Tags,
	}); err != nil {
		return err
	}

	return nil
}

func tagIamRoles(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	workerMachinePool := installConfig.Config.WorkerMachinePool()

	var iamRoleNames []*string
	if installConfig.Config.ControlPlane.Platform.AWS != nil {
		if installConfig.Config.ControlPlane.Platform.AWS.IAMRole != "" {
			iamRoleNames = append(iamRoleNames, &installConfig.Config.ControlPlane.Platform.AWS.IAMRole)
		}
	}

	if workerMachinePool.Platform.AWS.IAMRole != "" {
		iamRoleNames = append(iamRoleNames, &workerMachinePool.Platform.AWS.IAMRole)
	}

	if len(iamRoleNames) == 0 {
		return nil
	}

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return err
	}
	key, value := sharedTag(clusterID)
	iamTags := []*iam.Tag{{Key: &key, Value: &value}}
	iamClient := iam.New(session, aws.NewConfig().WithRegion(installConfig.Config.Platform.AWS.Region))

	for _, iamRoleName := range iamRoleNames {
		if _, err := iamClient.TagRoleWithContext(ctx, &iam.TagRoleInput{
			RoleName: iamRoleName,
			Tags:     iamTags,
		}); err != nil {
			return errors.Wrapf(err, "could not tag %s role", *iamRoleName)
		}
	}

	return nil
}

func sharedTag(clusterID string) (string, string) {
	return fmt.Sprintf("kubernetes.io/cluster/%s", clusterID), "shared"
}
