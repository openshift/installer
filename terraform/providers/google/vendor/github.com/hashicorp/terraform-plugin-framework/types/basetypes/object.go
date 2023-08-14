package basetypes

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ ObjectTypable  = ObjectType{}
	_ ObjectValuable = &ObjectValue{}
)

// ObjectTypable extends attr.Type for object types.
// Implement this interface to create a custom ObjectType type.
type ObjectTypable interface {
	attr.Type

	// ValueFromObject should convert the Object to an ObjectValuable type.
	ValueFromObject(context.Context, ObjectValue) (ObjectValuable, diag.Diagnostics)
}

// ObjectValuable extends attr.Value for object value types.
// Implement this interface to create a custom Object value type.
type ObjectValuable interface {
	attr.Value

	// ToObjectValue should convert the value type to an Object.
	ToObjectValue(ctx context.Context) (ObjectValue, diag.Diagnostics)
}

// ObjectType is an AttributeType representing an object.
type ObjectType struct {
	AttrTypes map[string]attr.Type
}

// WithAttributeTypes returns a new copy of the type with its attribute types
// set.
func (o ObjectType) WithAttributeTypes(typs map[string]attr.Type) attr.TypeWithAttributeTypes {
	return ObjectType{
		AttrTypes: typs,
	}
}

// AttributeTypes returns a copy of the type's attribute types.
func (o ObjectType) AttributeTypes() map[string]attr.Type {
	// Ensure callers cannot mutate the value
	result := make(map[string]attr.Type, len(o.AttrTypes))

	for key, value := range o.AttrTypes {
		result[key] = value
	}

	return result
}

// TerraformType returns the tftypes.Type that should be used to
// represent this type. This constrains what user input will be
// accepted and what kind of data can be set in state. The framework
// will use this to translate the AttributeType to something Terraform
// can understand.
func (o ObjectType) TerraformType(ctx context.Context) tftypes.Type {
	attributeTypes := map[string]tftypes.Type{}
	for k, v := range o.AttrTypes {
		attributeTypes[k] = v.TerraformType(ctx)
	}
	return tftypes.Object{
		AttributeTypes: attributeTypes,
	}
}

// ValueFromTerraform returns an attr.Value given a tftypes.Value.
// This is meant to convert the tftypes.Value into a more convenient Go
// type for the provider to consume the data with.
func (o ObjectType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewObjectNull(o.AttrTypes), nil
	}
	if !in.Type().Equal(o.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", o.TerraformType(ctx), in.Type())
	}
	if !in.IsKnown() {
		return NewObjectUnknown(o.AttrTypes), nil
	}
	if in.IsNull() {
		return NewObjectNull(o.AttrTypes), nil
	}
	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}
	err := in.As(&val)
	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := o.AttrTypes[k].ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, err
		}
		attributes[k] = a
	}
	// ValueFromTerraform above on each attribute should make this safe.
	// Otherwise, this will need to do some Diagnostics to error conversion.
	return NewObjectValueMust(o.AttrTypes, attributes), nil
}

// Equal returns true if `candidate` is also an ObjectType and has the same
// AttributeTypes.
func (o ObjectType) Equal(candidate attr.Type) bool {
	other, ok := candidate.(ObjectType)
	if !ok {
		return false
	}
	if len(other.AttrTypes) != len(o.AttrTypes) {
		return false
	}
	for k, v := range o.AttrTypes {
		attr, ok := other.AttrTypes[k]
		if !ok {
			return false
		}
		if !v.Equal(attr) {
			return false
		}
	}
	return true
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// object.
func (o ObjectType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	attrName, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to ObjectType", step)
	}

	attrType, ok := o.AttrTypes[string(attrName)]

	if !ok {
		return nil, fmt.Errorf("undefined attribute name %s in ObjectType", attrName)
	}

	return attrType, nil
}

// String returns a human-friendly description of the ObjectType.
func (o ObjectType) String() string {
	var res strings.Builder
	res.WriteString("types.ObjectType[")
	keys := make([]string, 0, len(o.AttrTypes))
	for k := range o.AttrTypes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for pos, key := range keys {
		if pos != 0 {
			res.WriteString(", ")
		}
		res.WriteString(`"` + key + `":`)
		res.WriteString(o.AttrTypes[key].String())
	}
	res.WriteString("]")
	return res.String()
}

