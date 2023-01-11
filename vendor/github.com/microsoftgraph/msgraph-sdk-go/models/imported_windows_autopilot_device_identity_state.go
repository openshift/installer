package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ImportedWindowsAutopilotDeviceIdentityState 
type ImportedWindowsAutopilotDeviceIdentityState struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Device error code reported by Device Directory Service(DDS).
    deviceErrorCode *int32
    // Device error name reported by Device Directory Service(DDS).
    deviceErrorName *string
    // The deviceImportStatus property
    deviceImportStatus *ImportedWindowsAutopilotDeviceIdentityImportStatus
    // Device Registration ID for successfully added device reported by Device Directory Service(DDS).
    deviceRegistrationId *string
    // The OdataType property
    odataType *string
}
// NewImportedWindowsAutopilotDeviceIdentityState instantiates a new importedWindowsAutopilotDeviceIdentityState and sets the default values.
func NewImportedWindowsAutopilotDeviceIdentityState()(*ImportedWindowsAutopilotDeviceIdentityState) {
    m := &ImportedWindowsAutopilotDeviceIdentityState{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateImportedWindowsAutopilotDeviceIdentityStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateImportedWindowsAutopilotDeviceIdentityStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewImportedWindowsAutopilotDeviceIdentityState(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ImportedWindowsAutopilotDeviceIdentityState) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDeviceErrorCode gets the deviceErrorCode property value. Device error code reported by Device Directory Service(DDS).
func (m *ImportedWindowsAutopilotDeviceIdentityState) GetDeviceErrorCode()(*int32) {
    return m.deviceErrorCode
}
// GetDeviceErrorName gets the deviceErrorName property value. Device error name reported by Device Directory Service(DDS).
func (m *ImportedWindowsAutopilotDeviceIdentityState) GetDeviceErrorName()(*string) {
    return m.deviceErrorName
}
// GetDeviceImportStatus gets the deviceImportStatus property value. The deviceImportStatus property
func (m *ImportedWindowsAutopilotDeviceIdentityState) GetDeviceImportStatus()(*ImportedWindowsAutopilotDeviceIdentityImportStatus) {
    return m.deviceImportStatus
}
// GetDeviceRegistrationId gets the deviceRegistrationId property value. Device Registration ID for successfully added device reported by Device Directory Service(DDS).
func (m *ImportedWindowsAutopilotDeviceIdentityState) GetDeviceRegistrationId()(*string) {
    return m.deviceRegistrationId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ImportedWindowsAutopilotDeviceIdentityState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["deviceErrorCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetDeviceErrorCode)
    res["deviceErrorName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceErrorName)
    res["deviceImportStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseImportedWindowsAutopilotDeviceIdentityImportStatus , m.SetDeviceImportStatus)
    res["deviceRegistrationId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceRegistrationId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ImportedWindowsAutopilotDeviceIdentityState) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ImportedWindowsAutopilotDeviceIdentityState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("deviceErrorCode", m.GetDeviceErrorCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deviceErrorName", m.GetDeviceErrorName())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceImportStatus() != nil {
        cast := (*m.GetDeviceImportStatus()).String()
        err := writer.WriteStringValue("deviceImportStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deviceRegistrationId", m.GetDeviceRegistrationId())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ImportedWindowsAutopilotDeviceIdentityState) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDeviceErrorCode sets the deviceErrorCode property value. Device error code reported by Device Directory Service(DDS).
func (m *ImportedWindowsAutopilotDeviceIdentityState) SetDeviceErrorCode(value *int32)() {
    m.deviceErrorCode = value
}
// SetDeviceErrorName sets the deviceErrorName property value. Device error name reported by Device Directory Service(DDS).
func (m *ImportedWindowsAutopilotDeviceIdentityState) SetDeviceErrorName(value *string)() {
    m.deviceErrorName = value
}
// SetDeviceImportStatus sets the deviceImportStatus property value. The deviceImportStatus property
func (m *ImportedWindowsAutopilotDeviceIdentityState) SetDeviceImportStatus(value *ImportedWindowsAutopilotDeviceIdentityImportStatus)() {
    m.deviceImportStatus = value
}
// SetDeviceRegistrationId sets the deviceRegistrationId property value. Device Registration ID for successfully added device reported by Device Directory Service(DDS).
func (m *ImportedWindowsAutopilotDeviceIdentityState) SetDeviceRegistrationId(value *string)() {
    m.deviceRegistrationId = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ImportedWindowsAutopilotDeviceIdentityState) SetOdataType(value *string)() {
    m.odataType = value
}
