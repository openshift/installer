package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SignInFrequencySessionControl 
type SignInFrequencySessionControl struct {
    ConditionalAccessSessionControl
    // Possible values are: days, hours.
    TypeEscaped *SigninFrequencyType
}
// NewSignInFrequencySessionControl instantiates a new SignInFrequencySessionControl and sets the default values.
func NewSignInFrequencySessionControl()(*SignInFrequencySessionControl) {
    m := &SignInFrequencySessionControl{
        ConditionalAccessSessionControl: *NewConditionalAccessSessionControl(),
    }
    odataTypeValue := "#microsoft.graph.signInFrequencySessionControl"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSignInFrequencySessionControlFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSignInFrequencySessionControlFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSignInFrequencySessionControl(), nil
}
// GetAuthenticationType gets the authenticationType property value. The possible values are primaryAndSecondaryAuthentication, secondaryAuthentication, unknownFutureValue.
func (m *SignInFrequencySessionControl) GetAuthenticationType()(*SignInFrequencyAuthenticationType) {
    val, err := m.GetBackingStore().Get("authenticationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SignInFrequencyAuthenticationType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SignInFrequencySessionControl) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConditionalAccessSessionControl.GetFieldDeserializers()
    res["authenticationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSignInFrequencyAuthenticationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationType(val.(*SignInFrequencyAuthenticationType))
        }
        return nil
    }
    res["frequencyInterval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSignInFrequencyInterval)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFrequencyInterval(val.(*SignInFrequencyInterval))
        }
        return nil
    }
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSigninFrequencyType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetType(val.(*SigninFrequencyType))
        }
        return nil
    }
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValue(val)
        }
        return nil
    }
    return res
}
// GetFrequencyInterval gets the frequencyInterval property value. The possible values are timeBased, everyTime, unknownFutureValue.
func (m *SignInFrequencySessionControl) GetFrequencyInterval()(*SignInFrequencyInterval) {
    val, err := m.GetBackingStore().Get("frequencyInterval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SignInFrequencyInterval)
    }
    return nil
}
// GetType gets the type property value. Possible values are: days, hours.
func (m *SignInFrequencySessionControl) GetType()(*SigninFrequencyType) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SigninFrequencyType)
    }
    return nil
}
// GetValue gets the value property value. The number of days or hours.
func (m *SignInFrequencySessionControl) GetValue()(*int32) {
    val, err := m.GetBackingStore().Get("value")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
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
    err := m.GetBackingStore().Set("authenticationType", value)
    if err != nil {
        panic(err)
    }
}
// SetFrequencyInterval sets the frequencyInterval property value. The possible values are timeBased, everyTime, unknownFutureValue.
func (m *SignInFrequencySessionControl) SetFrequencyInterval(value *SignInFrequencyInterval)() {
    err := m.GetBackingStore().Set("frequencyInterval", value)
    if err != nil {
        panic(err)
    }
}
// SetType sets the type property value. Possible values are: days, hours.
func (m *SignInFrequencySessionControl) SetType(value *SigninFrequencyType)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
// SetValue sets the value property value. The number of days or hours.
func (m *SignInFrequencySessionControl) SetValue(value *int32)() {
    err := m.GetBackingStore().Set("value", value)
    if err != nil {
        panic(err)
    }
}
// SignInFrequencySessionControlable 
type SignInFrequencySessionControlable interface {
    ConditionalAccessSessionControlable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationType()(*SignInFrequencyAuthenticationType)
    GetFrequencyInterval()(*SignInFrequencyInterval)
    GetType()(*SigninFrequencyType)
    GetValue()(*int32)
    SetAuthenticationType(value *SignInFrequencyAuthenticationType)()
    SetFrequencyInterval(value *SignInFrequencyInterval)()
    SetType(value *SigninFrequencyType)()
    SetValue(value *int32)()
}
