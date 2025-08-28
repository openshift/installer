// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package simulator

import (
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type CustomFieldsManager struct {
	mo.CustomFieldsManager

	nextKey int32
}

// Iterates through all entities of passed field type;
// Removes found field from their custom field properties.
func entitiesFieldRemove(ctx *Context, field types.CustomFieldDef) {
	entities := ctx.Map.All(field.ManagedObjectType)
	for _, e := range entities {
		entity := e.Entity()
		ctx.WithLock(entity, func() {
			aFields := entity.AvailableField
			for i, aField := range aFields {
				if aField.Key == field.Key {
					entity.AvailableField = append(aFields[:i], aFields[i+1:]...)
					break
				}
			}

			values := e.Entity().Value
			for i, value := range values {
				if value.(*types.CustomFieldStringValue).Key == field.Key {
					entity.Value = append(values[:i], values[i+1:]...)
					break
				}
			}

			cValues := e.Entity().CustomValue
			for i, cValue := range cValues {
				if cValue.(*types.CustomFieldStringValue).Key == field.Key {
					entity.CustomValue = append(cValues[:i], cValues[i+1:]...)
					break
				}
			}
		})
	}
}

// Iterates through all entities of passed field type;
// Renames found field in entity's AvailableField property.
func entitiesFieldRename(ctx *Context, field types.CustomFieldDef) {
	entities := ctx.Map.All(field.ManagedObjectType)
	for _, e := range entities {
		entity := e.Entity()
		ctx.WithLock(entity, func() {
			aFields := entity.AvailableField
			for i, aField := range aFields {
				if aField.Key == field.Key {
					aFields[i].Name = field.Name
					break
				}
			}
		})
	}
}

func (c *CustomFieldsManager) findByNameType(name, moType string) (int, *types.CustomFieldDef) {
	for i, field := range c.Field {
		if (field.ManagedObjectType == "" || field.ManagedObjectType == moType || moType == "") &&
			field.Name == name {
			return i, &c.Field[i]
		}
	}

	return -1, nil
}

func (c *CustomFieldsManager) findByKey(key int32) (int, *types.CustomFieldDef) {
	for i, field := range c.Field {
		if field.Key == key {
			return i, &c.Field[i]
		}
	}

	return -1, nil
}

func (c *CustomFieldsManager) AddCustomFieldDef(ctx *Context, req *types.AddCustomFieldDef) soap.HasFault {
	body := &methods.AddCustomFieldDefBody{}

	_, field := c.findByNameType(req.Name, req.MoType)
	if field != nil {
		body.Fault_ = Fault("", &types.DuplicateName{
			Name:   req.Name,
			Object: c.Reference(),
		})
		return body
	}

	def := types.CustomFieldDef{
		Key:                     c.nextKey,
		Name:                    req.Name,
		ManagedObjectType:       req.MoType,
		Type:                    req.MoType,
		FieldDefPrivileges:      req.FieldDefPolicy,
		FieldInstancePrivileges: req.FieldPolicy,
	}

	entities := ctx.Map.All(req.MoType)
	for _, e := range entities {
		entity := e.Entity()
		ctx.WithLock(entity, func() {
			entity.AvailableField = append(entity.AvailableField, def)
		})
	}

	c.Field = append(c.Field, def)
	c.nextKey++

	body.Res = &types.AddCustomFieldDefResponse{
		Returnval: def,
	}
	return body
}

func (c *CustomFieldsManager) RemoveCustomFieldDef(ctx *Context, req *types.RemoveCustomFieldDef) soap.HasFault {
	body := &methods.RemoveCustomFieldDefBody{}

	i, field := c.findByKey(req.Key)
	if field == nil {
		body.Fault_ = Fault("", &types.NotFound{})
		return body
	}

	entitiesFieldRemove(ctx, *field)

	c.Field = append(c.Field[:i], c.Field[i+1:]...)

	body.Res = &types.RemoveCustomFieldDefResponse{}
	return body
}

func (c *CustomFieldsManager) RenameCustomFieldDef(ctx *Context, req *types.RenameCustomFieldDef) soap.HasFault {
	body := &methods.RenameCustomFieldDefBody{}

	_, field := c.findByKey(req.Key)
	if field == nil {
		body.Fault_ = Fault("", &types.NotFound{})
		return body
	}

	field.Name = req.Name

	entitiesFieldRename(ctx, *field)

	body.Res = &types.RenameCustomFieldDefResponse{}
	return body
}

func (c *CustomFieldsManager) SetField(ctx *Context, req *types.SetField) soap.HasFault {
	body := &methods.SetFieldBody{}

	_, field := c.findByKey(req.Key)
	if field == nil {
		body.Fault_ = Fault("", &types.InvalidArgument{InvalidProperty: "key"})
		return body
	}

	newValue := &types.CustomFieldStringValue{
		CustomFieldValue: types.CustomFieldValue{Key: req.Key},
		Value:            req.Value,
	}

	removeIndex := func(s []types.BaseCustomFieldValue, i int) []types.BaseCustomFieldValue {
		new := make([]types.BaseCustomFieldValue, 0)
		new = append(new, s[:i]...)
		return append(new, s[i+1:]...)
	}

	removeExistingValues := func(s []types.BaseCustomFieldValue) []types.BaseCustomFieldValue {
		for i := 0; i < len(s); {
			if s[i].GetCustomFieldValue().Key == newValue.GetCustomFieldValue().Key {
				s = removeIndex(s, i)
			}
			i++
		}
		return s
	}

	entity := ctx.Map.Get(req.Entity).(mo.Entity).Entity()

	ctx.WithLock(entity, func() {
		// Check if custom value and value are already set. If so, remove them.
		entity.CustomValue = removeExistingValues(entity.CustomValue)
		entity.Value = removeExistingValues(entity.Value)

		// Add the new value
		entity.CustomValue = append(entity.CustomValue, newValue)
		entity.Value = append(entity.Value, newValue)
	})

	body.Res = &types.SetFieldResponse{}
	return body
}
