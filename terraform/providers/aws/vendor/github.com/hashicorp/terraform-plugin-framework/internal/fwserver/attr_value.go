package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func coerceListValue(schemaPath path.Path, value attr.Value) (types.List, diag.Diagnostics) {
	list, ok := value.(types.List)

	if !ok {
		return types.List{Null: true}, diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return list, nil
}

func coerceMapValue(schemaPath path.Path, value attr.Value) (types.Map, diag.Diagnostics) {
	m, ok := value.(types.Map)

	if !ok {
		return types.Map{Null: true}, diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return m, nil
}

func coerceObjectValue(schemaPath path.Path, value attr.Value) (types.Object, diag.Diagnostics) {
	object, ok := value.(types.Object)

	if !ok {
		return types.Object{Null: true}, diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return object, nil
}

func coerceSetValue(schemaPath path.Path, value attr.Value) (types.Set, diag.Diagnostics) {
	set, ok := value.(types.Set)

	if !ok {
		return types.Set{Null: true}, diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return set, nil
}

func listElemObject(ctx context.Context, schemaPath path.Path, list types.List, index int, description fwschemadata.DataDescription) (types.Object, diag.Diagnostics) {
	if list.IsNull() {
		return listElemObjectFromTerraformValue(ctx, schemaPath, list, description, nil)
	}

	if list.IsUnknown() {
		return listElemObjectFromTerraformValue(ctx, schemaPath, list, description, tftypes.UnknownValue)
	}

	if index >= len(list.Elems) {
		return listElemObjectFromTerraformValue(ctx, schemaPath, list, description, nil)
	}

	return coerceObjectValue(schemaPath, list.Elems[index])
}

func listElemObjectFromTerraformValue(ctx context.Context, schemaPath path.Path, list types.List, description fwschemadata.DataDescription, tfValue any) (types.Object, diag.Diagnostics) {
	elemValue, err := list.ElemType.ValueFromTerraform(ctx, tftypes.NewValue(list.ElemType.TerraformType(ctx), tfValue))

	if err != nil {
		return types.Object{Null: true}, diag.Diagnostics{
			attributePlanModificationValueError(ctx, list, description, err),
		}
	}

	return coerceObjectValue(schemaPath, elemValue)
}

func mapElemObject(ctx context.Context, schemaPath path.Path, m types.Map, key string, description fwschemadata.DataDescription) (types.Object, diag.Diagnostics) {
	if m.IsNull() {
		return mapElemObjectFromTerraformValue(ctx, schemaPath, m, description, nil)
	}

	if m.IsUnknown() {
		return mapElemObjectFromTerraformValue(ctx, schemaPath, m, description, tftypes.UnknownValue)
	}

	elemValue, ok := m.Elems[key]

	if !ok {
		return mapElemObjectFromTerraformValue(ctx, schemaPath, m, description, nil)
	}

	return coerceObjectValue(schemaPath, elemValue)
}

func mapElemObjectFromTerraformValue(ctx context.Context, schemaPath path.Path, m types.Map, description fwschemadata.DataDescription, tfValue any) (types.Object, diag.Diagnostics) {
	elemValue, err := m.ElemType.ValueFromTerraform(ctx, tftypes.NewValue(m.ElemType.TerraformType(ctx), tfValue))

	if err != nil {
		return types.Object{Null: true}, diag.Diagnostics{
			attributePlanModificationValueError(ctx, m, description, err),
		}
	}

	return coerceObjectValue(schemaPath, elemValue)
}

func objectAttributeValue(ctx context.Context, object types.Object, attributeName string, description fwschemadata.DataDescription) (attr.Value, diag.Diagnostics) {
	if object.IsNull() {
		return objectAttributeValueFromTerraformValue(ctx, object, attributeName, description, nil)
	}

	if object.IsUnknown() {
		return objectAttributeValueFromTerraformValue(ctx, object, attributeName, description, tftypes.UnknownValue)
	}

	// A panic here indicates a bug somewhere else in the framework or an
	// invalid test case.
	return object.Attrs[attributeName], nil
}

func objectAttributeValueFromTerraformValue(ctx context.Context, object types.Object, attributeName string, description fwschemadata.DataDescription, tfValue any) (attr.Value, diag.Diagnostics) {
	// A panic here indicates a bug somewhere else in the framework or an
	// invalid test case.
	attrType := object.AttrTypes[attributeName]

	elemValue, err := attrType.ValueFromTerraform(ctx, tftypes.NewValue(attrType.TerraformType(ctx), tfValue))

	if err != nil {
		return nil, diag.Diagnostics{
			attributePlanModificationValueError(ctx, object, description, err),
		}
	}

	return elemValue, nil
}

func setElemObject(ctx context.Context, schemaPath path.Path, set types.Set, index int, description fwschemadata.DataDescription) (types.Object, diag.Diagnostics) {
	if set.IsNull() {
		return setElemObjectFromTerraformValue(ctx, schemaPath, set, description, nil)
	}

	if set.IsUnknown() {
		return setElemObjectFromTerraformValue(ctx, schemaPath, set, description, tftypes.UnknownValue)
	}

	if index >= len(set.Elems) {
		return setElemObjectFromTerraformValue(ctx, schemaPath, set, description, nil)
	}

	return coerceObjectValue(schemaPath, set.Elems[index])
}

func setElemObjectFromTerraformValue(ctx context.Context, schemaPath path.Path, set types.Set, description fwschemadata.DataDescription, tfValue any) (types.Object, diag.Diagnostics) {
	elemValue, err := set.ElemType.ValueFromTerraform(ctx, tftypes.NewValue(set.ElemType.TerraformType(ctx), tfValue))

	if err != nil {
		return types.Object{Null: true}, diag.Diagnostics{
			attributePlanModificationValueError(ctx, set, description, err),
		}
	}

	return coerceObjectValue(schemaPath, elemValue)
}
