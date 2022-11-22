package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CalendarSharingMessage 
type CalendarSharingMessage struct {
    Message
    // The canAccept property
    canAccept *bool
    // The sharingMessageAction property
    sharingMessageAction CalendarSharingMessageActionable
    // The sharingMessageActions property
    sharingMessageActions []CalendarSharingMessageActionable
    // The suggestedCalendarName property
    suggestedCalendarName *string
}
// NewCalendarSharingMessage instantiates a new CalendarSharingMessage and sets the default values.
func NewCalendarSharingMessage()(*CalendarSharingMessage) {
    m := &CalendarSharingMessage{
        Message: *NewMessage(),
    }
    odataTypeValue := "#microsoft.graph.calendarSharingMessage";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateCalendarSharingMessageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCalendarSharingMessageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCalendarSharingMessage(), nil
}
// GetCanAccept gets the canAccept property value. The canAccept property
func (m *CalendarSharingMessage) GetCanAccept()(*bool) {
    return m.canAccept
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CalendarSharingMessage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Message.GetFieldDeserializers()
    res["canAccept"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetCanAccept)
    res["sharingMessageAction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCalendarSharingMessageActionFromDiscriminatorValue , m.SetSharingMessageAction)
    res["sharingMessageActions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCalendarSharingMessageActionFromDiscriminatorValue , m.SetSharingMessageActions)
    res["suggestedCalendarName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSuggestedCalendarName)
    return res
}
// GetSharingMessageAction gets the sharingMessageAction property value. The sharingMessageAction property
func (m *CalendarSharingMessage) GetSharingMessageAction()(CalendarSharingMessageActionable) {
    return m.sharingMessageAction
}
// GetSharingMessageActions gets the sharingMessageActions property value. The sharingMessageActions property
func (m *CalendarSharingMessage) GetSharingMessageActions()([]CalendarSharingMessageActionable) {
    return m.sharingMessageActions
}
// GetSuggestedCalendarName gets the suggestedCalendarName property value. The suggestedCalendarName property
func (m *CalendarSharingMessage) GetSuggestedCalendarName()(*string) {
    return m.suggestedCalendarName
}
// Serialize serializes information the current object
func (m *CalendarSharingMessage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Message.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("canAccept", m.GetCanAccept())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharingMessageAction", m.GetSharingMessageAction())
        if err != nil {
            return err
        }
    }
    if m.GetSharingMessageActions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSharingMessageActions())
        err = writer.WriteCollectionOfObjectValues("sharingMessageActions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("suggestedCalendarName", m.GetSuggestedCalendarName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCanAccept sets the canAccept property value. The canAccept property
func (m *CalendarSharingMessage) SetCanAccept(value *bool)() {
    m.canAccept = value
}
// SetSharingMessageAction sets the sharingMessageAction property value. The sharingMessageAction property
func (m *CalendarSharingMessage) SetSharingMessageAction(value CalendarSharingMessageActionable)() {
    m.sharingMessageAction = value
}
// SetSharingMessageActions sets the sharingMessageActions property value. The sharingMessageActions property
func (m *CalendarSharingMessage) SetSharingMessageActions(value []CalendarSharingMessageActionable)() {
    m.sharingMessageActions = value
}
// SetSuggestedCalendarName sets the suggestedCalendarName property value. The suggestedCalendarName property
func (m *CalendarSharingMessage) SetSuggestedCalendarName(value *string)() {
    m.suggestedCalendarName = value
}
