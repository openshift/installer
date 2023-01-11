package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppRoleable 
type AppRoleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedMemberTypes()([]string)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetId()(*string)
    GetIsEnabled()(*bool)
    GetOdataType()(*string)
    GetOrigin()(*string)
    GetValue()(*string)
    SetAllowedMemberTypes(value []string)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetId(value *string)()
    SetIsEnabled(value *bool)()
    SetOdataType(value *string)()
    SetOrigin(value *string)()
    SetValue(value *string)()
}
