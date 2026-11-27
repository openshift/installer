package vsphere

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func TestGenerateClusterAssetsUsesMachineManagementCredentials(t *testing.T) {
	config := &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		BaseDomain: "example.com",
		Platform: types.Platform{VSphere: &vspheretypes.Platform{
			CredentialType: vspheretypes.CredentialTypeComponentScoped,
			VCenters: []vspheretypes.VCenter{{
				Server: "vcenter.example.com",
				ComponentCredentials: &vspheretypes.ComponentCredentials{
					MachineManagement:      vspheretypes.Credential{User: "machine-user", Password: "machine-pass"},
					Storage:                vspheretypes.Credential{User: "storage-user", Password: "storage-pass"},
					CloudControllerManager: vspheretypes.Credential{User: "ccm-user", Password: "ccm-pass"},
					VSphereProblemDetector: vspheretypes.Credential{User: "problem-detector-user", Password: "problem-detector-pass"},
				},
			}},
		}},
	}

	output, err := GenerateClusterAssets(installconfig.MakeAsset(config), &installconfig.ClusterID{InfraID: "test-cluster"})
	if err != nil {
		t.Fatal(err)
	}
	if len(output.Manifests) == 0 {
		t.Fatal("expected CAPI manifests")
	}

	credentials, ok := output.Manifests[0].Object.(*corev1.Secret)
	if !ok {
		t.Fatalf("first manifest is %T, want *corev1.Secret", output.Manifests[0].Object)
	}
	if got := string(credentials.Data["username"]); got != "machine-user" {
		t.Errorf("username = %q, want %q", got, "machine-user")
	}
	if got := string(credentials.Data["password"]); got != "machine-pass" {
		t.Errorf("password = %q, want %q", got, "machine-pass")
	}
}

func TestGenerateClusterAssetsServerPort(t *testing.T) {
	tests := []struct {
		name       string
		port       int32
		wantServer string
	}{
		{
			name:       "custom port is included in cluster and deployment zone",
			port:       8443,
			wantServer: "https://vcenter.example.com:8443",
		},
		{
			name:       "default port 443 omits port",
			port:       443,
			wantServer: "https://vcenter.example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				BaseDomain: "example.com",
				Platform: types.Platform{VSphere: &vspheretypes.Platform{
					VCenters: []vspheretypes.VCenter{{
						Server:   "vcenter.example.com",
						Port:     tc.port,
						Username: "user",
						Password: "pass",
					}},
					FailureDomains: []vspheretypes.FailureDomain{{
						Name:     "fd-hostgroup",
						Server:   "vcenter.example.com",
						Region:   "region-1",
						Zone:     "zone-1",
						ZoneType: vspheretypes.HostGroupFailureDomain,
					}},
				}},
			}

			output, err := GenerateClusterAssets(installconfig.MakeAsset(config), &installconfig.ClusterID{InfraID: "test-cluster"})
			if err != nil {
				t.Fatal(err)
			}

			var foundCluster, foundDeploymentZone bool
			for _, f := range output.Manifests {
				switch obj := f.Object.(type) {
				case *capv.VSphereCluster:
					foundCluster = true
					if obj.Spec.Server != tc.wantServer {
						t.Errorf("VSphereCluster.Spec.Server = %q, want %q", obj.Spec.Server, tc.wantServer)
					}
				case *capv.VSphereDeploymentZone:
					foundDeploymentZone = true
					if obj.Spec.Server != tc.wantServer {
						t.Errorf("VSphereDeploymentZone.Spec.Server = %q, want %q", obj.Spec.Server, tc.wantServer)
					}
				}
			}

			if !foundCluster {
				t.Fatal("expected VSphereCluster manifest to be generated")
			}
			if !foundDeploymentZone {
				t.Fatal("expected VSphereDeploymentZone manifest to be generated")
			}
		})
	}
}
