package wipemanagedappregistrationsbydevicetag

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WipeManagedAppRegistrationsByDeviceTagPostRequestBody provides operations to call the wipeManagedAppRegistrationsByDeviceTag method.
type WipeManagedAppRegistrationsByDeviceTagPostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The deviceTag property
    deviceTag *string
}
// NewWipeManagedAppRegistrationsByDeviceTagPostRequestBody instantiates a new wipeManagedAppRegistrationsByDeviceTagPostRequestBody and sets the default values.
func NewWipeManagedAppRegistrationsByDeviceTagPostRequestBody()(*WipeManagedAppRegistrationsByDeviceTagPostRequestBody) {
    m := &WipeManagedAppRegistrationsByDeviceTagPostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWipeManagedAppRegistrationsByDeviceTagPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWipeManagedAppRegistrationsByDeviceTagPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWipeManagedAppRegistrationsByDeviceTagPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *WipeManagedAppRegistrationsByDeviceTagPostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDeviceTag gets the deviceTag property value. The deviceTag property
func (m *WipeManagedAppRegistrationsByDeviceTagPostRequestBody) GetDeviceTag()(*string) {
    return m.deviceTag
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WipeManagedAppRegistrationsByDeviceTagPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["deviceTag"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceTag)
    return res
}
// Serialize serializes information the current object
func (m *WipeManagedAppRegistrationsByDeviceTagPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("deviceTag", m.GetDeviceTag())
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
func (m *WipeManagedAppRegistrationsByDeviceTagPostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDeviceTag sets the deviceTag property value. The deviceTag property
func (m *WipeManagedAppRegistrationsByDeviceTagPostRequestBody) SetDeviceTag(value *string)() {
    m.deviceTag = value
}
