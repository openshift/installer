package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MacOSCompliancePolicy 
type MacOSCompliancePolicy struct {
    DeviceCompliancePolicy
    // Require that devices have enabled device threat protection.
    deviceThreatProtectionEnabled *bool
    // Device threat protection levels for the Device Threat Protection API.
    deviceThreatProtectionRequiredSecurityLevel *DeviceThreatProtectionLevel
    // Corresponds to the 'Block all incoming connections' option.
    firewallBlockAllIncoming *bool
    // Whether the firewall should be enabled or not.
    firewallEnabled *bool
    // Corresponds to 'Enable stealth mode.'
    firewallEnableStealthMode *bool
    // Maximum MacOS version.
    osMaximumVersion *string
    // Minimum MacOS version.
    osMinimumVersion *string
    // Indicates whether or not to block simple passwords.
    passwordBlockSimple *bool
    // Number of days before the password expires. Valid values 1 to 65535
    passwordExpirationDays *int32
    // The number of character sets required in the password.
    passwordMinimumCharacterSetCount *int32
    // Minimum length of password. Valid values 4 to 14
    passwordMinimumLength *int32
    // Minutes of inactivity before a password is required.
    passwordMinutesOfInactivityBeforeLock *int32
    // Number of previous passwords to block. Valid values 1 to 24
    passwordPreviousPasswordBlockCount *int32
    // Whether or not to require a password.
    passwordRequired *bool
    // Possible values of required passwords.
    passwordRequiredType *RequiredPasswordType
    // Require encryption on Mac OS devices.
    storageRequireEncryption *bool
    // Require that devices have enabled system integrity protection.
    systemIntegrityProtectionEnabled *bool
}
// NewMacOSCompliancePolicy instantiates a new MacOSCompliancePolicy and sets the default values.
func NewMacOSCompliancePolicy()(*MacOSCompliancePolicy) {
    m := &MacOSCompliancePolicy{
        DeviceCompliancePolicy: *NewDeviceCompliancePolicy(),
    }
    odataTypeValue := "#microsoft.graph.macOSCompliancePolicy";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateMacOSCompliancePolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMacOSCompliancePolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMacOSCompliancePolicy(), nil
}
// GetDeviceThreatProtectionEnabled gets the deviceThreatProtectionEnabled property value. Require that devices have enabled device threat protection.
func (m *MacOSCompliancePolicy) GetDeviceThreatProtectionEnabled()(*bool) {
    return m.deviceThreatProtectionEnabled
}
// GetDeviceThreatProtectionRequiredSecurityLevel gets the deviceThreatProtectionRequiredSecurityLevel property value. Device threat protection levels for the Device Threat Protection API.
func (m *MacOSCompliancePolicy) GetDeviceThreatProtectionRequiredSecurityLevel()(*DeviceThreatProtectionLevel) {
    return m.deviceThreatProtectionRequiredSecurityLevel
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MacOSCompliancePolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceCompliancePolicy.GetFieldDeserializers()
    res["deviceThreatProtectionEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetDeviceThreatProtectionEnabled)
    res["deviceThreatProtectionRequiredSecurityLevel"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceThreatProtectionLevel , m.SetDeviceThreatProtectionRequiredSecurityLevel)
    res["firewallBlockAllIncoming"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetFirewallBlockAllIncoming)
    res["firewallEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetFirewallEnabled)
    res["firewallEnableStealthMode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetFirewallEnableStealthMode)
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
    res["storageRequireEncryption"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetStorageRequireEncryption)
    res["systemIntegrityProtectionEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetSystemIntegrityProtectionEnabled)
    return res
}
// GetFirewallBlockAllIncoming gets the firewallBlockAllIncoming property value. Corresponds to the 'Block all incoming connections' option.
func (m *MacOSCompliancePolicy) GetFirewallBlockAllIncoming()(*bool) {
    return m.firewallBlockAllIncoming
}
// GetFirewallEnabled gets the firewallEnabled property value. Whether the firewall should be enabled or not.
func (m *MacOSCompliancePolicy) GetFirewallEnabled()(*bool) {
    return m.firewallEnabled
}
// GetFirewallEnableStealthMode gets the firewallEnableStealthMode property value. Corresponds to 'Enable stealth mode.'
func (m *MacOSCompliancePolicy) GetFirewallEnableStealthMode()(*bool) {
    return m.firewallEnableStealthMode
}
// GetOsMaximumVersion gets the osMaximumVersion property value. Maximum MacOS version.
func (m *MacOSCompliancePolicy) GetOsMaximumVersion()(*string) {
    return m.osMaximumVersion
}
// GetOsMinimumVersion gets the osMinimumVersion property value. Minimum MacOS version.
func (m *MacOSCompliancePolicy) GetOsMinimumVersion()(*string) {
    return m.osMinimumVersion
}
// GetPasswordBlockSimple gets the passwordBlockSimple property value. Indicates whether or not to block simple passwords.
func (m *MacOSCompliancePolicy) GetPasswordBlockSimple()(*bool) {
    return m.passwordBlockSimple
}
// GetPasswordExpirationDays gets the passwordExpirationDays property value. Number of days before the password expires. Valid values 1 to 65535
func (m *MacOSCompliancePolicy) GetPasswordExpirationDays()(*int32) {
    return m.passwordExpirationDays
}
// GetPasswordMinimumCharacterSetCount gets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
func (m *MacOSCompliancePolicy) GetPasswordMinimumCharacterSetCount()(*int32) {
    return m.passwordMinimumCharacterSetCount
}
// GetPasswordMinimumLength gets the passwordMinimumLength property value. Minimum length of password. Valid values 4 to 14
func (m *MacOSCompliancePolicy) GetPasswordMinimumLength()(*int32) {
    return m.passwordMinimumLength
}
// GetPasswordMinutesOfInactivityBeforeLock gets the passwordMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a password is required.
func (m *MacOSCompliancePolicy) GetPasswordMinutesOfInactivityBeforeLock()(*int32) {
    return m.passwordMinutesOfInactivityBeforeLock
}
// GetPasswordPreviousPasswordBlockCount gets the passwordPreviousPasswordBlockCount property value. Number of previous passwords to block. Valid values 1 to 24
func (m *MacOSCompliancePolicy) GetPasswordPreviousPasswordBlockCount()(*int32) {
    return m.passwordPreviousPasswordBlockCount
}
// GetPasswordRequired gets the passwordRequired property value. Whether or not to require a password.
func (m *MacOSCompliancePolicy) GetPasswordRequired()(*bool) {
    return m.passwordRequired
}
// GetPasswordRequiredType gets the passwordRequiredType property value. Possible values of required passwords.
func (m *MacOSCompliancePolicy) GetPasswordRequiredType()(*RequiredPasswordType) {
    return m.passwordRequiredType
}
// GetStorageRequireEncryption gets the storageRequireEncryption property value. Require encryption on Mac OS devices.
func (m *MacOSCompliancePolicy) GetStorageRequireEncryption()(*bool) {
    return m.storageRequireEncryption
}
// GetSystemIntegrityProtectionEnabled gets the systemIntegrityProtectionEnabled property value. Require that devices have enabled system integrity protection.
func (m *MacOSCompliancePolicy) GetSystemIntegrityProtectionEnabled()(*bool) {
    return m.systemIntegrityProtectionEnabled
}
// Serialize serializes information the current object
func (m *MacOSCompliancePolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceCompliancePolicy.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("deviceThreatProtectionEnabled", m.GetDeviceThreatProtectionEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceThreatProtectionRequiredSecurityLevel() != nil {
        cast := (*m.GetDeviceThreatProtectionRequiredSecurityLevel()).String()
        err = writer.WriteStringValue("deviceThreatProtectionRequiredSecurityLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallBlockAllIncoming", m.GetFirewallBlockAllIncoming())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallEnabled", m.GetFirewallEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("firewallEnableStealthMode", m.GetFirewallEnableStealthMode())
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
        err = writer.WriteBoolValue("storageRequireEncryption", m.GetStorageRequireEncryption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("systemIntegrityProtectionEnabled", m.GetSystemIntegrityProtectionEnabled())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceThreatProtectionEnabled sets the deviceThreatProtectionEnabled property value. Require that devices have enabled device threat protection.
func (m *MacOSCompliancePolicy) SetDeviceThreatProtectionEnabled(value *bool)() {
    m.deviceThreatProtectionEnabled = value
}
// SetDeviceThreatProtectionRequiredSecurityLevel sets the deviceThreatProtectionRequiredSecurityLevel property value. Device threat protection levels for the Device Threat Protection API.
func (m *MacOSCompliancePolicy) SetDeviceThreatProtectionRequiredSecurityLevel(value *DeviceThreatProtectionLevel)() {
    m.deviceThreatProtectionRequiredSecurityLevel = value
}
// SetFirewallBlockAllIncoming sets the firewallBlockAllIncoming property value. Corresponds to the 'Block all incoming connections' option.
func (m *MacOSCompliancePolicy) SetFirewallBlockAllIncoming(value *bool)() {
    m.firewallBlockAllIncoming = value
}
// SetFirewallEnabled sets the firewallEnabled property value. Whether the firewall should be enabled or not.
func (m *MacOSCompliancePolicy) SetFirewallEnabled(value *bool)() {
    m.firewallEnabled = value
}
// SetFirewallEnableStealthMode sets the firewallEnableStealthMode property value. Corresponds to 'Enable stealth mode.'
func (m *MacOSCompliancePolicy) SetFirewallEnableStealthMode(value *bool)() {
    m.firewallEnableStealthMode = value
}
// SetOsMaximumVersion sets the osMaximumVersion property value. Maximum MacOS version.
func (m *MacOSCompliancePolicy) SetOsMaximumVersion(value *string)() {
    m.osMaximumVersion = value
}
// SetOsMinimumVersion sets the osMinimumVersion property value. Minimum MacOS version.
func (m *MacOSCompliancePolicy) SetOsMinimumVersion(value *string)() {
    m.osMinimumVersion = value
}
// SetPasswordBlockSimple sets the passwordBlockSimple property value. Indicates whether or not to block simple passwords.
func (m *MacOSCompliancePolicy) SetPasswordBlockSimple(value *bool)() {
    m.passwordBlockSimple = value
}
// SetPasswordExpirationDays sets the passwordExpirationDays property value. Number of days before the password expires. Valid values 1 to 65535
func (m *MacOSCompliancePolicy) SetPasswordExpirationDays(value *int32)() {
    m.passwordExpirationDays = value
}
// SetPasswordMinimumCharacterSetCount sets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
func (m *MacOSCompliancePolicy) SetPasswordMinimumCharacterSetCount(value *int32)() {
    m.passwordMinimumCharacterSetCount = value
}
// SetPasswordMinimumLength sets the passwordMinimumLength property value. Minimum length of password. Valid values 4 to 14
func (m *MacOSCompliancePolicy) SetPasswordMinimumLength(value *int32)() {
    m.passwordMinimumLength = value
}
// SetPasswordMinutesOfInactivityBeforeLock sets the passwordMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a password is required.
func (m *MacOSCompliancePolicy) SetPasswordMinutesOfInactivityBeforeLock(value *int32)() {
    m.passwordMinutesOfInactivityBeforeLock = value
}
// SetPasswordPreviousPasswordBlockCount sets the passwordPreviousPasswordBlockCount property value. Number of previous passwords to block. Valid values 1 to 24
func (m *MacOSCompliancePolicy) SetPasswordPreviousPasswordBlockCount(value *int32)() {
    m.passwordPreviousPasswordBlockCount = value
}
// SetPasswordRequired sets the passwordRequired property value. Whether or not to require a password.
func (m *MacOSCompliancePolicy) SetPasswordRequired(value *bool)() {
    m.passwordRequired = value
}
// SetPasswordRequiredType sets the passwordRequiredType property value. Possible values of required passwords.
func (m *MacOSCompliancePolicy) SetPasswordRequiredType(value *RequiredPasswordType)() {
    m.passwordRequiredType = value
}
// SetStorageRequireEncryption sets the storageRequireEncryption property value. Require encryption on Mac OS devices.
func (m *MacOSCompliancePolicy) SetStorageRequireEncryption(value *bool)() {
    m.storageRequireEncryption = value
}
// SetSystemIntegrityProtectionEnabled sets the systemIntegrityProtectionEnabled property value. Require that devices have enabled system integrity protection.
func (m *MacOSCompliancePolicy) SetSystemIntegrityProtectionEnabled(value *bool)() {
    m.systemIntegrityProtectionEnabled = value
}
