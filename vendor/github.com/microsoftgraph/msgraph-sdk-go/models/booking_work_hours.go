package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BookingWorkHours this type represents the set of working hours in a single day of the week.
type BookingWorkHours struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The day property
    day *DayOfWeek
    // The OdataType property
    odataType *string
    // A list of start/end times during a day.
    timeSlots []BookingWorkTimeSlotable
}
// NewBookingWorkHours instantiates a new bookingWorkHours and sets the default values.
func NewBookingWorkHours()(*BookingWorkHours) {
    m := &BookingWorkHours{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateBookingWorkHoursFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateBookingWorkHoursFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBookingWorkHours(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BookingWorkHours) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDay gets the day property value. The day property
func (m *BookingWorkHours) GetDay()(*DayOfWeek) {
    return m.day
}
// GetFieldDeserializers the deserialization information for the current model
func (m *BookingWorkHours) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["day"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDayOfWeek , m.SetDay)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["timeSlots"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateBookingWorkTimeSlotFromDiscriminatorValue , m.SetTimeSlots)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *BookingWorkHours) GetOdataType()(*string) {
    return m.odataType
}
// GetTimeSlots gets the timeSlots property value. A list of start/end times during a day.
func (m *BookingWorkHours) GetTimeSlots()([]BookingWorkTimeSlotable) {
    return m.timeSlots
}
// Serialize serializes information the current object
func (m *BookingWorkHours) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDay() != nil {
        cast := (*m.GetDay()).String()
        err := writer.WriteStringValue("day", &cast)
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
    if m.GetTimeSlots() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTimeSlots())
        err := writer.WriteCollectionOfObjectValues("timeSlots", cast)
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
func (m *BookingWorkHours) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDay sets the day property value. The day property
func (m *BookingWorkHours) SetDay(value *DayOfWeek)() {
    m.day = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BookingWorkHours) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTimeSlots sets the timeSlots property value. A list of start/end times during a day.
func (m *BookingWorkHours) SetTimeSlots(value []BookingWorkTimeSlotable)() {
    m.timeSlots = value
}
