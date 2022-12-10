package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MeetingAttendanceReport provides operations to manage the collection of agreement entities.
type MeetingAttendanceReport struct {
    Entity
    // List of attendance records of an attendance report. Read-only.
    attendanceRecords []AttendanceRecordable
    // UTC time when the meeting ended. Read-only.
    meetingEndDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // UTC time when the meeting started. Read-only.
    meetingStartDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Total number of participants. Read-only.
    totalParticipantCount *int32
}
// NewMeetingAttendanceReport instantiates a new meetingAttendanceReport and sets the default values.
func NewMeetingAttendanceReport()(*MeetingAttendanceReport) {
    m := &MeetingAttendanceReport{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMeetingAttendanceReportFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMeetingAttendanceReportFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeetingAttendanceReport(), nil
}
// GetAttendanceRecords gets the attendanceRecords property value. List of attendance records of an attendance report. Read-only.
func (m *MeetingAttendanceReport) GetAttendanceRecords()([]AttendanceRecordable) {
    return m.attendanceRecords
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MeetingAttendanceReport) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attendanceRecords"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAttendanceRecordFromDiscriminatorValue , m.SetAttendanceRecords)
    res["meetingEndDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetMeetingEndDateTime)
    res["meetingStartDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetMeetingStartDateTime)
    res["totalParticipantCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetTotalParticipantCount)
    return res
}
// GetMeetingEndDateTime gets the meetingEndDateTime property value. UTC time when the meeting ended. Read-only.
func (m *MeetingAttendanceReport) GetMeetingEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.meetingEndDateTime
}
// GetMeetingStartDateTime gets the meetingStartDateTime property value. UTC time when the meeting started. Read-only.
func (m *MeetingAttendanceReport) GetMeetingStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.meetingStartDateTime
}
// GetTotalParticipantCount gets the totalParticipantCount property value. Total number of participants. Read-only.
func (m *MeetingAttendanceReport) GetTotalParticipantCount()(*int32) {
    return m.totalParticipantCount
}
// Serialize serializes information the current object
func (m *MeetingAttendanceReport) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAttendanceRecords() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAttendanceRecords())
        err = writer.WriteCollectionOfObjectValues("attendanceRecords", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("meetingEndDateTime", m.GetMeetingEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("meetingStartDateTime", m.GetMeetingStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalParticipantCount", m.GetTotalParticipantCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttendanceRecords sets the attendanceRecords property value. List of attendance records of an attendance report. Read-only.
func (m *MeetingAttendanceReport) SetAttendanceRecords(value []AttendanceRecordable)() {
    m.attendanceRecords = value
}
// SetMeetingEndDateTime sets the meetingEndDateTime property value. UTC time when the meeting ended. Read-only.
func (m *MeetingAttendanceReport) SetMeetingEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.meetingEndDateTime = value
}
// SetMeetingStartDateTime sets the meetingStartDateTime property value. UTC time when the meeting started. Read-only.
func (m *MeetingAttendanceReport) SetMeetingStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.meetingStartDateTime = value
}
// SetTotalParticipantCount sets the totalParticipantCount property value. Total number of participants. Read-only.
func (m *MeetingAttendanceReport) SetTotalParticipantCount(value *int32)() {
    m.totalParticipantCount = value
}
