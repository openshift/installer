package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtectionProxiedDomainCollection windows Information Protection Proxied Domain Collection
type WindowsInformationProtectionProxiedDomainCollection struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Display name
    displayName *string
    // The OdataType property
    odataType *string
    // Collection of proxied domains
    proxiedDomains []ProxiedDomainable
}
// NewWindowsInformationProtectionProxiedDomainCollection instantiates a new windowsInformationProtectionProxiedDomainCollection and sets the default values.
func NewWindowsInformationProtectionProxiedDomainCollection()(*WindowsInformationProtectionProxiedDomainCollection) {
    m := &WindowsInformationProtectionProxiedDomainCollection{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWindowsInformationProtectionProxiedDomainCollectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsInformationProtectionProxiedDomainCollectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsInformationProtectionProxiedDomainCollection(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *WindowsInformationProtectionProxiedDomainCollection) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDisplayName gets the displayName property value. Display name
func (m *WindowsInformationProtectionProxiedDomainCollection) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsInformationProtectionProxiedDomainCollection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["proxiedDomains"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateProxiedDomainFromDiscriminatorValue , m.SetProxiedDomains)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *WindowsInformationProtectionProxiedDomainCollection) GetOdataType()(*string) {
    return m.odataType
}
// GetProxiedDomains gets the proxiedDomains property value. Collection of proxied domains
func (m *WindowsInformationProtectionProxiedDomainCollection) GetProxiedDomains()([]ProxiedDomainable) {
    return m.proxiedDomains
}
// Serialize serializes information the current object
func (m *WindowsInformationProtectionProxiedDomainCollection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    if m.GetProxiedDomains() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetProxiedDomains())
        err := writer.WriteCollectionOfObjectValues("proxiedDomains", cast)
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
func (m *WindowsInformationProtectionProxiedDomainCollection) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDisplayName sets the displayName property value. Display name
func (m *WindowsInformationProtectionProxiedDomainCollection) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WindowsInformationProtectionProxiedDomainCollection) SetOdataType(value *string)() {
    m.odataType = value
}
// SetProxiedDomains sets the proxiedDomains property value. Collection of proxied domains
func (m *WindowsInformationProtectionProxiedDomainCollection) SetProxiedDomains(value []ProxiedDomainable)() {
    m.proxiedDomains = value
}
