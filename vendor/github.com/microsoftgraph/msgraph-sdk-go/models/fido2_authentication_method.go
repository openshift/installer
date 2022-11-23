package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Fido2AuthenticationMethod 
type Fido2AuthenticationMethod struct {
    AuthenticationMethod
    // Authenticator Attestation GUID, an identifier that indicates the type (e.g. make and model) of the authenticator.
    aaGuid *string
    // The attestation certificate(s) attached to this security key.
    attestationCertificates []string
    // The attestation level of this FIDO2 security key. Possible values are: attested, or notAttested.
    attestationLevel *AttestationLevel
    // The timestamp when this key was registered to the user.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The display name of the key as given by the user.
    displayName *string
    // The manufacturer-assigned model of the FIDO2 security key.
    model *string
}
// NewFido2AuthenticationMethod instantiates a new Fido2AuthenticationMethod and sets the default values.
func NewFido2AuthenticationMethod()(*Fido2AuthenticationMethod) {
    m := &Fido2AuthenticationMethod{
        AuthenticationMethod: *NewAuthenticationMethod(),
    }
    odataTypeValue := "#microsoft.graph.fido2AuthenticationMethod";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateFido2AuthenticationMethodFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateFido2AuthenticationMethodFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFido2AuthenticationMethod(), nil
}
// GetAaGuid gets the aaGuid property value. Authenticator Attestation GUID, an identifier that indicates the type (e.g. make and model) of the authenticator.
func (m *Fido2AuthenticationMethod) GetAaGuid()(*string) {
    return m.aaGuid
}
// GetAttestationCertificates gets the attestationCertificates property value. The attestation certificate(s) attached to this security key.
func (m *Fido2AuthenticationMethod) GetAttestationCertificates()([]string) {
    return m.attestationCertificates
}
// GetAttestationLevel gets the attestationLevel property value. The attestation level of this FIDO2 security key. Possible values are: attested, or notAttested.
func (m *Fido2AuthenticationMethod) GetAttestationLevel()(*AttestationLevel) {
    return m.attestationLevel
}
// GetCreatedDateTime gets the createdDateTime property value. The timestamp when this key was registered to the user.
func (m *Fido2AuthenticationMethod) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDisplayName gets the displayName property value. The display name of the key as given by the user.
func (m *Fido2AuthenticationMethod) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Fido2AuthenticationMethod) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethod.GetFieldDeserializers()
    res["aaGuid"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAaGuid)
    res["attestationCertificates"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetAttestationCertificates)
    res["attestationLevel"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseAttestationLevel , m.SetAttestationLevel)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["model"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetModel)
    return res
}
// GetModel gets the model property value. The manufacturer-assigned model of the FIDO2 security key.
func (m *Fido2AuthenticationMethod) GetModel()(*string) {
    return m.model
}
// Serialize serializes information the current object
func (m *Fido2AuthenticationMethod) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethod.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("aaGuid", m.GetAaGuid())
        if err != nil {
            return err
        }
    }
    if m.GetAttestationCertificates() != nil {
        err = writer.WriteCollectionOfStringValues("attestationCertificates", m.GetAttestationCertificates())
        if err != nil {
            return err
        }
    }
    if m.GetAttestationLevel() != nil {
        cast := (*m.GetAttestationLevel()).String()
        err = writer.WriteStringValue("attestationLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
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
        err = writer.WriteStringValue("model", m.GetModel())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAaGuid sets the aaGuid property value. Authenticator Attestation GUID, an identifier that indicates the type (e.g. make and model) of the authenticator.
func (m *Fido2AuthenticationMethod) SetAaGuid(value *string)() {
    m.aaGuid = value
}
// SetAttestationCertificates sets the attestationCertificates property value. The attestation certificate(s) attached to this security key.
func (m *Fido2AuthenticationMethod) SetAttestationCertificates(value []string)() {
    m.attestationCertificates = value
}
// SetAttestationLevel sets the attestationLevel property value. The attestation level of this FIDO2 security key. Possible values are: attested, or notAttested.
func (m *Fido2AuthenticationMethod) SetAttestationLevel(value *AttestationLevel)() {
    m.attestationLevel = value
}
// SetCreatedDateTime sets the createdDateTime property value. The timestamp when this key was registered to the user.
func (m *Fido2AuthenticationMethod) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDisplayName sets the displayName property value. The display name of the key as given by the user.
func (m *Fido2AuthenticationMethod) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetModel sets the model property value. The manufacturer-assigned model of the FIDO2 security key.
func (m *Fido2AuthenticationMethod) SetModel(value *string)() {
    m.model = value
}
