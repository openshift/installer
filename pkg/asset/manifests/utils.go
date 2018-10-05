package manifests

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Converts a platform to the cloudProvider that k8s understands
func tectonicCloudProvider(platform types.Platform) string {
	if platform.AWS != nil {
		return "aws"
	}
	if platform.Libvirt != nil {
		return "libvirt"
	}
	return ""
}

func getAPIServerURL(ic *types.InstallConfig) string {
	return fmt.Sprintf("https://%s-api.%s:6443", ic.ObjectMeta.Name, ic.BaseDomain)
}
