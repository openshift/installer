package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse 
type MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse struct {
    BaseCollectionPaginationCountResponse
    // The value property
    value []MicrosoftAuthenticatorAuthenticationMethodConfigurationable
}
// NewMicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse instantiates a new MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse and sets the default values.
func NewMicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse()(*MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse) {
    m := &MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateMicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMicrosoftAuthenticatorAuthenticationMethodConfigurationFromDiscriminatorValue , m.SetValue)
    return res
}
// GetValue gets the value property value. The value property
func (m *MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse) GetValue()([]MicrosoftAuthenticatorAuthenticationMethodConfigurationable) {
    return m.value
}
// Serialize serializes information the current object
func (m *MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseCollectionPaginationCountResponse.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetValue() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetValue())
        err = writer.WriteCollectionOfObjectValues("value", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. The value property
func (m *MicrosoftAuthenticatorAuthenticationMethodConfigurationCollectionResponse) SetValue(value []MicrosoftAuthenticatorAuthenticationMethodConfigurationable)() {
    m.value = value
}
