package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OnPremisesConditionalAccessSettingsable 
type OnPremisesConditionalAccessSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEnabled()(*bool)
    GetExcludedGroups()([]string)
    GetIncludedGroups()([]string)
    GetOverrideDefaultRule()(*bool)
    SetEnabled(value *bool)()
    SetExcludedGroups(value []string)()
    SetIncludedGroups(value []string)()
    SetOverrideDefaultRule(value *bool)()
}
