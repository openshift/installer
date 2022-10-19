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
	_ attr.Type  = ListType{}
	_ attr.Value = &List{}
)

// ListType is an AttributeType representing a list of values. All values must
// be of the same type, which the provider must specify as the ElemType
// property.
type ListType struct {
	ElemType attr.Type
}

// ElementType returns the attr.Type elements will be created from.
func (l ListType) ElementType() attr.Type {
	return l.ElemType
}

// WithElementType returns a ListType that is identical to `l`, but with the
// element type set to `typ`.
func (l ListType) WithElementType(typ attr.Type) attr.TypeWithElementType {
	return ListType{ElemType: typ}
}

// TerraformType returns the tftypes.Type that should be used to
// represent this type. This constrains what user input will be
// accepted and what kind of data can be set in state. The framework
// will use this to translate the AttributeType to something Terraform
// can understand.
func (l ListType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.List{
		ElementType: l.ElemType.TerraformType(ctx),
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value.
// This is meant to convert the tftypes.Value into a more convenient Go
// type for the provider to consume the data with.
func (l ListType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	list := List{
		ElemType: l.ElemType,
	}
	if in.Type() == nil {
		list.Null = true
		return list, nil
	}
	if !in.Type().Equal(l.TerraformType(ctx)) {
		return nil, fmt.Errorf("can't use %s as value of List with ElementType %T, can only use %s values", in.String(), l.ElemType, l.ElemType.TerraformType(ctx).String())
	}
	if !in.IsKnown() {
		list.Unknown = true
		return list, nil
	}
	if in.IsNull() {
		list.Null = true
		return list, nil
	}
	val := []tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}
	elems := make([]attr.Value, 0, len(val))
	for _, elem := range val {
		av, err := l.ElemType.ValueFromTerraform(ctx, elem)
		if err != nil {
			return nil, err
		}
		elems = append(elems, av)
	}
	list.Elems = elems
	return list, nil
}

// Equal returns true if `o` is also a ListType and has the same ElemType.
func (l ListType) Equal(o attr.Type) bool {
	if l.ElemType == nil {
		return false
	}
	other, ok := o.(ListType)
	if !ok {
		return false
	}
	return l.ElemType.Equal(other.ElemType)
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// list.
func (l ListType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	if _, ok := step.(tftypes.ElementKeyInt); !ok {
		return nil, fmt.Errorf("cannot apply step %T to ListType", step)
	}

	return l.ElemType, nil
}

// String returns a human-friendly description of the ListType.
func (l ListType) String() string {
	return "types.ListType[" + l.ElemType.String() + "]"
}

// Validate validates all elements of the list that are of type
// xattr.TypeWithValidate.
func (l ListType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.List{}) {
		err := fmt.Errorf("expected List value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"List Type Validation Error",
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
			"List Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	validatableType, isValidatable := l.ElemType.(xattr.TypeWithValidate)
	if !isValidatable {
		return diags
	}

	for index, elem := range elems {
		if !elem.IsFullyKnown() {
			continue
		}
		diags = append(diags, validatableType.Validate(ctx, elem, path.AtListIndex(index))...)
	}

	return diags
}

// ValueType returns the Value type.
func (t ListType) ValueType(_ context.Context) attr.Value {
	return List{
		ElemType: t.ElemType,
	}
}

// List represents a list of attr.Values, all of the same type, indicated
// by ElemType.
type List struct {
	// Unknown will be set to true if the entire list is an unknown value.
	// If only some of the elements in the list are unknown, their known or
	// unknown status will be represented however that attr.Value
	// surfaces that information. The List's Unknown property only tracks
	// if the number of elements in a List is known, not whether the
	// elements that are in the list are known.
	Unknown bool

	// Null will be set to true if the list is null, either because it was
	// omitted from the configuration, state, or plan, or because it was
	// explicitly set to null.
	Null bool

	// Elems are the elements in the list.
	Elems []attr.Value

	// ElemType is the tftypes.Type of the elements in the list. All
	// elements in the list must be of this type.
	ElemType attr.Type
}

// ElementsAs populates `target` with the elements of the List, throwing an
// error if the elements cannot be stored in `target`.
func (l List) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
	// we need a tftypes.Value for this List to be able to use it with our
	// reflection code
	values, err := l.ToTerraformValue(ctx)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"List Element Conversion Error",
				"An unexpected error was encountered trying to convert list elements. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			),
		}
	}
	return reflect.Into(ctx, ListType{ElemType: l.ElemType}, values, target, reflect.Options{
		UnhandledNullAsEmpty:    allowUnhandled,
		UnhandledUnknownAsEmpty: allowUnhandled,
	}, path.Empty())
}

// Type returns a ListType with the same element type as `l`.
func (l List) Type(ctx context.Context) attr.Type {
	return ListType{ElemType: l.ElemType}
}

// ToTerraformValue returns the data contained in the List as a tftypes.Value.
func (l List) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if l.ElemType == nil {
		return tftypes.Value{}, fmt.Errorf("cannot convert List to tftypes.Value if ElemType field is not set")
	}
	listType := tftypes.List{ElementType: l.ElemType.TerraformType(ctx)}
	if l.Unknown {
		return tftypes.NewValue(listType, tftypes.UnknownValue), nil
	}
	if l.Null {
		return tftypes.NewValue(listType, nil), nil
	}
	vals := make([]tftypes.Value, 0, len(l.Elems))
	for _, elem := range l.Elems {
		val, err := elem.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(listType, tftypes.UnknownValue), err
		}
		vals = append(vals, val)
	}
	if err := tftypes.ValidateValue(listType, vals); err != nil {
		return tftypes.NewValue(listType, tftypes.UnknownValue), err
	}
	return tftypes.NewValue(listType, vals), nil
}

// Equal returns true if the List is considered semantically equal
// (same type and same value) to the attr.Value passed as an argument.
func (l List) Equal(o attr.Value) bool {
	other, ok := o.(List)
	if !ok {
		return false
	}
	if l.Unknown != other.Unknown {
		return false
	}
	if l.Null != other.Null {
		return false
	}
	if l.ElemType == nil && other.ElemType != nil {
		return false
	}
	if l.ElemType != nil && !l.ElemType.Equal(other.ElemType) {
		return false
	}
	if len(l.Elems) != len(other.Elems) {
		return false
	}
	for pos, lElem := range l.Elems {
		oElem := other.Elems[pos]
		if !lElem.Equal(oElem) {
			return false
		}
	}
	return true
}

// IsNull returns true if the List represents a null value.
func (l List) IsNull() bool {
	return l.Null
}

// IsUnknown returns true if the List represents a currently unknown value.
func (l List) IsUnknown() bool {
	return l.Unknown
}

// String returns a human-readable representation of the List value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (l List) String() string {
	if l.Unknown {
		return attr.UnknownValueString
	}

	if l.Null {
		return attr.NullValueString
	}

	var res strings.Builder

	res.WriteString("[")
	for i, e := range l.Elems {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(e.String())
	}
	res.WriteString("]")

	return res.String()
}
