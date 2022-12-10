package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ConditionalAccessDevices 
type ConditionalAccessDevices struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Filter that defines the dynamic-device-syntax rule to include/exclude devices. A filter can use device properties (such as extension attributes) to include/exclude them.
    deviceFilter ConditionalAccessFilterable
    // The OdataType property
    odataType *string
}
// NewConditionalAccessDevices instantiates a new conditionalAccessDevices and sets the default values.
func NewConditionalAccessDevices()(*ConditionalAccessDevices) {
    m := &ConditionalAccessDevices{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateConditionalAccessDevicesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateConditionalAccessDevicesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessDevices(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ConditionalAccessDevices) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDeviceFilter gets the deviceFilter property value. Filter that defines the dynamic-device-syntax rule to include/exclude devices. A filter can use device properties (such as extension attributes) to include/exclude them.
func (m *ConditionalAccessDevices) GetDeviceFilter()(ConditionalAccessFilterable) {
    return m.deviceFilter
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ConditionalAccessDevices) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["deviceFilter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateConditionalAccessFilterFromDiscriminatorValue , m.SetDeviceFilter)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ConditionalAccessDevices) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ConditionalAccessDevices) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("deviceFilter", m.GetDeviceFilter())
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
func (m *ConditionalAccessDevices) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDeviceFilter sets the deviceFilter property value. Filter that defines the dynamic-device-syntax rule to include/exclude devices. A filter can use device properties (such as extension attributes) to include/exclude them.
func (m *ConditionalAccessDevices) SetDeviceFilter(value ConditionalAccessFilterable)() {
    m.deviceFilter = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessDevices) SetOdataType(value *string)() {
    m.odataType = value
}
