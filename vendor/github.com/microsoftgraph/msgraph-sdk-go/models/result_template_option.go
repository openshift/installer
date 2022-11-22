package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ResultTemplateOption 
type ResultTemplateOption struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Indicates whether search display layouts are enabled. If enabled, the user will get the result template to render the search results content in the resultTemplates property of the response. The result template is based on Adaptive Cards. Optional.
    enableResultTemplate *bool
    // The OdataType property
    odataType *string
}
// NewResultTemplateOption instantiates a new resultTemplateOption and sets the default values.
func NewResultTemplateOption()(*ResultTemplateOption) {
    m := &ResultTemplateOption{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateResultTemplateOptionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateResultTemplateOptionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewResultTemplateOption(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ResultTemplateOption) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetEnableResultTemplate gets the enableResultTemplate property value. Indicates whether search display layouts are enabled. If enabled, the user will get the result template to render the search results content in the resultTemplates property of the response. The result template is based on Adaptive Cards. Optional.
func (m *ResultTemplateOption) GetEnableResultTemplate()(*bool) {
    return m.enableResultTemplate
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ResultTemplateOption) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["enableResultTemplate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEnableResultTemplate)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ResultTemplateOption) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ResultTemplateOption) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("enableResultTemplate", m.GetEnableResultTemplate())
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
func (m *ResultTemplateOption) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetEnableResultTemplate sets the enableResultTemplate property value. Indicates whether search display layouts are enabled. If enabled, the user will get the result template to render the search results content in the resultTemplates property of the response. The result template is based on Adaptive Cards. Optional.
func (m *ResultTemplateOption) SetEnableResultTemplate(value *bool)() {
    m.enableResultTemplate = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ResultTemplateOption) SetOdataType(value *string)() {
    m.odataType = value
}
