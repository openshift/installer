package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EmailAuthenticationMethodConfiguration 
type EmailAuthenticationMethodConfiguration struct {
    AuthenticationMethodConfiguration
    // Determines whether email OTP is usable by external users for authentication. Possible values are: default, enabled, disabled, unknownFutureValue. Tenants in the default state who did not use public preview will automatically have email OTP enabled beginning in October 2021.
    allowExternalIdToUseEmailOtp *ExternalEmailOtpState
    // A collection of users or groups who are enabled to use the authentication method.
    includeTargets []AuthenticationMethodTargetable
}
// NewEmailAuthenticationMethodConfiguration instantiates a new EmailAuthenticationMethodConfiguration and sets the default values.
func NewEmailAuthenticationMethodConfiguration()(*EmailAuthenticationMethodConfiguration) {
    m := &EmailAuthenticationMethodConfiguration{
        AuthenticationMethodConfiguration: *NewAuthenticationMethodConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.emailAuthenticationMethodConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEmailAuthenticationMethodConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEmailAuthenticationMethodConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmailAuthenticationMethodConfiguration(), nil
}
// GetAllowExternalIdToUseEmailOtp gets the allowExternalIdToUseEmailOtp property value. Determines whether email OTP is usable by external users for authentication. Possible values are: default, enabled, disabled, unknownFutureValue. Tenants in the default state who did not use public preview will automatically have email OTP enabled beginning in October 2021.
func (m *EmailAuthenticationMethodConfiguration) GetAllowExternalIdToUseEmailOtp()(*ExternalEmailOtpState) {
    return m.allowExternalIdToUseEmailOtp
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EmailAuthenticationMethodConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethodConfiguration.GetFieldDeserializers()
    res["allowExternalIdToUseEmailOtp"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseExternalEmailOtpState , m.SetAllowExternalIdToUseEmailOtp)
    res["includeTargets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAuthenticationMethodTargetFromDiscriminatorValue , m.SetIncludeTargets)
    return res
}
// GetIncludeTargets gets the includeTargets property value. A collection of users or groups who are enabled to use the authentication method.
func (m *EmailAuthenticationMethodConfiguration) GetIncludeTargets()([]AuthenticationMethodTargetable) {
    return m.includeTargets
}
// Serialize serializes information the current object
func (m *EmailAuthenticationMethodConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethodConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowExternalIdToUseEmailOtp() != nil {
        cast := (*m.GetAllowExternalIdToUseEmailOtp()).String()
        err = writer.WriteStringValue("allowExternalIdToUseEmailOtp", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetIncludeTargets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetIncludeTargets())
        err = writer.WriteCollectionOfObjectValues("includeTargets", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowExternalIdToUseEmailOtp sets the allowExternalIdToUseEmailOtp property value. Determines whether email OTP is usable by external users for authentication. Possible values are: default, enabled, disabled, unknownFutureValue. Tenants in the default state who did not use public preview will automatically have email OTP enabled beginning in October 2021.
func (m *EmailAuthenticationMethodConfiguration) SetAllowExternalIdToUseEmailOtp(value *ExternalEmailOtpState)() {
    m.allowExternalIdToUseEmailOtp = value
}
// SetIncludeTargets sets the includeTargets property value. A collection of users or groups who are enabled to use the authentication method.
func (m *EmailAuthenticationMethodConfiguration) SetIncludeTargets(value []AuthenticationMethodTargetable)() {
    m.includeTargets = value
}
