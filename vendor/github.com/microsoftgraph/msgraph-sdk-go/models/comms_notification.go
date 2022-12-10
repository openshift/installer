package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CommsNotification 
type CommsNotification struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The changeType property
    changeType *ChangeType
    // The OdataType property
    odataType *string
    // URI of the resource that was changed.
    resourceUrl *string
}
// NewCommsNotification instantiates a new commsNotification and sets the default values.
func NewCommsNotification()(*CommsNotification) {
    m := &CommsNotification{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateCommsNotificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCommsNotificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCommsNotification(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CommsNotification) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetChangeType gets the changeType property value. The changeType property
func (m *CommsNotification) GetChangeType()(*ChangeType) {
    return m.changeType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CommsNotification) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["changeType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseChangeType , m.SetChangeType)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["resourceUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetResourceUrl)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *CommsNotification) GetOdataType()(*string) {
    return m.odataType
}
// GetResourceUrl gets the resourceUrl property value. URI of the resource that was changed.
func (m *CommsNotification) GetResourceUrl()(*string) {
    return m.resourceUrl
}
// Serialize serializes information the current object
func (m *CommsNotification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetChangeType() != nil {
        cast := (*m.GetChangeType()).String()
        err := writer.WriteStringValue("changeType", &cast)
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
        err := writer.WriteStringValue("resourceUrl", m.GetResourceUrl())
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
func (m *CommsNotification) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetChangeType sets the changeType property value. The changeType property
func (m *CommsNotification) SetChangeType(value *ChangeType)() {
    m.changeType = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *CommsNotification) SetOdataType(value *string)() {
    m.odataType = value
}
// SetResourceUrl sets the resourceUrl property value. URI of the resource that was changed.
func (m *CommsNotification) SetResourceUrl(value *string)() {
    m.resourceUrl = value
}
