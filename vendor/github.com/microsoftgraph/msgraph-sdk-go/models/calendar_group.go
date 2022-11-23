package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CalendarGroup provides operations to manage the collection of agreement entities.
type CalendarGroup struct {
    Entity
    // The calendars in the calendar group. Navigation property. Read-only. Nullable.
    calendars []Calendarable
    // Identifies the version of the calendar group. Every time the calendar group is changed, ChangeKey changes as well. This allows Exchange to apply changes to the correct version of the object. Read-only.
    changeKey *string
    // The class identifier. Read-only.
    classId *string
    // The group name.
    name *string
}
// NewCalendarGroup instantiates a new calendarGroup and sets the default values.
func NewCalendarGroup()(*CalendarGroup) {
    m := &CalendarGroup{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCalendarGroupFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCalendarGroupFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCalendarGroup(), nil
}
// GetCalendars gets the calendars property value. The calendars in the calendar group. Navigation property. Read-only. Nullable.
func (m *CalendarGroup) GetCalendars()([]Calendarable) {
    return m.calendars
}
// GetChangeKey gets the changeKey property value. Identifies the version of the calendar group. Every time the calendar group is changed, ChangeKey changes as well. This allows Exchange to apply changes to the correct version of the object. Read-only.
func (m *CalendarGroup) GetChangeKey()(*string) {
    return m.changeKey
}
// GetClassId gets the classId property value. The class identifier. Read-only.
func (m *CalendarGroup) GetClassId()(*string) {
    return m.classId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CalendarGroup) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["calendars"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCalendarFromDiscriminatorValue , m.SetCalendars)
    res["changeKey"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetChangeKey)
    res["classId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetClassId)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    return res
}
// GetName gets the name property value. The group name.
func (m *CalendarGroup) GetName()(*string) {
    return m.name
}
// Serialize serializes information the current object
func (m *CalendarGroup) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCalendars() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCalendars())
        err = writer.WriteCollectionOfObjectValues("calendars", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("changeKey", m.GetChangeKey())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("classId", m.GetClassId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCalendars sets the calendars property value. The calendars in the calendar group. Navigation property. Read-only. Nullable.
func (m *CalendarGroup) SetCalendars(value []Calendarable)() {
    m.calendars = value
}
// SetChangeKey sets the changeKey property value. Identifies the version of the calendar group. Every time the calendar group is changed, ChangeKey changes as well. This allows Exchange to apply changes to the correct version of the object. Read-only.
func (m *CalendarGroup) SetChangeKey(value *string)() {
    m.changeKey = value
}
// SetClassId sets the classId property value. The class identifier. Read-only.
func (m *CalendarGroup) SetClassId(value *string)() {
    m.classId = value
}
// SetName sets the name property value. The group name.
func (m *CalendarGroup) SetName(value *string)() {
    m.name = value
}
