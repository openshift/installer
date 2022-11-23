package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TimeSlot 
type TimeSlot struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The end property
    end DateTimeTimeZoneable
    // The OdataType property
    odataType *string
    // The start property
    start DateTimeTimeZoneable
}
// NewTimeSlot instantiates a new timeSlot and sets the default values.
func NewTimeSlot()(*TimeSlot) {
    m := &TimeSlot{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateTimeSlotFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTimeSlotFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTimeSlot(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TimeSlot) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetEnd gets the end property value. The end property
func (m *TimeSlot) GetEnd()(DateTimeTimeZoneable) {
    return m.end
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TimeSlot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["end"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue , m.SetEnd)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["start"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue , m.SetStart)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *TimeSlot) GetOdataType()(*string) {
    return m.odataType
}
// GetStart gets the start property value. The start property
func (m *TimeSlot) GetStart()(DateTimeTimeZoneable) {
    return m.start
}
// Serialize serializes information the current object
func (m *TimeSlot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("end", m.GetEnd())
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
        err := writer.WriteObjectValue("start", m.GetStart())
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
func (m *TimeSlot) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetEnd sets the end property value. The end property
func (m *TimeSlot) SetEnd(value DateTimeTimeZoneable)() {
    m.end = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TimeSlot) SetOdataType(value *string)() {
    m.odataType = value
}
// SetStart sets the start property value. The start property
func (m *TimeSlot) SetStart(value DateTimeTimeZoneable)() {
    m.start = value
}
