package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Reminder 
type Reminder struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Identifies the version of the reminder. Every time the reminder is changed, changeKey changes as well. This allows Exchange to apply changes to the correct version of the object.
    changeKey *string
    // The date, time and time zone that the event ends.
    eventEndTime DateTimeTimeZoneable
    // The unique ID of the event. Read only.
    eventId *string
    // The location of the event.
    eventLocation Locationable
    // The date, time, and time zone that the event starts.
    eventStartTime DateTimeTimeZoneable
    // The text of the event's subject line.
    eventSubject *string
    // The URL to open the event in Outlook on the web.The event will open in the browser if you are logged in to your mailbox via Outlook on the web. You will be prompted to login if you are not already logged in with the browser.This URL cannot be accessed from within an iFrame.
    eventWebLink *string
    // The OdataType property
    odataType *string
    // The date, time, and time zone that the reminder is set to occur.
    reminderFireTime DateTimeTimeZoneable
}
// NewReminder instantiates a new reminder and sets the default values.
func NewReminder()(*Reminder) {
    m := &Reminder{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateReminderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateReminderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewReminder(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Reminder) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetChangeKey gets the changeKey property value. Identifies the version of the reminder. Every time the reminder is changed, changeKey changes as well. This allows Exchange to apply changes to the correct version of the object.
func (m *Reminder) GetChangeKey()(*string) {
    return m.changeKey
}
// GetEventEndTime gets the eventEndTime property value. The date, time and time zone that the event ends.
func (m *Reminder) GetEventEndTime()(DateTimeTimeZoneable) {
    return m.eventEndTime
}
// GetEventId gets the eventId property value. The unique ID of the event. Read only.
func (m *Reminder) GetEventId()(*string) {
    return m.eventId
}
// GetEventLocation gets the eventLocation property value. The location of the event.
func (m *Reminder) GetEventLocation()(Locationable) {
    return m.eventLocation
}
// GetEventStartTime gets the eventStartTime property value. The date, time, and time zone that the event starts.
func (m *Reminder) GetEventStartTime()(DateTimeTimeZoneable) {
    return m.eventStartTime
}
// GetEventSubject gets the eventSubject property value. The text of the event's subject line.
func (m *Reminder) GetEventSubject()(*string) {
    return m.eventSubject
}
// GetEventWebLink gets the eventWebLink property value. The URL to open the event in Outlook on the web.The event will open in the browser if you are logged in to your mailbox via Outlook on the web. You will be prompted to login if you are not already logged in with the browser.This URL cannot be accessed from within an iFrame.
func (m *Reminder) GetEventWebLink()(*string) {
    return m.eventWebLink
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Reminder) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["changeKey"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetChangeKey)
    res["eventEndTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue , m.SetEventEndTime)
    res["eventId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEventId)
    res["eventLocation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateLocationFromDiscriminatorValue , m.SetEventLocation)
    res["eventStartTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue , m.SetEventStartTime)
    res["eventSubject"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEventSubject)
    res["eventWebLink"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEventWebLink)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["reminderFireTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue , m.SetReminderFireTime)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Reminder) GetOdataType()(*string) {
    return m.odataType
}
// GetReminderFireTime gets the reminderFireTime property value. The date, time, and time zone that the reminder is set to occur.
func (m *Reminder) GetReminderFireTime()(DateTimeTimeZoneable) {
    return m.reminderFireTime
}
// Serialize serializes information the current object
func (m *Reminder) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("changeKey", m.GetChangeKey())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("eventEndTime", m.GetEventEndTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("eventId", m.GetEventId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("eventLocation", m.GetEventLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("eventStartTime", m.GetEventStartTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("eventSubject", m.GetEventSubject())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("eventWebLink", m.GetEventWebLink())
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
        err := writer.WriteObjectValue("reminderFireTime", m.GetReminderFireTime())
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
func (m *Reminder) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetChangeKey sets the changeKey property value. Identifies the version of the reminder. Every time the reminder is changed, changeKey changes as well. This allows Exchange to apply changes to the correct version of the object.
func (m *Reminder) SetChangeKey(value *string)() {
    m.changeKey = value
}
// SetEventEndTime sets the eventEndTime property value. The date, time and time zone that the event ends.
func (m *Reminder) SetEventEndTime(value DateTimeTimeZoneable)() {
    m.eventEndTime = value
}
// SetEventId sets the eventId property value. The unique ID of the event. Read only.
func (m *Reminder) SetEventId(value *string)() {
    m.eventId = value
}
// SetEventLocation sets the eventLocation property value. The location of the event.
func (m *Reminder) SetEventLocation(value Locationable)() {
    m.eventLocation = value
}
// SetEventStartTime sets the eventStartTime property value. The date, time, and time zone that the event starts.
func (m *Reminder) SetEventStartTime(value DateTimeTimeZoneable)() {
    m.eventStartTime = value
}
// SetEventSubject sets the eventSubject property value. The text of the event's subject line.
func (m *Reminder) SetEventSubject(value *string)() {
    m.eventSubject = value
}
// SetEventWebLink sets the eventWebLink property value. The URL to open the event in Outlook on the web.The event will open in the browser if you are logged in to your mailbox via Outlook on the web. You will be prompted to login if you are not already logged in with the browser.This URL cannot be accessed from within an iFrame.
func (m *Reminder) SetEventWebLink(value *string)() {
    m.eventWebLink = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Reminder) SetOdataType(value *string)() {
    m.odataType = value
}
// SetReminderFireTime sets the reminderFireTime property value. The date, time, and time zone that the reminder is set to occur.
func (m *Reminder) SetReminderFireTime(value DateTimeTimeZoneable)() {
    m.reminderFireTime = value
}
