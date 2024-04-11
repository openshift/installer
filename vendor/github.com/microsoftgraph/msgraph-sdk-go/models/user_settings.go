package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserSettings 
type UserSettings struct {
    Entity
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
    val, err := m.GetBackingStore().Get("contributionToContentDiscoveryAsOrganizationDisabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetContributionToContentDiscoveryDisabled gets the contributionToContentDiscoveryDisabled property value. The contributionToContentDiscoveryDisabled property
func (m *UserSettings) GetContributionToContentDiscoveryDisabled()(*bool) {
    val, err := m.GetBackingStore().Get("contributionToContentDiscoveryDisabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["contributionToContentDiscoveryAsOrganizationDisabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContributionToContentDiscoveryAsOrganizationDisabled(val)
        }
        return nil
    }
    res["contributionToContentDiscoveryDisabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContributionToContentDiscoveryDisabled(val)
        }
        return nil
    }
    res["shiftPreferences"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateShiftPreferencesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShiftPreferences(val.(ShiftPreferencesable))
        }
        return nil
    }
    return res
}
// GetShiftPreferences gets the shiftPreferences property value. The shiftPreferences property
func (m *UserSettings) GetShiftPreferences()(ShiftPreferencesable) {
    val, err := m.GetBackingStore().Get("shiftPreferences")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ShiftPreferencesable)
    }
    return nil
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
    err := m.GetBackingStore().Set("contributionToContentDiscoveryAsOrganizationDisabled", value)
    if err != nil {
        panic(err)
    }
}
// SetContributionToContentDiscoveryDisabled sets the contributionToContentDiscoveryDisabled property value. The contributionToContentDiscoveryDisabled property
func (m *UserSettings) SetContributionToContentDiscoveryDisabled(value *bool)() {
    err := m.GetBackingStore().Set("contributionToContentDiscoveryDisabled", value)
    if err != nil {
        panic(err)
    }
}
// SetShiftPreferences sets the shiftPreferences property value. The shiftPreferences property
func (m *UserSettings) SetShiftPreferences(value ShiftPreferencesable)() {
    err := m.GetBackingStore().Set("shiftPreferences", value)
    if err != nil {
        panic(err)
    }
}
// UserSettingsable 
type UserSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContributionToContentDiscoveryAsOrganizationDisabled()(*bool)
    GetContributionToContentDiscoveryDisabled()(*bool)
    GetShiftPreferences()(ShiftPreferencesable)
    SetContributionToContentDiscoveryAsOrganizationDisabled(value *bool)()
    SetContributionToContentDiscoveryDisabled(value *bool)()
    SetShiftPreferences(value ShiftPreferencesable)()
}
