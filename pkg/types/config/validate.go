package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/openshift/installer/pkg/types/config/aws"

	"github.com/coreos/tectonic-config/config/tectonic-network"
	"github.com/sirupsen/logrus"
)

const (
	maxS3BucketNameLength = 63
)

// ErrUnmatchedNodePool is returned when a nodePool was specified but not found in the nodePools list.
type ErrUnmatchedNodePool struct {
	name string
}

// ErrUnmatchedNodePool implements the error interface.
func (e *ErrUnmatchedNodePool) Error() string {
	return fmt.Sprintf("no node pool named %q was found", e.name)
}

// ErrMissingNodePool is returned when a field that requires a nodePool does not specify one.
type ErrMissingNodePool struct {
	field string
}

// ErrMissingNodePool implements the error interface.
func (e *ErrMissingNodePool) Error() string {
	return fmt.Sprintf("the %s field requires at least one node pool to be specified", e.field)
}

// ErrMoreThanOneNodePool is returned when a field specifies more than one node pool.
type ErrMoreThanOneNodePool struct {
	field string
}

// ErrMoreThanOneNodePool implements the error interface.
func (e *ErrMoreThanOneNodePool) Error() string {
	return fmt.Sprintf("the %s field specifies more than one node pool; this is not currently allowed", e.field)
}

// ErrSharedNodePool is returned when two or more fields are defined to use the same nodePool.
type ErrSharedNodePool struct {
	name   string
	fields []string
}

// ErrSharedNodePool implements the error interface.
func (e *ErrSharedNodePool) Error() string {
	return fmt.Sprintf("node pools cannot be shared, but %q is used by %s", e.name, strings.Join(e.fields, ", "))
}

// ErrInvalidIgnConfig is returned when a invalid ign config is given.
type ErrInvalidIgnConfig struct {
	filePath string
	rpt      string
}

// ErrInvalidIgnConfig implements the error interface.
func (e *ErrInvalidIgnConfig) Error() string {
	return fmt.Sprintf("failed to parse ignition file %s: %s", e.filePath, e.rpt)
}

