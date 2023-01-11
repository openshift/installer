package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ControlScore 
type ControlScore struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Control action category (Identity, Data, Device, Apps, Infrastructure).
    controlCategory *string
    // Control unique name.
    controlName *string
    // Description of the control.
    description *string
    // The OdataType property
    odataType *string
    // Tenant achieved score for the control (it varies day by day depending on tenant operations on the control).
    score *float64
}
// NewControlScore instantiates a new controlScore and sets the default values.
func NewControlScore()(*ControlScore) {
    m := &ControlScore{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateControlScoreFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateControlScoreFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewControlScore(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ControlScore) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetControlCategory gets the controlCategory property value. Control action category (Identity, Data, Device, Apps, Infrastructure).
func (m *ControlScore) GetControlCategory()(*string) {
    return m.controlCategory
}
// GetControlName gets the controlName property value. Control unique name.
func (m *ControlScore) GetControlName()(*string) {
    return m.controlName
}
// GetDescription gets the description property value. Description of the control.
func (m *ControlScore) GetDescription()(*string) {
    return m.description
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ControlScore) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["controlCategory"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetControlCategory)
    res["controlName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetControlName)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["score"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetScore)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ControlScore) GetOdataType()(*string) {
    return m.odataType
}
// GetScore gets the score property value. Tenant achieved score for the control (it varies day by day depending on tenant operations on the control).
func (m *ControlScore) GetScore()(*float64) {
    return m.score
}
// Serialize serializes information the current object
func (m *ControlScore) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("controlCategory", m.GetControlCategory())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("controlName", m.GetControlName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("description", m.GetDescription())
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
        err := writer.WriteFloat64Value("score", m.GetScore())
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
func (m *ControlScore) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetControlCategory sets the controlCategory property value. Control action category (Identity, Data, Device, Apps, Infrastructure).
func (m *ControlScore) SetControlCategory(value *string)() {
    m.controlCategory = value
}
// SetControlName sets the controlName property value. Control unique name.
func (m *ControlScore) SetControlName(value *string)() {
    m.controlName = value
}
// SetDescription sets the description property value. Description of the control.
func (m *ControlScore) SetDescription(value *string)() {
    m.description = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ControlScore) SetOdataType(value *string)() {
    m.odataType = value
}
// SetScore sets the score property value. Tenant achieved score for the control (it varies day by day depending on tenant operations on the control).
func (m *ControlScore) SetScore(value *float64)() {
    m.score = value
}
