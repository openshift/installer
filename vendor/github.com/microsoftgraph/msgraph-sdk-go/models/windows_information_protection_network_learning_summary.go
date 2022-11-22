package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtectionNetworkLearningSummary windows Information Protection Network learning Summary entity.
type WindowsInformationProtectionNetworkLearningSummary struct {
    Entity
    // Device Count
    deviceCount *int32
    // Website url
    url *string
}
// NewWindowsInformationProtectionNetworkLearningSummary instantiates a new windowsInformationProtectionNetworkLearningSummary and sets the default values.
func NewWindowsInformationProtectionNetworkLearningSummary()(*WindowsInformationProtectionNetworkLearningSummary) {
    m := &WindowsInformationProtectionNetworkLearningSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWindowsInformationProtectionNetworkLearningSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsInformationProtectionNetworkLearningSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsInformationProtectionNetworkLearningSummary(), nil
}
// GetDeviceCount gets the deviceCount property value. Device Count
func (m *WindowsInformationProtectionNetworkLearningSummary) GetDeviceCount()(*int32) {
    return m.deviceCount
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsInformationProtectionNetworkLearningSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetDeviceCount)
    res["url"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUrl)
    return res
}
// GetUrl gets the url property value. Website url
func (m *WindowsInformationProtectionNetworkLearningSummary) GetUrl()(*string) {
    return m.url
}
// Serialize serializes information the current object
func (m *WindowsInformationProtectionNetworkLearningSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("deviceCount", m.GetDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("url", m.GetUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceCount sets the deviceCount property value. Device Count
func (m *WindowsInformationProtectionNetworkLearningSummary) SetDeviceCount(value *int32)() {
    m.deviceCount = value
}
// SetUrl sets the url property value. Website url
func (m *WindowsInformationProtectionNetworkLearningSummary) SetUrl(value *string)() {
    m.url = value
}
