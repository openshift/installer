package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosMobileAppConfiguration 
type IosMobileAppConfiguration struct {
    ManagedDeviceMobileAppConfiguration
    // mdm app configuration Base64 binary.
    encodedSettingXml []byte
    // app configuration setting items.
    settings []AppConfigurationSettingItemable
}
// NewIosMobileAppConfiguration instantiates a new IosMobileAppConfiguration and sets the default values.
func NewIosMobileAppConfiguration()(*IosMobileAppConfiguration) {
    m := &IosMobileAppConfiguration{
        ManagedDeviceMobileAppConfiguration: *NewManagedDeviceMobileAppConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.iosMobileAppConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateIosMobileAppConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosMobileAppConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosMobileAppConfiguration(), nil
}
// GetEncodedSettingXml gets the encodedSettingXml property value. mdm app configuration Base64 binary.
func (m *IosMobileAppConfiguration) GetEncodedSettingXml()([]byte) {
    return m.encodedSettingXml
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosMobileAppConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedDeviceMobileAppConfiguration.GetFieldDeserializers()
    res["encodedSettingXml"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteArrayValue(m.SetEncodedSettingXml)
    res["settings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAppConfigurationSettingItemFromDiscriminatorValue , m.SetSettings)
    return res
}
// GetSettings gets the settings property value. app configuration setting items.
func (m *IosMobileAppConfiguration) GetSettings()([]AppConfigurationSettingItemable) {
    return m.settings
}
// Serialize serializes information the current object
func (m *IosMobileAppConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedDeviceMobileAppConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteByteArrayValue("encodedSettingXml", m.GetEncodedSettingXml())
        if err != nil {
            return err
        }
    }
    if m.GetSettings() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSettings())
        err = writer.WriteCollectionOfObjectValues("settings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEncodedSettingXml sets the encodedSettingXml property value. mdm app configuration Base64 binary.
func (m *IosMobileAppConfiguration) SetEncodedSettingXml(value []byte)() {
    m.encodedSettingXml = value
}
// SetSettings sets the settings property value. app configuration setting items.
func (m *IosMobileAppConfiguration) SetSettings(value []AppConfigurationSettingItemable)() {
    m.settings = value
}
