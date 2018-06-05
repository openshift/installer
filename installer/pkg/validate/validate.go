package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dgrijalva/jwt-go"
)

func isMatch(re string, v string) bool {
	return regexp.MustCompile(re).MatchString(v)
}

// PrefixError wraps an error with a prefix or returns nil if there was no error.
// This is useful for wrapping errors returned by generic error funcs like `NonEmpty` so that the error includes the offending field name.
func PrefixError(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %v", prefix, err)
	}
	return nil
}

// JSONFile validates that the file at the given path is valid JSON.
func JSONFile(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if !json.Valid(b) {
		return fmt.Errorf("file %q contains invalid JSON", path)
	}
	return nil
}

// FileExists validates a file exists at the given path.
func FileExists(path string) error {
	_, err := os.Stat(path)
	return err
}

// License validates that the file at the given path is a valid license.
func License(path string) error {
	licenseBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read license file at %q: %v", path, err)
	}
	c := struct {
		ExpirationDate time.Time `json:"expirationDate"`
		jwt.StandardClaims
	}{}
	if _, _, err := (&jwt.Parser{}).ParseUnverified(string(licenseBytes), &c); err != nil {
		return fmt.Errorf("invalid JWT in license: %v; %q", err, string(licenseBytes))
	}
	if time.Now().After(c.ExpirationDate) {
		return fmt.Errorf("expired license %v", c.ExpirationDate)
	}
	return nil
}

// NonEmpty checks if the given string contains at least one non-whitespace character and returns an error if not.
func NonEmpty(v string) error {
	if utf8.RuneCountInString(strings.TrimSpace(v)) == 0 {
		return errors.New("cannot be empty")
	}
	return nil
}

// Int checks if the given string is a valid integer and returns an error if not.
func Int(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	if _, err := strconv.Atoi(v); err != nil {
		return errors.New("invalid integer")
	}
	return nil
}

// IntRange checks if the given string is a valid integer between `min` and `max` and returns an error if not.
func IntRange(v string, min int, max int) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return Int(v)
	}
	if i < min {
		return fmt.Errorf("cannot be less than %v", min)
	}
	if i > max {
		return fmt.Errorf("cannot be greater than %v", max)
	}
	return nil
}

// IntOdd checks if the given string is a valid integer and that it is odd and returns an error if not.
func IntOdd(v string) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return Int(v)
	}
	if i%2 != 1 {
		return errors.New("must be an odd integer")
	}
	return nil
}

// ClusterName checks if the given string is a valid name for a cluster and returns an error if not.
func ClusterName(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	if length := utf8.RuneCountInString(v); length < 1 || length > 253 {
		return errors.New("must be between 1 and 253 characters")
	}

	if strings.ToLower(v) != v {
		return errors.New("must be lower case")
	}

	if !isMatch("^[a-z0-9-.]*$", v) {
		return errors.New("only lower case alphanumeric [a-z0-9], dashes and dots are allowed")
	}

	isAlphaNum := regexp.MustCompile("^[a-z0-9]$").MatchString

	// If we got this far, we know the string is ASCII and has at least one character
	if !isAlphaNum(v[:1]) || !isAlphaNum(v[len(v)-1:]) {
		return errors.New("must start and end with a lower case alphanumeric character [a-z0-9]")
	}

	for _, segment := range strings.Split(v, ".") {
		// Each segment can have up to 63 characters
		if utf8.RuneCountInString(segment) > 63 {
			return errors.New("no segment between dots can be more than 63 characters")
		}
		if !isAlphaNum(segment[:1]) || !isAlphaNum(segment[len(segment)-1:]) {
			return errors.New("segments between dots must start and end with a lower case alphanumeric character [a-z0-9]")
		}
	}

	return nil
}

