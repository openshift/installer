package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AuditEventable 
type AuditEventable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivity()(*string)
    GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetActivityOperationType()(*string)
    GetActivityResult()(*string)
    GetActivityType()(*string)
    GetActor()(AuditActorable)
    GetCategory()(*string)
    GetComponentName()(*string)
    GetCorrelationId()(*string)
    GetDisplayName()(*string)
    GetResources()([]AuditResourceable)
    SetActivity(value *string)()
    SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetActivityOperationType(value *string)()
    SetActivityResult(value *string)()
    SetActivityType(value *string)()
    SetActor(value AuditActorable)()
    SetCategory(value *string)()
    SetComponentName(value *string)()
    SetCorrelationId(value *string)()
    SetDisplayName(value *string)()
    SetResources(value []AuditResourceable)()
}
