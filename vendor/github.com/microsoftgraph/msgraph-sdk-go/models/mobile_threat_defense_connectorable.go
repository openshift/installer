package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MobileThreatDefenseConnectorable 
type MobileThreatDefenseConnectorable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAndroidDeviceBlockedOnMissingPartnerData()(*bool)
    GetAndroidEnabled()(*bool)
    GetIosDeviceBlockedOnMissingPartnerData()(*bool)
    GetIosEnabled()(*bool)
    GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPartnerState()(*MobileThreatPartnerTenantState)
    GetPartnerUnresponsivenessThresholdInDays()(*int32)
    GetPartnerUnsupportedOsVersionBlocked()(*bool)
    SetAndroidDeviceBlockedOnMissingPartnerData(value *bool)()
    SetAndroidEnabled(value *bool)()
    SetIosDeviceBlockedOnMissingPartnerData(value *bool)()
    SetIosEnabled(value *bool)()
    SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPartnerState(value *MobileThreatPartnerTenantState)()
    SetPartnerUnresponsivenessThresholdInDays(value *int32)()
    SetPartnerUnsupportedOsVersionBlocked(value *bool)()
}
