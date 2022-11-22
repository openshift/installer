package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AttachmentInfo 
type AttachmentInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The attachmentType property
    attachmentType *AttachmentType
    // The contentType property
    contentType *string
    // The name property
    name *string
    // The OdataType property
    odataType *string
    // The size property
    size *int64
}
// NewAttachmentInfo instantiates a new attachmentInfo and sets the default values.
func NewAttachmentInfo()(*AttachmentInfo) {
    m := &AttachmentInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAttachmentInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAttachmentInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttachmentInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AttachmentInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAttachmentType gets the attachmentType property value. The attachmentType property
func (m *AttachmentInfo) GetAttachmentType()(*AttachmentType) {
    return m.attachmentType
}
// GetContentType gets the contentType property value. The contentType property
func (m *AttachmentInfo) GetContentType()(*string) {
    return m.contentType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AttachmentInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attachmentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseAttachmentType , m.SetAttachmentType)
    res["contentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetContentType)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["size"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt64Value(m.SetSize)
    return res
}
// GetName gets the name property value. The name property
func (m *AttachmentInfo) GetName()(*string) {
    return m.name
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AttachmentInfo) GetOdataType()(*string) {
    return m.odataType
}
// GetSize gets the size property value. The size property
func (m *AttachmentInfo) GetSize()(*int64) {
    return m.size
}
// Serialize serializes information the current object
func (m *AttachmentInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAttachmentType() != nil {
        cast := (*m.GetAttachmentType()).String()
        err := writer.WriteStringValue("attachmentType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("contentType", m.GetContentType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
        err := writer.WriteInt64Value("size", m.GetSize())
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
func (m *AttachmentInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAttachmentType sets the attachmentType property value. The attachmentType property
func (m *AttachmentInfo) SetAttachmentType(value *AttachmentType)() {
    m.attachmentType = value
}
// SetContentType sets the contentType property value. The contentType property
func (m *AttachmentInfo) SetContentType(value *string)() {
    m.contentType = value
}
// SetName sets the name property value. The name property
func (m *AttachmentInfo) SetName(value *string)() {
    m.name = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AttachmentInfo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetSize sets the size property value. The size property
func (m *AttachmentInfo) SetSize(value *int64)() {
    m.size = value
}
