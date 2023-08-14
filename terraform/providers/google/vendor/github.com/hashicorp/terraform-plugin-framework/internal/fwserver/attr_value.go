package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func coerceListValue(ctx context.Context, schemaPath path.Path, value attr.Value) (types.List, diag.Diagnostics) {
	listVal, ok := value.(basetypes.ListValuable)

	if !ok {
		return types.ListNull(nil), diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return listVal.ToListValue(ctx)
}

func coerceMapValue(ctx context.Context, schemaPath path.Path, value attr.Value) (types.Map, diag.Diagnostics) {
	mapVal, ok := value.(basetypes.MapValuable)

	if !ok {
		return types.MapNull(nil), diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return mapVal.ToMapValue(ctx)
}

func coerceObjectValue(ctx context.Context, schemaPath path.Path, value attr.Value) (types.Object, diag.Diagnostics) {
	objectVal, ok := value.(basetypes.ObjectValuable)

	if !ok {
		return types.ObjectNull(nil), diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return objectVal.ToObjectValue(ctx)
}

func coerceSetValue(ctx context.Context, schemaPath path.Path, value attr.Value) (types.Set, diag.Diagnostics) {
	setVal, ok := value.(basetypes.SetValuable)

	if !ok {
		return types.SetNull(nil), diag.Diagnostics{
			attributePlanModificationWalkError(schemaPath, value),
		}
	}

	return setVal.ToSetValue(ctx)
}

func listElemObject(ctx context.Context, schemaPath path.Path, list types.List, index int, description fwschemadata.DataDescription) (types.Object, diag.Diagnostics) {
	if list.IsNull() {
		return listElemObjectFromTerraformValue(ctx, schemaPath, list, description, nil)
	}

	if list.IsUnknown() {
		return listElemObjectFromTerraformValue(ctx, schemaPath, list, description, tftypes.UnknownValue)
	}

	if index >= len(list.Elements()) {
		return listElemObjectFromTerraformValue(ctx, schemaPath, list, description, nil)
	}

	return coerceObjectValue(ctx, schemaPath, list.Elements()[index])
}

func listElemObjectFromTerraformValue(ctx context.Context, schemaPath path.Path, list types.List, description fwschemadata.DataDescription, tfValue any) (types.Object, diag.Diagnostics) {
	elemType := list.ElementType(ctx)
	elemValue, err := elemType.ValueFromTerraform(ctx, tftypes.NewValue(elemType.TerraformType(ctx), tfValue))

	if err != nil {
		return types.ObjectNull(nil), diag.Diagnostics{
			attributePlanModificationValueError(ctx, list, description, err),
		}
	}

	return coerceObjectValue(ctx, schemaPath, elemValue)
}

func mapElemObject(ctx context.Context, schemaPath path.Path, m types.Map, key string, description fwschemadata.DataDescription) (types.Object, diag.Diagnostics) {
	if m.IsNull() {
		return mapElemObjectFromTerraformValue(ctx, schemaPath, m, description, nil)
	}

	if m.IsUnknown() {
		return mapElemObjectFromTerraformValue(ctx, schemaPath, m, description, tftypes.UnknownValue)
	}

	elemValue, ok := m.Elements()[key]

	if !ok {
		return mapElemObjectFromTerraformValue(ctx, schemaPath, m, description, nil)
	}

	return coerceObjectValue(ctx, schemaPath, elemValue)
}

func mapElemObjectFromTerraformValue(ctx context.Context, schemaPath path.Path, m types.Map, description fwschemadata.DataDescription, tfValue any) (types.Object, diag.Diagnostics) {
	elemType := m.ElementType(ctx)
	elemValue, err := elemType.ValueFromTerraform(ctx, tftypes.NewValue(elemType.TerraformType(ctx), tfValue))

	if err != nil {
		return types.ObjectNull(nil), diag.Diagnostics{
			attributePlanModificationValueError(ctx, m, description, err),
		}
	}

	return coerceObjectValue(ctx, schemaPath, elemValue)
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
	return object.Attributes()[attributeName], nil
}

func objectAttributeValueFromTerraformValue(ctx context.Context, object types.Object, attributeName string, description fwschemadata.DataDescription, tfValue any) (attr.Value, diag.Diagnostics) {
	// A panic here indicates a bug somewhere else in the framework or an
	// invalid test case.
	attrType := object.AttributeTypes(ctx)[attributeName]

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

	if index >= len(set.Elements()) {
		return setElemObjectFromTerraformValue(ctx, schemaPath, set, description, nil)
	}

	return coerceObjectValue(ctx, schemaPath, set.Elements()[index])
}

func setElemObjectFromTerraformValue(ctx context.Context, schemaPath path.Path, set types.Set, description fwschemadata.DataDescription, tfValue any) (types.Object, diag.Diagnostics) {
	elemType := set.ElementType(ctx)
	elemValue, err := elemType.ValueFromTerraform(ctx, tftypes.NewValue(elemType.TerraformType(ctx), tfValue))

	if err != nil {
		return types.ObjectNull(nil), diag.Diagnostics{
			attributePlanModificationValueError(ctx, set, description, err),
		}
	}

	return coerceObjectValue(ctx, schemaPath, elemValue)
}
