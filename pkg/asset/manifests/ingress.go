package manifests

import (
	"context"
	"fmt"
	"path/filepath"

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
	clusterIngressConfigFile     = filepath.Join(manifestDir, "cluster-ingress-02-config.yml")
	defaultIngressControllerFile = filepath.Join(manifestDir, "cluster-ingress-default-ingresscontroller.yaml")
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
	switch config.Publish {
	case types.MixedPublishingStrategy:
		if config.OperatorPublishingStrategy.Ingress != "Internal" {
			break
		}
		fallthrough
	case types.InternalPublishingStrategy:
		obj := &operatorv1.IngressController{
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
					Type: operatorv1.LoadBalancerServiceStrategyType,
					LoadBalancer: &operatorv1.LoadBalancerStrategy{
						Scope: operatorv1.InternalLoadBalancer,
					},
				},
			},
		}
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
