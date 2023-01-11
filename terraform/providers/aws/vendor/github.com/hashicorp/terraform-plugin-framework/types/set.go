package types

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var (
	_ attr.Type              = SetType{}
	_ xattr.TypeWithValidate = SetType{}
	_ attr.Value             = &Set{}
)

// SetType is an AttributeType representing a set of values. All values must
// be of the same type, which the provider must specify as the ElemType
// property.
type SetType struct {
	ElemType attr.Type
}

// ElementType returns the attr.Type elements will be created from.
func (st SetType) ElementType() attr.Type {
	return st.ElemType
}

// WithElementType returns a SetType that is identical to `l`, but with the
// element type set to `typ`.
func (st SetType) WithElementType(typ attr.Type) attr.TypeWithElementType {
	return SetType{ElemType: typ}
}

// TerraformType returns the tftypes.Type that should be used to
// represent this type. This constrains what user input will be
// accepted and what kind of data can be set in state. The framework
// will use this to translate the AttributeType to something Terraform
// can understand.
func (st SetType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.Set{
		ElementType: st.ElemType.TerraformType(ctx),
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value.
// This is meant to convert the tftypes.Value into a more convenient Go
// type for the provider to consume the data with.
func (st SetType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	set := Set{
		ElemType: st.ElemType,
	}
	if in.Type() == nil {
		set.Null = true
		return set, nil
	}
	if !in.Type().Equal(st.TerraformType(ctx)) {
		return nil, fmt.Errorf("can't use %s as value of Set with ElementType %T, can only use %s values", in.String(), st.ElemType, st.ElemType.TerraformType(ctx).String())
	}
	if !in.IsKnown() {
		set.Unknown = true
		return set, nil
	}
	if in.IsNull() {
		set.Null = true
		return set, nil
	}
	val := []tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}
	elems := make([]attr.Value, 0, len(val))
	for _, elem := range val {
		av, err := st.ElemType.ValueFromTerraform(ctx, elem)
		if err != nil {
			return nil, err
		}
		elems = append(elems, av)
	}
	set.Elems = elems
	return set, nil
}

// Equal returns true if `o` is also a SetType and has the same ElemType.
func (st SetType) Equal(o attr.Type) bool {
	if st.ElemType == nil {
		return false
	}
	other, ok := o.(SetType)
	if !ok {
		return false
	}
	return st.ElemType.Equal(other.ElemType)
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// set.
func (st SetType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	if _, ok := step.(tftypes.ElementKeyValue); !ok {
		return nil, fmt.Errorf("cannot apply step %T to SetType", step)
	}

	return st.ElemType, nil
}

// String returns a human-friendly description of the SetType.
func (st SetType) String() string {
	return "types.SetType[" + st.ElemType.String() + "]"
}

// Validate implements type validation. This type requires all elements to be
// unique.
func (st SetType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.Set{}) {
		err := fmt.Errorf("expected Set value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Set Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var elems []tftypes.Value

	if err := in.As(&elems); err != nil {
		diags.AddAttributeError(
			path,
			"Set Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	validatableType, isValidatable := st.ElemType.(xattr.TypeWithValidate)

	// Attempting to use map[tftypes.Value]struct{} for duplicate detection yields:
	//   panic: runtime error: hash of unhashable type tftypes.primitive
	// Instead, use for loops.
	for indexOuter, elemOuter := range elems {
		// Only evaluate fully known values for duplicates and validation.
		if !elemOuter.IsFullyKnown() {
			continue
		}

		// Validate the element first
		if isValidatable {
			elemValue, err := st.ElemType.ValueFromTerraform(ctx, elemOuter)
			if err != nil {
				diags.AddAttributeError(
					path,
					"Set Type Validation Error",
					"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
				)
				return diags
			}
			diags = append(diags, validatableType.Validate(ctx, elemOuter, path.AtSetValue(elemValue))...)
		}

		// Then check for duplicates
		for indexInner := indexOuter + 1; indexInner < len(elems); indexInner++ {
			elemInner := elems[indexInner]

			if !elemInner.Equal(elemOuter) {
				continue
			}

			// TODO: Point at element attr.Value when Validate method is converted to attr.Value
			// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/172
			diags.AddAttributeError(
				path,
				"Duplicate Set Element",
				fmt.Sprintf("This attribute contains duplicate values of: %s", elemInner),
			)
		}
	}

	return diags
}

// ValueType returns the Value type.
func (t SetType) ValueType(_ context.Context) attr.Value {
	return Set{
		ElemType: t.ElemType,
	}
}

// Set represents a set of attr.Value, all of the same type,
// indicated by ElemType.
type Set struct {
	// Unknown will be set to true if the entire set is an unknown value.
	// If only some of the elements in the set are unknown, their known or
	// unknown status will be represented however that attr.Value
	// surfaces that information. The Set's Unknown property only tracks
	// if the number of elements in a Set is known, not whether the
	// elements that are in the set are known.
	Unknown bool

	// Null will be set to true if the set is null, either because it was
	// omitted from the configuration, state, or plan, or because it was
	// explicitly set to null.
	Null bool

	// Elems are the elements in the set.
	Elems []attr.Value

	// ElemType is the tftypes.Type of the elements in the set. All
	// elements in the set must be of this type.
	ElemType attr.Type
}

// ElementsAs populates `target` with the elements of the Set, throwing an
// error if the elements cannot be stored in `target`.
func (s Set) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
	// we need a tftypes.Value for this Set to be able to use it with our
	// reflection code
	val, err := s.ToTerraformValue(ctx)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Set Element Conversion Error",
				"An unexpected error was encountered trying to convert set elements. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			),
		}
	}
	return reflect.Into(ctx, s.Type(ctx), val, target, reflect.Options{
		UnhandledNullAsEmpty:    allowUnhandled,
		UnhandledUnknownAsEmpty: allowUnhandled,
	}, path.Empty())
}

// Type returns a SetType with the same element type as `s`.
func (s Set) Type(ctx context.Context) attr.Type {
	return SetType{ElemType: s.ElemType}
}

// ToTerraformValue returns the data contained in the Set as a tftypes.Value.
func (s Set) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if s.ElemType == nil {
		return tftypes.Value{}, fmt.Errorf("cannot convert Set to tftypes.Value if ElemType field is not set")
	}
	setType := tftypes.Set{ElementType: s.ElemType.TerraformType(ctx)}
	if s.Unknown {
		return tftypes.NewValue(setType, tftypes.UnknownValue), nil
	}
	if s.Null {
		return tftypes.NewValue(setType, nil), nil
	}
	vals := make([]tftypes.Value, 0, len(s.Elems))
	for _, elem := range s.Elems {
		val, err := elem.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(setType, tftypes.UnknownValue), err
		}
		vals = append(vals, val)
	}
	if err := tftypes.ValidateValue(setType, vals); err != nil {
		return tftypes.NewValue(setType, tftypes.UnknownValue), err
	}
	return tftypes.NewValue(setType, vals), nil
}

