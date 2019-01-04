package validate

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterName(t *testing.T) {
	maxSizeName := strings.Repeat("123456789.", 25) + "123"

	cases := []struct {
		name        string
		clusterName string
		valid       bool
	}{
		{"empty", "", false},
		{"only whitespace", " ", false},
		{"single lowercase", "a", true},
		{"single uppercase", "A", false},
		{"contains whitespace", "abc D", false},
		{"single number", "1", true},
		{"single dot", ".", false},
		{"ends with dot", "a.", false},
		{"starts with dot", ".a", false},
		{"multiple labels", "a.a", true},
		{"starts with dash", "-a", false},
		{"ends with dash", "a-", false},
		{"label starts with dash", "a.-a", false},
		{"label ends with dash", "a-.a", false},
		{"invalid percent", "a%a", false},
		{"only non-ascii", "日本語", false},
		{"contains non-ascii", "a日本語a", false},
		{"max size", maxSizeName, true},
		{"too long", maxSizeName + "a", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ClusterName(tc.clusterName)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSubnetCIDR(t *testing.T) {
	cases := []struct {
		cidr  string
		valid bool
	}{
		{"0.0.0.0/32", false},
		{"1.2.3.4/0", false},
		{"1.2.3.4/1", false},
		{"1.2.3.4/31", true},
		{"1.2.3.4/32", true},
		{"0:0:0:0:0:1:102:304/116", false},
		{"0:0:0:0:0:ffff:102:304/116", true},
		{"172.17.1.2/20", false},
		{"172.17.1.2/8", false},
		{"255.255.255.255/1", false},
		{"255.255.255.255/32", true},
	}
	for _, tc := range cases {
		t.Run(tc.cidr, func(t *testing.T) {
			_, cidr, err := net.ParseCIDR(tc.cidr)
			if err != nil {
				t.Fatalf("could not parse cidr: %v", err)
			}
			err = SubnetCIDR(cidr)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDomainName(t *testing.T) {
	cases := []struct {
		domain string
		valid  bool
	}{
		{"", false},
		{" ", false},
		{"a", true},
		{".", false},
		{"日本語", false},
		{"日本語.com", false},
		{"abc.日本語.com", false},
		{"a日本語a.com", false},
		{"abc", true},
		{"ABC", false},
		{"ABC123", false},
		{"ABC123.COM123", false},
		{"1", true},
		{"0.0", true},
		{"1.2.3.4", true},
		{"1.2.3.4.", true},
		{"abc.", true},
		{"abc.com", true},
		{"abc.com.", true},
		{"a.b.c.d.e.f", true},
		{".abc", false},
		{".abc.com", false},
		{".abc.com", false},
	}
	for _, tc := range cases {
		t.Run(tc.domain, func(t *testing.T) {
			err := DomainName(tc.domain)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDoCIDRsOverlap(t *testing.T) {
	cases := []struct {
		a       string
		b       string
		overlap bool
	}{
		{
			a:       "192.168.0.0/30",
			b:       "192.168.0.3/30",
			overlap: true,
		},
		{
			a:       "192.168.0.0/30",
			b:       "192.168.0.4/30",
			overlap: false,
		},
		{
			a:       "192.168.0.0/29",
			b:       "192.168.0.4/30",
			overlap: true,
		},
		{
			a:       "0.0.0.0/0",
			b:       "192.168.0.0/24",
			overlap: true,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %s", tc.a, tc.b), func(t *testing.T) {
			_, a, err := net.ParseCIDR(tc.a)
			if err != nil {
				t.Fatalf("could not parse cidr %q: %v", tc.a, err)
			}
			_, b, err := net.ParseCIDR(tc.b)
			if err != nil {
				t.Fatalf("could not parse cidr %q: %v", tc.b, err)
			}
			actual := DoCIDRsOverlap(a, b)
			assert.Equal(t, tc.overlap, actual)
		})
	}
}

func TestImagePullSecret(t *testing.T) {
	cases := []struct {
		name   string
		secret string
		valid  bool
	}{
		{
			name:   "single entry with auth",
			secret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
			valid:  true,
		},
		{
			name:   "single entry with credsStore",
			secret: `{"auths":{"example.com":{"credsStore":"creds store value"}}}`,
			valid:  true,
		},
		{
			name:   "empty",
			secret: `{}`,
			valid:  false,
		},
		{
			name:   "no auths",
			secret: `{"not-auths":{"example.com":{"auth":"authorization value"}}}`,
			valid:  false,
		},
		{
			name:   "no auth or credsStore",
			secret: `{"auths":{"example.com":{"unrequired-field":"value"}}}`,
			valid:  false,
		},
		{
			name:   "additional fields",
			secret: `{"auths":{"example.com":{"auth":"authorization value","other-field":"other field value"}}}`,
			valid:  true,
		},
		{
			name:   "no entries",
			secret: `{"auths":{}}`,
			valid:  false,
		},
		{
			name:   "multiple valid entries",
			secret: `{"auths":{"example.com":{"auth":"authorization value"},"other-example.com":{"auth":"other auth value"}}}`,
			valid:  true,
		},
		{
			name:   "mix of valid and invalid entries",
			secret: `{"auths":{"example.com":{"auth":"authorization value"},"other-example.com":{"unrequired-field":"value"}}}`,
			valid:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ImagePullSecret(tc.secret)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestSSHPublicKey(t *testing.T) {
	cases := []struct {
		name  string
		key   string
		valid bool
	}{
		{
			name:  "valid",
			key:   "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q==",
			valid: true,
		},
		{
			name:  "valid with email",
			key:   "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q== name@example.com",
			valid: true,
		},
		{
			name:  "invalid format",
			key:   "bad-format AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q==",
			valid: true,
		},
		{
			name:  "invalid key",
			key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL",
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := SSHPublicKey(tc.key)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestURI(t *testing.T) {
	cases := []struct {
		name  string
		uri   string
		valid bool
	}{
		{
			name:  "valid",
			uri:   "https://example.com",
			valid: true,
		},
		{
			name:  "missing scheme",
			uri:   "example.com",
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := URI(tc.uri)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
