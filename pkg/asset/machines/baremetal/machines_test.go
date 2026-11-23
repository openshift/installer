package baremetal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	baremetalprovider "github.com/openshift/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	"github.com/openshift/installer/pkg/types"
)

func TestHostSelectorForStream(t *testing.T) {
	cases := []struct {
		name     string
		scos     bool
		stream   types.OSImageStream
		expected baremetalprovider.HostSelector
	}{
		{
			name:   "rhel-9 excludes rhel-10",
			stream: types.OSImageStreamRHCOS9,
			expected: baremetalprovider.HostSelector{
				MatchExpressions: []baremetalprovider.HostSelectorRequirement{
					{
						Key:      streamLabelKey,
						Operator: selection.NotIn,
						Values:   []string{string(types.OSImageStreamRHCOS10)},
					},
				},
			},
		},
		{
			name:   "rhel-10 excludes rhel-9",
			stream: types.OSImageStreamRHCOS10,
			expected: baremetalprovider.HostSelector{
				MatchExpressions: []baremetalprovider.HostSelectorRequirement{
					{
						Key:      streamLabelKey,
						Operator: selection.NotIn,
						Values:   []string{string(types.OSImageStreamRHCOS9)},
					},
				},
			},
		},
		{
			name:     "scos centos-10 has no other streams, selector is empty",
			scos:     true,
			stream:   types.OSImageStreamCentos10,
			expected: baremetalprovider.HostSelector{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			old := types.SCOS
			types.SCOS = tc.scos
			defer func() { types.SCOS = old }()

			actual := hostSelectorForStream(tc.stream)
			assert.Equal(t, tc.expected, actual)

			// Every generated requirement must be a valid label selector
			// requirement; the machine controller rejects in/notin
			// requirements with an empty values list.
			for _, req := range actual.MatchExpressions {
				_, err := labels.NewRequirement(req.Key, req.Operator, req.Values)
				assert.NoError(t, err, "invalid host selector requirement %+v", req)
			}
		})
	}
}
