package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SettingStateDeviceSummary device Compilance Policy and Configuration for a Setting State summary
type SettingStateDeviceSummary struct {
    Entity
    // Device Compliant count for the setting
    compliantDeviceCount *int32
    // Device conflict error count for the setting
    conflictDeviceCount *int32
    // Device error count for the setting
    errorDeviceCount *int32
    // Name of the InstancePath for the setting
    instancePath *string
    // Device NonCompliant count for the setting
    nonCompliantDeviceCount *int32
    // Device Not Applicable count for the setting
    notApplicableDeviceCount *int32
    // Device Compliant count for the setting
    remediatedDeviceCount *int32
    // Name of the setting
    settingName *string
    // Device Unkown count for the setting
    unknownDeviceCount *int32
}
// NewSettingStateDeviceSummary instantiates a new settingStateDeviceSummary and sets the default values.
func NewSettingStateDeviceSummary()(*SettingStateDeviceSummary) {
    m := &SettingStateDeviceSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSettingStateDeviceSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSettingStateDeviceSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSettingStateDeviceSummary(), nil
}
// GetCompliantDeviceCount gets the compliantDeviceCount property value. Device Compliant count for the setting
func (m *SettingStateDeviceSummary) GetCompliantDeviceCount()(*int32) {
    return m.compliantDeviceCount
}
// GetConflictDeviceCount gets the conflictDeviceCount property value. Device conflict error count for the setting
func (m *SettingStateDeviceSummary) GetConflictDeviceCount()(*int32) {
    return m.conflictDeviceCount
}
// GetErrorDeviceCount gets the errorDeviceCount property value. Device error count for the setting
func (m *SettingStateDeviceSummary) GetErrorDeviceCount()(*int32) {
    return m.errorDeviceCount
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SettingStateDeviceSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["compliantDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetCompliantDeviceCount)
    res["conflictDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetConflictDeviceCount)
    res["errorDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetErrorDeviceCount)
    res["instancePath"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInstancePath)
    res["nonCompliantDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetNonCompliantDeviceCount)
    res["notApplicableDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetNotApplicableDeviceCount)
    res["remediatedDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetRemediatedDeviceCount)
    res["settingName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSettingName)
    res["unknownDeviceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetUnknownDeviceCount)
    return res
}
// GetInstancePath gets the instancePath property value. Name of the InstancePath for the setting
func (m *SettingStateDeviceSummary) GetInstancePath()(*string) {
    return m.instancePath
}
// GetNonCompliantDeviceCount gets the nonCompliantDeviceCount property value. Device NonCompliant count for the setting
func (m *SettingStateDeviceSummary) GetNonCompliantDeviceCount()(*int32) {
    return m.nonCompliantDeviceCount
}
// GetNotApplicableDeviceCount gets the notApplicableDeviceCount property value. Device Not Applicable count for the setting
func (m *SettingStateDeviceSummary) GetNotApplicableDeviceCount()(*int32) {
    return m.notApplicableDeviceCount
}
// GetRemediatedDeviceCount gets the remediatedDeviceCount property value. Device Compliant count for the setting
func (m *SettingStateDeviceSummary) GetRemediatedDeviceCount()(*int32) {
    return m.remediatedDeviceCount
}
// GetSettingName gets the settingName property value. Name of the setting
func (m *SettingStateDeviceSummary) GetSettingName()(*string) {
    return m.settingName
}
// GetUnknownDeviceCount gets the unknownDeviceCount property value. Device Unkown count for the setting
func (m *SettingStateDeviceSummary) GetUnknownDeviceCount()(*int32) {
    return m.unknownDeviceCount
}
// Serialize serializes information the current object
func (m *SettingStateDeviceSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    {
        err = writer.WriteInt32Value("errorDeviceCount", m.GetErrorDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("instancePath", m.GetInstancePath())
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
    {
        err = writer.WriteInt32Value("remediatedDeviceCount", m.GetRemediatedDeviceCount())
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
// SetCompliantDeviceCount sets the compliantDeviceCount property value. Device Compliant count for the setting
func (m *SettingStateDeviceSummary) SetCompliantDeviceCount(value *int32)() {
    m.compliantDeviceCount = value
}
// SetConflictDeviceCount sets the conflictDeviceCount property value. Device conflict error count for the setting
func (m *SettingStateDeviceSummary) SetConflictDeviceCount(value *int32)() {
    m.conflictDeviceCount = value
}
// SetErrorDeviceCount sets the errorDeviceCount property value. Device error count for the setting
func (m *SettingStateDeviceSummary) SetErrorDeviceCount(value *int32)() {
    m.errorDeviceCount = value
}
// SetInstancePath sets the instancePath property value. Name of the InstancePath for the setting
func (m *SettingStateDeviceSummary) SetInstancePath(value *string)() {
    m.instancePath = value
}
// SetNonCompliantDeviceCount sets the nonCompliantDeviceCount property value. Device NonCompliant count for the setting
func (m *SettingStateDeviceSummary) SetNonCompliantDeviceCount(value *int32)() {
    m.nonCompliantDeviceCount = value
}
// SetNotApplicableDeviceCount sets the notApplicableDeviceCount property value. Device Not Applicable count for the setting
func (m *SettingStateDeviceSummary) SetNotApplicableDeviceCount(value *int32)() {
    m.notApplicableDeviceCount = value
}
// SetRemediatedDeviceCount sets the remediatedDeviceCount property value. Device Compliant count for the setting
func (m *SettingStateDeviceSummary) SetRemediatedDeviceCount(value *int32)() {
    m.remediatedDeviceCount = value
}
// SetSettingName sets the settingName property value. Name of the setting
func (m *SettingStateDeviceSummary) SetSettingName(value *string)() {
    m.settingName = value
}
// SetUnknownDeviceCount sets the unknownDeviceCount property value. Device Unkown count for the setting
func (m *SettingStateDeviceSummary) SetUnknownDeviceCount(value *int32)() {
    m.unknownDeviceCount = value
}
