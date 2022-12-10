package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookChart provides operations to manage the collection of agreement entities.
type WorkbookChart struct {
    Entity
    // Represents chart axes. Read-only.
    axes WorkbookChartAxesable
    // Represents the datalabels on the chart. Read-only.
    dataLabels WorkbookChartDataLabelsable
    // Encapsulates the format properties for the chart area. Read-only.
    format WorkbookChartAreaFormatable
    // Represents the height, in points, of the chart object.
    height *float64
    // The distance, in points, from the left side of the chart to the worksheet origin.
    left *float64
    // Represents the legend for the chart. Read-only.
    legend WorkbookChartLegendable
    // Represents the name of a chart object.
    name *string
    // Represents either a single series or collection of series in the chart. Read-only.
    series []WorkbookChartSeriesable
    // Represents the title of the specified chart, including the text, visibility, position and formating of the title. Read-only.
    title WorkbookChartTitleable
    // Represents the distance, in points, from the top edge of the object to the top of row 1 (on a worksheet) or the top of the chart area (on a chart).
    top *float64
    // Represents the width, in points, of the chart object.
    width *float64
    // The worksheet containing the current chart. Read-only.
    worksheet WorkbookWorksheetable
}
// NewWorkbookChart instantiates a new workbookChart and sets the default values.
func NewWorkbookChart()(*WorkbookChart) {
    m := &WorkbookChart{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookChartFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChart(), nil
}
// GetAxes gets the axes property value. Represents chart axes. Read-only.
func (m *WorkbookChart) GetAxes()(WorkbookChartAxesable) {
    return m.axes
}
// GetDataLabels gets the dataLabels property value. Represents the datalabels on the chart. Read-only.
func (m *WorkbookChart) GetDataLabels()(WorkbookChartDataLabelsable) {
    return m.dataLabels
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookChart) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["axes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartAxesFromDiscriminatorValue , m.SetAxes)
    res["dataLabels"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartDataLabelsFromDiscriminatorValue , m.SetDataLabels)
    res["format"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartAreaFormatFromDiscriminatorValue , m.SetFormat)
    res["height"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetHeight)
    res["left"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetLeft)
    res["legend"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartLegendFromDiscriminatorValue , m.SetLegend)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["series"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookChartSeriesFromDiscriminatorValue , m.SetSeries)
    res["title"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartTitleFromDiscriminatorValue , m.SetTitle)
    res["top"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetTop)
    res["width"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetWidth)
    res["worksheet"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookWorksheetFromDiscriminatorValue , m.SetWorksheet)
    return res
}
// GetFormat gets the format property value. Encapsulates the format properties for the chart area. Read-only.
func (m *WorkbookChart) GetFormat()(WorkbookChartAreaFormatable) {
    return m.format
}
// GetHeight gets the height property value. Represents the height, in points, of the chart object.
func (m *WorkbookChart) GetHeight()(*float64) {
    return m.height
}
// GetLeft gets the left property value. The distance, in points, from the left side of the chart to the worksheet origin.
func (m *WorkbookChart) GetLeft()(*float64) {
    return m.left
}
// GetLegend gets the legend property value. Represents the legend for the chart. Read-only.
func (m *WorkbookChart) GetLegend()(WorkbookChartLegendable) {
    return m.legend
}
// GetName gets the name property value. Represents the name of a chart object.
func (m *WorkbookChart) GetName()(*string) {
    return m.name
}
// GetSeries gets the series property value. Represents either a single series or collection of series in the chart. Read-only.
func (m *WorkbookChart) GetSeries()([]WorkbookChartSeriesable) {
    return m.series
}
// GetTitle gets the title property value. Represents the title of the specified chart, including the text, visibility, position and formating of the title. Read-only.
func (m *WorkbookChart) GetTitle()(WorkbookChartTitleable) {
    return m.title
}
// GetTop gets the top property value. Represents the distance, in points, from the top edge of the object to the top of row 1 (on a worksheet) or the top of the chart area (on a chart).
func (m *WorkbookChart) GetTop()(*float64) {
    return m.top
}
// GetWidth gets the width property value. Represents the width, in points, of the chart object.
func (m *WorkbookChart) GetWidth()(*float64) {
    return m.width
}
// GetWorksheet gets the worksheet property value. The worksheet containing the current chart. Read-only.
func (m *WorkbookChart) GetWorksheet()(WorkbookWorksheetable) {
    return m.worksheet
}
// Serialize serializes information the current object
func (m *WorkbookChart) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("axes", m.GetAxes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("dataLabels", m.GetDataLabels())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("format", m.GetFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("height", m.GetHeight())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("left", m.GetLeft())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("legend", m.GetLegend())
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
    if m.GetSeries() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSeries())
        err = writer.WriteCollectionOfObjectValues("series", cast)
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
    {
        err = writer.WriteFloat64Value("top", m.GetTop())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("width", m.GetWidth())
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
// SetAxes sets the axes property value. Represents chart axes. Read-only.
func (m *WorkbookChart) SetAxes(value WorkbookChartAxesable)() {
    m.axes = value
}
// SetDataLabels sets the dataLabels property value. Represents the datalabels on the chart. Read-only.
func (m *WorkbookChart) SetDataLabels(value WorkbookChartDataLabelsable)() {
    m.dataLabels = value
}
// SetFormat sets the format property value. Encapsulates the format properties for the chart area. Read-only.
func (m *WorkbookChart) SetFormat(value WorkbookChartAreaFormatable)() {
    m.format = value
}
// SetHeight sets the height property value. Represents the height, in points, of the chart object.
func (m *WorkbookChart) SetHeight(value *float64)() {
    m.height = value
}
// SetLeft sets the left property value. The distance, in points, from the left side of the chart to the worksheet origin.
func (m *WorkbookChart) SetLeft(value *float64)() {
    m.left = value
}
// SetLegend sets the legend property value. Represents the legend for the chart. Read-only.
func (m *WorkbookChart) SetLegend(value WorkbookChartLegendable)() {
    m.legend = value
}
// SetName sets the name property value. Represents the name of a chart object.
func (m *WorkbookChart) SetName(value *string)() {
    m.name = value
}
// SetSeries sets the series property value. Represents either a single series or collection of series in the chart. Read-only.
func (m *WorkbookChart) SetSeries(value []WorkbookChartSeriesable)() {
    m.series = value
}
// SetTitle sets the title property value. Represents the title of the specified chart, including the text, visibility, position and formating of the title. Read-only.
func (m *WorkbookChart) SetTitle(value WorkbookChartTitleable)() {
    m.title = value
}
// SetTop sets the top property value. Represents the distance, in points, from the top edge of the object to the top of row 1 (on a worksheet) or the top of the chart area (on a chart).
func (m *WorkbookChart) SetTop(value *float64)() {
    m.top = value
}
// SetWidth sets the width property value. Represents the width, in points, of the chart object.
func (m *WorkbookChart) SetWidth(value *float64)() {
    m.width = value
}
// SetWorksheet sets the worksheet property value. The worksheet containing the current chart. Read-only.
func (m *WorkbookChart) SetWorksheet(value WorkbookWorksheetable)() {
    m.worksheet = value
}
