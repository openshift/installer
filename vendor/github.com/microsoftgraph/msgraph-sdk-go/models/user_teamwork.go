package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserTeamwork 
type UserTeamwork struct {
    Entity
    // The list of associatedTeamInfo objects that a user is associated with.
    associatedTeams []AssociatedTeamInfoable
    // The apps installed in the personal scope of this user.
    installedApps []UserScopeTeamsAppInstallationable
}
// NewUserTeamwork instantiates a new userTeamwork and sets the default values.
func NewUserTeamwork()(*UserTeamwork) {
    m := &UserTeamwork{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserTeamworkFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserTeamworkFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserTeamwork(), nil
}
// GetAssociatedTeams gets the associatedTeams property value. The list of associatedTeamInfo objects that a user is associated with.
func (m *UserTeamwork) GetAssociatedTeams()([]AssociatedTeamInfoable) {
    return m.associatedTeams
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserTeamwork) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["associatedTeams"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAssociatedTeamInfoFromDiscriminatorValue , m.SetAssociatedTeams)
    res["installedApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateUserScopeTeamsAppInstallationFromDiscriminatorValue , m.SetInstalledApps)
    return res
}
// GetInstalledApps gets the installedApps property value. The apps installed in the personal scope of this user.
func (m *UserTeamwork) GetInstalledApps()([]UserScopeTeamsAppInstallationable) {
    return m.installedApps
}
// Serialize serializes information the current object
func (m *UserTeamwork) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssociatedTeams() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAssociatedTeams())
        err = writer.WriteCollectionOfObjectValues("associatedTeams", cast)
        if err != nil {
            return err
        }
    }
    if m.GetInstalledApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetInstalledApps())
        err = writer.WriteCollectionOfObjectValues("installedApps", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssociatedTeams sets the associatedTeams property value. The list of associatedTeamInfo objects that a user is associated with.
func (m *UserTeamwork) SetAssociatedTeams(value []AssociatedTeamInfoable)() {
    m.associatedTeams = value
}
// SetInstalledApps sets the installedApps property value. The apps installed in the personal scope of this user.
func (m *UserTeamwork) SetInstalledApps(value []UserScopeTeamsAppInstallationable)() {
    m.installedApps = value
}
