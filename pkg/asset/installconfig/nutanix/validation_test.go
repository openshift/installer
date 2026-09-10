package nutanix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findImagePrefix(t *testing.T) {
	tests := []struct {
		name           string
		imgSource      string
		expectedPrefix string
		expectedIndex  int
	}{
		{
			name:           "RHCOS image URL",
			imgSource:      "https://rhcos.mirror.openshift.com/art/storage/prod/streams/rhel-10.2/builds/10.2.20260715-0/x86_64/rhcos-10.2.20260715-0-nutanix.x86_64.qcow2",
			expectedPrefix: "/rhcos-",
			expectedIndex:  91,
		},
		{
			name:           "SCOS image URL",
			imgSource:      "https://rhcos.mirror.openshift.com/art/storage/prod/streams/c10s/builds/10.0.20251103-0/x86_64/scos-10.0.20251103-0-nutanix.x86_64.qcow2",
			expectedPrefix: "/scos-",
			expectedIndex:  86,
		},
		{
			name:           "unknown image prefix",
			imgSource:      "https://example.com/some-other-image-10.0.20251103-0-nutanix.x86_64.qcow2",
			expectedPrefix: "",
			expectedIndex:  -1,
		},
		{
			name:           "empty source URI",
			imgSource:      "",
			expectedPrefix: "",
			expectedIndex:  -1,
		},
		{
			name:           "RHCOS image from RHEL 9 stream",
			imgSource:      "https://rhcos.mirror.openshift.com/art/storage/prod/streams/rhel-9.8/builds/9.8.20260715-1/x86_64/rhcos-9.8.20260715-1-nutanix.x86_64.qcow2",
			expectedPrefix: "/rhcos-",
			expectedIndex:  89,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, index := findImagePrefix(tt.imgSource)
			assert.Equal(t, tt.expectedPrefix, prefix)
			if tt.expectedIndex == -1 {
				assert.Equal(t, -1, index)
			} else {
				assert.GreaterOrEqual(t, index, 0, "expected a valid index")
				assert.Equal(t, tt.expectedPrefix, tt.imgSource[index:index+len(prefix)])
			}
		})
	}
}
