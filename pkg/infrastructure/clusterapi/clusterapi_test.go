package clusterapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563
)

func TestCheckMachineReady(t *testing.T) {
	tests := []struct {
		name         string
		phase        string
		addresses    []clusterv1.MachineAddress
		requirePubIP bool
		wantReady    bool
		wantErr      bool
	}{
		{
			name:  "provisioned with internal IP",
			phase: string(clusterv1.MachinePhaseProvisioned),
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineInternalIP, Address: "10.0.0.1"},
			},
			wantReady: true,
		},
		{
			name:  "running with external IP",
			phase: string(clusterv1.MachinePhaseRunning),
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineExternalIP, Address: "1.2.3.4"},
			},
			wantReady: true,
		},
		{
			name:      "provisioning - not ready",
			phase:     string(clusterv1.MachinePhaseProvisioning),
			wantReady: false,
		},
		{
			name:    "failed - error",
			phase:   string(clusterv1.MachinePhaseFailed),
			wantErr: true,
		},
		{
			name:  "requires public IP but only has internal",
			phase: string(clusterv1.MachinePhaseProvisioned),
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineInternalIP, Address: "10.0.0.1"},
			},
			requirePubIP: true,
			wantReady:    false,
		},
		{
			name:  "requires public IP and has external",
			phase: string(clusterv1.MachinePhaseProvisioned),
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineExternalIP, Address: "1.2.3.4"},
			},
			requirePubIP: true,
			wantReady:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine := &clusterv1.Machine{
				Status: clusterv1.MachineStatus{
					Phase:     tt.phase,
					Addresses: tt.addresses,
				},
			}
			machine.Name = "test-machine"

			ready, err := checkMachineReady(machine, tt.requirePubIP)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantReady, ready)
			}
		})
	}
}

func TestHasRequiredIP(t *testing.T) {
	tests := []struct {
		name         string
		addresses    []clusterv1.MachineAddress
		requirePubIP bool
		want         bool
	}{
		{
			name: "external IP satisfies public requirement",
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineExternalIP, Address: "1.2.3.4"},
			},
			requirePubIP: true,
			want:         true,
		},
		{
			name: "internal IP does not satisfy public requirement",
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineInternalIP, Address: "10.0.0.1"},
			},
			requirePubIP: true,
			want:         false,
		},
		{
			name: "internal IP satisfies non-public requirement",
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineInternalIP, Address: "10.0.0.1"},
			},
			requirePubIP: false,
			want:         true,
		},
		{
			name:         "no addresses",
			addresses:    nil,
			requirePubIP: false,
			want:         false,
		},
		{
			name: "empty address ignored",
			addresses: []clusterv1.MachineAddress{
				{Type: clusterv1.MachineInternalIP, Address: ""},
			},
			requirePubIP: false,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine := &clusterv1.Machine{
				Status: clusterv1.MachineStatus{
					Addresses: tt.addresses,
				},
			}
			machine.Name = "test-machine"

			got := hasRequiredIP(machine, tt.requirePubIP)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckMachineReadyFailed(t *testing.T) {
	machine := &clusterv1.Machine{
		Status: clusterv1.MachineStatus{
			Phase:          string(clusterv1.MachinePhaseFailed),
			FailureMessage: strPtr("VM allocation failed"), //nolint:staticcheck
		},
	}
	machine.Name = "test-machine"

	ready, err := checkMachineReady(machine, false)
	assert.Error(t, err)
	assert.False(t, ready)
	assert.Contains(t, err.Error(), "VM allocation failed")
}

func strPtr(s string) *string { return &s }
