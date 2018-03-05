package validate

import (
	"strings"
	"testing"
)

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
	runTests(t, "NonEmpty", NonEmpty, tests)
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
	runTests(t, "Int", Int, tests)
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
		err := IntRange(test.in, test.min, test.max)
		if (err == nil && test.expected != "") || (err != nil && err.Error() != test.expected) {
			t.Errorf("For IntRange(%q, %v, %v), expected %q, got %q", test.in, test.min, test.max, test.expected, err)
		}
	}
}

func TestIntOdd(t *testing.T) {
	notOddMsg := "must be an odd integer"
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"0", notOddMsg},
		{"1", ""},
		{"2", notOddMsg},
		{"99", ""},
		{"100", notOddMsg},
		{"abc", invalidIntMsg},
		{"1 abc", invalidIntMsg},
		{"日本語", invalidIntMsg},
	}
	runTests(t, "IntOdd", IntOdd, tests)
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

func TestAWSClusterName(t *testing.T) {
	const charsMsg = "must be a lower case AWS Stack Name: [a-z][-a-z0-9]*"
	const lengthMsg = "must be between 1 and 28 characters"
	const hyphenMsg = "must not start or end with '-'"

	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", ""},
		{"A", caseMsg},
		{"abc D", caseMsg},
		{"1", charsMsg},
		{".", charsMsg},
		{"a.", charsMsg},
		{".a", charsMsg},
		{"a.a", charsMsg},
		{"a%a", charsMsg},
		{"a-a", ""},
		{"-abc", hyphenMsg},
		{"abc-", hyphenMsg},
		{"日本語", charsMsg},
		{"a日本語a", charsMsg},
		{"a234567890123456789012345678", ""},
		{"12345678901234567890123456789", lengthMsg},
		{"A2345678901234567890123456789", lengthMsg},
	}
	runTests(t, "AWSClusterName", AWSClusterName, tests)
}

func TestMAC(t *testing.T) {
	const invalidMsg = "invalid MAC Address"
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"abc", invalidMsg},
		{"12:34:45:78:9A:BC", ""},
		{"12-34-45-78-9A-BC", ""},
		{"12:34:45:78:9a:bc", ""},
		{"12:34:45:78:9X:YZ", invalidMsg},
		{"12.34.45.78.9A.BC", invalidMsg},
	}
	runTests(t, "MAC", MAC, tests)
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

func TestAWSsubnetCIDR(t *testing.T) {
	const awsNetmaskSizeMsg = "AWS subnets must be between /16 and /28"

	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"/20", invalidIPMsg},
		{"1.2.3.4", noCIDRNetmaskMsg},
		{"1.2.3.4/15", awsNetmaskSizeMsg},
		{"1.2.3.4/16", ""},
		{"1.2.3.4/28", ""},
		{"1.2.3.4/29", awsNetmaskSizeMsg},
	}
	runTests(t, "AWSSubnetCIDR", AWSSubnetCIDR, tests)
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

func TestHost(t *testing.T) {
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", ""},
		{".", invalidHostMsg},
		{"日本語", invalidHostMsg},
		{"日本語.com", invalidHostMsg},
		{"abc.日本語.com", invalidHostMsg},
		{"a日本語a.com", invalidHostMsg},
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
		{".abc", invalidHostMsg},
		{".abc.com", invalidHostMsg},
		{".abc.com", invalidHostMsg},
	}
	runTests(t, "Host", Host, tests)
}

func TestPort(t *testing.T) {
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", invalidPortMsg},
		{".", invalidPortMsg},
		{"日本語", invalidPortMsg},
		{"0", invalidPortMsg},
		{"1", ""},
		{"123", ""},
		{"12345", ""},
		{"65535", ""},
		{"65536", invalidPortMsg},
	}
	runTests(t, "Port", Port, tests)
}

