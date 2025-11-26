package clusterapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/sirupsen/logrus"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
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

	platformAWS := ic.Config.Platform.AWS
	client, err := awsconfig.NewIAMClient(ctx, awsconfig.EndpointOptions{
		Region:    platformAWS.Region,
		Endpoints: platformAWS.ServiceEndpoints,
	})
	if err != nil {
		return fmt.Errorf("failed to create iam client: %w", err)
	}

	// Create the IAM Roles for master and workers.
	tags := []iamtypes.Tag{
		{
			Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", infraID)),
			Value: aws.String("owned"),
		},
	}

	for k, v := range ic.Config.AWS.UserTags {
		tags = append(tags, iamtypes.Tag{
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
						getPartitionDNSSuffix(ic.AWS.Region),
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

	var defaultProfile string
	if dmp := ic.Config.AWS.DefaultMachinePlatform; dmp != nil && len(dmp.IAMProfile) > 0 {
		defaultProfile = dmp.IAMProfile
	}

	for _, role := range []string{master, worker} {
		instanceProfile := defaultProfile
		switch role {
		case master:
			if cp := ic.Config.ControlPlane; cp != nil && cp.Platform.AWS != nil && len(cp.Platform.AWS.IAMProfile) > 0 {
				instanceProfile = cp.Platform.AWS.IAMProfile
			}
		case worker:
			if w := ic.Config.Compute; len(w) > 0 && w[0].Platform.AWS != nil && len(w[0].Platform.AWS.IAMProfile) > 0 {
				instanceProfile = w[0].Platform.AWS.IAMProfile
			}
		}

		// A user-provided instance profile already has a role attached to it, so there is nothing else for the
		// Installer to do.
		if len(instanceProfile) > 0 {
			logrus.Debugf("Using existing %s instance profile %q", role, instanceProfile)
			continue
		}

		roleName, err := getOrCreateIAMRole(ctx, role, infraID, string(assumePolicyBytes), *ic, tags, client)
		if err != nil {
			return fmt.Errorf("failed to create IAM %s role: %w", role, err)
		}

		profileName := aws.String(fmt.Sprintf("%s-%s-profile", infraID, role))
		if _, err := client.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: profileName}); err != nil {
			var noSuchEntity *iamtypes.NoSuchEntityException
			if !errors.As(err, &noSuchEntity) {
				return fmt.Errorf("failed to get %s instance profile: %w", role, err)
			}
			// If the profile does not exist, create it.
			if _, err := client.CreateInstanceProfile(ctx, &iam.CreateInstanceProfileInput{
				InstanceProfileName: profileName,
				Tags:                tags,
			}); err != nil {
				return fmt.Errorf("failed to create %s instance profile: %w", role, err)
			}

			waiter := iam.NewInstanceProfileExistsWaiter(client)
			if err := waiter.Wait(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: profileName}, 15*time.Minute); err != nil {
				return fmt.Errorf("failed to wait for %s instance profile to exist: %w", role, err)
			}

			// Finally, attach the role to the profile.
			if _, err := client.AddRoleToInstanceProfile(ctx, &iam.AddRoleToInstanceProfileInput{
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
func getOrCreateIAMRole(ctx context.Context, nodeRole, infraID, assumePolicy string, ic installconfig.InstallConfig, tags []iamtypes.Tag, svc *iam.Client) (string, error) {
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

	if _, err := svc.GetRole(ctx, &iam.GetRoleInput{RoleName: roleName}); err != nil {
		var noSuchRoleError *iamtypes.NoSuchEntityException
		if !errors.As(err, &noSuchRoleError) {
			return "", fmt.Errorf("failed to get %s role: %w", nodeRole, err)
		}

		// If the role does not exist, create it.
		logrus.Infof("Creating IAM role for %s", nodeRole)
		createRoleInput := &iam.CreateRoleInput{
			RoleName:                 roleName,
			AssumeRolePolicyDocument: aws.String(assumePolicy),
			Tags:                     tags,
		}
		if _, err := svc.CreateRole(ctx, createRoleInput); err != nil {
			return "", fmt.Errorf("failed to create %s role: %w", nodeRole, err)
		}
		waiter := iam.NewRoleExistsWaiter(svc)
		if err := waiter.Wait(ctx, &iam.GetRoleInput{RoleName: roleName}, 15*time.Minute); err != nil {
			return "", fmt.Errorf("failed to wait for %s role to exist: %w", nodeRole, err)
		}
	}

	// Put the policy inline.
	policyName := aws.String(fmt.Sprintf("%s-%s-policy", infraID, nodeRole))
	b, err := json.Marshal(policies[nodeRole])
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s policy: %w", nodeRole, err)
	}
	if _, err := svc.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(string(b)),
		PolicyName:     policyName,
		RoleName:       roleName,
	}); err != nil {
		return "", fmt.Errorf("failed to create inline policy for role %s: %w", nodeRole, err)
	}

	return *roleName, nil
}

func getPartitionDNSSuffix(region string) string {
	endpoint, err := ec2.NewDefaultEndpointResolver().ResolveEndpoint(region, ec2.EndpointResolverOptions{})
	if err != nil {
		logrus.Errorf("failed to resolve AWS ec2 endpoint: %v", err)
		return ""
	}

	u, err := url.Parse(endpoint.URL)
	if err != nil {
		logrus.Errorf("failed to parse partition ID URL: %v", err)
		return ""
	}

	domain := "amazonaws.com"
	// Extract the hostname
	host := u.Hostname()
	// Split the hostname by "." to get the domain parts
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		domain = strings.Join(parts[2:], ".")
	}

	logrus.Debugf("Using domain name: %s", domain)
	return fmt.Sprintf("ec2.%s", domain)
}
