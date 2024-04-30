package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookRangeView 
type WorkbookRangeView struct {
    Entity
}
// NewWorkbookRangeView instantiates a new workbookRangeView and sets the default values.
func NewWorkbookRangeView()(*WorkbookRangeView) {
    m := &WorkbookRangeView{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookRangeViewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookRangeViewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookRangeView(), nil
}
// GetCellAddresses gets the cellAddresses property value. Represents the cell addresses
func (m *WorkbookRangeView) GetCellAddresses()(Jsonable) {
    val, err := m.GetBackingStore().Get("cellAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetColumnCount gets the columnCount property value. Returns the number of visible columns. Read-only.
func (m *WorkbookRangeView) GetColumnCount()(*int32) {
    val, err := m.GetBackingStore().Get("columnCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookRangeView) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["cellAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellAddresses(val.(Jsonable))
        }
        return nil
    }
    res["columnCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetColumnCount(val)
        }
        return nil
    }
    res["formulas"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormulas(val.(Jsonable))
        }
        return nil
    }
    res["formulasLocal"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormulasLocal(val.(Jsonable))
        }
        return nil
    }
    res["formulasR1C1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormulasR1C1(val.(Jsonable))
        }
        return nil
    }
    res["index"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIndex(val)
        }
        return nil
    }
    res["numberFormat"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNumberFormat(val.(Jsonable))
        }
        return nil
    }
    res["rowCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRowCount(val)
        }
        return nil
    }
    res["rows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookRangeViewFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookRangeViewable, len(val))
            for i, v := range val {
                res[i] = v.(WorkbookRangeViewable)
            }
            m.SetRows(res)
        }
        return nil
    }
    res["text"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetText(val.(Jsonable))
        }
        return nil
    }
    res["values"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValues(val.(Jsonable))
        }
        return nil
    }
    res["valueTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValueTypes(val.(Jsonable))
        }
        return nil
    }
    return res
}
// GetFormulas gets the formulas property value. Represents the formula in A1-style notation.
func (m *WorkbookRangeView) GetFormulas()(Jsonable) {
    val, err := m.GetBackingStore().Get("formulas")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetFormulasLocal gets the formulasLocal property value. Represents the formula in A1-style notation, in the user's language and number-formatting locale. For example, the English '=SUM(A1, 1.5)' formula would become '=SUMME(A1; 1,5)' in German.
func (m *WorkbookRangeView) GetFormulasLocal()(Jsonable) {
    val, err := m.GetBackingStore().Get("formulasLocal")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetFormulasR1C1 gets the formulasR1C1 property value. Represents the formula in R1C1-style notation.
func (m *WorkbookRangeView) GetFormulasR1C1()(Jsonable) {
    val, err := m.GetBackingStore().Get("formulasR1C1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetIndex gets the index property value. Index of the range.
func (m *WorkbookRangeView) GetIndex()(*int32) {
    val, err := m.GetBackingStore().Get("index")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNumberFormat gets the numberFormat property value. Represents Excel's number format code for the given cell. Read-only.
func (m *WorkbookRangeView) GetNumberFormat()(Jsonable) {
    val, err := m.GetBackingStore().Get("numberFormat")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetRowCount gets the rowCount property value. Returns the number of visible rows. Read-only.
func (m *WorkbookRangeView) GetRowCount()(*int32) {
    val, err := m.GetBackingStore().Get("rowCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRows gets the rows property value. Represents a collection of range views associated with the range. Read-only. Read-only.
func (m *WorkbookRangeView) GetRows()([]WorkbookRangeViewable) {
    val, err := m.GetBackingStore().Get("rows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookRangeViewable)
    }
    return nil
}
// GetText gets the text property value. Text values of the specified range. The Text value will not depend on the cell width. The # sign substitution that happens in Excel UI will not affect the text value returned by the API. Read-only.
func (m *WorkbookRangeView) GetText()(Jsonable) {
    val, err := m.GetBackingStore().Get("text")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetValues gets the values property value. Represents the raw values of the specified range view. The data returned could be of type string, number, or a boolean. Cell that contain an error will return the error string.
func (m *WorkbookRangeView) GetValues()(Jsonable) {
    val, err := m.GetBackingStore().Get("values")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetValueTypes gets the valueTypes property value. Represents the type of data of each cell. Read-only. The possible values are: Unknown, Empty, String, Integer, Double, Boolean, Error.
func (m *WorkbookRangeView) GetValueTypes()(Jsonable) {
    val, err := m.GetBackingStore().Get("valueTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookRangeView) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("cellAddresses", m.GetCellAddresses())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("columnCount", m.GetColumnCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("formulas", m.GetFormulas())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("formulasLocal", m.GetFormulasLocal())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("formulasR1C1", m.GetFormulasR1C1())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("index", m.GetIndex())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("numberFormat", m.GetNumberFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("rowCount", m.GetRowCount())
        if err != nil {
            return err
        }
    }
    if m.GetRows() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRows()))
        for i, v := range m.GetRows() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("rows", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("text", m.GetText())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("values", m.GetValues())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("valueTypes", m.GetValueTypes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCellAddresses sets the cellAddresses property value. Represents the cell addresses
func (m *WorkbookRangeView) SetCellAddresses(value Jsonable)() {
    err := m.GetBackingStore().Set("cellAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetColumnCount sets the columnCount property value. Returns the number of visible columns. Read-only.
func (m *WorkbookRangeView) SetColumnCount(value *int32)() {
    err := m.GetBackingStore().Set("columnCount", value)
    if err != nil {
        panic(err)
    }
}
// SetFormulas sets the formulas property value. Represents the formula in A1-style notation.
func (m *WorkbookRangeView) SetFormulas(value Jsonable)() {
    err := m.GetBackingStore().Set("formulas", value)
    if err != nil {
        panic(err)
    }
}
// SetFormulasLocal sets the formulasLocal property value. Represents the formula in A1-style notation, in the user's language and number-formatting locale. For example, the English '=SUM(A1, 1.5)' formula would become '=SUMME(A1; 1,5)' in German.
func (m *WorkbookRangeView) SetFormulasLocal(value Jsonable)() {
    err := m.GetBackingStore().Set("formulasLocal", value)
    if err != nil {
        panic(err)
    }
}
// SetFormulasR1C1 sets the formulasR1C1 property value. Represents the formula in R1C1-style notation.
func (m *WorkbookRangeView) SetFormulasR1C1(value Jsonable)() {
    err := m.GetBackingStore().Set("formulasR1C1", value)
    if err != nil {
        panic(err)
    }
}
// SetIndex sets the index property value. Index of the range.
func (m *WorkbookRangeView) SetIndex(value *int32)() {
    err := m.GetBackingStore().Set("index", value)
    if err != nil {
        panic(err)
    }
}
// SetNumberFormat sets the numberFormat property value. Represents Excel's number format code for the given cell. Read-only.
func (m *WorkbookRangeView) SetNumberFormat(value Jsonable)() {
    err := m.GetBackingStore().Set("numberFormat", value)
    if err != nil {
        panic(err)
    }
}
// SetRowCount sets the rowCount property value. Returns the number of visible rows. Read-only.
func (m *WorkbookRangeView) SetRowCount(value *int32)() {
    err := m.GetBackingStore().Set("rowCount", value)
    if err != nil {
        panic(err)
    }
}
// SetRows sets the rows property value. Represents a collection of range views associated with the range. Read-only. Read-only.
func (m *WorkbookRangeView) SetRows(value []WorkbookRangeViewable)() {
    err := m.GetBackingStore().Set("rows", value)
    if err != nil {
        panic(err)
    }
}
// SetText sets the text property value. Text values of the specified range. The Text value will not depend on the cell width. The # sign substitution that happens in Excel UI will not affect the text value returned by the API. Read-only.
func (m *WorkbookRangeView) SetText(value Jsonable)() {
    err := m.GetBackingStore().Set("text", value)
    if err != nil {
        panic(err)
    }
}
// SetValues sets the values property value. Represents the raw values of the specified range view. The data returned could be of type string, number, or a boolean. Cell that contain an error will return the error string.
func (m *WorkbookRangeView) SetValues(value Jsonable)() {
    err := m.GetBackingStore().Set("values", value)
    if err != nil {
        panic(err)
    }
}
// SetValueTypes sets the valueTypes property value. Represents the type of data of each cell. Read-only. The possible values are: Unknown, Empty, String, Integer, Double, Boolean, Error.
func (m *WorkbookRangeView) SetValueTypes(value Jsonable)() {
    err := m.GetBackingStore().Set("valueTypes", value)
    if err != nil {
        panic(err)
    }
}
// WorkbookRangeViewable 
type WorkbookRangeViewable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCellAddresses()(Jsonable)
    GetColumnCount()(*int32)
    GetFormulas()(Jsonable)
    GetFormulasLocal()(Jsonable)
    GetFormulasR1C1()(Jsonable)
    GetIndex()(*int32)
    GetNumberFormat()(Jsonable)
    GetRowCount()(*int32)
    GetRows()([]WorkbookRangeViewable)
    GetText()(Jsonable)
    GetValues()(Jsonable)
    GetValueTypes()(Jsonable)
    SetCellAddresses(value Jsonable)()
    SetColumnCount(value *int32)()
    SetFormulas(value Jsonable)()
    SetFormulasLocal(value Jsonable)()
    SetFormulasR1C1(value Jsonable)()
    SetIndex(value *int32)()
    SetNumberFormat(value Jsonable)()
    SetRowCount(value *int32)()
    SetRows(value []WorkbookRangeViewable)()
    SetText(value Jsonable)()
    SetValues(value Jsonable)()
    SetValueTypes(value Jsonable)()
}
