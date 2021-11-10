// Package validate contains validation utilities for installer types.
package validate

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	uuid "github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	k8serrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation"
)

var (
	// DockerBridgeCIDR is the network range that is used by default network for docker.
	DockerBridgeCIDR = func() *net.IPNet {
		_, cidr, _ := net.ParseCIDR("172.17.0.0/16")
		return cidr
	}()
)

// CABundle checks if the given string contains valid certificate(s) and returns an error if not.
func CABundle(v string) error {
	rest := []byte(v)
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			return fmt.Errorf("invalid block")
		}
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
		if len(rest) == 0 {
			break
		}
	}
	return nil
}
func validateSubdomain(v string) error {
	validationMessages := validation.IsDNS1123Subdomain(v)
	if len(validationMessages) == 0 {
		return nil
	}

	errs := make([]error, len(validationMessages))
	for i, m := range validationMessages {
		errs[i] = errors.New(m)
	}
	return k8serrors.NewAggregate(errs)
}

// DomainName checks if the given string is a valid domain name and returns an error if not.
func DomainName(v string, acceptTrailingDot bool) error {
	if acceptTrailingDot {
		v = strings.TrimSuffix(v, ".")
	}
	return validateSubdomain(v)
}

// NoProxyDomainName checks if the given string is a valid proxy noProxy domain name
// and returns an error if not. Example valid noProxy domains are ".foo.com", "bar.com",
// "bar.com." but not "*.foo.com".
func NoProxyDomainName(v string) error {
	v = strings.TrimSuffix(strings.TrimPrefix(v, "."), ".")
	return validateSubdomain(v)
}

type imagePullSecret struct {
	Auths map[string]map[string]interface{} `json:"auths"`
}

// ImagePullSecret checks if the given string is a valid image pull secret and returns an error if not.
func ImagePullSecret(secret string) error {
	var s imagePullSecret
	err := json.Unmarshal([]byte(secret), &s)
	if err != nil {
		return err
	}
	if len(s.Auths) == 0 {
		return fmt.Errorf("auths required")
	}
	errs := []error{}
	for d, a := range s.Auths {
		_, authPresent := a["auth"]
		_, credsStorePresnet := a["credsStore"]
		if !authPresent && !credsStorePresnet {
			errs = append(errs, fmt.Errorf("%q requires either auth or credsStore", d))
		}
	}
	return k8serrors.NewAggregate(errs)
}

// ClusterName1035 checks the provided cluster name matches RFC1035 naming requirements.
// Some platform resource names must comply with RFC1035 "[a-z]([-a-z0-9]*[a-z0-9])?". They
// are based on the InfraID, which is a truncated version of the cluster name where all non-
// alphanumeric characters "[^A-Za-z0-9-]" have been replaced with dashes "-". As a result,
// if we first verify the name starts with a lower-case letter "^[a-z]" then we can rely on
// the ClusterName function to confirm compliance with the rest. The resulting name will
// therefore match RFC1035 with the exception of possible periods ".", which will be
// translated into dashes "-" in the InfraID before being used to create cloud resources.
func ClusterName1035(v string) error {
	re := regexp.MustCompile("^[a-z]")
	if !re.MatchString(v) {
		return errors.New("cluster name must begin with a lower-case letter")
	}
	return ClusterName(v)
}

// GCPClusterName checks if the provided cluster name has words similar to the word 'google'
// since resources with that name are not allowed in GCP.
func GCPClusterName(v string) error {
	reStartsWith := regexp.MustCompile("^goog")
	reContains := regexp.MustCompile(".*g[o0]{2}gle.*")
	if reStartsWith.MatchString(v) || reContains.MatchString(v) {
		return errors.New("cluster name must not start with \"goog\" or contain variations of \"google\"")
	}
	return nil
}

// ClusterNameMaxLength validates if the string provided length is
// greater than maxlen argument.
func ClusterNameMaxLength(v string, maxlen int) error {
	if len(v) > maxlen {
		return errors.New(validation.MaxLenError(maxlen))
	}
	return nil
}

// ClusterName checks if the given string is a valid name for a cluster and returns an error if not.
// The max length of the DNS label is `DNS1123LabelMaxLength + 9` because the public DNS zones have records
// `api.clustername`, `*.apps.clustername`, and *.apps is rendered as the nine-character \052.apps in DNS records.
func ClusterName(v string) error {
	const maxlen = validation.DNS1123LabelMaxLength - 9
	err := ClusterNameMaxLength(v, maxlen)
	if err != nil {
		return err
	}
	return validateSubdomain(v)
}

