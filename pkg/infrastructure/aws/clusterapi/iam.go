package clusterapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
)

const (
	master = "master"
	worker = "worker"
)

var (
	policies = map[string]*iamv1.PolicyDocument{
		master: {
			Version: "2012-10-17",
			Statement: []iamv1.StatementEntry{
				{
					Effect: "Allow",
					Action: []string{
						"ec2:AttachVolume",
						"ec2:AuthorizeSecurityGroupIngress",
						"ec2:CreateSecurityGroup",
						"ec2:CreateTags",
						"ec2:CreateVolume",
						"ec2:DeleteSecurityGroup",
						"ec2:DeleteVolume",
						"ec2:Describe*",
						"ec2:DetachVolume",
						"ec2:ModifyInstanceAttribute",
						"ec2:ModifyVolume",
						"ec2:RevokeSecurityGroupIngress",
						"elasticloadbalancing:AddTags",
						"elasticloadbalancing:AttachLoadBalancerToSubnets",
						"elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
						"elasticloadbalancing:CreateListener",
						"elasticloadbalancing:CreateLoadBalancer",
						"elasticloadbalancing:CreateLoadBalancerPolicy",
						"elasticloadbalancing:CreateLoadBalancerListeners",
						"elasticloadbalancing:CreateTargetGroup",
						"elasticloadbalancing:ConfigureHealthCheck",
						"elasticloadbalancing:DeleteListener",
						"elasticloadbalancing:DeleteLoadBalancer",
						"elasticloadbalancing:DeleteLoadBalancerListeners",
						"elasticloadbalancing:DeleteTargetGroup",
						"elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
						"elasticloadbalancing:DeregisterTargets",
						"elasticloadbalancing:Describe*",
						"elasticloadbalancing:DetachLoadBalancerFromSubnets",
						"elasticloadbalancing:ModifyListener",
						"elasticloadbalancing:ModifyLoadBalancerAttributes",
						"elasticloadbalancing:ModifyTargetGroup",
						"elasticloadbalancing:ModifyTargetGroupAttributes",
						"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
						"elasticloadbalancing:RegisterTargets",
						"elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
						"elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
						"kms:DescribeKey",
					},
					Resource: iamv1.Resources{
						"*",
					},
				},
			},
		},
		worker: {
			Version: "2012-10-17",
			Statement: []iamv1.StatementEntry{
				{
					Effect: "Allow",
					Action: iamv1.Actions{
						"ec2:DescribeInstances",
						"ec2:DescribeRegions",
					},
					Resource: iamv1.Resources{"*"},
				},
			},
		},
	}
)

// createIAMRoles creates the roles used by control-plane and compute nodes.
func createIAMRoles(ctx context.Context, infraID string, ic *installconfig.InstallConfig) error {
	logrus.Infoln("Reconciling IAM roles for control-plane and compute nodes")
	// Create the IAM Role with the aws sdk.
	// https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.CreateRole
	session, err := ic.AWS.Session(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS session: %w", err)
	}
	svc := iam.New(session)

	// Create the IAM Roles for master and workers.
	tags := []*iam.Tag{
		{
			Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", infraID)),
			Value: aws.String("owned"),
		},
	}

	for k, v := range ic.Config.AWS.UserTags {
		tags = append(tags, &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	assumePolicy := &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: iamv1.Statements{
			{
				Effect: "Allow",
				Principal: iamv1.Principals{
					iamv1.PrincipalService: []string{
						getPartitionService(ic.AWS.Region),
					},
				},
				Action: iamv1.Actions{
					"sts:AssumeRole",
				},
			},
		},
	}
	assumePolicyBytes, err := json.Marshal(assumePolicy)
	if err != nil {
		return fmt.Errorf("failed to marshal assume policy: %w", err)
	}

	for _, role := range []string{master, worker} {
		roleName, err := getOrCreateIAMRole(ctx, role, infraID, string(assumePolicyBytes), *ic, tags, svc)
		if err != nil {
			return fmt.Errorf("failed to create IAM %s role: %w", role, err)
		}

		profileName := aws.String(fmt.Sprintf("%s-%s-profile", infraID, role))
		if _, err := svc.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != iam.ErrCodeNoSuchEntityException {
				return fmt.Errorf("failed to get %s instance profile: %w", role, err)
			}
			// If the profile does not exist, create it.
			if _, err := svc.CreateInstanceProfileWithContext(ctx, &iam.CreateInstanceProfileInput{
				InstanceProfileName: profileName,
				Tags:                tags,
			}); err != nil {
				return fmt.Errorf("failed to create %s instance profile: %w", role, err)
			}
			if err := svc.WaitUntilInstanceProfileExistsWithContext(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
				return fmt.Errorf("failed to wait for %s instance profile to exist: %w", role, err)
			}

			// Finally, attach the role to the profile.
			if _, err := svc.AddRoleToInstanceProfileWithContext(ctx, &iam.AddRoleToInstanceProfileInput{
				InstanceProfileName: profileName,
				RoleName:            aws.String(roleName),
			}); err != nil {
				return fmt.Errorf("failed to add %s role to instance profile: %w", role, err)
			}
		}
	}

	return nil
}

