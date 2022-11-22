package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SubscribedSkuable 
type SubscribedSkuable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppliesTo()(*string)
    GetCapabilityStatus()(*string)
    GetConsumedUnits()(*int32)
    GetPrepaidUnits()(LicenseUnitsDetailable)
    GetServicePlans()([]ServicePlanInfoable)
    GetSkuId()(*string)
    GetSkuPartNumber()(*string)
    SetAppliesTo(value *string)()
    SetCapabilityStatus(value *string)()
    SetConsumedUnits(value *int32)()
    SetPrepaidUnits(value LicenseUnitsDetailable)()
    SetServicePlans(value []ServicePlanInfoable)()
    SetSkuId(value *string)()
    SetSkuPartNumber(value *string)()
}
