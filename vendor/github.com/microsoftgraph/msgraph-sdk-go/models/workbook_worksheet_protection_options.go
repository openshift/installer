package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookWorksheetProtectionOptions 
type WorkbookWorksheetProtectionOptions struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Represents the worksheet protection option of allowing using auto filter feature.
    allowAutoFilter *bool
    // Represents the worksheet protection option of allowing deleting columns.
    allowDeleteColumns *bool
    // Represents the worksheet protection option of allowing deleting rows.
    allowDeleteRows *bool
    // Represents the worksheet protection option of allowing formatting cells.
    allowFormatCells *bool
    // Represents the worksheet protection option of allowing formatting columns.
    allowFormatColumns *bool
    // Represents the worksheet protection option of allowing formatting rows.
    allowFormatRows *bool
    // Represents the worksheet protection option of allowing inserting columns.
    allowInsertColumns *bool
    // Represents the worksheet protection option of allowing inserting hyperlinks.
    allowInsertHyperlinks *bool
    // Represents the worksheet protection option of allowing inserting rows.
    allowInsertRows *bool
    // Represents the worksheet protection option of allowing using pivot table feature.
    allowPivotTables *bool
    // Represents the worksheet protection option of allowing using sort feature.
    allowSort *bool
    // The OdataType property
    odataType *string
}
// NewWorkbookWorksheetProtectionOptions instantiates a new workbookWorksheetProtectionOptions and sets the default values.
func NewWorkbookWorksheetProtectionOptions()(*WorkbookWorksheetProtectionOptions) {
    m := &WorkbookWorksheetProtectionOptions{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWorkbookWorksheetProtectionOptionsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookWorksheetProtectionOptionsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookWorksheetProtectionOptions(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *WorkbookWorksheetProtectionOptions) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAllowAutoFilter gets the allowAutoFilter property value. Represents the worksheet protection option of allowing using auto filter feature.
func (m *WorkbookWorksheetProtectionOptions) GetAllowAutoFilter()(*bool) {
    return m.allowAutoFilter
}
// GetAllowDeleteColumns gets the allowDeleteColumns property value. Represents the worksheet protection option of allowing deleting columns.
func (m *WorkbookWorksheetProtectionOptions) GetAllowDeleteColumns()(*bool) {
    return m.allowDeleteColumns
}
// GetAllowDeleteRows gets the allowDeleteRows property value. Represents the worksheet protection option of allowing deleting rows.
func (m *WorkbookWorksheetProtectionOptions) GetAllowDeleteRows()(*bool) {
    return m.allowDeleteRows
}
// GetAllowFormatCells gets the allowFormatCells property value. Represents the worksheet protection option of allowing formatting cells.
func (m *WorkbookWorksheetProtectionOptions) GetAllowFormatCells()(*bool) {
    return m.allowFormatCells
}
// GetAllowFormatColumns gets the allowFormatColumns property value. Represents the worksheet protection option of allowing formatting columns.
func (m *WorkbookWorksheetProtectionOptions) GetAllowFormatColumns()(*bool) {
    return m.allowFormatColumns
}
// GetAllowFormatRows gets the allowFormatRows property value. Represents the worksheet protection option of allowing formatting rows.
func (m *WorkbookWorksheetProtectionOptions) GetAllowFormatRows()(*bool) {
    return m.allowFormatRows
}
// GetAllowInsertColumns gets the allowInsertColumns property value. Represents the worksheet protection option of allowing inserting columns.
func (m *WorkbookWorksheetProtectionOptions) GetAllowInsertColumns()(*bool) {
    return m.allowInsertColumns
}
// GetAllowInsertHyperlinks gets the allowInsertHyperlinks property value. Represents the worksheet protection option of allowing inserting hyperlinks.
func (m *WorkbookWorksheetProtectionOptions) GetAllowInsertHyperlinks()(*bool) {
    return m.allowInsertHyperlinks
}
// GetAllowInsertRows gets the allowInsertRows property value. Represents the worksheet protection option of allowing inserting rows.
func (m *WorkbookWorksheetProtectionOptions) GetAllowInsertRows()(*bool) {
    return m.allowInsertRows
}
// GetAllowPivotTables gets the allowPivotTables property value. Represents the worksheet protection option of allowing using pivot table feature.
func (m *WorkbookWorksheetProtectionOptions) GetAllowPivotTables()(*bool) {
    return m.allowPivotTables
}
// GetAllowSort gets the allowSort property value. Represents the worksheet protection option of allowing using sort feature.
func (m *WorkbookWorksheetProtectionOptions) GetAllowSort()(*bool) {
    return m.allowSort
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookWorksheetProtectionOptions) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowAutoFilter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowAutoFilter)
    res["allowDeleteColumns"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowDeleteColumns)
    res["allowDeleteRows"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowDeleteRows)
    res["allowFormatCells"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowFormatCells)
    res["allowFormatColumns"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowFormatColumns)
    res["allowFormatRows"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowFormatRows)
    res["allowInsertColumns"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowInsertColumns)
    res["allowInsertHyperlinks"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowInsertHyperlinks)
    res["allowInsertRows"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowInsertRows)
    res["allowPivotTables"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowPivotTables)
    res["allowSort"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowSort)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *WorkbookWorksheetProtectionOptions) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *WorkbookWorksheetProtectionOptions) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowAutoFilter", m.GetAllowAutoFilter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowDeleteColumns", m.GetAllowDeleteColumns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowDeleteRows", m.GetAllowDeleteRows())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowFormatCells", m.GetAllowFormatCells())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowFormatColumns", m.GetAllowFormatColumns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowFormatRows", m.GetAllowFormatRows())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowInsertColumns", m.GetAllowInsertColumns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowInsertHyperlinks", m.GetAllowInsertHyperlinks())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowInsertRows", m.GetAllowInsertRows())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowPivotTables", m.GetAllowPivotTables())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowSort", m.GetAllowSort())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *WorkbookWorksheetProtectionOptions) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAllowAutoFilter sets the allowAutoFilter property value. Represents the worksheet protection option of allowing using auto filter feature.
