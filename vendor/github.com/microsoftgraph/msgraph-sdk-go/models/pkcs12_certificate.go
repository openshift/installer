package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Pkcs12Certificate 
type Pkcs12Certificate struct {
    ApiAuthenticationConfigurationBase
    // The password for the pfx file. Required. If no password is used, you must still provide a value of ''.
    password *string
    // Represents the pfx content that is sent. The value should be a base-64 encoded version of the actual certificate content. Required.
    pkcs12Value *string
}
// NewPkcs12Certificate instantiates a new Pkcs12Certificate and sets the default values.
func NewPkcs12Certificate()(*Pkcs12Certificate) {
    m := &Pkcs12Certificate{
        ApiAuthenticationConfigurationBase: *NewApiAuthenticationConfigurationBase(),
    }
    odataTypeValue := "#microsoft.graph.pkcs12Certificate";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreatePkcs12CertificateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePkcs12CertificateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPkcs12Certificate(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Pkcs12Certificate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ApiAuthenticationConfigurationBase.GetFieldDeserializers()
    res["password"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPassword)
    res["pkcs12Value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPkcs12Value)
    return res
}
// GetPassword gets the password property value. The password for the pfx file. Required. If no password is used, you must still provide a value of ''.
func (m *Pkcs12Certificate) GetPassword()(*string) {
    return m.password
}
// GetPkcs12Value gets the pkcs12Value property value. Represents the pfx content that is sent. The value should be a base-64 encoded version of the actual certificate content. Required.
func (m *Pkcs12Certificate) GetPkcs12Value()(*string) {
    return m.pkcs12Value
}
// Serialize serializes information the current object
func (m *Pkcs12Certificate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteStringValue("pkcs12Value", m.GetPkcs12Value())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPassword sets the password property value. The password for the pfx file. Required. If no password is used, you must still provide a value of ''.
func (m *Pkcs12Certificate) SetPassword(value *string)() {
    m.password = value
}
// SetPkcs12Value sets the pkcs12Value property value. Represents the pfx content that is sent. The value should be a base-64 encoded version of the actual certificate content. Required.
func (m *Pkcs12Certificate) SetPkcs12Value(value *string)() {
    m.pkcs12Value = value
}
