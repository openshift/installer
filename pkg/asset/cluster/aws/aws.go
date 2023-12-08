// Package aws extracts AWS metadata from install configurations.
package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	awsic "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

type AWSProvisioner struct {
	ic *installconfig.InstallConfig
}

func InitAWSProvisioner(ic *installconfig.InstallConfig) AWSProvisioner {
	return AWSProvisioner{ic: ic}
}

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

	return nil
}

func tagSharedVPCResources(ctx context.Context, clusterID string, installConfig *installconfig.InstallConfig) error {
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

	ids := make([]*string, 0, len(privateSubnets)+len(publicSubnets))
	for id := range privateSubnets {
		ids = append(ids, aws.String(id))
	}
	for id := range publicSubnets {
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

func sharedTag(clusterID string) (string, string) {
	return fmt.Sprintf("kubernetes.io/cluster/%s", clusterID), "shared"
}

// PreProvision creates the control plane and compute instance profiles.
func (a AWSProvisioner) PreProvision(clusterID string) error {
	if err := putIAMRoles(clusterID, a.ic); err != nil {
		return fmt.Errorf("error putting IAM roles: %w", err)
	}
	return nil
}

// ValidControlPlaneEndpoint creates the DNS records for the cluster.
func (a AWSProvisioner) ValidControlPlaneEndpoint(cluster *clusterv1.Cluster) error {
	if err := createDNSRecords(a.ic, cluster); err != nil {
		return fmt.Errorf("error creating DNS records: %w", err)
	}
	return nil
}

func putIAMRoles(clusterID string, ic *installconfig.InstallConfig) error {
	// Create the IAM Role with the aws sdk.
	// https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.CreateRole
	session, err := ic.AWS.Session(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to load AWS session")
	}
	svc := iam.New(session)

	// Create the IAM Roles for master and workers.
	clusterOwnedIAMTag := &iam.Tag{
		Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", clusterID)),
		Value: aws.String("owned"),
	}

	assumePolicyBytes, err := json.Marshal(assumePolicy)
	if err != nil {
		return errors.Wrap(err, "failed to marshal assume policy")
	}

	for _, role := range []string{"master", "worker"} {
		roleName := aws.String(fmt.Sprintf("%s-%s-role", clusterID, role))
		logrus.Infof("Checking if %s IAM role exists", *roleName)
		if _, err := svc.GetRole(&iam.GetRoleInput{RoleName: roleName}); err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != iam.ErrCodeNoSuchEntityException {
				return errors.Wrapf(err, "failed to get %s role", role)
			}
			logrus.Infof("Creating %s IAM role because it was not found", *roleName)
			createRoleInput := &iam.CreateRoleInput{
				RoleName:                 roleName,
				AssumeRolePolicyDocument: aws.String(string(assumePolicyBytes)),
				Tags:                     []*iam.Tag{clusterOwnedIAMTag},
			}
			if _, err := svc.CreateRole(createRoleInput); err != nil {
				return errors.Wrapf(err, "failed to create %s role", role)
			}
			time.Sleep(10 * time.Second)
			if err := svc.WaitUntilRoleExists(&iam.GetRoleInput{RoleName: roleName}); err != nil {
				return errors.Wrapf(err, "failed to wait for %s role to exist", role)
			}
		}

		// Put the policy inline.
		policyName := aws.String(fmt.Sprintf("%s-%s-policy", clusterID, role))
		b, err := json.Marshal(policies[role])
		if err != nil {
			return errors.Wrapf(err, "failed to marshal %s policy", role)
		}
		if _, err := svc.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyDocument: aws.String(string(b)),
			PolicyName:     policyName,
			RoleName:       roleName,
		}); err != nil {
			return errors.Wrapf(err, "failed to create inline policy for role %s ", role)
		}

		profileName := aws.String(fmt.Sprintf("%s-%s-profile", clusterID, role))
		if _, err := svc.GetInstanceProfile(&iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != iam.ErrCodeNoSuchEntityException {
				return errors.Wrapf(err, "failed to get %s instance profile", role)
			}
			// If the profile does not exist, create it.
			if _, err := svc.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
				InstanceProfileName: profileName,
				Tags:                []*iam.Tag{clusterOwnedIAMTag},
			}); err != nil {
				return errors.Wrapf(err, "failed to create %s instance profile", role)
			}
			time.Sleep(10 * time.Second)
			if err := svc.WaitUntilInstanceProfileExists(&iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
				return errors.Wrapf(err, "failed to wait for %s role to exist", role)
			}

			// Finally, attach the role to the profile.
			if _, err := svc.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
				InstanceProfileName: profileName,
				RoleName:            roleName,
			}); err != nil {
				return errors.Wrapf(err, "failed to add %s role to instance profile", role)
			}
		}
	}

	return nil
}

func createDNSRecords(ic *installconfig.InstallConfig, cluster *clusterv1.Cluster) error {
	ssn, err := ic.AWS.Session(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	client := awsic.NewClient(ssn)
	r53cfg := awsic.GetR53ClientCfg(ssn, "")
	err = client.CreateOrUpdateRecord(ic.Config, cluster.Spec.ControlPlaneEndpoint.Host, r53cfg)
	if err != nil {
		return fmt.Errorf("failed to create route53 records: %w", err)
	}
	logrus.Infof("Created Route53 records to control plane load balancer.")
	return nil
}
