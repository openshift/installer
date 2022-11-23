package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Fido2AuthenticationMethodConfiguration 
type Fido2AuthenticationMethodConfiguration struct {
    AuthenticationMethodConfiguration
    // A collection of users or groups who are enabled to use the authentication method.
    includeTargets []AuthenticationMethodTargetable
    // Determines whether attestation must be enforced for FIDO2 security key registration.
    isAttestationEnforced *bool
    // Determines if users can register new FIDO2 security keys.
    isSelfServiceRegistrationAllowed *bool
    // Controls whether key restrictions are enforced on FIDO2 security keys, either allowing or disallowing certain key types as defined by Authenticator Attestation GUID (AAGUID), an identifier that indicates the type (e.g. make and model) of the authenticator.
    keyRestrictions Fido2KeyRestrictionsable
}
// NewFido2AuthenticationMethodConfiguration instantiates a new Fido2AuthenticationMethodConfiguration and sets the default values.
func NewFido2AuthenticationMethodConfiguration()(*Fido2AuthenticationMethodConfiguration) {
    m := &Fido2AuthenticationMethodConfiguration{
        AuthenticationMethodConfiguration: *NewAuthenticationMethodConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.fido2AuthenticationMethodConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateFido2AuthenticationMethodConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateFido2AuthenticationMethodConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFido2AuthenticationMethodConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Fido2AuthenticationMethodConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethodConfiguration.GetFieldDeserializers()
    res["includeTargets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAuthenticationMethodTargetFromDiscriminatorValue , m.SetIncludeTargets)
    res["isAttestationEnforced"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsAttestationEnforced)
    res["isSelfServiceRegistrationAllowed"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsSelfServiceRegistrationAllowed)
    res["keyRestrictions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateFido2KeyRestrictionsFromDiscriminatorValue , m.SetKeyRestrictions)
    return res
}
// GetIncludeTargets gets the includeTargets property value. A collection of users or groups who are enabled to use the authentication method.
func (m *Fido2AuthenticationMethodConfiguration) GetIncludeTargets()([]AuthenticationMethodTargetable) {
    return m.includeTargets
}
// GetIsAttestationEnforced gets the isAttestationEnforced property value. Determines whether attestation must be enforced for FIDO2 security key registration.
func (m *Fido2AuthenticationMethodConfiguration) GetIsAttestationEnforced()(*bool) {
    return m.isAttestationEnforced
}
// GetIsSelfServiceRegistrationAllowed gets the isSelfServiceRegistrationAllowed property value. Determines if users can register new FIDO2 security keys.
func (m *Fido2AuthenticationMethodConfiguration) GetIsSelfServiceRegistrationAllowed()(*bool) {
    return m.isSelfServiceRegistrationAllowed
}
// GetKeyRestrictions gets the keyRestrictions property value. Controls whether key restrictions are enforced on FIDO2 security keys, either allowing or disallowing certain key types as defined by Authenticator Attestation GUID (AAGUID), an identifier that indicates the type (e.g. make and model) of the authenticator.
func (m *Fido2AuthenticationMethodConfiguration) GetKeyRestrictions()(Fido2KeyRestrictionsable) {
    return m.keyRestrictions
}
// Serialize serializes information the current object
func (m *Fido2AuthenticationMethodConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethodConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetIncludeTargets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetIncludeTargets())
        err = writer.WriteCollectionOfObjectValues("includeTargets", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAttestationEnforced", m.GetIsAttestationEnforced())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSelfServiceRegistrationAllowed", m.GetIsSelfServiceRegistrationAllowed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("keyRestrictions", m.GetKeyRestrictions())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIncludeTargets sets the includeTargets property value. A collection of users or groups who are enabled to use the authentication method.
func (m *Fido2AuthenticationMethodConfiguration) SetIncludeTargets(value []AuthenticationMethodTargetable)() {
    m.includeTargets = value
}
// SetIsAttestationEnforced sets the isAttestationEnforced property value. Determines whether attestation must be enforced for FIDO2 security key registration.
func (m *Fido2AuthenticationMethodConfiguration) SetIsAttestationEnforced(value *bool)() {
    m.isAttestationEnforced = value
}
// SetIsSelfServiceRegistrationAllowed sets the isSelfServiceRegistrationAllowed property value. Determines if users can register new FIDO2 security keys.
func (m *Fido2AuthenticationMethodConfiguration) SetIsSelfServiceRegistrationAllowed(value *bool)() {
    m.isSelfServiceRegistrationAllowed = value
}
// SetKeyRestrictions sets the keyRestrictions property value. Controls whether key restrictions are enforced on FIDO2 security keys, either allowing or disallowing certain key types as defined by Authenticator Attestation GUID (AAGUID), an identifier that indicates the type (e.g. make and model) of the authenticator.
func (m *Fido2AuthenticationMethodConfiguration) SetKeyRestrictions(value Fido2KeyRestrictionsable)() {
    m.keyRestrictions = value
}