func (m *WorkbookWorksheetProtectionOptions) SetAllowAutoFilter(value *bool)() {
    m.allowAutoFilter = value
}
// SetAllowDeleteColumns sets the allowDeleteColumns property value. Represents the worksheet protection option of allowing deleting columns.
func (m *WorkbookWorksheetProtectionOptions) SetAllowDeleteColumns(value *bool)() {
    m.allowDeleteColumns = value
}
// SetAllowDeleteRows sets the allowDeleteRows property value. Represents the worksheet protection option of allowing deleting rows.
func (m *WorkbookWorksheetProtectionOptions) SetAllowDeleteRows(value *bool)() {
    m.allowDeleteRows = value
}
// SetAllowFormatCells sets the allowFormatCells property value. Represents the worksheet protection option of allowing formatting cells.
func (m *WorkbookWorksheetProtectionOptions) SetAllowFormatCells(value *bool)() {
    m.allowFormatCells = value
}
// SetAllowFormatColumns sets the allowFormatColumns property value. Represents the worksheet protection option of allowing formatting columns.
func (m *WorkbookWorksheetProtectionOptions) SetAllowFormatColumns(value *bool)() {
    m.allowFormatColumns = value
}
// SetAllowFormatRows sets the allowFormatRows property value. Represents the worksheet protection option of allowing formatting rows.
func (m *WorkbookWorksheetProtectionOptions) SetAllowFormatRows(value *bool)() {
    m.allowFormatRows = value
}
// SetAllowInsertColumns sets the allowInsertColumns property value. Represents the worksheet protection option of allowing inserting columns.
func (m *WorkbookWorksheetProtectionOptions) SetAllowInsertColumns(value *bool)() {
    m.allowInsertColumns = value
}
// SetAllowInsertHyperlinks sets the allowInsertHyperlinks property value. Represents the worksheet protection option of allowing inserting hyperlinks.
func (m *WorkbookWorksheetProtectionOptions) SetAllowInsertHyperlinks(value *bool)() {
    m.allowInsertHyperlinks = value
}
// SetAllowInsertRows sets the allowInsertRows property value. Represents the worksheet protection option of allowing inserting rows.
func (m *WorkbookWorksheetProtectionOptions) SetAllowInsertRows(value *bool)() {
    m.allowInsertRows = value
}
// SetAllowPivotTables sets the allowPivotTables property value. Represents the worksheet protection option of allowing using pivot table feature.
func (m *WorkbookWorksheetProtectionOptions) SetAllowPivotTables(value *bool)() {
    m.allowPivotTables = value
}
// SetAllowSort sets the allowSort property value. Represents the worksheet protection option of allowing using sort feature.
func (m *WorkbookWorksheetProtectionOptions) SetAllowSort(value *bool)() {
    m.allowSort = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WorkbookWorksheetProtectionOptions) SetOdataType(value *string)() {
    m.odataType = value
}
