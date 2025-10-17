// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// Metadata converts an install configuration to AWS metadata.
func Metadata(clusterID, infraID string, config *types.InstallConfig) *awstypes.Metadata {
	return &awstypes.Metadata{
		Region: config.Platform.AWS.Region,
		Identifier: []map[string]string{
			{fmt.Sprintf("kubernetes.io/cluster/%s", infraID): "owned"},
			{"openshiftClusterID": clusterID},
			{fmt.Sprintf("sigs.k8s.io/cluster-api-provider-aws/cluster/%s", infraID): "owned"},
		},
		ServiceEndpoints: config.AWS.ServiceEndpoints,
		ClusterDomain:    config.ClusterDomain(),
		HostedZoneRole:   config.AWS.HostedZoneRole,
	}
}

// PreTerraform performs any infrastructure initialization which must
// happen before Terraform creates the remaining infrastructure.
func PreTerraform(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	//if err := tagSharedVPCResources(ctx, clusterID, installConfig); err != nil {
	//	return err
	//}
	//
	//if err := tagSharedIAMRoles(ctx, clusterID, installConfig); err != nil {
	//	return err
	//}
	//
	//return tagSharedIAMProfiles(ctx, clusterID, installConfig)

	return nil
}

func tagSharedVPCResources(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	if len(installConfig.Config.Platform.AWS.VPC.Subnets) == 0 {
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

	edgeSubnets, err := installConfig.AWS.EdgeSubnets(ctx)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(privateSubnets)+len(publicSubnets)+len(edgeSubnets))
	for id := range privateSubnets {
		ids = append(ids, id)
	}
	for id := range publicSubnets {
		ids = append(ids, id)
	}
	for id := range edgeSubnets {
		ids = append(ids, id)
	}

	tagKey, tagValue := sharedTag(clusterID)

	ec2Client, err := awsconfig.NewEC2Client(ctx, awsconfig.EndpointOptions{
		Region:    installConfig.Config.Platform.AWS.Region,
		Endpoints: installConfig.Config.Platform.AWS.ServiceEndpoints,
	})
	if err != nil {
		return fmt.Errorf("failed to create EC2 client: %w", err)
	}

	if _, err = ec2Client.CreateTags(ctx, &ec2.CreateTagsInput{
		Resources: ids,
		Tags:      []ec2types.Tag{{Key: &tagKey, Value: &tagValue}},
	}); err != nil {
		return errors.Wrap(err, "could not add tags to subnets")
	}

	if zone := installConfig.Config.AWS.HostedZone; zone != "" {
		roleArn := installConfig.Config.AWS.HostedZoneRole
		route53Client, err := awsconfig.NewRoute53Client(ctx, awsconfig.EndpointOptions{
			Region:    installConfig.Config.Platform.AWS.Region,
			Endpoints: installConfig.Config.Platform.AWS.ServiceEndpoints,
		}, roleArn)
		if err != nil {
			return fmt.Errorf("failed to create Route 53 client: %w", err)
		}

		if _, err := route53Client.ChangeTagsForResource(ctx, &route53.ChangeTagsForResourceInput{
			ResourceType: r53types.TagResourceTypeHostedzone,
			ResourceId:   aws.String(zone),
			AddTags:      []r53types.Tag{{Key: &tagKey, Value: &tagValue}},
		}); err != nil {
			return errors.Wrap(err, "could not add tags to hosted zone")
		}
	}

	return nil
}

