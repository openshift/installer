package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrintSettings 
type PrintSettings struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Specifies whether document conversion is enabled for the tenant. If document conversion is enabled, Universal Print service will automatically convert documents into a format compatible with the printer (xps to pdf) when needed.
    documentConversionEnabled *bool
    // The OdataType property
    odataType *string
}
// NewPrintSettings instantiates a new printSettings and sets the default values.
func NewPrintSettings()(*PrintSettings) {
    m := &PrintSettings{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePrintSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrintSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PrintSettings) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDocumentConversionEnabled gets the documentConversionEnabled property value. Specifies whether document conversion is enabled for the tenant. If document conversion is enabled, Universal Print service will automatically convert documents into a format compatible with the printer (xps to pdf) when needed.
func (m *PrintSettings) GetDocumentConversionEnabled()(*bool) {
    return m.documentConversionEnabled
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PrintSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["documentConversionEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetDocumentConversionEnabled)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *PrintSettings) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *PrintSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("documentConversionEnabled", m.GetDocumentConversionEnabled())
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
func (m *PrintSettings) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDocumentConversionEnabled sets the documentConversionEnabled property value. Specifies whether document conversion is enabled for the tenant. If document conversion is enabled, Universal Print service will automatically convert documents into a format compatible with the printer (xps to pdf) when needed.
func (m *PrintSettings) SetDocumentConversionEnabled(value *bool)() {
    m.documentConversionEnabled = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PrintSettings) SetOdataType(value *string)() {
    m.odataType = value
}
