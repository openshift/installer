package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// LicenseDetails provides operations to manage the collection of agreement entities.
type LicenseDetails struct {
    Entity
    // Information about the service plans assigned with the license. Read-only, Not nullable
    servicePlans []ServicePlanInfoable
    // Unique identifier (GUID) for the service SKU. Equal to the skuId property on the related SubscribedSku object. Read-only
    skuId *string
    // Unique SKU display name. Equal to the skuPartNumber on the related SubscribedSku object; for example: 'AAD_Premium'. Read-only
    skuPartNumber *string
}
// NewLicenseDetails instantiates a new licenseDetails and sets the default values.
func NewLicenseDetails()(*LicenseDetails) {
    m := &LicenseDetails{
        Entity: *NewEntity(),
    }
    return m
}
// CreateLicenseDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateLicenseDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLicenseDetails(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *LicenseDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["servicePlans"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateServicePlanInfoFromDiscriminatorValue , m.SetServicePlans)
    res["skuId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSkuId)
    res["skuPartNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSkuPartNumber)
    return res
}
// GetServicePlans gets the servicePlans property value. Information about the service plans assigned with the license. Read-only, Not nullable
func (m *LicenseDetails) GetServicePlans()([]ServicePlanInfoable) {
    return m.servicePlans
}
// GetSkuId gets the skuId property value. Unique identifier (GUID) for the service SKU. Equal to the skuId property on the related SubscribedSku object. Read-only
func (m *LicenseDetails) GetSkuId()(*string) {
    return m.skuId
}
// GetSkuPartNumber gets the skuPartNumber property value. Unique SKU display name. Equal to the skuPartNumber on the related SubscribedSku object; for example: 'AAD_Premium'. Read-only
func (m *LicenseDetails) GetSkuPartNumber()(*string) {
    return m.skuPartNumber
}
// Serialize serializes information the current object
func (m *LicenseDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetServicePlans() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetServicePlans())
        err = writer.WriteCollectionOfObjectValues("servicePlans", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("skuId", m.GetSkuId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("skuPartNumber", m.GetSkuPartNumber())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetServicePlans sets the servicePlans property value. Information about the service plans assigned with the license. Read-only, Not nullable
func (m *LicenseDetails) SetServicePlans(value []ServicePlanInfoable)() {
    m.servicePlans = value
}
// SetSkuId sets the skuId property value. Unique identifier (GUID) for the service SKU. Equal to the skuId property on the related SubscribedSku object. Read-only
func (m *LicenseDetails) SetSkuId(value *string)() {
    m.skuId = value
}
// SetSkuPartNumber sets the skuPartNumber property value. Unique SKU display name. Equal to the skuPartNumber on the related SubscribedSku object; for example: 'AAD_Premium'. Read-only
func (m *LicenseDetails) SetSkuPartNumber(value *string)() {
    m.skuPartNumber = value
}
