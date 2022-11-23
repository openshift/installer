package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AddInable 
type AddInable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetId()(*string)
    GetOdataType()(*string)
    GetProperties()([]KeyValueable)
    GetType()(*string)
    SetId(value *string)()
    SetOdataType(value *string)()
    SetProperties(value []KeyValueable)()
    SetType(value *string)()
}