// ValueType returns the Value type.
func (o ObjectType) ValueType(_ context.Context) attr.Value {
	return ObjectValue{
		attributeTypes: o.AttrTypes,
	}
}

// ValueFromObject returns an ObjectValuable type given an Object.
func (o ObjectType) ValueFromObject(_ context.Context, obj ObjectValue) (ObjectValuable, diag.Diagnostics) {
	return obj, nil
}

// NewObjectNull creates a Object with a null value. Determine whether the value is
// null via the Object type IsNull method.
func NewObjectNull(attributeTypes map[string]attr.Type) ObjectValue {
	return ObjectValue{
		attributeTypes: attributeTypes,
		state:          attr.ValueStateNull,
	}
}

// NewObjectUnknown creates a Object with an unknown value. Determine whether the
// value is unknown via the Object type IsUnknown method.
func NewObjectUnknown(attributeTypes map[string]attr.Type) ObjectValue {
	return ObjectValue{
		attributeTypes: attributeTypes,
		state:          attr.ValueStateUnknown,
	}
}

// NewObjectValue creates a Object with a known value. Access the value via the Object
// type ElementsAs method.
func NewObjectValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing Object Attribute Value",
				"While creating a Object value, a missing attribute value was detected. "+
					"A Object must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Object Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid Object Attribute Type",
				"While creating a Object value, an invalid attribute value was detected. "+
					"A Object must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Object Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("Object Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra Object Attribute Value",
				"While creating a Object value, an extra attribute value was detected. "+
					"A Object must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra Object Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewObjectUnknown(attributeTypes), diags
	}

	return ObjectValue{
		attributeTypes: attributeTypes,
		attributes:     attributes,
		state:          attr.ValueStateKnown,
	}, nil
}

// NewObjectValueFrom creates a Object with a known value, using reflection rules.
// The attributes must be a map of string attribute names to attribute values
// which can convert into the given attribute type or a struct with tfsdk field
// tags. Access the value via the Object type Elements or ElementsAs methods.
func NewObjectValueFrom(ctx context.Context, attributeTypes map[string]attr.Type, attributes any) (ObjectValue, diag.Diagnostics) {
	attrValue, diags := reflect.FromValue(
		ctx,
		ObjectType{AttrTypes: attributeTypes},
		attributes,
		path.Empty(),
	)

	if diags.HasError() {
		return NewObjectUnknown(attributeTypes), diags
	}

	m, ok := attrValue.(ObjectValue)

	// This should not happen, but ensure there is an error if it does.
	if !ok {
		diags.AddError(
			"Unable to Convert Object Value",
			"An unexpected result occurred when creating a Object using ObjectValueFrom. "+
				"This is an issue with terraform-plugin-framework and should be reported to the provider developers.",
		)
	}

	return m, diags
}

