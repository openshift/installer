package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MeetingTimeSuggestionsResult 
type MeetingTimeSuggestionsResult struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // A reason for not returning any meeting suggestions. The possible values are: attendeesUnavailable, attendeesUnavailableOrUnknown, locationsUnavailable, organizerUnavailable, or unknown. This property is an empty string if the meetingTimeSuggestions property does include any meeting suggestions.
    emptySuggestionsReason *string
    // An array of meeting suggestions.
    meetingTimeSuggestions []MeetingTimeSuggestionable
    // The OdataType property
    odataType *string
}
// NewMeetingTimeSuggestionsResult instantiates a new meetingTimeSuggestionsResult and sets the default values.
func NewMeetingTimeSuggestionsResult()(*MeetingTimeSuggestionsResult) {
    m := &MeetingTimeSuggestionsResult{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateMeetingTimeSuggestionsResultFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMeetingTimeSuggestionsResultFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeetingTimeSuggestionsResult(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MeetingTimeSuggestionsResult) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetEmptySuggestionsReason gets the emptySuggestionsReason property value. A reason for not returning any meeting suggestions. The possible values are: attendeesUnavailable, attendeesUnavailableOrUnknown, locationsUnavailable, organizerUnavailable, or unknown. This property is an empty string if the meetingTimeSuggestions property does include any meeting suggestions.
func (m *MeetingTimeSuggestionsResult) GetEmptySuggestionsReason()(*string) {
    return m.emptySuggestionsReason
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MeetingTimeSuggestionsResult) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["emptySuggestionsReason"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEmptySuggestionsReason)
    res["meetingTimeSuggestions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMeetingTimeSuggestionFromDiscriminatorValue , m.SetMeetingTimeSuggestions)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetMeetingTimeSuggestions gets the meetingTimeSuggestions property value. An array of meeting suggestions.
func (m *MeetingTimeSuggestionsResult) GetMeetingTimeSuggestions()([]MeetingTimeSuggestionable) {
    return m.meetingTimeSuggestions
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *MeetingTimeSuggestionsResult) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *MeetingTimeSuggestionsResult) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("emptySuggestionsReason", m.GetEmptySuggestionsReason())
        if err != nil {
            return err
        }
    }
    if m.GetMeetingTimeSuggestions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMeetingTimeSuggestions())
        err := writer.WriteCollectionOfObjectValues("meetingTimeSuggestions", cast)
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
func (m *MeetingTimeSuggestionsResult) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetEmptySuggestionsReason sets the emptySuggestionsReason property value. A reason for not returning any meeting suggestions. The possible values are: attendeesUnavailable, attendeesUnavailableOrUnknown, locationsUnavailable, organizerUnavailable, or unknown. This property is an empty string if the meetingTimeSuggestions property does include any meeting suggestions.
func (m *MeetingTimeSuggestionsResult) SetEmptySuggestionsReason(value *string)() {
    m.emptySuggestionsReason = value
}
// SetMeetingTimeSuggestions sets the meetingTimeSuggestions property value. An array of meeting suggestions.
func (m *MeetingTimeSuggestionsResult) SetMeetingTimeSuggestions(value []MeetingTimeSuggestionable)() {
    m.meetingTimeSuggestions = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MeetingTimeSuggestionsResult) SetOdataType(value *string)() {
    m.odataType = value
}
