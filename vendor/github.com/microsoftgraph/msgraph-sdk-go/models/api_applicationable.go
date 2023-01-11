package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ApiApplicationable 
type ApiApplicationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAcceptMappedClaims()(*bool)
    GetKnownClientApplications()([]string)
    GetOauth2PermissionScopes()([]PermissionScopeable)
    GetOdataType()(*string)
    GetPreAuthorizedApplications()([]PreAuthorizedApplicationable)
    GetRequestedAccessTokenVersion()(*int32)
    SetAcceptMappedClaims(value *bool)()
    SetKnownClientApplications(value []string)()
    SetOauth2PermissionScopes(value []PermissionScopeable)()
    SetOdataType(value *string)()
    SetPreAuthorizedApplications(value []PreAuthorizedApplicationable)()
    SetRequestedAccessTokenVersion(value *int32)()
}
