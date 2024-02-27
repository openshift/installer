package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	mainCIDR := capiutils.CIDRFromInstallConfig(installConfig)

	zones, err := installConfig.AWS.AvailabilityZones(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get availability zones")
	}

	awsCluster := &capa.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capa.AWSClusterSpec{
			Region: installConfig.Config.AWS.Region,
			NetworkSpec: capa.NetworkSpec{
				VPC: capa.VPCSpec{
					CidrBlock:                  mainCIDR.String(),
					AvailabilityZoneUsageLimit: ptr.To(len(zones)),
					AvailabilityZoneSelection:  &capa.AZSelectionSchemeOrdered,
				},
				CNI: &capa.CNISpec{
					CNIIngressRules: capa.CNIIngressRules{
						{
							Description: "ICMP",
							Protocol:    capa.SecurityGroupProtocolICMP,
							FromPort:    -1,
							ToPort:      -1,
						},
						{
							Description: "Port 22 (TCP)",
							Protocol:    capa.SecurityGroupProtocolTCP,
							FromPort:    22,
							ToPort:      22,
						},
						{
							Description: "Port 4789 (UDP) for VXLAN",
							Protocol:    capa.SecurityGroupProtocolUDP,
							FromPort:    4789,
							ToPort:      4789,
						},
						{
							Description: "Port 6081 (UDP) for geneve",
							Protocol:    capa.SecurityGroupProtocolUDP,
							FromPort:    6081,
							ToPort:      6081,
						},
						{
							Description: "Port 500 (UDP) for IKE",
							Protocol:    capa.SecurityGroupProtocolUDP,
							FromPort:    500,
							ToPort:      500,
						},
						{
							Description: "Port 4500 (UDP) for IKE NAT",
							Protocol:    capa.SecurityGroupProtocolUDP,
							FromPort:    4500,
							ToPort:      4500,
						},
						{
							Description: "ESP",
							Protocol:    capa.SecurityGroupProtocolESP,
							FromPort:    -1,
							ToPort:      -1,
						},
						{
							Description: "Port 6441-6442 (TCP) for ovndb",
							Protocol:    capa.SecurityGroupProtocolTCP,
							FromPort:    6441,
							ToPort:      6442,
						},
						{
							Description: "Port 9000-9999 for node ports (TCP)",
							Protocol:    capa.SecurityGroupProtocolTCP,
							FromPort:    9000,
							ToPort:      9999,
						},
						{
							Description: "Port 9000-9999 for node ports (UDP)",
							Protocol:    capa.SecurityGroupProtocolUDP,
							FromPort:    9000,
							ToPort:      9999,
						},
						{
							Description: "Service node ports (TCP)",
							Protocol:    capa.SecurityGroupProtocolTCP,
							FromPort:    30000,
							ToPort:      32767,
						},
						{
							Description: "Service node ports (UDP)",
							Protocol:    capa.SecurityGroupProtocolUDP,
							FromPort:    30000,
							ToPort:      32767,
						},
						{
							Description: "Machine Config Server (MCS)",
							Protocol:    capa.SecurityGroupProtocolTCP,
							FromPort:    22623,
							ToPort:      22623,
						},
					},
				},
				AdditionalControlPlaneIngressRules: []capa.IngressRule{
					{
						Description: "MCS traffic from cluster network",
						Protocol:    capa.SecurityGroupProtocolTCP,
						FromPort:    22623,
						ToPort:      22623,
						//SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
						CidrBlocks: []string{"10.0.0.0/16"}, //TODO(padillon): figure out security group rules
					},
					{
						Description:              "controller-manager",
						Protocol:                 capa.SecurityGroupProtocolTCP,
						FromPort:                 10257,
						ToPort:                   10257,
						SourceSecurityGroupRoles: []capa.SecurityGroupRole{"controlplane", "node"},
					},
					{
						Description:              "kube-scheduler",
						Protocol:                 capa.SecurityGroupProtocolTCP,
						FromPort:                 10259,
						ToPort:                   10259,
						SourceSecurityGroupRoles: []capa.SecurityGroupRole{"controlplane", "node"},
					},
					{
						Description: "SSH everyone",
						Protocol:    capa.SecurityGroupProtocolTCP,
						FromPort:    22,
						ToPort:      22,
						CidrBlocks:  []string{"0.0.0.0/0"},
					},
					{
						Description: "public api", //TESTING
						Protocol:    capa.SecurityGroupProtocolTCP,
						FromPort:    6443,
						ToPort:      6443,
						CidrBlocks:  []string{"0.0.0.0/0"},
					},
				},
			},
			S3Bucket: &capa.S3Bucket{
				Name:                 fmt.Sprintf("openshift-bootstrap-data-%s", clusterID.InfraID),
				PresignedURLDuration: &metav1.Duration{Duration: 1 * time.Hour},
			},
			ControlPlaneLoadBalancer: &capa.AWSLoadBalancerSpec{
				Name:             ptr.To(clusterID.InfraID + "-int"),
				LoadBalancerType: capa.LoadBalancerTypeNLB,
				Scheme:           &capa.ELBSchemeInternal,
				AdditionalListeners: []capa.AdditionalListenerSpec{
					{
						Port:     22623,
						Protocol: capa.ELBProtocolTCP,
					},
				},
				IngressRules: []capa.IngressRule{
					{
						Description: "MCS traffic from cluster network",
						Protocol:    capa.SecurityGroupProtocolTCP,
						FromPort:    22623,
						ToPort:      22623,
						//SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
						CidrBlocks: []string{"10.0.0.0/16"}, //TODO(padillon): figure out security group rules
					},
					{
						Description: "public api", //TESTING. This doesn't really belong on the internal LB...
						Protocol:    capa.SecurityGroupProtocolTCP,
						FromPort:    6443,
						ToPort:      6443,
						CidrBlocks:  []string{"0.0.0.0/0"},
					},
				},
			},
			SecondaryControlPlaneLoadBalancer: &capa.AWSLoadBalancerSpec{
				Name:             ptr.To(clusterID.InfraID + "-ext"),
				LoadBalancerType: capa.LoadBalancerTypeNLB,
				Scheme:           &capa.ELBSchemeInternetFacing,
				// IngressRules: []capa.IngressRule{ //THIS doesn't seem to update LB security group
				// 	{
				// 		Description: "public api", //TESTING
				// 		Protocol:    capa.SecurityGroupProtocolTCP,
				// 		FromPort:    6443,
				// 		ToPort:      6443,
				// 		CidrBlocks:  []string{"0.0.0.0/0"},
				// 	},
				// },
			},
			AdditionalTags: capa.Tags{fmt.Sprintf("kubernetes.io/cluster/%s", clusterID.InfraID): "owned"},
		},
	}

	// If the install config has subnets, use them.
	if len(installConfig.AWS.Subnets) > 0 {
		privateSubnets, err := installConfig.AWS.PrivateSubnets(context.TODO())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get private subnets")
		}
		for _, subnet := range privateSubnets {
			awsCluster.Spec.NetworkSpec.Subnets = append(awsCluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
				ID:               subnet.ID,
				CidrBlock:        subnet.CIDR,
				AvailabilityZone: subnet.Zone.Name,
				IsPublic:         subnet.Public,
			})
		}
		publicSubnets, err := installConfig.AWS.PublicSubnets(context.TODO())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get public subnets")
		}

		for _, subnet := range publicSubnets {
			awsCluster.Spec.NetworkSpec.Subnets = append(awsCluster.Spec.NetworkSpec.Subnets, capa.SubnetSpec{
				ID:               subnet.ID,
				CidrBlock:        subnet.CIDR,
				AvailabilityZone: subnet.Zone.Name,
				IsPublic:         subnet.Public,
			})
		}

		vpc, err := installConfig.AWS.VPC(context.TODO())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get VPC")
		}
		awsCluster.Spec.NetworkSpec.VPC = capa.VPCSpec{
			ID: vpc,
		}
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: awsCluster,
		File:   asset.File{Filename: "02_infra-cluster.yaml"},
	})

	id := &capa.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: capiutils.Namespace,
		},
		Spec: capa.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: capa.AWSClusterIdentitySpec{
				AllowedNamespaces: &capa.AllowedNamespaces{}, // Allow all namespaces.
			},
		},
	}
	manifests = append(manifests, &asset.RuntimeFile{
		Object: id,
		File:   asset.File{Filename: "01_aws-cluster-controller-identity-default.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRef: &corev1.ObjectReference{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
			Kind:       "AWSCluster",
			Name:       awsCluster.Name,
			Namespace:  awsCluster.Namespace,
		},
	}, nil
}
