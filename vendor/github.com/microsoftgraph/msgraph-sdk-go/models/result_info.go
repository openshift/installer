package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ResultInfo 
type ResultInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The result code.
    code *int32
    // The message.
    message *string
    // The OdataType property
    odataType *string
    // The result sub-code.
    subcode *int32
}
// NewResultInfo instantiates a new resultInfo and sets the default values.
func NewResultInfo()(*ResultInfo) {
    m := &ResultInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateResultInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateResultInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewResultInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ResultInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCode gets the code property value. The result code.
func (m *ResultInfo) GetCode()(*int32) {
    return m.code
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ResultInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["code"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetCode)
    res["message"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMessage)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["subcode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetSubcode)
    return res
}
// GetMessage gets the message property value. The message.
func (m *ResultInfo) GetMessage()(*string) {
    return m.message
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ResultInfo) GetOdataType()(*string) {
    return m.odataType
}
// GetSubcode gets the subcode property value. The result sub-code.
func (m *ResultInfo) GetSubcode()(*int32) {
    return m.subcode
}
// Serialize serializes information the current object
func (m *ResultInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("code", m.GetCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("message", m.GetMessage())
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
        err := writer.WriteInt32Value("subcode", m.GetSubcode())
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
func (m *ResultInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCode sets the code property value. The result code.
func (m *ResultInfo) SetCode(value *int32)() {
    m.code = value
}
// SetMessage sets the message property value. The message.
func (m *ResultInfo) SetMessage(value *string)() {
    m.message = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ResultInfo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetSubcode sets the subcode property value. The result sub-code.
func (m *ResultInfo) SetSubcode(value *int32)() {
    m.subcode = value
}
