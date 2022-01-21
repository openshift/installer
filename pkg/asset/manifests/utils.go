package manifests

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
)

type configurationObject struct {
	metav1.TypeMeta

	Metadata metadata    `json:"metadata,omitempty"`
	Data     genericData `json:"data,omitempty"`
}

type metadata struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func configMap(namespace, name string, data genericData) *configurationObject {
	return &configurationObject{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Metadata: metadata{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

func getAPIServerURL(ic *types.InstallConfig) string {
	return fmt.Sprintf("https://api.%s:6443", ic.ClusterDomain())
}

func getInternalAPIServerURL(ic *types.InstallConfig) string {
	return fmt.Sprintf("https://api-int.%s:6443", ic.ClusterDomain())
}

func getControlPlaneTopology(ic *types.InstallConfig) configv1.TopologyMode {
	if ic.ControlPlane.Replicas != nil && *ic.ControlPlane.Replicas < 3 {
		return configv1.SingleReplicaTopologyMode
	}
	return configv1.HighlyAvailableTopologyMode
}

func getInfrastructureTopology(ic *types.InstallConfig) configv1.TopologyMode {
	numOfWorkers := int64(0)
	for _, mp := range ic.Compute {
		if mp.Replicas != nil {
			numOfWorkers += *mp.Replicas
		}
	}
	switch numOfWorkers {
	case 0:
		return getControlPlaneTopology(ic)
	case 1:
		return configv1.SingleReplicaTopologyMode
	default:
		return configv1.HighlyAvailableTopologyMode
	}
}
