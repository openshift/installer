package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MimeContent contains properties for a generic mime content.
type MimeContent struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // Indicates the content mime type.
    type_escaped *string
    // The byte array that contains the actual content.
    value []byte
}
// NewMimeContent instantiates a new mimeContent and sets the default values.
func NewMimeContent()(*MimeContent) {
    m := &MimeContent{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateMimeContentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMimeContentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMimeContent(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MimeContent) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MimeContent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetType)
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteArrayValue(m.SetValue)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *MimeContent) GetOdataType()(*string) {
    return m.odataType
}
// GetType gets the type property value. Indicates the content mime type.
func (m *MimeContent) GetType()(*string) {
    return m.type_escaped
}
// GetValue gets the value property value. The byte array that contains the actual content.
func (m *MimeContent) GetValue()([]byte) {
    return m.value
}
// Serialize serializes information the current object
func (m *MimeContent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("type", m.GetType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteByteArrayValue("value", m.GetValue())
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
func (m *MimeContent) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MimeContent) SetOdataType(value *string)() {
    m.odataType = value
}
// SetType sets the type property value. Indicates the content mime type.
func (m *MimeContent) SetType(value *string)() {
    m.type_escaped = value
}
// SetValue sets the value property value. The byte array that contains the actual content.
func (m *MimeContent) SetValue(value []byte)() {
    m.value = value
}
