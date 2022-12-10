package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppCatalogs provides operations to manage the appCatalogs singleton.
type AppCatalogs struct {
    Entity
    // The teamsApps property
    teamsApps []TeamsAppable
}
// NewAppCatalogs instantiates a new appCatalogs and sets the default values.
func NewAppCatalogs()(*AppCatalogs) {
    m := &AppCatalogs{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAppCatalogsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAppCatalogsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppCatalogs(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AppCatalogs) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["teamsApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamsAppFromDiscriminatorValue , m.SetTeamsApps)
    return res
}
// GetTeamsApps gets the teamsApps property value. The teamsApps property
func (m *AppCatalogs) GetTeamsApps()([]TeamsAppable) {
    return m.teamsApps
}
// Serialize serializes information the current object
func (m *AppCatalogs) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetTeamsApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTeamsApps())
        err = writer.WriteCollectionOfObjectValues("teamsApps", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetTeamsApps sets the teamsApps property value. The teamsApps property
func (m *AppCatalogs) SetTeamsApps(value []TeamsAppable)() {
    m.teamsApps = value
}
