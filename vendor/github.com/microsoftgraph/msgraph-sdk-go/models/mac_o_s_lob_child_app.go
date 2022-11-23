package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MacOSLobChildApp contains properties of a macOS .app in the package
type MacOSLobChildApp struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The build number of the app.
    buildNumber *string
    // The bundleId of the app.
    bundleId *string
    // The OdataType property
    odataType *string
    // The version number of the app.
    versionNumber *string
}
// NewMacOSLobChildApp instantiates a new macOSLobChildApp and sets the default values.
func NewMacOSLobChildApp()(*MacOSLobChildApp) {
    m := &MacOSLobChildApp{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateMacOSLobChildAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMacOSLobChildAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMacOSLobChildApp(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MacOSLobChildApp) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetBuildNumber gets the buildNumber property value. The build number of the app.
func (m *MacOSLobChildApp) GetBuildNumber()(*string) {
    return m.buildNumber
}
// GetBundleId gets the bundleId property value. The bundleId of the app.
func (m *MacOSLobChildApp) GetBundleId()(*string) {
    return m.bundleId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MacOSLobChildApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["buildNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetBuildNumber)
    res["bundleId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetBundleId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["versionNumber"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVersionNumber)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *MacOSLobChildApp) GetOdataType()(*string) {
    return m.odataType
}
// GetVersionNumber gets the versionNumber property value. The version number of the app.
func (m *MacOSLobChildApp) GetVersionNumber()(*string) {
    return m.versionNumber
}
// Serialize serializes information the current object
func (m *MacOSLobChildApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("buildNumber", m.GetBuildNumber())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("bundleId", m.GetBundleId())
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
        err := writer.WriteStringValue("versionNumber", m.GetVersionNumber())
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
func (m *MacOSLobChildApp) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetBuildNumber sets the buildNumber property value. The build number of the app.
func (m *MacOSLobChildApp) SetBuildNumber(value *string)() {
    m.buildNumber = value
}
// SetBundleId sets the bundleId property value. The bundleId of the app.
func (m *MacOSLobChildApp) SetBundleId(value *string)() {
    m.bundleId = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MacOSLobChildApp) SetOdataType(value *string)() {
    m.odataType = value
}
// SetVersionNumber sets the versionNumber property value. The version number of the app.
func (m *MacOSLobChildApp) SetVersionNumber(value *string)() {
    m.versionNumber = value
}
