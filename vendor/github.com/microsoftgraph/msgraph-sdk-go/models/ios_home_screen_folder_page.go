package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosHomeScreenFolderPage a page for a folder containing apps and web clips on the Home Screen.
type IosHomeScreenFolderPage struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // A list of apps and web clips to appear on a page within a folder. This collection can contain a maximum of 500 elements.
    apps []IosHomeScreenAppable
    // Name of the folder page
    displayName *string
    // The OdataType property
    odataType *string
}
// NewIosHomeScreenFolderPage instantiates a new iosHomeScreenFolderPage and sets the default values.
func NewIosHomeScreenFolderPage()(*IosHomeScreenFolderPage) {
    m := &IosHomeScreenFolderPage{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateIosHomeScreenFolderPageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosHomeScreenFolderPageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosHomeScreenFolderPage(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *IosHomeScreenFolderPage) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetApps gets the apps property value. A list of apps and web clips to appear on a page within a folder. This collection can contain a maximum of 500 elements.
func (m *IosHomeScreenFolderPage) GetApps()([]IosHomeScreenAppable) {
    return m.apps
}
// GetDisplayName gets the displayName property value. Name of the folder page
func (m *IosHomeScreenFolderPage) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosHomeScreenFolderPage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["apps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIosHomeScreenAppFromDiscriminatorValue , m.SetApps)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *IosHomeScreenFolderPage) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *IosHomeScreenFolderPage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetApps())
        err := writer.WriteCollectionOfObjectValues("apps", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
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
func (m *IosHomeScreenFolderPage) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetApps sets the apps property value. A list of apps and web clips to appear on a page within a folder. This collection can contain a maximum of 500 elements.
func (m *IosHomeScreenFolderPage) SetApps(value []IosHomeScreenAppable)() {
    m.apps = value
}
// SetDisplayName sets the displayName property value. Name of the folder page
func (m *IosHomeScreenFolderPage) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *IosHomeScreenFolderPage) SetOdataType(value *string)() {
    m.odataType = value
}
