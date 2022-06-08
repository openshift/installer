package mirror

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestRegistries_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name          string
		data          string
		fetchError    error
		expectedFound bool
		expectedError string
	}{
		{
			name: "valid-config-file",
			data: `
[[registry]]
prefix = "example.com/foo"
insecure = false
blocked = false
location = "internal-registry-for-example.com/bar"
mirror-by-digest-only = false

[[registry.mirror]]
location = "example-mirror-0.local/mirror-for-foo"

[[registry.mirror]]
location = "example-mirror-1.local/mirrors/foo"
insecure = true

[[registry]]
location = "registry.com"

[[registry.mirror]]
location = "mirror.registry.com`,
			expectedFound: true,
			expectedError: "",
		},
		{
			name:          "empty",
			data:          "",
			expectedFound: true,
			expectedError: "",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load mirror/registries.conf file: fetch failed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(registriesConfFilename).
				Return(
					&asset.File{
						Filename: registriesConfFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &RegistriesConf{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
