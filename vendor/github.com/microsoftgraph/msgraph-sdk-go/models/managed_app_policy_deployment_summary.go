package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAppPolicyDeploymentSummary 
type ManagedAppPolicyDeploymentSummary struct {
    Entity
    // Not yet documented
    configurationDeployedUserCount *int32
    // Not yet documented
    configurationDeploymentSummaryPerApp []ManagedAppPolicyDeploymentSummaryPerAppable
    // Not yet documented
    displayName *string
    // Not yet documented
    lastRefreshTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Version of the entity.
    version *string
}
// NewManagedAppPolicyDeploymentSummary instantiates a new managedAppPolicyDeploymentSummary and sets the default values.
func NewManagedAppPolicyDeploymentSummary()(*ManagedAppPolicyDeploymentSummary) {
    m := &ManagedAppPolicyDeploymentSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedAppPolicyDeploymentSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateManagedAppPolicyDeploymentSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedAppPolicyDeploymentSummary(), nil
}
// GetConfigurationDeployedUserCount gets the configurationDeployedUserCount property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) GetConfigurationDeployedUserCount()(*int32) {
    return m.configurationDeployedUserCount
}
// GetConfigurationDeploymentSummaryPerApp gets the configurationDeploymentSummaryPerApp property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) GetConfigurationDeploymentSummaryPerApp()([]ManagedAppPolicyDeploymentSummaryPerAppable) {
    return m.configurationDeploymentSummaryPerApp
}
// GetDisplayName gets the displayName property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ManagedAppPolicyDeploymentSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["configurationDeployedUserCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetConfigurationDeployedUserCount)
    res["configurationDeploymentSummaryPerApp"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateManagedAppPolicyDeploymentSummaryPerAppFromDiscriminatorValue , m.SetConfigurationDeploymentSummaryPerApp)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["lastRefreshTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastRefreshTime)
    res["version"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVersion)
    return res
}
// GetLastRefreshTime gets the lastRefreshTime property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) GetLastRefreshTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastRefreshTime
}
// GetVersion gets the version property value. Version of the entity.
func (m *ManagedAppPolicyDeploymentSummary) GetVersion()(*string) {
    return m.version
}
// Serialize serializes information the current object
func (m *ManagedAppPolicyDeploymentSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("configurationDeployedUserCount", m.GetConfigurationDeployedUserCount())
        if err != nil {
            return err
        }
    }
    if m.GetConfigurationDeploymentSummaryPerApp() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetConfigurationDeploymentSummaryPerApp())
        err = writer.WriteCollectionOfObjectValues("configurationDeploymentSummaryPerApp", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastRefreshTime", m.GetLastRefreshTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConfigurationDeployedUserCount sets the configurationDeployedUserCount property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetConfigurationDeployedUserCount(value *int32)() {
    m.configurationDeployedUserCount = value
}
// SetConfigurationDeploymentSummaryPerApp sets the configurationDeploymentSummaryPerApp property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetConfigurationDeploymentSummaryPerApp(value []ManagedAppPolicyDeploymentSummaryPerAppable)() {
    m.configurationDeploymentSummaryPerApp = value
}
// SetDisplayName sets the displayName property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetLastRefreshTime sets the lastRefreshTime property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetLastRefreshTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastRefreshTime = value
}
// SetVersion sets the version property value. Version of the entity.
func (m *ManagedAppPolicyDeploymentSummary) SetVersion(value *string)() {
    m.version = value
}
