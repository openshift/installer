package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BookingSchedulingPolicy this type represents the set of policies that dictate how bookings can be created in a Booking Calendar.
type BookingSchedulingPolicy struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // True if to allow customers to choose a specific person for the booking.
    allowStaffSelection *bool
    // Maximum number of days in advance that a booking can be made. It follows the ISO 8601 format.
    maximumAdvance *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration
    // The minimum amount of time before which bookings and cancellations must be made. It follows the ISO 8601 format.
    minimumLeadTime *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration
    // The OdataType property
    odataType *string
    // True to notify the business via email when a booking is created or changed. Use the email address specified in the email property of the bookingBusiness entity for the business.
    sendConfirmationsToOwner *bool
    // Duration of each time slot, denoted in ISO 8601 format.
    timeSlotInterval *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration
}
// NewBookingSchedulingPolicy instantiates a new bookingSchedulingPolicy and sets the default values.
func NewBookingSchedulingPolicy()(*BookingSchedulingPolicy) {
    m := &BookingSchedulingPolicy{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateBookingSchedulingPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateBookingSchedulingPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBookingSchedulingPolicy(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BookingSchedulingPolicy) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAllowStaffSelection gets the allowStaffSelection property value. True if to allow customers to choose a specific person for the booking.
func (m *BookingSchedulingPolicy) GetAllowStaffSelection()(*bool) {
    return m.allowStaffSelection
}
// GetFieldDeserializers the deserialization information for the current model
func (m *BookingSchedulingPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowStaffSelection"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowStaffSelection)
    res["maximumAdvance"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetISODurationValue(m.SetMaximumAdvance)
    res["minimumLeadTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetISODurationValue(m.SetMinimumLeadTime)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["sendConfirmationsToOwner"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetSendConfirmationsToOwner)
    res["timeSlotInterval"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetISODurationValue(m.SetTimeSlotInterval)
    return res
}
// GetMaximumAdvance gets the maximumAdvance property value. Maximum number of days in advance that a booking can be made. It follows the ISO 8601 format.
func (m *BookingSchedulingPolicy) GetMaximumAdvance()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    return m.maximumAdvance
}
// GetMinimumLeadTime gets the minimumLeadTime property value. The minimum amount of time before which bookings and cancellations must be made. It follows the ISO 8601 format.
func (m *BookingSchedulingPolicy) GetMinimumLeadTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    return m.minimumLeadTime
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *BookingSchedulingPolicy) GetOdataType()(*string) {
    return m.odataType
}
// GetSendConfirmationsToOwner gets the sendConfirmationsToOwner property value. True to notify the business via email when a booking is created or changed. Use the email address specified in the email property of the bookingBusiness entity for the business.
func (m *BookingSchedulingPolicy) GetSendConfirmationsToOwner()(*bool) {
    return m.sendConfirmationsToOwner
}
// GetTimeSlotInterval gets the timeSlotInterval property value. Duration of each time slot, denoted in ISO 8601 format.
func (m *BookingSchedulingPolicy) GetTimeSlotInterval()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    return m.timeSlotInterval
}
// Serialize serializes information the current object
func (m *BookingSchedulingPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowStaffSelection", m.GetAllowStaffSelection())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("maximumAdvance", m.GetMaximumAdvance())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("minimumLeadTime", m.GetMinimumLeadTime())
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
        err := writer.WriteBoolValue("sendConfirmationsToOwner", m.GetSendConfirmationsToOwner())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("timeSlotInterval", m.GetTimeSlotInterval())
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
func (m *BookingSchedulingPolicy) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAllowStaffSelection sets the allowStaffSelection property value. True if to allow customers to choose a specific person for the booking.
func (m *BookingSchedulingPolicy) SetAllowStaffSelection(value *bool)() {
    m.allowStaffSelection = value
}
// SetMaximumAdvance sets the maximumAdvance property value. Maximum number of days in advance that a booking can be made. It follows the ISO 8601 format.
func (m *BookingSchedulingPolicy) SetMaximumAdvance(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    m.maximumAdvance = value
}
// SetMinimumLeadTime sets the minimumLeadTime property value. The minimum amount of time before which bookings and cancellations must be made. It follows the ISO 8601 format.
func (m *BookingSchedulingPolicy) SetMinimumLeadTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    m.minimumLeadTime = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BookingSchedulingPolicy) SetOdataType(value *string)() {
    m.odataType = value
}
// SetSendConfirmationsToOwner sets the sendConfirmationsToOwner property value. True to notify the business via email when a booking is created or changed. Use the email address specified in the email property of the bookingBusiness entity for the business.
func (m *BookingSchedulingPolicy) SetSendConfirmationsToOwner(value *bool)() {
    m.sendConfirmationsToOwner = value
}
// SetTimeSlotInterval sets the timeSlotInterval property value. Duration of each time slot, denoted in ISO 8601 format.
func (m *BookingSchedulingPolicy) SetTimeSlotInterval(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    m.timeSlotInterval = value
}
