package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookWorksheet 
type WorkbookWorksheet struct {
    Entity
    // Returns collection of charts that are part of the worksheet. Read-only.
    charts []WorkbookChartable
    // The display name of the worksheet.
    name *string
    // Returns collection of names that are associated with the worksheet. Read-only.
    names []WorkbookNamedItemable
    // Collection of PivotTables that are part of the worksheet.
    pivotTables []WorkbookPivotTableable
    // The zero-based position of the worksheet within the workbook.
    position *int32
    // Returns sheet protection object for a worksheet. Read-only.
    protection WorkbookWorksheetProtectionable
    // Collection of tables that are part of the worksheet. Read-only.
    tables []WorkbookTableable
    // The Visibility of the worksheet. The possible values are: Visible, Hidden, VeryHidden.
    visibility *string
}
// NewWorkbookWorksheet instantiates a new workbookWorksheet and sets the default values.
func NewWorkbookWorksheet()(*WorkbookWorksheet) {
    m := &WorkbookWorksheet{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookWorksheetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookWorksheetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookWorksheet(), nil
}
// GetCharts gets the charts property value. Returns collection of charts that are part of the worksheet. Read-only.
func (m *WorkbookWorksheet) GetCharts()([]WorkbookChartable) {
    return m.charts
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookWorksheet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["charts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookChartFromDiscriminatorValue , m.SetCharts)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["names"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookNamedItemFromDiscriminatorValue , m.SetNames)
    res["pivotTables"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookPivotTableFromDiscriminatorValue , m.SetPivotTables)
    res["position"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPosition)
    res["protection"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookWorksheetProtectionFromDiscriminatorValue , m.SetProtection)
    res["tables"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookTableFromDiscriminatorValue , m.SetTables)
    res["visibility"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVisibility)
    return res
}
// GetName gets the name property value. The display name of the worksheet.
func (m *WorkbookWorksheet) GetName()(*string) {
    return m.name
}
// GetNames gets the names property value. Returns collection of names that are associated with the worksheet. Read-only.
func (m *WorkbookWorksheet) GetNames()([]WorkbookNamedItemable) {
    return m.names
}
// GetPivotTables gets the pivotTables property value. Collection of PivotTables that are part of the worksheet.
func (m *WorkbookWorksheet) GetPivotTables()([]WorkbookPivotTableable) {
    return m.pivotTables
}
// GetPosition gets the position property value. The zero-based position of the worksheet within the workbook.
func (m *WorkbookWorksheet) GetPosition()(*int32) {
    return m.position
}
// GetProtection gets the protection property value. Returns sheet protection object for a worksheet. Read-only.
func (m *WorkbookWorksheet) GetProtection()(WorkbookWorksheetProtectionable) {
    return m.protection
}
// GetTables gets the tables property value. Collection of tables that are part of the worksheet. Read-only.
func (m *WorkbookWorksheet) GetTables()([]WorkbookTableable) {
    return m.tables
}
// GetVisibility gets the visibility property value. The Visibility of the worksheet. The possible values are: Visible, Hidden, VeryHidden.
func (m *WorkbookWorksheet) GetVisibility()(*string) {
    return m.visibility
}
// Serialize serializes information the current object
func (m *WorkbookWorksheet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCharts() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCharts())
        err = writer.WriteCollectionOfObjectValues("charts", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetNames() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNames())
        err = writer.WriteCollectionOfObjectValues("names", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPivotTables() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPivotTables())
        err = writer.WriteCollectionOfObjectValues("pivotTables", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("position", m.GetPosition())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("protection", m.GetProtection())
        if err != nil {
            return err
        }
    }
    if m.GetTables() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTables())
        err = writer.WriteCollectionOfObjectValues("tables", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("visibility", m.GetVisibility())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCharts sets the charts property value. Returns collection of charts that are part of the worksheet. Read-only.
func (m *WorkbookWorksheet) SetCharts(value []WorkbookChartable)() {
    m.charts = value
}
// SetName sets the name property value. The display name of the worksheet.
func (m *WorkbookWorksheet) SetName(value *string)() {
    m.name = value
}
// SetNames sets the names property value. Returns collection of names that are associated with the worksheet. Read-only.
func (m *WorkbookWorksheet) SetNames(value []WorkbookNamedItemable)() {
    m.names = value
}
// SetPivotTables sets the pivotTables property value. Collection of PivotTables that are part of the worksheet.
func (m *WorkbookWorksheet) SetPivotTables(value []WorkbookPivotTableable)() {
    m.pivotTables = value
}
// SetPosition sets the position property value. The zero-based position of the worksheet within the workbook.
func (m *WorkbookWorksheet) SetPosition(value *int32)() {
    m.position = value
}
// SetProtection sets the protection property value. Returns sheet protection object for a worksheet. Read-only.
func (m *WorkbookWorksheet) SetProtection(value WorkbookWorksheetProtectionable)() {
    m.protection = value
}
// SetTables sets the tables property value. Collection of tables that are part of the worksheet. Read-only.
func (m *WorkbookWorksheet) SetTables(value []WorkbookTableable)() {
    m.tables = value
}
// SetVisibility sets the visibility property value. The Visibility of the worksheet. The possible values are: Visible, Hidden, VeryHidden.
func (m *WorkbookWorksheet) SetVisibility(value *string)() {
    m.visibility = value
}
