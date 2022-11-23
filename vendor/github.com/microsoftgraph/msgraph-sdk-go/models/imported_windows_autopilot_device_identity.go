package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ImportedWindowsAutopilotDeviceIdentity imported windows autopilot devices.
type ImportedWindowsAutopilotDeviceIdentity struct {
    Entity
    // UPN of the user the device will be assigned
    assignedUserPrincipalName *string
    // Group Tag of the Windows autopilot device.
    groupTag *string
    // Hardware Blob of the Windows autopilot device.
    hardwareIdentifier []byte
    // The Import Id of the Windows autopilot device.
    importId *string
    // Product Key of the Windows autopilot device.
    productKey *string
    // Serial number of the Windows autopilot device.
    serialNumber *string
    // Current state of the imported device.
    state ImportedWindowsAutopilotDeviceIdentityStateable
}
// NewImportedWindowsAutopilotDeviceIdentity instantiates a new importedWindowsAutopilotDeviceIdentity and sets the default values.
func NewImportedWindowsAutopilotDeviceIdentity()(*ImportedWindowsAutopilotDeviceIdentity) {
    m := &ImportedWindowsAutopilotDeviceIdentity{
        Entity: *NewEntity(),
    }
    return m
}
// CreateImportedWindowsAutopilotDeviceIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateImportedWindowsAutopilotDeviceIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewImportedWindowsAutopilotDeviceIdentity(), nil
}
// GetAssignedUserPrincipalName gets the assignedUserPrincipalName property value. UPN of the user the device will be assigned
func (m *ImportedWindowsAutopilotDeviceIdentity) GetAssignedUserPrincipalName()(*string) {
    return m.assignedUserPrincipalName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ImportedWindowsAutopilotDeviceIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignedUserPrincipalName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAssignedUserPrincipalName)
    res["groupTag"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetGroupTag)
    res["hardwareIdentifier"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteArrayValue(m.SetHardwareIdentifier)
    res["importId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetImportId)
    res["productKey"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetProductKey)
    res["serialNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSerialNumber)
    res["state"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateImportedWindowsAutopilotDeviceIdentityStateFromDiscriminatorValue , m.SetState)
    return res
}
// GetGroupTag gets the groupTag property value. Group Tag of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) GetGroupTag()(*string) {
    return m.groupTag
}
// GetHardwareIdentifier gets the hardwareIdentifier property value. Hardware Blob of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) GetHardwareIdentifier()([]byte) {
    return m.hardwareIdentifier
}
// GetImportId gets the importId property value. The Import Id of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) GetImportId()(*string) {
    return m.importId
}
// GetProductKey gets the productKey property value. Product Key of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) GetProductKey()(*string) {
    return m.productKey
}
// GetSerialNumber gets the serialNumber property value. Serial number of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) GetSerialNumber()(*string) {
    return m.serialNumber
}
// GetState gets the state property value. Current state of the imported device.
func (m *ImportedWindowsAutopilotDeviceIdentity) GetState()(ImportedWindowsAutopilotDeviceIdentityStateable) {
    return m.state
}
// Serialize serializes information the current object
func (m *ImportedWindowsAutopilotDeviceIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("assignedUserPrincipalName", m.GetAssignedUserPrincipalName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("groupTag", m.GetGroupTag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("hardwareIdentifier", m.GetHardwareIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("importId", m.GetImportId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productKey", m.GetProductKey())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("serialNumber", m.GetSerialNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignedUserPrincipalName sets the assignedUserPrincipalName property value. UPN of the user the device will be assigned
func (m *ImportedWindowsAutopilotDeviceIdentity) SetAssignedUserPrincipalName(value *string)() {
    m.assignedUserPrincipalName = value
}
// SetGroupTag sets the groupTag property value. Group Tag of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetGroupTag(value *string)() {
    m.groupTag = value
}
// SetHardwareIdentifier sets the hardwareIdentifier property value. Hardware Blob of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetHardwareIdentifier(value []byte)() {
    m.hardwareIdentifier = value
}
// SetImportId sets the importId property value. The Import Id of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetImportId(value *string)() {
    m.importId = value
}
// SetProductKey sets the productKey property value. Product Key of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetProductKey(value *string)() {
    m.productKey = value
}
// SetSerialNumber sets the serialNumber property value. Serial number of the Windows autopilot device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetSerialNumber(value *string)() {
    m.serialNumber = value
}
// SetState sets the state property value. Current state of the imported device.
func (m *ImportedWindowsAutopilotDeviceIdentity) SetState(value ImportedWindowsAutopilotDeviceIdentityStateable)() {
    m.state = value
}
