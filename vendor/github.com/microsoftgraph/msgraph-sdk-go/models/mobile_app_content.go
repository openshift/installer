package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MobileAppContent contains content properties for a specific app version. Each mobileAppContent can have multiple mobileAppContentFile.
type MobileAppContent struct {
    Entity
    // The list of files for this app content version.
    files []MobileAppContentFileable
}
// NewMobileAppContent instantiates a new mobileAppContent and sets the default values.
func NewMobileAppContent()(*MobileAppContent) {
    m := &MobileAppContent{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMobileAppContentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMobileAppContentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMobileAppContent(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MobileAppContent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["files"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMobileAppContentFileFromDiscriminatorValue , m.SetFiles)
    return res
}
// GetFiles gets the files property value. The list of files for this app content version.
func (m *MobileAppContent) GetFiles()([]MobileAppContentFileable) {
    return m.files
}
// Serialize serializes information the current object
func (m *MobileAppContent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetFiles() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetFiles())
        err = writer.WriteCollectionOfObjectValues("files", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFiles sets the files property value. The list of files for this app content version.
func (m *MobileAppContent) SetFiles(value []MobileAppContentFileable)() {
    m.files = value
}
