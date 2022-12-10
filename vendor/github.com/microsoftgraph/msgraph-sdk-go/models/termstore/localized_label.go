package termstore

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// LocalizedLabel 
type LocalizedLabel struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Indicates whether the label is the default label.
    isDefault *bool
    // The language tag for the label.
    languageTag *string
    // The name of the label.
    name *string
    // The OdataType property
    odataType *string
}
// NewLocalizedLabel instantiates a new localizedLabel and sets the default values.
func NewLocalizedLabel()(*LocalizedLabel) {
    m := &LocalizedLabel{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateLocalizedLabelFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateLocalizedLabelFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLocalizedLabel(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *LocalizedLabel) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *LocalizedLabel) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isDefault"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsDefault)
    res["languageTag"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLanguageTag)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetIsDefault gets the isDefault property value. Indicates whether the label is the default label.
func (m *LocalizedLabel) GetIsDefault()(*bool) {
    return m.isDefault
}
// GetLanguageTag gets the languageTag property value. The language tag for the label.
func (m *LocalizedLabel) GetLanguageTag()(*string) {
    return m.languageTag
}
// GetName gets the name property value. The name of the label.
func (m *LocalizedLabel) GetName()(*string) {
    return m.name
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *LocalizedLabel) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *LocalizedLabel) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isDefault", m.GetIsDefault())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("languageTag", m.GetLanguageTag())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
func (m *LocalizedLabel) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIsDefault sets the isDefault property value. Indicates whether the label is the default label.
func (m *LocalizedLabel) SetIsDefault(value *bool)() {
    m.isDefault = value
}
// SetLanguageTag sets the languageTag property value. The language tag for the label.
func (m *LocalizedLabel) SetLanguageTag(value *string)() {
    m.languageTag = value
}
// SetName sets the name property value. The name of the label.
func (m *LocalizedLabel) SetName(value *string)() {
    m.name = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *LocalizedLabel) SetOdataType(value *string)() {
    m.odataType = value
}
