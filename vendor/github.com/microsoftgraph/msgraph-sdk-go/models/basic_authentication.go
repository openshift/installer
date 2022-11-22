package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BasicAuthentication 
type BasicAuthentication struct {
    ApiAuthenticationConfigurationBase
    // The password. It is not returned in the responses.
    password *string
    // The username.
    username *string
}
// NewBasicAuthentication instantiates a new BasicAuthentication and sets the default values.
func NewBasicAuthentication()(*BasicAuthentication) {
    m := &BasicAuthentication{
        ApiAuthenticationConfigurationBase: *NewApiAuthenticationConfigurationBase(),
    }
    odataTypeValue := "#microsoft.graph.basicAuthentication";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateBasicAuthenticationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateBasicAuthenticationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBasicAuthentication(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *BasicAuthentication) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ApiAuthenticationConfigurationBase.GetFieldDeserializers()
    res["password"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPassword)
    res["username"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUsername)
    return res
}
// GetPassword gets the password property value. The password. It is not returned in the responses.
func (m *BasicAuthentication) GetPassword()(*string) {
    return m.password
}
// GetUsername gets the username property value. The username.
func (m *BasicAuthentication) GetUsername()(*string) {
    return m.username
}
// Serialize serializes information the current object
func (m *BasicAuthentication) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ApiAuthenticationConfigurationBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("password", m.GetPassword())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("username", m.GetUsername())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPassword sets the password property value. The password. It is not returned in the responses.
func (m *BasicAuthentication) SetPassword(value *string)() {
    m.password = value
}
// SetUsername sets the username property value. The username.
func (m *BasicAuthentication) SetUsername(value *string)() {
    m.username = value
}
