package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosMinimumOperatingSystemable 
type IosMinimumOperatingSystemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOdataType()(*string)
    GetV10_0()(*bool)
    GetV11_0()(*bool)
    GetV12_0()(*bool)
    GetV13_0()(*bool)
    GetV14_0()(*bool)
    GetV8_0()(*bool)
    GetV9_0()(*bool)
    SetOdataType(value *string)()
    SetV10_0(value *bool)()
    SetV11_0(value *bool)()
    SetV12_0(value *bool)()
    SetV13_0(value *bool)()
    SetV14_0(value *bool)()
    SetV8_0(value *bool)()
    SetV9_0(value *bool)()
}
