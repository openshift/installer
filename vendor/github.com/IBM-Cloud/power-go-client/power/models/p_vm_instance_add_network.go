// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PVMInstanceAddNetwork p VM instance add network
//
// swagger:model PVMInstanceAddNetwork
type PVMInstanceAddNetwork struct {

	// The requested ip address of this network interface
	IPAddress string `json:"ipAddress,omitempty"`

	// ID of the network
	// Required: true
	NetworkID *string `json:"networkID"`

	// Network security groups that the network interface is a member of. There is a limit of 1 network security group in the array. If not specified, default network security group is used.
	NetworkSecurityGroupIDs []string `json:"networkSecurityGroupIDs"`
}

// Validate validates this p VM instance add network
func (m *PVMInstanceAddNetwork) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateNetworkID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PVMInstanceAddNetwork) validateNetworkID(formats strfmt.Registry) error {

	if err := validate.Required("networkID", "body", m.NetworkID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this p VM instance add network based on context it is used
func (m *PVMInstanceAddNetwork) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PVMInstanceAddNetwork) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PVMInstanceAddNetwork) UnmarshalBinary(b []byte) error {
	var res PVMInstanceAddNetwork
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