// AWSClusterName checks if the given string is a valid name for a cluster on AWS and returns an error if not.
// See AWS docs:
//   http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-using-console-create-stack-parameters.html
//   http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-elasticloadbalancingv2-loadbalancer.html#cfn-elasticloadbalancingv2-loadbalancer-name
func AWSClusterName(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	if length := utf8.RuneCountInString(v); length < 1 || length > 28 {
		return errors.New("must be between 1 and 28 characters")
	}

	if strings.ToLower(v) != v {
		return errors.New("must be lower case")
	}

	if strings.HasPrefix(v, "-") || strings.HasSuffix(v, "-") {
		return errors.New("must not start or end with '-'")
	}

	if !isMatch("^[a-z][-a-z0-9]*$", v) {
		return errors.New("must be a lower case AWS Stack Name: [a-z][-a-z0-9]*")
	}

	return nil
}

// MAC checks if the given string is a valid MAC address and returns an error if not.
// Based on net.ParseMAC.
func MAC(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}
	if _, err := net.ParseMAC(v); err != nil {
		return errors.New("invalid MAC Address")
	}
	return nil
}

// IPv4 checks if the given string is a valid IP v4 address and returns an error if not.
// Based on net.ParseIP.
func IPv4(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}
	if ip := net.ParseIP(v); ip == nil || !strings.Contains(v, ".") {
		return errors.New("invalid IPv4 address")
	}
	return nil
}

// SubnetCIDR checks if the given string is a valid CIDR for a master nodes or worker nodes subnet and returns an error if not.
func SubnetCIDR(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	split := strings.Split(v, "/")

	if len(split) == 1 {
		return errors.New("must provide a CIDR netmask (eg, /24)")
	}

	if len(split) != 2 {
		return errors.New("invalid IPv4 address")
	}

	ip := split[0]

	if err := IPv4(ip); err != nil {
		return errors.New("invalid IPv4 address")
	}

	if mask, err := strconv.Atoi(split[1]); err != nil || mask < 0 || mask > 32 {
		return errors.New("invalid netmask size (must be between 0 and 32)")
	}

	// Catch any invalid CIDRs not caught by the checks above
	if _, _, err := net.ParseCIDR(v); err != nil {
		return errors.New("invalid CIDR")
	}

	if strings.HasPrefix(ip, "172.17.") {
		return errors.New("overlaps with default Docker Bridge subnet (172.17.0.0/16)")
	}

	return nil
}

// AWSSubnetCIDR checks if the given string is a valid CIDR for a master nodes or worker nodes subnet in an AWS VPC and returns an error if not.
func AWSSubnetCIDR(v string) error {
	if err := SubnetCIDR(v); err != nil {
		return err
	}

	_, network, err := net.ParseCIDR(v)
	if err != nil {
		return errors.New("invalid CIDR")
	}
	if mask, _ := network.Mask.Size(); mask < 16 || mask > 28 {
		return errors.New("AWS subnets must be between /16 and /28")
	}

	return nil
}

// DomainName checks if the given string is a valid domain name and returns an error if not.
func DomainName(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	split := strings.Split(v, ".")
	for i, segment := range split {
		// Trailing dot is OK
		if len(segment) == 0 && i == len(split)-1 {
			continue
		}
		if !isMatch("^[a-zA-Z0-9-]{1,63}$", segment) {
			return errors.New("invalid domain name")
		}
	}
	return nil
}

// Host checks if the given string is either a valid IPv4 address or a valid domain name and returns an error if not.
func Host(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	// Either a valid IP address or domain name
	if IPv4(v) != nil && DomainName(v) != nil {
		return errors.New("invalid host (must be a domain name or IP address)")
	}
	return nil
}

// Port checks if the given string is a valid port number and returns an error if not.
func Port(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}
	if IntRange(v, 1, 65535) != nil {
		return errors.New("invalid port number")
	}
	return nil
}

// HostPort checks if the given string is valid <host>:<port> format and returns an error if not.
func HostPort(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	split := strings.Split(v, ":")
	if len(split) != 2 {
		return errors.New("must use <host>:<port> format")
	}
	if err := Host(split[0]); err != nil {
		return err
	}
	return Port(split[1])
}

