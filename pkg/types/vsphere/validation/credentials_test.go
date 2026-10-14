package validation

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestValidateVCentersCredentialTypes(t *testing.T) {
	base := vsphere.VCenter{
		Server: "vcenter.example.com", Username: "global-user", Password: "global-pass", Datacenters: []string{"DC1"},
	}

	tests := []struct {
		name      string
		platform  vsphere.Platform
		wantError bool
	}{
		{
			name:     "global credentials remain valid",
			platform: vsphere.Platform{CredentialType: vsphere.CredentialTypeGlobal, VCenters: []vsphere.VCenter{base}},
		},
		{
			name: "global rejects component credentials",
			platform: vsphere.Platform{CredentialType: vsphere.CredentialTypeGlobal, VCenters: []vsphere.VCenter{{
				Server: "vcenter.example.com", Username: "global-user", Password: "global-pass", Datacenters: []string{"DC1"},
				ComponentCredentials: &vsphere.ComponentCredentials{},
			}}},
			wantError: true,
		},
		{
			name: "component scoped requires component credentials",
			platform: vsphere.Platform{CredentialType: vsphere.CredentialTypeComponentScoped, VCenters: []vsphere.VCenter{{
				Server: "vcenter.example.com", Datacenters: []string{"DC1"},
				ComponentCredentials: &vsphere.ComponentCredentials{
					MachineManagement:      vsphere.Credential{User: "machine", Password: "machine-pass"},
					Storage:                vsphere.Credential{User: "csi", Password: "csi-pass"},
					CloudControllerManager: vsphere.Credential{User: "ccm", Password: "ccm-pass"},
					VSphereProblemDetector: vsphere.Credential{User: "diag", Password: "diag-pass"},
				},
			}}},
		},
		{
			name: "component scoped rejects legacy credentials",
			platform: vsphere.Platform{CredentialType: vsphere.CredentialTypeComponentScoped, VCenters: []vsphere.VCenter{{
				Server: "vcenter.example.com", Username: "legacy", Datacenters: []string{"DC1"},
				ComponentCredentials: &vsphere.ComponentCredentials{
					MachineManagement:      vsphere.Credential{User: "machine", Password: "machine-pass"},
					Storage:                vsphere.Credential{User: "csi", Password: "csi-pass"},
					CloudControllerManager: vsphere.Credential{User: "ccm", Password: "ccm-pass"},
					VSphereProblemDetector: vsphere.Credential{User: "diag", Password: "diag-pass"},
				},
			}}},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateVCenters(&tt.platform, field.NewPath("vcenters"))
			if (len(errs) != 0) != tt.wantError {
				t.Fatalf("validateVCenters() errors = %v, wantError %v", errs, tt.wantError)
			}
		})
	}
}

func TestValidatePlatformRejectsUnknownCredentialType(t *testing.T) {
	platform := &vsphere.Platform{
		CredentialType: "per-component",
		VCenters: []vsphere.VCenter{{
			Server: "vcenter.example.com", Username: "user", Password: "pass", Datacenters: []string{"DC1"},
		}},
	}
	installConfig := &types.InstallConfig{Platform: types.Platform{VSphere: platform}}
	errs := ValidatePlatform(platform, true, field.NewPath("platform", "vsphere"), installConfig)
	if len(errs) == 0 || !strings.Contains(errs[0].Field, "credentialType") {
		t.Fatalf("expected credentialType validation error, got %v", errs)
	}
}
