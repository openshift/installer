package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AgreementFileLocalization provides operations to manage the collection of agreement entities.
type AgreementFileLocalization struct {
    AgreementFileProperties
    // Read-only. Customized versions of the terms of use agreement in the Azure AD tenant.
    versions []AgreementFileVersionable
}
// NewAgreementFileLocalization instantiates a new agreementFileLocalization and sets the default values.
func NewAgreementFileLocalization()(*AgreementFileLocalization) {
    m := &AgreementFileLocalization{
        AgreementFileProperties: *NewAgreementFileProperties(),
    }
    return m
}
// CreateAgreementFileLocalizationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAgreementFileLocalizationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreementFileLocalization(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AgreementFileLocalization) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AgreementFileProperties.GetFieldDeserializers()
    res["versions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAgreementFileVersionFromDiscriminatorValue , m.SetVersions)
    return res
}
// GetVersions gets the versions property value. Read-only. Customized versions of the terms of use agreement in the Azure AD tenant.
func (m *AgreementFileLocalization) GetVersions()([]AgreementFileVersionable) {
    return m.versions
}
// Serialize serializes information the current object
func (m *AgreementFileLocalization) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AgreementFileProperties.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetVersions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetVersions())
        err = writer.WriteCollectionOfObjectValues("versions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetVersions sets the versions property value. Read-only. Customized versions of the terms of use agreement in the Azure AD tenant.
func (m *AgreementFileLocalization) SetVersions(value []AgreementFileVersionable)() {
    m.versions = value
}
