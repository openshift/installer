package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10MobileCompliancePolicy 
type Windows10MobileCompliancePolicy struct {
    DeviceCompliancePolicy
    // Require devices to be reported healthy by Windows Device Health Attestation - bit locker is enabled
    bitLockerEnabled *bool
    // Require devices to be reported as healthy by Windows Device Health Attestation.
    codeIntegrityEnabled *bool
    // Require devices to be reported as healthy by Windows Device Health Attestation - early launch antimalware driver is enabled.
    earlyLaunchAntiMalwareDriverEnabled *bool
    // Maximum Windows Phone version.
    osMaximumVersion *string
    // Minimum Windows Phone version.
    osMinimumVersion *string
    // Whether or not to block syncing the calendar.
    passwordBlockSimple *bool
    // Number of days before password expiration. Valid values 1 to 255
    passwordExpirationDays *int32
    // The number of character sets required in the password.
    passwordMinimumCharacterSetCount *int32
    // Minimum password length. Valid values 4 to 16
    passwordMinimumLength *int32
    // Minutes of inactivity before a password is required.
    passwordMinutesOfInactivityBeforeLock *int32
    // The number of previous passwords to prevent re-use of.
    passwordPreviousPasswordBlockCount *int32
    // Require a password to unlock Windows Phone device.
    passwordRequired *bool
    // Possible values of required passwords.
    passwordRequiredType *RequiredPasswordType
    // Require a password to unlock an idle device.
    passwordRequireToUnlockFromIdle *bool
    // Require devices to be reported as healthy by Windows Device Health Attestation - secure boot is enabled.
    secureBootEnabled *bool
    // Require encryption on windows devices.
    storageRequireEncryption *bool
}
// NewWindows10MobileCompliancePolicy instantiates a new Windows10MobileCompliancePolicy and sets the default values.
func NewWindows10MobileCompliancePolicy()(*Windows10MobileCompliancePolicy) {
    m := &Windows10MobileCompliancePolicy{
        DeviceCompliancePolicy: *NewDeviceCompliancePolicy(),
    }
    odataTypeValue := "#microsoft.graph.windows10MobileCompliancePolicy";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWindows10MobileCompliancePolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindows10MobileCompliancePolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10MobileCompliancePolicy(), nil
}
// GetBitLockerEnabled gets the bitLockerEnabled property value. Require devices to be reported healthy by Windows Device Health Attestation - bit locker is enabled
func (m *Windows10MobileCompliancePolicy) GetBitLockerEnabled()(*bool) {
    return m.bitLockerEnabled
}
// GetCodeIntegrityEnabled gets the codeIntegrityEnabled property value. Require devices to be reported as healthy by Windows Device Health Attestation.
func (m *Windows10MobileCompliancePolicy) GetCodeIntegrityEnabled()(*bool) {
    return m.codeIntegrityEnabled
}
// GetEarlyLaunchAntiMalwareDriverEnabled gets the earlyLaunchAntiMalwareDriverEnabled property value. Require devices to be reported as healthy by Windows Device Health Attestation - early launch antimalware driver is enabled.
func (m *Windows10MobileCompliancePolicy) GetEarlyLaunchAntiMalwareDriverEnabled()(*bool) {
    return m.earlyLaunchAntiMalwareDriverEnabled
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Windows10MobileCompliancePolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceCompliancePolicy.GetFieldDeserializers()
    res["bitLockerEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetBitLockerEnabled)
    res["codeIntegrityEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetCodeIntegrityEnabled)
    res["earlyLaunchAntiMalwareDriverEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEarlyLaunchAntiMalwareDriverEnabled)
    res["osMaximumVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOsMaximumVersion)
    res["osMinimumVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOsMinimumVersion)
    res["passwordBlockSimple"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetPasswordBlockSimple)
    res["passwordExpirationDays"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPasswordExpirationDays)
    res["passwordMinimumCharacterSetCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPasswordMinimumCharacterSetCount)
    res["passwordMinimumLength"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPasswordMinimumLength)
    res["passwordMinutesOfInactivityBeforeLock"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPasswordMinutesOfInactivityBeforeLock)
    res["passwordPreviousPasswordBlockCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPasswordPreviousPasswordBlockCount)
    res["passwordRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetPasswordRequired)
    res["passwordRequiredType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseRequiredPasswordType , m.SetPasswordRequiredType)
    res["passwordRequireToUnlockFromIdle"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetPasswordRequireToUnlockFromIdle)
    res["secureBootEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetSecureBootEnabled)
    res["storageRequireEncryption"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetStorageRequireEncryption)
    return res
}
// GetOsMaximumVersion gets the osMaximumVersion property value. Maximum Windows Phone version.
func (m *Windows10MobileCompliancePolicy) GetOsMaximumVersion()(*string) {
    return m.osMaximumVersion
}
// GetOsMinimumVersion gets the osMinimumVersion property value. Minimum Windows Phone version.
func (m *Windows10MobileCompliancePolicy) GetOsMinimumVersion()(*string) {
    return m.osMinimumVersion
}
// GetPasswordBlockSimple gets the passwordBlockSimple property value. Whether or not to block syncing the calendar.
func (m *Windows10MobileCompliancePolicy) GetPasswordBlockSimple()(*bool) {
    return m.passwordBlockSimple
}
// GetPasswordExpirationDays gets the passwordExpirationDays property value. Number of days before password expiration. Valid values 1 to 255
func (m *Windows10MobileCompliancePolicy) GetPasswordExpirationDays()(*int32) {
    return m.passwordExpirationDays
}
// GetPasswordMinimumCharacterSetCount gets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
func (m *Windows10MobileCompliancePolicy) GetPasswordMinimumCharacterSetCount()(*int32) {
    return m.passwordMinimumCharacterSetCount
}
// GetPasswordMinimumLength gets the passwordMinimumLength property value. Minimum password length. Valid values 4 to 16
func (m *Windows10MobileCompliancePolicy) GetPasswordMinimumLength()(*int32) {
    return m.passwordMinimumLength
}
// GetPasswordMinutesOfInactivityBeforeLock gets the passwordMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a password is required.
func (m *Windows10MobileCompliancePolicy) GetPasswordMinutesOfInactivityBeforeLock()(*int32) {
    return m.passwordMinutesOfInactivityBeforeLock
}
// GetPasswordPreviousPasswordBlockCount gets the passwordPreviousPasswordBlockCount property value. The number of previous passwords to prevent re-use of.
func (m *Windows10MobileCompliancePolicy) GetPasswordPreviousPasswordBlockCount()(*int32) {
    return m.passwordPreviousPasswordBlockCount
}
// GetPasswordRequired gets the passwordRequired property value. Require a password to unlock Windows Phone device.
func (m *Windows10MobileCompliancePolicy) GetPasswordRequired()(*bool) {
    return m.passwordRequired
}
// GetPasswordRequiredType gets the passwordRequiredType property value. Possible values of required passwords.
func (m *Windows10MobileCompliancePolicy) GetPasswordRequiredType()(*RequiredPasswordType) {
    return m.passwordRequiredType
}
// GetPasswordRequireToUnlockFromIdle gets the passwordRequireToUnlockFromIdle property value. Require a password to unlock an idle device.
func (m *Windows10MobileCompliancePolicy) GetPasswordRequireToUnlockFromIdle()(*bool) {
    return m.passwordRequireToUnlockFromIdle
}
// GetSecureBootEnabled gets the secureBootEnabled property value. Require devices to be reported as healthy by Windows Device Health Attestation - secure boot is enabled.
func (m *Windows10MobileCompliancePolicy) GetSecureBootEnabled()(*bool) {
    return m.secureBootEnabled
}
// GetStorageRequireEncryption gets the storageRequireEncryption property value. Require encryption on windows devices.
func (m *Windows10MobileCompliancePolicy) GetStorageRequireEncryption()(*bool) {
    return m.storageRequireEncryption
}
// Serialize serializes information the current object
func (m *Windows10MobileCompliancePolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceCompliancePolicy.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("bitLockerEnabled", m.GetBitLockerEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("codeIntegrityEnabled", m.GetCodeIntegrityEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("earlyLaunchAntiMalwareDriverEnabled", m.GetEarlyLaunchAntiMalwareDriverEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("osMaximumVersion", m.GetOsMaximumVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("osMinimumVersion", m.GetOsMinimumVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordBlockSimple", m.GetPasswordBlockSimple())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordExpirationDays", m.GetPasswordExpirationDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinimumCharacterSetCount", m.GetPasswordMinimumCharacterSetCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinimumLength", m.GetPasswordMinimumLength())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinutesOfInactivityBeforeLock", m.GetPasswordMinutesOfInactivityBeforeLock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordPreviousPasswordBlockCount", m.GetPasswordPreviousPasswordBlockCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordRequired", m.GetPasswordRequired())
        if err != nil {
            return err
        }
    }
    if m.GetPasswordRequiredType() != nil {
        cast := (*m.GetPasswordRequiredType()).String()
        err = writer.WriteStringValue("passwordRequiredType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordRequireToUnlockFromIdle", m.GetPasswordRequireToUnlockFromIdle())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("secureBootEnabled", m.GetSecureBootEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageRequireEncryption", m.GetStorageRequireEncryption())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBitLockerEnabled sets the bitLockerEnabled property value. Require devices to be reported healthy by Windows Device Health Attestation - bit locker is enabled
func (m *Windows10MobileCompliancePolicy) SetBitLockerEnabled(value *bool)() {
    m.bitLockerEnabled = value
}
// SetCodeIntegrityEnabled sets the codeIntegrityEnabled property value. Require devices to be reported as healthy by Windows Device Health Attestation.
func (m *Windows10MobileCompliancePolicy) SetCodeIntegrityEnabled(value *bool)() {
    m.codeIntegrityEnabled = value
}
// SetEarlyLaunchAntiMalwareDriverEnabled sets the earlyLaunchAntiMalwareDriverEnabled property value. Require devices to be reported as healthy by Windows Device Health Attestation - early launch antimalware driver is enabled.
func (m *Windows10MobileCompliancePolicy) SetEarlyLaunchAntiMalwareDriverEnabled(value *bool)() {
    m.earlyLaunchAntiMalwareDriverEnabled = value
}
// SetOsMaximumVersion sets the osMaximumVersion property value. Maximum Windows Phone version.
func (m *Windows10MobileCompliancePolicy) SetOsMaximumVersion(value *string)() {
    m.osMaximumVersion = value
}
// SetOsMinimumVersion sets the osMinimumVersion property value. Minimum Windows Phone version.
func (m *Windows10MobileCompliancePolicy) SetOsMinimumVersion(value *string)() {
    m.osMinimumVersion = value
}
// SetPasswordBlockSimple sets the passwordBlockSimple property value. Whether or not to block syncing the calendar.
func (m *Windows10MobileCompliancePolicy) SetPasswordBlockSimple(value *bool)() {
    m.passwordBlockSimple = value
}
// SetPasswordExpirationDays sets the passwordExpirationDays property value. Number of days before password expiration. Valid values 1 to 255
func (m *Windows10MobileCompliancePolicy) SetPasswordExpirationDays(value *int32)() {
    m.passwordExpirationDays = value
}
// SetPasswordMinimumCharacterSetCount sets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
func (m *Windows10MobileCompliancePolicy) SetPasswordMinimumCharacterSetCount(value *int32)() {
    m.passwordMinimumCharacterSetCount = value
}
// SetPasswordMinimumLength sets the passwordMinimumLength property value. Minimum password length. Valid values 4 to 16
func (m *Windows10MobileCompliancePolicy) SetPasswordMinimumLength(value *int32)() {
    m.passwordMinimumLength = value
}
// SetPasswordMinutesOfInactivityBeforeLock sets the passwordMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a password is required.
func (m *Windows10MobileCompliancePolicy) SetPasswordMinutesOfInactivityBeforeLock(value *int32)() {
    m.passwordMinutesOfInactivityBeforeLock = value
}
// SetPasswordPreviousPasswordBlockCount sets the passwordPreviousPasswordBlockCount property value. The number of previous passwords to prevent re-use of.
func (m *Windows10MobileCompliancePolicy) SetPasswordPreviousPasswordBlockCount(value *int32)() {
    m.passwordPreviousPasswordBlockCount = value
}
// SetPasswordRequired sets the passwordRequired property value. Require a password to unlock Windows Phone device.
func (m *Windows10MobileCompliancePolicy) SetPasswordRequired(value *bool)() {
    m.passwordRequired = value
}
// SetPasswordRequiredType sets the passwordRequiredType property value. Possible values of required passwords.
func (m *Windows10MobileCompliancePolicy) SetPasswordRequiredType(value *RequiredPasswordType)() {
    m.passwordRequiredType = value
}
// SetPasswordRequireToUnlockFromIdle sets the passwordRequireToUnlockFromIdle property value. Require a password to unlock an idle device.
func (m *Windows10MobileCompliancePolicy) SetPasswordRequireToUnlockFromIdle(value *bool)() {
    m.passwordRequireToUnlockFromIdle = value
}
// SetSecureBootEnabled sets the secureBootEnabled property value. Require devices to be reported as healthy by Windows Device Health Attestation - secure boot is enabled.
func (m *Windows10MobileCompliancePolicy) SetSecureBootEnabled(value *bool)() {
    m.secureBootEnabled = value
}
// SetStorageRequireEncryption sets the storageRequireEncryption property value. Require encryption on windows devices.
func (m *Windows10MobileCompliancePolicy) SetStorageRequireEncryption(value *bool)() {
    m.storageRequireEncryption = value
}
