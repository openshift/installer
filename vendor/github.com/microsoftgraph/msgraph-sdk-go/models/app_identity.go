package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppIdentity 
type AppIdentity struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Refers to the Unique GUID representing Application Id in the Azure Active Directory.
    appId *string
    // Refers to the Application Name displayed in the Azure Portal.
    displayName *string
    // The OdataType property
    odataType *string
    // Refers to the Unique GUID indicating Service Principal Id in Azure Active Directory for the corresponding App.
    servicePrincipalId *string
    // Refers to the Service Principal Name is the Application name in the tenant.
    servicePrincipalName *string
}
// NewAppIdentity instantiates a new appIdentity and sets the default values.
func NewAppIdentity()(*AppIdentity) {
    m := &AppIdentity{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAppIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAppIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppIdentity(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AppIdentity) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAppId gets the appId property value. Refers to the Unique GUID representing Application Id in the Azure Active Directory.
func (m *AppIdentity) GetAppId()(*string) {
    return m.appId
}
// GetDisplayName gets the displayName property value. Refers to the Application Name displayed in the Azure Portal.
func (m *AppIdentity) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AppIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["appId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppId)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["servicePrincipalId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetServicePrincipalId)
    res["servicePrincipalName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetServicePrincipalName)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AppIdentity) GetOdataType()(*string) {
    return m.odataType
}
// GetServicePrincipalId gets the servicePrincipalId property value. Refers to the Unique GUID indicating Service Principal Id in Azure Active Directory for the corresponding App.
func (m *AppIdentity) GetServicePrincipalId()(*string) {
    return m.servicePrincipalId
}
// GetServicePrincipalName gets the servicePrincipalName property value. Refers to the Service Principal Name is the Application name in the tenant.
func (m *AppIdentity) GetServicePrincipalName()(*string) {
    return m.servicePrincipalName
}
// Serialize serializes information the current object
func (m *AppIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("appId", m.GetAppId())
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
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("servicePrincipalId", m.GetServicePrincipalId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("servicePrincipalName", m.GetServicePrincipalName())
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
func (m *AppIdentity) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAppId sets the appId property value. Refers to the Unique GUID representing Application Id in the Azure Active Directory.
func (m *AppIdentity) SetAppId(value *string)() {
    m.appId = value
}
// SetDisplayName sets the displayName property value. Refers to the Application Name displayed in the Azure Portal.
func (m *AppIdentity) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AppIdentity) SetOdataType(value *string)() {
    m.odataType = value
}
// SetServicePrincipalId sets the servicePrincipalId property value. Refers to the Unique GUID indicating Service Principal Id in Azure Active Directory for the corresponding App.
func (m *AppIdentity) SetServicePrincipalId(value *string)() {
    m.servicePrincipalId = value
}
// SetServicePrincipalName sets the servicePrincipalName property value. Refers to the Service Principal Name is the Application name in the tenant.
func (m *AppIdentity) SetServicePrincipalName(value *string)() {
    m.servicePrincipalName = value
}
