package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OutlookUser 
type OutlookUser struct {
    Entity
    // A list of categories defined for the user.
    masterCategories []OutlookCategoryable
}
// NewOutlookUser instantiates a new outlookUser and sets the default values.
func NewOutlookUser()(*OutlookUser) {
    m := &OutlookUser{
        Entity: *NewEntity(),
    }
    return m
}
// CreateOutlookUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOutlookUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOutlookUser(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *OutlookUser) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["masterCategories"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOutlookCategoryFromDiscriminatorValue , m.SetMasterCategories)
    return res
}
// GetMasterCategories gets the masterCategories property value. A list of categories defined for the user.
func (m *OutlookUser) GetMasterCategories()([]OutlookCategoryable) {
    return m.masterCategories
}
// Serialize serializes information the current object
func (m *OutlookUser) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetMasterCategories() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMasterCategories())
        err = writer.WriteCollectionOfObjectValues("masterCategories", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMasterCategories sets the masterCategories property value. A list of categories defined for the user.
func (m *OutlookUser) SetMasterCategories(value []OutlookCategoryable)() {
    m.masterCategories = value
}
