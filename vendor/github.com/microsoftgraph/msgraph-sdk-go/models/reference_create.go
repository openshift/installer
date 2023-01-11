package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ReferenceCreate 
type ReferenceCreate struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataId property
    odataId *string
}
// NewReferenceCreate instantiates a new ReferenceCreate and sets the default values.
func NewReferenceCreate()(*ReferenceCreate) {
    m := &ReferenceCreate{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateReferenceCreateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateReferenceCreateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewReferenceCreate(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ReferenceCreate) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ReferenceCreate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.id"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataId)
    return res
}
// GetOdataId gets the @odata.id property value. The OdataId property
func (m *ReferenceCreate) GetOdataId()(*string) {
    return m.odataId
}
// Serialize serializes information the current object
func (m *ReferenceCreate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.id", m.GetOdataId())
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
func (m *ReferenceCreate) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataId sets the @odata.id property value. The OdataId property
func (m *ReferenceCreate) SetOdataId(value *string)() {
    m.odataId = value
}
