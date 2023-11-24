package manifests

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
)

// OvnKubeConfig creates a config file for the OVNKubernetes CNI provider
func OvnKubeConfig(cns []configv1.ClusterNetworkEntry, sn []string, useHostRouting bool) ([]byte, error) {

	operCNs := []operatorv1.ClusterNetworkEntry{}
	for _, cn := range cns {
		ocn := operatorv1.ClusterNetworkEntry{
			CIDR:       cn.CIDR,
			HostPrefix: cn.HostPrefix,
		}
		operCNs = append(operCNs, ocn)
	}
	ovnConfig := operatorv1.Network{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "operator.openshift.io/v1",
			Kind:       "Network",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: operatorv1.NetworkSpec{
			OperatorSpec:   operatorv1.OperatorSpec{ManagementState: operatorv1.Managed},
			ClusterNetwork: operCNs,
			ServiceNetwork: sn,
			DefaultNetwork: operatorv1.DefaultNetworkDefinition{
				Type: operatorv1.NetworkTypeOVNKubernetes,
				OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
					GatewayConfig: &operatorv1.GatewayConfig{
						RoutingViaHost: useHostRouting,
						// (chocobomb) This is only for test, to see if the Network CRD contains this field.
						// If yes, this will succeed. If not, it will crash with "field not declared in schema".
						// Used to debug https://github.com/openshift/installer/pull/7748 together with
						// https://github.com/openshift/cluster-network-operator/pull/2118.
						IPv4: operatorv1.IPv4GatewayConfig{
							InternalMasqueradeSubnet: "169.254.169.0/29",
						},
					},
				},
			},
		},
		Status: operatorv1.NetworkStatus{},
	}

	return yaml.Marshal(ovnConfig)
}
