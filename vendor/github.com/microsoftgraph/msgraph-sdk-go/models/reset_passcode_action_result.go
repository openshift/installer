package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ResetPasscodeActionResult 
type ResetPasscodeActionResult struct {
    DeviceActionResult
    // Newly generated passcode for the device
    passcode *string
}
// NewResetPasscodeActionResult instantiates a new ResetPasscodeActionResult and sets the default values.
func NewResetPasscodeActionResult()(*ResetPasscodeActionResult) {
    m := &ResetPasscodeActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateResetPasscodeActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateResetPasscodeActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewResetPasscodeActionResult(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ResetPasscodeActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["passcode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPasscode)
    return res
}
// GetPasscode gets the passcode property value. Newly generated passcode for the device
func (m *ResetPasscodeActionResult) GetPasscode()(*string) {
    return m.passcode
}
// Serialize serializes information the current object
func (m *ResetPasscodeActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("passcode", m.GetPasscode())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPasscode sets the passcode property value. Newly generated passcode for the device
func (m *ResetPasscodeActionResult) SetPasscode(value *string)() {
    m.passcode = value
}
