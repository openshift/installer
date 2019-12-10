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
	"strings"

	"golang.org/x/crypto/ssh"
	k8serrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation"
)

var (
	dockerBridgeCIDR = func() *net.IPNet {
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

// ClusterName checks if the given string is a valid name for a cluster and returns an error if not.
// The max length of the DNS label is `DNS1123LabelMaxLength + 9` because the public DNS zones have records
// `api.clustername`, `*.apps.clustername`, and *.apps is rendered as the nine-character \052.apps in DNS records.
func ClusterName(v string) error {
	maxlen := validation.DNS1123LabelMaxLength - 9
	if len(v) > maxlen {
		return errors.New(validation.MaxLenError(maxlen))
	}
	return validateSubdomain(v)
}

// SubnetCIDR checks if the given IP net is a valid CIDR.
func SubnetCIDR(cidr *net.IPNet, allowIPv6, requireIPv6 bool) error {
	if allowIPv6 == false && cidr.IP.To4() == nil {
		return errors.New("must use IPv4")
	}
	if requireIPv6 == true && cidr.IP.To4() != nil {
		return errors.New("must use IPv6")
	}
	if cidr.IP.IsUnspecified() {
		return errors.New("address must be specified")
	}
	nip := cidr.IP.Mask(cidr.Mask)
	if nip.String() != cidr.IP.String() {
		return fmt.Errorf("invalid network address. got %s, expecting %s", cidr.String(), (&net.IPNet{IP: nip, Mask: cidr.Mask}).String())
	}
	if DoCIDRsOverlap(cidr, dockerBridgeCIDR) {
		return fmt.Errorf("overlaps with default Docker Bridge subnet (%v)", cidr.String())
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

// MAC validates that a value is a valid mac address
func MAC(addr string) error {
	_, err := net.ParseMAC(addr)
	return err
}
