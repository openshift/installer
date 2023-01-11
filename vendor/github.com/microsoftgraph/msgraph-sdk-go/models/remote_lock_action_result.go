package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RemoteLockActionResult 
type RemoteLockActionResult struct {
    DeviceActionResult
    // Pin to unlock the client
    unlockPin *string
}
// NewRemoteLockActionResult instantiates a new RemoteLockActionResult and sets the default values.
func NewRemoteLockActionResult()(*RemoteLockActionResult) {
    m := &RemoteLockActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateRemoteLockActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRemoteLockActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRemoteLockActionResult(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RemoteLockActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["unlockPin"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUnlockPin)
    return res
}
// GetUnlockPin gets the unlockPin property value. Pin to unlock the client
func (m *RemoteLockActionResult) GetUnlockPin()(*string) {
    return m.unlockPin
}
// Serialize serializes information the current object
func (m *RemoteLockActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("unlockPin", m.GetUnlockPin())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUnlockPin sets the unlockPin property value. Pin to unlock the client
func (m *RemoteLockActionResult) SetUnlockPin(value *string)() {
    m.unlockPin = value
}
