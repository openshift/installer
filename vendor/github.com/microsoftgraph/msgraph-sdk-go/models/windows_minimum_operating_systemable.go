package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsMinimumOperatingSystemable 
type WindowsMinimumOperatingSystemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOdataType()(*string)
    GetV10_0()(*bool)
    GetV8_0()(*bool)
    GetV8_1()(*bool)
    SetOdataType(value *string)()
    SetV10_0(value *bool)()
    SetV8_0(value *bool)()
    SetV8_1(value *bool)()
}
