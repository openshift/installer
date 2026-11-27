package vsphere

import (
	"context"
	"reflect"
	"testing"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func TestMetadataAddCredentials(t *testing.T) {
	tests := []struct {
		name       string
		server     string
		port       int32
		wantServer string
	}{
		{
			name:       "custom port is included in session params",
			server:     "vcsim.lab.local",
			port:       8443,
			wantServer: "vcsim.lab.local:8443",
		},
		{
			name:       "default port 443 leaves server unchanged",
			server:     "vcsim.lab.local",
			port:       443,
			wantServer: "vcsim.lab.local",
		},
		{
			name:       "unset port 0 leaves server unchanged",
			server:     "vcsim.lab.local",
			port:       0,
			wantServer: "vcsim.lab.local",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMetadata()
			got := m.AddCredentials(tc.server, tc.port, "user", "pass")

			stored, ok := m.VCenterCredentials[tc.server]
			if !ok {
				t.Fatalf("expected VCenterCredentials to be keyed by bare server %q", tc.server)
			}
			if stored.Port != tc.port {
				t.Errorf("stored Port = %d, want %d", stored.Port, tc.port)
			}

			want := session.NewParams().WithServer(tc.wantServer).WithUserInfo("user", "pass")
			if !reflect.DeepEqual(got, want) {
				t.Errorf("AddCredentials() params = %+v, want %+v", got, want)
			}
		})
	}
}

func TestMetadataSessionRebuildsCredentialsWithPort(t *testing.T) {
	// Simulate Metadata restored from .openshift_install_state.json
	// where VCenterCredentials is populated but in-memory maps are nil.
	m := &Metadata{
		VCenterCredentials: map[string]VCenterCredential{
			"vcsim.lab.local": {
				Username: "user",
				Password: "pass",
				Port:     8443,
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Session fails fast due to the canceled context, after rebuilding m.credentials.
	if _, err := m.Session(ctx, "vcsim.lab.local"); err == nil {
		t.Fatal("expected Session() to fail with canceled context")
	}

	got, ok := m.credentials["vcsim.lab.local"]
	if !ok {
		t.Fatal("expected credentials to be rebuilt under bare server key")
	}
	want := session.NewParams().WithServer("vcsim.lab.local:8443").WithUserInfo("user", "pass")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("rebuilt params = %+v, want %+v", got, want)
	}
}