// Email checks if the given string is a valid email address and returns an error if not.
func Email(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	invalidError := errors.New("invalid email address")

	split := strings.Split(v, "@")
	if len(split) != 2 {
		return invalidError
	}
	localPart := split[0]
	domain := split[1]

	if NonEmpty(localPart) != nil {
		return invalidError
	}

	// No whitespace allowed in local-part
	if isMatch(`\s`, localPart) {
		return invalidError
	}

	return DomainName(domain)
}

const base64RegExp = `[A-Za-z0-9+\/]+={0,2}`

// Certificate checks if the given string is a valid certificate in PEM format and returns an error if not.
// Ignores leading and trailing whitespace.
func Certificate(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	trimmed := strings.TrimSpace(v)

	// Don't let users hang themselves
	if isMatch(`-BEGIN [\w-]+ PRIVATE KEY-`, trimmed) {
		return errors.New("invalid certificate (appears to be a private key)")
	}

	if !isMatch("(?s:^-----BEGIN CERTIFICATE-----\n"+base64RegExp+"\n-----END CERTIFICATE-----$)", trimmed) {
		return errors.New("invalid certificate")
	}
	return nil
}

// PrivateKey checks if the given string is a valid private key in PEM format and returns an error if not.
// Ignores leading and trailing whitespace.
func PrivateKey(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	trimmed := strings.TrimSpace(v)

	if !isMatch("(?s:^-----BEGIN [A-Z]{2,10} PRIVATE KEY-----\n"+base64RegExp+"\n-----END [A-Z]{2,10} PRIVATE KEY-----$)", trimmed) {
		return errors.New("invalid private key")
	}
	return nil
}

// OpenSSHPublicKey checks if the given string is a valid OpenSSH public key and returns an error if not.
// Ignores leading and trailing whitespace.
func OpenSSHPublicKey(v string) error {
	if err := NonEmpty(v); err != nil {
		return err
	}

	trimmed := strings.TrimSpace(v)

	// Don't let users hang themselves
	if isMatch(`-BEGIN [\w-]+ PRIVATE KEY-`, trimmed) {
		return errors.New("invalid SSH public key (appears to be a private key)")
	}

	if strings.Contains(trimmed, "\n") {
		return errors.New("invalid SSH public key (should not contain any newline characters)")
	}

	invalidError := errors.New("invalid SSH public key")

	keyParts := regexp.MustCompile(`\s+`).Split(trimmed, -1)
	if len(keyParts) < 2 {
		return invalidError
	}

	keyType := keyParts[0]
	keyBase64 := keyParts[1]
	if !isMatch(`^[\w-]+$`, keyType) || !isMatch("^"+base64RegExp+"$", keyBase64) {
		return invalidError
	}

	return nil
}

// CIDRsOverlap checks whether two given CIDRs overlap
// with one another. CIDR starting IPs should be canonicalized
// before being compared.
func CIDRsOverlap(a, b *net.IPNet) bool {
	// IPs are of different families.
	if len(a.IP) != len(b.IP) {
		return false
	}
	if a.Contains(b.IP) {
		return true
	}
	if a.Contains(lastIP(b)) {
		return true
	}
	if b.Contains(a.IP) {
		return true
	}
	if b.Contains(lastIP(a)) {
		return true
	}
	return false
}

// CanonicalizeIP ensures that the given IP is in standard form
// and returns an error otherwise.
func CanonicalizeIP(ip *net.IP) error {
	if ip.To4 != nil {
		*ip = ip.To4()
		return nil
	}
	if ip.To16 != nil {
		*ip = ip.To16()
		return nil
	}
	return fmt.Errorf("IP %q is of unknown type", ip)
}

func lastIP(cidr *net.IPNet) net.IP {
	var last net.IP
	for i := 0; i < len(cidr.IP); i++ {
		last = append(last, cidr.IP[i]|^cidr.Mask[i])
	}
	return last
}
