package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OutlookGeoCoordinates 
type OutlookGeoCoordinates struct {
    // The accuracy of the latitude and longitude. As an example, the accuracy can be measured in meters, such as the latitude and longitude are accurate to within 50 meters.
    accuracy *float64
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The altitude of the location.
    altitude *float64
    // The accuracy of the altitude.
    altitudeAccuracy *float64
    // The latitude of the location.
    latitude *float64
    // The longitude of the location.
    longitude *float64
    // The OdataType property
    odataType *string
}
// NewOutlookGeoCoordinates instantiates a new outlookGeoCoordinates and sets the default values.
func NewOutlookGeoCoordinates()(*OutlookGeoCoordinates) {
    m := &OutlookGeoCoordinates{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateOutlookGeoCoordinatesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOutlookGeoCoordinatesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOutlookGeoCoordinates(), nil
}
// GetAccuracy gets the accuracy property value. The accuracy of the latitude and longitude. As an example, the accuracy can be measured in meters, such as the latitude and longitude are accurate to within 50 meters.
func (m *OutlookGeoCoordinates) GetAccuracy()(*float64) {
    return m.accuracy
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *OutlookGeoCoordinates) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAltitude gets the altitude property value. The altitude of the location.
func (m *OutlookGeoCoordinates) GetAltitude()(*float64) {
    return m.altitude
}
// GetAltitudeAccuracy gets the altitudeAccuracy property value. The accuracy of the altitude.
func (m *OutlookGeoCoordinates) GetAltitudeAccuracy()(*float64) {
    return m.altitudeAccuracy
}
// GetFieldDeserializers the deserialization information for the current model
func (m *OutlookGeoCoordinates) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["accuracy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetAccuracy)
    res["altitude"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetAltitude)
    res["altitudeAccuracy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetAltitudeAccuracy)
    res["latitude"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetLatitude)
    res["longitude"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetLongitude)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetLatitude gets the latitude property value. The latitude of the location.
func (m *OutlookGeoCoordinates) GetLatitude()(*float64) {
    return m.latitude
}
// GetLongitude gets the longitude property value. The longitude of the location.
func (m *OutlookGeoCoordinates) GetLongitude()(*float64) {
    return m.longitude
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *OutlookGeoCoordinates) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *OutlookGeoCoordinates) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteFloat64Value("accuracy", m.GetAccuracy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("altitude", m.GetAltitude())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("altitudeAccuracy", m.GetAltitudeAccuracy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("latitude", m.GetLatitude())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("longitude", m.GetLongitude())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccuracy sets the accuracy property value. The accuracy of the latitude and longitude. As an example, the accuracy can be measured in meters, such as the latitude and longitude are accurate to within 50 meters.
func (m *OutlookGeoCoordinates) SetAccuracy(value *float64)() {
    m.accuracy = value
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *OutlookGeoCoordinates) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAltitude sets the altitude property value. The altitude of the location.
func (m *OutlookGeoCoordinates) SetAltitude(value *float64)() {
    m.altitude = value
}
// SetAltitudeAccuracy sets the altitudeAccuracy property value. The accuracy of the altitude.
func (m *OutlookGeoCoordinates) SetAltitudeAccuracy(value *float64)() {
    m.altitudeAccuracy = value
}
// SetLatitude sets the latitude property value. The latitude of the location.
func (m *OutlookGeoCoordinates) SetLatitude(value *float64)() {
    m.latitude = value
}
// SetLongitude sets the longitude property value. The longitude of the location.
func (m *OutlookGeoCoordinates) SetLongitude(value *float64)() {
    m.longitude = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *OutlookGeoCoordinates) SetOdataType(value *string)() {
    m.odataType = value
}
