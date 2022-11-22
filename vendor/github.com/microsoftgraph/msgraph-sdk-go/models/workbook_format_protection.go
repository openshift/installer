package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookFormatProtection 
type WorkbookFormatProtection struct {
    Entity
    // Indicates if Excel hides the formula for the cells in the range. A null value indicates that the entire range doesn't have uniform formula hidden setting.
    formulaHidden *bool
    // Indicates if Excel locks the cells in the object. A null value indicates that the entire range doesn't have uniform lock setting.
    locked *bool
}
// NewWorkbookFormatProtection instantiates a new WorkbookFormatProtection and sets the default values.
func NewWorkbookFormatProtection()(*WorkbookFormatProtection) {
    m := &WorkbookFormatProtection{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookFormatProtectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookFormatProtectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookFormatProtection(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookFormatProtection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["formulaHidden"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetFormulaHidden)
    res["locked"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetLocked)
    return res
}
// GetFormulaHidden gets the formulaHidden property value. Indicates if Excel hides the formula for the cells in the range. A null value indicates that the entire range doesn't have uniform formula hidden setting.
func (m *WorkbookFormatProtection) GetFormulaHidden()(*bool) {
    return m.formulaHidden
}
// GetLocked gets the locked property value. Indicates if Excel locks the cells in the object. A null value indicates that the entire range doesn't have uniform lock setting.
func (m *WorkbookFormatProtection) GetLocked()(*bool) {
    return m.locked
}
// Serialize serializes information the current object
func (m *WorkbookFormatProtection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("formulaHidden", m.GetFormulaHidden())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("locked", m.GetLocked())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormulaHidden sets the formulaHidden property value. Indicates if Excel hides the formula for the cells in the range. A null value indicates that the entire range doesn't have uniform formula hidden setting.
func (m *WorkbookFormatProtection) SetFormulaHidden(value *bool)() {
    m.formulaHidden = value
}
// SetLocked sets the locked property value. Indicates if Excel locks the cells in the object. A null value indicates that the entire range doesn't have uniform lock setting.
func (m *WorkbookFormatProtection) SetLocked(value *bool)() {
    m.locked = value
}
