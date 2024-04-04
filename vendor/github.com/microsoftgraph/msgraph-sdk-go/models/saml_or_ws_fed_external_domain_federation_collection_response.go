package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SamlOrWsFedExternalDomainFederationCollectionResponse 
type SamlOrWsFedExternalDomainFederationCollectionResponse struct {
    BaseCollectionPaginationCountResponse
}
// NewSamlOrWsFedExternalDomainFederationCollectionResponse instantiates a new SamlOrWsFedExternalDomainFederationCollectionResponse and sets the default values.
func NewSamlOrWsFedExternalDomainFederationCollectionResponse()(*SamlOrWsFedExternalDomainFederationCollectionResponse) {
    m := &SamlOrWsFedExternalDomainFederationCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateSamlOrWsFedExternalDomainFederationCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSamlOrWsFedExternalDomainFederationCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSamlOrWsFedExternalDomainFederationCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SamlOrWsFedExternalDomainFederationCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSamlOrWsFedExternalDomainFederationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SamlOrWsFedExternalDomainFederationable, len(val))
            for i, v := range val {
                res[i] = v.(SamlOrWsFedExternalDomainFederationable)
            }
            m.SetValue(res)
        }
        return nil
    }
    return res
}
// GetValue gets the value property value. The value property
func (m *SamlOrWsFedExternalDomainFederationCollectionResponse) GetValue()([]SamlOrWsFedExternalDomainFederationable) {
    val, err := m.GetBackingStore().Get("value")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SamlOrWsFedExternalDomainFederationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SamlOrWsFedExternalDomainFederationCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseCollectionPaginationCountResponse.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetValue() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetValue()))
        for i, v := range m.GetValue() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("value", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. The value property
func (m *SamlOrWsFedExternalDomainFederationCollectionResponse) SetValue(value []SamlOrWsFedExternalDomainFederationable)() {
    err := m.GetBackingStore().Set("value", value)
    if err != nil {
        panic(err)
    }
}
// SamlOrWsFedExternalDomainFederationCollectionResponseable 
type SamlOrWsFedExternalDomainFederationCollectionResponseable interface {
    BaseCollectionPaginationCountResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetValue()([]SamlOrWsFedExternalDomainFederationable)
    SetValue(value []SamlOrWsFedExternalDomainFederationable)()
}
