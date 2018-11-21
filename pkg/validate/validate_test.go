package validate

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastIP(t *testing.T) {
	cases := []struct {
		in  net.IPNet
		out net.IP
	}{
		{
			in: net.IPNet{
				IP:   net.ParseIP("192.168.0.0").To4(),
				Mask: net.CIDRMask(24, 32),
			},
			out: net.ParseIP("192.168.0.255"),
		},
		{
			in: net.IPNet{
				IP:   net.ParseIP("192.168.0.0").To4(),
				Mask: net.CIDRMask(22, 32),
			},
			out: net.ParseIP("192.168.3.255"),
		},
		{
			in: net.IPNet{
				IP:   net.ParseIP("192.168.0.0").To4(),
				Mask: net.CIDRMask(32, 32),
			},
			out: net.ParseIP("192.168.0.0"),
		},
		{
			in: net.IPNet{
				IP:   net.ParseIP("0.0.0.0").To4(),
				Mask: net.CIDRMask(0, 32),
			},
			out: net.ParseIP("255.255.255.255"),
		},
	}

	var out net.IP
	for i, c := range cases {
		if out = lastIP(&c.in); out.String() != c.out.String() {
			t.Errorf("test case %d: expected %s but got %s", i, c.out, out)
		}
	}
}

const caseMsg = "must be lower case"
const emptyMsg = "cannot be empty"
const invalidDomainMsg = "invalid domain name"
const invalidHostMsg = "invalid host (must be a domain name or IP address)"
const invalidIPMsg = "invalid IPv4 address"
const invalidIntMsg = "invalid integer"
const invalidPortMsg = "invalid port number"
const noCIDRNetmaskMsg = "must provide a CIDR netmask (eg, /24)"

type test struct {
	in       string
	expected string
}

type validator func(string) error

func runTests(t *testing.T, funcName string, fn validator, tests []test) {
	for _, test := range tests {
		err := fn(test.in)
		if (err == nil && test.expected != "") || (err != nil && err.Error() != test.expected) {
			t.Errorf("For %s(%q), expected %q, got %q", funcName, test.in, test.expected, err)
		}
	}
}

func TestNonEmpty(t *testing.T) {
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", ""},
		{".", ""},
		{"日本語", ""},
	}
	runTests(t, "NonEmpty", nonEmpty, tests)
}

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

func TestIPv4(t *testing.T) {
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"0.0.0.0", ""},
		{"1.2.3.4", ""},
		{"1.2.3.", invalidIPMsg},
		{"1.2.3.4.", invalidIPMsg},
		{"1.2.3.a", invalidIPMsg},
		{"255.255.255.255", ""},
	}
	runTests(t, "IPv4", IPv4, tests)
}

func TestSubnetCIDR(t *testing.T) {
	const netmaskSizeMsg = "invalid netmask size (must be between 0 and 32)"

	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"/16", invalidIPMsg},
		{"0.0.0.0/0", ""},
		{"0.0.0.0/32", ""},
		{"1.2.3.4", noCIDRNetmaskMsg},
		{"1.2.3.", noCIDRNetmaskMsg},
		{"1.2.3.4.", noCIDRNetmaskMsg},
		{"1.2.3.4/0", ""},
		{"1.2.3.4/1", ""},
		{"1.2.3.4/31", ""},
		{"1.2.3.4/32", ""},
		{"1.2.3./16", invalidIPMsg},
		{"1.2.3.4./16", invalidIPMsg},
		{"1.2.3.4/33", netmaskSizeMsg},
		{"1.2.3.4/-1", netmaskSizeMsg},
		{"1.2.3.4/abc", netmaskSizeMsg},
		{"172.17.1.2", noCIDRNetmaskMsg},
		{"172.17.1.2/", netmaskSizeMsg},
		{"172.17.1.2/33", netmaskSizeMsg},
		{"172.17.1.2/20", "overlaps with default Docker Bridge subnet (172.17.0.0/16)"},
		{"255.255.255.255/1", ""},
		{"255.255.255.255/32", ""},
	}
	runTests(t, "SubnetCIDR", SubnetCIDR, tests)
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

