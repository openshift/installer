package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ItemBody 
type ItemBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The content of the item.
    content *string
    // The type of the content. Possible values are text and html.
    contentType *BodyType
    // The OdataType property
    odataType *string
}
// NewItemBody instantiates a new itemBody and sets the default values.
func NewItemBody()(*ItemBody) {
    m := &ItemBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateItemBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateItemBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ItemBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetContent gets the content property value. The content of the item.
func (m *ItemBody) GetContent()(*string) {
    return m.content
}
// GetContentType gets the contentType property value. The type of the content. Possible values are text and html.
func (m *ItemBody) GetContentType()(*BodyType) {
    return m.contentType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ItemBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["content"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetContent)
    res["contentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseBodyType , m.SetContentType)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ItemBody) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ItemBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("content", m.GetContent())
        if err != nil {
            return err
        }
    }
    if m.GetContentType() != nil {
        cast := (*m.GetContentType()).String()
        err := writer.WriteStringValue("contentType", &cast)
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
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ItemBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetContent sets the content property value. The content of the item.
func (m *ItemBody) SetContent(value *string)() {
    m.content = value
}
// SetContentType sets the contentType property value. The type of the content. Possible values are text and html.
func (m *ItemBody) SetContentType(value *BodyType)() {
    m.contentType = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ItemBody) SetOdataType(value *string)() {
    m.odataType = value
}
