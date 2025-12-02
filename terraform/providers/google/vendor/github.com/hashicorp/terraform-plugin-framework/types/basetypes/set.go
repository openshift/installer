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
	_ SetTypable             = SetType{}
	_ xattr.TypeWithValidate = SetType{}
	_ SetValuable            = &SetValue{}
)

// SetTypable extends attr.Type for set types.
// Implement this interface to create a custom SetType type.
type SetTypable interface {
	attr.Type

	// ValueFromSet should convert the Set to a SetValuable type.
	ValueFromSet(context.Context, SetValue) (SetValuable, diag.Diagnostics)
}

// SetValuable extends attr.Value for set value types.
// Implement this interface to create a custom Set value type.
type SetValuable interface {
	attr.Value

	// ToSetValue should convert the value type to a Set.
	ToSetValue(ctx context.Context) (SetValue, diag.Diagnostics)
}

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
	if in.Type() == nil {
		return NewSetNull(st.ElemType), nil
	}
	if !in.Type().Equal(st.TerraformType(ctx)) {
		return nil, fmt.Errorf("can't use %s as value of Set with ElementType %T, can only use %s values", in.String(), st.ElemType, st.ElemType.TerraformType(ctx).String())
	}
	if !in.IsKnown() {
		return NewSetUnknown(st.ElemType), nil
	}
	if in.IsNull() {
		return NewSetNull(st.ElemType), nil
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
	// ValueFromTerraform above on each element should make this safe.
	// Otherwise, this will need to do some Diagnostics to error conversion.
	return NewSetValueMust(st.ElemType, elems), nil
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
func (st SetType) ValueType(_ context.Context) attr.Value {
	return SetValue{
		elementType: st.ElemType,
	}
}

// ValueFromSet returns a SetValuable type given a Set.
func (st SetType) ValueFromSet(_ context.Context, set SetValue) (SetValuable, diag.Diagnostics) {
	return set, nil
}

// NewSetNull creates a Set with a null value. Determine whether the value is
// null via the Set type IsNull method.
func NewSetNull(elementType attr.Type) SetValue {
	return SetValue{
		elementType: elementType,
		state:       attr.ValueStateNull,
	}
}

// NewSetUnknown creates a Set with an unknown value. Determine whether the
// value is unknown via the Set type IsUnknown method.
func NewSetUnknown(elementType attr.Type) SetValue {
	return SetValue{
		elementType: elementType,
		state:       attr.ValueStateUnknown,
	}
}

// NewSetValue creates a Set with a known value. Access the value via the Set
// type Elements or ElementsAs methods.
func NewSetValue(elementType attr.Type, elements []attr.Value) (SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for idx, element := range elements {
		if !elementType.Equal(element.Type(ctx)) {
			diags.AddError(
				"Invalid Set Element Type",
				"While creating a Set value, an invalid element was detected. "+
					"A Set must use the single, given element type. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Set Element Type: %s\n", elementType.String())+
					fmt.Sprintf("Set Index (%d) Element Type: %s", idx, element.Type(ctx)),
			)
		}
	}

	if diags.HasError() {
		return NewSetUnknown(elementType), diags
	}

	return SetValue{
		elementType: elementType,
		elements:    elements,
		state:       attr.ValueStateKnown,
	}, nil
}

// NewSetValueFrom creates a Set with a known value, using reflection rules.
// The elements must be a slice which can convert into the given element type.
// Access the value via the Set type Elements or ElementsAs methods.
func NewSetValueFrom(ctx context.Context, elementType attr.Type, elements any) (SetValue, diag.Diagnostics) {
	attrValue, diags := reflect.FromValue(
		ctx,
		SetType{ElemType: elementType},
		elements,
		path.Empty(),
	)

	if diags.HasError() {
		return NewSetUnknown(elementType), diags
	}

	set, ok := attrValue.(SetValue)

	// This should not happen, but ensure there is an error if it does.
	if !ok {
		diags.AddError(
			"Unable to Convert Set Value",
			"An unexpected result occurred when creating a Set using SetValueFrom. "+
				"This is an issue with terraform-plugin-framework and should be reported to the provider developers.",
		)
	}

	return set, diags
}

// NewSetValueMust creates a Set with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Set
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create Set values which will
// not potentially effect practitioners, such as testing, or exhaustively
// tested provider logic.
func NewSetValueMust(elementType attr.Type, elements []attr.Value) SetValue {
	set, diags := NewSetValue(elementType, elements)

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

		panic("SetValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return set
}

// SetValue represents a set of attr.Value, all of the same type,
// indicated by ElemType.
type SetValue struct {
	// elements is the collection of known values in the Set.
	elements []attr.Value

	// elementType is the type of the elements in the Set.
	elementType attr.Type

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// Elements returns a copy of the collection of elements for the Set.
func (s SetValue) Elements() []attr.Value {
	// Ensure callers cannot mutate the internal elements
	result := make([]attr.Value, 0, len(s.elements))
	result = append(result, s.elements...)

	return result
}

// ElementsAs populates `target` with the elements of the SetValue, throwing an
// error if the elements cannot be stored in `target`.
func (s SetValue) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
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

// ElementType returns the element type for the Set.
func (s SetValue) ElementType(_ context.Context) attr.Type {
	return s.elementType
}

// Type returns a SetType with the same element type as `s`.
func (s SetValue) Type(ctx context.Context) attr.Type {
	return SetType{ElemType: s.ElementType(ctx)}
}

// ToTerraformValue returns the data contained in the Set as a tftypes.Value.
func (s SetValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	setType := tftypes.Set{ElementType: s.ElementType(ctx).TerraformType(ctx)}

	switch s.state {
	case attr.ValueStateKnown:
		vals := make([]tftypes.Value, 0, len(s.elements))

		for _, elem := range s.elements {
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
	case attr.ValueStateNull:
		return tftypes.NewValue(setType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(setType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Set state in ToTerraformValue: %s", s.state))
	}
}

// Equal returns true if the Set is considered semantically equal
// (same type and same value) to the attr.Value passed as an argument.
func (s SetValue) Equal(o attr.Value) bool {
	other, ok := o.(SetValue)

	if !ok {
		return false
	}

	if !s.elementType.Equal(other.elementType) {
		return false
	}

	if s.state != other.state {
		return false
	}

	if s.state != attr.ValueStateKnown {
		return true
	}

	if len(s.elements) != len(other.elements) {
		return false
	}

	for _, elem := range s.elements {
		if !other.contains(elem) {
			return false
		}
	}

	return true
}

func (s SetValue) contains(v attr.Value) bool {
	for _, elem := range s.Elements() {
		if elem.Equal(v) {
			return true
		}
	}

	return false
}

// IsNull returns true if the Set represents a null value.
func (s SetValue) IsNull() bool {
	return s.state == attr.ValueStateNull
}

// IsUnknown returns true if the Set represents a currently unknown value.
// Returns false if the Set has a known number of elements, even if all are
// unknown values.
func (s SetValue) IsUnknown() bool {
	return s.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Set value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (s SetValue) String() string {
	if s.IsUnknown() {
		return attr.UnknownValueString
	}

	if s.IsNull() {
		return attr.NullValueString
	}

	var res strings.Builder

	res.WriteString("[")
	for i, e := range s.Elements() {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(e.String())
	}
	res.WriteString("]")

	return res.String()
}

// ToSetValue returns the Set.
func (s SetValue) ToSetValue(context.Context) (SetValue, diag.Diagnostics) {
	return s, nil
}
