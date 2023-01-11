package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CertificateBasedAuthConfiguration provides operations to manage the collection of certificateBasedAuthConfiguration entities.
type CertificateBasedAuthConfiguration struct {
    Entity
    // Collection of certificate authorities which creates a trusted certificate chain.
    certificateAuthorities []CertificateAuthorityable
}
// NewCertificateBasedAuthConfiguration instantiates a new certificateBasedAuthConfiguration and sets the default values.
func NewCertificateBasedAuthConfiguration()(*CertificateBasedAuthConfiguration) {
    m := &CertificateBasedAuthConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCertificateBasedAuthConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCertificateBasedAuthConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCertificateBasedAuthConfiguration(), nil
}
// GetCertificateAuthorities gets the certificateAuthorities property value. Collection of certificate authorities which creates a trusted certificate chain.
func (m *CertificateBasedAuthConfiguration) GetCertificateAuthorities()([]CertificateAuthorityable) {
    return m.certificateAuthorities
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CertificateBasedAuthConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["certificateAuthorities"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCertificateAuthorityFromDiscriminatorValue , m.SetCertificateAuthorities)
    return res
}
// Serialize serializes information the current object
func (m *CertificateBasedAuthConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCertificateAuthorities() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCertificateAuthorities())
        err = writer.WriteCollectionOfObjectValues("certificateAuthorities", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCertificateAuthorities sets the certificateAuthorities property value. Collection of certificate authorities which creates a trusted certificate chain.
func (m *CertificateBasedAuthConfiguration) SetCertificateAuthorities(value []CertificateAuthorityable)() {
    m.certificateAuthorities = value
}
