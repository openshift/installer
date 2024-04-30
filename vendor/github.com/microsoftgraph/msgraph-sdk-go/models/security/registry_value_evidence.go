package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RegistryValueEvidence 
type RegistryValueEvidence struct {
    AlertEvidence
}
// NewRegistryValueEvidence instantiates a new RegistryValueEvidence and sets the default values.
func NewRegistryValueEvidence()(*RegistryValueEvidence) {
    m := &RegistryValueEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    return m
}
// CreateRegistryValueEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRegistryValueEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRegistryValueEvidence(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RegistryValueEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["registryHive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistryHive(val)
        }
        return nil
    }
    res["registryKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistryKey(val)
        }
        return nil
    }
    res["registryValue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistryValue(val)
        }
        return nil
    }
    res["registryValueName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistryValueName(val)
        }
        return nil
    }
    res["registryValueType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistryValueType(val)
        }
        return nil
    }
    return res
}
// GetRegistryHive gets the registryHive property value. Registry hive of the key that the recorded action was applied to.
func (m *RegistryValueEvidence) GetRegistryHive()(*string) {
    val, err := m.GetBackingStore().Get("registryHive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistryKey gets the registryKey property value. Registry key that the recorded action was applied to.
func (m *RegistryValueEvidence) GetRegistryKey()(*string) {
    val, err := m.GetBackingStore().Get("registryKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistryValue gets the registryValue property value. Data of the registry value that the recorded action was applied to.
func (m *RegistryValueEvidence) GetRegistryValue()(*string) {
    val, err := m.GetBackingStore().Get("registryValue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistryValueName gets the registryValueName property value. Name of the registry value that the recorded action was applied to.
func (m *RegistryValueEvidence) GetRegistryValueName()(*string) {
    val, err := m.GetBackingStore().Get("registryValueName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistryValueType gets the registryValueType property value. Data type, such as binary or string, of the registry value that the recorded action was applied to.
func (m *RegistryValueEvidence) GetRegistryValueType()(*string) {
    val, err := m.GetBackingStore().Get("registryValueType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RegistryValueEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("registryHive", m.GetRegistryHive())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("registryKey", m.GetRegistryKey())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("registryValue", m.GetRegistryValue())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("registryValueName", m.GetRegistryValueName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("registryValueType", m.GetRegistryValueType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRegistryHive sets the registryHive property value. Registry hive of the key that the recorded action was applied to.
func (m *RegistryValueEvidence) SetRegistryHive(value *string)() {
    err := m.GetBackingStore().Set("registryHive", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistryKey sets the registryKey property value. Registry key that the recorded action was applied to.
func (m *RegistryValueEvidence) SetRegistryKey(value *string)() {
    err := m.GetBackingStore().Set("registryKey", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistryValue sets the registryValue property value. Data of the registry value that the recorded action was applied to.
func (m *RegistryValueEvidence) SetRegistryValue(value *string)() {
    err := m.GetBackingStore().Set("registryValue", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistryValueName sets the registryValueName property value. Name of the registry value that the recorded action was applied to.
func (m *RegistryValueEvidence) SetRegistryValueName(value *string)() {
    err := m.GetBackingStore().Set("registryValueName", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistryValueType sets the registryValueType property value. Data type, such as binary or string, of the registry value that the recorded action was applied to.
func (m *RegistryValueEvidence) SetRegistryValueType(value *string)() {
    err := m.GetBackingStore().Set("registryValueType", value)
    if err != nil {
        panic(err)
    }
}
// RegistryValueEvidenceable 
type RegistryValueEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRegistryHive()(*string)
    GetRegistryKey()(*string)
    GetRegistryValue()(*string)
    GetRegistryValueName()(*string)
    GetRegistryValueType()(*string)
    SetRegistryHive(value *string)()
    SetRegistryKey(value *string)()
    SetRegistryValue(value *string)()
    SetRegistryValueName(value *string)()
    SetRegistryValueType(value *string)()
}
