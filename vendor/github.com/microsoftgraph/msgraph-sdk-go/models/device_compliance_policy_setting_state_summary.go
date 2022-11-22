package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceCompliancePolicySettingStateSummary device Compilance Policy Setting State summary across the account.
type DeviceCompliancePolicySettingStateSummary struct {
    Entity
    // Number of compliant devices
    compliantDeviceCount *int32
    // Number of conflict devices
    conflictDeviceCount *int32
    // Not yet documented
    deviceComplianceSettingStates []DeviceComplianceSettingStateable
    // Number of error devices
    errorDeviceCount *int32
    // Number of NonCompliant devices
    nonCompliantDeviceCount *int32
    // Number of not applicable devices
    notApplicableDeviceCount *int32
    // Supported platform types for policies.
    platformType *PolicyPlatformType
    // Number of remediated devices
    remediatedDeviceCount *int32
    // The setting class name and property name.
    setting *string
    // Name of the setting.
    settingName *string
    // Number of unknown devices
    unknownDeviceCount *int32
}
// NewDeviceCompliancePolicySettingStateSummary instantiates a new deviceCompliancePolicySettingStateSummary and sets the default values.
func NewDeviceCompliancePolicySettingStateSummary()(*DeviceCompliancePolicySettingStateSummary) {
    m := &DeviceCompliancePolicySettingStateSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceCompliancePolicySettingStateSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDeviceCompliancePolicySettingStateSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceCompliancePolicySettingStateSummary(), nil
}
// GetCompliantDeviceCount gets the compliantDeviceCount property value. Number of compliant devices
func (m *DeviceCompliancePolicySettingStateSummary) GetCompliantDeviceCount()(*int32) {
    return m.compliantDeviceCount
}
// GetConflictDeviceCount gets the conflictDeviceCount property value. Number of conflict devices
func (m *DeviceCompliancePolicySettingStateSummary) GetConflictDeviceCount()(*int32) {
    return m.conflictDeviceCount
}
// GetDeviceComplianceSettingStates gets the deviceComplianceSettingStates property value. Not yet documented
func (m *DeviceCompliancePolicySettingStateSummary) GetDeviceComplianceSettingStates()([]DeviceComplianceSettingStateable) {
    return m.deviceComplianceSettingStates
}
// GetErrorDeviceCount gets the errorDeviceCount property value. Number of error devices
func (m *DeviceCompliancePolicySettingStateSummary) GetErrorDeviceCount()(*int32) {
    return m.errorDeviceCount
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceCompliancePolicySettingStateSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["compliantDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetCompliantDeviceCount)
    res["conflictDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetConflictDeviceCount)
    res["deviceComplianceSettingStates"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceComplianceSettingStateFromDiscriminatorValue , m.SetDeviceComplianceSettingStates)
    res["errorDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetErrorDeviceCount)
    res["nonCompliantDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetNonCompliantDeviceCount)
    res["notApplicableDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetNotApplicableDeviceCount)
    res["platformType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePolicyPlatformType , m.SetPlatformType)
    res["remediatedDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetRemediatedDeviceCount)
    res["setting"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSetting)
    res["settingName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSettingName)
    res["unknownDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetUnknownDeviceCount)
    return res
}
// GetNonCompliantDeviceCount gets the nonCompliantDeviceCount property value. Number of NonCompliant devices
func (m *DeviceCompliancePolicySettingStateSummary) GetNonCompliantDeviceCount()(*int32) {
    return m.nonCompliantDeviceCount
}
// GetNotApplicableDeviceCount gets the notApplicableDeviceCount property value. Number of not applicable devices
func (m *DeviceCompliancePolicySettingStateSummary) GetNotApplicableDeviceCount()(*int32) {
    return m.notApplicableDeviceCount
}
// GetPlatformType gets the platformType property value. Supported platform types for policies.
func (m *DeviceCompliancePolicySettingStateSummary) GetPlatformType()(*PolicyPlatformType) {
    return m.platformType
}
// GetRemediatedDeviceCount gets the remediatedDeviceCount property value. Number of remediated devices
func (m *DeviceCompliancePolicySettingStateSummary) GetRemediatedDeviceCount()(*int32) {
    return m.remediatedDeviceCount
}
// GetSetting gets the setting property value. The setting class name and property name.
func (m *DeviceCompliancePolicySettingStateSummary) GetSetting()(*string) {
    return m.setting
}
// GetSettingName gets the settingName property value. Name of the setting.
func (m *DeviceCompliancePolicySettingStateSummary) GetSettingName()(*string) {
    return m.settingName
}
// GetUnknownDeviceCount gets the unknownDeviceCount property value. Number of unknown devices
func (m *DeviceCompliancePolicySettingStateSummary) GetUnknownDeviceCount()(*int32) {
    return m.unknownDeviceCount
}
// Serialize serializes information the current object
func (m *DeviceCompliancePolicySettingStateSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("compliantDeviceCount", m.GetCompliantDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("conflictDeviceCount", m.GetConflictDeviceCount())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceComplianceSettingStates() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceComplianceSettingStates())
        err = writer.WriteCollectionOfObjectValues("deviceComplianceSettingStates", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("errorDeviceCount", m.GetErrorDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("nonCompliantDeviceCount", m.GetNonCompliantDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("notApplicableDeviceCount", m.GetNotApplicableDeviceCount())
        if err != nil {
            return err
        }
    }
    if m.GetPlatformType() != nil {
        cast := (*m.GetPlatformType()).String()
        err = writer.WriteStringValue("platformType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("remediatedDeviceCount", m.GetRemediatedDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("setting", m.GetSetting())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("settingName", m.GetSettingName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("unknownDeviceCount", m.GetUnknownDeviceCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompliantDeviceCount sets the compliantDeviceCount property value. Number of compliant devices
func (m *DeviceCompliancePolicySettingStateSummary) SetCompliantDeviceCount(value *int32)() {
    m.compliantDeviceCount = value
}
// SetConflictDeviceCount sets the conflictDeviceCount property value. Number of conflict devices
func (m *DeviceCompliancePolicySettingStateSummary) SetConflictDeviceCount(value *int32)() {
    m.conflictDeviceCount = value
}
// SetDeviceComplianceSettingStates sets the deviceComplianceSettingStates property value. Not yet documented
func (m *DeviceCompliancePolicySettingStateSummary) SetDeviceComplianceSettingStates(value []DeviceComplianceSettingStateable)() {
    m.deviceComplianceSettingStates = value
}
// SetErrorDeviceCount sets the errorDeviceCount property value. Number of error devices
func (m *DeviceCompliancePolicySettingStateSummary) SetErrorDeviceCount(value *int32)() {
    m.errorDeviceCount = value
}
// SetNonCompliantDeviceCount sets the nonCompliantDeviceCount property value. Number of NonCompliant devices
func (m *DeviceCompliancePolicySettingStateSummary) SetNonCompliantDeviceCount(value *int32)() {
    m.nonCompliantDeviceCount = value
}
// SetNotApplicableDeviceCount sets the notApplicableDeviceCount property value. Number of not applicable devices
func (m *DeviceCompliancePolicySettingStateSummary) SetNotApplicableDeviceCount(value *int32)() {
    m.notApplicableDeviceCount = value
}
// SetPlatformType sets the platformType property value. Supported platform types for policies.
func (m *DeviceCompliancePolicySettingStateSummary) SetPlatformType(value *PolicyPlatformType)() {
    m.platformType = value
}
// SetRemediatedDeviceCount sets the remediatedDeviceCount property value. Number of remediated devices
func (m *DeviceCompliancePolicySettingStateSummary) SetRemediatedDeviceCount(value *int32)() {
    m.remediatedDeviceCount = value
}
// SetSetting sets the setting property value. The setting class name and property name.
func (m *DeviceCompliancePolicySettingStateSummary) SetSetting(value *string)() {
    m.setting = value
}
// SetSettingName sets the settingName property value. Name of the setting.
func (m *DeviceCompliancePolicySettingStateSummary) SetSettingName(value *string)() {
    m.settingName = value
}
// SetUnknownDeviceCount sets the unknownDeviceCount property value. Number of unknown devices
func (m *DeviceCompliancePolicySettingStateSummary) SetUnknownDeviceCount(value *int32)() {
    m.unknownDeviceCount = value
}
