package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Schedule 
type Schedule struct {
    Entity
    // Indicates whether the schedule is enabled for the team. Required.
    enabled *bool
    // The offerShiftRequests property
    offerShiftRequests []OfferShiftRequestable
    // Indicates whether offer shift requests are enabled for the schedule.
    offerShiftRequestsEnabled *bool
    // The openShiftChangeRequests property
    openShiftChangeRequests []OpenShiftChangeRequestable
    // The openShifts property
    openShifts []OpenShiftable
    // Indicates whether open shifts are enabled for the schedule.
    openShiftsEnabled *bool
    // The status of the schedule provisioning. The possible values are notStarted, running, completed, failed.
    provisionStatus *OperationStatus
    // Additional information about why schedule provisioning failed.
    provisionStatusCode *string
    // The logical grouping of users in the schedule (usually by role).
    schedulingGroups []SchedulingGroupable
    // The shifts in the schedule.
    shifts []Shiftable
    // The swapShiftsChangeRequests property
    swapShiftsChangeRequests []SwapShiftsChangeRequestable
    // Indicates whether swap shifts requests are enabled for the schedule.
    swapShiftsRequestsEnabled *bool
    // Indicates whether time clock is enabled for the schedule.
    timeClockEnabled *bool
    // The set of reasons for a time off in the schedule.
    timeOffReasons []TimeOffReasonable
    // The timeOffRequests property
    timeOffRequests []TimeOffRequestable
    // Indicates whether time off requests are enabled for the schedule.
    timeOffRequestsEnabled *bool
    // The instances of times off in the schedule.
    timesOff []TimeOffable
    // Indicates the time zone of the schedule team using tz database format. Required.
    timeZone *string
    // The workforceIntegrationIds property
    workforceIntegrationIds []string
}
// NewSchedule instantiates a new schedule and sets the default values.
func NewSchedule()(*Schedule) {
    m := &Schedule{
        Entity: *NewEntity(),
    }
    return m
}
// CreateScheduleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateScheduleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSchedule(), nil
}
// GetEnabled gets the enabled property value. Indicates whether the schedule is enabled for the team. Required.
func (m *Schedule) GetEnabled()(*bool) {
    return m.enabled
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Schedule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["enabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEnabled)
    res["offerShiftRequests"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOfferShiftRequestFromDiscriminatorValue , m.SetOfferShiftRequests)
    res["offerShiftRequestsEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetOfferShiftRequestsEnabled)
    res["openShiftChangeRequests"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOpenShiftChangeRequestFromDiscriminatorValue , m.SetOpenShiftChangeRequests)
    res["openShifts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOpenShiftFromDiscriminatorValue , m.SetOpenShifts)
    res["openShiftsEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetOpenShiftsEnabled)
    res["provisionStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseOperationStatus , m.SetProvisionStatus)
    res["provisionStatusCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetProvisionStatusCode)
    res["schedulingGroups"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSchedulingGroupFromDiscriminatorValue , m.SetSchedulingGroups)
    res["shifts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateShiftFromDiscriminatorValue , m.SetShifts)
    res["swapShiftsChangeRequests"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSwapShiftsChangeRequestFromDiscriminatorValue , m.SetSwapShiftsChangeRequests)
    res["swapShiftsRequestsEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetSwapShiftsRequestsEnabled)
    res["timeClockEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetTimeClockEnabled)
    res["timeOffReasons"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTimeOffReasonFromDiscriminatorValue , m.SetTimeOffReasons)
    res["timeOffRequests"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTimeOffRequestFromDiscriminatorValue , m.SetTimeOffRequests)
    res["timeOffRequestsEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetTimeOffRequestsEnabled)
    res["timesOff"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTimeOffFromDiscriminatorValue , m.SetTimesOff)
    res["timeZone"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTimeZone)
    res["workforceIntegrationIds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetWorkforceIntegrationIds)
    return res
}
// GetOfferShiftRequests gets the offerShiftRequests property value. The offerShiftRequests property
func (m *Schedule) GetOfferShiftRequests()([]OfferShiftRequestable) {
    return m.offerShiftRequests
}
// GetOfferShiftRequestsEnabled gets the offerShiftRequestsEnabled property value. Indicates whether offer shift requests are enabled for the schedule.
func (m *Schedule) GetOfferShiftRequestsEnabled()(*bool) {
    return m.offerShiftRequestsEnabled
}
// GetOpenShiftChangeRequests gets the openShiftChangeRequests property value. The openShiftChangeRequests property
func (m *Schedule) GetOpenShiftChangeRequests()([]OpenShiftChangeRequestable) {
    return m.openShiftChangeRequests
}
// GetOpenShifts gets the openShifts property value. The openShifts property
func (m *Schedule) GetOpenShifts()([]OpenShiftable) {
    return m.openShifts
}
// GetOpenShiftsEnabled gets the openShiftsEnabled property value. Indicates whether open shifts are enabled for the schedule.
func (m *Schedule) GetOpenShiftsEnabled()(*bool) {
    return m.openShiftsEnabled
}
// GetProvisionStatus gets the provisionStatus property value. The status of the schedule provisioning. The possible values are notStarted, running, completed, failed.
func (m *Schedule) GetProvisionStatus()(*OperationStatus) {
    return m.provisionStatus
}
// GetProvisionStatusCode gets the provisionStatusCode property value. Additional information about why schedule provisioning failed.
func (m *Schedule) GetProvisionStatusCode()(*string) {
    return m.provisionStatusCode
}
// GetSchedulingGroups gets the schedulingGroups property value. The logical grouping of users in the schedule (usually by role).
func (m *Schedule) GetSchedulingGroups()([]SchedulingGroupable) {
    return m.schedulingGroups
}
// GetShifts gets the shifts property value. The shifts in the schedule.
func (m *Schedule) GetShifts()([]Shiftable) {
    return m.shifts
}
// GetSwapShiftsChangeRequests gets the swapShiftsChangeRequests property value. The swapShiftsChangeRequests property
func (m *Schedule) GetSwapShiftsChangeRequests()([]SwapShiftsChangeRequestable) {
    return m.swapShiftsChangeRequests
}
// GetSwapShiftsRequestsEnabled gets the swapShiftsRequestsEnabled property value. Indicates whether swap shifts requests are enabled for the schedule.
func (m *Schedule) GetSwapShiftsRequestsEnabled()(*bool) {
    return m.swapShiftsRequestsEnabled
}
// GetTimeClockEnabled gets the timeClockEnabled property value. Indicates whether time clock is enabled for the schedule.
func (m *Schedule) GetTimeClockEnabled()(*bool) {
    return m.timeClockEnabled
}
// GetTimeOffReasons gets the timeOffReasons property value. The set of reasons for a time off in the schedule.
func (m *Schedule) GetTimeOffReasons()([]TimeOffReasonable) {
    return m.timeOffReasons
}
// GetTimeOffRequests gets the timeOffRequests property value. The timeOffRequests property
func (m *Schedule) GetTimeOffRequests()([]TimeOffRequestable) {
    return m.timeOffRequests
}
// GetTimeOffRequestsEnabled gets the timeOffRequestsEnabled property value. Indicates whether time off requests are enabled for the schedule.
func (m *Schedule) GetTimeOffRequestsEnabled()(*bool) {
    return m.timeOffRequestsEnabled
}
// GetTimesOff gets the timesOff property value. The instances of times off in the schedule.
func (m *Schedule) GetTimesOff()([]TimeOffable) {
    return m.timesOff
}
// GetTimeZone gets the timeZone property value. Indicates the time zone of the schedule team using tz database format. Required.
func (m *Schedule) GetTimeZone()(*string) {
    return m.timeZone
}
// GetWorkforceIntegrationIds gets the workforceIntegrationIds property value. The workforceIntegrationIds property
func (m *Schedule) GetWorkforceIntegrationIds()([]string) {
    return m.workforceIntegrationIds
}
// Serialize serializes information the current object
func (m *Schedule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("enabled", m.GetEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetOfferShiftRequests() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOfferShiftRequests())
        err = writer.WriteCollectionOfObjectValues("offerShiftRequests", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("offerShiftRequestsEnabled", m.GetOfferShiftRequestsEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetOpenShiftChangeRequests() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOpenShiftChangeRequests())
        err = writer.WriteCollectionOfObjectValues("openShiftChangeRequests", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOpenShifts() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOpenShifts())
        err = writer.WriteCollectionOfObjectValues("openShifts", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("openShiftsEnabled", m.GetOpenShiftsEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetSchedulingGroups() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSchedulingGroups())
        err = writer.WriteCollectionOfObjectValues("schedulingGroups", cast)
        if err != nil {
            return err
        }
    }
    if m.GetShifts() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetShifts())
        err = writer.WriteCollectionOfObjectValues("shifts", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSwapShiftsChangeRequests() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSwapShiftsChangeRequests())
        err = writer.WriteCollectionOfObjectValues("swapShiftsChangeRequests", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("swapShiftsRequestsEnabled", m.GetSwapShiftsRequestsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("timeClockEnabled", m.GetTimeClockEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetTimeOffReasons() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTimeOffReasons())
        err = writer.WriteCollectionOfObjectValues("timeOffReasons", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTimeOffRequests() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTimeOffRequests())
        err = writer.WriteCollectionOfObjectValues("timeOffRequests", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("timeOffRequestsEnabled", m.GetTimeOffRequestsEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetTimesOff() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTimesOff())
        err = writer.WriteCollectionOfObjectValues("timesOff", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("timeZone", m.GetTimeZone())
        if err != nil {
            return err
        }
    }
    if m.GetWorkforceIntegrationIds() != nil {
        err = writer.WriteCollectionOfStringValues("workforceIntegrationIds", m.GetWorkforceIntegrationIds())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEnabled sets the enabled property value. Indicates whether the schedule is enabled for the team. Required.
func (m *Schedule) SetEnabled(value *bool)() {
    m.enabled = value
}
// SetOfferShiftRequests sets the offerShiftRequests property value. The offerShiftRequests property
func (m *Schedule) SetOfferShiftRequests(value []OfferShiftRequestable)() {
    m.offerShiftRequests = value
}
// SetOfferShiftRequestsEnabled sets the offerShiftRequestsEnabled property value. Indicates whether offer shift requests are enabled for the schedule.
func (m *Schedule) SetOfferShiftRequestsEnabled(value *bool)() {
    m.offerShiftRequestsEnabled = value
}
// SetOpenShiftChangeRequests sets the openShiftChangeRequests property value. The openShiftChangeRequests property
func (m *Schedule) SetOpenShiftChangeRequests(value []OpenShiftChangeRequestable)() {
    m.openShiftChangeRequests = value
}
// SetOpenShifts sets the openShifts property value. The openShifts property
func (m *Schedule) SetOpenShifts(value []OpenShiftable)() {
    m.openShifts = value
}
// SetOpenShiftsEnabled sets the openShiftsEnabled property value. Indicates whether open shifts are enabled for the schedule.
func (m *Schedule) SetOpenShiftsEnabled(value *bool)() {
    m.openShiftsEnabled = value
}
// SetProvisionStatus sets the provisionStatus property value. The status of the schedule provisioning. The possible values are notStarted, running, completed, failed.
func (m *Schedule) SetProvisionStatus(value *OperationStatus)() {
    m.provisionStatus = value
}
// SetProvisionStatusCode sets the provisionStatusCode property value. Additional information about why schedule provisioning failed.
func (m *Schedule) SetProvisionStatusCode(value *string)() {
    m.provisionStatusCode = value
}
// SetSchedulingGroups sets the schedulingGroups property value. The logical grouping of users in the schedule (usually by role).
func (m *Schedule) SetSchedulingGroups(value []SchedulingGroupable)() {
    m.schedulingGroups = value
}
// SetShifts sets the shifts property value. The shifts in the schedule.
func (m *Schedule) SetShifts(value []Shiftable)() {
    m.shifts = value
}
// SetSwapShiftsChangeRequests sets the swapShiftsChangeRequests property value. The swapShiftsChangeRequests property
func (m *Schedule) SetSwapShiftsChangeRequests(value []SwapShiftsChangeRequestable)() {
    m.swapShiftsChangeRequests = value
}
// SetSwapShiftsRequestsEnabled sets the swapShiftsRequestsEnabled property value. Indicates whether swap shifts requests are enabled for the schedule.
func (m *Schedule) SetSwapShiftsRequestsEnabled(value *bool)() {
    m.swapShiftsRequestsEnabled = value
}
// SetTimeClockEnabled sets the timeClockEnabled property value. Indicates whether time clock is enabled for the schedule.
func (m *Schedule) SetTimeClockEnabled(value *bool)() {
    m.timeClockEnabled = value
}
// SetTimeOffReasons sets the timeOffReasons property value. The set of reasons for a time off in the schedule.
func (m *Schedule) SetTimeOffReasons(value []TimeOffReasonable)() {
    m.timeOffReasons = value
}
// SetTimeOffRequests sets the timeOffRequests property value. The timeOffRequests property
func (m *Schedule) SetTimeOffRequests(value []TimeOffRequestable)() {
    m.timeOffRequests = value
}
// SetTimeOffRequestsEnabled sets the timeOffRequestsEnabled property value. Indicates whether time off requests are enabled for the schedule.
func (m *Schedule) SetTimeOffRequestsEnabled(value *bool)() {
    m.timeOffRequestsEnabled = value
}
// SetTimesOff sets the timesOff property value. The instances of times off in the schedule.
func (m *Schedule) SetTimesOff(value []TimeOffable)() {
    m.timesOff = value
}
// SetTimeZone sets the timeZone property value. Indicates the time zone of the schedule team using tz database format. Required.
func (m *Schedule) SetTimeZone(value *string)() {
    m.timeZone = value
}
// SetWorkforceIntegrationIds sets the workforceIntegrationIds property value. The workforceIntegrationIds property
func (m *Schedule) SetWorkforceIntegrationIds(value []string)() {
    m.workforceIntegrationIds = value
}
