package configimage

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/imagebased"
)

const (
	testSSHKey = "ssh-rsa AAAAB3NzaC1y1LJe3zew1ghc= root@localhost.localdomain"

	testSecret = `{"auths":{"cloud.openshift.com":{"auth":"b3BlUTA=","email":"test@redhat.com"}}}` //nolint:gosec // not real credentials
)

func TestClusterConfiguration_Generate(t *testing.T) {
	cases := []struct {
		name         string
		dependencies []asset.Asset

		expectedError  string
		expectedConfig *imagebased.SeedReconfiguration
	}{
		{
			name: "missing install config",
			dependencies: []asset.Asset{
				clusterID(),
				&InstallConfig{},
				kubeadminPassword(),
				lbCertKey(),
				localhostCertKey(),
				serviceNetworkCertKey(),
				adminKubeConfigCertKey(),
				ingressCertKey(),
				imageBasedConfig(),
			},
			expectedError: "missing configuration or manifest file",
		},
		{
			name: "valid configuration",
			dependencies: []asset.Asset{
				clusterID(),
				kubeadminPassword(),
				lbCertKey(),
				localhostCertKey(),
				serviceNetworkCertKey(),
				adminKubeConfigCertKey(),
				ingressCertKey(),
				installConfig().build(),
				imageBasedConfig(),
			},

			expectedConfig: clusterConfiguration().build().Config,
		},
		{
			name: "valid configuration with proxy",
			dependencies: []asset.Asset{
				clusterID(),
				kubeadminPassword(),
				lbCertKey(),
				localhostCertKey(),
				serviceNetworkCertKey(),
				adminKubeConfigCertKey(),
				ingressCertKey(),
				installConfig().proxy(proxy()).build(),
				imageBasedConfig(),
			},

			expectedConfig: clusterConfiguration().proxy(proxy()).build().Config,
		},
		{
			name: "valid configuration with additionalTrustBundle, policyAlways without proxy",
			dependencies: []asset.Asset{
				clusterID(),
				kubeadminPassword(),
				lbCertKey(),
				localhostCertKey(),
				serviceNetworkCertKey(),
				adminKubeConfigCertKey(),
				ingressCertKey(),
				installConfig().
					additionalTrustBundle(testCert).
					additionalTrustBundlePolicy(types.PolicyAlways).build(),
				imageBasedConfig(),
			},

			expectedConfig: clusterConfiguration().additionalTrustBundle(imagebased.AdditionalTrustBundle{
				UserCaBundle:         testCert,
				ProxyConfigmapBundle: testCert,
				ProxyConfigmapName:   "user-ca-bundle",
			}).build().Config,
		},
		{
			name: "valid configuration with additionalTrustBundle, policyProxyOnly without proxy",
			dependencies: []asset.Asset{
				clusterID(),
				kubeadminPassword(),
				lbCertKey(),
				localhostCertKey(),
				serviceNetworkCertKey(),
				adminKubeConfigCertKey(),
				ingressCertKey(),
				installConfig().proxy(proxy()).additionalTrustBundle(testCert).build(),
				imageBasedConfig(),
			},

			expectedConfig: clusterConfiguration().
				proxy(proxy()).
				additionalTrustBundle(imagebased.AdditionalTrustBundle{
					UserCaBundle: testCert,
				}).
				build().Config,
		},
		{
			name: "valid configuration with additionalTrustBundle with proxy",
			dependencies: []asset.Asset{
				clusterID(),
				kubeadminPassword(),
				lbCertKey(),
				localhostCertKey(),
				serviceNetworkCertKey(),
				adminKubeConfigCertKey(),
				ingressCertKey(),
				installConfig().
					proxy(proxy()).
					additionalTrustBundle(testCert).
					additionalTrustBundlePolicy(types.PolicyProxyOnly).build(),
				imageBasedConfig(),
			},

			expectedConfig: clusterConfiguration().
				proxy(proxy()).
				additionalTrustBundle(imagebased.AdditionalTrustBundle{
					UserCaBundle:         testCert,
					ProxyConfigmapBundle: testCert,
					ProxyConfigmapName:   "user-ca-bundle",
				}).
				build().Config,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &ClusterConfiguration{}
			err := asset.Generate(context.TODO(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-configuration/manifest.json", configFile.Filename)

				var actualConfig imagebased.SeedReconfiguration
				err = yaml.Unmarshal(configFile.Data, &actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedConfig, actualConfig)
			}
		})
	}
}

