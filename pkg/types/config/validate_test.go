package config

import (
	"net"
	"strings"
	"testing"

	"github.com/openshift/installer/pkg/types/config/aws"
	"github.com/openshift/installer/pkg/types/config/libvirt"
)

func TestMissingNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"", "", ""},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 0,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrMissingNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d missing node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestMoreThanOneNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master", "master2"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master", "master2"},
				},
				Worker: Worker{
					NodePools: []string{"worker", "worker2"},
				},
			},
			errs: 2,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrMoreThanOneNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d more-than-one node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestUnmatchedNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master", "master2"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
				NodePools: NodePools{
					{
						Name:  "master",
						Count: 1,
					},
					{
						Name:  "worker",
						Count: 1,
					},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
				NodePools: NodePools{
					{
						Name:  "master",
						Count: 1,
					},
					{
						Name:  "worker",
						Count: 1,
					},
				},
			},
			errs: 0,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrUnmatchedNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d unmatched node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestSharedNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared"},
				},
				Worker: Worker{
					NodePools: []string{"shared"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared"},
				},
				Worker: Worker{
					NodePools: []string{"shared"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared", "shared2"},
				},
				Worker: Worker{
					NodePools: []string{"shared", "shared2"},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared", "shared2"},
				},
				Worker: Worker{
					NodePools: []string{"shared", "shared2", "shared3"},
				},
			},
			errs: 2,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrSharedNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d shared node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestAWSEndpoints(t *testing.T) {
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: Cluster{},
			err:     true,
		},
		{
			cluster: defaultCluster,
			err:     false,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: "foo",
				},
			},
			err: true,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: aws.EndpointsAll,
				},
			},
			err: false,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: aws.EndpointsPrivate,
				},
			},
			err: false,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: aws.EndpointsPublic,
				},
			},
			err: false,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateAWSEndpoints(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestS3BucketNames(t *testing.T) {
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: defaultCluster,
			err:     true,
		},
		{
			cluster: Cluster{},
			err:     true,
		},
		{
			cluster: Cluster{
				Name:       "foo",
				BaseDomain: "example.com",
			},
			err: false,
		},
		{
			cluster: Cluster{
				Name:       ".foo",
				BaseDomain: "example.com",
			},
			err: true,
		},
		{
			cluster: Cluster{
				Name:       "foo",
				BaseDomain: "example.com.",
			},
			err: false,
		},
		{
			cluster: Cluster{
				Name:       "foo",
				BaseDomain: "012345678901234567890123456789012345678901234567890123456789.com",
			},
			err: true,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateS3Bucket(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestValidateLibvirt(t *testing.T) {
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: Cluster{},
			err:     true,
		},
		{
			cluster: defaultCluster,
			err:     true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{},
					Image:   "",
					URI:     "",
				},
				Networking: defaultCluster.Networking,
			},
			err: true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "10.0.1.0/24",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: false,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "10.2.1.0/24",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "x",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "192.168.0.1/24",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: false,
		},
	}

	for i, c := range cases {
		c.cluster.Platform = PlatformLibvirt
		if err := c.cluster.validateLibvirt(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestValidateAWS(t *testing.T) {
	d1 := defaultCluster
	d1.Platform = PlatformAWS
	d2 := d1
	d2.Name = "test"
	d2.BaseDomain = "example.com"
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: Cluster{},
			err:     false,
		},
		{
			cluster: Cluster{
				Platform: PlatformAWS,
			},
			err: true,
		},
		{
			cluster: d1,
			err:     true,
		},
		{
			cluster: d2,
			err:     false,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateAWS(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestValidateOverlapWithPodOrServiceCIDR(t *testing.T) {
	cases := []struct {
		cidr    string
		cluster Cluster
		err     bool
	}{
		{
			cidr:    "192.168.0.1/24",
			cluster: Cluster{},
			err:     true,
		},
		{
			cidr:    "192.168.0.1/24",
			cluster: defaultCluster,
			err:     false,
		},
		{
			cidr:    "10.1.0.0/16",
			cluster: defaultCluster,
			err:     false,
		},
		{
			cidr:    "10.2.0.0/16",
			cluster: defaultCluster,
			err:     true,
		},
		{
			cidr: "10.1.0.0/16",
			cluster: Cluster{
				Networking: Networking{
					PodCIDR:     "10.3.0.0/16",
					ServiceCIDR: "10.4.0.0/16",
				},
			},
			err: false,
		},
		{
			cidr: "10.3.0.0/24",
			cluster: Cluster{
				Networking: Networking{
					PodCIDR:     "10.3.0.0/16",
					ServiceCIDR: "10.4.0.0/16",
				},
			},
			err: true,
		},
		{
			cidr: "0.0.0.0/0",
			cluster: Cluster{
				Networking: Networking{
					PodCIDR:     "10.3.0.0/16",
					ServiceCIDR: "10.4.0.0/16",
				},
			},
			err: true,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateOverlapWithPodOrServiceCIDR(c.cidr, "test"); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

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
	runTests(t, "NonEmpty", validateNonEmpty, tests)
}

func TestInt(t *testing.T) {
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"2 3", invalidIntMsg},
		{"1.1", invalidIntMsg},
		{"abc", invalidIntMsg},
		{"日本語", invalidIntMsg},
		{"1 abc", invalidIntMsg},
		{"日本語2", invalidIntMsg},
		{"0", ""},
		{"1", ""},
		{"999999", ""},
		{"-1", ""},
	}
	runTests(t, "Int", validateInt, tests)
}

func TestIntRange(t *testing.T) {
	tests := []struct {
		in       string
		min      int
		max      int
		expected string
	}{
		{"", 4, 6, emptyMsg},
		{" ", 4, 6, emptyMsg},
		{"2 3", 1, 2, invalidIntMsg},
		{"1.1", 0, 0, invalidIntMsg},
		{"abc", -2, -1, invalidIntMsg},
		{"日本語", 99, 100, invalidIntMsg},
		{"5", 4, 6, ""},
		{"5", 5, 5, ""},
		{"5", 6, 8, "cannot be less than 6"},
		{"5", 6, 4, "cannot be less than 6"},
		{"5", 2, 4, "cannot be greater than 4"},
	}

	for _, test := range tests {
		err := validateIntRange(test.in, test.min, test.max)
		if (err == nil && test.expected != "") || (err != nil && err.Error() != test.expected) {
			t.Errorf("For IntRange(%q, %v, %v), expected %q, got %q", test.in, test.min, test.max, test.expected, err)
		}
	}
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
	runTests(t, "ClusterName", validateClusterName, tests)
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
	runTests(t, "IPv4", validateIPv4, tests)
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
	runTests(t, "SubnetCIDR", validateSubnetCIDR, tests)
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
	runTests(t, "DomainName", ValidateDomainName, tests)
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
	runTests(t, "Email", ValidateEmail, tests)
}

func TestCIDRsDontOverlap(t *testing.T) {
	cases := []struct {
		a   string
		b   string
		err bool
	}{
		{
			a:   "192.168.0.0/24",
			b:   "192.168.0.0/24",
			err: true,
		},
		{
			a:   "192.168.0.0/24",
			b:   "192.168.0.3/24",
			err: true,
		},
		{
			a:   "192.168.0.0/30",
			b:   "192.168.0.3/30",
			err: true,
		},
		{
			a:   "192.168.0.0/30",
			b:   "192.168.0.4/30",
			err: false,
		},
		{
			a:   "0.0.0.0/0",
			b:   "192.168.0.0/24",
			err: true,
		},
	}

	for i, c := range cases {
		if err := validateCIDRsDontOverlap(c.a, c.b); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}
