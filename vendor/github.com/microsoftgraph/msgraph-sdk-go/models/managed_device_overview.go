package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedDeviceOverview 
type ManagedDeviceOverview struct {
    Entity
    // Distribution of Exchange Access State in Intune
    deviceExchangeAccessStateSummary DeviceExchangeAccessStateSummaryable
    // Device operating system summary.
    deviceOperatingSystemSummary DeviceOperatingSystemSummaryable
    // The number of devices enrolled in both MDM and EAS
    dualEnrolledDeviceCount *int32
    // Total enrolled device count. Does not include PC devices managed via Intune PC Agent
    enrolledDeviceCount *int32
    // The number of devices enrolled in MDM
    mdmEnrolledCount *int32
}
// NewManagedDeviceOverview instantiates a new managedDeviceOverview and sets the default values.
func NewManagedDeviceOverview()(*ManagedDeviceOverview) {
    m := &ManagedDeviceOverview{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedDeviceOverviewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateManagedDeviceOverviewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedDeviceOverview(), nil
}
// GetDeviceExchangeAccessStateSummary gets the deviceExchangeAccessStateSummary property value. Distribution of Exchange Access State in Intune
func (m *ManagedDeviceOverview) GetDeviceExchangeAccessStateSummary()(DeviceExchangeAccessStateSummaryable) {
    return m.deviceExchangeAccessStateSummary
}
// GetDeviceOperatingSystemSummary gets the deviceOperatingSystemSummary property value. Device operating system summary.
func (m *ManagedDeviceOverview) GetDeviceOperatingSystemSummary()(DeviceOperatingSystemSummaryable) {
    return m.deviceOperatingSystemSummary
}
// GetDualEnrolledDeviceCount gets the dualEnrolledDeviceCount property value. The number of devices enrolled in both MDM and EAS
func (m *ManagedDeviceOverview) GetDualEnrolledDeviceCount()(*int32) {
    return m.dualEnrolledDeviceCount
}
// GetEnrolledDeviceCount gets the enrolledDeviceCount property value. Total enrolled device count. Does not include PC devices managed via Intune PC Agent
func (m *ManagedDeviceOverview) GetEnrolledDeviceCount()(*int32) {
    return m.enrolledDeviceCount
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ManagedDeviceOverview) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deviceExchangeAccessStateSummary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceExchangeAccessStateSummaryFromDiscriminatorValue , m.SetDeviceExchangeAccessStateSummary)
    res["deviceOperatingSystemSummary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceOperatingSystemSummaryFromDiscriminatorValue , m.SetDeviceOperatingSystemSummary)
    res["dualEnrolledDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetDualEnrolledDeviceCount)
    res["enrolledDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetEnrolledDeviceCount)
    res["mdmEnrolledCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetMdmEnrolledCount)
    return res
}
// GetMdmEnrolledCount gets the mdmEnrolledCount property value. The number of devices enrolled in MDM
func (m *ManagedDeviceOverview) GetMdmEnrolledCount()(*int32) {
    return m.mdmEnrolledCount
}
// Serialize serializes information the current object
func (m *ManagedDeviceOverview) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("deviceExchangeAccessStateSummary", m.GetDeviceExchangeAccessStateSummary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deviceOperatingSystemSummary", m.GetDeviceOperatingSystemSummary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("dualEnrolledDeviceCount", m.GetDualEnrolledDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("enrolledDeviceCount", m.GetEnrolledDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("mdmEnrolledCount", m.GetMdmEnrolledCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceExchangeAccessStateSummary sets the deviceExchangeAccessStateSummary property value. Distribution of Exchange Access State in Intune
func (m *ManagedDeviceOverview) SetDeviceExchangeAccessStateSummary(value DeviceExchangeAccessStateSummaryable)() {
    m.deviceExchangeAccessStateSummary = value
}
// SetDeviceOperatingSystemSummary sets the deviceOperatingSystemSummary property value. Device operating system summary.
func (m *ManagedDeviceOverview) SetDeviceOperatingSystemSummary(value DeviceOperatingSystemSummaryable)() {
    m.deviceOperatingSystemSummary = value
}
// SetDualEnrolledDeviceCount sets the dualEnrolledDeviceCount property value. The number of devices enrolled in both MDM and EAS
func (m *ManagedDeviceOverview) SetDualEnrolledDeviceCount(value *int32)() {
    m.dualEnrolledDeviceCount = value
}
// SetEnrolledDeviceCount sets the enrolledDeviceCount property value. Total enrolled device count. Does not include PC devices managed via Intune PC Agent
func (m *ManagedDeviceOverview) SetEnrolledDeviceCount(value *int32)() {
    m.enrolledDeviceCount = value
}
// SetMdmEnrolledCount sets the mdmEnrolledCount property value. The number of devices enrolled in MDM
func (m *ManagedDeviceOverview) SetMdmEnrolledCount(value *int32)() {
    m.mdmEnrolledCount = value
}
