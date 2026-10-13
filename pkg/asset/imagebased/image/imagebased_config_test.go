package image

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/imagebased"
)

var (
	rawNMStateConfig = `
interfaces:
  - name: eth0
    type: ethernet
    state: up
    mac-address: 00:00:00:00:00:00
    ipv4:
      enabled: true
      address:
        - ip: 192.168.122.2
          prefix-length: 23
      dhcp: false`
)

func TestImageBasedInstallConfig_LoadedFromDisk(t *testing.T) {
	skipTestIfnmstatectlIsMissing(t)

	cases := []struct {
		name       string
		data       string
		fetchError error

		expectedError  string
		expectedFound  bool
		expectedConfig *ImageBasedInstallationConfigBuilder
	}{
		{
			name: "valid-config-single-node",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedImage: quay.io/openshift-kni/seed-image:4.16.0
seedVersion: 4.16.0
installationDisk: /dev/vda
releaseRegistry: mirror.quay.io
networkConfig:
  interfaces:
    - name: eth0
      type: ethernet
      state: up
      mac-address: 00:00:00:00:00:00
      ipv4:
        enabled: true
        address:
          - ip: 192.168.122.2
            prefix-length: 23
        dhcp: false
`,

			expectedFound:  true,
			expectedConfig: ibiConfig(),
		},
		{
			name: "not-yaml",
			data: `This is not a yaml file`,

			expectedFound: false,
			expectedError: "failed to unmarshal image-based-installation-config.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type imagebased.InstallationConfig",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},

			expectedFound: false,
		},
		{
			name:       "error-fetching-file",
			fetchError: errors.New("fetch failed"),

			expectedFound: false,
			expectedError: "failed to load image-based-installation-config.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
wrongField: wrongValue`,

			expectedFound: false,
			expectedError: "failed to unmarshal image-based-installation-config.yaml: error unmarshaling JSON: while decoding JSON: json: unknown field \"wrongField\"",
		},
		{
			name: "empty-pullSecret",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: pullSecret: Required value: you must specify a pullSecret",
		},
		{
			name: "invalid-pullSecret",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"missing\":\"auths\"}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: pullSecret: Invalid value: \"{\\\"missing\\\":\\\"auths\\\"}\": auths required",
		},
		{
			name: "invalid-sshKey",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
sshKey: invalid_ssh_key`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: sshKey: Invalid value: \"invalid_ssh_key\": ssh: no key found",
		},
		{
			name: "empty-seedVersion",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: seedVersion: Required value: you must specify a seedVersion",
		},
		{
			name: "empty-seedImage",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
installationDisk: /dev/vda`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: seedImage: Required value: you must specify a seedImage",
		},
		{
			name: "empty-installationDisk",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedImage: "quay.io/openshift-kni/seed-image:4.16.0"
seedVersion: "4.16.0"`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: installationDisk: Required value: you must specify an installationDisk",
		},
		{
			name: "invalid-additionalTrustBundle",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
additionalTrustBundle: invalid_cert
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: additionalTrustBundle: Invalid value: \"invalid_cert\": invalid block",
		},
		{
			name: "invalid-networkConfig",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
networkConfig:
  invalid: config
`,

			expectedFound: false,
			expectedError: "unknown field `invalid`",
		},
		{
			name: "empty networkConfig",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
networkConfig: null
`,
			expectedFound: true,
			expectedError: "",
		},
		{
			name: "invalid-imageDigestSources",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
imageDigestSources:
- source: quay.io
  mirrors:
  - Registry.lab.redhat.com:5000
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: imageDigestSources[0].mirrors[0]: Invalid value: \"Registry.lab.redhat.com:5000\": failed to parse: invalid reference format: repository name must be lowercase",
		},
		{
			name: "invalid-proxy-schemes",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
proxy:
  httpProxy: ""
  httpsProxy: ""
  noProxy: ""
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: proxy: Required value: must include httpProxy or httpsProxy",
		},
		{
			name: "invalid-proxy-http-uri",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
proxy:
  httpProxy: "invalidscheme://"
  httpsProxy: ""
  noProxy: ""
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: proxy.httpProxy: Unsupported value: \"invalidscheme\": supported values: \"http\"",
		},
		{
			name: "invalid-proxy-https-uri",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
proxy:
  httpProxy: ""
  httpsProxy: "invalidscheme://"
  noProxy: ""
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: proxy.httpsProxy: Unsupported value: \"invalidscheme\": supported values: \"http\", \"https\"",
		},
		{
			name: "invalid-coreosInstallerArgs arg - not allowed arg",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
coreosInstallerArgs:
- "--test"
- "3"
`,

			expectedFound: false,
			expectedError: "found unexpected flag --test for coreosInstallerArgs - allowed flags are [--append-karg --delete-karg --save-partlabel --save-partindex",
		},
		{
			name: "valid-coreosInstallerArgs arg",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
coreosInstallerArgs:
- "--save-partindex"
- "5"
`,

			expectedFound: true,
			expectedError: "",
		},
		{
			name: "valid-extraPartitionStart-zero",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
extraPartitionStart: "0"
`,

			expectedFound: true,
			expectedError: "",
		},
		{
			name: "valid-extraPartitionStart-empty-prefix",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
extraPartitionStart: "10M"
`,

			expectedFound: true,
			expectedError: "",
		},
		{
			name: "valid-extraPartitionStart-with-prefix",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
extraPartitionStart: "+10M"
`,

			expectedFound: true,
			expectedError: "",
		},
		{
			name: "invalid-extraPartitionStart-empty-suffix",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
extraPartitionStart: "-10"
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: ExtraPartitionStart: Invalid value: \"-10\": partition start must be '0' or match pattern [+-]?<number>[KMGTP]",
		},
		{
			name: "invalid-extraPartitionStart-invalid-suffix",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
extraPartitionStart: "-10L"
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: ExtraPartitionStart: Invalid value: \"-10L\": partition start must be '0' or match pattern [+-]?<number>[KMGTP]",
		},
		{
			name: "invalid-extraPartitionStart-random-string",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
extraPartitionStart: "invalid"
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: ExtraPartitionStart: Invalid value: \"invalid\": partition start must be '0' or match pattern [+-]?<number>[KMGTP]",
		},
		{
			name: "invalid-extraPartitionStart-empty-string",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
extraPartitionStart: ""
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: ExtraPartitionStart: Required value: partition start sector cannot be empty",
		},
		{
			name: "valid-architecture-amd64",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
architecture: amd64
`,

			expectedFound: true,
			expectedError: "",
		},
		{
			name: "valid-architecture-arm64",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
architecture: arm64
`,

			expectedFound: true,
			expectedError: "",
		},
		{
			name: "invalid-architecture-random-string",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-installation-config
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
seedVersion: 4.16.0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
installationDisk: /dev/vda
architecture: "invalid"
`,

			expectedFound: false,
			expectedError: "invalid Image-based Installation ISO Config: architecture: Invalid value: \"invalid\": architecture must be one of [amd64 arm64]",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(configFilename).
				Return(
					&asset.File{
						Filename: configFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ImageBasedInstallationConfig{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				if tc.expectedConfig != nil {
					assert.Equal(t, tc.expectedConfig.build(), asset.Config, "unexpected Config in ImageBasedInstallConfig")
				}
			}
		})
	}
}

// ImageBasedInstallationConfigBuilder it's a builder class to make it easier
// creating imagebased.InstallationConfig instance used in the test cases.
type ImageBasedInstallationConfigBuilder struct {
	imagebased.InstallationConfig
}

func ibiConfig() *ImageBasedInstallationConfigBuilder {
	return &ImageBasedInstallationConfigBuilder{
		InstallationConfig: imagebased.InstallationConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "image-based-installation-config",
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: imagebased.ImageBasedConfigVersion,
			},
			SeedImage:            "quay.io/openshift-kni/seed-image:4.16.0",
			SeedVersion:          "4.16.0",
			InstallationDisk:     "/dev/vda",
			ExtraPartitionStart:  "-40G",
			ExtraPartitionLabel:  defaultExtraPartitionLabel,
			ExtraPartitionNumber: 5,
			ReleaseRegistry:      "mirror.quay.io",
			Shutdown:             false,
			SSHKey:               "",
			PullSecret:           "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}",
			NetworkConfig: &aiv1beta1.NetConfig{
				Raw: unmarshalJSON([]byte(rawNMStateConfig)),
			},
		},
	}
}

func (icb *ImageBasedInstallationConfigBuilder) build() *imagebased.InstallationConfig {
	return &icb.InstallationConfig
}

func (icb *ImageBasedInstallationConfigBuilder) additionalTrustBundle(atb string) *ImageBasedInstallationConfigBuilder {
	icb.InstallationConfig.AdditionalTrustBundle = atb
	return icb
}

func (icb *ImageBasedInstallationConfigBuilder) networkConfig(nc *aiv1beta1.NetConfig) *ImageBasedInstallationConfigBuilder {
	icb.InstallationConfig.NetworkConfig = nc
	return icb
}

func (icb *ImageBasedInstallationConfigBuilder) imageDigestSources(ids []types.ImageDigestSource) *ImageBasedInstallationConfigBuilder {
	icb.InstallationConfig.ImageDigestSources = ids
	return icb
}

func (icb *ImageBasedInstallationConfigBuilder) ignitionConfigOverride(ignitionConfigOverride string) *ImageBasedInstallationConfigBuilder {
	icb.InstallationConfig.IgnitionConfigOverride = ignitionConfigOverride
	return icb
}

func (icb *ImageBasedInstallationConfigBuilder) architecture(architecture string) *ImageBasedInstallationConfigBuilder {
	icb.InstallationConfig.Architecture = architecture
	return icb
}

func unmarshalJSON(b []byte) []byte {
	output, err := yaml.JSONToYAML(b)
	if err != nil {
		return nil
	}
	return output
}

func skipTestIfnmstatectlIsMissing(t *testing.T) {
	t.Helper()

	_, execErr := exec.LookPath("nmstatectl")
	if execErr != nil {
		t.Skip("No nmstatectl binary available")
	}
}