func TestEmail(t *testing.T) {
	const invalidMsg = "invalid email address"
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", invalidMsg},
		{".", invalidMsg},
		{"日本語", invalidMsg},
		{"a@abc.com", ""},
		{"A@abc.com", ""},
		{"1@abc.com", ""},
		{"a.B.1.あ@abc.com", ""},
		{"ア@abc.com", ""},
		{"中文@abc.com", ""},
		{"a@abc.com", ""},
		{"a@123.com", ""},
		{"@abc.com", invalidMsg},
	}
	runTests(t, "Email", Email, tests)
}

func TestCIDRsDontOverlap(t *testing.T) {
	cases := []struct {
		a   string
		b   string
		err *regexp.Regexp
	}{
		{
			a:   "192.168.0.0/24",
			b:   "192.168.0.0/24",
			err: regexp.MustCompile("^\"192.168.0.0/24\" and \"192.168.0.0/24\" overlap$"),
		},
		{
			a:   "192.168.0.0/24",
			b:   "192.168.0.3/24",
			err: regexp.MustCompile("^\"192.168.0.0/24\" and \"192.168.0.3/24\" overlap$"),
		},
		{
			a:   "192.168.0.0/30",
			b:   "192.168.0.3/30",
			err: regexp.MustCompile("^\"192.168.0.0/30\" and \"192.168.0.3/30\" overlap$"),
		},
		{
			a: "192.168.0.0/30",
			b: "192.168.0.4/30",
		},
		{
			a:   "0.0.0.0/0",
			b:   "192.168.0.0/24",
			err: regexp.MustCompile("^\"0.0.0.0/0\" and \"192.168.0.0/24\" overlap$"),
		},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%s %s", testCase.a, testCase.b), func(t *testing.T) {
			err := CIDRsDontOverlap(testCase.a, testCase.b)
			if testCase.err == nil {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				assert.Regexp(t, testCase.err, err)
			}
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

func TestOpenSSHPublicKey(t *testing.T) {
	const invalidMsg = "invalid SSH public key"
	const multiLineMsg = "invalid SSH public key (should not contain any newline characters)"
	const privateKeyMsg = "invalid SSH public key (appears to be a private key)"
	tests := []struct {
		in       string
		expected string
	}{
		{"a", invalidMsg},
		{".", invalidMsg},
		{"日本語", invalidMsg},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", ""},
		{"ssh-rsa \t AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com", ""},
		{"\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com\n", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", multiLineMsg},
		{"ssh-rsa\nAAAAB3NzaC1yc2EAAAADAQABAAACAQDxL you@example.com", multiLineMsg},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL\nyou@example.com", multiLineMsg},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL", ""},
		{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCt3BebCHqnSsgpLjo4kVvyfY/z2BS8t27r/7du+O2pb4xYkr7n+KFpbOz523vMTpQ+o1jY4u4TgexglyT9nqasWgLOvo1qjD1agHme8LlTPQSk07rXqOB85Uq5p7ig2zoOejF6qXhcc3n1c7+HkxHrgpBENjLVHOBpzPBIAHkAGaZcl07OCqbsG5yxqEmSGiAlh/IiUVOZgdDMaGjCRFy0wk0mQaGD66DmnFc1H5CzcPjsxr0qO65e7lTGsE930KkO1Vc+RHCVwvhdXs+c2NhJ2/3740Kpes9n1/YullaWZUzlCPDXtRuy6JRbFbvy39JUgHWGWzB3d+3f8oJ/N4qZ cardno:000603633110", ""},
		{"-----BEGIN CERTIFICATE-----abcd-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN RSA PRIVATE KEY-----\nabc\n-----END RSA PRIVATE KEY-----", privateKeyMsg},
	}

	for _, test := range tests {
		err := SSHPublicKey(test.in)
		if (err == nil && test.expected != "") || (err != nil && err.Error() != test.expected) {
			t.Errorf("For %q, expected %q, got %q", test.in, test.expected, err)
		}
	}
}
