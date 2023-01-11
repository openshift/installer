package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagementPartner entity which represents a connection to device management partner.
type DeviceManagementPartner struct {
    Entity
    // Partner display name
    displayName *string
    // Whether device management partner is configured or not
    isConfigured *bool
    // Timestamp of last heartbeat after admin enabled option Connect to Device management Partner
    lastHeartbeatDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Partner App Type.
    partnerAppType *DeviceManagementPartnerAppType
    // Partner state of this tenant.
    partnerState *DeviceManagementPartnerTenantState
    // Partner Single tenant App id
    singleTenantAppId *string
    // DateTime in UTC when PartnerDevices will be marked as NonCompliant
    whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // DateTime in UTC when PartnerDevices will be removed
    whenPartnerDevicesWillBeRemovedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
}
// NewDeviceManagementPartner instantiates a new deviceManagementPartner and sets the default values.
func NewDeviceManagementPartner()(*DeviceManagementPartner) {
    m := &DeviceManagementPartner{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceManagementPartnerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDeviceManagementPartnerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceManagementPartner(), nil
}
// GetDisplayName gets the displayName property value. Partner display name
func (m *DeviceManagementPartner) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceManagementPartner) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["isConfigured"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsConfigured)
    res["lastHeartbeatDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastHeartbeatDateTime)
    res["partnerAppType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceManagementPartnerAppType , m.SetPartnerAppType)
    res["partnerState"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceManagementPartnerTenantState , m.SetPartnerState)
    res["singleTenantAppId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSingleTenantAppId)
    res["whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime)
    res["whenPartnerDevicesWillBeRemovedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetWhenPartnerDevicesWillBeRemovedDateTime)
    return res
}
// GetIsConfigured gets the isConfigured property value. Whether device management partner is configured or not
func (m *DeviceManagementPartner) GetIsConfigured()(*bool) {
    return m.isConfigured
}
// GetLastHeartbeatDateTime gets the lastHeartbeatDateTime property value. Timestamp of last heartbeat after admin enabled option Connect to Device management Partner
func (m *DeviceManagementPartner) GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastHeartbeatDateTime
}
// GetPartnerAppType gets the partnerAppType property value. Partner App Type.
func (m *DeviceManagementPartner) GetPartnerAppType()(*DeviceManagementPartnerAppType) {
    return m.partnerAppType
}
// GetPartnerState gets the partnerState property value. Partner state of this tenant.
func (m *DeviceManagementPartner) GetPartnerState()(*DeviceManagementPartnerTenantState) {
    return m.partnerState
}
// GetSingleTenantAppId gets the singleTenantAppId property value. Partner Single tenant App id
func (m *DeviceManagementPartner) GetSingleTenantAppId()(*string) {
    return m.singleTenantAppId
}
// GetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime gets the whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime property value. DateTime in UTC when PartnerDevices will be marked as NonCompliant
func (m *DeviceManagementPartner) GetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime
}
// GetWhenPartnerDevicesWillBeRemovedDateTime gets the whenPartnerDevicesWillBeRemovedDateTime property value. DateTime in UTC when PartnerDevices will be removed
func (m *DeviceManagementPartner) GetWhenPartnerDevicesWillBeRemovedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.whenPartnerDevicesWillBeRemovedDateTime
}
// Serialize serializes information the current object
func (m *DeviceManagementPartner) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isConfigured", m.GetIsConfigured())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastHeartbeatDateTime", m.GetLastHeartbeatDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetPartnerAppType() != nil {
        cast := (*m.GetPartnerAppType()).String()
        err = writer.WriteStringValue("partnerAppType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetPartnerState() != nil {
        cast := (*m.GetPartnerState()).String()
        err = writer.WriteStringValue("partnerState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("singleTenantAppId", m.GetSingleTenantAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime", m.GetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("whenPartnerDevicesWillBeRemovedDateTime", m.GetWhenPartnerDevicesWillBeRemovedDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. Partner display name
func (m *DeviceManagementPartner) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetIsConfigured sets the isConfigured property value. Whether device management partner is configured or not
func (m *DeviceManagementPartner) SetIsConfigured(value *bool)() {
    m.isConfigured = value
}
// SetLastHeartbeatDateTime sets the lastHeartbeatDateTime property value. Timestamp of last heartbeat after admin enabled option Connect to Device management Partner
func (m *DeviceManagementPartner) SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastHeartbeatDateTime = value
}
// SetPartnerAppType sets the partnerAppType property value. Partner App Type.
func (m *DeviceManagementPartner) SetPartnerAppType(value *DeviceManagementPartnerAppType)() {
    m.partnerAppType = value
}
// SetPartnerState sets the partnerState property value. Partner state of this tenant.
func (m *DeviceManagementPartner) SetPartnerState(value *DeviceManagementPartnerTenantState)() {
    m.partnerState = value
}
// SetSingleTenantAppId sets the singleTenantAppId property value. Partner Single tenant App id
func (m *DeviceManagementPartner) SetSingleTenantAppId(value *string)() {
    m.singleTenantAppId = value
}
// SetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime sets the whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime property value. DateTime in UTC when PartnerDevices will be marked as NonCompliant
func (m *DeviceManagementPartner) SetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime = value
}
// SetWhenPartnerDevicesWillBeRemovedDateTime sets the whenPartnerDevicesWillBeRemovedDateTime property value. DateTime in UTC when PartnerDevices will be removed
func (m *DeviceManagementPartner) SetWhenPartnerDevicesWillBeRemovedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.whenPartnerDevicesWillBeRemovedDateTime = value
}
