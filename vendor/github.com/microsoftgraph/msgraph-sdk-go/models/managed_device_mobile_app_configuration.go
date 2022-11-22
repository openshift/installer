package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedDeviceMobileAppConfiguration an abstract class for Mobile app configuration for enrolled devices.
type ManagedDeviceMobileAppConfiguration struct {
    Entity
    // The list of group assignemenets for app configration.
    assignments []ManagedDeviceMobileAppConfigurationAssignmentable
    // DateTime the object was created.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Admin provided description of the Device Configuration.
    description *string
    // List of ManagedDeviceMobileAppConfigurationDeviceStatus.
    deviceStatuses []ManagedDeviceMobileAppConfigurationDeviceStatusable
    // App configuration device status summary.
    deviceStatusSummary ManagedDeviceMobileAppConfigurationDeviceSummaryable
    // Admin provided name of the device configuration.
    displayName *string
    // DateTime the object was last modified.
    lastModifiedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // the associated app.
    targetedMobileApps []string
    // List of ManagedDeviceMobileAppConfigurationUserStatus.
    userStatuses []ManagedDeviceMobileAppConfigurationUserStatusable
    // App configuration user status summary.
    userStatusSummary ManagedDeviceMobileAppConfigurationUserSummaryable
    // Version of the device configuration.
    version *int32
}
// NewManagedDeviceMobileAppConfiguration instantiates a new managedDeviceMobileAppConfiguration and sets the default values.
func NewManagedDeviceMobileAppConfiguration()(*ManagedDeviceMobileAppConfiguration) {
    m := &ManagedDeviceMobileAppConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedDeviceMobileAppConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateManagedDeviceMobileAppConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.iosMobileAppConfiguration":
                        return NewIosMobileAppConfiguration(), nil
                }
            }
        }
    }
    return NewManagedDeviceMobileAppConfiguration(), nil
}
// GetAssignments gets the assignments property value. The list of group assignemenets for app configration.
func (m *ManagedDeviceMobileAppConfiguration) GetAssignments()([]ManagedDeviceMobileAppConfigurationAssignmentable) {
    return m.assignments
}
// GetCreatedDateTime gets the createdDateTime property value. DateTime the object was created.
func (m *ManagedDeviceMobileAppConfiguration) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDescription gets the description property value. Admin provided description of the Device Configuration.
func (m *ManagedDeviceMobileAppConfiguration) GetDescription()(*string) {
    return m.description
}
// GetDeviceStatuses gets the deviceStatuses property value. List of ManagedDeviceMobileAppConfigurationDeviceStatus.
func (m *ManagedDeviceMobileAppConfiguration) GetDeviceStatuses()([]ManagedDeviceMobileAppConfigurationDeviceStatusable) {
    return m.deviceStatuses
}
// GetDeviceStatusSummary gets the deviceStatusSummary property value. App configuration device status summary.
func (m *ManagedDeviceMobileAppConfiguration) GetDeviceStatusSummary()(ManagedDeviceMobileAppConfigurationDeviceSummaryable) {
    return m.deviceStatusSummary
}
// GetDisplayName gets the displayName property value. Admin provided name of the device configuration.
func (m *ManagedDeviceMobileAppConfiguration) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ManagedDeviceMobileAppConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateManagedDeviceMobileAppConfigurationAssignmentFromDiscriminatorValue , m.SetAssignments)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["deviceStatuses"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateManagedDeviceMobileAppConfigurationDeviceStatusFromDiscriminatorValue , m.SetDeviceStatuses)
    res["deviceStatusSummary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateManagedDeviceMobileAppConfigurationDeviceSummaryFromDiscriminatorValue , m.SetDeviceStatusSummary)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["lastModifiedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastModifiedDateTime)
    res["targetedMobileApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetTargetedMobileApps)
    res["userStatuses"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateManagedDeviceMobileAppConfigurationUserStatusFromDiscriminatorValue , m.SetUserStatuses)
    res["userStatusSummary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateManagedDeviceMobileAppConfigurationUserSummaryFromDiscriminatorValue , m.SetUserStatusSummary)
    res["version"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetVersion)
    return res
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. DateTime the object was last modified.
func (m *ManagedDeviceMobileAppConfiguration) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastModifiedDateTime
}
// GetTargetedMobileApps gets the targetedMobileApps property value. the associated app.
func (m *ManagedDeviceMobileAppConfiguration) GetTargetedMobileApps()([]string) {
    return m.targetedMobileApps
}
// GetUserStatuses gets the userStatuses property value. List of ManagedDeviceMobileAppConfigurationUserStatus.
func (m *ManagedDeviceMobileAppConfiguration) GetUserStatuses()([]ManagedDeviceMobileAppConfigurationUserStatusable) {
    return m.userStatuses
}
// GetUserStatusSummary gets the userStatusSummary property value. App configuration user status summary.
func (m *ManagedDeviceMobileAppConfiguration) GetUserStatusSummary()(ManagedDeviceMobileAppConfigurationUserSummaryable) {
    return m.userStatusSummary
}
// GetVersion gets the version property value. Version of the device configuration.
func (m *ManagedDeviceMobileAppConfiguration) GetVersion()(*int32) {
    return m.version
}
// Serialize serializes information the current object
func (m *ManagedDeviceMobileAppConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignments() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAssignments())
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceStatuses() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceStatuses())
        err = writer.WriteCollectionOfObjectValues("deviceStatuses", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deviceStatusSummary", m.GetDeviceStatusSummary())
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
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetTargetedMobileApps() != nil {
        err = writer.WriteCollectionOfStringValues("targetedMobileApps", m.GetTargetedMobileApps())
        if err != nil {
            return err
        }
    }
    if m.GetUserStatuses() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetUserStatuses())
        err = writer.WriteCollectionOfObjectValues("userStatuses", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("userStatusSummary", m.GetUserStatusSummary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignments sets the assignments property value. The list of group assignemenets for app configration.
func (m *ManagedDeviceMobileAppConfiguration) SetAssignments(value []ManagedDeviceMobileAppConfigurationAssignmentable)() {
    m.assignments = value
}
// SetCreatedDateTime sets the createdDateTime property value. DateTime the object was created.
func (m *ManagedDeviceMobileAppConfiguration) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDescription sets the description property value. Admin provided description of the Device Configuration.
func (m *ManagedDeviceMobileAppConfiguration) SetDescription(value *string)() {
    m.description = value
}
// SetDeviceStatuses sets the deviceStatuses property value. List of ManagedDeviceMobileAppConfigurationDeviceStatus.
func (m *ManagedDeviceMobileAppConfiguration) SetDeviceStatuses(value []ManagedDeviceMobileAppConfigurationDeviceStatusable)() {
    m.deviceStatuses = value
}
// SetDeviceStatusSummary sets the deviceStatusSummary property value. App configuration device status summary.
func (m *ManagedDeviceMobileAppConfiguration) SetDeviceStatusSummary(value ManagedDeviceMobileAppConfigurationDeviceSummaryable)() {
    m.deviceStatusSummary = value
}
// SetDisplayName sets the displayName property value. Admin provided name of the device configuration.
func (m *ManagedDeviceMobileAppConfiguration) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. DateTime the object was last modified.
func (m *ManagedDeviceMobileAppConfiguration) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastModifiedDateTime = value
}
// SetTargetedMobileApps sets the targetedMobileApps property value. the associated app.
func (m *ManagedDeviceMobileAppConfiguration) SetTargetedMobileApps(value []string)() {
    m.targetedMobileApps = value
}
// SetUserStatuses sets the userStatuses property value. List of ManagedDeviceMobileAppConfigurationUserStatus.
func (m *ManagedDeviceMobileAppConfiguration) SetUserStatuses(value []ManagedDeviceMobileAppConfigurationUserStatusable)() {
    m.userStatuses = value
}
// SetUserStatusSummary sets the userStatusSummary property value. App configuration user status summary.
func (m *ManagedDeviceMobileAppConfiguration) SetUserStatusSummary(value ManagedDeviceMobileAppConfigurationUserSummaryable)() {
    m.userStatusSummary = value
}
// SetVersion sets the version property value. Version of the device configuration.
func (m *ManagedDeviceMobileAppConfiguration) SetVersion(value *int32)() {
    m.version = value
}
