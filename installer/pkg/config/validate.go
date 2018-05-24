package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/coreos/tectonic-installer/installer/pkg/validate"

	log "github.com/Sirupsen/logrus"
	ignconfig "github.com/coreos/ignition/config/v2_0"
	"github.com/coreos/tectonic-config/config/tectonic-network"
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
	errs = append(errs, c.validateIgnitionFiles()...)
	errs = append(errs, c.validateNetworking()...)
	if err := validate.PrefixError("cluster name", validate.ClusterName(c.Name)); err != nil {
		errs = append(errs, err)
	}
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

	var podOK, serviceOK bool
	_, pod, err := net.ParseCIDR(c.Networking.PodCIDR)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid pod CIDR %q: %v", c.Networking.PodCIDR, err))
	} else if err := validate.CanonicalizeIP(&pod.IP); err != nil {
		errs = append(errs, fmt.Errorf("invalid pod CIDR %q: %v", c.Networking.PodCIDR, err))
	} else {
		podOK = true
	}
	_, service, err := net.ParseCIDR(c.Networking.ServiceCIDR)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid service CIDR %q: %v", c.Networking.ServiceCIDR, err))
	} else if err := validate.CanonicalizeIP(&service.IP); err != nil {
		errs = append(errs, fmt.Errorf("invalid service CIDR %q: %v", c.Networking.ServiceCIDR, err))
	} else {
		serviceOK = true
	}
	if podOK && serviceOK && validate.CIDRsOverlap(pod, service) {
		errs = append(errs, errors.New("pod and service CIDRs overlap"))
	}
	return errs
}

func (c *Cluster) validateNetworkType() error {
	switch tectonicnetwork.NetworkType(c.Networking.Type) {
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

func (c *Cluster) validateIgnitionFiles() []error {
	var errs []error
	for _, n := range c.NodePools {
		if n.IgnitionFile == "" {
			continue
		}

		if err := validateFileExist(n.IgnitionFile); err != nil {
			errs = append(errs, err)
			continue
		}

		if err := validateIgnitionConfig(n.IgnitionFile); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func validateFileExist(ignitionFile string) error {
	_, err := os.Stat(ignitionFile)
	return err
}

func validateIgnitionConfig(filePath string) error {
	blob, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, rpt, _ := ignconfig.Parse(blob)
	if len(rpt.Entries) > 0 {
		return &ErrInvalidIgnConfig{
			filePath,
			rpt.String(),
		}
	}
	return nil
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
		{pools: c.Etcd.NodePools, field: "etcd"},
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
			for j := range c.Etcd.NodePools {
				if c.Master.NodePools[i] == c.Etcd.NodePools[j] {
					if fields[c.Master.NodePools[i]] == nil {
						fields[c.Master.NodePools[i]] = make(map[string]struct{})
					}
					fields[c.Master.NodePools[i]]["master"] = struct{}{}
					fields[c.Master.NodePools[i]]["etcd"] = struct{}{}
				}
			}
		}
	}
	for i := range c.Worker.NodePools {
		if c.Worker.NodePools[i] != "" {
			for j := range c.Etcd.NodePools {
				if c.Worker.NodePools[i] == c.Etcd.NodePools[j] {
					if fields[c.Worker.NodePools[i]] == nil {
						fields[c.Worker.NodePools[i]] = make(map[string]struct{})
					}
					fields[c.Worker.NodePools[i]]["worker"] = struct{}{}
					fields[c.Worker.NodePools[i]]["etcd"] = struct{}{}
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
