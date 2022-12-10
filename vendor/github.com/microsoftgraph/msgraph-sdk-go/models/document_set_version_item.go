package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DocumentSetVersionItem 
type DocumentSetVersionItem struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The unique identifier for the item.
    itemId *string
    // The OdataType property
    odataType *string
    // The title of the item.
    title *string
    // The version ID of the item.
    versionId *string
}
// NewDocumentSetVersionItem instantiates a new documentSetVersionItem and sets the default values.
func NewDocumentSetVersionItem()(*DocumentSetVersionItem) {
    m := &DocumentSetVersionItem{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateDocumentSetVersionItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDocumentSetVersionItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDocumentSetVersionItem(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *DocumentSetVersionItem) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DocumentSetVersionItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["itemId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetItemId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["title"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTitle)
    res["versionId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVersionId)
    return res
}
// GetItemId gets the itemId property value. The unique identifier for the item.
func (m *DocumentSetVersionItem) GetItemId()(*string) {
    return m.itemId
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *DocumentSetVersionItem) GetOdataType()(*string) {
    return m.odataType
}
// GetTitle gets the title property value. The title of the item.
func (m *DocumentSetVersionItem) GetTitle()(*string) {
    return m.title
}
// GetVersionId gets the versionId property value. The version ID of the item.
func (m *DocumentSetVersionItem) GetVersionId()(*string) {
    return m.versionId
}
// Serialize serializes information the current object
func (m *DocumentSetVersionItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("itemId", m.GetItemId())
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
        err := writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("versionId", m.GetVersionId())
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
func (m *DocumentSetVersionItem) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetItemId sets the itemId property value. The unique identifier for the item.
func (m *DocumentSetVersionItem) SetItemId(value *string)() {
    m.itemId = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DocumentSetVersionItem) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTitle sets the title property value. The title of the item.
func (m *DocumentSetVersionItem) SetTitle(value *string)() {
    m.title = value
}
// SetVersionId sets the versionId property value. The version ID of the item.
func (m *DocumentSetVersionItem) SetVersionId(value *string)() {
    m.versionId = value
}
