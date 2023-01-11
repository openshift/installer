package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OmaSettingBoolean 
type OmaSettingBoolean struct {
    OmaSetting
    // Value.
    value *bool
}
// NewOmaSettingBoolean instantiates a new OmaSettingBoolean and sets the default values.
func NewOmaSettingBoolean()(*OmaSettingBoolean) {
    m := &OmaSettingBoolean{
        OmaSetting: *NewOmaSetting(),
    }
    odataTypeValue := "#microsoft.graph.omaSettingBoolean";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateOmaSettingBooleanFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOmaSettingBooleanFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOmaSettingBoolean(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *OmaSettingBoolean) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OmaSetting.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetValue)
    return res
}
// GetValue gets the value property value. Value.
func (m *OmaSettingBoolean) GetValue()(*bool) {
    return m.value
}
// Serialize serializes information the current object
func (m *OmaSettingBoolean) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OmaSetting.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("value", m.GetValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. Value.
func (m *OmaSettingBoolean) SetValue(value *bool)() {
    m.value = value
}
