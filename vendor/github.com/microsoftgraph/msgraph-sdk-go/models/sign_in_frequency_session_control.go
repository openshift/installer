package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SignInFrequencySessionControl 
type SignInFrequencySessionControl struct {
    ConditionalAccessSessionControl
    // The possible values are primaryAndSecondaryAuthentication, secondaryAuthentication, unknownFutureValue.
    authenticationType *SignInFrequencyAuthenticationType
    // The possible values are timeBased, everyTime, unknownFutureValue.
    frequencyInterval *SignInFrequencyInterval
    // Possible values are: days, hours.
    type_escaped *SigninFrequencyType
    // The number of days or hours.
    value *int32
}
// NewSignInFrequencySessionControl instantiates a new SignInFrequencySessionControl and sets the default values.
func NewSignInFrequencySessionControl()(*SignInFrequencySessionControl) {
    m := &SignInFrequencySessionControl{
        ConditionalAccessSessionControl: *NewConditionalAccessSessionControl(),
    }
    odataTypeValue := "#microsoft.graph.signInFrequencySessionControl";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateSignInFrequencySessionControlFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSignInFrequencySessionControlFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSignInFrequencySessionControl(), nil
}
// GetAuthenticationType gets the authenticationType property value. The possible values are primaryAndSecondaryAuthentication, secondaryAuthentication, unknownFutureValue.
func (m *SignInFrequencySessionControl) GetAuthenticationType()(*SignInFrequencyAuthenticationType) {
    return m.authenticationType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SignInFrequencySessionControl) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConditionalAccessSessionControl.GetFieldDeserializers()
    res["authenticationType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseSignInFrequencyAuthenticationType , m.SetAuthenticationType)
    res["frequencyInterval"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseSignInFrequencyInterval , m.SetFrequencyInterval)
    res["type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseSigninFrequencyType , m.SetType)
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetValue)
    return res
}
// GetFrequencyInterval gets the frequencyInterval property value. The possible values are timeBased, everyTime, unknownFutureValue.
func (m *SignInFrequencySessionControl) GetFrequencyInterval()(*SignInFrequencyInterval) {
    return m.frequencyInterval
}
// GetType gets the type property value. Possible values are: days, hours.
func (m *SignInFrequencySessionControl) GetType()(*SigninFrequencyType) {
    return m.type_escaped
}
// GetValue gets the value property value. The number of days or hours.
func (m *SignInFrequencySessionControl) GetValue()(*int32) {
    return m.value
}
// Serialize serializes information the current object
func (m *SignInFrequencySessionControl) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConditionalAccessSessionControl.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuthenticationType() != nil {
        cast := (*m.GetAuthenticationType()).String()
        err = writer.WriteStringValue("authenticationType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFrequencyInterval() != nil {
        cast := (*m.GetFrequencyInterval()).String()
        err = writer.WriteStringValue("frequencyInterval", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetType() != nil {
        cast := (*m.GetType()).String()
        err = writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("value", m.GetValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationType sets the authenticationType property value. The possible values are primaryAndSecondaryAuthentication, secondaryAuthentication, unknownFutureValue.
func (m *SignInFrequencySessionControl) SetAuthenticationType(value *SignInFrequencyAuthenticationType)() {
    m.authenticationType = value
}
// SetFrequencyInterval sets the frequencyInterval property value. The possible values are timeBased, everyTime, unknownFutureValue.
func (m *SignInFrequencySessionControl) SetFrequencyInterval(value *SignInFrequencyInterval)() {
    m.frequencyInterval = value
}
// SetType sets the type property value. Possible values are: days, hours.
func (m *SignInFrequencySessionControl) SetType(value *SigninFrequencyType)() {
    m.type_escaped = value
}
// SetValue sets the value property value. The number of days or hours.
func (m *SignInFrequencySessionControl) SetValue(value *int32)() {
    m.value = value
}
