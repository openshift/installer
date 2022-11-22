package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ServicePlanInfoable 
type ServicePlanInfoable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppliesTo()(*string)
    GetOdataType()(*string)
    GetProvisioningStatus()(*string)
    GetServicePlanId()(*string)
    GetServicePlanName()(*string)
    SetAppliesTo(value *string)()
    SetOdataType(value *string)()
    SetProvisioningStatus(value *string)()
    SetServicePlanId(value *string)()
    SetServicePlanName(value *string)()
}
