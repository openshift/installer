package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/openshift/installer/installer/pkg/validate"
	"github.com/openshift/installer/pkg/types/config/aws"

	"github.com/coreos/tectonic-config/config/tectonic-network"
	log "github.com/sirupsen/logrus"
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
	errs = append(errs, c.validatePullSecret()...)
	errs = append(errs, c.validateLibvirt()...)
	if err := validate.PrefixError("cluster name", validate.ClusterName(c.Name)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("base domain", validate.DomainName(c.BaseDomain)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("admin password", validate.NonEmpty(c.Admin.Password)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("admin email", validate.Email(c.Admin.Email)); err != nil {
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
	if err := validate.PrefixError("aws vpcCIDRBlock", validate.SubnetCIDR(c.AWS.VPCCIDRBlock)); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, c.validateOverlapWithPodOrServiceCIDR(c.AWS.VPCCIDRBlock, "aws vpcCIDRBlock")...)
	if err := validate.PrefixError("aws profile", validate.NonEmpty(c.AWS.Profile)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("aws region", validate.NonEmpty(c.AWS.Region)); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// validateOverlapWithPodOrServiceCIDR ensures that the given CIDR does not
// overlap with the pod or service CIDRs of the cluster config.
func (c *Cluster) validateOverlapWithPodOrServiceCIDR(cidr, name string) []error {
	var errs []error
	if err := validate.PrefixError(fmt.Sprintf("%s and podCIDR", name), validate.CIDRsDontOverlap(cidr, c.Networking.PodCIDR)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError(fmt.Sprintf("%s and serviceCIDR", name), validate.CIDRsDontOverlap(cidr, c.Networking.ServiceCIDR)); err != nil {
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
	if err := validate.PrefixError("libvirt network ipRange", validate.SubnetCIDR(c.Libvirt.Network.IPRange)); err != nil {
		errs = append(errs, err)
	}
	if len(c.Libvirt.MasterIPs) > 0 {
		if len(c.Libvirt.MasterIPs) != c.NodeCount(c.Master.NodePools) {
			errs = append(errs, fmt.Errorf("length of masterIPs does't match master count"))
		}
		for i, ip := range c.Libvirt.MasterIPs {
			if err := validate.PrefixError(fmt.Sprintf("libvirt masterIPs[%d] %q", i, ip), validate.IPv4(ip)); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if err := validate.PrefixError("libvirt uri", validate.NonEmpty(c.Libvirt.URI)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("libvirt network name", validate.NonEmpty(c.Libvirt.Network.Name)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("libvirt network ifName", validate.NonEmpty(c.Libvirt.Network.IfName)); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, c.validateOverlapWithPodOrServiceCIDR(c.Libvirt.Network.IPRange, "libvirt ipRange")...)
	return errs
}

func (c *Cluster) validateNetworking() []error {
	var errs []error
	// https://en.wikipedia.org/wiki/Maximum_transmission_unit#MTUs_for_common_media
	if err := validate.PrefixError("mtu", validate.IntRange(c.Networking.MTU, 68, 64*1024)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("podCIDR", validate.SubnetCIDR(c.Networking.PodCIDR)); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("serviceCIDR", validate.SubnetCIDR(c.Networking.ServiceCIDR)); err != nil {
		errs = append(errs, err)
	}
	if err := c.validateNetworkType(); err != nil {
		errs = append(errs, err)
	}
	if err := validate.PrefixError("pod and service CIDRs", validate.CIDRsDontOverlap(c.Networking.PodCIDR, c.Networking.ServiceCIDR)); err != nil {
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
		s := ""
		if len(errs) != 1 {
			s = "s"
		}
		log.Errorf("Found %d error%s in the cluster definition:", len(errs), s)
		for i, err := range errs {
			log.Errorf("error %d: %v", i+1, err)
		}
		return fmt.Errorf("found %d cluster definition error%s", len(errs), s)
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
func (c *Cluster) validateS3Bucket() error {
	bucket := fmt.Sprintf("%s.%s", c.Name, c.BaseDomain)
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
	if err := validate.JSON([]byte(c.PullSecret)); err != nil {
		errs = append(errs, err)
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