// getOrCreateRole returns the name of the IAM role to be used,
// creating it when not specified by the user in the install config.
func getOrCreateIAMRole(ctx context.Context, nodeRole, infraID, assumePolicy string, ic installconfig.InstallConfig, tags []*iam.Tag, svc *iam.IAM) (string, error) {
	roleName := aws.String(fmt.Sprintf("%s-%s-role", infraID, nodeRole))

	var defaultRole string
	if dmp := ic.Config.AWS.DefaultMachinePlatform; dmp != nil && len(dmp.IAMRole) > 0 {
		defaultRole = dmp.IAMRole
	}

	masterRole := defaultRole
	if cp := ic.Config.ControlPlane; cp != nil && cp.Platform.AWS != nil && len(cp.Platform.AWS.IAMRole) > 0 {
		masterRole = cp.Platform.AWS.IAMRole
	}

	workerRole := defaultRole
	if w := ic.Config.Compute; len(w) > 0 && w[0].Platform.AWS != nil && len(w[0].Platform.AWS.IAMRole) > 0 {
		workerRole = w[0].Platform.AWS.IAMRole
	}

	switch {
	case nodeRole == master && len(masterRole) > 0:
		return masterRole, nil
	case nodeRole == worker && len(workerRole) > 0:
		return workerRole, nil
	}

	if _, err := svc.GetRoleWithContext(ctx, &iam.GetRoleInput{RoleName: roleName}); err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() != iam.ErrCodeNoSuchEntityException {
			return "", fmt.Errorf("failed to get %s role: %w", nodeRole, err)
		}
		// If the role does not exist, create it.
		logrus.Infof("Creating IAM role for %s", nodeRole)
		createRoleInput := &iam.CreateRoleInput{
			RoleName:                 roleName,
			AssumeRolePolicyDocument: aws.String(assumePolicy),
			Tags:                     tags,
		}
		if _, err := svc.CreateRoleWithContext(ctx, createRoleInput); err != nil {
			return "", fmt.Errorf("failed to create %s role: %w", nodeRole, err)
		}

		if err := svc.WaitUntilRoleExistsWithContext(ctx, &iam.GetRoleInput{RoleName: roleName}); err != nil {
			return "", fmt.Errorf("failed to wait for %s role to exist: %w", nodeRole, err)
		}
	}

	// Put the policy inline.
	policyName := aws.String(fmt.Sprintf("%s-%s-policy", infraID, nodeRole))
	b, err := json.Marshal(policies[nodeRole])
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s policy: %w", nodeRole, err)
	}
	if _, err := svc.PutRolePolicyWithContext(ctx, &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(string(b)),
		PolicyName:     policyName,
		RoleName:       roleName,
	}); err != nil {
		return "", fmt.Errorf("failed to create inline policy for role %s: %w", nodeRole, err)
	}

	return *roleName, nil
}

func getPartitionService(region string) string {
	partitionDNSSuffix := "amazonaws.com"
	if ps, found := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region); found {
		partitionDNSSuffix = ps.DNSSuffix()
	}
	return fmt.Sprintf("ec2.%s", partitionDNSSuffix)
}
