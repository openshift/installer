package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AttendeeBase 
type AttendeeBase struct {
    Recipient
    // The type of attendee. The possible values are: required, optional, resource. Currently if the attendee is a person, findMeetingTimes always considers the person is of the Required type.
    type_escaped *AttendeeType
}
// NewAttendeeBase instantiates a new AttendeeBase and sets the default values.
func NewAttendeeBase()(*AttendeeBase) {
    m := &AttendeeBase{
        Recipient: *NewRecipient(),
    }
    odataTypeValue := "#microsoft.graph.attendeeBase";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateAttendeeBaseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAttendeeBaseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.attendee":
                        return NewAttendee(), nil
                }
            }
        }
    }
    return NewAttendeeBase(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AttendeeBase) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Recipient.GetFieldDeserializers()
    res["type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseAttendeeType , m.SetType)
    return res
}
// GetType gets the type property value. The type of attendee. The possible values are: required, optional, resource. Currently if the attendee is a person, findMeetingTimes always considers the person is of the Required type.
func (m *AttendeeBase) GetType()(*AttendeeType) {
    return m.type_escaped
}
// Serialize serializes information the current object
func (m *AttendeeBase) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Recipient.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetType() != nil {
        cast := (*m.GetType()).String()
        err = writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetType sets the type property value. The type of attendee. The possible values are: required, optional, resource. Currently if the attendee is a person, findMeetingTimes always considers the person is of the Required type.
func (m *AttendeeBase) SetType(value *AttendeeType)() {
    m.type_escaped = value
}
