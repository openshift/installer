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
	const charsMsg = "only lower case alphanumeric [a-z0-9], dashes and dots are allowed"
	const lengthMsg = "must be between 1 and 253 characters"
	const segmentLengthMsg = "no segment between dots can be more than 63 characters"
	const startEndCharMsg = "must start and end with a lower case alphanumeric character [a-z0-9]"
	const segmentStartEndCharMsg = "segments between dots must start and end with a lower case alphanumeric character [a-z0-9]"

	maxSizeName := strings.Repeat("123456789.", 25) + "123"
	maxSizeSegment := strings.Repeat("1234567890", 6) + "123"

	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", ""},
		{"A", caseMsg},
		{"abc D", caseMsg},
		{"1", ""},
		{".", startEndCharMsg},
		{"a.", startEndCharMsg},
		{".a", startEndCharMsg},
		{"a.a", ""},
		{"-a", startEndCharMsg},
		{"a-", startEndCharMsg},
		{"a.-a", segmentStartEndCharMsg},
		{"a-.a", segmentStartEndCharMsg},
		{"a%a", charsMsg},
		{"日本語", charsMsg},
		{"a日本語a", charsMsg},
		{maxSizeName, ""},
		{maxSizeName + "a", lengthMsg},
		{maxSizeSegment + ".abc", ""},
		{maxSizeSegment + "a.abc", segmentLengthMsg},
	}
	runTests(t, "ClusterName", ClusterName, tests)
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
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", ""},
		{".", invalidDomainMsg},
		{"日本語", invalidDomainMsg},
		{"日本語.com", invalidDomainMsg},
		{"abc.日本語.com", invalidDomainMsg},
		{"a日本語a.com", invalidDomainMsg},
		{"abc", ""},
		{"ABC", ""},
		{"ABC123", ""},
		{"ABC123.COM123", ""},
		{"1", ""},
		{"0.0", ""},
		{"1.2.3.4", ""},
		{"1.2.3.4.", ""},
		{"abc.", ""},
		{"abc.com", ""},
		{"abc.com.", ""},
		{"a.b.c.d.e.f", ""},
		{".abc", invalidDomainMsg},
		{".abc.com", invalidDomainMsg},
		{".abc.com", invalidDomainMsg},
	}
	runTests(t, "DomainName", DomainName, tests)
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
		{"a@ABC.com", ""},
		{"a@123.com", ""},
		{"a@日本語.com", invalidDomainMsg},
		{"a@.com", invalidDomainMsg},
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
