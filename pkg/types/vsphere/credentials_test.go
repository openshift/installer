package vsphere

import (
	"os"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func TestCredentialsForVCenter(t *testing.T) {
	tests := []struct {
		name       string
		platform   Platform
		server     string
		wantUser   string
		wantPass   string
		wantErrSub string
	}{
		{
			name: "global",
			platform: Platform{VCenters: []VCenter{{
				Server: "vcenter.example.com", Username: "global-user", Password: "global-pass",
			}}},
			server: "vcenter.example.com", wantUser: "global-user", wantPass: "global-pass",
		},
		{
			name: "component scoped uses machine API credentials",
			platform: Platform{
				CredentialType: CredentialTypeComponentScoped,
				VCenters: []VCenter{{
					Server:               "vcenter.example.com",
					ComponentCredentials: &ComponentCredentials{MachineManagement: Credential{User: "machine-user", Password: "machine-pass"}},
				}},
			},
			server: "vcenter.example.com", wantUser: "machine-user", wantPass: "machine-pass",
		},
		{
			name:     "unknown server",
			platform: Platform{VCenters: []VCenter{{Server: "vcenter.example.com", Username: "user", Password: "pass"}}},
			server:   "other.example.com", wantErrSub: "other.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.platform.CredentialsForVCenter(tt.server)
			if tt.wantErrSub != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrSub) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErrSub, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.User != tt.wantUser || got.Password != tt.wantPass {
				t.Errorf("got credentials %q/%q, want %q/%q", got.User, got.Password, tt.wantUser, tt.wantPass)
			}
		})
	}
}

// TestVCenterUserPasswordNotRequiredInCRD pins the published install-config
// schema: user and password are optional because component-scoped
// installations must omit them; presence is enforced for global credentials
// by install-config validation instead.
func TestVCenterUserPasswordNotRequiredInCRD(t *testing.T) {
	raw, err := os.ReadFile("../../../data/data/install.openshift.io_installconfigs.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var crd map[string]interface{}
	if err := yaml.Unmarshal(raw, &crd); err != nil {
		t.Fatal(err)
	}

	versions, found, err := unstructured.NestedSlice(crd, "spec", "versions")
	if err != nil || !found || len(versions) == 0 {
		t.Fatalf("expected spec.versions in CRD, found=%v err=%v", found, err)
	}
	for _, v := range versions {
		version, ok := v.(map[string]interface{})
		if !ok {
			t.Fatalf("expected version object, got %T", v)
		}
		required, found, err := unstructured.NestedStringSlice(version,
			"schema", "openAPIV3Schema", "properties", "platform", "properties", "vsphere",
			"properties", "vcenters", "items", "required")
		if err != nil || !found {
			t.Fatalf("expected vcenters required list in CRD, found=%v err=%v", found, err)
		}
		// Sanity: this is the vcenters entry list.
		if !strings.Contains(strings.Join(required, ","), "server") {
			t.Fatalf("unexpected vcenters required list: %v", required)
		}
		for _, field := range []string{"user", "password"} {
			if strings.Contains(strings.Join(required, ","), field) {
				t.Errorf("vcenters entry must not require %q for component-scoped credentials, required list: %v", field, required)
			}
		}
	}
}

// TestFailureDomainDescriptionInCRD pins the published description of the
// failureDomains entries: the FailureDomain doc comment must survive
// regeneration of the install-config CRD intact, because the description is
// user-facing API documentation.
func TestFailureDomainDescriptionInCRD(t *testing.T) {
	raw, err := os.ReadFile("../../../data/data/install.openshift.io_installconfigs.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var crd map[string]interface{}
	if err := yaml.Unmarshal(raw, &crd); err != nil {
		t.Fatal(err)
	}

	versions, found, err := unstructured.NestedSlice(crd, "spec", "versions")
	if err != nil || !found || len(versions) == 0 {
		t.Fatalf("expected spec.versions in CRD, found=%v err=%v", found, err)
	}
	for _, v := range versions {
		version, ok := v.(map[string]interface{})
		if !ok {
			t.Fatalf("expected version object, got %T", v)
		}
		description, found, err := unstructured.NestedString(version,
			"schema", "openAPIV3Schema", "properties", "platform", "properties", "vsphere",
			"properties", "failureDomains", "items", "description")
		if err != nil || !found {
			t.Fatalf("expected failureDomains items description in CRD, found=%v err=%v", found, err)
		}
		if !strings.HasPrefix(description, "FailureDomain holds the region and zone failure domain and") {
			t.Errorf("failureDomains items description lost its first line, got: %q", description)
		}
	}
}
