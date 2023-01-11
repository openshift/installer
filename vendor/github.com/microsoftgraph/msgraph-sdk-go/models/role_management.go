package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RoleManagement 
type RoleManagement struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The directory property
    directory RbacApplicationable
    // Container for roles and assignments for entitlement management resources.
    entitlementManagement RbacApplicationable
    // The OdataType property
    odataType *string
}
// NewRoleManagement instantiates a new RoleManagement and sets the default values.
func NewRoleManagement()(*RoleManagement) {
    m := &RoleManagement{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateRoleManagementFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRoleManagementFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRoleManagement(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RoleManagement) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDirectory gets the directory property value. The directory property
func (m *RoleManagement) GetDirectory()(RbacApplicationable) {
    return m.directory
}
// GetEntitlementManagement gets the entitlementManagement property value. Container for roles and assignments for entitlement management resources.
func (m *RoleManagement) GetEntitlementManagement()(RbacApplicationable) {
    return m.entitlementManagement
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RoleManagement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["directory"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRbacApplicationFromDiscriminatorValue , m.SetDirectory)
    res["entitlementManagement"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRbacApplicationFromDiscriminatorValue , m.SetEntitlementManagement)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *RoleManagement) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *RoleManagement) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("directory", m.GetDirectory())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("entitlementManagement", m.GetEntitlementManagement())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RoleManagement) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDirectory sets the directory property value. The directory property
func (m *RoleManagement) SetDirectory(value RbacApplicationable)() {
    m.directory = value
}
// SetEntitlementManagement sets the entitlementManagement property value. Container for roles and assignments for entitlement management resources.
func (m *RoleManagement) SetEntitlementManagement(value RbacApplicationable)() {
    m.entitlementManagement = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RoleManagement) SetOdataType(value *string)() {
    m.odataType = value
}
