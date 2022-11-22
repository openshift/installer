package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AuthenticationMethodsPolicy provides operations to manage the authenticationMethodsPolicy singleton.
type AuthenticationMethodsPolicy struct {
    Entity
    // Represents the settings for each authentication method. Automatically expanded on GET /policies/authenticationMethodsPolicy.
    authenticationMethodConfigurations []AuthenticationMethodConfigurationable
    // A description of the policy. Read-only.
    description *string
    // The name of the policy. Read-only.
    displayName *string
    // The date and time of the last update to the policy. Read-only.
    lastModifiedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The version of the policy in use. Read-only.
    policyVersion *string
    // The reconfirmationInDays property
    reconfirmationInDays *int32
    // Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods.
    registrationEnforcement RegistrationEnforcementable
}
// NewAuthenticationMethodsPolicy instantiates a new authenticationMethodsPolicy and sets the default values.
func NewAuthenticationMethodsPolicy()(*AuthenticationMethodsPolicy) {
    m := &AuthenticationMethodsPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuthenticationMethodsPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAuthenticationMethodsPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthenticationMethodsPolicy(), nil
}
// GetAuthenticationMethodConfigurations gets the authenticationMethodConfigurations property value. Represents the settings for each authentication method. Automatically expanded on GET /policies/authenticationMethodsPolicy.
func (m *AuthenticationMethodsPolicy) GetAuthenticationMethodConfigurations()([]AuthenticationMethodConfigurationable) {
    return m.authenticationMethodConfigurations
}
// GetDescription gets the description property value. A description of the policy. Read-only.
func (m *AuthenticationMethodsPolicy) GetDescription()(*string) {
    return m.description
}
// GetDisplayName gets the displayName property value. The name of the policy. Read-only.
func (m *AuthenticationMethodsPolicy) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AuthenticationMethodsPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationMethodConfigurations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAuthenticationMethodConfigurationFromDiscriminatorValue , m.SetAuthenticationMethodConfigurations)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["lastModifiedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastModifiedDateTime)
    res["policyVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPolicyVersion)
    res["reconfirmationInDays"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetReconfirmationInDays)
    res["registrationEnforcement"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRegistrationEnforcementFromDiscriminatorValue , m.SetRegistrationEnforcement)
    return res
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time of the last update to the policy. Read-only.
func (m *AuthenticationMethodsPolicy) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastModifiedDateTime
}
// GetPolicyVersion gets the policyVersion property value. The version of the policy in use. Read-only.
func (m *AuthenticationMethodsPolicy) GetPolicyVersion()(*string) {
    return m.policyVersion
}
// GetReconfirmationInDays gets the reconfirmationInDays property value. The reconfirmationInDays property
func (m *AuthenticationMethodsPolicy) GetReconfirmationInDays()(*int32) {
    return m.reconfirmationInDays
}
// GetRegistrationEnforcement gets the registrationEnforcement property value. Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods.
func (m *AuthenticationMethodsPolicy) GetRegistrationEnforcement()(RegistrationEnforcementable) {
    return m.registrationEnforcement
}
// Serialize serializes information the current object
func (m *AuthenticationMethodsPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuthenticationMethodConfigurations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAuthenticationMethodConfigurations())
        err = writer.WriteCollectionOfObjectValues("authenticationMethodConfigurations", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("policyVersion", m.GetPolicyVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("reconfirmationInDays", m.GetReconfirmationInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("registrationEnforcement", m.GetRegistrationEnforcement())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationMethodConfigurations sets the authenticationMethodConfigurations property value. Represents the settings for each authentication method. Automatically expanded on GET /policies/authenticationMethodsPolicy.
func (m *AuthenticationMethodsPolicy) SetAuthenticationMethodConfigurations(value []AuthenticationMethodConfigurationable)() {
    m.authenticationMethodConfigurations = value
}
// SetDescription sets the description property value. A description of the policy. Read-only.
func (m *AuthenticationMethodsPolicy) SetDescription(value *string)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. The name of the policy. Read-only.
func (m *AuthenticationMethodsPolicy) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time of the last update to the policy. Read-only.
func (m *AuthenticationMethodsPolicy) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastModifiedDateTime = value
}
// SetPolicyVersion sets the policyVersion property value. The version of the policy in use. Read-only.
func (m *AuthenticationMethodsPolicy) SetPolicyVersion(value *string)() {
    m.policyVersion = value
}
// SetReconfirmationInDays sets the reconfirmationInDays property value. The reconfirmationInDays property
func (m *AuthenticationMethodsPolicy) SetReconfirmationInDays(value *int32)() {
    m.reconfirmationInDays = value
}
// SetRegistrationEnforcement sets the registrationEnforcement property value. Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods.
func (m *AuthenticationMethodsPolicy) SetRegistrationEnforcement(value RegistrationEnforcementable)() {
    m.registrationEnforcement = value
}
