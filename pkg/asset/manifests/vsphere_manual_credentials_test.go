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
		namespace string
		user      string
		password  string
	}{
		"99_vsphere-machine-api-credentials.yaml":      {"openshift-machine-api", "machine", "machine-pass"},
		"99_vsphere-csi-credentials.yaml":              {"openshift-cluster-csi-drivers", "csi", "csi-pass"},
		"99_vsphere-cloud-controller-credentials.yaml": {"openshift-cloud-controller-manager", "ccm", "ccm-pass"},
		"99_vsphere-problem-detector-credentials.yaml": {"openshift-cluster-storage-operator", "diag", "diag-pass"},
	}
	for filename, expected := range want {
		secret, ok := secrets[filename]
		if !ok {
			t.Fatalf("missing %s", filename)
		}
		if secret.Namespace != expected.namespace {
			t.Errorf("%s namespace = %q, want %q", filename, secret.Namespace, expected.namespace)
		}
		if got := string(secret.Data["vcenter.example.com.username"]); got != expected.user {
			t.Errorf("%s username = %q, want %q", filename, got, expected.user)
		}
		if got := string(secret.Data["vcenter.example.com.password"]); got != expected.password {
			t.Errorf("%s password = %q, want %q", filename, got, expected.password)
		}
	}
}