func TestClusterConfiguration_LoadedFromDisk(t *testing.T) {
	cases := []struct {
		name       string
		data       string
		fetchError error

		expectedFound  bool
		expectedError  string
		expectedConfig *imagebased.SeedReconfiguration
	}{
		{
			name: "valid-config-file",
			data: `
{
  "api_version": 1,
  "base_domain": "testing.com",
  "cluster_name": "ocp-ibi",
  "cluster_id": "6065edc6-939c-4dc3-81c7-1c20d840d064",
  "infra_id": "an-infra-id",
  "release_registry": "quay.io",
  "hostname": "somehostname",
  "KubeconfigCryptoRetention": {
    "KubeAPICrypto": {
      "ServingCrypto": {
        "localhost_signer_private_key": "localhost-key",
        "service_network_signer_private_key": "service-network-key",
        "loadbalancer_external_signer_private_key": "lb-key"
      },
      "ClientAuthCrypto": {
        "admin_ca_certificate": "admin-kubeconfig-cert"
      }
    },
    "IngresssCrypto": {
      "ingress_ca": "ingress-key"
    }
  },
  "ssh_key": "ssh-rsa AAAAB3NzaC1y1LJe3zew1ghc= root@localhost.localdomain",
  "kubeadmin_password_hash": "a-password-hash",
  "raw_nm_state_config": "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.2\n      prefix-length: 23\n    dhcp: false\n    enabled: true\n  mac-address: \"00:00:00:00:00:00\"\n  name: eth0\n  state: up\n  type: ethernet\n",
  "pull_secret": "{\"auths\":{\"cloud.openshift.com\":{\"auth\":\"b3BlUTA=\",\"email\":\"test@redhat.com\"}}}",
  "machine_network": "10.10.11.0/24"
}`,

			expectedFound:  true,
			expectedConfig: clusterConfiguration().Config,
		},
		{
			name: "not-json",
			data: `This is not a JSON file`,

			expectedError: "failed to unmarshal cluster-configuration/manifest.json: invalid JSON syntax",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:       "error-fetching-file",
			fetchError: errors.New("fetch failed"),

			expectedError: "failed to load cluster-configuration/manifest.json file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `{"some-unknown-field":"withsomevalue"}`,

			expectedError: "failed to unmarshal cluster-configuration/manifest.json: unknown field \"some-unknown-field\"",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(clusterConfigurationFilename).
				Return(
					&asset.File{
						Filename: clusterConfigurationFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ClusterConfiguration{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.Equal(t, nil, err)
			}

			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in ClusterConfiguration")
			}
		})
	}
}

type ClusterConfigurationBuilder struct {
	ClusterConfiguration
}

func (ccb *ClusterConfigurationBuilder) build() *ClusterConfiguration {
	return &ccb.ClusterConfiguration
}

func (ccb *ClusterConfigurationBuilder) proxy(proxy *types.Proxy) *ClusterConfigurationBuilder {
	ccb.ClusterConfiguration.Config.Proxy = proxy
	return ccb
}

func (ccb *ClusterConfigurationBuilder) additionalTrustBundle(additionalTrustBundle imagebased.AdditionalTrustBundle) *ClusterConfigurationBuilder {
	ccb.ClusterConfiguration.Config.AdditionalTrustBundle = additionalTrustBundle
	return ccb
}

func clusterConfiguration() *ClusterConfigurationBuilder {
	ccb := &ClusterConfigurationBuilder{}

	cc := &ClusterConfiguration{}

	clusterID := clusterID()
	installConfig := installConfig().build()
	imageBasedConfig := imageBasedConfig()

	cc.Config = &imagebased.SeedReconfiguration{
		APIVersion:            imagebased.SeedReconfigurationVersion,
		BaseDomain:            installConfig.Config.BaseDomain,
		ClusterID:             clusterID.UUID,
		ClusterName:           installConfig.ClusterName(),
		Hostname:              imageBasedConfig.Config.Hostname,
		InfraID:               clusterID.InfraID,
		KubeadminPasswordHash: string(kubeadminPassword().PasswordHash),
		PullSecret:            installConfig.Config.PullSecret,
		RawNMStateConfig:      imageBasedConfig.Config.NetworkConfig.String(),
		ReleaseRegistry:       imageBasedConfig.Config.ReleaseRegistry,
		SSHKey:                installConfig.Config.SSHKey,
	}

	cc.Config.KubeconfigCryptoRetention = imagebased.KubeConfigCryptoRetention{
		KubeAPICrypto: imagebased.KubeAPICrypto{
			ServingCrypto: imagebased.ServingCrypto{
				LoadbalancerSignerPrivateKey:   string(lbCertKey().Key()),
				LocalhostSignerPrivateKey:      string(localhostCertKey().Key()),
				ServiceNetworkSignerPrivateKey: string(serviceNetworkCertKey().Key()),
			},
			ClientAuthCrypto: imagebased.ClientAuthCrypto{
				AdminCACertificate: string(adminKubeConfigCertKey().Cert()),
			},
		},
		IngresssCrypto: imagebased.IngresssCrypto{
			IngressCA: string(ingressCertKey().SelfSignedCertKey.Key()),
		},
	}

	cc.Config.MachineNetwork = installConfig.Config.Networking.MachineNetwork[0].CIDR.String()

	ccb.ClusterConfiguration = *cc
	return ccb
}

