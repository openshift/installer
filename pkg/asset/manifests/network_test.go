package manifests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/powervs"
)

func validInstallConfigWithMTU(ic *types.InstallConfig) *installconfig.InstallConfig {
	return &installconfig.InstallConfig{
		AssetBase: installconfig.AssetBase{
			Config: ic,
		},
	}
}

func stubDefaultNetworkConfigOVN() *operatorv1.Network {
	return &operatorv1.Network{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "operator.openshift.io/v1",
			Kind:       "Network",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: operatorv1.NetworkSpec{
			ClusterNetwork: []operatorv1.ClusterNetworkEntry{},
			ServiceNetwork: []string{},
			OperatorSpec:   operatorv1.OperatorSpec{ManagementState: operatorv1.Managed},
			DefaultNetwork: operatorv1.DefaultNetworkDefinition{
				Type:                operatorv1.NetworkTypeOVNKubernetes,
				OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{},
			},
		},
		Status: operatorv1.NetworkStatus{},
	}
}

func TestNetworking_GenerateCustomNetworkConfigMTU(t *testing.T) {
	tests := []struct {
		name       string
		no         *Networking
		ic         *installconfig.InstallConfig
		want       *operatorv1.Network
		wantErr    bool
		wantErrExp string
	}{
		{
			name: "mtu unset OVNKubernetes edge defaults",
			no:   &Networking{},
			ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
				},
				Compute: []types.MachinePool{
					{
						Name: "edge",
					},
				},
				Platform: types.Platform{AWS: &aws.Platform{}},
			}),
			want: func() *operatorv1.Network {
				ovn := stubDefaultNetworkConfigOVN()
				ovn.Spec.DefaultNetwork.OVNKubernetesConfig.MTU = ptr.To(uint32(1200))
				return ovn
			}(),
		}, {
			name: "mtu set OVNKubernetes",
			no:   &Networking{},
			ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetworkMTU: 1500,
					NetworkType:       "OVNKubernetes",
				},
				Platform: types.Platform{AWS: &aws.Platform{}},
			}),
			want: func() *operatorv1.Network {
				ovn := stubDefaultNetworkConfigOVN()
				ovn.Spec.DefaultNetwork.OVNKubernetesConfig.MTU = ptr.To(uint32(1500))
				return ovn
			}(),
		}, {
			name: "mtu set OVNKubernetes edge",
			no:   &Networking{},
			ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetworkMTU: 1500,
					NetworkType:       "OVNKubernetes",
				},
				Compute: []types.MachinePool{
					{
						Name: "edge",
					},
				},
				Platform: types.Platform{AWS: &aws.Platform{}},
			}),
			want: func() *operatorv1.Network {
				ovn := stubDefaultNetworkConfigOVN()
				ovn.Spec.DefaultNetwork.OVNKubernetesConfig.MTU = ptr.To(uint32(1500))
				return ovn
			}(),
		}, {
			name: "mtu set OVNKubernetes edge high",
			no:   &Networking{},
			ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetworkMTU: 8000,
					NetworkType:       "OVNKubernetes",
				},
				Compute: []types.MachinePool{
					{
						Name: "edge",
					},
				},
				Platform: types.Platform{AWS: &aws.Platform{}},
			}),
			want: func() *operatorv1.Network {
				ovn := stubDefaultNetworkConfigOVN()
				ovn.Spec.DefaultNetwork.OVNKubernetesConfig.MTU = ptr.To(uint32(8000))
				return ovn
			}(),
		},
		{
			name: "PowerVS has host routing enabled",
			no:   &Networking{},
			ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
				},
				Platform: types.Platform{PowerVS: &powervs.Platform{}},
			}),
			want: func() *operatorv1.Network {
				ovn := stubDefaultNetworkConfigOVN()
				ovn.Spec.DefaultNetwork.OVNKubernetesConfig.GatewayConfig = &operatorv1.GatewayConfig{RoutingViaHost: true}
				return ovn
			}(),
		},
		{
			name: "Custom OVN V4 InternalJoinSubnet",
			no:   &Networking{},
			ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
					OVNKubernetesConfig: &types.OVNKubernetesConfig{
						IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("101.64.0.0/16")},
					},
				},
				Platform: types.Platform{AWS: &aws.Platform{}},
			}),
			want: func() *operatorv1.Network {
				ovn := stubDefaultNetworkConfigOVN()
				ovn.Spec.DefaultNetwork.OVNKubernetesConfig.IPv4 = &operatorv1.IPv4OVNKubernetesConfig{InternalJoinSubnet: "101.64.0.0/16"}
				return ovn
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := clusterNetworkOperatorConfig(tt.ic, []configv1.ClusterNetworkEntry{}, []string{})
			if (err != nil) != tt.wantErr {
				if tt.wantErrExp != "" {
					assert.Regexp(t, tt.wantErrExp, err)
					return
				}
				t.Errorf("Networking.GenerateCustomNetworkConfigMTU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
