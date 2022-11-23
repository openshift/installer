package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SchemaExtension 
type SchemaExtension struct {
    Entity
    // Description for the schema extension. Supports $filter (eq).
    description *string
    // The appId of the application that is the owner of the schema extension. This property can be supplied on creation, to set the owner.  If not supplied, then the calling application's appId will be set as the owner. In either case, the signed-in user must be the owner of the application. So, for example, if creating a new schema extension definition using Graph Explorer, you must supply the owner property. Once set, this property is read-only and cannot be changed. Supports $filter (eq).
    owner *string
    // The collection of property names and types that make up the schema extension definition.
    properties []ExtensionSchemaPropertyable
    // The lifecycle state of the schema extension. Possible states are InDevelopment, Available, and Deprecated. Automatically set to InDevelopment on creation. For more information about the possible state transitions and behaviors, see Schema extensions lifecycle. Supports $filter (eq).
    status *string
    // Set of Microsoft Graph types (that can support extensions) that the schema extension can be applied to. Select from administrativeUnit, contact, device, event, group, message, organization, post, todoTask, todoTaskList, or user.
    targetTypes []string
}
// NewSchemaExtension instantiates a new SchemaExtension and sets the default values.
func NewSchemaExtension()(*SchemaExtension) {
    m := &SchemaExtension{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSchemaExtensionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSchemaExtensionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSchemaExtension(), nil
}
// GetDescription gets the description property value. Description for the schema extension. Supports $filter (eq).
func (m *SchemaExtension) GetDescription()(*string) {
    return m.description
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SchemaExtension) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["owner"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOwner)
    res["properties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateExtensionSchemaPropertyFromDiscriminatorValue , m.SetProperties)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetStatus)
    res["targetTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetTargetTypes)
    return res
}
// GetOwner gets the owner property value. The appId of the application that is the owner of the schema extension. This property can be supplied on creation, to set the owner.  If not supplied, then the calling application's appId will be set as the owner. In either case, the signed-in user must be the owner of the application. So, for example, if creating a new schema extension definition using Graph Explorer, you must supply the owner property. Once set, this property is read-only and cannot be changed. Supports $filter (eq).
func (m *SchemaExtension) GetOwner()(*string) {
    return m.owner
}
// GetProperties gets the properties property value. The collection of property names and types that make up the schema extension definition.
func (m *SchemaExtension) GetProperties()([]ExtensionSchemaPropertyable) {
    return m.properties
}
// GetStatus gets the status property value. The lifecycle state of the schema extension. Possible states are InDevelopment, Available, and Deprecated. Automatically set to InDevelopment on creation. For more information about the possible state transitions and behaviors, see Schema extensions lifecycle. Supports $filter (eq).
func (m *SchemaExtension) GetStatus()(*string) {
    return m.status
}
// GetTargetTypes gets the targetTypes property value. Set of Microsoft Graph types (that can support extensions) that the schema extension can be applied to. Select from administrativeUnit, contact, device, event, group, message, organization, post, todoTask, todoTaskList, or user.
func (m *SchemaExtension) GetTargetTypes()([]string) {
    return m.targetTypes
}
// Serialize serializes information the current object
func (m *SchemaExtension) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("owner", m.GetOwner())
        if err != nil {
            return err
        }
    }
    if m.GetProperties() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetProperties())
        err = writer.WriteCollectionOfObjectValues("properties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    if m.GetTargetTypes() != nil {
        err = writer.WriteCollectionOfStringValues("targetTypes", m.GetTargetTypes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. Description for the schema extension. Supports $filter (eq).
func (m *SchemaExtension) SetDescription(value *string)() {
    m.description = value
}
// SetOwner sets the owner property value. The appId of the application that is the owner of the schema extension. This property can be supplied on creation, to set the owner.  If not supplied, then the calling application's appId will be set as the owner. In either case, the signed-in user must be the owner of the application. So, for example, if creating a new schema extension definition using Graph Explorer, you must supply the owner property. Once set, this property is read-only and cannot be changed. Supports $filter (eq).
func (m *SchemaExtension) SetOwner(value *string)() {
    m.owner = value
}
// SetProperties sets the properties property value. The collection of property names and types that make up the schema extension definition.
func (m *SchemaExtension) SetProperties(value []ExtensionSchemaPropertyable)() {
    m.properties = value
}
// SetStatus sets the status property value. The lifecycle state of the schema extension. Possible states are InDevelopment, Available, and Deprecated. Automatically set to InDevelopment on creation. For more information about the possible state transitions and behaviors, see Schema extensions lifecycle. Supports $filter (eq).
func (m *SchemaExtension) SetStatus(value *string)() {
    m.status = value
}
// SetTargetTypes sets the targetTypes property value. Set of Microsoft Graph types (that can support extensions) that the schema extension can be applied to. Select from administrativeUnit, contact, device, event, group, message, organization, post, todoTask, todoTaskList, or user.
func (m *SchemaExtension) SetTargetTypes(value []string)() {
    m.targetTypes = value
}
