// Package validate contains validation utilities for installer types.
package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"golang.org/x/crypto/ssh"
	k8serrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation"
)

var (
	dockerBridgeCIDR = func() *net.IPNet {
		_, cidr, _ := net.ParseCIDR("172.17.0.0/16")
		return cidr
	}()

	// ValidNetworkTypes is a collection of the valid network types.
	ValidNetworkTypes = map[netopv1.NetworkType]bool{
		netopv1.NetworkTypeOpenshiftSDN:  true,
		netopv1.NetworkTypeOVNKubernetes: true,
		netopv1.NetworkTypeCalico:        true,
		netopv1.NetworkTypeKuryr:         true,
	}

	// ValidNetworkTypeValues is a slice filled with the valid network types as
	// strings.
	ValidNetworkTypeValues = func() []string {
		validValues := make([]string, len(ValidNetworkTypes))
		i := 0
		for t := range ValidNetworkTypes {
			validValues[i] = string(t)
			i++
		}
		return validValues
	}()
)

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
func DomainName(v string) error {
	// Trailing dot is OK
	return validateSubdomain(strings.TrimSuffix(v, "."))
}

// Email checks if the given string is a valid email address and returns an error if not.
func Email(v string) error {
	if err := nonEmpty(v); err != nil {
		return err
	}

	invalidError := errors.New("invalid email address")

	split := strings.Split(v, "@")
	if len(split) != 2 {
		return invalidError
	}
	localPart := split[0]
	domain := split[1]

	if nonEmpty(localPart) != nil {
		return invalidError
	}

	// No whitespace allowed in local-part
	if isMatch(`\s`, localPart) {
		return invalidError
	}

	return DomainName(domain)
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

func isMatch(re string, v string) bool {
	return regexp.MustCompile(re).MatchString(v)
}

// ClusterName checks if the given string is a valid name for a cluster and returns an error if not.
func ClusterName(v string) error {
	return validateSubdomain(v)
}

// nonEmpty checks if the given string contains at least one non-whitespace character and returns an error if not.
func nonEmpty(v string) error {
	if utf8.RuneCountInString(strings.TrimSpace(v)) == 0 {
		return errors.New("cannot be empty")
	}
	return nil
}

// SubnetCIDR checks if the given IP net is a valid CIDR for a master nodes or worker nodes subnet and returns an error if not.
func SubnetCIDR(cidr *net.IPNet) error {
	if cidr.IP.To4() == nil {
		return errors.New("must use IPv4")
	}
	if cidr.IP.IsUnspecified() {
		return errors.New("address must be specified")
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

// IPv4 checks if the given string is a valid IP v4 address and returns an error if not.
// Based on net.ParseIP.
func IPv4(v string) error {
	if err := nonEmpty(v); err != nil {
		return err
	}
	if ip := net.ParseIP(v); ip == nil || !strings.Contains(v, ".") {
		return errors.New("invalid IPv4 address")
	}
	return nil
}

// canonicalizeIP ensures that the given IP is in standard form
// and returns an error otherwise.
func canonicalizeIP(ip *net.IP) error {
	if ip.To4() != nil {
		*ip = ip.To4()
		return nil
	}
	if ip.To16() != nil {
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