func clusterID() *ClusterID {
	return &ClusterID{
		installconfig.ClusterID{
			UUID:    "6065edc6-939c-4dc3-81c7-1c20d840d064",
			InfraID: "an-infra-id",
		},
	}
}

func proxy() *types.Proxy {
	return &types.Proxy{
		HTTPProxy:  "http://10.10.10.11:80",
		HTTPSProxy: "http://my-lab-proxy.org:443",
		NoProxy:    "internal.com",
	}
}

func kubeadminPassword() *password.KubeadminPassword {
	return &password.KubeadminPassword{
		PasswordHash: []byte("a-password-hash"),
	}
}

func lbCertKey() *tls.KubeAPIServerLBSignerCertKey {
	return &tls.KubeAPIServerLBSignerCertKey{
		SelfSignedCertKey: tls.SelfSignedCertKey{
			CertKey: tls.CertKey{
				CertRaw: []byte("lb-cert"),
				KeyRaw:  []byte("lb-key"),
			},
		},
	}
}

func localhostCertKey() *tls.KubeAPIServerLocalhostSignerCertKey {
	return &tls.KubeAPIServerLocalhostSignerCertKey{
		SelfSignedCertKey: tls.SelfSignedCertKey{
			CertKey: tls.CertKey{
				CertRaw: []byte("localhost-cert"),
				KeyRaw:  []byte("localhost-key"),
			},
		},
	}
}

func serviceNetworkCertKey() *tls.KubeAPIServerServiceNetworkSignerCertKey {
	return &tls.KubeAPIServerServiceNetworkSignerCertKey{
		SelfSignedCertKey: tls.SelfSignedCertKey{
			CertKey: tls.CertKey{
				CertRaw: []byte("service-network-cert"),
				KeyRaw:  []byte("service-network-key"),
			},
		},
	}
}

func adminKubeConfigCertKey() *tls.AdminKubeConfigSignerCertKey {
	return &tls.AdminKubeConfigSignerCertKey{
		SelfSignedCertKey: tls.SelfSignedCertKey{
			CertKey: tls.CertKey{
				CertRaw: []byte("admin-kubeconfig-cert"),
				KeyRaw:  []byte("admin-kubeconfig-key"),
			},
		},
	}
}

func ingressCertKey() *IngressOperatorSignerCertKey {
	return &IngressOperatorSignerCertKey{
		SelfSignedCertKey: tls.SelfSignedCertKey{
			CertKey: tls.CertKey{
				CertRaw: []byte("ingress-cert"),
				KeyRaw:  []byte("ingress-key"),
			},
		},
	}
}

type InstallConfigBuilder struct {
	InstallConfig
}

func (icb *InstallConfigBuilder) build() *InstallConfig {
	return &icb.InstallConfig
}

func (icb *InstallConfigBuilder) proxy(proxy *types.Proxy) *InstallConfigBuilder {
	icb.InstallConfig.Config.Proxy = proxy
	return icb
}

func (icb *InstallConfigBuilder) additionalTrustBundle(additionalTrustBundle string) *InstallConfigBuilder {
	icb.InstallConfig.Config.AdditionalTrustBundle = additionalTrustBundle
	return icb
}

func (icb *InstallConfigBuilder) additionalTrustBundlePolicy(policy types.PolicyType) *InstallConfigBuilder {
	icb.InstallConfig.Config.AdditionalTrustBundlePolicy = policy
	return icb
}

func installConfig() *InstallConfigBuilder {
	return &InstallConfigBuilder{
		InstallConfig: *defaultInstallConfig(),
	}
}
