package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SettingTemplateValue 
type SettingTemplateValue struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Default value for the setting.
    defaultValue *string
    // Description of the setting.
    description *string
    // Name of the setting.
    name *string
    // The OdataType property
    odataType *string
    // Type of the setting.
    type_escaped *string
}
// NewSettingTemplateValue instantiates a new settingTemplateValue and sets the default values.
func NewSettingTemplateValue()(*SettingTemplateValue) {
    m := &SettingTemplateValue{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSettingTemplateValueFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSettingTemplateValueFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSettingTemplateValue(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SettingTemplateValue) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDefaultValue gets the defaultValue property value. Default value for the setting.
func (m *SettingTemplateValue) GetDefaultValue()(*string) {
    return m.defaultValue
}
// GetDescription gets the description property value. Description of the setting.
func (m *SettingTemplateValue) GetDescription()(*string) {
    return m.description
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SettingTemplateValue) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["defaultValue"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDefaultValue)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetType)
    return res
}
// GetName gets the name property value. Name of the setting.
func (m *SettingTemplateValue) GetName()(*string) {
    return m.name
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SettingTemplateValue) GetOdataType()(*string) {
    return m.odataType
}
// GetType gets the type property value. Type of the setting.
func (m *SettingTemplateValue) GetType()(*string) {
    return m.type_escaped
}
// Serialize serializes information the current object
func (m *SettingTemplateValue) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("defaultValue", m.GetDefaultValue())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("description", m.GetDescription())
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
        err := writer.WriteStringValue("type", m.GetType())
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
func (m *SettingTemplateValue) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDefaultValue sets the defaultValue property value. Default value for the setting.
func (m *SettingTemplateValue) SetDefaultValue(value *string)() {
    m.defaultValue = value
}
// SetDescription sets the description property value. Description of the setting.
func (m *SettingTemplateValue) SetDescription(value *string)() {
    m.description = value
}
// SetName sets the name property value. Name of the setting.
func (m *SettingTemplateValue) SetName(value *string)() {
    m.name = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SettingTemplateValue) SetOdataType(value *string)() {
    m.odataType = value
}
// SetType sets the type property value. Type of the setting.
func (m *SettingTemplateValue) SetType(value *string)() {
    m.type_escaped = value
}
