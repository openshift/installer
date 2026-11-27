package vsphere

import "testing"

func TestFormatServerWithPort(t *testing.T) {
	tests := []struct {
		name   string
		server string
		port   int32
		want   string
	}{
		{
			name:   "empty server",
			server: "",
			port:   8443,
			want:   "",
		},
		{
			name:   "port zero leaves server unchanged",
			server: "vcsim.lab.local",
			port:   0,
			want:   "vcsim.lab.local",
		},
		{
			name:   "default HTTPS port 443 leaves server unchanged",
			server: "vcsim.lab.local",
			port:   443,
			want:   "vcsim.lab.local",
		},
		{
			name:   "custom port is appended",
			server: "vcsim.lab.local",
			port:   8443,
			want:   "vcsim.lab.local:8443",
		},
		{
			name:   "server already containing port is unchanged",
			server: "vcsim.lab.local:8443",
			port:   8443,
			want:   "vcsim.lab.local:8443",
		},
		{
			name:   "IPv4 with custom port",
			server: "192.168.1.10",
			port:   8443,
			want:   "192.168.1.10:8443",
		},
		{
			name:   "IPv6 with custom port is bracketed",
			server: "2001:db8::1",
			port:   8443,
			want:   "[2001:db8::1]:8443",
		},
		{
			name:   "IPv6 with default port 443 is unchanged",
			server: "2001:db8::1",
			port:   443,
			want:   "2001:db8::1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := FormatServerWithPort(tc.server, tc.port); got != tc.want {
				t.Errorf("FormatServerWithPort(%q, %d) = %q, want %q", tc.server, tc.port, got, tc.want)
			}
		})
	}
}

func TestVCenterServerWithPort(t *testing.T) {
	vcCustom := VCenter{Server: "vcsim.lab.local", Port: 8443}
	if got := vcCustom.ServerWithPort(); got != "vcsim.lab.local:8443" {
		t.Errorf("VCenter.ServerWithPort() = %q, want %q", got, "vcsim.lab.local:8443")
	}

	vcDefault := VCenter{Server: "vcsim.lab.local", Port: 443}
	if got := vcDefault.ServerWithPort(); got != "vcsim.lab.local" {
		t.Errorf("VCenter.ServerWithPort() = %q, want %q", got, "vcsim.lab.local")
	}
}

func TestPlatformServerWithPort(t *testing.T) {
	var nilPlatform *Platform
	if got := nilPlatform.ServerWithPort("vcsim.lab.local"); got != "vcsim.lab.local" {
		t.Errorf("nil Platform.ServerWithPort() = %q, want %q", got, "vcsim.lab.local")
	}

	p := &Platform{
		VCenters: []VCenter{
			{Server: "vc-custom.lab.local", Port: 8443},
			{Server: "vc-default.lab.local", Port: 443},
		},
	}

	tests := []struct {
		name   string
		server string
		want   string
	}{
		{name: "matching custom port", server: "vc-custom.lab.local", want: "vc-custom.lab.local:8443"},
		{name: "matching default port", server: "vc-default.lab.local", want: "vc-default.lab.local"},
		{name: "unknown server returned unchanged", server: "other.lab.local", want: "other.lab.local"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := p.ServerWithPort(tc.server); got != tc.want {
				t.Errorf("Platform.ServerWithPort(%q) = %q, want %q", tc.server, got, tc.want)
			}
		})
	}
}
