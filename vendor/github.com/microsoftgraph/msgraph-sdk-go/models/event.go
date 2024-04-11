package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Event 
type Event struct {
    OutlookItem
    // The type property
    TypeEscaped *EventType
}
// NewEvent instantiates a new Event and sets the default values.
func NewEvent()(*Event) {
    m := &Event{
        OutlookItem: *NewOutlookItem(),
    }
    odataTypeValue := "#microsoft.graph.event"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEventFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEventFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEvent(), nil
}
// GetAllowNewTimeProposals gets the allowNewTimeProposals property value. true if the meeting organizer allows invitees to propose a new time when responding; otherwise, false. Optional. Default is true.
func (m *Event) GetAllowNewTimeProposals()(*bool) {
    val, err := m.GetBackingStore().Get("allowNewTimeProposals")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAttachments gets the attachments property value. The collection of FileAttachment, ItemAttachment, and referenceAttachment attachments for the event. Navigation property. Read-only. Nullable.
func (m *Event) GetAttachments()([]Attachmentable) {
    val, err := m.GetBackingStore().Get("attachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Attachmentable)
    }
    return nil
}
// GetAttendees gets the attendees property value. The collection of attendees for the event.
func (m *Event) GetAttendees()([]Attendeeable) {
    val, err := m.GetBackingStore().Get("attendees")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Attendeeable)
    }
    return nil
}
// GetBody gets the body property value. The body of the message associated with the event. It can be in HTML or text format.
func (m *Event) GetBody()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("body")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// GetBodyPreview gets the bodyPreview property value. The preview of the message associated with the event. It is in text format.
func (m *Event) GetBodyPreview()(*string) {
    val, err := m.GetBackingStore().Get("bodyPreview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCalendar gets the calendar property value. The calendar that contains the event. Navigation property. Read-only.
func (m *Event) GetCalendar()(Calendarable) {
    val, err := m.GetBackingStore().Get("calendar")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Calendarable)
    }
    return nil
}
// GetEnd gets the end property value. The date, time, and time zone that the event ends. By default, the end time is in UTC.
func (m *Event) GetEnd()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("end")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetExtensions gets the extensions property value. The collection of open extensions defined for the event. Nullable.
func (m *Event) GetExtensions()([]Extensionable) {
    val, err := m.GetBackingStore().Get("extensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Extensionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Event) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OutlookItem.GetFieldDeserializers()
    res["allowNewTimeProposals"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowNewTimeProposals(val)
        }
        return nil
    }
    res["attachments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttachmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Attachmentable, len(val))
            for i, v := range val {
                res[i] = v.(Attachmentable)
            }
            m.SetAttachments(res)
        }
        return nil
    }
    res["attendees"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttendeeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Attendeeable, len(val))
            for i, v := range val {
                res[i] = v.(Attendeeable)
            }
            m.SetAttendees(res)
        }
        return nil
    }
    res["body"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemBodyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBody(val.(ItemBodyable))
        }
        return nil
    }
    res["bodyPreview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBodyPreview(val)
        }
        return nil
    }
    res["calendar"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCalendarFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCalendar(val.(Calendarable))
        }
        return nil
    }
    res["end"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnd(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["extensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExtensionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Extensionable, len(val))
            for i, v := range val {
                res[i] = v.(Extensionable)
            }
            m.SetExtensions(res)
        }
        return nil
    }
    res["hasAttachments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasAttachments(val)
        }
        return nil
    }
    res["hideAttendees"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHideAttendees(val)
        }
        return nil
    }
    res["iCalUId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICalUId(val)
        }
        return nil
    }
    res["importance"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseImportance)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImportance(val.(*Importance))
        }
        return nil
    }
    res["instances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Eventable, len(val))
            for i, v := range val {
                res[i] = v.(Eventable)
            }
            m.SetInstances(res)
        }
        return nil
    }
    res["isAllDay"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAllDay(val)
        }
        return nil
    }
    res["isCancelled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsCancelled(val)
        }
        return nil
    }
    res["isDraft"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDraft(val)
        }
        return nil
    }
    res["isOnlineMeeting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsOnlineMeeting(val)
        }
        return nil
    }
    res["isOrganizer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsOrganizer(val)
        }
        return nil
    }
    res["isReminderOn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsReminderOn(val)
        }
        return nil
    }
    res["location"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocation(val.(Locationable))
        }
        return nil
    }
    res["locations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Locationable, len(val))
            for i, v := range val {
                res[i] = v.(Locationable)
            }
            m.SetLocations(res)
        }
        return nil
    }
    res["multiValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMultiValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MultiValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                res[i] = v.(MultiValueLegacyExtendedPropertyable)
            }
            m.SetMultiValueExtendedProperties(res)
        }
        return nil
    }
    res["onlineMeeting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnlineMeetingInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnlineMeeting(val.(OnlineMeetingInfoable))
        }
        return nil
    }
    res["onlineMeetingProvider"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOnlineMeetingProviderType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnlineMeetingProvider(val.(*OnlineMeetingProviderType))
        }
        return nil
    }
    res["onlineMeetingUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOnlineMeetingUrl(val)
        }
        return nil
    }
    res["organizer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRecipientFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrganizer(val.(Recipientable))
        }
        return nil
    }
    res["originalEndTimeZone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOriginalEndTimeZone(val)
        }
        return nil
    }
    res["originalStart"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOriginalStart(val)
        }
        return nil
    }
    res["originalStartTimeZone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOriginalStartTimeZone(val)
        }
        return nil
    }
    res["recurrence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePatternedRecurrenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecurrence(val.(PatternedRecurrenceable))
        }
        return nil
    }
    res["reminderMinutesBeforeStart"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReminderMinutesBeforeStart(val)
        }
        return nil
    }
    res["responseRequested"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResponseRequested(val)
        }
        return nil
    }
    res["responseStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResponseStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResponseStatus(val.(ResponseStatusable))
        }
        return nil
    }
    res["sensitivity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSensitivity)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSensitivity(val.(*Sensitivity))
        }
        return nil
    }
    res["seriesMasterId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSeriesMasterId(val)
        }
        return nil
    }
    res["showAs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFreeBusyStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowAs(val.(*FreeBusyStatus))
        }
        return nil
    }
    res["singleValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SingleValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                res[i] = v.(SingleValueLegacyExtendedPropertyable)
            }
            m.SetSingleValueExtendedProperties(res)
        }
        return nil
    }
    res["start"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStart(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["subject"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubject(val)
        }
        return nil
    }
    res["transactionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTransactionId(val)
        }
        return nil
    }
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEventType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetType(val.(*EventType))
        }
        return nil
    }
    res["webLink"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebLink(val)
        }
        return nil
    }
    return res
}
// GetHasAttachments gets the hasAttachments property value. Set to true if the event has attachments.
func (m *Event) GetHasAttachments()(*bool) {
    val, err := m.GetBackingStore().Get("hasAttachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHideAttendees gets the hideAttendees property value. When set to true, each attendee only sees themselves in the meeting request and meeting Tracking list. Default is false.
func (m *Event) GetHideAttendees()(*bool) {
    val, err := m.GetBackingStore().Get("hideAttendees")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICalUId gets the iCalUId property value. A unique identifier for an event across calendars. This ID is different for each occurrence in a recurring series. Read-only.
func (m *Event) GetICalUId()(*string) {
    val, err := m.GetBackingStore().Get("iCalUId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetImportance gets the importance property value. The importance property
func (m *Event) GetImportance()(*Importance) {
    val, err := m.GetBackingStore().Get("importance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Importance)
    }
    return nil
}
// GetInstances gets the instances property value. The occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
func (m *Event) GetInstances()([]Eventable) {
    val, err := m.GetBackingStore().Get("instances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Eventable)
    }
    return nil
}
// GetIsAllDay gets the isAllDay property value. The isAllDay property
func (m *Event) GetIsAllDay()(*bool) {
    val, err := m.GetBackingStore().Get("isAllDay")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsCancelled gets the isCancelled property value. The isCancelled property
func (m *Event) GetIsCancelled()(*bool) {
    val, err := m.GetBackingStore().Get("isCancelled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsDraft gets the isDraft property value. The isDraft property
func (m *Event) GetIsDraft()(*bool) {
    val, err := m.GetBackingStore().Get("isDraft")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsOnlineMeeting gets the isOnlineMeeting property value. The isOnlineMeeting property
func (m *Event) GetIsOnlineMeeting()(*bool) {
    val, err := m.GetBackingStore().Get("isOnlineMeeting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsOrganizer gets the isOrganizer property value. The isOrganizer property
func (m *Event) GetIsOrganizer()(*bool) {
    val, err := m.GetBackingStore().Get("isOrganizer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsReminderOn gets the isReminderOn property value. The isReminderOn property
func (m *Event) GetIsReminderOn()(*bool) {
    val, err := m.GetBackingStore().Get("isReminderOn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLocation gets the location property value. The location property
func (m *Event) GetLocation()(Locationable) {
    val, err := m.GetBackingStore().Get("location")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Locationable)
    }
    return nil
}
// GetLocations gets the locations property value. The locations property
func (m *Event) GetLocations()([]Locationable) {
    val, err := m.GetBackingStore().Get("locations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Locationable)
    }
    return nil
}
// GetMultiValueExtendedProperties gets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("multiValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MultiValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetOnlineMeeting gets the onlineMeeting property value. The onlineMeeting property
func (m *Event) GetOnlineMeeting()(OnlineMeetingInfoable) {
    val, err := m.GetBackingStore().Get("onlineMeeting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnlineMeetingInfoable)
    }
    return nil
}
// GetOnlineMeetingProvider gets the onlineMeetingProvider property value. The onlineMeetingProvider property
func (m *Event) GetOnlineMeetingProvider()(*OnlineMeetingProviderType) {
    val, err := m.GetBackingStore().Get("onlineMeetingProvider")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OnlineMeetingProviderType)
    }
    return nil
}
// GetOnlineMeetingUrl gets the onlineMeetingUrl property value. The onlineMeetingUrl property
func (m *Event) GetOnlineMeetingUrl()(*string) {
    val, err := m.GetBackingStore().Get("onlineMeetingUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrganizer gets the organizer property value. The organizer property
func (m *Event) GetOrganizer()(Recipientable) {
    val, err := m.GetBackingStore().Get("organizer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Recipientable)
    }
    return nil
}
// GetOriginalEndTimeZone gets the originalEndTimeZone property value. The originalEndTimeZone property
func (m *Event) GetOriginalEndTimeZone()(*string) {
    val, err := m.GetBackingStore().Get("originalEndTimeZone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOriginalStart gets the originalStart property value. The originalStart property
func (m *Event) GetOriginalStart()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("originalStart")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetOriginalStartTimeZone gets the originalStartTimeZone property value. The originalStartTimeZone property
func (m *Event) GetOriginalStartTimeZone()(*string) {
    val, err := m.GetBackingStore().Get("originalStartTimeZone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecurrence gets the recurrence property value. The recurrence property
func (m *Event) GetRecurrence()(PatternedRecurrenceable) {
    val, err := m.GetBackingStore().Get("recurrence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PatternedRecurrenceable)
    }
    return nil
}
// GetReminderMinutesBeforeStart gets the reminderMinutesBeforeStart property value. The reminderMinutesBeforeStart property
func (m *Event) GetReminderMinutesBeforeStart()(*int32) {
    val, err := m.GetBackingStore().Get("reminderMinutesBeforeStart")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetResponseRequested gets the responseRequested property value. The responseRequested property
func (m *Event) GetResponseRequested()(*bool) {
    val, err := m.GetBackingStore().Get("responseRequested")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetResponseStatus gets the responseStatus property value. The responseStatus property
func (m *Event) GetResponseStatus()(ResponseStatusable) {
    val, err := m.GetBackingStore().Get("responseStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResponseStatusable)
    }
    return nil
}
// GetSensitivity gets the sensitivity property value. The sensitivity property
func (m *Event) GetSensitivity()(*Sensitivity) {
    val, err := m.GetBackingStore().Get("sensitivity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Sensitivity)
    }
    return nil
}
// GetSeriesMasterId gets the seriesMasterId property value. The seriesMasterId property
func (m *Event) GetSeriesMasterId()(*string) {
    val, err := m.GetBackingStore().Get("seriesMasterId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetShowAs gets the showAs property value. The showAs property
func (m *Event) GetShowAs()(*FreeBusyStatus) {
    val, err := m.GetBackingStore().Get("showAs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FreeBusyStatus)
    }
    return nil
}
// GetSingleValueExtendedProperties gets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("singleValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SingleValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetStart gets the start property value. The start property
func (m *Event) GetStart()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("start")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetSubject gets the subject property value. The subject property
func (m *Event) GetSubject()(*string) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTransactionId gets the transactionId property value. The transactionId property
func (m *Event) GetTransactionId()(*string) {
    val, err := m.GetBackingStore().Get("transactionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetType gets the type property value. The type property
func (m *Event) GetType()(*EventType) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EventType)
    }
    return nil
}
// GetWebLink gets the webLink property value. The webLink property
func (m *Event) GetWebLink()(*string) {
    val, err := m.GetBackingStore().Get("webLink")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Event) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OutlookItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowNewTimeProposals", m.GetAllowNewTimeProposals())
        if err != nil {
            return err
        }
    }
    if m.GetAttachments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttachments()))
        for i, v := range m.GetAttachments() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("attachments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAttendees() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttendees()))
        for i, v := range m.GetAttendees() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("attendees", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("body", m.GetBody())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("bodyPreview", m.GetBodyPreview())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("calendar", m.GetCalendar())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("end", m.GetEnd())
        if err != nil {
            return err
        }
    }
    if m.GetExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExtensions()))
        for i, v := range m.GetExtensions() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("extensions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasAttachments", m.GetHasAttachments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hideAttendees", m.GetHideAttendees())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("iCalUId", m.GetICalUId())
        if err != nil {
            return err
        }
    }
    if m.GetImportance() != nil {
        cast := (*m.GetImportance()).String()
        err = writer.WriteStringValue("importance", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetInstances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInstances()))
        for i, v := range m.GetInstances() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("instances", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAllDay", m.GetIsAllDay())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isCancelled", m.GetIsCancelled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDraft", m.GetIsDraft())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isOnlineMeeting", m.GetIsOnlineMeeting())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isOrganizer", m.GetIsOrganizer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isReminderOn", m.GetIsReminderOn())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("location", m.GetLocation())
        if err != nil {
            return err
        }
    }
    if m.GetLocations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLocations()))
        for i, v := range m.GetLocations() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("locations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMultiValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMultiValueExtendedProperties()))
        for i, v := range m.GetMultiValueExtendedProperties() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("multiValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("onlineMeeting", m.GetOnlineMeeting())
        if err != nil {
            return err
        }
    }
    if m.GetOnlineMeetingProvider() != nil {
        cast := (*m.GetOnlineMeetingProvider()).String()
        err = writer.WriteStringValue("onlineMeetingProvider", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("onlineMeetingUrl", m.GetOnlineMeetingUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("organizer", m.GetOrganizer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("originalEndTimeZone", m.GetOriginalEndTimeZone())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("originalStart", m.GetOriginalStart())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("originalStartTimeZone", m.GetOriginalStartTimeZone())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("recurrence", m.GetRecurrence())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("reminderMinutesBeforeStart", m.GetReminderMinutesBeforeStart())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("responseRequested", m.GetResponseRequested())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("responseStatus", m.GetResponseStatus())
        if err != nil {
            return err
        }
    }
    if m.GetSensitivity() != nil {
        cast := (*m.GetSensitivity()).String()
        err = writer.WriteStringValue("sensitivity", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("seriesMasterId", m.GetSeriesMasterId())
        if err != nil {
            return err
        }
    }
    if m.GetShowAs() != nil {
        cast := (*m.GetShowAs()).String()
        err = writer.WriteStringValue("showAs", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetSingleValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSingleValueExtendedProperties()))
        for i, v := range m.GetSingleValueExtendedProperties() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("singleValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("start", m.GetStart())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("subject", m.GetSubject())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("transactionId", m.GetTransactionId())
        if err != nil {
            return err
        }
    }
    if m.GetType() != nil {
        cast := (*m.GetType()).String()
        err = writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webLink", m.GetWebLink())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowNewTimeProposals sets the allowNewTimeProposals property value. true if the meeting organizer allows invitees to propose a new time when responding; otherwise, false. Optional. Default is true.
func (m *Event) SetAllowNewTimeProposals(value *bool)() {
    err := m.GetBackingStore().Set("allowNewTimeProposals", value)
    if err != nil {
        panic(err)
    }
}
// SetAttachments sets the attachments property value. The collection of FileAttachment, ItemAttachment, and referenceAttachment attachments for the event. Navigation property. Read-only. Nullable.
func (m *Event) SetAttachments(value []Attachmentable)() {
    err := m.GetBackingStore().Set("attachments", value)
    if err != nil {
        panic(err)
    }
}
// SetAttendees sets the attendees property value. The collection of attendees for the event.
func (m *Event) SetAttendees(value []Attendeeable)() {
    err := m.GetBackingStore().Set("attendees", value)
    if err != nil {
        panic(err)
    }
}
// SetBody sets the body property value. The body of the message associated with the event. It can be in HTML or text format.
func (m *Event) SetBody(value ItemBodyable)() {
    err := m.GetBackingStore().Set("body", value)
    if err != nil {
        panic(err)
    }
}
// SetBodyPreview sets the bodyPreview property value. The preview of the message associated with the event. It is in text format.
func (m *Event) SetBodyPreview(value *string)() {
    err := m.GetBackingStore().Set("bodyPreview", value)
    if err != nil {
        panic(err)
    }
}
// SetCalendar sets the calendar property value. The calendar that contains the event. Navigation property. Read-only.
func (m *Event) SetCalendar(value Calendarable)() {
    err := m.GetBackingStore().Set("calendar", value)
    if err != nil {
        panic(err)
    }
}
// SetEnd sets the end property value. The date, time, and time zone that the event ends. By default, the end time is in UTC.
func (m *Event) SetEnd(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("end", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensions sets the extensions property value. The collection of open extensions defined for the event. Nullable.
func (m *Event) SetExtensions(value []Extensionable)() {
    err := m.GetBackingStore().Set("extensions", value)
    if err != nil {
        panic(err)
    }
}
// SetHasAttachments sets the hasAttachments property value. Set to true if the event has attachments.
func (m *Event) SetHasAttachments(value *bool)() {
    err := m.GetBackingStore().Set("hasAttachments", value)
    if err != nil {
        panic(err)
    }
}
// SetHideAttendees sets the hideAttendees property value. When set to true, each attendee only sees themselves in the meeting request and meeting Tracking list. Default is false.
func (m *Event) SetHideAttendees(value *bool)() {
    err := m.GetBackingStore().Set("hideAttendees", value)
    if err != nil {
        panic(err)
    }
}
// SetICalUId sets the iCalUId property value. A unique identifier for an event across calendars. This ID is different for each occurrence in a recurring series. Read-only.
func (m *Event) SetICalUId(value *string)() {
    err := m.GetBackingStore().Set("iCalUId", value)
    if err != nil {
        panic(err)
    }
}
// SetImportance sets the importance property value. The importance property
func (m *Event) SetImportance(value *Importance)() {
    err := m.GetBackingStore().Set("importance", value)
    if err != nil {
        panic(err)
    }
}
// SetInstances sets the instances property value. The occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
func (m *Event) SetInstances(value []Eventable)() {
    err := m.GetBackingStore().Set("instances", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAllDay sets the isAllDay property value. The isAllDay property
func (m *Event) SetIsAllDay(value *bool)() {
    err := m.GetBackingStore().Set("isAllDay", value)
    if err != nil {
        panic(err)
    }
}
// SetIsCancelled sets the isCancelled property value. The isCancelled property
func (m *Event) SetIsCancelled(value *bool)() {
    err := m.GetBackingStore().Set("isCancelled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDraft sets the isDraft property value. The isDraft property
func (m *Event) SetIsDraft(value *bool)() {
    err := m.GetBackingStore().Set("isDraft", value)
    if err != nil {
        panic(err)
    }
}
// SetIsOnlineMeeting sets the isOnlineMeeting property value. The isOnlineMeeting property
func (m *Event) SetIsOnlineMeeting(value *bool)() {
    err := m.GetBackingStore().Set("isOnlineMeeting", value)
    if err != nil {
        panic(err)
    }
}
// SetIsOrganizer sets the isOrganizer property value. The isOrganizer property
func (m *Event) SetIsOrganizer(value *bool)() {
    err := m.GetBackingStore().Set("isOrganizer", value)
    if err != nil {
        panic(err)
    }
}
// SetIsReminderOn sets the isReminderOn property value. The isReminderOn property
func (m *Event) SetIsReminderOn(value *bool)() {
    err := m.GetBackingStore().Set("isReminderOn", value)
    if err != nil {
        panic(err)
    }
}
// SetLocation sets the location property value. The location property
func (m *Event) SetLocation(value Locationable)() {
    err := m.GetBackingStore().Set("location", value)
    if err != nil {
        panic(err)
    }
}
// SetLocations sets the locations property value. The locations property
func (m *Event) SetLocations(value []Locationable)() {
    err := m.GetBackingStore().Set("locations", value)
    if err != nil {
        panic(err)
    }
}
// SetMultiValueExtendedProperties sets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("multiValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetOnlineMeeting sets the onlineMeeting property value. The onlineMeeting property
func (m *Event) SetOnlineMeeting(value OnlineMeetingInfoable)() {
    err := m.GetBackingStore().Set("onlineMeeting", value)
    if err != nil {
        panic(err)
    }
}
// SetOnlineMeetingProvider sets the onlineMeetingProvider property value. The onlineMeetingProvider property
func (m *Event) SetOnlineMeetingProvider(value *OnlineMeetingProviderType)() {
    err := m.GetBackingStore().Set("onlineMeetingProvider", value)
    if err != nil {
        panic(err)
    }
}
// SetOnlineMeetingUrl sets the onlineMeetingUrl property value. The onlineMeetingUrl property
func (m *Event) SetOnlineMeetingUrl(value *string)() {
    err := m.GetBackingStore().Set("onlineMeetingUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetOrganizer sets the organizer property value. The organizer property
func (m *Event) SetOrganizer(value Recipientable)() {
    err := m.GetBackingStore().Set("organizer", value)
    if err != nil {
        panic(err)
    }
}
// SetOriginalEndTimeZone sets the originalEndTimeZone property value. The originalEndTimeZone property
func (m *Event) SetOriginalEndTimeZone(value *string)() {
    err := m.GetBackingStore().Set("originalEndTimeZone", value)
    if err != nil {
        panic(err)
    }
}
// SetOriginalStart sets the originalStart property value. The originalStart property
func (m *Event) SetOriginalStart(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("originalStart", value)
    if err != nil {
        panic(err)
    }
}
// SetOriginalStartTimeZone sets the originalStartTimeZone property value. The originalStartTimeZone property
func (m *Event) SetOriginalStartTimeZone(value *string)() {
    err := m.GetBackingStore().Set("originalStartTimeZone", value)
    if err != nil {
        panic(err)
    }
}
// SetRecurrence sets the recurrence property value. The recurrence property
func (m *Event) SetRecurrence(value PatternedRecurrenceable)() {
    err := m.GetBackingStore().Set("recurrence", value)
    if err != nil {
        panic(err)
    }
}
// SetReminderMinutesBeforeStart sets the reminderMinutesBeforeStart property value. The reminderMinutesBeforeStart property
func (m *Event) SetReminderMinutesBeforeStart(value *int32)() {
    err := m.GetBackingStore().Set("reminderMinutesBeforeStart", value)
    if err != nil {
        panic(err)
    }
}
// SetResponseRequested sets the responseRequested property value. The responseRequested property
func (m *Event) SetResponseRequested(value *bool)() {
    err := m.GetBackingStore().Set("responseRequested", value)
    if err != nil {
        panic(err)
    }
}
// SetResponseStatus sets the responseStatus property value. The responseStatus property
func (m *Event) SetResponseStatus(value ResponseStatusable)() {
    err := m.GetBackingStore().Set("responseStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetSensitivity sets the sensitivity property value. The sensitivity property
func (m *Event) SetSensitivity(value *Sensitivity)() {
    err := m.GetBackingStore().Set("sensitivity", value)
    if err != nil {
        panic(err)
    }
}
// SetSeriesMasterId sets the seriesMasterId property value. The seriesMasterId property
func (m *Event) SetSeriesMasterId(value *string)() {
    err := m.GetBackingStore().Set("seriesMasterId", value)
    if err != nil {
        panic(err)
    }
}
// SetShowAs sets the showAs property value. The showAs property
func (m *Event) SetShowAs(value *FreeBusyStatus)() {
    err := m.GetBackingStore().Set("showAs", value)
    if err != nil {
        panic(err)
    }
}
// SetSingleValueExtendedProperties sets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("singleValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetStart sets the start property value. The start property
func (m *Event) SetStart(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("start", value)
    if err != nil {
        panic(err)
    }
}
// SetSubject sets the subject property value. The subject property
func (m *Event) SetSubject(value *string)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
// SetTransactionId sets the transactionId property value. The transactionId property
func (m *Event) SetTransactionId(value *string)() {
    err := m.GetBackingStore().Set("transactionId", value)
    if err != nil {
        panic(err)
    }
}
// SetType sets the type property value. The type property
func (m *Event) SetType(value *EventType)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
// SetWebLink sets the webLink property value. The webLink property
func (m *Event) SetWebLink(value *string)() {
    err := m.GetBackingStore().Set("webLink", value)
    if err != nil {
        panic(err)
    }
}
// Eventable 
type Eventable interface {
    OutlookItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowNewTimeProposals()(*bool)
    GetAttachments()([]Attachmentable)
    GetAttendees()([]Attendeeable)
    GetBody()(ItemBodyable)
    GetBodyPreview()(*string)
    GetCalendar()(Calendarable)
    GetEnd()(DateTimeTimeZoneable)
    GetExtensions()([]Extensionable)
    GetHasAttachments()(*bool)
    GetHideAttendees()(*bool)
    GetICalUId()(*string)
    GetImportance()(*Importance)
    GetInstances()([]Eventable)
    GetIsAllDay()(*bool)
    GetIsCancelled()(*bool)
    GetIsDraft()(*bool)
    GetIsOnlineMeeting()(*bool)
    GetIsOrganizer()(*bool)
    GetIsReminderOn()(*bool)
    GetLocation()(Locationable)
    GetLocations()([]Locationable)
    GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable)
    GetOnlineMeeting()(OnlineMeetingInfoable)
    GetOnlineMeetingProvider()(*OnlineMeetingProviderType)
    GetOnlineMeetingUrl()(*string)
    GetOrganizer()(Recipientable)
    GetOriginalEndTimeZone()(*string)
    GetOriginalStart()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetOriginalStartTimeZone()(*string)
    GetRecurrence()(PatternedRecurrenceable)
    GetReminderMinutesBeforeStart()(*int32)
    GetResponseRequested()(*bool)
    GetResponseStatus()(ResponseStatusable)
    GetSensitivity()(*Sensitivity)
    GetSeriesMasterId()(*string)
    GetShowAs()(*FreeBusyStatus)
    GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable)
    GetStart()(DateTimeTimeZoneable)
    GetSubject()(*string)
    GetTransactionId()(*string)
    GetType()(*EventType)
    GetWebLink()(*string)
    SetAllowNewTimeProposals(value *bool)()
    SetAttachments(value []Attachmentable)()
    SetAttendees(value []Attendeeable)()
    SetBody(value ItemBodyable)()
    SetBodyPreview(value *string)()
    SetCalendar(value Calendarable)()
    SetEnd(value DateTimeTimeZoneable)()
    SetExtensions(value []Extensionable)()
    SetHasAttachments(value *bool)()
    SetHideAttendees(value *bool)()
    SetICalUId(value *string)()
    SetImportance(value *Importance)()
    SetInstances(value []Eventable)()
    SetIsAllDay(value *bool)()
    SetIsCancelled(value *bool)()
    SetIsDraft(value *bool)()
    SetIsOnlineMeeting(value *bool)()
    SetIsOrganizer(value *bool)()
    SetIsReminderOn(value *bool)()
    SetLocation(value Locationable)()
    SetLocations(value []Locationable)()
    SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)()
    SetOnlineMeeting(value OnlineMeetingInfoable)()
    SetOnlineMeetingProvider(value *OnlineMeetingProviderType)()
    SetOnlineMeetingUrl(value *string)()
    SetOrganizer(value Recipientable)()
    SetOriginalEndTimeZone(value *string)()
    SetOriginalStart(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetOriginalStartTimeZone(value *string)()
    SetRecurrence(value PatternedRecurrenceable)()
    SetReminderMinutesBeforeStart(value *int32)()
    SetResponseRequested(value *bool)()
    SetResponseStatus(value ResponseStatusable)()
    SetSensitivity(value *Sensitivity)()
    SetSeriesMasterId(value *string)()
    SetShowAs(value *FreeBusyStatus)()
    SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)()
    SetStart(value DateTimeTimeZoneable)()
    SetSubject(value *string)()
    SetTransactionId(value *string)()
    SetType(value *EventType)()
    SetWebLink(value *string)()
}