// Equal returns true if the Set is considered semantically equal
// (same type and same value) to the attr.Value passed as an argument.
func (s Set) Equal(o attr.Value) bool {
	other, ok := o.(Set)
	if !ok {
		return false
	}
	if s.Unknown != other.Unknown {
		return false
	}
	if s.Null != other.Null {
		return false
	}
	if s.ElemType == nil && other.ElemType != nil {
		return false
	}
	if s.ElemType != nil && !s.ElemType.Equal(other.ElemType) {
		return false
	}
	if len(s.Elems) != len(other.Elems) {
		return false
	}
	for _, elem := range s.Elems {
		if !other.contains(elem) {
			return false
		}
	}
	return true
}

func (s Set) contains(v attr.Value) bool {
	for _, elem := range s.Elems {
		if elem.Equal(v) {
			return true
		}
	}

	return false
}

// IsNull returns true if the Set represents a null value.
func (s Set) IsNull() bool {
	return s.Null
}

// IsUnknown returns true if the Set represents a currently unknown value.
func (s Set) IsUnknown() bool {
	return s.Unknown
}

// String returns a human-readable representation of the Set value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (s Set) String() string {
	if s.Unknown {
		return attr.UnknownValueString
	}

	if s.Null {
		return attr.NullValueString
	}

	var res strings.Builder

	res.WriteString("[")
	for i, e := range s.Elems {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(e.String())
	}
	res.WriteString("]")

	return res.String()
}
