// Code generated by go-swagger; DO NOT EDIT.

package openapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Authorization authorization
//
// swagger:model Authorization
type Authorization struct {

	// id
	ID string `json:"id,omitempty"`

	// requesting party
	// Required: true
	RequestingParty *string `json:"requestingParty"`

	// scope
	// Required: true
	Scope *AuthorizationScope `json:"scope"`

	// zcap
	Zcap string `json:"zcap,omitempty"`
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *Authorization) UnmarshalJSON(raw []byte) error {
	var data struct {
		ID string `json:"id,omitempty"`

		RequestingParty *string `json:"requestingParty"`

		Scope *AuthorizationScope `json:"scope"`

		Zcap string `json:"zcap,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var result Authorization

	// id
	result.ID = data.ID

	// requestingParty
	result.RequestingParty = data.RequestingParty

	// scope
	result.Scope = data.Scope

	// zcap
	result.Zcap = data.Zcap

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m Authorization) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {
		ID string `json:"id,omitempty"`

		RequestingParty *string `json:"requestingParty"`

		Scope *AuthorizationScope `json:"scope"`

		Zcap string `json:"zcap,omitempty"`
	}{

		ID: m.ID,

		RequestingParty: m.RequestingParty,

		Scope: m.Scope,

		Zcap: m.Zcap,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
	}{})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this authorization
func (m *Authorization) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRequestingParty(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateScope(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Authorization) validateRequestingParty(formats strfmt.Registry) error {

	if err := validate.Required("requestingParty", "body", m.RequestingParty); err != nil {
		return err
	}

	return nil
}

func (m *Authorization) validateScope(formats strfmt.Registry) error {

	if err := validate.Required("scope", "body", m.Scope); err != nil {
		return err
	}

	if m.Scope != nil {
		if err := m.Scope.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("scope")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this authorization based on the context it is used
func (m *Authorization) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateScope(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Authorization) contextValidateScope(ctx context.Context, formats strfmt.Registry) error {

	if m.Scope != nil {
		if err := m.Scope.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("scope")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Authorization) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Authorization) UnmarshalBinary(b []byte) error {
	var res Authorization
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// AuthorizationScope authorization scope
//
// swagger:model AuthorizationScope
type AuthorizationScope struct {

	// action
	// Required: true
	Action []string `json:"action"`

	caveatsField []Caveat

	// resource ID
	// Required: true
	ResourceID *string `json:"resourceID"`

	// resource type
	// Required: true
	ResourceType *string `json:"resourceType"`
}

// Caveats gets the caveats of this base type
func (m *AuthorizationScope) Caveats() []Caveat {
	return m.caveatsField
}

// SetCaveats sets the caveats of this base type
func (m *AuthorizationScope) SetCaveats(val []Caveat) {
	m.caveatsField = val
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *AuthorizationScope) UnmarshalJSON(raw []byte) error {
	var data struct {
		Action []string `json:"action"`

		Caveats json.RawMessage `json:"caveats"`

		ResourceID *string `json:"resourceID"`

		ResourceType *string `json:"resourceType"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var propCaveats []Caveat
	if string(data.Caveats) != "null" {
		caveats, err := UnmarshalCaveatSlice(bytes.NewBuffer(data.Caveats), runtime.JSONConsumer())
		if err != nil && err != io.EOF {
			return err
		}
		propCaveats = caveats
	}

	var result AuthorizationScope

	// action
	result.Action = data.Action

	// caveats
	result.caveatsField = propCaveats

	// resourceID
	result.ResourceID = data.ResourceID

	// resourceType
	result.ResourceType = data.ResourceType

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m AuthorizationScope) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {
		Action []string `json:"action"`

		ResourceID *string `json:"resourceID"`

		ResourceType *string `json:"resourceType"`
	}{

		Action: m.Action,

		ResourceID: m.ResourceID,

		ResourceType: m.ResourceType,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		Caveats []Caveat `json:"caveats"`
	}{

		Caveats: m.caveatsField,
	})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this authorization scope
func (m *AuthorizationScope) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAction(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCaveats(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var authorizationScopeActionItemsEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["read","reference"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		authorizationScopeActionItemsEnum = append(authorizationScopeActionItemsEnum, v)
	}
}

func (m *AuthorizationScope) validateActionItemsEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, authorizationScopeActionItemsEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *AuthorizationScope) validateAction(formats strfmt.Registry) error {

	if err := validate.Required("scope"+"."+"action", "body", m.Action); err != nil {
		return err
	}

	for i := 0; i < len(m.Action); i++ {

		// value enum
		if err := m.validateActionItemsEnum("scope"+"."+"action"+"."+strconv.Itoa(i), "body", m.Action[i]); err != nil {
			return err
		}

	}

	return nil
}

func (m *AuthorizationScope) validateCaveats(formats strfmt.Registry) error {
	if swag.IsZero(m.Caveats()) { // not required
		return nil
	}

	for i := 0; i < len(m.Caveats()); i++ {

		if err := m.caveatsField[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("scope" + "." + "caveats" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *AuthorizationScope) validateResourceID(formats strfmt.Registry) error {

	if err := validate.Required("scope"+"."+"resourceID", "body", m.ResourceID); err != nil {
		return err
	}

	return nil
}

func (m *AuthorizationScope) validateResourceType(formats strfmt.Registry) error {

	if err := validate.Required("scope"+"."+"resourceType", "body", m.ResourceType); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this authorization scope based on the context it is used
func (m *AuthorizationScope) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCaveats(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AuthorizationScope) contextValidateCaveats(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Caveats()); i++ {

		if err := m.caveatsField[i].ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("scope" + "." + "caveats" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *AuthorizationScope) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AuthorizationScope) UnmarshalBinary(b []byte) error {
	var res AuthorizationScope
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
