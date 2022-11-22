package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SamlOrWsFedExternalDomainFederation 
type SamlOrWsFedExternalDomainFederation struct {
    SamlOrWsFedProvider
    // Collection of domain names of the external organizations that the tenant is federating with. Supports $filter (eq).
    domains []ExternalDomainNameable
}
// NewSamlOrWsFedExternalDomainFederation instantiates a new SamlOrWsFedExternalDomainFederation and sets the default values.
func NewSamlOrWsFedExternalDomainFederation()(*SamlOrWsFedExternalDomainFederation) {
    m := &SamlOrWsFedExternalDomainFederation{
        SamlOrWsFedProvider: *NewSamlOrWsFedProvider(),
    }
    odataTypeValue := "#microsoft.graph.samlOrWsFedExternalDomainFederation";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateSamlOrWsFedExternalDomainFederationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSamlOrWsFedExternalDomainFederationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSamlOrWsFedExternalDomainFederation(), nil
}
// GetDomains gets the domains property value. Collection of domain names of the external organizations that the tenant is federating with. Supports $filter (eq).
func (m *SamlOrWsFedExternalDomainFederation) GetDomains()([]ExternalDomainNameable) {
    return m.domains
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SamlOrWsFedExternalDomainFederation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SamlOrWsFedProvider.GetFieldDeserializers()
    res["domains"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateExternalDomainNameFromDiscriminatorValue , m.SetDomains)
    return res
}
// Serialize serializes information the current object
func (m *SamlOrWsFedExternalDomainFederation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SamlOrWsFedProvider.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDomains() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDomains())
        err = writer.WriteCollectionOfObjectValues("domains", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDomains sets the domains property value. Collection of domain names of the external organizations that the tenant is federating with. Supports $filter (eq).
func (m *SamlOrWsFedExternalDomainFederation) SetDomains(value []ExternalDomainNameable)() {
    m.domains = value
}
