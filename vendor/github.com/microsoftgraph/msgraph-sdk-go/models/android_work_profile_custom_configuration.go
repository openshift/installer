package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AndroidWorkProfileCustomConfiguration 
type AndroidWorkProfileCustomConfiguration struct {
    DeviceConfiguration
    // OMA settings. This collection can contain a maximum of 500 elements.
    omaSettings []OmaSettingable
}
// NewAndroidWorkProfileCustomConfiguration instantiates a new AndroidWorkProfileCustomConfiguration and sets the default values.
func NewAndroidWorkProfileCustomConfiguration()(*AndroidWorkProfileCustomConfiguration) {
    m := &AndroidWorkProfileCustomConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.androidWorkProfileCustomConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateAndroidWorkProfileCustomConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAndroidWorkProfileCustomConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAndroidWorkProfileCustomConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AndroidWorkProfileCustomConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["omaSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOmaSettingFromDiscriminatorValue , m.SetOmaSettings)
    return res
}
// GetOmaSettings gets the omaSettings property value. OMA settings. This collection can contain a maximum of 500 elements.
func (m *AndroidWorkProfileCustomConfiguration) GetOmaSettings()([]OmaSettingable) {
    return m.omaSettings
}
// Serialize serializes information the current object
func (m *AndroidWorkProfileCustomConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetOmaSettings() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOmaSettings())
        err = writer.WriteCollectionOfObjectValues("omaSettings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOmaSettings sets the omaSettings property value. OMA settings. This collection can contain a maximum of 500 elements.
func (m *AndroidWorkProfileCustomConfiguration) SetOmaSettings(value []OmaSettingable)() {
    m.omaSettings = value
}