func TestHostPort(t *testing.T) {
	const invalidHostPortMsg = "must use <host>:<port> format"
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{".", invalidHostPortMsg},
		{"日本語", invalidHostPortMsg},
		{"abc.com", invalidHostPortMsg},
		{"abc.com:0", invalidPortMsg},
		{"abc.com:1", ""},
		{"abc.com:65535", ""},
		{"abc.com:65536", invalidPortMsg},
		{"abc.com:abc", invalidPortMsg},
		{"1.2.3.4:1234", ""},
		{"1.2.3.4:abc", invalidPortMsg},
		{"日本語:1234", invalidHostMsg},
	}
	runTests(t, "HostPort", HostPort, tests)
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

func TestCertificate(t *testing.T) {
	const invalidMsg = "invalid certificate"
	const privateKeyMsg = "invalid certificate (appears to be a private key)"
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", invalidMsg},
		{".", invalidMsg},
		{"日本語", invalidMsg},
		{"-----BEGIN CERTIFICATE-----\na\n-----END CERTIFICATE-----", ""},
		{"-----BEGIN CERTIFICATE-----\nabc\n-----END CERTIFICATE-----", ""},
		{"-----BEGIN CERTIFICATE-----\nabc=\n-----END CERTIFICATE-----", ""},
		{"-----BEGIN CERTIFICATE-----\nabc==\n-----END CERTIFICATE-----", ""},
		{"-----BEGIN CERTIFICATE-----\nabc===\n-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN CERTIFICATE-----\na%a\n-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN CERTIFICATE-----\n\nab\n-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN CERTIFICATE-----\nab\n\n-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN CERTIFICATE-----\na\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\na\n-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN RSA PRIVATE KEY-----\nabc\n-----END RSA PRIVATE KEY-----", privateKeyMsg},
	}
	runTests(t, "Certificate", Certificate, tests)
}

func TestPrivateKey(t *testing.T) {
	const invalidMsg = "invalid private key"
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
		{"a", invalidMsg},
		{".", invalidMsg},
		{"日本語", invalidMsg},
		{"-----BEGIN RSA PRIVATE KEY-----\na\n-----END RSA PRIVATE KEY-----", ""},
		{"-----BEGIN RSA PRIVATE KEY-----\nabc\n-----END RSA PRIVATE KEY-----", ""},
		{"-----BEGIN RSA PRIVATE KEY-----\nabc=\n-----END RSA PRIVATE KEY-----", ""},
		{"-----BEGIN RSA PRIVATE KEY-----\nabc==\n-----END RSA PRIVATE KEY-----", ""},
		{"-----BEGIN RSA PRIVATE KEY-----\nabc===\n-----END RSA PRIVATE KEY-----", invalidMsg},
		{"-----BEGIN EC PRIVATE KEY-----\nabc\n-----END EC PRIVATE KEY-----", ""},
		{"-----BEGIN RSA PRIVATE KEY-----\na%a\n-----END RSA PRIVATE KEY-----", invalidMsg},
		{"-----BEGIN RSA PRIVATE KEY-----\n\nab\n-----END RSA PRIVATE KEY-----", invalidMsg},
		{"-----BEGIN RSA PRIVATE KEY-----\nab\n\n-----END RSA PRIVATE KEY-----", invalidMsg},
		{"-----BEGIN RSA PRIVATE KEY-----\na\n-----END RSA PRIVATE KEY-----\n-----BEGIN CERTIFICATE-----\na\n-----END CERTIFICATE-----", invalidMsg},
		{"-----BEGIN CERTIFICATE-----\na\n-----END CERTIFICATE-----", invalidMsg},
	}
	runTests(t, "PrivateKey", PrivateKey, tests)
}

func TestOpenSSHPublicKey(t *testing.T) {
	const invalidMsg = "invalid SSH public key"
	const multiLineMsg = "invalid SSH public key (should not contain any newline characters)"
	const privateKeyMsg = "invalid SSH public key (appears to be a private key)"
	tests := []test{
		{"", emptyMsg},
		{" ", emptyMsg},
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
	runTests(t, "OpenSSHPublicKey", OpenSSHPublicKey, tests)
}
