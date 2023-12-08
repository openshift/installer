package aws

import iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"

var (
	policies = map[string]*iamv1.PolicyDocument{
		"master": {
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
		"worker": {
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

	assumePolicy = &iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: iamv1.Statements{
			{
				Effect: "Allow",
				Principal: iamv1.Principals{
					iamv1.PrincipalService: []string{
						"ec2.amazonaws.com",
					},
				},
				Action: iamv1.Actions{
					"sts:AssumeRole",
				},
			},
		},
	}
)
