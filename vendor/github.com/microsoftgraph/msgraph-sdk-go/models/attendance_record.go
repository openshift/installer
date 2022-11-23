package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AttendanceRecord provides operations to manage the collection of agreement entities.
type AttendanceRecord struct {
    Entity
    // List of time periods between joining and leaving a meeting.
    attendanceIntervals []AttendanceIntervalable
    // Email address of the user associated with this atttendance record.
    emailAddress *string
    // Identity of the user associated with this atttendance record.
    identity Identityable
    // Role of the attendee. Possible values are: None, Attendee, Presenter, and Organizer.
    role *string
    // Total duration of the attendances in seconds.
    totalAttendanceInSeconds *int32
}
// NewAttendanceRecord instantiates a new attendanceRecord and sets the default values.
func NewAttendanceRecord()(*AttendanceRecord) {
    m := &AttendanceRecord{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAttendanceRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAttendanceRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttendanceRecord(), nil
}
// GetAttendanceIntervals gets the attendanceIntervals property value. List of time periods between joining and leaving a meeting.
func (m *AttendanceRecord) GetAttendanceIntervals()([]AttendanceIntervalable) {
    return m.attendanceIntervals
}
// GetEmailAddress gets the emailAddress property value. Email address of the user associated with this atttendance record.
func (m *AttendanceRecord) GetEmailAddress()(*string) {
    return m.emailAddress
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AttendanceRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attendanceIntervals"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAttendanceIntervalFromDiscriminatorValue , m.SetAttendanceIntervals)
    res["emailAddress"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEmailAddress)
    res["identity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentityFromDiscriminatorValue , m.SetIdentity)
    res["role"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRole)
    res["totalAttendanceInSeconds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetTotalAttendanceInSeconds)
    return res
}
// GetIdentity gets the identity property value. Identity of the user associated with this atttendance record.
func (m *AttendanceRecord) GetIdentity()(Identityable) {
    return m.identity
}
// GetRole gets the role property value. Role of the attendee. Possible values are: None, Attendee, Presenter, and Organizer.
func (m *AttendanceRecord) GetRole()(*string) {
    return m.role
}
// GetTotalAttendanceInSeconds gets the totalAttendanceInSeconds property value. Total duration of the attendances in seconds.
func (m *AttendanceRecord) GetTotalAttendanceInSeconds()(*int32) {
    return m.totalAttendanceInSeconds
}
// Serialize serializes information the current object
func (m *AttendanceRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAttendanceIntervals() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAttendanceIntervals())
        err = writer.WriteCollectionOfObjectValues("attendanceIntervals", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("identity", m.GetIdentity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("role", m.GetRole())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalAttendanceInSeconds", m.GetTotalAttendanceInSeconds())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttendanceIntervals sets the attendanceIntervals property value. List of time periods between joining and leaving a meeting.
func (m *AttendanceRecord) SetAttendanceIntervals(value []AttendanceIntervalable)() {
    m.attendanceIntervals = value
}
// SetEmailAddress sets the emailAddress property value. Email address of the user associated with this atttendance record.
func (m *AttendanceRecord) SetEmailAddress(value *string)() {
    m.emailAddress = value
}
// SetIdentity sets the identity property value. Identity of the user associated with this atttendance record.
func (m *AttendanceRecord) SetIdentity(value Identityable)() {
    m.identity = value
}
// SetRole sets the role property value. Role of the attendee. Possible values are: None, Attendee, Presenter, and Organizer.
func (m *AttendanceRecord) SetRole(value *string)() {
    m.role = value
}
// SetTotalAttendanceInSeconds sets the totalAttendanceInSeconds property value. Total duration of the attendances in seconds.
func (m *AttendanceRecord) SetTotalAttendanceInSeconds(value *int32)() {
    m.totalAttendanceInSeconds = value
}
