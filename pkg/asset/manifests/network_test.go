package manifests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

const stubDefaultNetworkConfigOVNWithMTU = `apiVersion: operator.openshift.io/v1
kind: Network
metadata:
  creationTimestamp: null
  name: cluster
spec:
  clusterNetwork: null
  defaultNetwork:
    ovnKubernetesConfig:
      egressIPConfig: {}
      mtu: 1000
    type: OVNKubernetes
  disableNetworkDiagnostics: false
  managementState: Managed
  observedConfig: null
  serviceNetwork: null
  unsupportedConfigOverrides: null
status:
  readyReplicas: 0
`

func validInstallConfigWithMTU(ic *types.InstallConfig) *installconfig.InstallConfig {
	return &installconfig.InstallConfig{
		AssetBase: installconfig.AssetBase{
			Config: ic,
		},
	}
}

func TestNetworking_GenerateCustomNetworkConfigMTU(t *testing.T) {
	type args struct {
		ic *installconfig.InstallConfig
	}
	tests := []struct {
		name       string
		no         *Networking
		args       args
		want       *operatorv1.DefaultNetworkDefinition
		wantErr    bool
		wantErrExp string
	}{
		{
			name:    "mtu unset OVNKubernetes",
			no:      &Networking{},
			args:    args{ic: &installconfig.InstallConfig{}},
			wantErr: false,
		}, {
			name: "mtu unset OVNKubernetes edge defaults",
			no:   &Networking{},
			args: args{ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
				},
				Compute: []types.MachinePool{
					{
						Name: "edge",
					},
				},
				Platform: types.Platform{AWS: &awstypes.Platform{}},
			})},
			want: &operatorv1.DefaultNetworkDefinition{
				Type: operatorv1.NetworkTypeOVNKubernetes,
				OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
					MTU: ptr.To(uint32(1200)),
				},
			},
		}, {
			name: "mtu set OVNKubernetes",
			no:   &Networking{},
			args: args{ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetworkMTU: 1500,
					NetworkType:       "OVNKubernetes",
				},
				Platform: types.Platform{AWS: &awstypes.Platform{}},
			})},
			want: &operatorv1.DefaultNetworkDefinition{
				Type: operatorv1.NetworkTypeOVNKubernetes,
				OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
					MTU: ptr.To(uint32(1500)),
				},
			},
		}, {
			name: "mtu set OVNKubernetes edge",
			no:   &Networking{},
			args: args{ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetworkMTU: 1500,
					NetworkType:       "OVNKubernetes",
				},
				Compute: []types.MachinePool{
					{
						Name: "edge",
					},
				},
				Platform: types.Platform{AWS: &awstypes.Platform{}},
			})},
			want: &operatorv1.DefaultNetworkDefinition{
				Type: operatorv1.NetworkTypeOVNKubernetes,
				OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
					MTU: ptr.To(uint32(1500)),
				},
			},
		}, {
			name: "mtu set OVNKubernetes edge high",
			no:   &Networking{},
			args: args{ic: validInstallConfigWithMTU(&types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetworkMTU: 8000,
					NetworkType:       "OVNKubernetes",
				},
				Compute: []types.MachinePool{
					{
						Name: "edge",
					},
				},
				Platform: types.Platform{AWS: &awstypes.Platform{}},
			})},
			want: &operatorv1.DefaultNetworkDefinition{
				Type: operatorv1.NetworkTypeOVNKubernetes,
				OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
					MTU: ptr.To(uint32(8000)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &Networking{
				Config:   tt.no.Config,
				FileList: tt.no.FileList,
			}
			got, err := no.GenerateCustomNetworkConfigMTU(tt.args.ic)
			if (err != nil) != tt.wantErr {
				if tt.wantErrExp != "" {
					assert.Regexp(t, tt.wantErrExp, err)
					return
				}
				t.Errorf("Networking.GenerateCustomNetworkConfigMTU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				dumpDN := ""
				if got.OVNKubernetesConfig != nil {
					t.Errorf("Networking.GenerateCustomNetworkConfigMTU() = %v (%v), want %v", got, *got.OVNKubernetesConfig, tt.want)
					dumpDN = fmt.Sprintf("%v", *got.OVNKubernetesConfig)
					if got.OVNKubernetesConfig.MTU != nil {
						dumpDN += fmt.Sprintf(" mtu[%v]", *got.OVNKubernetesConfig.MTU)
					}
				} else if got.OpenShiftSDNConfig != nil {
					t.Errorf("Networking.GenerateCustomNetworkConfigMTU() = %v (%v), want %v", got, *got.OpenShiftSDNConfig, tt.want)
					dumpDN = fmt.Sprintf("%v", *got.OpenShiftSDNConfig)
					if got.OpenShiftSDNConfig.MTU != nil {
						dumpDN += fmt.Sprintf(" mtu[%v]", *got.OpenShiftSDNConfig.MTU)
					}
				}

				t.Errorf("Networking.GenerateCustomNetworkConfigMTU() = %v (debug: %s), want %v", got, dumpDN, tt.want)
			}
		})
	}
}

func TestNetworking_generateCustomCnoConfig(t *testing.T) {
	type fields struct {
		Config   *configv1.Network
		FileList []*asset.File
	}
	type args struct {
		defaultNetwork *operatorv1.DefaultNetworkDefinition
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []byte
		wantErr    bool
		wantErrExp string
	}{
		{
			name:   "custom OVNKubernetes",
			fields: fields{},
			args: args{defaultNetwork: &operatorv1.DefaultNetworkDefinition{
				Type: operatorv1.NetworkTypeOVNKubernetes,
				OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
					MTU: ptr.To(uint32(1000)),
				},
			}},
			want: []byte(stubDefaultNetworkConfigOVNWithMTU),
		}, {
			name:       "invalid config",
			fields:     fields{},
			args:       args{defaultNetwork: nil},
			wantErr:    true,
			wantErrExp: `defaultNetwork must be specified`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &Networking{
				Config:   tt.fields.Config,
				FileList: tt.fields.FileList,
			}
			got, err := no.generateCustomCnoConfig(tt.args.defaultNetwork)
			if (err != nil) != tt.wantErr {
				if tt.wantErrExp != "" {
					assert.Regexp(t, tt.wantErrExp, err)
					return
				}
				t.Errorf("Networking.generateCustomCnoConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Networking.generateCustomCnoConfig() got =>\n[%v],\n want =>\n[%v]", string(got), string(tt.want))
			}
		})
	}
}
