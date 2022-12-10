package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// LobbyBypassSettings 
type LobbyBypassSettings struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Specifies whether or not to always let dial-in callers bypass the lobby. Optional.
    isDialInBypassEnabled *bool
    // The OdataType property
    odataType *string
    // Specifies the type of participants that are automatically admitted into a meeting, bypassing the lobby. Optional.
    scope *LobbyBypassScope
}
// NewLobbyBypassSettings instantiates a new lobbyBypassSettings and sets the default values.
func NewLobbyBypassSettings()(*LobbyBypassSettings) {
    m := &LobbyBypassSettings{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateLobbyBypassSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateLobbyBypassSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLobbyBypassSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *LobbyBypassSettings) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *LobbyBypassSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isDialInBypassEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsDialInBypassEnabled)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["scope"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseLobbyBypassScope , m.SetScope)
    return res
}
// GetIsDialInBypassEnabled gets the isDialInBypassEnabled property value. Specifies whether or not to always let dial-in callers bypass the lobby. Optional.
func (m *LobbyBypassSettings) GetIsDialInBypassEnabled()(*bool) {
    return m.isDialInBypassEnabled
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *LobbyBypassSettings) GetOdataType()(*string) {
    return m.odataType
}
// GetScope gets the scope property value. Specifies the type of participants that are automatically admitted into a meeting, bypassing the lobby. Optional.
func (m *LobbyBypassSettings) GetScope()(*LobbyBypassScope) {
    return m.scope
}
// Serialize serializes information the current object
func (m *LobbyBypassSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isDialInBypassEnabled", m.GetIsDialInBypassEnabled())
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
    if m.GetScope() != nil {
        cast := (*m.GetScope()).String()
        err := writer.WriteStringValue("scope", &cast)
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
func (m *LobbyBypassSettings) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIsDialInBypassEnabled sets the isDialInBypassEnabled property value. Specifies whether or not to always let dial-in callers bypass the lobby. Optional.
func (m *LobbyBypassSettings) SetIsDialInBypassEnabled(value *bool)() {
    m.isDialInBypassEnabled = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *LobbyBypassSettings) SetOdataType(value *string)() {
    m.odataType = value
}
// SetScope sets the scope property value. Specifies the type of participants that are automatically admitted into a meeting, bypassing the lobby. Optional.
func (m *LobbyBypassSettings) SetScope(value *LobbyBypassScope)() {
    m.scope = value
}
