package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AuditResource a class containing the properties for Audit Resource.
type AuditResource struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Audit resource's type.
    auditResourceType *string
    // Display name.
    displayName *string
    // List of modified properties.
    modifiedProperties []AuditPropertyable
    // The OdataType property
    odataType *string
    // Audit resource's Id.
    resourceId *string
}
// NewAuditResource instantiates a new auditResource and sets the default values.
func NewAuditResource()(*AuditResource) {
    m := &AuditResource{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAuditResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAuditResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuditResource(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AuditResource) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAuditResourceType gets the auditResourceType property value. Audit resource's type.
func (m *AuditResource) GetAuditResourceType()(*string) {
    return m.auditResourceType
}
// GetDisplayName gets the displayName property value. Display name.
func (m *AuditResource) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AuditResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["auditResourceType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAuditResourceType)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["modifiedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAuditPropertyFromDiscriminatorValue , m.SetModifiedProperties)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["resourceId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetResourceId)
    return res
}
// GetModifiedProperties gets the modifiedProperties property value. List of modified properties.
func (m *AuditResource) GetModifiedProperties()([]AuditPropertyable) {
    return m.modifiedProperties
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AuditResource) GetOdataType()(*string) {
    return m.odataType
}
// GetResourceId gets the resourceId property value. Audit resource's Id.
func (m *AuditResource) GetResourceId()(*string) {
    return m.resourceId
}
// Serialize serializes information the current object
func (m *AuditResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("auditResourceType", m.GetAuditResourceType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetModifiedProperties() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetModifiedProperties())
        err := writer.WriteCollectionOfObjectValues("modifiedProperties", cast)
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
        err := writer.WriteStringValue("resourceId", m.GetResourceId())
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
func (m *AuditResource) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAuditResourceType sets the auditResourceType property value. Audit resource's type.
func (m *AuditResource) SetAuditResourceType(value *string)() {
    m.auditResourceType = value
}
// SetDisplayName sets the displayName property value. Display name.
func (m *AuditResource) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetModifiedProperties sets the modifiedProperties property value. List of modified properties.
func (m *AuditResource) SetModifiedProperties(value []AuditPropertyable)() {
    m.modifiedProperties = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AuditResource) SetOdataType(value *string)() {
    m.odataType = value
}
// SetResourceId sets the resourceId property value. Audit resource's Id.
func (m *AuditResource) SetResourceId(value *string)() {
    m.resourceId = value
}
