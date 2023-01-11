package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Album 
type Album struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Unique identifier of the [driveItem][] that is the cover of the album.
    coverImageItemId *string
    // The OdataType property
    odataType *string
}
// NewAlbum instantiates a new album and sets the default values.
func NewAlbum()(*Album) {
    m := &Album{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAlbumFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAlbumFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAlbum(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Album) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCoverImageItemId gets the coverImageItemId property value. Unique identifier of the [driveItem][] that is the cover of the album.
func (m *Album) GetCoverImageItemId()(*string) {
    return m.coverImageItemId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Album) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["coverImageItemId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCoverImageItemId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Album) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *Album) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("coverImageItemId", m.GetCoverImageItemId())
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
func (m *Album) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCoverImageItemId sets the coverImageItemId property value. Unique identifier of the [driveItem][] that is the cover of the album.
func (m *Album) SetCoverImageItemId(value *string)() {
    m.coverImageItemId = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Album) SetOdataType(value *string)() {
    m.odataType = value
}
