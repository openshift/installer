package clusterapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563

	"github.com/openshift/installer/pkg/asset/installconfig"
)

type fakeProvider struct {
	name string
}

func (f *fakeProvider) Name() string { return f.name }
func (f *fakeProvider) PublicGatherEndpoint() GatherEndpoint {
	return ExternalIP
}

type fakeBootstrapReadyProvider struct {
	fakeProvider
	called bool
	err    error
}

func (f *fakeBootstrapReadyProvider) BootstrapReady(_ context.Context, _ BootstrapReadyInput) error {
	f.called = true
	return f.err
}

func makeMachine(name, phase string, addresses []clusterv1.MachineAddress) *clusterv1.Machine {
	return &clusterv1.Machine{
		Status: clusterv1.MachineStatus{
			Phase:     phase,
			Addresses: addresses,
		},
	}
}

func TestCallBootstrapReadyHook(t *testing.T) {
	internalAddr := []clusterv1.MachineAddress{{
		Type:    clusterv1.MachineInternalIP,
		Address: "10.0.0.1",
	}}

	tests := []struct {
		name            string
		impl            Provider
		machines        []*clusterv1.Machine
		bootstrapName   string
		bootstrapPhase  string
		wantCalled      bool
		wantReturn      bool
		hookErr         error
		expectEarlyOK   bool
		reqBootstrapPub bool
	}{
		{
			name:           "bootstrap ready triggers hook",
			impl:           &fakeBootstrapReadyProvider{fakeProvider: fakeProvider{name: "test"}},
			bootstrapName:  "test-bootstrap",
			bootstrapPhase: string(clusterv1.MachinePhaseRunning),
			wantCalled:     true,
			wantReturn:     true,
		},
		{
			name:           "bootstrap not ready does not trigger hook",
			impl:           &fakeBootstrapReadyProvider{fakeProvider: fakeProvider{name: "test"}},
			bootstrapName:  "test-bootstrap",
			bootstrapPhase: string(clusterv1.MachinePhaseProvisioning),
			wantCalled:     false,
			wantReturn:     false,
		},
		{
			name:           "provider without hook returns true",
			impl:           &fakeProvider{name: "test"},
			bootstrapName:  "test-bootstrap",
			bootstrapPhase: string(clusterv1.MachinePhaseRunning),
			wantCalled:     false,
			wantReturn:     true,
		},
		{
			name:           "hook error returns true to prevent retries",
			impl:           &fakeBootstrapReadyProvider{fakeProvider: fakeProvider{name: "test"}, err: fmt.Errorf("ssh rule failed")},
			bootstrapName:  "test-bootstrap",
			bootstrapPhase: string(clusterv1.MachinePhaseRunning),
			wantCalled:     true,
			wantReturn:     true,
		},
		{
			name:           "bootstrap machine not in list returns false",
			impl:           &fakeBootstrapReadyProvider{fakeProvider: fakeProvider{name: "test"}},
			bootstrapName:  "test-bootstrap",
			bootstrapPhase: "",
			wantCalled:     false,
			wantReturn:     false,
			expectEarlyOK:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := &InfraProvider{impl: tt.impl}

			var machines []*clusterv1.Machine
			if !tt.expectEarlyOK {
				m := makeMachine(tt.bootstrapName, tt.bootstrapPhase, internalAddr)
				m.Name = tt.bootstrapName
				machines = []*clusterv1.Machine{m}
			}

			got := ip.callBootstrapReadyHook(
				context.Background(),
				nil,
				machines,
				tt.bootstrapName,
				tt.reqBootstrapPub,
				&installconfig.InstallConfig{},
				"test-infra",
			)

			assert.Equal(t, tt.wantReturn, got, "return value")

			if brp, ok := tt.impl.(*fakeBootstrapReadyProvider); ok {
				assert.Equal(t, tt.wantCalled, brp.called, "hook called")
			}
		})
	}
}

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
			name:      "failed - not ready",
			phase:     string(clusterv1.MachinePhaseFailed),
			wantReady: false,
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

func TestCheckMachineReadyFailedNotReady(t *testing.T) {
	machine := &clusterv1.Machine{
		Status: clusterv1.MachineStatus{
			Phase:          string(clusterv1.MachinePhaseFailed),
			FailureMessage: strPtr("VM allocation failed"), //nolint:staticcheck
		},
	}
	machine.Name = "test-machine"

	ready, err := checkMachineReady(machine, false)
	assert.NoError(t, err)
	assert.False(t, ready)
}

func strPtr(s string) *string { return &s }
