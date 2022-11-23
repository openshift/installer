package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserSettings 
type UserSettings struct {
    Entity
    // The contributionToContentDiscoveryAsOrganizationDisabled property
    contributionToContentDiscoveryAsOrganizationDisabled *bool
    // The contributionToContentDiscoveryDisabled property
    contributionToContentDiscoveryDisabled *bool
    // The shiftPreferences property
    shiftPreferences ShiftPreferencesable
}
// NewUserSettings instantiates a new userSettings and sets the default values.
func NewUserSettings()(*UserSettings) {
    m := &UserSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserSettings(), nil
}
// GetContributionToContentDiscoveryAsOrganizationDisabled gets the contributionToContentDiscoveryAsOrganizationDisabled property value. The contributionToContentDiscoveryAsOrganizationDisabled property
func (m *UserSettings) GetContributionToContentDiscoveryAsOrganizationDisabled()(*bool) {
    return m.contributionToContentDiscoveryAsOrganizationDisabled
}
// GetContributionToContentDiscoveryDisabled gets the contributionToContentDiscoveryDisabled property value. The contributionToContentDiscoveryDisabled property
func (m *UserSettings) GetContributionToContentDiscoveryDisabled()(*bool) {
    return m.contributionToContentDiscoveryDisabled
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["contributionToContentDiscoveryAsOrganizationDisabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetContributionToContentDiscoveryAsOrganizationDisabled)
    res["contributionToContentDiscoveryDisabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetContributionToContentDiscoveryDisabled)
    res["shiftPreferences"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateShiftPreferencesFromDiscriminatorValue , m.SetShiftPreferences)
    return res
}
// GetShiftPreferences gets the shiftPreferences property value. The shiftPreferences property
func (m *UserSettings) GetShiftPreferences()(ShiftPreferencesable) {
    return m.shiftPreferences
}
// Serialize serializes information the current object
func (m *UserSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("contributionToContentDiscoveryAsOrganizationDisabled", m.GetContributionToContentDiscoveryAsOrganizationDisabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("contributionToContentDiscoveryDisabled", m.GetContributionToContentDiscoveryDisabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("shiftPreferences", m.GetShiftPreferences())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContributionToContentDiscoveryAsOrganizationDisabled sets the contributionToContentDiscoveryAsOrganizationDisabled property value. The contributionToContentDiscoveryAsOrganizationDisabled property
func (m *UserSettings) SetContributionToContentDiscoveryAsOrganizationDisabled(value *bool)() {
    m.contributionToContentDiscoveryAsOrganizationDisabled = value
}
// SetContributionToContentDiscoveryDisabled sets the contributionToContentDiscoveryDisabled property value. The contributionToContentDiscoveryDisabled property
func (m *UserSettings) SetContributionToContentDiscoveryDisabled(value *bool)() {
    m.contributionToContentDiscoveryDisabled = value
}
// SetShiftPreferences sets the shiftPreferences property value. The shiftPreferences property
func (m *UserSettings) SetShiftPreferences(value ShiftPreferencesable)() {
    m.shiftPreferences = value
}