func tagSharedIAMRoles(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	iamRoles := sets.New[string]()
	{
		mpool := awstypes.MachinePool{}
		mpool.Set(installConfig.Config.AWS.DefaultMachinePlatform)
		if mp := installConfig.Config.ControlPlane; mp != nil {
			mpool.Set(mp.Platform.AWS)
		}
		if len(mpool.IAMRole) > 0 {
			iamRoles.Insert(mpool.IAMRole)
		}
	}

	for _, compute := range installConfig.Config.Compute {
		mpool := awstypes.MachinePool{}
		mpool.Set(installConfig.Config.AWS.DefaultMachinePlatform)
		mpool.Set(compute.Platform.AWS)
		if len(mpool.IAMRole) > 0 {
			iamRoles.Insert(mpool.IAMRole)
		}
	}

	// If compute stanza was not defined, it will inherit from DefaultMachinePlatform later on.
	if installConfig.Config.Compute == nil {
		mpool := installConfig.Config.AWS.DefaultMachinePlatform
		if mpool != nil && len(mpool.IAMRole) > 0 {
			iamRoles.Insert(mpool.IAMRole)
		}
	}

	if iamRoles.Len() == 0 {
		return nil
	}

	logrus.Debugf("Tagging shared instance roles: %v", sets.List(iamRoles))

	tagKey, tagValue := sharedTag(clusterID)

	iamClient, err := awsconfig.NewIAMClient(ctx, awsconfig.EndpointOptions{
		Region:    installConfig.Config.Platform.AWS.Region,
		Endpoints: installConfig.Config.Platform.AWS.ServiceEndpoints,
	})
	if err != nil {
		return fmt.Errorf("failed to create IAM client: %w", err)
	}

	for role := range iamRoles {
		if _, err := iamClient.TagRole(ctx, &iam.TagRoleInput{
			RoleName: aws.String(role),
			Tags: []iamtypes.Tag{
				{Key: aws.String(tagKey), Value: aws.String(tagValue)},
			},
		}); err != nil {
			return fmt.Errorf("could not tag %q instance role: %w", role, err)
		}
	}

	return nil
}

// tagSharedIAMProfiles tags users BYO instance profiles so they are not destroyed by the Installer.
func tagSharedIAMProfiles(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	iamProfileNames := sets.New[string]()

	{
		mpool := awstypes.MachinePool{}
		mpool.Set(installConfig.Config.AWS.DefaultMachinePlatform)

		if mp := installConfig.Config.ControlPlane; mp != nil {
			mpool.Set(mp.Platform.AWS)
		}
		if len(mpool.IAMProfile) > 0 {
			iamProfileNames.Insert(mpool.IAMProfile)
		}
	}

	for _, compute := range installConfig.Config.Compute {
		mpool := awstypes.MachinePool{}
		mpool.Set(installConfig.Config.AWS.DefaultMachinePlatform)
		mpool.Set(compute.Platform.AWS)
		if len(mpool.IAMProfile) > 0 {
			iamProfileNames.Insert(mpool.IAMProfile)
		}
	}

	// If compute stanza was not defined in the install-config.yaml, it will inherit from the
	// DefaultMachinePlatform later on.
	if installConfig.Config.Compute == nil {
		mpool := installConfig.Config.AWS.DefaultMachinePlatform
		if mpool != nil && len(mpool.IAMProfile) > 0 {
			iamProfileNames.Insert(mpool.IAMProfile)
		}
	}

	if iamProfileNames.Len() == 0 {
		return nil
	}

	logrus.Debugf("Tagging shared instance profiles: %v", sets.List(iamProfileNames))

	iamClient, err := awsconfig.NewIAMClient(ctx, awsconfig.EndpointOptions{
		Region:    installConfig.Config.Platform.AWS.Region,
		Endpoints: installConfig.Config.Platform.AWS.ServiceEndpoints,
	})
	if err != nil {
		return fmt.Errorf("failed to create IAM client: %w", err)
	}

	tagKey, tagValue := sharedTag(clusterID)
	for name := range iamProfileNames {
		if _, err := iamClient.TagInstanceProfile(ctx, &iam.TagInstanceProfileInput{
			InstanceProfileName: aws.String(name),
			Tags: []iamtypes.Tag{
				{Key: aws.String(tagKey), Value: aws.String(tagValue)},
			},
		}); err != nil {
			return fmt.Errorf("could not tag %q instance profile: %w", name, err)
		}
	}

	return nil
}

func sharedTag(clusterID string) (string, string) {
	return fmt.Sprintf("kubernetes.io/cluster/%s", clusterID), "shared"
}
