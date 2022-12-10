package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AttendanceInterval 
type AttendanceInterval struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Duration of the meeting interval in seconds; that is, the difference between joinDateTime and leaveDateTime.
    durationInSeconds *int32
    // The time the attendee joined in UTC.
    joinDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The time the attendee left in UTC.
    leaveDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The OdataType property
    odataType *string
}
// NewAttendanceInterval instantiates a new attendanceInterval and sets the default values.
func NewAttendanceInterval()(*AttendanceInterval) {
    m := &AttendanceInterval{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAttendanceIntervalFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAttendanceIntervalFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttendanceInterval(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AttendanceInterval) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDurationInSeconds gets the durationInSeconds property value. Duration of the meeting interval in seconds; that is, the difference between joinDateTime and leaveDateTime.
func (m *AttendanceInterval) GetDurationInSeconds()(*int32) {
    return m.durationInSeconds
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AttendanceInterval) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["durationInSeconds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetDurationInSeconds)
    res["joinDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetJoinDateTime)
    res["leaveDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLeaveDateTime)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetJoinDateTime gets the joinDateTime property value. The time the attendee joined in UTC.
func (m *AttendanceInterval) GetJoinDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.joinDateTime
}
// GetLeaveDateTime gets the leaveDateTime property value. The time the attendee left in UTC.
func (m *AttendanceInterval) GetLeaveDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.leaveDateTime
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AttendanceInterval) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *AttendanceInterval) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("durationInSeconds", m.GetDurationInSeconds())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("joinDateTime", m.GetJoinDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("leaveDateTime", m.GetLeaveDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AttendanceInterval) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDurationInSeconds sets the durationInSeconds property value. Duration of the meeting interval in seconds; that is, the difference between joinDateTime and leaveDateTime.
func (m *AttendanceInterval) SetDurationInSeconds(value *int32)() {
    m.durationInSeconds = value
}
// SetJoinDateTime sets the joinDateTime property value. The time the attendee joined in UTC.
func (m *AttendanceInterval) SetJoinDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.joinDateTime = value
}
// SetLeaveDateTime sets the leaveDateTime property value. The time the attendee left in UTC.
func (m *AttendanceInterval) SetLeaveDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.leaveDateTime = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AttendanceInterval) SetOdataType(value *string)() {
    m.odataType = value
}
