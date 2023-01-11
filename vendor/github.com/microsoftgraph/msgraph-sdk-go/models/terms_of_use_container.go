package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TermsOfUseContainer 
type TermsOfUseContainer struct {
    Entity
    // Represents the current status of a user's response to a company's customizable terms of use agreement.
    agreementAcceptances []AgreementAcceptanceable
    // Represents a tenant's customizable terms of use agreement that's created and managed with Azure Active Directory (Azure AD).
    agreements []Agreementable
}
// NewTermsOfUseContainer instantiates a new TermsOfUseContainer and sets the default values.
func NewTermsOfUseContainer()(*TermsOfUseContainer) {
    m := &TermsOfUseContainer{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTermsOfUseContainerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTermsOfUseContainerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTermsOfUseContainer(), nil
}
// GetAgreementAcceptances gets the agreementAcceptances property value. Represents the current status of a user's response to a company's customizable terms of use agreement.
func (m *TermsOfUseContainer) GetAgreementAcceptances()([]AgreementAcceptanceable) {
    return m.agreementAcceptances
}
// GetAgreements gets the agreements property value. Represents a tenant's customizable terms of use agreement that's created and managed with Azure Active Directory (Azure AD).
func (m *TermsOfUseContainer) GetAgreements()([]Agreementable) {
    return m.agreements
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TermsOfUseContainer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["agreementAcceptances"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAgreementAcceptanceFromDiscriminatorValue , m.SetAgreementAcceptances)
    res["agreements"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAgreementFromDiscriminatorValue , m.SetAgreements)
    return res
}
// Serialize serializes information the current object
func (m *TermsOfUseContainer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAgreementAcceptances() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAgreementAcceptances())
        err = writer.WriteCollectionOfObjectValues("agreementAcceptances", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAgreements() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAgreements())
        err = writer.WriteCollectionOfObjectValues("agreements", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAgreementAcceptances sets the agreementAcceptances property value. Represents the current status of a user's response to a company's customizable terms of use agreement.
func (m *TermsOfUseContainer) SetAgreementAcceptances(value []AgreementAcceptanceable)() {
    m.agreementAcceptances = value
}
// SetAgreements sets the agreements property value. Represents a tenant's customizable terms of use agreement that's created and managed with Azure Active Directory (Azure AD).
func (m *TermsOfUseContainer) SetAgreements(value []Agreementable)() {
    m.agreements = value
}
