// custom.go file has custom models for assisted-service that are not
// auto-generated via the swagger.yaml file due to the need for custom
// validation or fields
package models

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net"
	"net/netip"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	modelvalidations "github.com/openshift/assisted-service/models/validations"
)

// DomainResolutionRequestDomain is a struct to hold the domain resolution request domain
type DomainResolutionRequestDomain struct {

	// The domain name that should be resolved
	// Required: true
	DomainName *string `json:"domain_name"`
}

// Validate is a function required for interfaces derived from swagger models and it
// validates this domain resolution request domain
func (m *DomainResolutionRequestDomain) Validate(formats strfmt.Registry) error {
	if err := m.validateDomainName(formats); err != nil {
		return err
	}

	return nil
}

// validateDomainName ensures that the required DomainName field exists and that the
// DomainName is valid
func (m *DomainResolutionRequestDomain) validateDomainName(_ strfmt.Registry) error {
	if err := validate.Required("domain_name", "body", m.DomainName); err != nil {
		return err
	}

	if _, err := modelvalidations.ValidateDomainNameFormat(*m.DomainName); err != nil {
		return err
	}

	return nil
}

// The following functions (ContextValidate, MarshalBinary, UnmarshalBinary) are required for
// interfaces derived from swagger models

// ContextValidate validates this domain resolution request domain based on context it is used
func (m *DomainResolutionRequestDomain) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DomainResolutionRequestDomain) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DomainResolutionRequestDomain) UnmarshalBinary(b []byte) error {
	var res DomainResolutionRequestDomain
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Normalize returns the IP in canonical form (e.g. 2620:0052:0009:164d:0000:0000:0000:0046 → 2620:52:9:164d::46).
func (m IP) Normalize() (IP, error) {
	if m == "" {
		return m, nil
	}
	addr, err := netip.ParseAddr(string(m))
	if err != nil {
		return m, fmt.Errorf("invalid IP address %q: %w", m, err)
	}
	return IP(addr.String()), nil
}

// Equal returns true if two IPs represent the same address, handling different IPv6 representations.
// Falls back to string comparison if either value cannot be parsed.
func (m IP) Equal(other IP) bool {
	normA, errA := m.Normalize()
	normB, errB := other.Normalize()
	if errA != nil || errB != nil {
		return string(m) == string(other)
	}
	return normA == normB
}

// Value implements driver.Valuer for IP to convert Go string to PostgreSQL inet.
// Returns inet text (e.g. 203.0.113.1/32); works with database/sql drivers that accept string parameters.
func (m IP) Value() (driver.Value, error) {
	if m == "" {
		return nil, nil
	}
	addr, err := netip.ParseAddr(string(m))
	if err != nil {
		return nil, fmt.Errorf("invalid IP address %q: %w", m, err)
	}
	prefix := netip.PrefixFrom(addr, addr.BitLen())
	return prefix.String(), nil
}

// Scan implements sql.Scanner for IP to convert PostgreSQL inet to Go string
func (m *IP) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*m = IP(v)
	case []byte:
		*m = IP(string(v))
	case *net.IPNet:
		if v != nil {
			*m = IP(v.IP.String())
		} else {
			*m = ""
		}
	case net.IP:
		*m = IP(v.String())
	case netip.Prefix:
		if !v.IsValid() {
			*m = ""
			return nil
		}
		*m = IP(v.Addr().String())
	case *netip.Prefix:
		if v == nil || !v.IsValid() {
			*m = ""
			return nil
		}
		*m = IP(v.Addr().String())
	default:
		return fmt.Errorf("cannot convert %T to IP", value)
	}
	return nil
}

// Normalize returns the Subnet in canonical form (e.g. fcff:0069:0001::/64 → fcff:69:1::/64).
func (m Subnet) Normalize() (Subnet, error) {
	if m == "" {
		return m, nil
	}
	_, ipnet, err := net.ParseCIDR(string(m))
	if err != nil {
		return m, fmt.Errorf("invalid CIDR %q: %w", m, err)
	}
	return Subnet(ipnet.String()), nil
}

// Equal returns true if two Subnets represent the same network, handling different IPv6 representations.
// Falls back to string comparison if either value cannot be parsed.
func (m Subnet) Equal(other Subnet) bool {
	normA, errA := m.Normalize()
	normB, errB := other.Normalize()
	if errA != nil || errB != nil {
		return string(m) == string(other)
	}
	return string(normA) == string(normB)
}

// Value implements driver.Valuer for Subnet to convert Go string to PostgreSQL cidr.
func (m Subnet) Value() (driver.Value, error) {
	if m == "" {
		return nil, nil
	}
	prefix, err := netip.ParsePrefix(string(m))
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR %q: %w", m, err)
	}
	prefix = prefix.Masked()
	return prefix.String(), nil
}

// Scan implements sql.Scanner for Subnet to convert PostgreSQL cidr to Go string
func (m *Subnet) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		*m = Subnet(v)
	case []byte:
		*m = Subnet(string(v))
	case *net.IPNet:
		if v != nil {
			*m = Subnet(v.String())
		} else {
			*m = ""
		}
	case net.IPNet:
		*m = Subnet(v.String())
	case netip.Prefix:
		if !v.IsValid() {
			*m = ""
			return nil
		}
		*m = Subnet(v.Masked().String())
	case *netip.Prefix:
		if v == nil || !v.IsValid() {
			*m = ""
			return nil
		}
		*m = Subnet(v.Masked().String())
	default:
		return fmt.Errorf("cannot convert %T to Subnet", value)
	}
	return nil
}
