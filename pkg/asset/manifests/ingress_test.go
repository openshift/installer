package manifests

import (
	"testing"

	"github.com/ghodss/yaml"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

func TestGenerateDefaultIngressController(t *testing.T) {
	cases := []struct {
		name                      string
		installConfig             *types.InstallConfig
		expectedIngressController *operatorv1.IngressController
	}{{
		name: "public",
		installConfig: icBuild.build(
			icBuild.withPublish(types.ExternalPublishingStrategy),
		),
		expectedIngressController: nil,
	}, {
		name: "private highly available cluster",
		installConfig: icBuild.build(
			icBuild.withPublish(types.InternalPublishingStrategy),
		),
		expectedIngressController: ingresscontrollerBuild.build(
			ingresscontrollerBuild.withLoadBalancer(),
			ingresscontrollerBuild.withScope(operatorv1.InternalLoadBalancer),
		),
	}, {
		name: "private single-node cluster",
		installConfig: icBuild.build(
			icBuild.withPublish(types.InternalPublishingStrategy),
			icBuild.withControlPlaneReplicas(1),
		),
		expectedIngressController: ingresscontrollerBuild.build(
			ingresscontrollerBuild.withLoadBalancer(),
			ingresscontrollerBuild.withScope(operatorv1.InternalLoadBalancer),
			ingresscontrollerBuild.withReplicas(1),
		),
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(
				&installconfig.ClusterID{
					UUID:    "test-uuid",
					InfraID: "test-infra-id",
				},
				&installconfig.InstallConfig{Config: tc.installConfig},
				&CloudProviderConfig{},
				&AdditionalTrustBundleConfig{},
			)
			ingressAsset := &Ingress{}
			err := ingressAsset.Generate(parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}
			if tc.expectedIngressController == nil {
				assert.Len(t, ingressAsset.FileList, 1, "expected only one file to be generated")
				return
			}
			if !assert.Len(t, ingressAsset.FileList, 2, "expected two files to be generated") {
				return
			}
			assert.Equal(t, ingressAsset.FileList[1].Filename, "manifests/cluster-ingress-default-ingresscontroller.yaml")
			var actualIngressController operatorv1.IngressController
			err = yaml.Unmarshal(ingressAsset.FileList[1].Data, &actualIngressController)
			if !assert.NoError(t, err, "failed to unmarshal ingress manifest") {
				return
			}
			assert.Equal(t, tc.expectedIngressController, &actualIngressController)
		})
	}
}

func (b icBuildNamespace) withPublish(publish types.PublishingStrategy) icOption {
	return func(ic *types.InstallConfig) {
		ic.Publish = publish
	}
}

type ingresscontrollerOption func(*operatorv1.IngressController)

type ingresscontrollerBuildNamespace struct{}

var ingresscontrollerBuild ingresscontrollerBuildNamespace

func (ingresscontrollerBuildNamespace) build(opts ...ingresscontrollerOption) *operatorv1.IngressController {
	ingresscontroller := &operatorv1.IngressController{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1.GroupVersion.String(),
			Kind:       "IngressController",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-ingress-operator",
			Name:      "default",
		},
	}
	for _, opt := range opts {
		opt(ingresscontroller)
	}
	return ingresscontroller
}

func (b ingresscontrollerBuildNamespace) withLoadBalancer() ingresscontrollerOption {
	return func(ingresscontroller *operatorv1.IngressController) {
		if ingresscontroller.Spec.EndpointPublishingStrategy != nil && ingresscontroller.Spec.EndpointPublishingStrategy.LoadBalancer != nil {
			return
		}
		ingresscontroller.Spec.EndpointPublishingStrategy = &operatorv1.EndpointPublishingStrategy{
			Type:         operatorv1.LoadBalancerServiceStrategyType,
			LoadBalancer: &operatorv1.LoadBalancerStrategy{},
		}
	}
}

func (b ingresscontrollerBuildNamespace) withReplicas(replicas int) ingresscontrollerOption {
	return func(ingresscontroller *operatorv1.IngressController) {
		i := int32(replicas)
		ingresscontroller.Spec.Replicas = &i
	}
}

func (b ingresscontrollerBuildNamespace) withScope(scope operatorv1.LoadBalancerScope) ingresscontrollerOption {
	return func(ingresscontroller *operatorv1.IngressController) {
		ingresscontroller.Spec.EndpointPublishingStrategy.LoadBalancer.Scope = scope
	}
}
