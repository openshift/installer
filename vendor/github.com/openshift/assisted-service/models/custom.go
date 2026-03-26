// custom.go file has custom models for assisted-service that are not
// auto-generated via the swagger.yaml file due to the need for custom
// validation or fields
package models

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/jackc/pgtype"
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

// Value implements driver.Valuer for IP to convert Go string to PostgreSQL inet
func (m IP) Value() (driver.Value, error) {
	if m == "" {
		return nil, nil
	}

	// Create a pgtype.Inet and set the value - it knows how to encode for PostgreSQL
	var inet pgtype.Inet
	if err := inet.Set(string(m)); err != nil {
		return nil, fmt.Errorf("invalid IP address %q: %w", m, err)
	}

	// Return the pgtype.Inet's driver.Value
	return inet.Value()
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
		// pgx may return *net.IPNet for inet columns
		if v != nil {
			*m = IP(v.IP.String())
		} else {
			*m = ""
		}
	case net.IP:
		*m = IP(v.String())
	case *pgtype.Inet:
		// pgtype.Inet is used by pgx for inet/cidr columns
		if v != nil && v.Status == pgtype.Present && v.IPNet != nil {
			*m = IP(v.IPNet.IP.String())
		} else {
			*m = ""
		}
	case pgtype.Inet:
		if v.Status == pgtype.Present && v.IPNet != nil {
			*m = IP(v.IPNet.IP.String())
		} else {
			*m = ""
		}
	default:
		return fmt.Errorf("cannot convert %T to IP", value)
	}
	return nil
}

// Value implements driver.Valuer for Subnet to convert Go string to PostgreSQL cidr
func (m Subnet) Value() (driver.Value, error) {
	if m == "" {
		return nil, nil
	}

	// Parse and normalize the CIDR to ensure host bits are zero
	// This is required for PostgreSQL cidr type which is strict about network addresses
	_, ipnet, err := net.ParseCIDR(string(m))
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR %q: %w", m, err)
	}

	// Create a pgtype.Inet with the normalized network address
	var inet pgtype.Inet
	if err := inet.Set(ipnet.String()); err != nil {
		return nil, fmt.Errorf("invalid CIDR %q: %w", m, err)
	}

	// Return the pgtype.Inet's driver.Value
	return inet.Value()
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
		// pgx may return *net.IPNet for cidr columns
		if v != nil {
			*m = Subnet(v.String())
		} else {
			*m = ""
		}
	case net.IPNet:
		*m = Subnet(v.String())
	case *pgtype.Inet:
		// pgtype.Inet is used by pgx for inet/cidr columns
		if v != nil && v.Status == pgtype.Present && v.IPNet != nil {
			*m = Subnet(v.IPNet.String())
		} else {
			*m = ""
		}
	case pgtype.Inet:
		if v.Status == pgtype.Present && v.IPNet != nil {
			*m = Subnet(v.IPNet.String())
		} else {
			*m = ""
		}
	default:
		return fmt.Errorf("cannot convert %T to Subnet", value)
	}
	return nil
}
