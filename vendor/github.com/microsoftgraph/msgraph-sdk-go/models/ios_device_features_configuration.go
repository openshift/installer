package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosDeviceFeaturesConfiguration 
type IosDeviceFeaturesConfiguration struct {
    AppleDeviceFeaturesConfigurationBase
    // Asset tag information for the device, displayed on the login window and lock screen.
    assetTagTemplate *string
    // A list of app and folders to appear on the Home Screen Dock. This collection can contain a maximum of 500 elements.
    homeScreenDockIcons []IosHomeScreenItemable
    // A list of pages on the Home Screen. This collection can contain a maximum of 500 elements.
    homeScreenPages []IosHomeScreenPageable
    // A footnote displayed on the login window and lock screen. Available in iOS 9.3.1 and later.
    lockScreenFootnote *string
    // Notification settings for each bundle id. Applicable to devices in supervised mode only (iOS 9.3 and later). This collection can contain a maximum of 500 elements.
    notificationSettings []IosNotificationSettingsable
}
// NewIosDeviceFeaturesConfiguration instantiates a new IosDeviceFeaturesConfiguration and sets the default values.
func NewIosDeviceFeaturesConfiguration()(*IosDeviceFeaturesConfiguration) {
    m := &IosDeviceFeaturesConfiguration{
        AppleDeviceFeaturesConfigurationBase: *NewAppleDeviceFeaturesConfigurationBase(),
    }
    odataTypeValue := "#microsoft.graph.iosDeviceFeaturesConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateIosDeviceFeaturesConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosDeviceFeaturesConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosDeviceFeaturesConfiguration(), nil
}
// GetAssetTagTemplate gets the assetTagTemplate property value. Asset tag information for the device, displayed on the login window and lock screen.
func (m *IosDeviceFeaturesConfiguration) GetAssetTagTemplate()(*string) {
    return m.assetTagTemplate
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosDeviceFeaturesConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AppleDeviceFeaturesConfigurationBase.GetFieldDeserializers()
    res["assetTagTemplate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAssetTagTemplate)
    res["homeScreenDockIcons"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIosHomeScreenItemFromDiscriminatorValue , m.SetHomeScreenDockIcons)
    res["homeScreenPages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIosHomeScreenPageFromDiscriminatorValue , m.SetHomeScreenPages)
    res["lockScreenFootnote"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLockScreenFootnote)
    res["notificationSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIosNotificationSettingsFromDiscriminatorValue , m.SetNotificationSettings)
    return res
}
// GetHomeScreenDockIcons gets the homeScreenDockIcons property value. A list of app and folders to appear on the Home Screen Dock. This collection can contain a maximum of 500 elements.
func (m *IosDeviceFeaturesConfiguration) GetHomeScreenDockIcons()([]IosHomeScreenItemable) {
    return m.homeScreenDockIcons
}
// GetHomeScreenPages gets the homeScreenPages property value. A list of pages on the Home Screen. This collection can contain a maximum of 500 elements.
func (m *IosDeviceFeaturesConfiguration) GetHomeScreenPages()([]IosHomeScreenPageable) {
    return m.homeScreenPages
}
// GetLockScreenFootnote gets the lockScreenFootnote property value. A footnote displayed on the login window and lock screen. Available in iOS 9.3.1 and later.
func (m *IosDeviceFeaturesConfiguration) GetLockScreenFootnote()(*string) {
    return m.lockScreenFootnote
}
// GetNotificationSettings gets the notificationSettings property value. Notification settings for each bundle id. Applicable to devices in supervised mode only (iOS 9.3 and later). This collection can contain a maximum of 500 elements.
func (m *IosDeviceFeaturesConfiguration) GetNotificationSettings()([]IosNotificationSettingsable) {
    return m.notificationSettings
}
// Serialize serializes information the current object
func (m *IosDeviceFeaturesConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AppleDeviceFeaturesConfigurationBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("assetTagTemplate", m.GetAssetTagTemplate())
        if err != nil {
            return err
        }
    }
    if m.GetHomeScreenDockIcons() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetHomeScreenDockIcons())
        err = writer.WriteCollectionOfObjectValues("homeScreenDockIcons", cast)
        if err != nil {
            return err
        }
    }
    if m.GetHomeScreenPages() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetHomeScreenPages())
        err = writer.WriteCollectionOfObjectValues("homeScreenPages", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("lockScreenFootnote", m.GetLockScreenFootnote())
        if err != nil {
            return err
        }
    }
    if m.GetNotificationSettings() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNotificationSettings())
        err = writer.WriteCollectionOfObjectValues("notificationSettings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssetTagTemplate sets the assetTagTemplate property value. Asset tag information for the device, displayed on the login window and lock screen.
func (m *IosDeviceFeaturesConfiguration) SetAssetTagTemplate(value *string)() {
    m.assetTagTemplate = value
}
// SetHomeScreenDockIcons sets the homeScreenDockIcons property value. A list of app and folders to appear on the Home Screen Dock. This collection can contain a maximum of 500 elements.
func (m *IosDeviceFeaturesConfiguration) SetHomeScreenDockIcons(value []IosHomeScreenItemable)() {
    m.homeScreenDockIcons = value
}
// SetHomeScreenPages sets the homeScreenPages property value. A list of pages on the Home Screen. This collection can contain a maximum of 500 elements.
func (m *IosDeviceFeaturesConfiguration) SetHomeScreenPages(value []IosHomeScreenPageable)() {
    m.homeScreenPages = value
}
// SetLockScreenFootnote sets the lockScreenFootnote property value. A footnote displayed on the login window and lock screen. Available in iOS 9.3.1 and later.
func (m *IosDeviceFeaturesConfiguration) SetLockScreenFootnote(value *string)() {
    m.lockScreenFootnote = value
}
// SetNotificationSettings sets the notificationSettings property value. Notification settings for each bundle id. Applicable to devices in supervised mode only (iOS 9.3 and later). This collection can contain a maximum of 500 elements.
func (m *IosDeviceFeaturesConfiguration) SetNotificationSettings(value []IosNotificationSettingsable)() {
    m.notificationSettings = value
}
