package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EmailIdentity 
type EmailIdentity struct {
    Identity
    // Email address of the user.
    email *string
}
// NewEmailIdentity instantiates a new EmailIdentity and sets the default values.
func NewEmailIdentity()(*EmailIdentity) {
    m := &EmailIdentity{
        Identity: *NewIdentity(),
    }
    odataTypeValue := "#microsoft.graph.emailIdentity";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEmailIdentityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEmailIdentityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmailIdentity(), nil
}
// GetEmail gets the email property value. Email address of the user.
func (m *EmailIdentity) GetEmail()(*string) {
    return m.email
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EmailIdentity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Identity.GetFieldDeserializers()
    res["email"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEmail)
    return res
}
// Serialize serializes information the current object
func (m *EmailIdentity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Identity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEmail sets the email property value. Email address of the user.
func (m *EmailIdentity) SetEmail(value *string)() {
    m.email = value
}
