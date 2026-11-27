package vsphere

import (
	"context"
	"fmt"
	"testing"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta2"

	icvsphere "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func TestGenerateMachinesServerPort(t *testing.T) {
	tests := []struct {
		name string
		port int32
	}{
		{
			name: "custom port 8443 is included on control plane and bootstrap machines",
			port: 8443,
		},
		{
			name: "default port 443 omits port",
			port: 443,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, err := parseInstallConfig()
			if err != nil {
				t.Fatalf("parseInstallConfig() failed: %v", err)
			}

			for i := range config.VSphere.VCenters {
				config.VSphere.VCenters[i].Port = tc.port
			}
			for i := range config.VSphere.FailureDomains {
				config.VSphere.FailureDomains[i].Topology.Networks = []string{"/dc1/network/vm-network"}
			}

			wantServer := vspheretypes.FormatServerWithPort(config.VSphere.VCenters[0].Server, tc.port)
			wantClusterName := fmt.Sprintf("%s-0", testClusterID)

			files, err := GenerateMachines(
				context.Background(),
				testClusterID,
				config,
				&machinePoolValidZones,
				masterRole,
				icvsphere.NewMetadata(),
			)
			if err != nil {
				t.Fatalf("GenerateMachines() failed: %v", err)
			}

			var vsphereMachineCount int
			var foundBootstrapInfra, foundBootstrapMachine bool

			for _, f := range files {
				switch obj := f.Object.(type) {
				case *capv.VSphereMachine:
					vsphereMachineCount++
					if obj.Spec.Server != wantServer {
						t.Errorf("VSphereMachine %s Spec.Server = %q, want %q", obj.Name, obj.Spec.Server, wantServer)
					}
					if obj.Name == fmt.Sprintf("%s-bootstrap", testClusterID) {
						foundBootstrapInfra = true
					}
				case *capi.Machine:
					if obj.Spec.ClusterName != wantClusterName {
						t.Errorf("Machine %s Spec.ClusterName = %q, want %q", obj.Name, obj.Spec.ClusterName, wantClusterName)
					}
					if obj.Name == fmt.Sprintf("%s-bootstrap", testClusterID) {
						foundBootstrapMachine = true
					}
				}
			}

			if vsphereMachineCount == 0 {
				t.Fatal("expected at least one VSphereMachine manifest")
			}
			if !foundBootstrapInfra {
				t.Fatal("expected bootstrap VSphereMachine manifest")
			}
			if !foundBootstrapMachine {
				t.Fatal("expected bootstrap CAPI Machine manifest")
			}
		})
	}
}
