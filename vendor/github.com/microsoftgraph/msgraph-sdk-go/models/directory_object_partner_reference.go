package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DirectoryObjectPartnerReference 
type DirectoryObjectPartnerReference struct {
    DirectoryObject
    // Description of the object returned. Read-only.
    description *string
    // Name of directory object being returned, like group or application. Read-only.
    displayName *string
    // The tenant identifier for the partner tenant. Read-only.
    externalPartnerTenantId *string
    // The type of the referenced object in the partner tenant. Read-only.
    objectType *string
}
// NewDirectoryObjectPartnerReference instantiates a new DirectoryObjectPartnerReference and sets the default values.
func NewDirectoryObjectPartnerReference()(*DirectoryObjectPartnerReference) {
    m := &DirectoryObjectPartnerReference{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.directoryObjectPartnerReference";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateDirectoryObjectPartnerReferenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDirectoryObjectPartnerReferenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectoryObjectPartnerReference(), nil
}
// GetDescription gets the description property value. Description of the object returned. Read-only.
func (m *DirectoryObjectPartnerReference) GetDescription()(*string) {
    return m.description
}
// GetDisplayName gets the displayName property value. Name of directory object being returned, like group or application. Read-only.
func (m *DirectoryObjectPartnerReference) GetDisplayName()(*string) {
    return m.displayName
}
// GetExternalPartnerTenantId gets the externalPartnerTenantId property value. The tenant identifier for the partner tenant. Read-only.
func (m *DirectoryObjectPartnerReference) GetExternalPartnerTenantId()(*string) {
    return m.externalPartnerTenantId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DirectoryObjectPartnerReference) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["externalPartnerTenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetExternalPartnerTenantId)
    res["objectType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetObjectType)
    return res
}
// GetObjectType gets the objectType property value. The type of the referenced object in the partner tenant. Read-only.
func (m *DirectoryObjectPartnerReference) GetObjectType()(*string) {
    return m.objectType
}
// Serialize serializes information the current object
func (m *DirectoryObjectPartnerReference) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
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
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalPartnerTenantId", m.GetExternalPartnerTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("objectType", m.GetObjectType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. Description of the object returned. Read-only.
func (m *DirectoryObjectPartnerReference) SetDescription(value *string)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. Name of directory object being returned, like group or application. Read-only.
func (m *DirectoryObjectPartnerReference) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetExternalPartnerTenantId sets the externalPartnerTenantId property value. The tenant identifier for the partner tenant. Read-only.
func (m *DirectoryObjectPartnerReference) SetExternalPartnerTenantId(value *string)() {
    m.externalPartnerTenantId = value
}
// SetObjectType sets the objectType property value. The type of the referenced object in the partner tenant. Read-only.
func (m *DirectoryObjectPartnerReference) SetObjectType(value *string)() {
    m.objectType = value
}
