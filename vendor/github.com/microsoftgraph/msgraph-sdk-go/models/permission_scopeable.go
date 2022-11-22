package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PermissionScopeable 
type PermissionScopeable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAdminConsentDescription()(*string)
    GetAdminConsentDisplayName()(*string)
    GetId()(*string)
    GetIsEnabled()(*bool)
    GetOdataType()(*string)
    GetOrigin()(*string)
    GetType()(*string)
    GetUserConsentDescription()(*string)
    GetUserConsentDisplayName()(*string)
    GetValue()(*string)
    SetAdminConsentDescription(value *string)()
    SetAdminConsentDisplayName(value *string)()
    SetId(value *string)()
    SetIsEnabled(value *bool)()
    SetOdataType(value *string)()
    SetOrigin(value *string)()
    SetType(value *string)()
    SetUserConsentDescription(value *string)()
    SetUserConsentDisplayName(value *string)()
    SetValue(value *string)()
}
