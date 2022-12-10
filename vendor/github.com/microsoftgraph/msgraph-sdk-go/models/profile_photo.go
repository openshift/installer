package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ProfilePhoto 
type ProfilePhoto struct {
    Entity
    // The height of the photo. Read-only.
    height *int32
    // The width of the photo. Read-only.
    width *int32
}
// NewProfilePhoto instantiates a new profilePhoto and sets the default values.
func NewProfilePhoto()(*ProfilePhoto) {
    m := &ProfilePhoto{
        Entity: *NewEntity(),
    }
    return m
}
// CreateProfilePhotoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateProfilePhotoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewProfilePhoto(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ProfilePhoto) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["height"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetHeight)
    res["width"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetWidth)
    return res
}
// GetHeight gets the height property value. The height of the photo. Read-only.
func (m *ProfilePhoto) GetHeight()(*int32) {
    return m.height
}
// GetWidth gets the width property value. The width of the photo. Read-only.
func (m *ProfilePhoto) GetWidth()(*int32) {
    return m.width
}
// Serialize serializes information the current object
func (m *ProfilePhoto) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("height", m.GetHeight())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("width", m.GetWidth())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetHeight sets the height property value. The height of the photo. Read-only.
func (m *ProfilePhoto) SetHeight(value *int32)() {
    m.height = value
}
// SetWidth sets the width property value. The width of the photo. Read-only.
func (m *ProfilePhoto) SetWidth(value *int32)() {
    m.width = value
}
