package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkingHours 
type WorkingHours struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The days of the week on which the user works.
    daysOfWeek []DayOfWeek
    // The time of the day that the user stops working.
    endTime *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly
    // The OdataType property
    odataType *string
    // The time of the day that the user starts working.
    startTime *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly
    // The time zone to which the working hours apply.
    timeZone TimeZoneBaseable
}
// NewWorkingHours instantiates a new workingHours and sets the default values.
func NewWorkingHours()(*WorkingHours) {
    m := &WorkingHours{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWorkingHoursFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkingHoursFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkingHours(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *WorkingHours) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDaysOfWeek gets the daysOfWeek property value. The days of the week on which the user works.
func (m *WorkingHours) GetDaysOfWeek()([]DayOfWeek) {
    return m.daysOfWeek
}
// GetEndTime gets the endTime property value. The time of the day that the user stops working.
func (m *WorkingHours) GetEndTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly) {
    return m.endTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkingHours) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["daysOfWeek"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParseDayOfWeek , m.SetDaysOfWeek)
    res["endTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeOnlyValue(m.SetEndTime)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["startTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeOnlyValue(m.SetStartTime)
    res["timeZone"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTimeZoneBaseFromDiscriminatorValue , m.SetTimeZone)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *WorkingHours) GetOdataType()(*string) {
    return m.odataType
}
// GetStartTime gets the startTime property value. The time of the day that the user starts working.
func (m *WorkingHours) GetStartTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly) {
    return m.startTime
}
// GetTimeZone gets the timeZone property value. The time zone to which the working hours apply.
func (m *WorkingHours) GetTimeZone()(TimeZoneBaseable) {
    return m.timeZone
}
// Serialize serializes information the current object
func (m *WorkingHours) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDaysOfWeek() != nil {
        err := writer.WriteCollectionOfStringValues("daysOfWeek", SerializeDayOfWeek(m.GetDaysOfWeek()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeOnlyValue("endTime", m.GetEndTime())
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
        err := writer.WriteTimeOnlyValue("startTime", m.GetStartTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("timeZone", m.GetTimeZone())
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
func (m *WorkingHours) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDaysOfWeek sets the daysOfWeek property value. The days of the week on which the user works.
func (m *WorkingHours) SetDaysOfWeek(value []DayOfWeek)() {
    m.daysOfWeek = value
}
// SetEndTime sets the endTime property value. The time of the day that the user stops working.
func (m *WorkingHours) SetEndTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)() {
    m.endTime = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WorkingHours) SetOdataType(value *string)() {
    m.odataType = value
}
// SetStartTime sets the startTime property value. The time of the day that the user starts working.
func (m *WorkingHours) SetStartTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)() {
    m.startTime = value
}
// SetTimeZone sets the timeZone property value. The time zone to which the working hours apply.
func (m *WorkingHours) SetTimeZone(value TimeZoneBaseable)() {
    m.timeZone = value
}
