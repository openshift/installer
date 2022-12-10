package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AgreementFile 
type AgreementFile struct {
    AgreementFileProperties
    // The localized version of the terms of use agreement files attached to the agreement.
    localizations []AgreementFileLocalizationable
}
// NewAgreementFile instantiates a new AgreementFile and sets the default values.
func NewAgreementFile()(*AgreementFile) {
    m := &AgreementFile{
        AgreementFileProperties: *NewAgreementFileProperties(),
    }
    return m
}
// CreateAgreementFileFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAgreementFileFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreementFile(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AgreementFile) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AgreementFileProperties.GetFieldDeserializers()
    res["localizations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAgreementFileLocalizationFromDiscriminatorValue , m.SetLocalizations)
    return res
}
// GetLocalizations gets the localizations property value. The localized version of the terms of use agreement files attached to the agreement.
func (m *AgreementFile) GetLocalizations()([]AgreementFileLocalizationable) {
    return m.localizations
}
// Serialize serializes information the current object
func (m *AgreementFile) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AgreementFileProperties.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetLocalizations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetLocalizations())
        err = writer.WriteCollectionOfObjectValues("localizations", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLocalizations sets the localizations property value. The localized version of the terms of use agreement files attached to the agreement.
func (m *AgreementFile) SetLocalizations(value []AgreementFileLocalizationable)() {
    m.localizations = value
}
