package aws

import (
	"fmt"
	"time"

	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/types"
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
func GenerateClusterAssets(ic *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}

	tags, err := aws.CapaTagsFromUserTags(clusterID.InfraID, ic.Config.AWS.UserTags)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tags: %w", err)
	}

	awsCluster := &capa.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capa.AWSClusterSpec{
			Region:      ic.Config.AWS.Region,
			NetworkSpec: capa.NetworkSpec{},
			S3Bucket: &capa.S3Bucket{
				Name:                 fmt.Sprintf("openshift-bootstrap-data-%s", clusterID.InfraID),
				PresignedURLDuration: &metav1.Duration{Duration: 1 * time.Hour},
			},
			AdditionalTags: tags,
		},
	}

	// Setup Security Group
	awsCluster.Spec.NetworkSpec.CNI = &capa.CNISpec{CNIIngressRules: getDefaultNetworkCNIIngressRules()}
	awsCluster.Spec.NetworkSpec.AdditionalControlPlaneIngressRules = getDefaultNetworkAdditionalControlPlaneIngressRules()

	// Setup Load Balancers
	// Primary == internal LB
	mainCIDR := capiutils.CIDRFromInstallConfig(ic)
	awsCluster.Spec.ControlPlaneLoadBalancer = &capa.AWSLoadBalancerSpec{
		Name:             ptr.To(clusterID.InfraID + "-int"),
		LoadBalancerType: capa.LoadBalancerTypeNLB,
		Scheme:           &capa.ELBSchemeInternal,
		HealthCheck: &capa.TargetGroupHealthCheck{
			Protocol:                ptr.To("HTTPS"),
			Path:                    ptr.To("/readyz"),
			IntervalSeconds:         ptr.To(int64(10)),
			TimeoutSeconds:          ptr.To(int64(10)),
			ThresholdCount:          ptr.To(int64(2)),
			UnhealthyThresholdCount: ptr.To(int64(2)),
		},
		AdditionalListeners: []capa.AdditionalListenerSpec{
			{
				Port:     22623,
				Protocol: capa.ELBProtocolTCP,
				HealthCheck: &capa.TargetGroupHealthCheck{
					Protocol:                ptr.To("HTTPS"),
					Path:                    ptr.To("/healthz"),
					IntervalSeconds:         ptr.To(int64(10)),
					TimeoutSeconds:          ptr.To(int64(10)),
					ThresholdCount:          ptr.To(int64(2)),
					UnhealthyThresholdCount: ptr.To(int64(2)),
				},
			},
		},
		IngressRules: []capa.IngressRule{
			{
				Description: "Machine Config Server internal traffic from cluster",
				Protocol:    capa.SecurityGroupProtocolTCP,
				FromPort:    22623,
				ToPort:      22623,
				CidrBlocks:  []string{mainCIDR.String()},
			},
		},
	}

	if ic.Config.Publish == types.ExternalPublishingStrategy {
		awsCluster.Spec.ControlPlaneLoadBalancer.IngressRules = append(awsCluster.Spec.ControlPlaneLoadBalancer.IngressRules, capa.IngressRule{
			Description: "Kubernetes API Server traffic for public access",
			Protocol:    capa.SecurityGroupProtocolTCP,
			FromPort:    6443,
			ToPort:      6443,
			CidrBlocks:  []string{"0.0.0.0/0"},
		})
		awsCluster.Spec.SecondaryControlPlaneLoadBalancer = &capa.AWSLoadBalancerSpec{
			Name:                   ptr.To(clusterID.InfraID + "-ext"),
			LoadBalancerType:       capa.LoadBalancerTypeNLB,
			Scheme:                 &capa.ELBSchemeInternetFacing,
			CrossZoneLoadBalancing: true,
		}
	}

	// Set the VPC and zones (managed) or subnets (BYO VPC) based in the
	// install-config.yaml.
	err = setZones(&zoneConfigInput{
		InstallConfig: ic,
		Config:        ic.Config,
		Meta:          ic.AWS,
		ClusterID:     clusterID,
		Cluster:       awsCluster,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to set cluster zones or subnets")
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
