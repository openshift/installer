package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Event 
type Event struct {
    OutlookItem
    // true if the meeting organizer allows invitees to propose a new time when responding; otherwise, false. Optional. Default is true.
    allowNewTimeProposals *bool
    // The collection of FileAttachment, ItemAttachment, and referenceAttachment attachments for the event. Navigation property. Read-only. Nullable.
    attachments []Attachmentable
    // The collection of attendees for the event.
    attendees []Attendeeable
    // The body of the message associated with the event. It can be in HTML or text format.
    body ItemBodyable
    // The preview of the message associated with the event. It is in text format.
    bodyPreview *string
    // The calendar that contains the event. Navigation property. Read-only.
    calendar Calendarable
    // The date, time, and time zone that the event ends. By default, the end time is in UTC.
    end DateTimeTimeZoneable
    // The collection of open extensions defined for the event. Nullable.
    extensions []Extensionable
    // Set to true if the event has attachments.
    hasAttachments *bool
    // When set to true, each attendee only sees themselves in the meeting request and meeting Tracking list. Default is false.
    hideAttendees *bool
    // A unique identifier for an event across calendars. This ID is different for each occurrence in a recurring series. Read-only.
    iCalUId *string
    // The importance property
    importance *Importance
    // The occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
    instances []Eventable
    // The isAllDay property
    isAllDay *bool
    // The isCancelled property
    isCancelled *bool
    // The isDraft property
    isDraft *bool
    // The isOnlineMeeting property
    isOnlineMeeting *bool
    // The isOrganizer property
    isOrganizer *bool
    // The isReminderOn property
    isReminderOn *bool
    // The location property
    location Locationable
    // The locations property
    locations []Locationable
    // The collection of multi-value extended properties defined for the event. Read-only. Nullable.
    multiValueExtendedProperties []MultiValueLegacyExtendedPropertyable
    // The onlineMeeting property
    onlineMeeting OnlineMeetingInfoable
    // The onlineMeetingProvider property
    onlineMeetingProvider *OnlineMeetingProviderType
    // The onlineMeetingUrl property
    onlineMeetingUrl *string
    // The organizer property
    organizer Recipientable
    // The originalEndTimeZone property
    originalEndTimeZone *string
    // The originalStart property
    originalStart *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The originalStartTimeZone property
    originalStartTimeZone *string
    // The recurrence property
    recurrence PatternedRecurrenceable
    // The reminderMinutesBeforeStart property
    reminderMinutesBeforeStart *int32
    // The responseRequested property
    responseRequested *bool
    // The responseStatus property
    responseStatus ResponseStatusable
    // The sensitivity property
    sensitivity *Sensitivity
    // The seriesMasterId property
    seriesMasterId *string
    // The showAs property
    showAs *FreeBusyStatus
    // The collection of single-value extended properties defined for the event. Read-only. Nullable.
    singleValueExtendedProperties []SingleValueLegacyExtendedPropertyable
    // The start property
    start DateTimeTimeZoneable
    // The subject property
    subject *string
    // The transactionId property
    transactionId *string
    // The type property
    type_escaped *EventType
    // The webLink property
    webLink *string
}
// NewEvent instantiates a new Event and sets the default values.
func NewEvent()(*Event) {
    m := &Event{
        OutlookItem: *NewOutlookItem(),
    }
    odataTypeValue := "#microsoft.graph.event";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEventFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEventFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEvent(), nil
}
// GetAllowNewTimeProposals gets the allowNewTimeProposals property value. true if the meeting organizer allows invitees to propose a new time when responding; otherwise, false. Optional. Default is true.
func (m *Event) GetAllowNewTimeProposals()(*bool) {
    return m.allowNewTimeProposals
}
// GetAttachments gets the attachments property value. The collection of FileAttachment, ItemAttachment, and referenceAttachment attachments for the event. Navigation property. Read-only. Nullable.
func (m *Event) GetAttachments()([]Attachmentable) {
    return m.attachments
}
// GetAttendees gets the attendees property value. The collection of attendees for the event.
func (m *Event) GetAttendees()([]Attendeeable) {
    return m.attendees
}
// GetBody gets the body property value. The body of the message associated with the event. It can be in HTML or text format.
func (m *Event) GetBody()(ItemBodyable) {
    return m.body
}
// GetBodyPreview gets the bodyPreview property value. The preview of the message associated with the event. It is in text format.
func (m *Event) GetBodyPreview()(*string) {
    return m.bodyPreview
}
// GetCalendar gets the calendar property value. The calendar that contains the event. Navigation property. Read-only.
func (m *Event) GetCalendar()(Calendarable) {
    return m.calendar
}
// GetEnd gets the end property value. The date, time, and time zone that the event ends. By default, the end time is in UTC.
func (m *Event) GetEnd()(DateTimeTimeZoneable) {
    return m.end
}
// GetExtensions gets the extensions property value. The collection of open extensions defined for the event. Nullable.
func (m *Event) GetExtensions()([]Extensionable) {
    return m.extensions
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Event) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OutlookItem.GetFieldDeserializers()
    res["allowNewTimeProposals"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowNewTimeProposals)
    res["attachments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAttachmentFromDiscriminatorValue , m.SetAttachments)
    res["attendees"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAttendeeFromDiscriminatorValue , m.SetAttendees)
    res["body"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateItemBodyFromDiscriminatorValue , m.SetBody)
    res["bodyPreview"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetBodyPreview)
    res["calendar"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCalendarFromDiscriminatorValue , m.SetCalendar)
    res["end"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue , m.SetEnd)
    res["extensions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateExtensionFromDiscriminatorValue , m.SetExtensions)
    res["hasAttachments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHasAttachments)
    res["hideAttendees"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHideAttendees)
    res["iCalUId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetICalUId)
    res["importance"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseImportance , m.SetImportance)
    res["instances"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEventFromDiscriminatorValue , m.SetInstances)
    res["isAllDay"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsAllDay)
    res["isCancelled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsCancelled)
    res["isDraft"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsDraft)
    res["isOnlineMeeting"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsOnlineMeeting)
    res["isOrganizer"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsOrganizer)
    res["isReminderOn"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsReminderOn)
    res["location"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateLocationFromDiscriminatorValue , m.SetLocation)
    res["locations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateLocationFromDiscriminatorValue , m.SetLocations)
    res["multiValueExtendedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMultiValueLegacyExtendedPropertyFromDiscriminatorValue , m.SetMultiValueExtendedProperties)
    res["onlineMeeting"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateOnlineMeetingInfoFromDiscriminatorValue , m.SetOnlineMeeting)
    res["onlineMeetingProvider"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseOnlineMeetingProviderType , m.SetOnlineMeetingProvider)
    res["onlineMeetingUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOnlineMeetingUrl)
    res["organizer"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRecipientFromDiscriminatorValue , m.SetOrganizer)
    res["originalEndTimeZone"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOriginalEndTimeZone)
    res["originalStart"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetOriginalStart)
    res["originalStartTimeZone"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOriginalStartTimeZone)
    res["recurrence"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePatternedRecurrenceFromDiscriminatorValue , m.SetRecurrence)
    res["reminderMinutesBeforeStart"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetReminderMinutesBeforeStart)
    res["responseRequested"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetResponseRequested)
    res["responseStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateResponseStatusFromDiscriminatorValue , m.SetResponseStatus)
    res["sensitivity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseSensitivity , m.SetSensitivity)
    res["seriesMasterId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSeriesMasterId)
    res["showAs"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseFreeBusyStatus , m.SetShowAs)
    res["singleValueExtendedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue , m.SetSingleValueExtendedProperties)
    res["start"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue , m.SetStart)
    res["subject"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSubject)
    res["transactionId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTransactionId)
    res["type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseEventType , m.SetType)
    res["webLink"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetWebLink)
    return res
}
// GetHasAttachments gets the hasAttachments property value. Set to true if the event has attachments.
func (m *Event) GetHasAttachments()(*bool) {
    return m.hasAttachments
}
// GetHideAttendees gets the hideAttendees property value. When set to true, each attendee only sees themselves in the meeting request and meeting Tracking list. Default is false.
func (m *Event) GetHideAttendees()(*bool) {
    return m.hideAttendees
}
// GetICalUId gets the iCalUId property value. A unique identifier for an event across calendars. This ID is different for each occurrence in a recurring series. Read-only.
func (m *Event) GetICalUId()(*string) {
    return m.iCalUId
}
// GetImportance gets the importance property value. The importance property
func (m *Event) GetImportance()(*Importance) {
    return m.importance
}
// GetInstances gets the instances property value. The occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
func (m *Event) GetInstances()([]Eventable) {
    return m.instances
}
// GetIsAllDay gets the isAllDay property value. The isAllDay property
func (m *Event) GetIsAllDay()(*bool) {
    return m.isAllDay
}
// GetIsCancelled gets the isCancelled property value. The isCancelled property
func (m *Event) GetIsCancelled()(*bool) {
    return m.isCancelled
}
// GetIsDraft gets the isDraft property value. The isDraft property
func (m *Event) GetIsDraft()(*bool) {
    return m.isDraft
}
// GetIsOnlineMeeting gets the isOnlineMeeting property value. The isOnlineMeeting property
func (m *Event) GetIsOnlineMeeting()(*bool) {
    return m.isOnlineMeeting
}
// GetIsOrganizer gets the isOrganizer property value. The isOrganizer property
func (m *Event) GetIsOrganizer()(*bool) {
    return m.isOrganizer
}
// GetIsReminderOn gets the isReminderOn property value. The isReminderOn property
func (m *Event) GetIsReminderOn()(*bool) {
    return m.isReminderOn
}
// GetLocation gets the location property value. The location property
func (m *Event) GetLocation()(Locationable) {
    return m.location
}
// GetLocations gets the locations property value. The locations property
func (m *Event) GetLocations()([]Locationable) {
    return m.locations
}
// GetMultiValueExtendedProperties gets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable) {
    return m.multiValueExtendedProperties
}
// GetOnlineMeeting gets the onlineMeeting property value. The onlineMeeting property
func (m *Event) GetOnlineMeeting()(OnlineMeetingInfoable) {
    return m.onlineMeeting
}
// GetOnlineMeetingProvider gets the onlineMeetingProvider property value. The onlineMeetingProvider property
func (m *Event) GetOnlineMeetingProvider()(*OnlineMeetingProviderType) {
    return m.onlineMeetingProvider
}
// GetOnlineMeetingUrl gets the onlineMeetingUrl property value. The onlineMeetingUrl property
func (m *Event) GetOnlineMeetingUrl()(*string) {
    return m.onlineMeetingUrl
}
// GetOrganizer gets the organizer property value. The organizer property
func (m *Event) GetOrganizer()(Recipientable) {
    return m.organizer
}
// GetOriginalEndTimeZone gets the originalEndTimeZone property value. The originalEndTimeZone property
func (m *Event) GetOriginalEndTimeZone()(*string) {
    return m.originalEndTimeZone
}
// GetOriginalStart gets the originalStart property value. The originalStart property
func (m *Event) GetOriginalStart()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.originalStart
}
// GetOriginalStartTimeZone gets the originalStartTimeZone property value. The originalStartTimeZone property
func (m *Event) GetOriginalStartTimeZone()(*string) {
    return m.originalStartTimeZone
}
// GetRecurrence gets the recurrence property value. The recurrence property
func (m *Event) GetRecurrence()(PatternedRecurrenceable) {
    return m.recurrence
}
// GetReminderMinutesBeforeStart gets the reminderMinutesBeforeStart property value. The reminderMinutesBeforeStart property
func (m *Event) GetReminderMinutesBeforeStart()(*int32) {
    return m.reminderMinutesBeforeStart
}
// GetResponseRequested gets the responseRequested property value. The responseRequested property
func (m *Event) GetResponseRequested()(*bool) {
    return m.responseRequested
}
// GetResponseStatus gets the responseStatus property value. The responseStatus property
func (m *Event) GetResponseStatus()(ResponseStatusable) {
    return m.responseStatus
}
// GetSensitivity gets the sensitivity property value. The sensitivity property
func (m *Event) GetSensitivity()(*Sensitivity) {
    return m.sensitivity
}
// GetSeriesMasterId gets the seriesMasterId property value. The seriesMasterId property
func (m *Event) GetSeriesMasterId()(*string) {
    return m.seriesMasterId
}
// GetShowAs gets the showAs property value. The showAs property
func (m *Event) GetShowAs()(*FreeBusyStatus) {
    return m.showAs
}
// GetSingleValueExtendedProperties gets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable) {
    return m.singleValueExtendedProperties
}
// GetStart gets the start property value. The start property
func (m *Event) GetStart()(DateTimeTimeZoneable) {
    return m.start
}
// GetSubject gets the subject property value. The subject property
func (m *Event) GetSubject()(*string) {
    return m.subject
}
// GetTransactionId gets the transactionId property value. The transactionId property
func (m *Event) GetTransactionId()(*string) {
    return m.transactionId
}
// GetType gets the type property value. The type property
func (m *Event) GetType()(*EventType) {
    return m.type_escaped
}
// GetWebLink gets the webLink property value. The webLink property
func (m *Event) GetWebLink()(*string) {
    return m.webLink
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAttachments())
        err = writer.WriteCollectionOfObjectValues("attachments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAttendees() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAttendees())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetExtensions())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetInstances())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetLocations())
        err = writer.WriteCollectionOfObjectValues("locations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMultiValueExtendedProperties() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMultiValueExtendedProperties())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSingleValueExtendedProperties())
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
    m.allowNewTimeProposals = value
}
// SetAttachments sets the attachments property value. The collection of FileAttachment, ItemAttachment, and referenceAttachment attachments for the event. Navigation property. Read-only. Nullable.
func (m *Event) SetAttachments(value []Attachmentable)() {
    m.attachments = value
}
// SetAttendees sets the attendees property value. The collection of attendees for the event.
func (m *Event) SetAttendees(value []Attendeeable)() {
    m.attendees = value
}
// SetBody sets the body property value. The body of the message associated with the event. It can be in HTML or text format.
func (m *Event) SetBody(value ItemBodyable)() {
    m.body = value
}
// SetBodyPreview sets the bodyPreview property value. The preview of the message associated with the event. It is in text format.
func (m *Event) SetBodyPreview(value *string)() {
    m.bodyPreview = value
}
// SetCalendar sets the calendar property value. The calendar that contains the event. Navigation property. Read-only.
func (m *Event) SetCalendar(value Calendarable)() {
    m.calendar = value
}
// SetEnd sets the end property value. The date, time, and time zone that the event ends. By default, the end time is in UTC.
func (m *Event) SetEnd(value DateTimeTimeZoneable)() {
    m.end = value
}
// SetExtensions sets the extensions property value. The collection of open extensions defined for the event. Nullable.
func (m *Event) SetExtensions(value []Extensionable)() {
    m.extensions = value
}
// SetHasAttachments sets the hasAttachments property value. Set to true if the event has attachments.
func (m *Event) SetHasAttachments(value *bool)() {
    m.hasAttachments = value
}
// SetHideAttendees sets the hideAttendees property value. When set to true, each attendee only sees themselves in the meeting request and meeting Tracking list. Default is false.
func (m *Event) SetHideAttendees(value *bool)() {
    m.hideAttendees = value
}
// SetICalUId sets the iCalUId property value. A unique identifier for an event across calendars. This ID is different for each occurrence in a recurring series. Read-only.
func (m *Event) SetICalUId(value *string)() {
    m.iCalUId = value
}
// SetImportance sets the importance property value. The importance property
func (m *Event) SetImportance(value *Importance)() {
    m.importance = value
}
// SetInstances sets the instances property value. The occurrences of a recurring series, if the event is a series master. This property includes occurrences that are part of the recurrence pattern, and exceptions that have been modified, but does not include occurrences that have been cancelled from the series. Navigation property. Read-only. Nullable.
func (m *Event) SetInstances(value []Eventable)() {
    m.instances = value
}
// SetIsAllDay sets the isAllDay property value. The isAllDay property
func (m *Event) SetIsAllDay(value *bool)() {
    m.isAllDay = value
}
// SetIsCancelled sets the isCancelled property value. The isCancelled property
func (m *Event) SetIsCancelled(value *bool)() {
    m.isCancelled = value
}
// SetIsDraft sets the isDraft property value. The isDraft property
func (m *Event) SetIsDraft(value *bool)() {
    m.isDraft = value
}
// SetIsOnlineMeeting sets the isOnlineMeeting property value. The isOnlineMeeting property
func (m *Event) SetIsOnlineMeeting(value *bool)() {
    m.isOnlineMeeting = value
}
// SetIsOrganizer sets the isOrganizer property value. The isOrganizer property
func (m *Event) SetIsOrganizer(value *bool)() {
    m.isOrganizer = value
}
// SetIsReminderOn sets the isReminderOn property value. The isReminderOn property
func (m *Event) SetIsReminderOn(value *bool)() {
    m.isReminderOn = value
}
// SetLocation sets the location property value. The location property
func (m *Event) SetLocation(value Locationable)() {
    m.location = value
}
// SetLocations sets the locations property value. The locations property
func (m *Event) SetLocations(value []Locationable)() {
    m.locations = value
}
// SetMultiValueExtendedProperties sets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)() {
    m.multiValueExtendedProperties = value
}
// SetOnlineMeeting sets the onlineMeeting property value. The onlineMeeting property
func (m *Event) SetOnlineMeeting(value OnlineMeetingInfoable)() {
    m.onlineMeeting = value
}
// SetOnlineMeetingProvider sets the onlineMeetingProvider property value. The onlineMeetingProvider property
func (m *Event) SetOnlineMeetingProvider(value *OnlineMeetingProviderType)() {
    m.onlineMeetingProvider = value
}
// SetOnlineMeetingUrl sets the onlineMeetingUrl property value. The onlineMeetingUrl property
func (m *Event) SetOnlineMeetingUrl(value *string)() {
    m.onlineMeetingUrl = value
}
// SetOrganizer sets the organizer property value. The organizer property
func (m *Event) SetOrganizer(value Recipientable)() {
    m.organizer = value
}
// SetOriginalEndTimeZone sets the originalEndTimeZone property value. The originalEndTimeZone property
func (m *Event) SetOriginalEndTimeZone(value *string)() {
    m.originalEndTimeZone = value
}
// SetOriginalStart sets the originalStart property value. The originalStart property
func (m *Event) SetOriginalStart(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.originalStart = value
}
// SetOriginalStartTimeZone sets the originalStartTimeZone property value. The originalStartTimeZone property
func (m *Event) SetOriginalStartTimeZone(value *string)() {
    m.originalStartTimeZone = value
}
// SetRecurrence sets the recurrence property value. The recurrence property
func (m *Event) SetRecurrence(value PatternedRecurrenceable)() {
    m.recurrence = value
}
// SetReminderMinutesBeforeStart sets the reminderMinutesBeforeStart property value. The reminderMinutesBeforeStart property
func (m *Event) SetReminderMinutesBeforeStart(value *int32)() {
    m.reminderMinutesBeforeStart = value
}
// SetResponseRequested sets the responseRequested property value. The responseRequested property
func (m *Event) SetResponseRequested(value *bool)() {
    m.responseRequested = value
}
// SetResponseStatus sets the responseStatus property value. The responseStatus property
func (m *Event) SetResponseStatus(value ResponseStatusable)() {
    m.responseStatus = value
}
// SetSensitivity sets the sensitivity property value. The sensitivity property
func (m *Event) SetSensitivity(value *Sensitivity)() {
    m.sensitivity = value
}
// SetSeriesMasterId sets the seriesMasterId property value. The seriesMasterId property
func (m *Event) SetSeriesMasterId(value *string)() {
    m.seriesMasterId = value
}
// SetShowAs sets the showAs property value. The showAs property
func (m *Event) SetShowAs(value *FreeBusyStatus)() {
    m.showAs = value
}
// SetSingleValueExtendedProperties sets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the event. Read-only. Nullable.
func (m *Event) SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)() {
    m.singleValueExtendedProperties = value
}
// SetStart sets the start property value. The start property
func (m *Event) SetStart(value DateTimeTimeZoneable)() {
    m.start = value
}
// SetSubject sets the subject property value. The subject property
func (m *Event) SetSubject(value *string)() {
    m.subject = value
}
// SetTransactionId sets the transactionId property value. The transactionId property
func (m *Event) SetTransactionId(value *string)() {
    m.transactionId = value
}
// SetType sets the type property value. The type property
func (m *Event) SetType(value *EventType)() {
    m.type_escaped = value
}
// SetWebLink sets the webLink property value. The webLink property
func (m *Event) SetWebLink(value *string)() {
    m.webLink = value
}
