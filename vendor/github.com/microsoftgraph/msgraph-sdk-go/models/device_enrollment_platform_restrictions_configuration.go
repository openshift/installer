package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceEnrollmentPlatformRestrictionsConfiguration 
type DeviceEnrollmentPlatformRestrictionsConfiguration struct {
    DeviceEnrollmentConfiguration
    // Android restrictions based on platform, platform operating system version, and device ownership
    androidRestriction DeviceEnrollmentPlatformRestrictionable
    // Ios restrictions based on platform, platform operating system version, and device ownership
    iosRestriction DeviceEnrollmentPlatformRestrictionable
    // Mac restrictions based on platform, platform operating system version, and device ownership
    macOSRestriction DeviceEnrollmentPlatformRestrictionable
    // Windows mobile restrictions based on platform, platform operating system version, and device ownership
    windowsMobileRestriction DeviceEnrollmentPlatformRestrictionable
    // Windows restrictions based on platform, platform operating system version, and device ownership
    windowsRestriction DeviceEnrollmentPlatformRestrictionable
}
// NewDeviceEnrollmentPlatformRestrictionsConfiguration instantiates a new DeviceEnrollmentPlatformRestrictionsConfiguration and sets the default values.
func NewDeviceEnrollmentPlatformRestrictionsConfiguration()(*DeviceEnrollmentPlatformRestrictionsConfiguration) {
    m := &DeviceEnrollmentPlatformRestrictionsConfiguration{
        DeviceEnrollmentConfiguration: *NewDeviceEnrollmentConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.deviceEnrollmentPlatformRestrictionsConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateDeviceEnrollmentPlatformRestrictionsConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDeviceEnrollmentPlatformRestrictionsConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceEnrollmentPlatformRestrictionsConfiguration(), nil
}
// GetAndroidRestriction gets the androidRestriction property value. Android restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) GetAndroidRestriction()(DeviceEnrollmentPlatformRestrictionable) {
    return m.androidRestriction
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceEnrollmentConfiguration.GetFieldDeserializers()
    res["androidRestriction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceEnrollmentPlatformRestrictionFromDiscriminatorValue , m.SetAndroidRestriction)
    res["iosRestriction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceEnrollmentPlatformRestrictionFromDiscriminatorValue , m.SetIosRestriction)
    res["macOSRestriction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceEnrollmentPlatformRestrictionFromDiscriminatorValue , m.SetMacOSRestriction)
    res["windowsMobileRestriction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceEnrollmentPlatformRestrictionFromDiscriminatorValue , m.SetWindowsMobileRestriction)
    res["windowsRestriction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceEnrollmentPlatformRestrictionFromDiscriminatorValue , m.SetWindowsRestriction)
    return res
}
// GetIosRestriction gets the iosRestriction property value. Ios restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) GetIosRestriction()(DeviceEnrollmentPlatformRestrictionable) {
    return m.iosRestriction
}
// GetMacOSRestriction gets the macOSRestriction property value. Mac restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) GetMacOSRestriction()(DeviceEnrollmentPlatformRestrictionable) {
    return m.macOSRestriction
}
// GetWindowsMobileRestriction gets the windowsMobileRestriction property value. Windows mobile restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) GetWindowsMobileRestriction()(DeviceEnrollmentPlatformRestrictionable) {
    return m.windowsMobileRestriction
}
// GetWindowsRestriction gets the windowsRestriction property value. Windows restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) GetWindowsRestriction()(DeviceEnrollmentPlatformRestrictionable) {
    return m.windowsRestriction
}
// Serialize serializes information the current object
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceEnrollmentConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("androidRestriction", m.GetAndroidRestriction())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("iosRestriction", m.GetIosRestriction())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("macOSRestriction", m.GetMacOSRestriction())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("windowsMobileRestriction", m.GetWindowsMobileRestriction())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("windowsRestriction", m.GetWindowsRestriction())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAndroidRestriction sets the androidRestriction property value. Android restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) SetAndroidRestriction(value DeviceEnrollmentPlatformRestrictionable)() {
    m.androidRestriction = value
}
// SetIosRestriction sets the iosRestriction property value. Ios restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) SetIosRestriction(value DeviceEnrollmentPlatformRestrictionable)() {
    m.iosRestriction = value
}
// SetMacOSRestriction sets the macOSRestriction property value. Mac restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) SetMacOSRestriction(value DeviceEnrollmentPlatformRestrictionable)() {
    m.macOSRestriction = value
}
// SetWindowsMobileRestriction sets the windowsMobileRestriction property value. Windows mobile restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) SetWindowsMobileRestriction(value DeviceEnrollmentPlatformRestrictionable)() {
    m.windowsMobileRestriction = value
}
// SetWindowsRestriction sets the windowsRestriction property value. Windows restrictions based on platform, platform operating system version, and device ownership
func (m *DeviceEnrollmentPlatformRestrictionsConfiguration) SetWindowsRestriction(value DeviceEnrollmentPlatformRestrictionable)() {
    m.windowsRestriction = value
}
