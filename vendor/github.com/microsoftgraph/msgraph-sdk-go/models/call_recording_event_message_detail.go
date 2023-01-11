package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CallRecordingEventMessageDetail 
type CallRecordingEventMessageDetail struct {
    EventMessageDetail
    // Unique identifier of the call.
    callId *string
    // Display name for the call recording.
    callRecordingDisplayName *string
    // Duration of the call recording.
    callRecordingDuration *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration
    // Status of the call recording. Possible values are: success, failure, initial, chunkFinished, unknownFutureValue.
    callRecordingStatus *CallRecordingStatus
    // Call recording URL.
    callRecordingUrl *string
    // Initiator of the event.
    initiator IdentitySetable
    // Organizer of the meeting.
    meetingOrganizer IdentitySetable
}
// NewCallRecordingEventMessageDetail instantiates a new CallRecordingEventMessageDetail and sets the default values.
func NewCallRecordingEventMessageDetail()(*CallRecordingEventMessageDetail) {
    m := &CallRecordingEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.callRecordingEventMessageDetail";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateCallRecordingEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCallRecordingEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCallRecordingEventMessageDetail(), nil
}
// GetCallId gets the callId property value. Unique identifier of the call.
func (m *CallRecordingEventMessageDetail) GetCallId()(*string) {
    return m.callId
}
// GetCallRecordingDisplayName gets the callRecordingDisplayName property value. Display name for the call recording.
func (m *CallRecordingEventMessageDetail) GetCallRecordingDisplayName()(*string) {
    return m.callRecordingDisplayName
}
// GetCallRecordingDuration gets the callRecordingDuration property value. Duration of the call recording.
func (m *CallRecordingEventMessageDetail) GetCallRecordingDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    return m.callRecordingDuration
}
// GetCallRecordingStatus gets the callRecordingStatus property value. Status of the call recording. Possible values are: success, failure, initial, chunkFinished, unknownFutureValue.
func (m *CallRecordingEventMessageDetail) GetCallRecordingStatus()(*CallRecordingStatus) {
    return m.callRecordingStatus
}
// GetCallRecordingUrl gets the callRecordingUrl property value. Call recording URL.
func (m *CallRecordingEventMessageDetail) GetCallRecordingUrl()(*string) {
    return m.callRecordingUrl
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CallRecordingEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["callId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCallId)
    res["callRecordingDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCallRecordingDisplayName)
    res["callRecordingDuration"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetISODurationValue(m.SetCallRecordingDuration)
    res["callRecordingStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseCallRecordingStatus , m.SetCallRecordingStatus)
    res["callRecordingUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCallRecordingUrl)
    res["initiator"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetInitiator)
    res["meetingOrganizer"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetMeetingOrganizer)
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
func (m *CallRecordingEventMessageDetail) GetInitiator()(IdentitySetable) {
    return m.initiator
}
// GetMeetingOrganizer gets the meetingOrganizer property value. Organizer of the meeting.
func (m *CallRecordingEventMessageDetail) GetMeetingOrganizer()(IdentitySetable) {
    return m.meetingOrganizer
}
// Serialize serializes information the current object
func (m *CallRecordingEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("callId", m.GetCallId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callRecordingDisplayName", m.GetCallRecordingDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("callRecordingDuration", m.GetCallRecordingDuration())
        if err != nil {
            return err
        }
    }
    if m.GetCallRecordingStatus() != nil {
        cast := (*m.GetCallRecordingStatus()).String()
        err = writer.WriteStringValue("callRecordingStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callRecordingUrl", m.GetCallRecordingUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("meetingOrganizer", m.GetMeetingOrganizer())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCallId sets the callId property value. Unique identifier of the call.
func (m *CallRecordingEventMessageDetail) SetCallId(value *string)() {
    m.callId = value
}
// SetCallRecordingDisplayName sets the callRecordingDisplayName property value. Display name for the call recording.
func (m *CallRecordingEventMessageDetail) SetCallRecordingDisplayName(value *string)() {
    m.callRecordingDisplayName = value
}
// SetCallRecordingDuration sets the callRecordingDuration property value. Duration of the call recording.
func (m *CallRecordingEventMessageDetail) SetCallRecordingDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    m.callRecordingDuration = value
}
// SetCallRecordingStatus sets the callRecordingStatus property value. Status of the call recording. Possible values are: success, failure, initial, chunkFinished, unknownFutureValue.
func (m *CallRecordingEventMessageDetail) SetCallRecordingStatus(value *CallRecordingStatus)() {
    m.callRecordingStatus = value
}
// SetCallRecordingUrl sets the callRecordingUrl property value. Call recording URL.
func (m *CallRecordingEventMessageDetail) SetCallRecordingUrl(value *string)() {
    m.callRecordingUrl = value
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *CallRecordingEventMessageDetail) SetInitiator(value IdentitySetable)() {
    m.initiator = value
}
// SetMeetingOrganizer sets the meetingOrganizer property value. Organizer of the meeting.
func (m *CallRecordingEventMessageDetail) SetMeetingOrganizer(value IdentitySetable)() {
    m.meetingOrganizer = value
}
