package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// LocateDeviceActionResult 
type LocateDeviceActionResult struct {
    DeviceActionResult
    // device location
    deviceLocation DeviceGeoLocationable
}
// NewLocateDeviceActionResult instantiates a new LocateDeviceActionResult and sets the default values.
func NewLocateDeviceActionResult()(*LocateDeviceActionResult) {
    m := &LocateDeviceActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateLocateDeviceActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateLocateDeviceActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLocateDeviceActionResult(), nil
}
// GetDeviceLocation gets the deviceLocation property value. device location
func (m *LocateDeviceActionResult) GetDeviceLocation()(DeviceGeoLocationable) {
    return m.deviceLocation
}
// GetFieldDeserializers the deserialization information for the current model
func (m *LocateDeviceActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["deviceLocation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceGeoLocationFromDiscriminatorValue , m.SetDeviceLocation)
    return res
}
// Serialize serializes information the current object
func (m *LocateDeviceActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("deviceLocation", m.GetDeviceLocation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceLocation sets the deviceLocation property value. device location
func (m *LocateDeviceActionResult) SetDeviceLocation(value DeviceGeoLocationable)() {
    m.deviceLocation = value
}
