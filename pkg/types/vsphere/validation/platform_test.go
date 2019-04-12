package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenter:          "test-vcenter",
		Username:         "test-username",
		Password:         "test-password",
		Datacenter:       "test-datacenter",
		DefaultDatastore: "test-datastore",
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		platform      *vsphere.Platform
		expectedError string
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "missing vCenter name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = ""
				return p
			}(),
			expectedError: `^test-path\.vCenter: Required value: must specify the name of the vCenter$`,
		},
		{
			name: "missing username",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Username = ""
				return p
			}(),
			expectedError: `^test-path\.username: Required value: must specify the username$`,
		},
		{
			name: "missing password",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Password = ""
				return p
			}(),
			expectedError: `^test-path\.password: Required value: must specify the password$`,
		},
		{
			name: "missing datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Datacenter = ""
				return p
			}(),
			expectedError: `^test-path\.datacenter: Required value: must specify the datacenter$`,
		},
		{
			name: "missing default datastore",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DefaultDatastore = ""
				return p
			}(),
			expectedError: `^test-path\.defaultDatastore: Required value: must specify the default datastore$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
