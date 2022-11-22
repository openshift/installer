package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsDefenderScanActionResult 
type WindowsDefenderScanActionResult struct {
    DeviceActionResult
    // Scan type either full scan or quick scan
    scanType *string
}
// NewWindowsDefenderScanActionResult instantiates a new WindowsDefenderScanActionResult and sets the default values.
func NewWindowsDefenderScanActionResult()(*WindowsDefenderScanActionResult) {
    m := &WindowsDefenderScanActionResult{
        DeviceActionResult: *NewDeviceActionResult(),
    }
    return m
}
// CreateWindowsDefenderScanActionResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsDefenderScanActionResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsDefenderScanActionResult(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsDefenderScanActionResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceActionResult.GetFieldDeserializers()
    res["scanType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetScanType)
    return res
}
// GetScanType gets the scanType property value. Scan type either full scan or quick scan
func (m *WindowsDefenderScanActionResult) GetScanType()(*string) {
    return m.scanType
}
// Serialize serializes information the current object
func (m *WindowsDefenderScanActionResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceActionResult.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("scanType", m.GetScanType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetScanType sets the scanType property value. Scan type either full scan or quick scan
func (m *WindowsDefenderScanActionResult) SetScanType(value *string)() {
    m.scanType = value
}
