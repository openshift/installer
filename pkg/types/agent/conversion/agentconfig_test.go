package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types/agent"
)

func TestConvertAgentConfig(t *testing.T) {
	cases := []struct {
		name          string
		config        *agent.Config
		expected      *agent.Config
		expectedError string
	}{
		{
			name: "empty",
			config: &agent.Config{
				TypeMeta: metav1.TypeMeta{
					APIVersion: agent.AgentConfigVersion,
				},
			},
			expected: &agent.Config{
				TypeMeta: metav1.TypeMeta{
					APIVersion: agent.AgentConfigVersion,
				},
			},
		},
		{
			name: "v1alpha1",
			config: &agent.Config{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1alpha1",
				},
			},
			expected: &agent.Config{
				TypeMeta: metav1.TypeMeta{
					APIVersion: agent.AgentConfigVersion,
				},
			},
		},
		{
			name:          "no version",
			config:        &agent.Config{},
			expectedError: "no version was provided",
		},
		{
			name: "bad version",
			config: &agent.Config{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1alpha0",
				},
			},
			expectedError: "cannot upconvert from version v1alpha0",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ConvertAgentConfig(tc.config)
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, tc.config, "unexpected install config")
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
