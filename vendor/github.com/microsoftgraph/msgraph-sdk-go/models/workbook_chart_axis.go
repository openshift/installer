package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookChartAxis 
type WorkbookChartAxis struct {
    Entity
    // Represents the formatting of a chart object, which includes line and font formatting. Read-only.
    format WorkbookChartAxisFormatable
    // Returns a gridlines object that represents the major gridlines for the specified axis. Read-only.
    majorGridlines WorkbookChartGridlinesable
    // Represents the interval between two major tick marks. Can be set to a numeric value or an empty string.  The returned value is always a number.
    majorUnit Jsonable
    // Represents the maximum value on the value axis.  Can be set to a numeric value or an empty string (for automatic axis values).  The returned value is always a number.
    maximum Jsonable
    // Represents the minimum value on the value axis. Can be set to a numeric value or an empty string (for automatic axis values).  The returned value is always a number.
    minimum Jsonable
    // Returns a Gridlines object that represents the minor gridlines for the specified axis. Read-only.
    minorGridlines WorkbookChartGridlinesable
    // Represents the interval between two minor tick marks. 'Can be set to a numeric value or an empty string (for automatic axis values). The returned value is always a number.
    minorUnit Jsonable
    // Represents the axis title. Read-only.
    title WorkbookChartAxisTitleable
}
// NewWorkbookChartAxis instantiates a new workbookChartAxis and sets the default values.
func NewWorkbookChartAxis()(*WorkbookChartAxis) {
    m := &WorkbookChartAxis{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartAxisFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookChartAxisFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartAxis(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookChartAxis) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["format"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartAxisFormatFromDiscriminatorValue , m.SetFormat)
    res["majorGridlines"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartGridlinesFromDiscriminatorValue , m.SetMajorGridlines)
    res["majorUnit"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateJsonFromDiscriminatorValue , m.SetMajorUnit)
    res["maximum"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateJsonFromDiscriminatorValue , m.SetMaximum)
    res["minimum"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateJsonFromDiscriminatorValue , m.SetMinimum)
    res["minorGridlines"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartGridlinesFromDiscriminatorValue , m.SetMinorGridlines)
    res["minorUnit"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateJsonFromDiscriminatorValue , m.SetMinorUnit)
    res["title"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartAxisTitleFromDiscriminatorValue , m.SetTitle)
    return res
}
// GetFormat gets the format property value. Represents the formatting of a chart object, which includes line and font formatting. Read-only.
func (m *WorkbookChartAxis) GetFormat()(WorkbookChartAxisFormatable) {
    return m.format
}
// GetMajorGridlines gets the majorGridlines property value. Returns a gridlines object that represents the major gridlines for the specified axis. Read-only.
func (m *WorkbookChartAxis) GetMajorGridlines()(WorkbookChartGridlinesable) {
    return m.majorGridlines
}
// GetMajorUnit gets the majorUnit property value. Represents the interval between two major tick marks. Can be set to a numeric value or an empty string.  The returned value is always a number.
func (m *WorkbookChartAxis) GetMajorUnit()(Jsonable) {
    return m.majorUnit
}
// GetMaximum gets the maximum property value. Represents the maximum value on the value axis.  Can be set to a numeric value or an empty string (for automatic axis values).  The returned value is always a number.
func (m *WorkbookChartAxis) GetMaximum()(Jsonable) {
    return m.maximum
}
// GetMinimum gets the minimum property value. Represents the minimum value on the value axis. Can be set to a numeric value or an empty string (for automatic axis values).  The returned value is always a number.
func (m *WorkbookChartAxis) GetMinimum()(Jsonable) {
    return m.minimum
}
// GetMinorGridlines gets the minorGridlines property value. Returns a Gridlines object that represents the minor gridlines for the specified axis. Read-only.
func (m *WorkbookChartAxis) GetMinorGridlines()(WorkbookChartGridlinesable) {
    return m.minorGridlines
}
// GetMinorUnit gets the minorUnit property value. Represents the interval between two minor tick marks. 'Can be set to a numeric value or an empty string (for automatic axis values). The returned value is always a number.
func (m *WorkbookChartAxis) GetMinorUnit()(Jsonable) {
    return m.minorUnit
}
// GetTitle gets the title property value. Represents the axis title. Read-only.
func (m *WorkbookChartAxis) GetTitle()(WorkbookChartAxisTitleable) {
    return m.title
}
// Serialize serializes information the current object
func (m *WorkbookChartAxis) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("format", m.GetFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("majorGridlines", m.GetMajorGridlines())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("majorUnit", m.GetMajorUnit())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("maximum", m.GetMaximum())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("minimum", m.GetMinimum())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("minorGridlines", m.GetMinorGridlines())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("minorUnit", m.GetMinorUnit())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormat sets the format property value. Represents the formatting of a chart object, which includes line and font formatting. Read-only.
func (m *WorkbookChartAxis) SetFormat(value WorkbookChartAxisFormatable)() {
    m.format = value
}
// SetMajorGridlines sets the majorGridlines property value. Returns a gridlines object that represents the major gridlines for the specified axis. Read-only.
func (m *WorkbookChartAxis) SetMajorGridlines(value WorkbookChartGridlinesable)() {
    m.majorGridlines = value
}
// SetMajorUnit sets the majorUnit property value. Represents the interval between two major tick marks. Can be set to a numeric value or an empty string.  The returned value is always a number.
func (m *WorkbookChartAxis) SetMajorUnit(value Jsonable)() {
    m.majorUnit = value
}
// SetMaximum sets the maximum property value. Represents the maximum value on the value axis.  Can be set to a numeric value or an empty string (for automatic axis values).  The returned value is always a number.
func (m *WorkbookChartAxis) SetMaximum(value Jsonable)() {
    m.maximum = value
}
// SetMinimum sets the minimum property value. Represents the minimum value on the value axis. Can be set to a numeric value or an empty string (for automatic axis values).  The returned value is always a number.
func (m *WorkbookChartAxis) SetMinimum(value Jsonable)() {
    m.minimum = value
}
// SetMinorGridlines sets the minorGridlines property value. Returns a Gridlines object that represents the minor gridlines for the specified axis. Read-only.
func (m *WorkbookChartAxis) SetMinorGridlines(value WorkbookChartGridlinesable)() {
    m.minorGridlines = value
}
// SetMinorUnit sets the minorUnit property value. Represents the interval between two minor tick marks. 'Can be set to a numeric value or an empty string (for automatic axis values). The returned value is always a number.
func (m *WorkbookChartAxis) SetMinorUnit(value Jsonable)() {
    m.minorUnit = value
}
// SetTitle sets the title property value. Represents the axis title. Read-only.
func (m *WorkbookChartAxis) SetTitle(value WorkbookChartAxisTitleable)() {
    m.title = value
}
