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

func TestCaBundle_LoadedFromDisk(t *testing.T) {

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
-----BEGIN CERTIFICATE-----
MIIDvTCCAqWgAwIBAgICEAA....
----- END CERTIFICATE -----`,
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
			expectedError: "failed to load mirror/ca-bundle.crt file: fetch failed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(caBundleFilename).
				Return(
					&asset.File{
						Filename: caBundleFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &CaBundle{}
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