// NewObjectValueMust creates a Object with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Object
// type Elements or ElementsAs methods.
//
// This creation function is only recommended to create Object values which will
// not potentially effect practitioners, such as testing, or exhaustively
// tested provider logic.
func NewObjectValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ObjectValue {
	object, diags := NewObjectValue(attributeTypes, attributes)

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

		panic("ObjectValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

// ObjectValue represents an object
type ObjectValue struct {
	// attributes is the mapping of known attribute values in the Object.
	attributes map[string]attr.Value

	// attributeTypes is the type of the attributes in the Object.
	attributeTypes map[string]attr.Type

	// state represents whether the value is null, unknown, or known. The
	// zero-value is null.
	state attr.ValueState
}

// ObjectAsOptions is a collection of toggles to control the behavior of
// Object.As.
type ObjectAsOptions struct {
	// UnhandledNullAsEmpty controls what happens when As needs to put a
	// null value in a type that has no way to preserve that distinction.
	// When set to true, the type's empty value will be used.  When set to
	// false, an error will be returned.
	UnhandledNullAsEmpty bool

	// UnhandledUnknownAsEmpty controls what happens when As needs to put
	// an unknown value in a type that has no way to preserve that
	// distinction. When set to true, the type's empty value will be used.
	// When set to false, an error will be returned.
	UnhandledUnknownAsEmpty bool
}

// As populates `target` with the data in the ObjectValue, throwing an error if the
// data cannot be stored in `target`.
func (o ObjectValue) As(ctx context.Context, target interface{}, opts ObjectAsOptions) diag.Diagnostics {
	// we need a tftypes.Value for this Object to be able to use it with
	// our reflection code
	obj := ObjectType{AttrTypes: o.attributeTypes}
	val, err := o.ToTerraformValue(ctx)
	if err != nil {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Object Conversion Error",
				"An unexpected error was encountered trying to convert object. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			),
		}
	}
	return reflect.Into(ctx, obj, val, target, reflect.Options{
		UnhandledNullAsEmpty:    opts.UnhandledNullAsEmpty,
		UnhandledUnknownAsEmpty: opts.UnhandledUnknownAsEmpty,
	}, path.Empty())
}

// Attributes returns a copy of the mapping of known attribute values for the Object.
func (o ObjectValue) Attributes() map[string]attr.Value {
	// Ensure callers cannot mutate the internal attributes
	result := make(map[string]attr.Value, len(o.attributes))

	for name, value := range o.attributes {
		result[name] = value
	}

	return result
}

// AttributeTypes returns a copy of the mapping of attribute types for the Object.
func (o ObjectValue) AttributeTypes(_ context.Context) map[string]attr.Type {
	// Ensure callers cannot mutate the internal attribute types
	result := make(map[string]attr.Type, len(o.attributeTypes))

	for name, typ := range o.attributeTypes {
		result[name] = typ
	}

	return result
}

// Type returns an ObjectType with the same attribute types as `o`.
func (o ObjectValue) Type(ctx context.Context) attr.Type {
	return ObjectType{AttrTypes: o.AttributeTypes(ctx)}
}

// ToTerraformValue returns the data contained in the attr.Value as
// a tftypes.Value.
func (o ObjectValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := map[string]tftypes.Type{}
	for attr, typ := range o.AttributeTypes(ctx) {
		attrTypes[attr] = typ.TerraformType(ctx)
	}
	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch o.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, len(o.attributes))

		for name, v := range o.attributes {
			val, err := v.ToTerraformValue(ctx)

			if err != nil {
				return tftypes.NewValue(objectType, tftypes.UnknownValue), err
			}

			vals[name] = val
		}

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", o.state))
	}
}

// Equal returns true if the Object is considered semantically equal
// (same type and same value) to the attr.Value passed as an argument.
func (o ObjectValue) Equal(c attr.Value) bool {
	other, ok := c.(ObjectValue)

	if !ok {
		return false
	}

	if o.state != other.state {
		return false
	}

	if o.state != attr.ValueStateKnown {
		return true
	}

	if len(o.attributeTypes) != len(other.attributeTypes) {
		return false
	}

	for name, oAttributeType := range o.attributeTypes {
		otherAttributeType, ok := other.attributeTypes[name]

		if !ok {
			return false
		}

		if !oAttributeType.Equal(otherAttributeType) {
			return false
		}
	}

	if len(o.attributes) != len(other.attributes) {
		return false
	}

	for name, oAttribute := range o.attributes {
		otherAttribute, ok := other.attributes[name]

		if !ok {
			return false
		}

		if !oAttribute.Equal(otherAttribute) {
			return false
		}
	}

	return true
}

// IsNull returns true if the Object represents a null value.
func (o ObjectValue) IsNull() bool {
	return o.state == attr.ValueStateNull
}

// IsUnknown returns true if the Object represents a currently unknown value.
func (o ObjectValue) IsUnknown() bool {
	return o.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the Object value.
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (o ObjectValue) String() string {
	if o.IsUnknown() {
		return attr.UnknownValueString
	}

	if o.IsNull() {
		return attr.NullValueString
	}

	// We want the output to be consistent, so we sort the output by key
	keys := make([]string, 0, len(o.Attributes()))
	for k := range o.Attributes() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var res strings.Builder

	res.WriteString("{")
	for i, k := range keys {
		if i != 0 {
			res.WriteString(",")
		}
		res.WriteString(fmt.Sprintf(`"%s":%s`, k, o.Attributes()[k].String()))
	}
	res.WriteString("}")

	return res.String()
}

// ToObjectValue returns the Object.
func (o ObjectValue) ToObjectValue(context.Context) (ObjectValue, diag.Diagnostics) {
	return o, nil
}
