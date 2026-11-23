package manifests

import (
	"testing"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestManualVSphereCredentialSecrets(t *testing.T) {
	secrets, err := manualVSphereCredentialSecrets(&vsphere.Platform{
		CredentialType: vsphere.CredentialTypeComponentScoped,
		VCenters: []vsphere.VCenter{{
			Server: "vcenter.example.com",
			ComponentCredentials: &vsphere.ComponentCredentials{
				MachineManagement:      vsphere.Credential{User: "machine", Password: "machine-pass"},
				Storage:                vsphere.Credential{User: "csi", Password: "csi-pass"},
				CloudControllerManager: vsphere.Credential{User: "ccm", Password: "ccm-pass"},
				VSphereProblemDetector: vsphere.Credential{User: "diag", Password: "diag-pass"},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]struct {
		name      string
		namespace string
		user      string
		password  string
	}{
		"99_vsphere-machine-api-credentials.yaml":             {"vsphere-cloud-credentials", "openshift-machine-api", "machine", "machine-pass"},
		"99_vsphere-machine-config-operator-credentials.yaml": {"vsphere-cloud-credentials", "openshift-machine-config-operator", "machine", "machine-pass"},
		"99_vsphere-csi-credentials.yaml":                     {"vmware-vsphere-cloud-credentials", "openshift-cluster-csi-drivers", "csi", "csi-pass"},
		"99_vsphere-cloud-controller-credentials.yaml":        {"vsphere-cloud-credentials", "openshift-cloud-controller-manager", "ccm", "ccm-pass"},
		"99_vsphere-problem-detector-credentials.yaml":        {"vsphere-cloud-credentials", "openshift-cluster-storage-operator", "diag", "diag-pass"},
	}
	if len(secrets) != len(want) {
		t.Errorf("got %d secrets, want %d", len(secrets), len(want))
	}
	for filename, expected := range want {
		secret, ok := secrets[filename]
		if !ok {
			t.Fatalf("missing %s", filename)
		}
		if secret.Namespace != expected.namespace {
			t.Errorf("%s namespace = %q, want %q", filename, secret.Namespace, expected.namespace)
		}
		if secret.Name != expected.name {
			t.Errorf("%s name = %q, want %q", filename, secret.Name, expected.name)
		}
		if got := string(secret.Data["vcenter.example.com.username"]); got != expected.user {
			t.Errorf("%s username = %q, want %q", filename, got, expected.user)
		}
		if got := string(secret.Data["vcenter.example.com.password"]); got != expected.password {
			t.Errorf("%s password = %q, want %q", filename, got, expected.password)
		}
	}
}

// TestManualVSphereCredentialSecretsMultiVCenter pins the accumulation of
// per-server credential keys: every vCenter entry contributes its own
// <server>.username and <server>.password keys to each component Secret.
func TestManualVSphereCredentialSecretsMultiVCenter(t *testing.T) {
	credentials := func(suffix string) *vsphere.ComponentCredentials {
		return &vsphere.ComponentCredentials{
			MachineManagement:      vsphere.Credential{User: "machine-" + suffix, Password: "machine-pass-" + suffix},
			Storage:                vsphere.Credential{User: "csi-" + suffix, Password: "csi-pass-" + suffix},
			CloudControllerManager: vsphere.Credential{User: "ccm-" + suffix, Password: "ccm-pass-" + suffix},
			VSphereProblemDetector: vsphere.Credential{User: "diag-" + suffix, Password: "diag-pass-" + suffix},
		}
	}
	secrets, err := manualVSphereCredentialSecrets(&vsphere.Platform{
		CredentialType: vsphere.CredentialTypeComponentScoped,
		VCenters: []vsphere.VCenter{
			{Server: "vcenter1.example.com", ComponentCredentials: credentials("one")},
			{Server: "vcenter2.example.com", ComponentCredentials: credentials("two")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(secrets) != 5 {
		t.Errorf("got %d secrets, want 5", len(secrets))
	}
	for filename, secret := range secrets {
		if len(secret.Data) != 4 {
			t.Errorf("%s got %d data keys, want 4 (two servers x username/password): %v", filename, len(secret.Data), secret.Data)
		}
		for _, server := range []string{"vcenter1.example.com", "vcenter2.example.com"} {
			for _, key := range []string{"username", "password"} {
				if _, ok := secret.Data[server+"."+key]; !ok {
					t.Errorf("%s missing %s.%s key", filename, server, key)
				}
			}
		}
	}
	// Spot-check one component across both servers.
	csi := secrets["99_vsphere-csi-credentials.yaml"]
	if got := string(csi.Data["vcenter1.example.com.username"]); got != "csi-one" {
		t.Errorf("csi vcenter1 username = %q, want %q", got, "csi-one")
	}
	if got := string(csi.Data["vcenter2.example.com.username"]); got != "csi-two" {
		t.Errorf("csi vcenter2 username = %q, want %q", got, "csi-two")
	}
}
