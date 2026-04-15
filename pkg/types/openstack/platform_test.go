package openstack

import (
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

// TestBootstrapFlavorYAMLParsing verifies that the bootstrapFlavor field in
// the Platform struct is correctly parsed from YAML in all expected scenarios.
func TestBootstrapFlavorYAMLParsing(t *testing.T) {
	tests := []struct {
		name         string
		yaml         string
		wantFlavor   string
		wantParseErr bool
	}{
		{
			name: "bootstrapFlavor with valid value",
			yaml: `
cloud: mycloud
bootstrapFlavor: m1.xlarge
`,
			wantFlavor: "m1.xlarge",
		},
		{
			name: "bootstrapFlavor empty string",
			yaml: `
cloud: mycloud
bootstrapFlavor: ""
`,
			wantFlavor: "",
		},
		{
			name: "bootstrapFlavor omitted (null)",
			yaml: `
cloud: mycloud
`,
			wantFlavor: "",
		},
		{
			name: "case sensitivity preservation - all caps",
			yaml: `
cloud: mycloud
bootstrapFlavor: M1.XLARGE
`,
			wantFlavor: "M1.XLARGE",
		},
		{
			name: "case sensitivity preservation - exact mixed case",
			yaml: `
cloud: mycloud
bootstrapFlavor: Bootstrap-Flavor-Large
`,
			wantFlavor: "Bootstrap-Flavor-Large",
		},
		{
			name: "case sensitivity preservation - partial caps",
			yaml: `
cloud: mycloud
bootstrapFlavor: m1.XLARGE
`,
			wantFlavor: "m1.XLARGE",
		},
		{
			name: "bootstrapFlavor with spaces in name",
			yaml: `
cloud: mycloud
bootstrapFlavor: "my flavor large"
`,
			wantFlavor: "my flavor large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var platform Platform
			err := yaml.Unmarshal([]byte(tt.yaml), &platform)
			if tt.wantParseErr {
				if err == nil {
					t.Errorf("expected parse error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			if platform.BootstrapFlavor != tt.wantFlavor {
				t.Errorf("BootstrapFlavor = %q, want %q", platform.BootstrapFlavor, tt.wantFlavor)
			}
		})
	}
}

// TestBootstrapFlavorCaseSensitivityDistinct verifies that different
// capitalizations of the same name produce distinct (not equal) values,
// confirming exact case preservation.
func TestBootstrapFlavorCaseSensitivityDistinct(t *testing.T) {
	yamlLower := `
cloud: mycloud
bootstrapFlavor: m1.xlarge
`
	yamlUpper := `
cloud: mycloud
bootstrapFlavor: M1.XLARGE
`

	var pLower, pUpper Platform
	if err := yaml.Unmarshal([]byte(yamlLower), &pLower); err != nil {
		t.Fatalf("unexpected parse error for lower: %v", err)
	}
	if err := yaml.Unmarshal([]byte(yamlUpper), &pUpper); err != nil {
		t.Fatalf("unexpected parse error for upper: %v", err)
	}

	if pLower.BootstrapFlavor == pUpper.BootstrapFlavor {
		t.Errorf("expected case-distinct flavors to differ, but both = %q", pLower.BootstrapFlavor)
	}
}

// TestBootstrapFlavorJSONTagOmitempty verifies that when BootstrapFlavor is
// empty, the field is omitted from the JSON/YAML output (omitempty behavior).
func TestBootstrapFlavorJSONTagOmitempty(t *testing.T) {
	platform := Platform{
		Cloud: "mycloud",
	}
	data, err := yaml.Marshal(platform)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	output := string(data)
	if strings.Contains(output, "bootstrapFlavor") {
		t.Errorf("expected bootstrapFlavor to be omitted when empty, but got output: %s", output)
	}
}

// TestBootstrapFlavorJSONTagPresent verifies that when BootstrapFlavor is
// set, it appears in the JSON/YAML output under the correct key name.
func TestBootstrapFlavorJSONTagPresent(t *testing.T) {
	platform := Platform{
		Cloud:           "mycloud",
		BootstrapFlavor: "m1.xlarge",
	}
	data, err := yaml.Marshal(platform)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	output := string(data)
	if !strings.Contains(output, "bootstrapFlavor: m1.xlarge") {
		t.Errorf("expected bootstrapFlavor key in output, got: %s", output)
	}
}
