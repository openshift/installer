package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AzureActiveDirectoryTenant 
type AzureActiveDirectoryTenant struct {
    IdentitySource
    // The name of the Azure Active Directory tenant. Read only.
    displayName *string
    // The ID of the Azure Active Directory tenant. Read only.
    tenantId *string
}
// NewAzureActiveDirectoryTenant instantiates a new AzureActiveDirectoryTenant and sets the default values.
func NewAzureActiveDirectoryTenant()(*AzureActiveDirectoryTenant) {
    m := &AzureActiveDirectoryTenant{
        IdentitySource: *NewIdentitySource(),
    }
    odataTypeValue := "#microsoft.graph.azureActiveDirectoryTenant";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateAzureActiveDirectoryTenantFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAzureActiveDirectoryTenantFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAzureActiveDirectoryTenant(), nil
}
// GetDisplayName gets the displayName property value. The name of the Azure Active Directory tenant. Read only.
func (m *AzureActiveDirectoryTenant) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AzureActiveDirectoryTenant) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IdentitySource.GetFieldDeserializers()
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["tenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTenantId)
    return res
}
// GetTenantId gets the tenantId property value. The ID of the Azure Active Directory tenant. Read only.
func (m *AzureActiveDirectoryTenant) GetTenantId()(*string) {
    return m.tenantId
}
// Serialize serializes information the current object
func (m *AzureActiveDirectoryTenant) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IdentitySource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The name of the Azure Active Directory tenant. Read only.
func (m *AzureActiveDirectoryTenant) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetTenantId sets the tenantId property value. The ID of the Azure Active Directory tenant. Read only.
func (m *AzureActiveDirectoryTenant) SetTenantId(value *string)() {
    m.tenantId = value
}
