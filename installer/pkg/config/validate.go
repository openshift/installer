package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	ignconfig "github.com/coreos/ignition/config/v2_0"
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
	return errs
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
