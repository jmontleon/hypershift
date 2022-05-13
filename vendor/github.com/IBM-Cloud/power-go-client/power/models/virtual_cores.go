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

// VirtualCores virtual cores
//
// swagger:model VirtualCores
type VirtualCores struct {

	// The active virtual Cores
	// Required: true
	// Minimum: 1
	Assigned *int64 `json:"assigned"`

	// The maximum DLPAR range for virtual Cores (Display only support)
	Max int64 `json:"max,omitempty"`

	// The minimum DLPAR range for virtual Cores (Display only support)
	Min int64 `json:"min,omitempty"`
}

// Validate validates this virtual cores
func (m *VirtualCores) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAssigned(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VirtualCores) validateAssigned(formats strfmt.Registry) error {

	if err := validate.Required("assigned", "body", m.Assigned); err != nil {
		return err
	}

	if err := validate.MinimumInt("assigned", "body", *m.Assigned, 1, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this virtual cores based on context it is used
func (m *VirtualCores) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *VirtualCores) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VirtualCores) UnmarshalBinary(b []byte) error {
	var res VirtualCores
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}