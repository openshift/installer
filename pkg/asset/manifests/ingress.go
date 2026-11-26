package manifests

import (
	"context"
	"fmt"
	"path"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	clusterIngressConfigFile     = path.Join(manifestDir, "cluster-ingress-02-config.yml")
	defaultIngressControllerFile = path.Join(manifestDir, "cluster-ingress-default-ingresscontroller.yaml")
)

// Ingress generates the cluster-ingress-*.yml files.
type Ingress struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*Ingress)(nil)

// Name returns a human friendly name for the asset.
func (*Ingress) Name() string {
	return "Ingress Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*Ingress) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the ingress cluster config and default ingresscontroller.
//
// A cluster ingress config is always created.
//
// A default ingresscontroller is only created if the cluster is using an internal
// publishing strategy. In this case, the default ingresscontroller is also set
// to use the internal publishing strategy.
func (ing *Ingress) Generate(_ context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)

	ing.FileList = []*asset.File{}

	clusterConfig, err := ing.generateClusterConfig(installConfig.Config)
	if err != nil {
		return errors.Wrap(err, "failed to create cluster config")
	}
	ing.FileList = append(ing.FileList, &asset.File{
		Filename: clusterIngressConfigFile,
		Data:     clusterConfig,
	})

	defaultIngressController, err := ing.generateDefaultIngressController(installConfig.Config)
	if err != nil {
		return errors.Wrap(err, "failed to create default ingresscontroller")
	}
	if len(defaultIngressController) > 0 {
		ing.FileList = append(ing.FileList, &asset.File{
			Filename: defaultIngressControllerFile,
			Data:     defaultIngressController,
		})
	}

	return nil
}

func (ing *Ingress) generateClusterConfig(config *types.InstallConfig) ([]byte, error) {
	controlPlaneTopology, _ := determineTopologies(config)

	isSingleControlPlaneNode := controlPlaneTopology == configv1.SingleReplicaTopologyMode

	defaultPlacement := configv1.DefaultPlacementWorkers
	if config.Platform.None != nil && isSingleControlPlaneNode {
		// A none-platform single control-plane node cluster doesn't need a
		// load balancer, the API and ingress traffic for such cluster can be
		// directed at the single node directly. We want to maintain that even
		// when worker nodes are added to such cluster. We do that by asking
		// the Cluster Ingress Operator to place the ingress pod on the single
		// control plane node. This would ensure that the ingress pod won't
		// ever be scheduled on the added workers, as that would create a
		// requirement for a load-balancer. Even when a single control-plane
		// node cluster is installed with one or more worker nodes since day 1,
		// we still want control-plane ingress default placement, for the sake
		// of consistency. Users can override this decision manually if they
		// wish.
		defaultPlacement = configv1.DefaultPlacementControlPlane
	}

	obj := &configv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.GroupVersion.String(),
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.IngressSpec{
			Domain: fmt.Sprintf("apps.%s", config.ClusterDomain()),
		},
		Status: configv1.IngressStatus{
			DefaultPlacement: defaultPlacement,
		},
	}

	switch config.Platform.Name() {
	case aws.Name:
		lbType := configv1.Classic
		if config.AWS.LBType == configv1.NLB {
			lbType = configv1.NLB
		}
		obj.Spec.LoadBalancer = configv1.LoadBalancer{
			Platform: configv1.IngressPlatformSpec{
				AWS: &configv1.AWSIngressSpec{
					Type: lbType,
				},
				Type: configv1.AWSPlatformType,
			},
		}
	}
	return yaml.Marshal(obj)
}

func (ing *Ingress) generateDefaultIngressController(config *types.InstallConfig) ([]byte, error) {
	// getDefaultIngressController returns an object representing the default Ingress Controller
	// with empty LoadBalancer spec.
	getDefaultIngressController := func() *operatorv1.IngressController {
		return &operatorv1.IngressController{
			TypeMeta: metav1.TypeMeta{
				APIVersion: operatorv1.GroupVersion.String(),
				Kind:       "IngressController",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-ingress-operator",
				Name:      "default",
			},
			Spec: operatorv1.IngressControllerSpec{
				EndpointPublishingStrategy: &operatorv1.EndpointPublishingStrategy{
					Type:         operatorv1.LoadBalancerServiceStrategyType,
					LoadBalancer: &operatorv1.LoadBalancerStrategy{},
				},
			},
		}
	}

	var obj *operatorv1.IngressController

	switch config.Platform.Name() {
	case aws.Name:
		subnetIDsByRole := make(map[aws.SubnetRoleType][]operatorv1.AWSSubnetID)
		for _, subnet := range config.AWS.VPC.Subnets {
			for _, role := range subnet.Roles {
				subnetIDsByRole[role.Type] = append(subnetIDsByRole[role.Type], operatorv1.AWSSubnetID(subnet.ID))
			}
		}

		// BYO-subnet install case and subnet roles are specified.
		if len(subnetIDsByRole) > 0 {
			obj = getDefaultIngressController()
			lbSpec := obj.Spec.EndpointPublishingStrategy.LoadBalancer

			if config.PublicIngress() {
				lbSpec.Scope = operatorv1.ExternalLoadBalancer
			} else {
				lbSpec.Scope = operatorv1.InternalLoadBalancer
			}

			if config.AWS.LBType == configv1.NLB {
				lbSpec.ProviderParameters = &operatorv1.ProviderLoadBalancerParameters{
					Type: operatorv1.AWSLoadBalancerProvider,
					AWS: &operatorv1.AWSLoadBalancerParameters{
						Type: operatorv1.AWSNetworkLoadBalancer,
						NetworkLoadBalancerParameters: &operatorv1.AWSNetworkLoadBalancerParameters{
							Subnets: &operatorv1.AWSSubnets{
								IDs: subnetIDsByRole[aws.IngressControllerLBSubnetRole],
							},
						},
					},
				}
			} else {
				lbSpec.ProviderParameters = &operatorv1.ProviderLoadBalancerParameters{
					Type: operatorv1.AWSLoadBalancerProvider,
					AWS: &operatorv1.AWSLoadBalancerParameters{
						Type: operatorv1.AWSClassicLoadBalancer,
						ClassicLoadBalancerParameters: &operatorv1.AWSClassicLoadBalancerParameters{
							Subnets: &operatorv1.AWSSubnets{
								IDs: subnetIDsByRole[aws.IngressControllerLBSubnetRole],
							},
						},
					},
				}
			}
			break
		}
		// Fall back to existing logic similar to other platforms otherwise.
		// i.e. managed subnets or no subnet roles are specified.
		fallthrough
	default:
		// The Ingress Operator creates the default Ingress Controller with scope External if none is provided.
		// https://github.com/openshift/cluster-ingress-operator/blob/81e314f2e9b41b6616ad2c3db657e869915577a8/pkg/operator/operator.go#L470-L516
		// Thus, if ingress LB scope is configured to be internal, we need to generate the default Ingress Controller with scope Internal.
		if !config.PublicIngress() {
			obj = getDefaultIngressController()
			obj.Spec.EndpointPublishingStrategy.LoadBalancer.Scope = operatorv1.InternalLoadBalancer
		}
	}

	if obj != nil {
		return yaml.Marshal(obj)
	}
	return nil, nil
}

// Files returns the files generated by the asset.
func (ing *Ingress) Files() []*asset.File {
	return ing.FileList
}

// Load returns false since this asset is not written to disk by the installer.
func (ing *Ingress) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}
