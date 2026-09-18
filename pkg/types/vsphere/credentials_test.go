package vsphere

import (
	"strings"
	"testing"
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
