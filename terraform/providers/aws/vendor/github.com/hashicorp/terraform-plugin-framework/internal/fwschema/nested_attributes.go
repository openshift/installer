package fwschema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// NestingMode is an enum type of the ways nested attributes can be nested in
// an attribute. They can be a list, a set, a map (with string
// keys), or they can be nested directly, like an object.
type NestingMode uint8

const (
	// NestingModeUnknown is an invalid nesting mode, used to catch when a
	// nesting mode is expected and not set.
	NestingModeUnknown NestingMode = 0

	// NestingModeSingle is for attributes that represent a struct or
	// object, a single instance of those attributes directly nested under
	// another attribute.
	NestingModeSingle NestingMode = 1

	// NestingModeList is for attributes that represent a list of objects,
	// with multiple instances of those attributes nested inside a list
	// under another attribute.
	NestingModeList NestingMode = 2

	// NestingModeSet is for attributes that represent a set of objects,
	// with multiple, unique instances of those attributes nested inside a
	// set under another attribute.
	NestingModeSet NestingMode = 3

	// NestingModeMap is for attributes that represent a map of objects,
	// with multiple instances of those attributes, each associated with a
	// unique string key, nested inside a map under another attribute.
	NestingModeMap NestingMode = 4
)

// NestedAttributes surfaces a group of attributes to nest beneath another
// attribute, and how that nesting should behave. Nesting can have the
// following modes:
//
// * SingleNestedAttributes are nested attributes that represent a struct or
// object; there should only be one instance of them nested beneath that
// specific attribute.
//
// * ListNestedAttributes are nested attributes that represent a list of
// structs or objects; there can be multiple instances of them beneath that
// specific attribute.
//
// * SetNestedAttributes are nested attributes that represent a set of structs
// or objects; there can be multiple instances of them beneath that specific
// attribute. Unlike ListNestedAttributes, these nested attributes must have
// unique values.
//
// * MapNestedAttributes are nested attributes that represent a string-indexed
// map of structs or objects; there can be multiple instances of them beneath
// that specific attribute. Unlike ListNestedAttributes, these nested
// attributes must be associated with a unique key. Unlike SetNestedAttributes,
// the key must be explicitly set by the user.
type NestedAttributes interface {
	// Implementations should include the tftypes.AttributePathStepper
	// interface methods for proper path and data handling.
	tftypes.AttributePathStepper

	// AttributeType should return the framework type of the nested attributes.
	// This method should be deprecated in preference of Type().
	AttributeType() attr.Type

	// Equal should return true if the other NestedAttributes is equivalent.
	Equal(NestedAttributes) bool

	// GetNestingMode should return the nesting mode (list, map, set, or
	// single) of the nested attributes.
	GetNestingMode() NestingMode

	// GetAttributes() should return the mapping of names to nested attributes.
	GetAttributes() map[string]Attribute

	// Type should return the framework type of the nested attributes.
	Type() attr.Type
}

type UnderlyingAttributes map[string]Attribute

func (n UnderlyingAttributes) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	a, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("can't apply %T to Attributes", step)
	}

	res, ok := n[string(a)]

	if !ok {
		return nil, fmt.Errorf("no attribute %q on Attributes", a)
	}

	return res, nil
}

// Type returns the framework type of the nested attributes.
func (n UnderlyingAttributes) Type() attr.Type {
	attrTypes := map[string]attr.Type{}
	for name, attr := range n {
		attrTypes[name] = attr.FrameworkType()
	}
	return types.ObjectType{
		AttrTypes: attrTypes,
	}
}

type SingleNestedAttributes struct {
	UnderlyingAttributes
}

func (s SingleNestedAttributes) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	a, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("can't apply %T to Attributes", step)
	}

	res, ok := s.UnderlyingAttributes[string(a)]

	if !ok {
		return nil, fmt.Errorf("no attribute %q on Attributes", a)
	}

	return res, nil
}

// Deprecated: Use Type() instead.
func (s SingleNestedAttributes) AttributeType() attr.Type {
	return s.Type()
}

func (s SingleNestedAttributes) GetAttributes() map[string]Attribute {
	return s.UnderlyingAttributes
}

func (s SingleNestedAttributes) GetNestingMode() NestingMode {
	return NestingModeSingle
}

func (s SingleNestedAttributes) Equal(o NestedAttributes) bool {
	other, ok := o.(SingleNestedAttributes)
	if !ok {
		return false
	}
	if len(other.UnderlyingAttributes) != len(s.UnderlyingAttributes) {
		return false
	}
	for k, v := range s.UnderlyingAttributes {
		otherV, ok := other.UnderlyingAttributes[k]
		if !ok {
			return false
		}
		if !v.Equal(otherV) {
			return false
		}
	}
	return true
}

