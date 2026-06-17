package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

// TestBootstrapFlavorYAMLDeserialization verifies that the bootstrapFlavor field
// correctly round-trips through YAML serialization and deserialization. The
// field must flow through unchanged (no defaulting) so that an empty value is
// preserved as empty and a specified value is preserved as-is.
func TestBootstrapFlavorYAMLDeserialization(t *testing.T) {
	cases := []struct {
		name     string
		yaml     string
		expected string
	}{
		{
			name: "bootstrapFlavor specified",
			yaml: `
cloud: mycloud
bootstrapFlavor: m1.xlarge
`,
			expected: "m1.xlarge",
		},
		{
			name: "bootstrapFlavor not set",
			yaml: `
cloud: mycloud
`,
			expected: "",
		},
		{
			name: "bootstrapFlavor empty string",
			yaml: `
cloud: mycloud
bootstrapFlavor: ""
`,
			expected: "",
		},
		{
			name: "bootstrapFlavor with hyphen and digits",
			yaml: `
cloud: mycloud
bootstrapFlavor: ocp.bootstrap-4cpu
`,
			expected: "ocp.bootstrap-4cpu",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var p Platform
			if err := yaml.Unmarshal([]byte(tc.yaml), &p); err != nil {
				t.Fatalf("unexpected error unmarshaling platform YAML: %v", err)
			}
			assert.Equal(t, tc.expected, p.BootstrapFlavor,
				"BootstrapFlavor should be deserialized correctly from YAML without any defaulting")
		})
	}
}

// TestBootstrapFlavorYAMLSerialization verifies that the bootstrapFlavor field
// is correctly serialized to YAML. When empty it should be omitted (omitempty),
// when set it should appear with the exact value.
func TestBootstrapFlavorYAMLSerialization(t *testing.T) {
	t.Run("empty bootstrapFlavor omitted from YAML output", func(t *testing.T) {
		p := Platform{Cloud: "mycloud"}
		data, err := yaml.Marshal(&p)
		if err != nil {
			t.Fatalf("unexpected error marshaling platform: %v", err)
		}
		// With omitempty, an empty BootstrapFlavor should not appear in the output.
		assert.NotContains(t, string(data), "bootstrapFlavor",
			"empty bootstrapFlavor should be omitted from YAML output")
	})

	t.Run("non-empty bootstrapFlavor included in YAML output", func(t *testing.T) {
		p := Platform{Cloud: "mycloud", BootstrapFlavor: "m1.large"}
		data, err := yaml.Marshal(&p)
		if err != nil {
			t.Fatalf("unexpected error marshaling platform: %v", err)
		}
		assert.Contains(t, string(data), "bootstrapFlavor: m1.large",
			"non-empty bootstrapFlavor should appear in YAML output")
	})
}
