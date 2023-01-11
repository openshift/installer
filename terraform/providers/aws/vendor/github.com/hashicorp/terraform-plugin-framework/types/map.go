package types

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var (
	_ attr.Type  = MapType{}
	_ attr.Value = &Map{}
)

// MapType is an AttributeType representing a map of values. All values must
// be of the same type, which the provider must specify as the ElemType
// property. Keys will always be strings.
type MapType struct {
	ElemType attr.Type
}

// WithElementType returns a new copy of the type with its element type set.
func (m MapType) WithElementType(typ attr.Type) attr.TypeWithElementType {
	return MapType{
		ElemType: typ,
	}
}

// ElementType returns the type's element type.
func (m MapType) ElementType() attr.Type {
	return m.ElemType
}

// TerraformType returns the tftypes.Type that should be used to represent this
// type. This constrains what user input will be accepted and what kind of data
// can be set in state. The framework will use this to translate the
// AttributeType to something Terraform can understand.
func (m MapType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.Map{
		ElementType: m.ElemType.TerraformType(ctx),
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value. This is
// meant to convert the tftypes.Value into a more convenient Go type for the
// provider to consume the data with.
func (m MapType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	ma := Map{
		ElemType: m.ElemType,
	}
	if in.Type() == nil {
		ma.Null = true
		return ma, nil
	}
	if !in.Type().Is(tftypes.Map{}) {
		return nil, fmt.Errorf("can't use %s as value of Map, can only use tftypes.Map values", in.String())
	}
	if !in.Type().Equal(tftypes.Map{ElementType: m.ElemType.TerraformType(ctx)}) {
		return nil, fmt.Errorf("can't use %s as value of Map with ElementType %T, can only use %s values", in.String(), m.ElemType, m.ElemType.TerraformType(ctx).String())
	}
	if !in.IsKnown() {
		ma.Unknown = true
		return ma, nil
	}
	if in.IsNull() {
		ma.Null = true
		return ma, nil
	}
	val := map[string]tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}
	elems := make(map[string]attr.Value, len(val))
	for key, elem := range val {
		av, err := m.ElemType.ValueFromTerraform(ctx, elem)
		if err != nil {
			return nil, err
		}
		elems[key] = av
	}
	ma.Elems = elems
	return ma, nil
}

// Equal returns true if `o` is also a MapType and has the same ElemType.
func (m MapType) Equal(o attr.Type) bool {
	if m.ElemType == nil {
		return false
	}
	other, ok := o.(MapType)
	if !ok {
		return false
	}
	return m.ElemType.Equal(other.ElemType)
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// map.
func (m MapType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	if _, ok := step.(tftypes.ElementKeyString); !ok {
		return nil, fmt.Errorf("cannot apply step %T to MapType", step)
	}

	return m.ElemType, nil
}

// String returns a human-friendly description of the MapType.
func (m MapType) String() string {
	return "types.MapType[" + m.ElemType.String() + "]"
}

// Validate validates all elements of the map that are of type
// xattr.TypeWithValidate.
func (m MapType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.Map{}) {
		err := fmt.Errorf("expected Map value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Map Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var elems map[string]tftypes.Value

	if err := in.As(&elems); err != nil {
		diags.AddAttributeError(
			path,
			"Map Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	validatableType, isValidatable := m.ElemType.(xattr.TypeWithValidate)
	if !isValidatable {
		return diags
	}

	for index, elem := range elems {
		if !elem.IsFullyKnown() {
			continue
		}
		diags = append(diags, validatableType.Validate(ctx, elem, path.AtMapKey(index))...)
	}

	return diags
}

// ValueType returns the Value type.
func (t MapType) ValueType(_ context.Context) attr.Value {
	return Map{
		ElemType: t.ElemType,
	}
}

// Map represents a map of attr.Values, all of the same type, indicated by
// ElemType. Keys for the map will always be strings.
type Map struct {
	// Unknown will be set to true if the entire map is an unknown value.
	// If only some of the elements in the map are unknown, their known or
	// unknown status will be represented however that attr.Value
	// surfaces that information. The Map's Unknown property only tracks if
	// the number of elements in a Map is known, not whether the elements
	// that are in the map are known.
	Unknown bool

	// Null will be set to true if the map is null, either because it was
	// omitted from the configuration, state, or plan, or because it was
	// explicitly set to null.
	Null bool

	// Elems are the elements in the map.
	Elems map[string]attr.Value

	// ElemType is the AttributeType of the elements in the map. All
	// elements in the map must be of this type.
	ElemType attr.Type
}

// ElementsAs populates `target` with the elements of the Map, throwing an
// error if the elements cannot be stored in `target`.
func (m Map) ElementsAs(ctx context.Context, target interface{}, allowUnhandled bool) diag.Diagnostics {
	// we need a tftypes.Value for this Map to be able to use it with our
	// reflection code
	val, err := m.ToTerraformValue(ctx)
	if err != nil {
		err := fmt.Errorf("error getting Terraform value for map: %w", err)
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Map Conversion Error",
				"An unexpected error was encountered trying to convert the map into an equivalent Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			),
		}
	}

	return reflect.Into(ctx, MapType{ElemType: m.ElemType}, val, target, reflect.Options{
		UnhandledNullAsEmpty:    allowUnhandled,
		UnhandledUnknownAsEmpty: allowUnhandled,
	}, path.Empty())
}

// Type returns a MapType with the same element type as `m`.
func (m Map) Type(ctx context.Context) attr.Type {
	return MapType{ElemType: m.ElemType}
}

// ToTerraformValue returns the data contained in the List as a tftypes.Value.
func (m Map) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if m.ElemType == nil {
		return tftypes.Value{}, fmt.Errorf("cannot convert Map to tftypes.Value if ElemType field is not set")
	}
	mapType := tftypes.Map{ElementType: m.ElemType.TerraformType(ctx)}
	if m.Unknown {
		return tftypes.NewValue(mapType, tftypes.UnknownValue), nil
	}
	if m.Null {
		return tftypes.NewValue(mapType, nil), nil
	}
	vals := make(map[string]tftypes.Value, len(m.Elems))
	for key, elem := range m.Elems {
		val, err := elem.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(mapType, tftypes.UnknownValue), err
		}
		vals[key] = val
	}
	if err := tftypes.ValidateValue(mapType, vals); err != nil {
		return tftypes.NewValue(mapType, tftypes.UnknownValue), err
	}
	return tftypes.NewValue(mapType, vals), nil
}

// Equal returns true if the Map is considered semantically equal
// (same type and same value) to the attr.Value passed as an argument.
func (m Map) Equal(o attr.Value) bool {
	other, ok := o.(Map)
	if !ok {
		return false
	}
	if m.Unknown != other.Unknown {
		return false
	}
	if m.Null != other.Null {
		return false
	}
	if !m.ElemType.Equal(other.ElemType) {
		return false
	}
	if len(m.Elems) != len(other.Elems) {
		return false
	}
	for key, mElem := range m.Elems {
		oElem, ok := other.Elems[key]
		if !ok {
			return false
		}
		if !mElem.Equal(oElem) {
			return false
		}
	}
	return true
}

// IsNull returns true if the Map represents a null value.
func (m Map) IsNull() bool {
	return m.Null
}

// IsUnknown returns true if the Map represents a currently unknown value.
func (m Map) IsUnknown() bool {
	return m.Unknown
}

// String returns a human-readable representation of the Map value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (m Map) String() string {
	if m.Unknown {
		return attr.UnknownValueString
	}

	if m.Null {
		return attr.NullValueString
	}

	// We want the output to be consistent, so we sort the output by key
	keys := make([]string, 0, len(m.Elems))
	for k := range m.Elems {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var res strings.Builder

	res.WriteString("{")
	for i, k := range keys {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(fmt.Sprintf("%q:%s", k, m.Elems[k].String()))
	}
	res.WriteString("}")

	return res.String()
}
