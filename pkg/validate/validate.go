// Package validate contains validation utilities for installer types.
package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"

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
		sort.Strings(validValues)
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
func ClusterName(v string) error {
	return validateSubdomain(v)
}

// SubnetCIDR checks if the given IP net is a valid CIDR.
func SubnetCIDR(cidr *net.IPNet) error {
	if cidr.IP.To4() == nil {
		return errors.New("must use IPv4")
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
