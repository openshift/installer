package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Contractable 
type Contractable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContractType()(*string)
    GetCustomerId()(*string)
    GetDefaultDomainName()(*string)
    GetDisplayName()(*string)
    SetContractType(value *string)()
    SetCustomerId(value *string)()
    SetDefaultDomainName(value *string)()
    SetDisplayName(value *string)()
}
