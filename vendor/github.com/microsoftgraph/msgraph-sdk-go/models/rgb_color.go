package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RgbColor color in RGB.
type RgbColor struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Blue value
    b *byte
    // Green value
    g *byte
    // The OdataType property
    odataType *string
    // Red value
    r *byte
}
// NewRgbColor instantiates a new rgbColor and sets the default values.
func NewRgbColor()(*RgbColor) {
    m := &RgbColor{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateRgbColorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRgbColorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRgbColor(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RgbColor) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetB gets the b property value. Blue value
func (m *RgbColor) GetB()(*byte) {
    return m.b
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RgbColor) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["b"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteValue(m.SetB)
    res["g"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteValue(m.SetG)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["r"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteValue(m.SetR)
    return res
}
// GetG gets the g property value. Green value
func (m *RgbColor) GetG()(*byte) {
    return m.g
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *RgbColor) GetOdataType()(*string) {
    return m.odataType
}
// GetR gets the r property value. Red value
func (m *RgbColor) GetR()(*byte) {
    return m.r
}
// Serialize serializes information the current object
func (m *RgbColor) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteByteValue("b", m.GetB())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteByteValue("g", m.GetG())
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
        err := writer.WriteByteValue("r", m.GetR())
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
func (m *RgbColor) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetB sets the b property value. Blue value
func (m *RgbColor) SetB(value *byte)() {
    m.b = value
}
// SetG sets the g property value. Green value
func (m *RgbColor) SetG(value *byte)() {
    m.g = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RgbColor) SetOdataType(value *string)() {
    m.odataType = value
}
// SetR sets the r property value. Red value
func (m *RgbColor) SetR(value *byte)() {
    m.r = value
}
