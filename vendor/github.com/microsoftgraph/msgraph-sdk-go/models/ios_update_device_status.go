package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosUpdateDeviceStatus provides operations to manage the collection of agreement entities.
type IosUpdateDeviceStatus struct {
    Entity
    // The DateTime when device compliance grace period expires
    complianceGracePeriodExpirationDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Device name of the DevicePolicyStatus.
    deviceDisplayName *string
    // The device id that is being reported.
    deviceId *string
    // The device model that is being reported
    deviceModel *string
    // The installStatus property
    installStatus *IosUpdatesInstallStatus
    // Last modified date time of the policy report.
    lastReportedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The device version that is being reported.
    osVersion *string
    // The status property
    status *ComplianceStatus
    // The User id that is being reported.
    userId *string
    // The User Name that is being reported
    userName *string
    // UserPrincipalName.
    userPrincipalName *string
}
// NewIosUpdateDeviceStatus instantiates a new iosUpdateDeviceStatus and sets the default values.
func NewIosUpdateDeviceStatus()(*IosUpdateDeviceStatus) {
    m := &IosUpdateDeviceStatus{
        Entity: *NewEntity(),
    }
    return m
}
// CreateIosUpdateDeviceStatusFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosUpdateDeviceStatusFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosUpdateDeviceStatus(), nil
}
// GetComplianceGracePeriodExpirationDateTime gets the complianceGracePeriodExpirationDateTime property value. The DateTime when device compliance grace period expires
func (m *IosUpdateDeviceStatus) GetComplianceGracePeriodExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.complianceGracePeriodExpirationDateTime
}
// GetDeviceDisplayName gets the deviceDisplayName property value. Device name of the DevicePolicyStatus.
func (m *IosUpdateDeviceStatus) GetDeviceDisplayName()(*string) {
    return m.deviceDisplayName
}
// GetDeviceId gets the deviceId property value. The device id that is being reported.
func (m *IosUpdateDeviceStatus) GetDeviceId()(*string) {
    return m.deviceId
}
// GetDeviceModel gets the deviceModel property value. The device model that is being reported
func (m *IosUpdateDeviceStatus) GetDeviceModel()(*string) {
    return m.deviceModel
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosUpdateDeviceStatus) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["complianceGracePeriodExpirationDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetComplianceGracePeriodExpirationDateTime)
    res["deviceDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceDisplayName)
    res["deviceId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceId)
    res["deviceModel"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDeviceModel)
    res["installStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseIosUpdatesInstallStatus , m.SetInstallStatus)
    res["lastReportedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastReportedDateTime)
    res["osVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOsVersion)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseComplianceStatus , m.SetStatus)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    res["userName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserName)
    res["userPrincipalName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserPrincipalName)
    return res
}
// GetInstallStatus gets the installStatus property value. The installStatus property
func (m *IosUpdateDeviceStatus) GetInstallStatus()(*IosUpdatesInstallStatus) {
    return m.installStatus
}
// GetLastReportedDateTime gets the lastReportedDateTime property value. Last modified date time of the policy report.
func (m *IosUpdateDeviceStatus) GetLastReportedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastReportedDateTime
}
// GetOsVersion gets the osVersion property value. The device version that is being reported.
func (m *IosUpdateDeviceStatus) GetOsVersion()(*string) {
    return m.osVersion
}
// GetStatus gets the status property value. The status property
func (m *IosUpdateDeviceStatus) GetStatus()(*ComplianceStatus) {
    return m.status
}
// GetUserId gets the userId property value. The User id that is being reported.
func (m *IosUpdateDeviceStatus) GetUserId()(*string) {
    return m.userId
}
// GetUserName gets the userName property value. The User Name that is being reported
func (m *IosUpdateDeviceStatus) GetUserName()(*string) {
    return m.userName
}
// GetUserPrincipalName gets the userPrincipalName property value. UserPrincipalName.
func (m *IosUpdateDeviceStatus) GetUserPrincipalName()(*string) {
    return m.userPrincipalName
}
// Serialize serializes information the current object
func (m *IosUpdateDeviceStatus) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("complianceGracePeriodExpirationDateTime", m.GetComplianceGracePeriodExpirationDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceDisplayName", m.GetDeviceDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceId", m.GetDeviceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceModel", m.GetDeviceModel())
        if err != nil {
            return err
        }
    }
    if m.GetInstallStatus() != nil {
        cast := (*m.GetInstallStatus()).String()
        err = writer.WriteStringValue("installStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastReportedDateTime", m.GetLastReportedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("osVersion", m.GetOsVersion())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userName", m.GetUserName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetComplianceGracePeriodExpirationDateTime sets the complianceGracePeriodExpirationDateTime property value. The DateTime when device compliance grace period expires
func (m *IosUpdateDeviceStatus) SetComplianceGracePeriodExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.complianceGracePeriodExpirationDateTime = value
}
// SetDeviceDisplayName sets the deviceDisplayName property value. Device name of the DevicePolicyStatus.
func (m *IosUpdateDeviceStatus) SetDeviceDisplayName(value *string)() {
    m.deviceDisplayName = value
}
// SetDeviceId sets the deviceId property value. The device id that is being reported.
func (m *IosUpdateDeviceStatus) SetDeviceId(value *string)() {
    m.deviceId = value
}
// SetDeviceModel sets the deviceModel property value. The device model that is being reported
func (m *IosUpdateDeviceStatus) SetDeviceModel(value *string)() {
    m.deviceModel = value
}
// SetInstallStatus sets the installStatus property value. The installStatus property
func (m *IosUpdateDeviceStatus) SetInstallStatus(value *IosUpdatesInstallStatus)() {
    m.installStatus = value
}
// SetLastReportedDateTime sets the lastReportedDateTime property value. Last modified date time of the policy report.
func (m *IosUpdateDeviceStatus) SetLastReportedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastReportedDateTime = value
}
// SetOsVersion sets the osVersion property value. The device version that is being reported.
func (m *IosUpdateDeviceStatus) SetOsVersion(value *string)() {
    m.osVersion = value
}
// SetStatus sets the status property value. The status property
func (m *IosUpdateDeviceStatus) SetStatus(value *ComplianceStatus)() {
    m.status = value
}
// SetUserId sets the userId property value. The User id that is being reported.
func (m *IosUpdateDeviceStatus) SetUserId(value *string)() {
    m.userId = value
}
// SetUserName sets the userName property value. The User Name that is being reported
func (m *IosUpdateDeviceStatus) SetUserName(value *string)() {
    m.userName = value
}
// SetUserPrincipalName sets the userPrincipalName property value. UserPrincipalName.
func (m *IosUpdateDeviceStatus) SetUserPrincipalName(value *string)() {
    m.userPrincipalName = value
}
