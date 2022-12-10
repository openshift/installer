package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MeetingTimeSuggestion 
type MeetingTimeSuggestion struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // An array that shows the availability status of each attendee for this meeting suggestion.
    attendeeAvailability []AttendeeAvailabilityable
    // A percentage that represents the likelhood of all the attendees attending.
    confidence *float64
    // An array that specifies the name and geographic location of each meeting location for this meeting suggestion.
    locations []Locationable
    // A time period suggested for the meeting.
    meetingTimeSlot TimeSlotable
    // The OdataType property
    odataType *string
    // Order of meeting time suggestions sorted by their computed confidence value from high to low, then by chronology if there are suggestions with the same confidence.
    order *int32
    // Availability of the meeting organizer for this meeting suggestion. The possible values are: free, tentative, busy, oof, workingElsewhere, unknown.
    organizerAvailability *FreeBusyStatus
    // Reason for suggesting the meeting time.
    suggestionReason *string
}
// NewMeetingTimeSuggestion instantiates a new meetingTimeSuggestion and sets the default values.
func NewMeetingTimeSuggestion()(*MeetingTimeSuggestion) {
    m := &MeetingTimeSuggestion{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateMeetingTimeSuggestionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMeetingTimeSuggestionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeetingTimeSuggestion(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MeetingTimeSuggestion) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAttendeeAvailability gets the attendeeAvailability property value. An array that shows the availability status of each attendee for this meeting suggestion.
func (m *MeetingTimeSuggestion) GetAttendeeAvailability()([]AttendeeAvailabilityable) {
    return m.attendeeAvailability
}
// GetConfidence gets the confidence property value. A percentage that represents the likelhood of all the attendees attending.
func (m *MeetingTimeSuggestion) GetConfidence()(*float64) {
    return m.confidence
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MeetingTimeSuggestion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attendeeAvailability"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAttendeeAvailabilityFromDiscriminatorValue , m.SetAttendeeAvailability)
    res["confidence"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetConfidence)
    res["locations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateLocationFromDiscriminatorValue , m.SetLocations)
    res["meetingTimeSlot"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTimeSlotFromDiscriminatorValue , m.SetMeetingTimeSlot)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["order"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetOrder)
    res["organizerAvailability"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseFreeBusyStatus , m.SetOrganizerAvailability)
    res["suggestionReason"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSuggestionReason)
    return res
}
// GetLocations gets the locations property value. An array that specifies the name and geographic location of each meeting location for this meeting suggestion.
func (m *MeetingTimeSuggestion) GetLocations()([]Locationable) {
    return m.locations
}
// GetMeetingTimeSlot gets the meetingTimeSlot property value. A time period suggested for the meeting.
func (m *MeetingTimeSuggestion) GetMeetingTimeSlot()(TimeSlotable) {
    return m.meetingTimeSlot
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *MeetingTimeSuggestion) GetOdataType()(*string) {
    return m.odataType
}
// GetOrder gets the order property value. Order of meeting time suggestions sorted by their computed confidence value from high to low, then by chronology if there are suggestions with the same confidence.
func (m *MeetingTimeSuggestion) GetOrder()(*int32) {
    return m.order
}
// GetOrganizerAvailability gets the organizerAvailability property value. Availability of the meeting organizer for this meeting suggestion. The possible values are: free, tentative, busy, oof, workingElsewhere, unknown.
func (m *MeetingTimeSuggestion) GetOrganizerAvailability()(*FreeBusyStatus) {
    return m.organizerAvailability
}
// GetSuggestionReason gets the suggestionReason property value. Reason for suggesting the meeting time.
func (m *MeetingTimeSuggestion) GetSuggestionReason()(*string) {
    return m.suggestionReason
}
// Serialize serializes information the current object
func (m *MeetingTimeSuggestion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAttendeeAvailability() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAttendeeAvailability())
        err := writer.WriteCollectionOfObjectValues("attendeeAvailability", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("confidence", m.GetConfidence())
        if err != nil {
            return err
        }
    }
    if m.GetLocations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetLocations())
        err := writer.WriteCollectionOfObjectValues("locations", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("meetingTimeSlot", m.GetMeetingTimeSlot())
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
        err := writer.WriteInt32Value("order", m.GetOrder())
        if err != nil {
            return err
        }
    }
    if m.GetOrganizerAvailability() != nil {
        cast := (*m.GetOrganizerAvailability()).String()
        err := writer.WriteStringValue("organizerAvailability", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("suggestionReason", m.GetSuggestionReason())
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
func (m *MeetingTimeSuggestion) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAttendeeAvailability sets the attendeeAvailability property value. An array that shows the availability status of each attendee for this meeting suggestion.
func (m *MeetingTimeSuggestion) SetAttendeeAvailability(value []AttendeeAvailabilityable)() {
    m.attendeeAvailability = value
}
// SetConfidence sets the confidence property value. A percentage that represents the likelhood of all the attendees attending.
func (m *MeetingTimeSuggestion) SetConfidence(value *float64)() {
    m.confidence = value
}
// SetLocations sets the locations property value. An array that specifies the name and geographic location of each meeting location for this meeting suggestion.
func (m *MeetingTimeSuggestion) SetLocations(value []Locationable)() {
    m.locations = value
}
// SetMeetingTimeSlot sets the meetingTimeSlot property value. A time period suggested for the meeting.
func (m *MeetingTimeSuggestion) SetMeetingTimeSlot(value TimeSlotable)() {
    m.meetingTimeSlot = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MeetingTimeSuggestion) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOrder sets the order property value. Order of meeting time suggestions sorted by their computed confidence value from high to low, then by chronology if there are suggestions with the same confidence.
func (m *MeetingTimeSuggestion) SetOrder(value *int32)() {
    m.order = value
}
// SetOrganizerAvailability sets the organizerAvailability property value. Availability of the meeting organizer for this meeting suggestion. The possible values are: free, tentative, busy, oof, workingElsewhere, unknown.
func (m *MeetingTimeSuggestion) SetOrganizerAvailability(value *FreeBusyStatus)() {
    m.organizerAvailability = value
}
// SetSuggestionReason sets the suggestionReason property value. Reason for suggesting the meeting time.
func (m *MeetingTimeSuggestion) SetSuggestionReason(value *string)() {
    m.suggestionReason = value
}
