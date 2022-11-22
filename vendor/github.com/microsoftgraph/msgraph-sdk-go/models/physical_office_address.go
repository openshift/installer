package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PhysicalOfficeAddress 
type PhysicalOfficeAddress struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The city.
    city *string
    // The country or region. It's a free-format string value, for example, 'United States'.
    countryOrRegion *string
    // The OdataType property
    odataType *string
    // Office location such as building and office number for an organizational contact.
    officeLocation *string
    // The postal code.
    postalCode *string
    // The state.
    state *string
    // The street.
    street *string
}
// NewPhysicalOfficeAddress instantiates a new physicalOfficeAddress and sets the default values.
func NewPhysicalOfficeAddress()(*PhysicalOfficeAddress) {
    m := &PhysicalOfficeAddress{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePhysicalOfficeAddressFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePhysicalOfficeAddressFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPhysicalOfficeAddress(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PhysicalOfficeAddress) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCity gets the city property value. The city.
func (m *PhysicalOfficeAddress) GetCity()(*string) {
    return m.city
}
// GetCountryOrRegion gets the countryOrRegion property value. The country or region. It's a free-format string value, for example, 'United States'.
func (m *PhysicalOfficeAddress) GetCountryOrRegion()(*string) {
    return m.countryOrRegion
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PhysicalOfficeAddress) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["city"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCity)
    res["countryOrRegion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCountryOrRegion)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["officeLocation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOfficeLocation)
    res["postalCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPostalCode)
    res["state"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetState)
    res["street"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetStreet)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *PhysicalOfficeAddress) GetOdataType()(*string) {
    return m.odataType
}
// GetOfficeLocation gets the officeLocation property value. Office location such as building and office number for an organizational contact.
func (m *PhysicalOfficeAddress) GetOfficeLocation()(*string) {
    return m.officeLocation
}
// GetPostalCode gets the postalCode property value. The postal code.
func (m *PhysicalOfficeAddress) GetPostalCode()(*string) {
    return m.postalCode
}
// GetState gets the state property value. The state.
func (m *PhysicalOfficeAddress) GetState()(*string) {
    return m.state
}
// GetStreet gets the street property value. The street.
func (m *PhysicalOfficeAddress) GetStreet()(*string) {
    return m.street
}
// Serialize serializes information the current object
func (m *PhysicalOfficeAddress) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("city", m.GetCity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("countryOrRegion", m.GetCountryOrRegion())
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
        err := writer.WriteStringValue("officeLocation", m.GetOfficeLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("postalCode", m.GetPostalCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("street", m.GetStreet())
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
func (m *PhysicalOfficeAddress) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCity sets the city property value. The city.
func (m *PhysicalOfficeAddress) SetCity(value *string)() {
    m.city = value
}
// SetCountryOrRegion sets the countryOrRegion property value. The country or region. It's a free-format string value, for example, 'United States'.
func (m *PhysicalOfficeAddress) SetCountryOrRegion(value *string)() {
    m.countryOrRegion = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PhysicalOfficeAddress) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOfficeLocation sets the officeLocation property value. Office location such as building and office number for an organizational contact.
func (m *PhysicalOfficeAddress) SetOfficeLocation(value *string)() {
    m.officeLocation = value
}
// SetPostalCode sets the postalCode property value. The postal code.
func (m *PhysicalOfficeAddress) SetPostalCode(value *string)() {
    m.postalCode = value
}
// SetState sets the state property value. The state.
func (m *PhysicalOfficeAddress) SetState(value *string)() {
    m.state = value
}
// SetStreet sets the street property value. The street.
func (m *PhysicalOfficeAddress) SetStreet(value *string)() {
    m.street = value
}
