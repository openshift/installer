// custom.go file has custom models for assisted-service that are not
// auto-generated via the swagger.yaml file due to the need for custom
// validation or fields
package models

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"github.com/openshift/assisted-service/pkg/validations"
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
func (m *DomainResolutionRequestDomain) validateDomainName(formats strfmt.Registry) error {
	if err := validate.Required("domain_name", "body", m.DomainName); err != nil {
		return err
	}

	if _, err := validations.ValidateDomainNameFormat(*m.DomainName); err != nil {
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
