package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppListItem represents an app in the list of managed applications
type AppListItem struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The application or bundle identifier of the application
    appId *string
    // The Store URL of the application
    appStoreUrl *string
    // The application name
    name *string
    // The OdataType property
    odataType *string
    // The publisher of the application
    publisher *string
}
// NewAppListItem instantiates a new appListItem and sets the default values.
func NewAppListItem()(*AppListItem) {
    m := &AppListItem{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAppListItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAppListItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppListItem(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AppListItem) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAppId gets the appId property value. The application or bundle identifier of the application
func (m *AppListItem) GetAppId()(*string) {
    return m.appId
}
// GetAppStoreUrl gets the appStoreUrl property value. The Store URL of the application
func (m *AppListItem) GetAppStoreUrl()(*string) {
    return m.appStoreUrl
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AppListItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["appId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppId)
    res["appStoreUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppStoreUrl)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["publisher"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPublisher)
    return res
}
// GetName gets the name property value. The application name
func (m *AppListItem) GetName()(*string) {
    return m.name
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AppListItem) GetOdataType()(*string) {
    return m.odataType
}
// GetPublisher gets the publisher property value. The publisher of the application
func (m *AppListItem) GetPublisher()(*string) {
    return m.publisher
}
// Serialize serializes information the current object
func (m *AppListItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("appStoreUrl", m.GetAppStoreUrl())
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
        err := writer.WriteStringValue("publisher", m.GetPublisher())
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
func (m *AppListItem) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAppId sets the appId property value. The application or bundle identifier of the application
func (m *AppListItem) SetAppId(value *string)() {
    m.appId = value
}
// SetAppStoreUrl sets the appStoreUrl property value. The Store URL of the application
func (m *AppListItem) SetAppStoreUrl(value *string)() {
    m.appStoreUrl = value
}
// SetName sets the name property value. The application name
func (m *AppListItem) SetName(value *string)() {
    m.name = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AppListItem) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPublisher sets the publisher property value. The publisher of the application
func (m *AppListItem) SetPublisher(value *string)() {
    m.publisher = value
}
