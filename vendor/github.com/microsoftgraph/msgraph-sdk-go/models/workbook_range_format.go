package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookRangeFormat 
type WorkbookRangeFormat struct {
    Entity
    // Collection of border objects that apply to the overall range selected Read-only.
    borders []WorkbookRangeBorderable
    // Gets or sets the width of all colums within the range. If the column widths are not uniform, null will be returned.
    columnWidth *float64
    // Returns the fill object defined on the overall range. Read-only.
    fill WorkbookRangeFillable
    // Returns the font object defined on the overall range selected Read-only.
    font WorkbookRangeFontable
    // Represents the horizontal alignment for the specified object. The possible values are: General, Left, Center, Right, Fill, Justify, CenterAcrossSelection, Distributed.
    horizontalAlignment *string
    // Returns the format protection object for a range. Read-only.
    protection WorkbookFormatProtectionable
    // Gets or sets the height of all rows in the range. If the row heights are not uniform null will be returned.
    rowHeight *float64
    // Represents the vertical alignment for the specified object. The possible values are: Top, Center, Bottom, Justify, Distributed.
    verticalAlignment *string
    // Indicates if Excel wraps the text in the object. A null value indicates that the entire range doesn't have uniform wrap setting
    wrapText *bool
}
// NewWorkbookRangeFormat instantiates a new workbookRangeFormat and sets the default values.
func NewWorkbookRangeFormat()(*WorkbookRangeFormat) {
    m := &WorkbookRangeFormat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookRangeFormatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookRangeFormatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookRangeFormat(), nil
}
// GetBorders gets the borders property value. Collection of border objects that apply to the overall range selected Read-only.
func (m *WorkbookRangeFormat) GetBorders()([]WorkbookRangeBorderable) {
    return m.borders
}
// GetColumnWidth gets the columnWidth property value. Gets or sets the width of all colums within the range. If the column widths are not uniform, null will be returned.
func (m *WorkbookRangeFormat) GetColumnWidth()(*float64) {
    return m.columnWidth
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookRangeFormat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["borders"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookRangeBorderFromDiscriminatorValue , m.SetBorders)
    res["columnWidth"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetColumnWidth)
    res["fill"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookRangeFillFromDiscriminatorValue , m.SetFill)
    res["font"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookRangeFontFromDiscriminatorValue , m.SetFont)
    res["horizontalAlignment"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetHorizontalAlignment)
    res["protection"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookFormatProtectionFromDiscriminatorValue , m.SetProtection)
    res["rowHeight"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetRowHeight)
    res["verticalAlignment"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVerticalAlignment)
    res["wrapText"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetWrapText)
    return res
}
// GetFill gets the fill property value. Returns the fill object defined on the overall range. Read-only.
func (m *WorkbookRangeFormat) GetFill()(WorkbookRangeFillable) {
    return m.fill
}
// GetFont gets the font property value. Returns the font object defined on the overall range selected Read-only.
func (m *WorkbookRangeFormat) GetFont()(WorkbookRangeFontable) {
    return m.font
}
// GetHorizontalAlignment gets the horizontalAlignment property value. Represents the horizontal alignment for the specified object. The possible values are: General, Left, Center, Right, Fill, Justify, CenterAcrossSelection, Distributed.
func (m *WorkbookRangeFormat) GetHorizontalAlignment()(*string) {
    return m.horizontalAlignment
}
// GetProtection gets the protection property value. Returns the format protection object for a range. Read-only.
func (m *WorkbookRangeFormat) GetProtection()(WorkbookFormatProtectionable) {
    return m.protection
}
// GetRowHeight gets the rowHeight property value. Gets or sets the height of all rows in the range. If the row heights are not uniform null will be returned.
func (m *WorkbookRangeFormat) GetRowHeight()(*float64) {
    return m.rowHeight
}
// GetVerticalAlignment gets the verticalAlignment property value. Represents the vertical alignment for the specified object. The possible values are: Top, Center, Bottom, Justify, Distributed.
func (m *WorkbookRangeFormat) GetVerticalAlignment()(*string) {
    return m.verticalAlignment
}
// GetWrapText gets the wrapText property value. Indicates if Excel wraps the text in the object. A null value indicates that the entire range doesn't have uniform wrap setting
func (m *WorkbookRangeFormat) GetWrapText()(*bool) {
    return m.wrapText
}
// Serialize serializes information the current object
func (m *WorkbookRangeFormat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetBorders() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetBorders())
        err = writer.WriteCollectionOfObjectValues("borders", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("columnWidth", m.GetColumnWidth())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("fill", m.GetFill())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("font", m.GetFont())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("horizontalAlignment", m.GetHorizontalAlignment())
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
    {
        err = writer.WriteFloat64Value("rowHeight", m.GetRowHeight())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("verticalAlignment", m.GetVerticalAlignment())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wrapText", m.GetWrapText())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBorders sets the borders property value. Collection of border objects that apply to the overall range selected Read-only.
func (m *WorkbookRangeFormat) SetBorders(value []WorkbookRangeBorderable)() {
    m.borders = value
}
// SetColumnWidth sets the columnWidth property value. Gets or sets the width of all colums within the range. If the column widths are not uniform, null will be returned.
func (m *WorkbookRangeFormat) SetColumnWidth(value *float64)() {
    m.columnWidth = value
}
// SetFill sets the fill property value. Returns the fill object defined on the overall range. Read-only.
func (m *WorkbookRangeFormat) SetFill(value WorkbookRangeFillable)() {
    m.fill = value
}
// SetFont sets the font property value. Returns the font object defined on the overall range selected Read-only.
func (m *WorkbookRangeFormat) SetFont(value WorkbookRangeFontable)() {
    m.font = value
}
// SetHorizontalAlignment sets the horizontalAlignment property value. Represents the horizontal alignment for the specified object. The possible values are: General, Left, Center, Right, Fill, Justify, CenterAcrossSelection, Distributed.
func (m *WorkbookRangeFormat) SetHorizontalAlignment(value *string)() {
    m.horizontalAlignment = value
}
// SetProtection sets the protection property value. Returns the format protection object for a range. Read-only.
func (m *WorkbookRangeFormat) SetProtection(value WorkbookFormatProtectionable)() {
    m.protection = value
}
// SetRowHeight sets the rowHeight property value. Gets or sets the height of all rows in the range. If the row heights are not uniform null will be returned.
func (m *WorkbookRangeFormat) SetRowHeight(value *float64)() {
    m.rowHeight = value
}
// SetVerticalAlignment sets the verticalAlignment property value. Represents the vertical alignment for the specified object. The possible values are: Top, Center, Bottom, Justify, Distributed.
func (m *WorkbookRangeFormat) SetVerticalAlignment(value *string)() {
    m.verticalAlignment = value
}
// SetWrapText sets the wrapText property value. Indicates if Excel wraps the text in the object. A null value indicates that the entire range doesn't have uniform wrap setting
func (m *WorkbookRangeFormat) SetWrapText(value *bool)() {
    m.wrapText = value
}
