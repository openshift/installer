package quotas

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

type detailResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a Quota resource.
func (r commonResult) Extract() (*Quota, error) {
	var s struct {
		Quota *Quota `json:"quota"`
	}
	err := r.ExtractInto(&s)
	return s.Quota, err
}

// Extract is a function that accepts a result and extracts a QuotaDetailSet resource.
func (r detailResult) Extract() (*QuotaDetailSet, error) {
	var s struct {
		Quota *QuotaDetailSet `json:"quota"`
	}
	err := r.ExtractInto(&s)
	return s.Quota, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Quota.
type GetResult struct {
	commonResult
}

// GetDetailResult represents the detailed result of a get operation. Call its Extract
// method to interpret it as a Quota.
type GetDetailResult struct {
	detailResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Quota.
type UpdateResult struct {
	commonResult
}

// Quota contains Networking quotas for a project.
type Quota struct {
	// FloatingIP represents a number of floating IPs. A "-1" value means no limit.
	FloatingIP int `json:"floatingip"`

	// Network represents a number of networks. A "-1" value means no limit.
	Network int `json:"network"`

	// Port represents a number of ports. A "-1" value means no limit.
	Port int `json:"port"`

	// RBACPolicy represents a number of RBAC policies. A "-1" value means no limit.
	RBACPolicy int `json:"rbac_policy"`

	// Router represents a number of routers. A "-1" value means no limit.
	Router int `json:"router"`

	// SecurityGroup represents a number of security groups. A "-1" value means no limit.
	SecurityGroup int `json:"security_group"`

	// SecurityGroupRule represents a number of security group rules. A "-1" value means no limit.
	SecurityGroupRule int `json:"security_group_rule"`

	// Subnet represents a number of subnets. A "-1" value means no limit.
	Subnet int `json:"subnet"`

	// SubnetPool represents a number of subnet pools. A "-1" value means no limit.
	SubnetPool int `json:"subnetpool"`

	// Trunk represents a number of trunks. A "-1" value means no limit.
	Trunk int `json:"trunk"`
}

// QuotaDetailSet represents details of both operational limits of Networking resources for a project
// and the current usage of those resources.
type QuotaDetailSet struct {
	// FloatingIP represents a number of floating IPs. A "-1" value means no limit.
	FloatingIP QuotaDetail `json:"floatingip"`

	// Network represents a number of networks. A "-1" value means no limit.
	Network QuotaDetail `json:"network"`

	// Port represents a number of ports. A "-1" value means no limit.
	Port QuotaDetail `json:"port"`

	// RBACPolicy represents a number of RBAC policies. A "-1" value means no limit.
	RBACPolicy QuotaDetail `json:"rbac_policy"`

	// Router represents a number of routers. A "-1" value means no limit.
	Router QuotaDetail `json:"router"`

	// SecurityGroup represents a number of security groups. A "-1" value means no limit.
	SecurityGroup QuotaDetail `json:"security_group"`

	// SecurityGroupRule represents a number of security group rules. A "-1" value means no limit.
	SecurityGroupRule QuotaDetail `json:"security_group_rule"`

	// Subnet represents a number of subnets. A "-1" value means no limit.
	Subnet QuotaDetail `json:"subnet"`

	// SubnetPool represents a number of subnet pools. A "-1" value means no limit.
	SubnetPool QuotaDetail `json:"subnetpool"`

	// Trunk represents a number of trunks. A "-1" value means no limit.
	Trunk QuotaDetail `json:"trunk"`
}

// QuotaDetail is a set of details about a single operational limit that allows
// for control of networking usage.
type QuotaDetail struct {
	// Used is the current number of provisioned/allocated resources of the
	// given type.
	Used int `json:"used"`

	// Reserved is a transitional state when a claim against quota has been made
	// but the resource is not yet fully online.
	Reserved int `json:"reserved"`

	// Limit is the maximum number of a given resource that can be
	// allocated/provisioned.  This is what "quota" usually refers to.
	Limit int `json:"limit"`
}

// UnmarshalJSON overrides the default unmarshalling function to accept
// Reserved as a string.
//
// Due to a bug in Neutron, under some conditions Reserved is returned as a
// string.
//
// This method is left for compatibility with unpatched versions of Neutron.
//
// cf. https://bugs.launchpad.net/neutron/+bug/1918565
func (q *QuotaDetail) UnmarshalJSON(b []byte) error {
	type tmp QuotaDetail
	var s struct {
		tmp
		Reserved interface{} `json:"reserved"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*q = QuotaDetail(s.tmp)

	switch t := s.Reserved.(type) {
	case float64:
		q.Reserved = int(t)
	case string:
		if q.Reserved, err = strconv.Atoi(t); err != nil {
			return err
		}
	default:
		return fmt.Errorf("reserved has unexpected type: %T", t)
	}

	return nil
}
