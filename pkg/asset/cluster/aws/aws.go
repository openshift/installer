// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset/installconfig"
	awsic "github.com/openshift/installer/pkg/asset/installconfig/aws"
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
	if err := tagSharedVPCResources(ctx, clusterID, installConfig); err != nil {
		return err
	}

	if err := tagSharedIAMRoles(ctx, clusterID, installConfig); err != nil {
		return err
	}

	return tagSharedIAMProfiles(ctx, clusterID, installConfig)
}

func tagSharedVPCResources(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
	if len(installConfig.Config.Platform.AWS.DeprecatedSubnets) == 0 {
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

	ids := make([]*string, 0, len(privateSubnets)+len(publicSubnets)+len(edgeSubnets))
	for id := range privateSubnets {
		ids = append(ids, aws.String(id))
	}
	for id := range publicSubnets {
		ids = append(ids, aws.String(id))
	}
	for id := range edgeSubnets {
		ids = append(ids, aws.String(id))
	}

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return errors.Wrap(err, "could not create AWS session")
	}

	tagKey, tagValue := sharedTag(clusterID)

	ec2Client := ec2.New(session, aws.NewConfig().WithRegion(installConfig.Config.Platform.AWS.Region))
	if _, err = ec2Client.CreateTagsWithContext(ctx, &ec2.CreateTagsInput{
		Resources: ids,
		Tags:      []*ec2.Tag{{Key: &tagKey, Value: &tagValue}},
	}); err != nil {
		return errors.Wrap(err, "could not add tags to subnets")
	}

	if zone := installConfig.Config.AWS.HostedZone; zone != "" {
		r53cfg := awsic.GetR53ClientCfg(session, installConfig.Config.AWS.HostedZoneRole)
		route53Client := route53.New(session, r53cfg)
		if _, err := route53Client.ChangeTagsForResourceWithContext(ctx, &route53.ChangeTagsForResourceInput{
			ResourceType: aws.String("hostedzone"),
			ResourceId:   aws.String(zone),
			AddTags:      []*route53.Tag{{Key: &tagKey, Value: &tagValue}},
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

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return fmt.Errorf("could not create AWS session: %w", err)
	}

	tagKey, tagValue := sharedTag(clusterID)

	iamClient := iam.New(session, aws.NewConfig().WithRegion(installConfig.Config.Platform.AWS.Region))
	for role := range iamRoles {
		if _, err := iamClient.TagRoleWithContext(ctx, &iam.TagRoleInput{
			RoleName: aws.String(role),
			Tags: []*iam.Tag{
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

	session, err := installConfig.AWS.Session(ctx)
	if err != nil {
		return errors.Wrap(err, "could not create AWS session")
	}
	iamClient := iam.New(session, aws.NewConfig().WithRegion(installConfig.Config.AWS.Region))

	tagKey, tagValue := sharedTag(clusterID)
	for name := range iamProfileNames {
		if _, err := iamClient.TagInstanceProfileWithContext(ctx, &iam.TagInstanceProfileInput{
			InstanceProfileName: aws.String(name),
			Tags: []*iam.Tag{
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
