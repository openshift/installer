package aws

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/network"
)

// BootstrapSSHDescription is the description for the
// ingress rule that provides SSH access to the bootstrap node
// & identifies the rule for removal during bootstrap destroy.
const BootstrapSSHDescription = "Bootstrap SSH Access"

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(ic *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	platformAWS := ic.Config.AWS
	enableIPv6 := platformAWS.IPFamily.DualStackEnabled()

	tags, err := aws.CapaTagsFromUserTags(clusterID.InfraID, platformAWS.UserTags)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tags: %w", err)
	}

	var sshRuleCidrs []ipnet.IPNet
	if !ic.Config.PublicAPI() {
		sshRuleCidrs = capiutils.MachineCIDRsFromInstallConfig(ic)
	} else {
		sshRuleCidrs = []ipnet.IPNet{*capiutils.AnyIPv4CidrBlock}
		if enableIPv6 {
			sshRuleCidrs = append(sshRuleCidrs, *capiutils.AnyIPv6CidrBlock)
		}
	}

	targetGroupIPType := GetTargetGroupIPType(platformAWS.IPFamily)
	awsCluster := &capa.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capa.AWSClusterSpec{
			Region: platformAWS.Region,
			NetworkSpec: capa.NetworkSpec{
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
					},
				},
				AdditionalControlPlaneIngressRules: []capa.IngressRule{
					{
						Description:              "MCS traffic from cluster network",
						Protocol:                 capa.SecurityGroupProtocolTCP,
						FromPort:                 22623,
						ToPort:                   22623,
						SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane", "apiserver-lb"},
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
						Description:    BootstrapSSHDescription,
						Protocol:       capa.SecurityGroupProtocolTCP,
						FromPort:       22,
						ToPort:         22,
						CidrBlocks:     capiutils.CIDRsToString(capiutils.GetIPv4CIDRs(sshRuleCidrs)),
						IPv6CidrBlocks: capiutils.CIDRsToString(capiutils.GetIPv6CIDRs(sshRuleCidrs)),
					},
				},
				// If the installer provisions the VPC, VPC IPv6 CIDR is unknown at install time and added after infraReady
				NodePortIngressRuleCidrBlocks: capiutils.CIDRsToString(capiutils.MachineCIDRsFromInstallConfig(ic)),
			},
			S3Bucket: &capa.S3Bucket{
				Name:                    GetIgnitionBucketName(clusterID.InfraID),
				PresignedURLDuration:    &metav1.Duration{Duration: 1 * time.Hour},
				BestEffortDeleteObjects: ptr.To(ic.Config.AWS.BestEffortDeleteIgnition),
			},
			ControlPlaneLoadBalancer: &capa.AWSLoadBalancerSpec{
				Name:                   ptr.To(clusterID.InfraID + "-int"),
				LoadBalancerType:       capa.LoadBalancerTypeNLB,
				Scheme:                 &capa.ELBSchemeInternal,
				CrossZoneLoadBalancing: true,
				HealthCheckProtocol:    &capa.ELBProtocolHTTPS,
				HealthCheck: &capa.TargetGroupHealthCheckAPISpec{
					IntervalSeconds:         ptr.To[int64](10),
					TimeoutSeconds:          ptr.To[int64](10),
					ThresholdCount:          ptr.To[int64](2),
					UnhealthyThresholdCount: ptr.To[int64](2),
				},
				TargetGroupIPType: targetGroupIPType,
				AdditionalListeners: []capa.AdditionalListenerSpec{
					{
						Port:     22623,
						Protocol: capa.ELBProtocolTCP,
						HealthCheck: &capa.TargetGroupHealthCheckAdditionalSpec{
							Protocol:                ptr.To[string](capa.ELBProtocolHTTPS.String()),
							Port:                    ptr.To[string]("22623"),
							Path:                    ptr.To[string]("/healthz"),
							IntervalSeconds:         ptr.To[int64](10),
							TimeoutSeconds:          ptr.To[int64](10),
							ThresholdCount:          ptr.To[int64](2),
							UnhealthyThresholdCount: ptr.To[int64](2),
						},
						TargetGroupIPType: targetGroupIPType,
					},
				},
				IngressRules: []capa.IngressRule{
					{
						Description:              "Machine Config Server internal traffic from cluster",
						Protocol:                 capa.SecurityGroupProtocolTCP,
						FromPort:                 22623,
						ToPort:                   22623,
						SourceSecurityGroupRoles: []capa.SecurityGroupRole{"node", "controlplane"},
					},
				},
			},
			AdditionalTags: tags,
		},
	}
	awsCluster.SetGroupVersionKind(capa.GroupVersion.WithKind("AWSCluster"))

	if enableIPv6 {
		awsCluster.Spec.NetworkSpec.CNI.CNIIngressRules = append(awsCluster.Spec.NetworkSpec.CNI.CNIIngressRules,
			capa.CNIIngressRule{
				Description: "ICMPv6",
				Protocol:    capa.SecurityGroupProtocolICMPv6,
				FromPort:    -1,
				ToPort:      -1,
			},
		)
	}

	// Create a ingress rule to allow acccess to the API LB.
	apiLBIngressRule := capa.IngressRule{
		Description: "Kubernetes API Server traffic",
		Protocol:    capa.SecurityGroupProtocolTCP,
		FromPort:    6443,
		ToPort:      6443,
		CidrBlocks:  []string{capiutils.AnyIPv4CidrBlock.String()},
	}
	if enableIPv6 {
		apiLBIngressRule.IPv6CidrBlocks = []string{capiutils.AnyIPv6CidrBlock.String()}
	}

	if ic.Config.PublicAPI() {
		apiLBIngressRule.Description = "Kubernetes API Server traffic for public access"
		awsCluster.Spec.SecondaryControlPlaneLoadBalancer = &capa.AWSLoadBalancerSpec{
			Name:                   ptr.To(clusterID.InfraID + "-ext"),
			LoadBalancerType:       capa.LoadBalancerTypeNLB,
			Scheme:                 &capa.ELBSchemeInternetFacing,
			CrossZoneLoadBalancing: true,
			HealthCheckProtocol:    &capa.ELBProtocolHTTPS,
			HealthCheck: &capa.TargetGroupHealthCheckAPISpec{
				IntervalSeconds:         ptr.To[int64](10),
				TimeoutSeconds:          ptr.To[int64](10),
				ThresholdCount:          ptr.To[int64](2),
				UnhealthyThresholdCount: ptr.To[int64](2),
			},
			TargetGroupIPType: targetGroupIPType,
			IngressRules:      []capa.IngressRule{apiLBIngressRule},
		}
	} else {
		awsCluster.Spec.ControlPlaneLoadBalancer.IngressRules = append(
			awsCluster.Spec.ControlPlaneLoadBalancer.IngressRules,
			apiLBIngressRule,
		)
	}

	// Set the NetworkSpec.Subnets from VPC and zones (managed) or subnets (BYO VPC) based in the install-config.yaml.
	// If subnet roles are assigned, set subnets for the ControlPlane LBs.
	err = setSubnets(context.TODO(), &networkInput{
		InstallConfig: ic,
		ClusterID:     clusterID,
		Cluster:       awsCluster,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set cluster zones or subnets: %w", err)
	}

	// Enable BYO Public IPv4 when defined on install-config.yaml
	if len(ic.Config.Platform.AWS.PublicIpv4Pool) > 0 {
		awsCluster.Spec.NetworkSpec.VPC.ElasticIPPool = &capa.ElasticIPPool{
			PublicIpv4Pool:              ptr.To(ic.Config.Platform.AWS.PublicIpv4Pool),
			PublicIpv4PoolFallBackOrder: ptr.To(capa.PublicIpv4PoolFallbackOrderAmazonPool),
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
	id.SetGroupVersionKind(capa.GroupVersion.WithKind("AWSClusterControllerIdentity"))
	manifests = append(manifests, &asset.RuntimeFile{
		Object: id,
		File:   asset.File{Filename: "01_aws-cluster-controller-identity-default.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: capa.GroupVersion.String(),
				Kind:       "AWSCluster",
				Name:       awsCluster.Name,
				Namespace:  awsCluster.Namespace,
			},
		},
	}, nil
}

// GetIgnitionBucketName returns the name of the bucket for the given cluster.
func GetIgnitionBucketName(infraID string) string {
	return fmt.Sprintf("openshift-bootstrap-data-%s", infraID)
}

// GetTargetGroupIPType returns the ipType of the target group based on ipFamily.
func GetTargetGroupIPType(ipFamily network.IPFamily) *capa.TargetGroupIPType {
	var tgIPType capa.TargetGroupIPType
	switch ipFamily {
	case network.DualStackIPv6Primary:
		tgIPType = capa.TargetGroupIPTypeIPv6
	case network.DualStackIPv4Primary:
		tgIPType = capa.TargetGroupIPTypeIPv4
	default:
		// Default to IPv4 if not specified or invalid
		tgIPType = capa.TargetGroupIPTypeIPv4
	}

	return &tgIPType
}
