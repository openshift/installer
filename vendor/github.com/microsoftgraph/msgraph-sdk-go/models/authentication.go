package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Authentication 
type Authentication struct {
    Entity
    // The email address registered to a user for authentication.
    emailMethods []EmailAuthenticationMethodable
    // Represents the FIDO2 security keys registered to a user for authentication.
    fido2Methods []Fido2AuthenticationMethodable
    // Represents all authentication methods registered to a user.
    methods []AuthenticationMethodable
    // The details of the Microsoft Authenticator app registered to a user for authentication.
    microsoftAuthenticatorMethods []MicrosoftAuthenticatorAuthenticationMethodable
    // Represents the status of a long-running operation.
    operations []LongRunningOperationable
    // Represents the password that's registered to a user for authentication. For security, the password itself will never be returned in the object, but action can be taken to reset a password.
    passwordMethods []PasswordAuthenticationMethodable
    // The phone numbers registered to a user for authentication.
    phoneMethods []PhoneAuthenticationMethodable
    // The software OATH TOTP applications registered to a user for authentication.
    softwareOathMethods []SoftwareOathAuthenticationMethodable
    // Represents a Temporary Access Pass registered to a user for authentication through time-limited passcodes.
    temporaryAccessPassMethods []TemporaryAccessPassAuthenticationMethodable
    // Represents the Windows Hello for Business authentication method registered to a user for authentication.
    windowsHelloForBusinessMethods []WindowsHelloForBusinessAuthenticationMethodable
}
// NewAuthentication instantiates a new authentication and sets the default values.
func NewAuthentication()(*Authentication) {
    m := &Authentication{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAuthenticationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthentication(), nil
}
// GetEmailMethods gets the emailMethods property value. The email address registered to a user for authentication.
func (m *Authentication) GetEmailMethods()([]EmailAuthenticationMethodable) {
    return m.emailMethods
}
// GetFido2Methods gets the fido2Methods property value. Represents the FIDO2 security keys registered to a user for authentication.
func (m *Authentication) GetFido2Methods()([]Fido2AuthenticationMethodable) {
    return m.fido2Methods
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Authentication) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["emailMethods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEmailAuthenticationMethodFromDiscriminatorValue , m.SetEmailMethods)
    res["fido2Methods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateFido2AuthenticationMethodFromDiscriminatorValue , m.SetFido2Methods)
    res["methods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAuthenticationMethodFromDiscriminatorValue , m.SetMethods)
    res["microsoftAuthenticatorMethods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMicrosoftAuthenticatorAuthenticationMethodFromDiscriminatorValue , m.SetMicrosoftAuthenticatorMethods)
    res["operations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateLongRunningOperationFromDiscriminatorValue , m.SetOperations)
    res["passwordMethods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePasswordAuthenticationMethodFromDiscriminatorValue , m.SetPasswordMethods)
    res["phoneMethods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePhoneAuthenticationMethodFromDiscriminatorValue , m.SetPhoneMethods)
    res["softwareOathMethods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSoftwareOathAuthenticationMethodFromDiscriminatorValue , m.SetSoftwareOathMethods)
    res["temporaryAccessPassMethods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTemporaryAccessPassAuthenticationMethodFromDiscriminatorValue , m.SetTemporaryAccessPassMethods)
    res["windowsHelloForBusinessMethods"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsHelloForBusinessAuthenticationMethodFromDiscriminatorValue , m.SetWindowsHelloForBusinessMethods)
    return res
}
// GetMethods gets the methods property value. Represents all authentication methods registered to a user.
func (m *Authentication) GetMethods()([]AuthenticationMethodable) {
    return m.methods
}
// GetMicrosoftAuthenticatorMethods gets the microsoftAuthenticatorMethods property value. The details of the Microsoft Authenticator app registered to a user for authentication.
func (m *Authentication) GetMicrosoftAuthenticatorMethods()([]MicrosoftAuthenticatorAuthenticationMethodable) {
    return m.microsoftAuthenticatorMethods
}
// GetOperations gets the operations property value. Represents the status of a long-running operation.
func (m *Authentication) GetOperations()([]LongRunningOperationable) {
    return m.operations
}
// GetPasswordMethods gets the passwordMethods property value. Represents the password that's registered to a user for authentication. For security, the password itself will never be returned in the object, but action can be taken to reset a password.
func (m *Authentication) GetPasswordMethods()([]PasswordAuthenticationMethodable) {
    return m.passwordMethods
}
// GetPhoneMethods gets the phoneMethods property value. The phone numbers registered to a user for authentication.
func (m *Authentication) GetPhoneMethods()([]PhoneAuthenticationMethodable) {
    return m.phoneMethods
}
// GetSoftwareOathMethods gets the softwareOathMethods property value. The software OATH TOTP applications registered to a user for authentication.
func (m *Authentication) GetSoftwareOathMethods()([]SoftwareOathAuthenticationMethodable) {
    return m.softwareOathMethods
}
// GetTemporaryAccessPassMethods gets the temporaryAccessPassMethods property value. Represents a Temporary Access Pass registered to a user for authentication through time-limited passcodes.
func (m *Authentication) GetTemporaryAccessPassMethods()([]TemporaryAccessPassAuthenticationMethodable) {
    return m.temporaryAccessPassMethods
}
// GetWindowsHelloForBusinessMethods gets the windowsHelloForBusinessMethods property value. Represents the Windows Hello for Business authentication method registered to a user for authentication.
func (m *Authentication) GetWindowsHelloForBusinessMethods()([]WindowsHelloForBusinessAuthenticationMethodable) {
    return m.windowsHelloForBusinessMethods
}
// Serialize serializes information the current object
func (m *Authentication) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetEmailMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEmailMethods())
        err = writer.WriteCollectionOfObjectValues("emailMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetFido2Methods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetFido2Methods())
        err = writer.WriteCollectionOfObjectValues("fido2Methods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMethods())
        err = writer.WriteCollectionOfObjectValues("methods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMicrosoftAuthenticatorMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMicrosoftAuthenticatorMethods())
        err = writer.WriteCollectionOfObjectValues("microsoftAuthenticatorMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOperations())
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPasswordMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPasswordMethods())
        err = writer.WriteCollectionOfObjectValues("passwordMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPhoneMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPhoneMethods())
        err = writer.WriteCollectionOfObjectValues("phoneMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSoftwareOathMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSoftwareOathMethods())
        err = writer.WriteCollectionOfObjectValues("softwareOathMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTemporaryAccessPassMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTemporaryAccessPassMethods())
        err = writer.WriteCollectionOfObjectValues("temporaryAccessPassMethods", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsHelloForBusinessMethods() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetWindowsHelloForBusinessMethods())
        err = writer.WriteCollectionOfObjectValues("windowsHelloForBusinessMethods", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEmailMethods sets the emailMethods property value. The email address registered to a user for authentication.
func (m *Authentication) SetEmailMethods(value []EmailAuthenticationMethodable)() {
    m.emailMethods = value
}
// SetFido2Methods sets the fido2Methods property value. Represents the FIDO2 security keys registered to a user for authentication.
func (m *Authentication) SetFido2Methods(value []Fido2AuthenticationMethodable)() {
    m.fido2Methods = value
}
// SetMethods sets the methods property value. Represents all authentication methods registered to a user.
func (m *Authentication) SetMethods(value []AuthenticationMethodable)() {
    m.methods = value
}
// SetMicrosoftAuthenticatorMethods sets the microsoftAuthenticatorMethods property value. The details of the Microsoft Authenticator app registered to a user for authentication.
func (m *Authentication) SetMicrosoftAuthenticatorMethods(value []MicrosoftAuthenticatorAuthenticationMethodable)() {
    m.microsoftAuthenticatorMethods = value
}
// SetOperations sets the operations property value. Represents the status of a long-running operation.
func (m *Authentication) SetOperations(value []LongRunningOperationable)() {
    m.operations = value
}
// SetPasswordMethods sets the passwordMethods property value. Represents the password that's registered to a user for authentication. For security, the password itself will never be returned in the object, but action can be taken to reset a password.
func (m *Authentication) SetPasswordMethods(value []PasswordAuthenticationMethodable)() {
    m.passwordMethods = value
}
// SetPhoneMethods sets the phoneMethods property value. The phone numbers registered to a user for authentication.
func (m *Authentication) SetPhoneMethods(value []PhoneAuthenticationMethodable)() {
    m.phoneMethods = value
}
// SetSoftwareOathMethods sets the softwareOathMethods property value. The software OATH TOTP applications registered to a user for authentication.
func (m *Authentication) SetSoftwareOathMethods(value []SoftwareOathAuthenticationMethodable)() {
    m.softwareOathMethods = value
}
// SetTemporaryAccessPassMethods sets the temporaryAccessPassMethods property value. Represents a Temporary Access Pass registered to a user for authentication through time-limited passcodes.
func (m *Authentication) SetTemporaryAccessPassMethods(value []TemporaryAccessPassAuthenticationMethodable)() {
    m.temporaryAccessPassMethods = value
}
// SetWindowsHelloForBusinessMethods sets the windowsHelloForBusinessMethods property value. Represents the Windows Hello for Business authentication method registered to a user for authentication.
func (m *Authentication) SetWindowsHelloForBusinessMethods(value []WindowsHelloForBusinessAuthenticationMethodable)() {
    m.windowsHelloForBusinessMethods = value
}
