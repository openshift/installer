package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

func TestValidateCustomEndpoint(t *testing.T) {
	cases := []struct {
		name           string
		customEndpoint *aws.CustomEndpoint
		expected       string
	}{
		{
			name:           "empty",
			customEndpoint: &aws.CustomEndpoint{},
		},
		{
			name: "Valid Endpoint",
			customEndpoint: &aws.CustomEndpoint{
				Service: "ec2",
				URL:     "https://ec2.amazonaws.com",
			},
		},
		{
			name: "UnSupported Endpoint",
			customEndpoint: &aws.CustomEndpoint{
				Service: "sns",
				URL:     "https://sns.amazonaws.com",
			},
			expected: `^test-path\.Service: Unsupported value: "sns": supported values:.`,
		},
		{
			name: "Invalid Endpoint",
			customEndpoint: &aws.CustomEndpoint{
				Service: "ec2",
				URL:     "",
			},
			expected: `^test-path\.URL: Required value: Service Override URL cannot be empty`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCustomEndpoint(tc.customEndpoint, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
