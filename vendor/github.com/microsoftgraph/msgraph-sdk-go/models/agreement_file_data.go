package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AgreementFileData 
type AgreementFileData struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Data that represents the terms of use PDF document. Read-only.
    data []byte
    // The OdataType property
    odataType *string
}
// NewAgreementFileData instantiates a new agreementFileData and sets the default values.
func NewAgreementFileData()(*AgreementFileData) {
    m := &AgreementFileData{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAgreementFileDataFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAgreementFileDataFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreementFileData(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AgreementFileData) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetData gets the data property value. Data that represents the terms of use PDF document. Read-only.
func (m *AgreementFileData) GetData()([]byte) {
    return m.data
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AgreementFileData) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["data"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteArrayValue(m.SetData)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AgreementFileData) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *AgreementFileData) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteByteArrayValue("data", m.GetData())
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
func (m *AgreementFileData) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetData sets the data property value. Data that represents the terms of use PDF document. Read-only.
func (m *AgreementFileData) SetData(value []byte)() {
    m.data = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AgreementFileData) SetOdataType(value *string)() {
    m.odataType = value
}
