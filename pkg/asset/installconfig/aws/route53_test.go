package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetSOAMinimumTTL(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected string
		wantErr  bool
	}{
		{
			name:     "lowers the minimum field only",
			value:    "ns-2048.awsdns-64.net. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400",
			expected: "ns-2048.awsdns-64.net. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 60",
		},
		{
			name:     "preserves an already-minimal record other than the minimum field",
			value:    "ns-2048.awsdns-64.net. awsdns-hostmaster.amazon.com. 5 7200 900 1209600 60",
			expected: "ns-2048.awsdns-64.net. awsdns-hostmaster.amazon.com. 5 7200 900 1209600 60",
		},
		{
			name:    "too few fields",
			value:   "ns-2048.awsdns-64.net. awsdns-hostmaster.amazon.com. 1 7200 900 1209600",
			wantErr: true,
		},
		{
			name:    "too many fields",
			value:   "ns-2048.awsdns-64.net. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400 extra",
			wantErr: true,
		},
		{
			name:    "empty value",
			value:   "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := setSOAMinimumTTL(tc.value)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
