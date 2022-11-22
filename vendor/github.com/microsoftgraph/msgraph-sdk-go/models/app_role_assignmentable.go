package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppRoleAssignmentable 
type AppRoleAssignmentable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppRoleId()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPrincipalDisplayName()(*string)
    GetPrincipalId()(*string)
    GetPrincipalType()(*string)
    GetResourceDisplayName()(*string)
    GetResourceId()(*string)
    SetAppRoleId(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPrincipalDisplayName(value *string)()
    SetPrincipalId(value *string)()
    SetPrincipalType(value *string)()
    SetResourceDisplayName(value *string)()
    SetResourceId(value *string)()
}
