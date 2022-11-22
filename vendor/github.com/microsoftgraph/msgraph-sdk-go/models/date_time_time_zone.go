package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DateTimeTimeZone 
type DateTimeTimeZone struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // A single point of time in a combined date and time representation ({date}T{time}; for example, 2017-08-29T04:00:00.0000000).
    dateTime *string
    // The OdataType property
    odataType *string
    // Represents a time zone, for example, 'Pacific Standard Time'. See below for more possible values.
    timeZone *string
}
// NewDateTimeTimeZone instantiates a new dateTimeTimeZone and sets the default values.
func NewDateTimeTimeZone()(*DateTimeTimeZone) {
    m := &DateTimeTimeZone{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateDateTimeTimeZoneFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDateTimeTimeZoneFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDateTimeTimeZone(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *DateTimeTimeZone) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDateTime gets the dateTime property value. A single point of time in a combined date and time representation ({date}T{time}; for example, 2017-08-29T04:00:00.0000000).
func (m *DateTimeTimeZone) GetDateTime()(*string) {
    return m.dateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DateTimeTimeZone) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["dateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDateTime)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["timeZone"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTimeZone)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *DateTimeTimeZone) GetOdataType()(*string) {
    return m.odataType
}
// GetTimeZone gets the timeZone property value. Represents a time zone, for example, 'Pacific Standard Time'. See below for more possible values.
func (m *DateTimeTimeZone) GetTimeZone()(*string) {
    return m.timeZone
}
// Serialize serializes information the current object
func (m *DateTimeTimeZone) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("dateTime", m.GetDateTime())
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
        err := writer.WriteStringValue("timeZone", m.GetTimeZone())
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
func (m *DateTimeTimeZone) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDateTime sets the dateTime property value. A single point of time in a combined date and time representation ({date}T{time}; for example, 2017-08-29T04:00:00.0000000).
func (m *DateTimeTimeZone) SetDateTime(value *string)() {
    m.dateTime = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DateTimeTimeZone) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTimeZone sets the timeZone property value. Represents a time zone, for example, 'Pacific Standard Time'. See below for more possible values.
func (m *DateTimeTimeZone) SetTimeZone(value *string)() {
    m.timeZone = value
}
