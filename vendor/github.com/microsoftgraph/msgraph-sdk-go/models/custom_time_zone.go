package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CustomTimeZone 
type CustomTimeZone struct {
    TimeZoneBase
    // The time offset of the time zone from Coordinated Universal Time (UTC). This value is in minutes. Time zones that are ahead of UTC have a positive offset; time zones that are behind UTC have a negative offset.
    bias *int32
    // Specifies when the time zone switches from standard time to daylight saving time.
    daylightOffset DaylightTimeZoneOffsetable
    // Specifies when the time zone switches from daylight saving time to standard time.
    standardOffset StandardTimeZoneOffsetable
}
// NewCustomTimeZone instantiates a new CustomTimeZone and sets the default values.
func NewCustomTimeZone()(*CustomTimeZone) {
    m := &CustomTimeZone{
        TimeZoneBase: *NewTimeZoneBase(),
    }
    odataTypeValue := "#microsoft.graph.customTimeZone";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateCustomTimeZoneFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCustomTimeZoneFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCustomTimeZone(), nil
}
// GetBias gets the bias property value. The time offset of the time zone from Coordinated Universal Time (UTC). This value is in minutes. Time zones that are ahead of UTC have a positive offset; time zones that are behind UTC have a negative offset.
func (m *CustomTimeZone) GetBias()(*int32) {
    return m.bias
}
// GetDaylightOffset gets the daylightOffset property value. Specifies when the time zone switches from standard time to daylight saving time.
func (m *CustomTimeZone) GetDaylightOffset()(DaylightTimeZoneOffsetable) {
    return m.daylightOffset
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CustomTimeZone) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TimeZoneBase.GetFieldDeserializers()
    res["bias"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetBias)
    res["daylightOffset"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDaylightTimeZoneOffsetFromDiscriminatorValue , m.SetDaylightOffset)
    res["standardOffset"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateStandardTimeZoneOffsetFromDiscriminatorValue , m.SetStandardOffset)
    return res
}
// GetStandardOffset gets the standardOffset property value. Specifies when the time zone switches from daylight saving time to standard time.
func (m *CustomTimeZone) GetStandardOffset()(StandardTimeZoneOffsetable) {
    return m.standardOffset
}
// Serialize serializes information the current object
func (m *CustomTimeZone) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TimeZoneBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("bias", m.GetBias())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("daylightOffset", m.GetDaylightOffset())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("standardOffset", m.GetStandardOffset())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBias sets the bias property value. The time offset of the time zone from Coordinated Universal Time (UTC). This value is in minutes. Time zones that are ahead of UTC have a positive offset; time zones that are behind UTC have a negative offset.
func (m *CustomTimeZone) SetBias(value *int32)() {
    m.bias = value
}
// SetDaylightOffset sets the daylightOffset property value. Specifies when the time zone switches from standard time to daylight saving time.
func (m *CustomTimeZone) SetDaylightOffset(value DaylightTimeZoneOffsetable)() {
    m.daylightOffset = value
}
// SetStandardOffset sets the standardOffset property value. Specifies when the time zone switches from daylight saving time to standard time.
func (m *CustomTimeZone) SetStandardOffset(value StandardTimeZoneOffsetable)() {
    m.standardOffset = value
}
