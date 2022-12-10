package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PermissionGrantPolicy 
type PermissionGrantPolicy struct {
    PolicyBase
    // Condition sets which are excluded in this permission grant policy. Automatically expanded on GET.
    excludes []PermissionGrantConditionSetable
    // Condition sets which are included in this permission grant policy. Automatically expanded on GET.
    includes []PermissionGrantConditionSetable
}
// NewPermissionGrantPolicy instantiates a new PermissionGrantPolicy and sets the default values.
func NewPermissionGrantPolicy()(*PermissionGrantPolicy) {
    m := &PermissionGrantPolicy{
        PolicyBase: *NewPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.permissionGrantPolicy";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreatePermissionGrantPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePermissionGrantPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPermissionGrantPolicy(), nil
}
// GetExcludes gets the excludes property value. Condition sets which are excluded in this permission grant policy. Automatically expanded on GET.
func (m *PermissionGrantPolicy) GetExcludes()([]PermissionGrantConditionSetable) {
    return m.excludes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PermissionGrantPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PolicyBase.GetFieldDeserializers()
    res["excludes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePermissionGrantConditionSetFromDiscriminatorValue , m.SetExcludes)
    res["includes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePermissionGrantConditionSetFromDiscriminatorValue , m.SetIncludes)
    return res
}
// GetIncludes gets the includes property value. Condition sets which are included in this permission grant policy. Automatically expanded on GET.
func (m *PermissionGrantPolicy) GetIncludes()([]PermissionGrantConditionSetable) {
    return m.includes
}
// Serialize serializes information the current object
func (m *PermissionGrantPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PolicyBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetExcludes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetExcludes())
        err = writer.WriteCollectionOfObjectValues("excludes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetIncludes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetIncludes())
        err = writer.WriteCollectionOfObjectValues("includes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExcludes sets the excludes property value. Condition sets which are excluded in this permission grant policy. Automatically expanded on GET.
func (m *PermissionGrantPolicy) SetExcludes(value []PermissionGrantConditionSetable)() {
    m.excludes = value
}
// SetIncludes sets the includes property value. Condition sets which are included in this permission grant policy. Automatically expanded on GET.
func (m *PermissionGrantPolicy) SetIncludes(value []PermissionGrantConditionSetable)() {
    m.includes = value
}
