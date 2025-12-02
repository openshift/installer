package basetypes

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
	_ ListTypable  = ListType{}
	_ ListValuable = &ListValue{}
)

// ListTypable extends attr.Type for list types.
// Implement this interface to create a custom ListType type.
type ListTypable interface {
	attr.Type

	// ValueFromList should convert the List to a ListValuable type.
	ValueFromList(context.Context, ListValue) (ListValuable, diag.Diagnostics)
}

// ListValuable extends attr.Value for list value types.
// Implement this interface to create a custom List value type.
type ListValuable interface {
	attr.Value

	// ToListValue should convert the value type to a List.
	ToListValue(ctx context.Context) (ListValue, diag.Diagnostics)
}

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
	if in.Type() == nil {
		return NewListNull(l.ElemType), nil
	}
	if !in.Type().Equal(l.TerraformType(ctx)) {
		return nil, fmt.Errorf("can't use %s as value of List with ElementType %T, can only use %s values", in.String(), l.ElemType, l.ElemType.TerraformType(ctx).String())
	}
	if !in.IsKnown() {
		return NewListUnknown(l.ElemType), nil
	}
	if in.IsNull() {
		return NewListNull(l.ElemType), nil
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
	// ValueFromTerraform above on each element should make this safe.
	// Otherwise, this will need to do some Diagnostics to error conversion.
	return NewListValueMust(l.ElemType, elems), nil
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
func (l ListType) ValueType(_ context.Context) attr.Value {
	return ListValue{
		elementType: l.ElemType,
	}
}

// ValueFromList returns a ListValuable type given a List.
func (l ListType) ValueFromList(_ context.Context, list ListValue) (ListValuable, diag.Diagnostics) {
	return list, nil
}

// NewListNull creates a List with a null value. Determine whether the value is
// null via the List type IsNull method.
func NewListNull(elementType attr.Type) ListValue {
	return ListValue{
		elementType: elementType,
		state:       attr.ValueStateNull,
	}
}

// NewListUnknown creates a List with an unknown value. Determine whether the
// value is unknown via the List type IsUnknown method.
func NewListUnknown(elementType attr.Type) ListValue {
	return ListValue{
		elementType: elementType,
		state:       attr.ValueStateUnknown,
	}
}

// NewListValue creates a List with a known value. Access the value via the List
// type Elements or ElementsAs methods.
func NewListValue(elementType attr.Type, elements []attr.Value) (ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for idx, element := range elements {
		if !elementType.Equal(element.Type(ctx)) {
			diags.AddError(
				"Invalid List Element Type",
				"While creating a List value, an invalid element was detected. "+
					"A List must use the single, given element type. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("List Element Type: %s\n", elementType.String())+
					fmt.Sprintf("List Index (%d) Element Type: %s", idx, element.Type(ctx)),
			)
		}
	}

	if diags.HasError() {
		return NewListUnknown(elementType), diags
	}

	return ListValue{
		elementType: elementType,
		elements:    elements,
		state:       attr.ValueStateKnown,
	}, nil
}

// NewListValueFrom creates a List with a known value, using reflection rules.
// The elements must be a slice which can convert into the given element type.
// Access the value via the List type Elements or ElementsAs methods.
func NewListValueFrom(ctx context.Context, elementType attr.Type, elements any) (ListValue, diag.Diagnostics) {
	attrValue, diags := reflect.FromValue(
		ctx,
		ListType{ElemType: elementType},
		elements,
		path.Empty(),
	)

	if diags.HasError() {
		return NewListUnknown(elementType), diags
	}

	list, ok := attrValue.(ListValue)

	// This should not happen, but ensure there is an error if it does.
	if !ok {
		diags.AddError(
			"Unable to Convert List Value",
			"An unexpected result occurred when creating a List using NewListValueFrom. "+
				"This is an issue with terraform-plugin-framework and should be reported to the provider developers.",
		)
	}

	return list, diags
}

// NewListValueMust creates a List with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the List
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create List values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func NewListValueMust(elementType attr.Type, elements []attr.Value) ListValue {
	list, diags := NewListValue(elementType, elements)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewListValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return list
}

// ListValue represents a list of attr.Values, all of the same type, indicated
// by ElemType.
type ListValue struct {
	// elements is the collection of known values in the List.
	elements []attr.Value

	// elementType is the type of the elements in the List.
	elementType attr.Type

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// Elements returns a copy of the collection of elements for the List.
func (l ListValue) Elements() []attr.Value {
	// Ensure callers cannot mutate the internal elements
	result := make([]attr.Value, 0, len(l.elements))
	result = append(result, l.elements...)

	return result
}

// ElementsAs populates `target` with the elements of the ListValue, throwing an
// error if the elements cannot be stored in `target`.
func (l ListValue) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
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
	return reflect.Into(ctx, ListType{ElemType: l.elementType}, values, target, reflect.Options{
		UnhandledNullAsEmpty:    allowUnhandled,
		UnhandledUnknownAsEmpty: allowUnhandled,
	}, path.Empty())
}

// ElementType returns the element type for the List.
func (l ListValue) ElementType(_ context.Context) attr.Type {
	return l.elementType
}

// Type returns a ListType with the same element type as `l`.
func (l ListValue) Type(ctx context.Context) attr.Type {
	return ListType{ElemType: l.ElementType(ctx)}
}

// ToTerraformValue returns the data contained in the List as a tftypes.Value.
func (l ListValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	listType := tftypes.List{ElementType: l.ElementType(ctx).TerraformType(ctx)}

	switch l.state {
	case attr.ValueStateKnown:
		vals := make([]tftypes.Value, 0, len(l.elements))

		for _, elem := range l.elements {
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
	case attr.ValueStateNull:
		return tftypes.NewValue(listType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(listType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled List state in ToTerraformValue: %s", l.state))
	}
}

// Equal returns true if the List is considered semantically equal
// (same type and same value) to the attr.Value passed as an argument.
func (l ListValue) Equal(o attr.Value) bool {
	other, ok := o.(ListValue)

	if !ok {
		return false
	}

	if !l.elementType.Equal(other.elementType) {
		return false
	}

	if l.state != other.state {
		return false
	}

	if l.state != attr.ValueStateKnown {
		return true
	}

	if len(l.elements) != len(other.elements) {
		return false
	}

	for idx, lElem := range l.elements {
		otherElem := other.elements[idx]

		if !lElem.Equal(otherElem) {
			return false
		}
	}

	return true
}

// IsNull returns true if the List represents a null value.
func (l ListValue) IsNull() bool {
	return l.state == attr.ValueStateNull
}

// IsUnknown returns true if the List represents a currently unknown value.
// Returns false if the List has a known number of elements, even if all are
// unknown values.
func (l ListValue) IsUnknown() bool {
	return l.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the List value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (l ListValue) String() string {
	if l.IsUnknown() {
		return attr.UnknownValueString
	}

	if l.IsNull() {
		return attr.NullValueString
	}

	var res strings.Builder

	res.WriteString("[")
	for i, e := range l.Elements() {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(e.String())
	}
	res.WriteString("]")

	return res.String()
}

// ToListValue returns the List.
func (l ListValue) ToListValue(context.Context) (ListValue, diag.Diagnostics) {
	return l, nil
}
