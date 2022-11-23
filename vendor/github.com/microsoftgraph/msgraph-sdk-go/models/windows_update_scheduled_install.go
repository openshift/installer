package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsUpdateScheduledInstall 
type WindowsUpdateScheduledInstall struct {
    WindowsUpdateInstallScheduleType
    // Possible values for a weekly schedule.
    scheduledInstallDay *WeeklySchedule
    // Scheduled Install Time during day
    scheduledInstallTime *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly
}
// NewWindowsUpdateScheduledInstall instantiates a new WindowsUpdateScheduledInstall and sets the default values.
func NewWindowsUpdateScheduledInstall()(*WindowsUpdateScheduledInstall) {
    m := &WindowsUpdateScheduledInstall{
        WindowsUpdateInstallScheduleType: *NewWindowsUpdateInstallScheduleType(),
    }
    odataTypeValue := "#microsoft.graph.windowsUpdateScheduledInstall";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWindowsUpdateScheduledInstallFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsUpdateScheduledInstallFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsUpdateScheduledInstall(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsUpdateScheduledInstall) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WindowsUpdateInstallScheduleType.GetFieldDeserializers()
    res["scheduledInstallDay"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWeeklySchedule , m.SetScheduledInstallDay)
    res["scheduledInstallTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeOnlyValue(m.SetScheduledInstallTime)
    return res
}
// GetScheduledInstallDay gets the scheduledInstallDay property value. Possible values for a weekly schedule.
func (m *WindowsUpdateScheduledInstall) GetScheduledInstallDay()(*WeeklySchedule) {
    return m.scheduledInstallDay
}
// GetScheduledInstallTime gets the scheduledInstallTime property value. Scheduled Install Time during day
func (m *WindowsUpdateScheduledInstall) GetScheduledInstallTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly) {
    return m.scheduledInstallTime
}
// Serialize serializes information the current object
func (m *WindowsUpdateScheduledInstall) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WindowsUpdateInstallScheduleType.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetScheduledInstallDay() != nil {
        cast := (*m.GetScheduledInstallDay()).String()
        err = writer.WriteStringValue("scheduledInstallDay", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeOnlyValue("scheduledInstallTime", m.GetScheduledInstallTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetScheduledInstallDay sets the scheduledInstallDay property value. Possible values for a weekly schedule.
func (m *WindowsUpdateScheduledInstall) SetScheduledInstallDay(value *WeeklySchedule)() {
    m.scheduledInstallDay = value
}
// SetScheduledInstallTime sets the scheduledInstallTime property value. Scheduled Install Time during day
func (m *WindowsUpdateScheduledInstall) SetScheduledInstallTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)() {
    m.scheduledInstallTime = value
}
