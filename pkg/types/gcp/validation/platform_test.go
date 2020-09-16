package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/gcp"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *gcp.Platform
		valid    bool
	}{
		{
			name: "minimal",
			platform: &gcp.Platform{
				Region: "us-east1",
			},
			valid: true,
		},
		{
			name: "invalid region",
			platform: &gcp.Platform{
				Region: "bad-region",
			},
			valid: false,
		},
		{
			name: "valid machine pool",
			platform: &gcp.Platform{
				Region:                 "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{},
			},
			valid: true,
		},
		{
			name: "valid subnets & network",
			platform: &gcp.Platform{
				Region:             "us-east1",
				Network:            "valid-vpc",
				ComputeSubnet:      "valid-compute-subnet",
				ControlPlaneSubnet: "valid-cp-subnet",
			},
			valid: true,
		},
		{
			name: "missing subnets",
			platform: &gcp.Platform{
				Region:  "us-east1",
				Network: "valid-vpc",
			},
			valid: false,
		},
		{
			name: "subnets missing network",
			platform: &gcp.Platform{
				Region:        "us-east1",
				ComputeSubnet: "valid-compute-subnet",
			},
			valid: false,
		},
		{
			name: "unsupported GCP disk type",
			platform: &gcp.Platform{
				Region: "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						DiskType: "pd-standard",
					},
				},
			},
			valid: false,
		},

		{
			name: "supported GCP disk type",
			platform: &gcp.Platform{
				Region: "us-east1",
				DefaultMachinePlatform: &gcp.MachinePool{
					OSDisk: gcp.OSDisk{
						DiskType: "pd-ssd",
					},
				},
			},
			valid: true,
		},
		{
			name: "valid labels",
			platform: &gcp.Platform{
				Region: "us-east1",
				Labels: map[string]string{
					"foo":       "bar",
					"test-name": "v4l1d",
				},
			},
			valid: true,
		},
		{
			name: "invalid label",
			platform: &gcp.Platform{
				Region: "us-east1",
				Labels: map[string]string{
					"UpperCase": "disallowed",
				},
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateLabels(t *testing.T) {
	cases := []struct {
		name   string
		labels map[string]string
		valid  bool
	}{
		{
			name: "valid labels",
			labels: map[string]string{
				"owner":         "1337_one",
				"this-is-fine":  "",
				"i_like_snakes": "hissssss---__",
				"a123":          "0",
				"тестовая":      "да",
				"κυβερνήτης":    "its-all-greek-to-me",

				"key-that-is-really-long-but-is-still-only-63-characters-so-fine": "value-that-is-also-63-characters-long-hooray_______123456789abc",
			},
			valid: true,
		},
		{
			name: "invalid: uppercase in key",
			labels: map[string]string{
				"camelCase": "disallowed",
			},
			valid: false,
		},
		{
			name: "invalid: uppercase in value",
			labels: map[string]string{
				"camel": "Case",
			},
			valid: false,
		},
		{
			name: "invalid: disallowed symbols in key",
			labels: map[string]string{
				"test.openshift.io/foo": "bar",
			},
			valid: false,
		},
		{
			name: "invalid: disallowed symbols in value",
			labels: map[string]string{
				"localhost": "127.0.0.1",
			},
			valid: false,
		},
		{
			name: "invalid: empty key",
			labels: map[string]string{
				"": "empty",
			},
			valid: false,
		},
		{
			name: "invalid: key too long",
			labels: map[string]string{
				"k234567890123456789012345678901234567890123456789012345678901234": "test",
			},
			valid: false,
		},
		{
			name: "invalid: value too long",
			labels: map[string]string{
				"test": "v234567890123456789012345678901234567890123456789012345678901234",
			},
			valid: false,
		},
		{
			name: "invalid: key doesn't start with letter",
			labels: map[string]string{
				"_private": "fails",
			},
			valid: false,
		},
		{
			name: "invalid: uppercase greek",
			labels: map[string]string{
				"Κυβερνήτης": "Kubernetes",
			},
			valid: false,
		},
		{
			name: "valid and invalid labels",
			labels: map[string]string{
				"a-valid-label":  "perfectly-fine",
				"4n1nv4l1dl4b3l": "1337haxx0r",
			},
			valid: false,
		},
		{
			name: "too many labels",
			labels: map[string]string{
				"a1": "value",
				"a2": "value",
				"a3": "value",
				"a4": "value",
				"a5": "value",
				"a6": "value",
				"a7": "value",
				"a8": "value",
				"a9": "value",
				"a0": "value",
				"b1": "value",
				"b2": "value",
				"b3": "value",
				"b4": "value",
				"b5": "value",
				"b6": "value",
				"b7": "value",
				"b8": "value",
				"b9": "value",
				"b0": "value",
				"c1": "value",
				"c2": "value",
				"c3": "value",
				"c4": "value",
				"c5": "value",
				"c6": "value",
				"c7": "value",
				"c8": "value",
				"c9": "value",
				"c0": "value",
				"d1": "value",
				"d2": "value",
				"d3": "value",
				"d4": "value",
				"d5": "value",
				"d6": "value",
				"d7": "value",
				"d8": "value",
				"d9": "value",
				"d0": "value",
				"e1": "value",
				"e2": "value",
				"e3": "value",
				"e4": "value",
				"e5": "value",
				"e6": "value",
				"e7": "value",
				"e8": "value",
				"e9": "value",
				"e0": "value",
				"f1": "value",
				"f2": "value",
				"f3": "value",
				"f4": "value",
				"f5": "value",
				"f6": "value",
				"f7": "value",
				"f8": "value",
				"f9": "value",
				"f0": "value",
				"g1": "value",
				"g2": "value",
				"g3": "value",
				"g4": "value",
				"g5": "value",
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLabels(tc.labels, field.NewPath("test-labels")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
