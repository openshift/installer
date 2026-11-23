package powervs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOSImageURL(t *testing.T) {
	cases := []struct {
		name           string
		url            string
		expectedBucket string
		expectedObject string
		expectedRegion string
		expectedErr    bool
	}{
		{
			name:           "valid",
			url:            "https://s3.us-south.cloud-object-storage.appdomain.cloud/rhcos-powervs-images-us-south/rhcos-10-2-20260423-0-ppc64le-powervs.ova.gz",
			expectedBucket: "rhcos-powervs-images-us-south",
			expectedObject: "rhcos-10-2-20260423-0-ppc64le-powervs.ova.gz",
			expectedRegion: "us-south",
		},
		{
			name:           "object with folder",
			url:            "https://s3.eu-de.cloud-object-storage.appdomain.cloud/bucket/folder/image.ova.gz",
			expectedBucket: "bucket",
			expectedObject: "folder/image.ova.gz",
			expectedRegion: "eu-de",
		},
		{
			name:        "not https",
			url:         "http://s3.us-south.cloud-object-storage.appdomain.cloud/bucket/image.ova.gz",
			expectedErr: true,
		},
		{
			name:        "wrong host",
			url:         "https://example.com/bucket/image.ova.gz",
			expectedErr: true,
		},
		{
			name:        "missing region",
			url:         "https://s3..cloud-object-storage.appdomain.cloud/bucket/image.ova.gz",
			expectedErr: true,
		},
		{
			name:        "missing object",
			url:         "https://s3.us-south.cloud-object-storage.appdomain.cloud/bucket",
			expectedErr: true,
		},
		{
			name:        "not a url",
			url:         "my-workspace-image",
			expectedErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bucket, object, region, err := ParseOSImageURL(tc.url)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBucket, bucket)
			assert.Equal(t, tc.expectedObject, object)
			assert.Equal(t, tc.expectedRegion, region)
		})
	}
}
