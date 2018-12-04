package installconfig

import (
	"errors"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/golang/mock/gomock"
	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		BaseDomain: "test-domain",
		Networking: types.Networking{
			Type:        "OpenshiftSDN",
			MachineCIDR: *defaultMachineCIDR,
			ServiceCIDR: *defaultServiceCIDR,
			ClusterNetworks: []netopv1.ClusterNetwork{
				{
					CIDR:             "192.168.1.0/24",
					HostSubnetLength: 4,
				},
			},
		},
		Machines: []types.MachinePool{
			{
				Name: "master",
			},
			{
				Name: "worker",
			},
		},
		Platform: types.Platform{
			AWS: &aws.Platform{
				Region: "us-east-1",
			},
		},
		PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
	}
}

func TestInstallConfigLoad(t *testing.T) {
	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  bool
		expectedConfig *types.InstallConfig
	}{
		{
			name: "valid InstallConfig",
			data: func() string {
				ic := validInstallConfig()
				data, _ := yaml.Marshal(ic)
				return string(data)
			}(),
			expectedFound:  true,
			expectedConfig: validInstallConfig(),
		},
		{
			name: "invalid InstallConfig",
			data: `
metadata:
  name: test-cluster
`,
			expectedError: true,
		},
		{
			name:          "empty",
			data:          "",
			expectedError: true,
		},
		{
			name:          "not yaml",
			data:          "This is not yaml.",
			expectedError: true,
		},
		{
			name:       "file not found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error fetching file",
			fetchError:    errors.New("fetch failed"),
			expectedError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(installConfigFilename).
				Return(
					&asset.File{
						Filename: "test-filename",
						Data:     []byte(tc.data)},
					tc.fetchError,
				)
			fileFetcher.EXPECT().FetchByName(deprecatedInstallConfigFilename).
				Return(nil, &os.PathError{Err: os.ErrNotExist}).
				AnyTimes()

			ic := &InstallConfig{}
			found, err := ic.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError {
				assert.Error(t, err, "expected error from Load")
			} else {
				assert.NoError(t, err, "unexpected error from Load")
			}
			if tc.expectedFound {
				assert.Equal(t, "test-filename", ic.File.Filename, "unexpected File.Filename in InstallConfig")
				assert.Equal(t, tc.data, string(ic.File.Data), "unexpected File.Data in InstallConfig")
			} else {
				assert.Nil(t, ic.File, "expected File in InstallConfig to be nil")
			}
			assert.Equal(t, tc.expectedConfig, ic.Config, "unexpected Config in InstallConfig")
		})
	}
}

func TestFetchInstallConfigFile(t *testing.T) {
	cases := []struct {
		name                         string
		installConfigFile            *asset.File
		installConfigError           error
		deprecatedInstallConfigFile  *asset.File
		deprecatedInstallConfigError error
		expectedFile                 *asset.File
		expectedError                error
	}{
		{
			name: "no files",
		},
		{
			name: "only new file",
			installConfigFile: &asset.File{
				Filename: installConfigFilename,
				Data:     []byte("test-data"),
			},
			expectedFile: &asset.File{
				Filename: installConfigFilename,
				Data:     []byte("test-data"),
			},
		},
		{
			name:               "only deprecated file",
			installConfigError: &os.PathError{Err: os.ErrNotExist},
			deprecatedInstallConfigFile: &asset.File{
				Filename: deprecatedInstallConfigFilename,
				Data:     []byte("test-data"),
			},
			expectedFile: &asset.File{
				Filename: installConfigFilename,
				Data:     []byte("test-data"),
			},
		},
		{
			name: "both files",
			installConfigFile: &asset.File{
				Filename: installConfigFilename,
				Data:     []byte("test-data"),
			},
			deprecatedInstallConfigFile: &asset.File{
				Filename: deprecatedInstallConfigFilename,
				Data:     []byte("deprecated-test-data"),
			},
			expectedFile: &asset.File{
				Filename: installConfigFilename,
				Data:     []byte("test-data"),
			},
		},
		{
			name:               "error reading new file",
			installConfigError: errors.New("test-error"),
			deprecatedInstallConfigFile: &asset.File{
				Filename: deprecatedInstallConfigFilename,
				Data:     []byte("deprecated-test-data"),
			},
			expectedError: errors.New("test-error"),
		},
		{
			name:                         "error reading deprecated file",
			installConfigError:           &os.PathError{Err: os.ErrNotExist},
			deprecatedInstallConfigError: errors.New("test-error"),
			expectedError:                errors.New("test-error"),
		},
		{
			name:                         "error reading both files",
			installConfigError:           errors.New("test-error-new"),
			deprecatedInstallConfigError: errors.New("test-error-deprecated"),
			expectedError:                errors.New("test-error-new"),
		},
		{
			name: "new file with error reading deprecated file",
			installConfigFile: &asset.File{
				Filename: installConfigFilename,
				Data:     []byte("test-data"),
			},
			deprecatedInstallConfigError: errors.New("test-error"),
			expectedFile: &asset.File{
				Filename: installConfigFilename,
				Data:     []byte("test-data"),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(installConfigFilename).
				Return(tc.installConfigFile, tc.installConfigError)
			fileFetcher.EXPECT().FetchByName(deprecatedInstallConfigFilename).
				Return(tc.deprecatedInstallConfigFile, tc.deprecatedInstallConfigError).
				AnyTimes()

			file, err := fetchInstallConfigFile(fileFetcher)
			assert.Equal(t, tc.expectedFile, file, "unexpected file")
			assert.Equal(t, tc.expectedError, err, "unexpected error")
		})
	}
}
