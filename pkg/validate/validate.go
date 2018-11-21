// Package validate contains validation utilities for installer types.
package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	k8serrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation"
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

// SubnetCIDR checks if the given string is a valid CIDR for a master nodes or worker nodes subnet and returns an error if not.
func SubnetCIDR(v string) error {
	if err := nonEmpty(v); err != nil {
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

// CIDRsDontOverlap ensures two given CIDRs don't overlap
// with one another. CIDR starting IPs are canonicalized
// before being compared.
func CIDRsDontOverlap(acidr, bcidr string) error {
	_, a, err := net.ParseCIDR(acidr)
	if err != nil {
		return fmt.Errorf("invalid CIDR %q: %v", acidr, err)
	}
	if err := canonicalizeIP(&a.IP); err != nil {
		return fmt.Errorf("invalid CIDR %q: %v", acidr, err)
	}
	_, b, err := net.ParseCIDR(bcidr)
	if err != nil {
		return fmt.Errorf("invalid CIDR %q: %v", bcidr, err)
	}
	if err := canonicalizeIP(&b.IP); err != nil {
		return fmt.Errorf("invalid CIDR %q: %v", bcidr, err)
	}
	err = fmt.Errorf("%q and %q overlap", acidr, bcidr)
	// IPs are of different families.
	if len(a.IP) != len(b.IP) {
		return nil
	}
	if a.Contains(b.IP) {
		return err
	}
	if a.Contains(lastIP(b)) {
		return err
	}
	if b.Contains(a.IP) {
		return err
	}
	if b.Contains(lastIP(a)) {
		return err
	}
	return nil
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

// SSHPublicKey checks if the given string is a valid OpenSSH public key
// and returns an error if not.
func SSHPublicKey(v string) error {
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
	if !isMatch(`^[\w-]+$`, keyType) || !isMatch(`^[A-Za-z0-9+\/]+={0,2}$`, keyBase64) {
		return invalidError
	}

	return nil
}