// Validate ensures that the Cluster is semantically correct and returns an error if not.
func (c *Cluster) Validate() []error {
	var errs []error
	errs = append(errs, c.validateNodePools()...)
	errs = append(errs, c.validateNetworking()...)
	errs = append(errs, c.validateAWS()...)
	errs = append(errs, c.validateOpenStack()...)
	errs = append(errs, c.validatePullSecret()...)
	errs = append(errs, c.validateLibvirt()...)
	if err := prefixError("cluster name", validateClusterName(c.Name)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("base domain", ValidateDomainName(c.BaseDomain)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("admin password", validateNonEmpty(c.Admin.Password)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("admin email", ValidateEmail(c.Admin.Email)); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// validateAWS validates all fields specific to AWS.
func (c *Cluster) validateAWS() []error {
	var errs []error
	if c.Platform != PlatformAWS {
		return errs
	}
	if err := c.validateAWSEndpoints(); err != nil {
		errs = append(errs, err)
	}
	if err := c.validateS3Bucket(); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("aws vpcCIDRBlock", validateSubnetCIDR(c.AWS.VPCCIDRBlock)); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, c.validateOverlapWithPodOrServiceCIDR(c.AWS.VPCCIDRBlock, "aws vpcCIDRBlock")...)
	if err := prefixError("aws profile", validateNonEmpty(c.AWS.Profile)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("aws region", validateNonEmpty(c.AWS.Region)); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// validateOpenStack validates all fields specific to OpenStack.
func (c *Cluster) validateOpenStack() []error {
	var errs []error
	if c.Platform != PlatformOpenStack {
		return errs
	}
	return errs
}

// validateOverlapWithPodOrServiceCIDR ensures that the given CIDR does not
// overlap with the pod or service CIDRs of the cluster config.
func (c *Cluster) validateOverlapWithPodOrServiceCIDR(cidr, name string) []error {
	var errs []error
	if err := prefixError(fmt.Sprintf("%s and podCIDR", name), validateCIDRsDontOverlap(cidr, c.Networking.PodCIDR)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError(fmt.Sprintf("%s and serviceCIDR", name), validateCIDRsDontOverlap(cidr, c.Networking.ServiceCIDR)); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// validateLibvirt validates all fields specific to libvirt.
func (c *Cluster) validateLibvirt() []error {
	var errs []error
	if c.Platform != PlatformLibvirt {
		return errs
	}
	if err := prefixError("libvirt network ipRange", validateSubnetCIDR(c.Libvirt.Network.IPRange)); err != nil {
		errs = append(errs, err)
	}
	if len(c.Libvirt.MasterIPs) > 0 {
		if len(c.Libvirt.MasterIPs) != c.NodeCount(c.Master.NodePools) {
			errs = append(errs, fmt.Errorf("length of masterIPs does't match master count"))
		}
		for i, ip := range c.Libvirt.MasterIPs {
			if err := prefixError(fmt.Sprintf("libvirt masterIPs[%d] %q", i, ip), validateIPv4(ip)); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if err := prefixError("libvirt uri", validateNonEmpty(c.Libvirt.URI)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("libvirt network name", validateNonEmpty(c.Libvirt.Network.Name)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("libvirt network ifName", validateNonEmpty(c.Libvirt.Network.IfName)); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, c.validateOverlapWithPodOrServiceCIDR(c.Libvirt.Network.IPRange, "libvirt ipRange")...)
	return errs
}

func (c *Cluster) validateNetworking() []error {
	var errs []error
	// https://en.wikipedia.org/wiki/Maximum_transmission_unit#MTUs_for_common_media
	if err := prefixError("mtu", validateIntRange(c.Networking.MTU, 68, 64*1024)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("podCIDR", validateSubnetCIDR(c.Networking.PodCIDR)); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("serviceCIDR", validateSubnetCIDR(c.Networking.ServiceCIDR)); err != nil {
		errs = append(errs, err)
	}
	if err := c.validateNetworkType(); err != nil {
		errs = append(errs, err)
	}
	if err := prefixError("pod and service CIDRs", validateCIDRsDontOverlap(c.Networking.PodCIDR, c.Networking.ServiceCIDR)); err != nil {
		errs = append(errs, err)
	}
	return errs
}

func (c *Cluster) validateNetworkType() error {
	switch c.Networking.Type {
	case tectonicnetwork.NetworkNone:
		fallthrough
	case tectonicnetwork.NetworkCanal:
		fallthrough
	case tectonicnetwork.NetworkFlannel:
		fallthrough
	case tectonicnetwork.NetworkCalicoIPIP:
		return nil
	default:
		return fmt.Errorf("invalid network type %q", c.Networking.Type)
	}
}

// ValidateAndLog performs cluster configuration validation using `Validate`
// but rather than return a slice of errors, it logs any errors and returns
// a single error for convenience.
func (c *Cluster) ValidateAndLog() error {
	if errs := c.Validate(); len(errs) != 0 {
		logrus.Errorf("Found %d error(s) in the cluster definition:", len(errs))
		for i, err := range errs {
			logrus.Errorf("  Error %d: %v", i+1, err)
		}
		return fmt.Errorf("found %d cluster definition error(s)", len(errs))
	}
	return nil
}

// validateAWSEndpoints ensures that the value of the endpoints field is one of:
// 'all', 'public', or 'private'.
func (c *Cluster) validateAWSEndpoints() error {
	switch c.AWS.Endpoints {
	case aws.EndpointsAll:
		fallthrough
	case aws.EndpointsPrivate:
		fallthrough
	case aws.EndpointsPublic:
		return nil
	default:
		return fmt.Errorf("invalid AWS endpoints %q", c.AWS.Endpoints)
	}
}

// validateS3Bucket does some basic validation to ensure that the S3 bucket
// matches the S3 bucket naming rules. Not all rules are checked
// because Tectonic controls the generation of S3 bucket names, creating
// buckets of the form: <cluster-name>.<domain-name>
// If domain-name contains a trailing dot, it's removed from the bucket name.
func (c *Cluster) validateS3Bucket() error {
	bucket := fmt.Sprintf("%s.%s", c.Name, strings.TrimRight(c.BaseDomain, "."))
	if len(bucket) > maxS3BucketNameLength {
		return fmt.Errorf("the S3 bucket name %q, generated from the cluster name and base domain, is too long; S3 bucket names must be less than 63 characters; please choose a shorter cluster name or base domain", bucket)
	}
	if !regexp.MustCompile("^[a-z0-9][a-z0-9-.]{1,61}[a-z0-9]$").MatchString(bucket) {
		return errors.New("invalid characters in S3 bucket name")
	}
	return nil
}

func (c *Cluster) validatePullSecret() []error {
	var errs []error
	if err := ValidateJSON([]byte(c.PullSecret)); err != nil {
		errs = append(errs, prefixError("pull secret", err))
	}
	return errs
}

func (c *Cluster) validateNodePools() []error {
	var errs []error
	n := c.NodePools.Map()
	fields := []struct {
		pools []string
		field string
	}{
		{pools: c.Master.NodePools, field: "master"},
		{pools: c.Worker.NodePools, field: "worker"},
	}
	for _, f := range fields {
		var found bool
		for _, p := range f.pools {
			if p == "" {
				continue
			}
			found = true
			if _, ok := n[p]; !ok {
				errs = append(errs, &ErrUnmatchedNodePool{p})
			}
		}
		if !found {
			errs = append(errs, &ErrMissingNodePool{f.field})
		}
		if len(f.pools) > 1 {
			errs = append(errs, &ErrMoreThanOneNodePool{f.field})
		}
	}

	errs = append(errs, c.validateNoSharedNodePools()...)

	return errs
}

func (c *Cluster) validateNoSharedNodePools() []error {
	var errs []error
	fields := make(map[string]map[string]struct{})
	for i := range c.Master.NodePools {
		if c.Master.NodePools[i] != "" {
			for j := range c.Worker.NodePools {
				if c.Master.NodePools[i] == c.Worker.NodePools[j] {
					if fields[c.Master.NodePools[i]] == nil {
						fields[c.Master.NodePools[i]] = make(map[string]struct{})
					}
					fields[c.Master.NodePools[i]]["master"] = struct{}{}
					fields[c.Master.NodePools[i]]["worker"] = struct{}{}
				}
			}
		}
	}
	for k, v := range fields {
		err := &ErrSharedNodePool{name: k}
		for f := range v {
			err.fields = append(err.fields, f)
		}
		errs = append(errs, err)
	}
	return errs
}

// ValidateDomainName checks if the given string is a valid domain name and returns an error if not.
func ValidateDomainName(v string) error {
	if err := validateNonEmpty(v); err != nil {
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

// ValidateEmail checks if the given string is a valid email address and returns an error if not.
func ValidateEmail(v string) error {
	if err := validateNonEmpty(v); err != nil {
		return err
	}

	invalidError := errors.New("invalid email address")

	split := strings.Split(v, "@")
	if len(split) != 2 {
		return invalidError
	}
	localPart := split[0]
	domain := split[1]

	if validateNonEmpty(localPart) != nil {
		return invalidError
	}

	// No whitespace allowed in local-part
	if isMatch(`\s`, localPart) {
		return invalidError
	}

	return ValidateDomainName(domain)
}

// ValidateJSON validates that the given data is valid JSON.
func ValidateJSON(data []byte) error {
	var dummy interface{}
	return json.Unmarshal(data, &dummy)
}

// prefixError wraps an error with a prefix or returns nil if there was no error.
// This is useful for wrapping errors returned by generic error funcs like `validateNonEmpty` so that the error includes the offending field name.
func prefixError(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %v", prefix, err)
	}
	return nil
}

func isMatch(re string, v string) bool {
	return regexp.MustCompile(re).MatchString(v)
}

// validateClusterName checks if the given string is a valid name for a cluster and returns an error if not.
func validateClusterName(v string) error {
	if err := validateNonEmpty(v); err != nil {
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

// validateNonEmpty checks if the given string contains at least one non-whitespace character and returns an error if not.
func validateNonEmpty(v string) error {
	if utf8.RuneCountInString(strings.TrimSpace(v)) == 0 {
		return errors.New("cannot be empty")
	}
	return nil
}

// validateSubnetCIDR checks if the given string is a valid CIDR for a master nodes or worker nodes subnet and returns an error if not.
func validateSubnetCIDR(v string) error {
	if err := validateNonEmpty(v); err != nil {
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

	if err := validateIPv4(ip); err != nil {
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

// validateCIDRsDontOverlap ensures two given CIDRs don't overlap
// with one another. CIDR starting IPs are canonicalized
// before being compared.
func validateCIDRsDontOverlap(acidr, bcidr string) error {
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

// validateIPv4 checks if the given string is a valid IP v4 address and returns an error if not.
// Based on net.ParseIP.
func validateIPv4(v string) error {
	if err := validateNonEmpty(v); err != nil {
		return err
	}
	if ip := net.ParseIP(v); ip == nil || !strings.Contains(v, ".") {
		return errors.New("invalid IPv4 address")
	}
	return nil
}

// validateIntRange checks if the given string is a valid integer between `min` and `max` and returns an error if not.
func validateIntRange(v string, min int, max int) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return validateInt(v)
	}
	if i < min {
		return fmt.Errorf("cannot be less than %v", min)
	}
	if i > max {
		return fmt.Errorf("cannot be greater than %v", max)
	}
	return nil
}

// validateInt checks if the given string is a valid integer and returns an error if not.
func validateInt(v string) error {
	if err := validateNonEmpty(v); err != nil {
		return err
	}

	if _, err := strconv.Atoi(v); err != nil {
		return errors.New("invalid integer")
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