// Type returns the framework type of the nested attributes.
func (s SingleNestedAttributes) Type() attr.Type {
	return s.UnderlyingAttributes.Type()
}

type ListNestedAttributes struct {
	UnderlyingAttributes
}

func (l ListNestedAttributes) GetAttributes() map[string]Attribute {
	return l.UnderlyingAttributes
}

func (l ListNestedAttributes) GetNestingMode() NestingMode {
	return NestingModeList
}

// AttributeType returns an attr.Type corresponding to the nested attributes.
// Deprecated: Use Type() instead.
func (l ListNestedAttributes) AttributeType() attr.Type {
	return l.Type()
}

func (l ListNestedAttributes) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	_, ok := step.(tftypes.ElementKeyInt)
	if !ok {
		return nil, fmt.Errorf("can't apply %T to ListNestedAttributes", step)
	}
	return l.UnderlyingAttributes, nil
}

func (l ListNestedAttributes) Equal(o NestedAttributes) bool {
	other, ok := o.(ListNestedAttributes)
	if !ok {
		return false
	}
	if len(other.UnderlyingAttributes) != len(l.UnderlyingAttributes) {
		return false
	}
	for k, v := range l.UnderlyingAttributes {
		otherV, ok := other.UnderlyingAttributes[k]
		if !ok {
			return false
		}
		if !v.Equal(otherV) {
			return false
		}
	}
	return true
}

// Type returns the framework type of the nested attributes.
func (l ListNestedAttributes) Type() attr.Type {
	return types.ListType{
		ElemType: l.UnderlyingAttributes.Type(),
	}
}

type SetNestedAttributes struct {
	UnderlyingAttributes
}

func (s SetNestedAttributes) GetAttributes() map[string]Attribute {
	return s.UnderlyingAttributes
}

func (s SetNestedAttributes) GetNestingMode() NestingMode {
	return NestingModeSet
}

// AttributeType returns an attr.Type corresponding to the nested attributes.
// Deprecated: Use Type() instead.
func (s SetNestedAttributes) AttributeType() attr.Type {
	return s.Type()
}

func (s SetNestedAttributes) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	_, ok := step.(tftypes.ElementKeyValue)
	if !ok {
		return nil, fmt.Errorf("can't use %T on sets", step)
	}
	return s.UnderlyingAttributes, nil
}

func (s SetNestedAttributes) Equal(o NestedAttributes) bool {
	other, ok := o.(SetNestedAttributes)
	if !ok {
		return false
	}
	if len(other.UnderlyingAttributes) != len(s.UnderlyingAttributes) {
		return false
	}
	for k, v := range s.UnderlyingAttributes {
		otherV, ok := other.UnderlyingAttributes[k]
		if !ok {
			return false
		}
		if !v.Equal(otherV) {
			return false
		}
	}
	return true
}

// Type returns the framework type of the nested attributes.
func (s SetNestedAttributes) Type() attr.Type {
	return types.SetType{
		ElemType: s.UnderlyingAttributes.Type(),
	}
}

type MapNestedAttributes struct {
	UnderlyingAttributes
}

func (m MapNestedAttributes) GetAttributes() map[string]Attribute {
	return m.UnderlyingAttributes
}

func (m MapNestedAttributes) GetNestingMode() NestingMode {
	return NestingModeMap
}

// AttributeType returns an attr.Type corresponding to the nested attributes.
// Deprecated: Use Type() instead.
func (m MapNestedAttributes) AttributeType() attr.Type {
	return m.Type()
}

func (m MapNestedAttributes) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	_, ok := step.(tftypes.ElementKeyString)
	if !ok {
		return nil, fmt.Errorf("can't use %T on maps", step)
	}
	return m.UnderlyingAttributes, nil
}

func (m MapNestedAttributes) Equal(o NestedAttributes) bool {
	other, ok := o.(MapNestedAttributes)
	if !ok {
		return false
	}
	if len(other.UnderlyingAttributes) != len(m.UnderlyingAttributes) {
		return false
	}
	for k, v := range m.UnderlyingAttributes {
		otherV, ok := other.UnderlyingAttributes[k]
		if !ok {
			return false
		}
		if !v.Equal(otherV) {
			return false
		}
	}
	return true
}

// Type returns the framework type of the nested attributes.
func (m MapNestedAttributes) Type() attr.Type {
	return types.MapType{
		ElemType: m.UnderlyingAttributes.Type(),
	}
}
