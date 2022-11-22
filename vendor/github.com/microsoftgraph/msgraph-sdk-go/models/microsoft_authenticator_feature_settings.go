package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MicrosoftAuthenticatorFeatureSettings 
type MicrosoftAuthenticatorFeatureSettings struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Determines whether the user's Authenticator app will show them the client app they are signing into.
    displayAppInformationRequiredState AuthenticationMethodFeatureConfigurationable
    // Determines whether the user's Authenticator app will show them the geographic location of where the authentication request originated from.
    displayLocationInformationRequiredState AuthenticationMethodFeatureConfigurationable
    // The OdataType property
    odataType *string
}
// NewMicrosoftAuthenticatorFeatureSettings instantiates a new microsoftAuthenticatorFeatureSettings and sets the default values.
func NewMicrosoftAuthenticatorFeatureSettings()(*MicrosoftAuthenticatorFeatureSettings) {
    m := &MicrosoftAuthenticatorFeatureSettings{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateMicrosoftAuthenticatorFeatureSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMicrosoftAuthenticatorFeatureSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMicrosoftAuthenticatorFeatureSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MicrosoftAuthenticatorFeatureSettings) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDisplayAppInformationRequiredState gets the displayAppInformationRequiredState property value. Determines whether the user's Authenticator app will show them the client app they are signing into.
func (m *MicrosoftAuthenticatorFeatureSettings) GetDisplayAppInformationRequiredState()(AuthenticationMethodFeatureConfigurationable) {
    return m.displayAppInformationRequiredState
}
// GetDisplayLocationInformationRequiredState gets the displayLocationInformationRequiredState property value. Determines whether the user's Authenticator app will show them the geographic location of where the authentication request originated from.
func (m *MicrosoftAuthenticatorFeatureSettings) GetDisplayLocationInformationRequiredState()(AuthenticationMethodFeatureConfigurationable) {
    return m.displayLocationInformationRequiredState
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MicrosoftAuthenticatorFeatureSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["displayAppInformationRequiredState"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateAuthenticationMethodFeatureConfigurationFromDiscriminatorValue , m.SetDisplayAppInformationRequiredState)
    res["displayLocationInformationRequiredState"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateAuthenticationMethodFeatureConfigurationFromDiscriminatorValue , m.SetDisplayLocationInformationRequiredState)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *MicrosoftAuthenticatorFeatureSettings) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *MicrosoftAuthenticatorFeatureSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("displayAppInformationRequiredState", m.GetDisplayAppInformationRequiredState())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("displayLocationInformationRequiredState", m.GetDisplayLocationInformationRequiredState())
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
func (m *MicrosoftAuthenticatorFeatureSettings) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDisplayAppInformationRequiredState sets the displayAppInformationRequiredState property value. Determines whether the user's Authenticator app will show them the client app they are signing into.
func (m *MicrosoftAuthenticatorFeatureSettings) SetDisplayAppInformationRequiredState(value AuthenticationMethodFeatureConfigurationable)() {
    m.displayAppInformationRequiredState = value
}
// SetDisplayLocationInformationRequiredState sets the displayLocationInformationRequiredState property value. Determines whether the user's Authenticator app will show them the geographic location of where the authentication request originated from.
func (m *MicrosoftAuthenticatorFeatureSettings) SetDisplayLocationInformationRequiredState(value AuthenticationMethodFeatureConfigurationable)() {
    m.displayLocationInformationRequiredState = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MicrosoftAuthenticatorFeatureSettings) SetOdataType(value *string)() {
    m.odataType = value
}