// SubnetCIDR checks if the given IP net is a valid CIDR.
func SubnetCIDR(cidr *net.IPNet) error {
	if cidr.IP.IsUnspecified() {
		return errors.New("address must be specified")
	}
	nip := cidr.IP.Mask(cidr.Mask)
	if nip.String() != cidr.IP.String() {
		return fmt.Errorf("invalid network address. got %s, expecting %s", cidr.String(), (&net.IPNet{IP: nip, Mask: cidr.Mask}).String())
	}
	return nil
}

// ServiceSubnetCIDR checks if the given IP net is a valid CIDR for the Kubernetes service network
func ServiceSubnetCIDR(cidr *net.IPNet) error {
	if cidr.IP.IsUnspecified() {
		return errors.New("address must be specified")
	}
	nip := cidr.IP.Mask(cidr.Mask)
	if nip.String() != cidr.IP.String() {
		return fmt.Errorf("invalid network address. got %s, expecting %s", cidr.String(), (&net.IPNet{IP: nip, Mask: cidr.Mask}).String())
	}
	maskLen, addrLen := cidr.Mask.Size()
	if addrLen == 32 && maskLen < 12 {
		return fmt.Errorf("subnet size for IPv4 service network must be /12 or greater (/16 is recommended)")
	} else if addrLen == 128 && maskLen < 108 {
		// Kubernetes allows any length greater than 108 (and so do we, for
		// backward compat), but for various reasons there is no point in
		// using any value other than 112.
		return fmt.Errorf("subnet size for IPv6 service network should be /112")
	}
	return nil
}

// DoCIDRsOverlap returns true if one of the CIDRs is a subset of the other.
func DoCIDRsOverlap(acidr, bcidr *net.IPNet) bool {
	return acidr.Contains(bcidr.IP) || bcidr.Contains(acidr.IP)
}

// SSHPublicKey checks if the given string is a valid SSH public key
// and returns an error if not.
func SSHPublicKey(v string) error {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(v))
	return err
}

// URI validates if the URI is a valid absolute URI.
func URI(uri string) error {
	parsed, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if !parsed.IsAbs() {
		return fmt.Errorf("invalid URI %q (no scheme)", uri)
	}
	return nil
}

// URIWithProtocol validates that the URI specifies a certain
// protocol scheme (e.g. "https")
func URIWithProtocol(uri string, protocol string) error {
	parsed, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if parsed.Scheme != protocol {
		return fmt.Errorf("must use %s protocol", protocol)
	}
	return nil
}

// IP validates if a string is a valid IP.
func IP(ip string) error {
	addr := net.ParseIP(ip)
	if addr == nil {
		return fmt.Errorf("%q is not a valid IP", ip)
	}
	return nil
}

// MAC validates that a value is a valid unicast EUI-48 MAC address
func MAC(addr string) error {
	hwAddr, err := net.ParseMAC(addr)
	if err != nil {
		return err
	}

	// net.ParseMAC checks for any valid mac, including 20-octet infiniband
	// MAC's. Let's make sure we have an EUI-48 MAC, consisting of 6 octets
	if len(hwAddr) != 6 {
		return fmt.Errorf("invalid MAC address")
	}

	// We also need to check that the MAC is a valid unicast address. A multicast
	// address is an address where the least significant bit of the most significant
	// byte is 1.
	//
	//      Example 1: Multicast MAC
	//      ------------------------
	//      7D:CE:E3:29:35:6F
	//       ^--> most significant byte
	//
	//      0x7D is 0b11111101
	//                       ^--> this is a multicast MAC
	//
	//      Example 2: Unicast MAC
	//      ----------------------
	//      7A:CE:E3:29:35:6F
	//       ^--> most significant byte
	//
	//      0x7A is 0b11111010
	//                       ^--> this is a unicast MAC
	if hwAddr[0]&1 == 1 {
		return fmt.Errorf("expected unicast mac address, found multicast")
	}

	return nil
}

// UUID validates that a uuid is non-empty and a valid uuid.
func UUID(val string) error {
	_, err := uuid.Parse(val)
	return err
}

// Host validates that a given string is a valid URI host.
func Host(v string) error {
	proxyIP := net.ParseIP(v)
	if proxyIP != nil {
		return nil
	}
	return validateSubdomain(v)
}
