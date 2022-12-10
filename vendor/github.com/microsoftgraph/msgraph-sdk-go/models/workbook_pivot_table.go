package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookPivotTable provides operations to manage the collection of agreement entities.
type WorkbookPivotTable struct {
    Entity
    // Name of the PivotTable.
    name *string
    // The worksheet containing the current PivotTable. Read-only.
    worksheet WorkbookWorksheetable
}
// NewWorkbookPivotTable instantiates a new workbookPivotTable and sets the default values.
func NewWorkbookPivotTable()(*WorkbookPivotTable) {
    m := &WorkbookPivotTable{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookPivotTableFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookPivotTableFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookPivotTable(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookPivotTable) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["worksheet"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookWorksheetFromDiscriminatorValue , m.SetWorksheet)
    return res
}
// GetName gets the name property value. Name of the PivotTable.
func (m *WorkbookPivotTable) GetName()(*string) {
    return m.name
}
// GetWorksheet gets the worksheet property value. The worksheet containing the current PivotTable. Read-only.
func (m *WorkbookPivotTable) GetWorksheet()(WorkbookWorksheetable) {
    return m.worksheet
}
// Serialize serializes information the current object
func (m *WorkbookPivotTable) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("worksheet", m.GetWorksheet())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetName sets the name property value. Name of the PivotTable.
func (m *WorkbookPivotTable) SetName(value *string)() {
    m.name = value
}
// SetWorksheet sets the worksheet property value. The worksheet containing the current PivotTable. Read-only.
func (m *WorkbookPivotTable) SetWorksheet(value WorkbookWorksheetable)() {
    m.worksheet = value
}
