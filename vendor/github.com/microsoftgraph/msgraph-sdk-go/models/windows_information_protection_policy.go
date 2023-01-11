package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtectionPolicy 
type WindowsInformationProtectionPolicy struct {
    WindowsInformationProtection
    // Offline interval before app data is wiped (days)
    daysWithoutContactBeforeUnenroll *int32
    // Enrollment url for the MDM
    mdmEnrollmentUrl *string
    // Specifies the maximum amount of time (in minutes) allowed after the device is idle that will cause the device to become PIN or password locked.   Range is an integer X where 0 <= X <= 999.
    minutesOfInactivityBeforeDeviceLock *int32
    // Integer value that specifies the number of past PINs that can be associated to a user account that can't be reused. The largest number you can configure for this policy setting is 50. The lowest number you can configure for this policy setting is 0. If this policy is set to 0, then storage of previous PINs is not required. This node was added in Windows 10, version 1511. Default is 0.
    numberOfPastPinsRemembered *int32
    // The number of authentication failures allowed before the device will be wiped. A value of 0 disables device wipe functionality. Range is an integer X where 4 <= X <= 16 for desktop and 0 <= X <= 999 for mobile devices.
    passwordMaximumAttemptCount *int32
    // Integer value specifies the period of time (in days) that a PIN can be used before the system requires the user to change it. The largest number you can configure for this policy setting is 730. The lowest number you can configure for this policy setting is 0. If this policy is set to 0, then the user's PIN will never expire. This node was added in Windows 10, version 1511. Default is 0.
    pinExpirationDays *int32
    // Pin Character Requirements
    pinLowercaseLetters *WindowsInformationProtectionPinCharacterRequirements
    // Integer value that sets the minimum number of characters required for the PIN. Default value is 4. The lowest number you can configure for this policy setting is 4. The largest number you can configure must be less than the number configured in the Maximum PIN length policy setting or the number 127, whichever is the lowest.
    pinMinimumLength *int32
    // Pin Character Requirements
    pinSpecialCharacters *WindowsInformationProtectionPinCharacterRequirements
    // Pin Character Requirements
    pinUppercaseLetters *WindowsInformationProtectionPinCharacterRequirements
    // New property in RS2, pending documentation
    revokeOnMdmHandoffDisabled *bool
    // Boolean value that sets Windows Hello for Business as a method for signing into Windows.
    windowsHelloForBusinessBlocked *bool
}
// NewWindowsInformationProtectionPolicy instantiates a new WindowsInformationProtectionPolicy and sets the default values.
func NewWindowsInformationProtectionPolicy()(*WindowsInformationProtectionPolicy) {
    m := &WindowsInformationProtectionPolicy{
        WindowsInformationProtection: *NewWindowsInformationProtection(),
    }
    odataTypeValue := "#microsoft.graph.windowsInformationProtectionPolicy";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWindowsInformationProtectionPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsInformationProtectionPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsInformationProtectionPolicy(), nil
}
// GetDaysWithoutContactBeforeUnenroll gets the daysWithoutContactBeforeUnenroll property value. Offline interval before app data is wiped (days)
func (m *WindowsInformationProtectionPolicy) GetDaysWithoutContactBeforeUnenroll()(*int32) {
    return m.daysWithoutContactBeforeUnenroll
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsInformationProtectionPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WindowsInformationProtection.GetFieldDeserializers()
    res["daysWithoutContactBeforeUnenroll"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetDaysWithoutContactBeforeUnenroll)
    res["mdmEnrollmentUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMdmEnrollmentUrl)
    res["minutesOfInactivityBeforeDeviceLock"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetMinutesOfInactivityBeforeDeviceLock)
    res["numberOfPastPinsRemembered"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetNumberOfPastPinsRemembered)
    res["passwordMaximumAttemptCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPasswordMaximumAttemptCount)
    res["pinExpirationDays"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPinExpirationDays)
    res["pinLowercaseLetters"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWindowsInformationProtectionPinCharacterRequirements , m.SetPinLowercaseLetters)
    res["pinMinimumLength"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPinMinimumLength)
    res["pinSpecialCharacters"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWindowsInformationProtectionPinCharacterRequirements , m.SetPinSpecialCharacters)
    res["pinUppercaseLetters"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWindowsInformationProtectionPinCharacterRequirements , m.SetPinUppercaseLetters)
    res["revokeOnMdmHandoffDisabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetRevokeOnMdmHandoffDisabled)
    res["windowsHelloForBusinessBlocked"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetWindowsHelloForBusinessBlocked)
    return res
}
// GetMdmEnrollmentUrl gets the mdmEnrollmentUrl property value. Enrollment url for the MDM
func (m *WindowsInformationProtectionPolicy) GetMdmEnrollmentUrl()(*string) {
    return m.mdmEnrollmentUrl
}
// GetMinutesOfInactivityBeforeDeviceLock gets the minutesOfInactivityBeforeDeviceLock property value. Specifies the maximum amount of time (in minutes) allowed after the device is idle that will cause the device to become PIN or password locked.   Range is an integer X where 0 <= X <= 999.
func (m *WindowsInformationProtectionPolicy) GetMinutesOfInactivityBeforeDeviceLock()(*int32) {
    return m.minutesOfInactivityBeforeDeviceLock
}
// GetNumberOfPastPinsRemembered gets the numberOfPastPinsRemembered property value. Integer value that specifies the number of past PINs that can be associated to a user account that can't be reused. The largest number you can configure for this policy setting is 50. The lowest number you can configure for this policy setting is 0. If this policy is set to 0, then storage of previous PINs is not required. This node was added in Windows 10, version 1511. Default is 0.
func (m *WindowsInformationProtectionPolicy) GetNumberOfPastPinsRemembered()(*int32) {
    return m.numberOfPastPinsRemembered
}
// GetPasswordMaximumAttemptCount gets the passwordMaximumAttemptCount property value. The number of authentication failures allowed before the device will be wiped. A value of 0 disables device wipe functionality. Range is an integer X where 4 <= X <= 16 for desktop and 0 <= X <= 999 for mobile devices.
func (m *WindowsInformationProtectionPolicy) GetPasswordMaximumAttemptCount()(*int32) {
    return m.passwordMaximumAttemptCount
}
// GetPinExpirationDays gets the pinExpirationDays property value. Integer value specifies the period of time (in days) that a PIN can be used before the system requires the user to change it. The largest number you can configure for this policy setting is 730. The lowest number you can configure for this policy setting is 0. If this policy is set to 0, then the user's PIN will never expire. This node was added in Windows 10, version 1511. Default is 0.
func (m *WindowsInformationProtectionPolicy) GetPinExpirationDays()(*int32) {
    return m.pinExpirationDays
}
// GetPinLowercaseLetters gets the pinLowercaseLetters property value. Pin Character Requirements
func (m *WindowsInformationProtectionPolicy) GetPinLowercaseLetters()(*WindowsInformationProtectionPinCharacterRequirements) {
    return m.pinLowercaseLetters
}
// GetPinMinimumLength gets the pinMinimumLength property value. Integer value that sets the minimum number of characters required for the PIN. Default value is 4. The lowest number you can configure for this policy setting is 4. The largest number you can configure must be less than the number configured in the Maximum PIN length policy setting or the number 127, whichever is the lowest.
func (m *WindowsInformationProtectionPolicy) GetPinMinimumLength()(*int32) {
    return m.pinMinimumLength
}
// GetPinSpecialCharacters gets the pinSpecialCharacters property value. Pin Character Requirements
func (m *WindowsInformationProtectionPolicy) GetPinSpecialCharacters()(*WindowsInformationProtectionPinCharacterRequirements) {
    return m.pinSpecialCharacters
}
// GetPinUppercaseLetters gets the pinUppercaseLetters property value. Pin Character Requirements
func (m *WindowsInformationProtectionPolicy) GetPinUppercaseLetters()(*WindowsInformationProtectionPinCharacterRequirements) {
    return m.pinUppercaseLetters
}
// GetRevokeOnMdmHandoffDisabled gets the revokeOnMdmHandoffDisabled property value. New property in RS2, pending documentation
func (m *WindowsInformationProtectionPolicy) GetRevokeOnMdmHandoffDisabled()(*bool) {
    return m.revokeOnMdmHandoffDisabled
}
// GetWindowsHelloForBusinessBlocked gets the windowsHelloForBusinessBlocked property value. Boolean value that sets Windows Hello for Business as a method for signing into Windows.
func (m *WindowsInformationProtectionPolicy) GetWindowsHelloForBusinessBlocked()(*bool) {
    return m.windowsHelloForBusinessBlocked
}
// Serialize serializes information the current object
func (m *WindowsInformationProtectionPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WindowsInformationProtection.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("daysWithoutContactBeforeUnenroll", m.GetDaysWithoutContactBeforeUnenroll())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mdmEnrollmentUrl", m.GetMdmEnrollmentUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("minutesOfInactivityBeforeDeviceLock", m.GetMinutesOfInactivityBeforeDeviceLock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("numberOfPastPinsRemembered", m.GetNumberOfPastPinsRemembered())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMaximumAttemptCount", m.GetPasswordMaximumAttemptCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("pinExpirationDays", m.GetPinExpirationDays())
        if err != nil {
            return err
        }
    }
    if m.GetPinLowercaseLetters() != nil {
        cast := (*m.GetPinLowercaseLetters()).String()
        err = writer.WriteStringValue("pinLowercaseLetters", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("pinMinimumLength", m.GetPinMinimumLength())
        if err != nil {
            return err
        }
    }
    if m.GetPinSpecialCharacters() != nil {
        cast := (*m.GetPinSpecialCharacters()).String()
        err = writer.WriteStringValue("pinSpecialCharacters", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetPinUppercaseLetters() != nil {
        cast := (*m.GetPinUppercaseLetters()).String()
        err = writer.WriteStringValue("pinUppercaseLetters", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("revokeOnMdmHandoffDisabled", m.GetRevokeOnMdmHandoffDisabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsHelloForBusinessBlocked", m.GetWindowsHelloForBusinessBlocked())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDaysWithoutContactBeforeUnenroll sets the daysWithoutContactBeforeUnenroll property value. Offline interval before app data is wiped (days)
func (m *WindowsInformationProtectionPolicy) SetDaysWithoutContactBeforeUnenroll(value *int32)() {
    m.daysWithoutContactBeforeUnenroll = value
}
// SetMdmEnrollmentUrl sets the mdmEnrollmentUrl property value. Enrollment url for the MDM
func (m *WindowsInformationProtectionPolicy) SetMdmEnrollmentUrl(value *string)() {
    m.mdmEnrollmentUrl = value
}
// SetMinutesOfInactivityBeforeDeviceLock sets the minutesOfInactivityBeforeDeviceLock property value. Specifies the maximum amount of time (in minutes) allowed after the device is idle that will cause the device to become PIN or password locked.   Range is an integer X where 0 <= X <= 999.
func (m *WindowsInformationProtectionPolicy) SetMinutesOfInactivityBeforeDeviceLock(value *int32)() {
    m.minutesOfInactivityBeforeDeviceLock = value
}
// SetNumberOfPastPinsRemembered sets the numberOfPastPinsRemembered property value. Integer value that specifies the number of past PINs that can be associated to a user account that can't be reused. The largest number you can configure for this policy setting is 50. The lowest number you can configure for this policy setting is 0. If this policy is set to 0, then storage of previous PINs is not required. This node was added in Windows 10, version 1511. Default is 0.
func (m *WindowsInformationProtectionPolicy) SetNumberOfPastPinsRemembered(value *int32)() {
    m.numberOfPastPinsRemembered = value
}
// SetPasswordMaximumAttemptCount sets the passwordMaximumAttemptCount property value. The number of authentication failures allowed before the device will be wiped. A value of 0 disables device wipe functionality. Range is an integer X where 4 <= X <= 16 for desktop and 0 <= X <= 999 for mobile devices.
func (m *WindowsInformationProtectionPolicy) SetPasswordMaximumAttemptCount(value *int32)() {
    m.passwordMaximumAttemptCount = value
}
// SetPinExpirationDays sets the pinExpirationDays property value. Integer value specifies the period of time (in days) that a PIN can be used before the system requires the user to change it. The largest number you can configure for this policy setting is 730. The lowest number you can configure for this policy setting is 0. If this policy is set to 0, then the user's PIN will never expire. This node was added in Windows 10, version 1511. Default is 0.
func (m *WindowsInformationProtectionPolicy) SetPinExpirationDays(value *int32)() {
    m.pinExpirationDays = value
}
// SetPinLowercaseLetters sets the pinLowercaseLetters property value. Pin Character Requirements
func (m *WindowsInformationProtectionPolicy) SetPinLowercaseLetters(value *WindowsInformationProtectionPinCharacterRequirements)() {
    m.pinLowercaseLetters = value
}
// SetPinMinimumLength sets the pinMinimumLength property value. Integer value that sets the minimum number of characters required for the PIN. Default value is 4. The lowest number you can configure for this policy setting is 4. The largest number you can configure must be less than the number configured in the Maximum PIN length policy setting or the number 127, whichever is the lowest.
func (m *WindowsInformationProtectionPolicy) SetPinMinimumLength(value *int32)() {
    m.pinMinimumLength = value
}
// SetPinSpecialCharacters sets the pinSpecialCharacters property value. Pin Character Requirements
func (m *WindowsInformationProtectionPolicy) SetPinSpecialCharacters(value *WindowsInformationProtectionPinCharacterRequirements)() {
    m.pinSpecialCharacters = value
}
// SetPinUppercaseLetters sets the pinUppercaseLetters property value. Pin Character Requirements
func (m *WindowsInformationProtectionPolicy) SetPinUppercaseLetters(value *WindowsInformationProtectionPinCharacterRequirements)() {
    m.pinUppercaseLetters = value
}
// SetRevokeOnMdmHandoffDisabled sets the revokeOnMdmHandoffDisabled property value. New property in RS2, pending documentation
func (m *WindowsInformationProtectionPolicy) SetRevokeOnMdmHandoffDisabled(value *bool)() {
    m.revokeOnMdmHandoffDisabled = value
}
// SetWindowsHelloForBusinessBlocked sets the windowsHelloForBusinessBlocked property value. Boolean value that sets Windows Hello for Business as a method for signing into Windows.
func (m *WindowsInformationProtectionPolicy) SetWindowsHelloForBusinessBlocked(value *bool)() {
    m.windowsHelloForBusinessBlocked = value
}
