package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookChartSeries provides operations to manage the collection of agreement entities.
type WorkbookChartSeries struct {
    Entity
    // Represents the formatting of a chart series, which includes fill and line formatting. Read-only.
    format WorkbookChartSeriesFormatable
    // Represents the name of a series in a chart.
    name *string
    // Represents a collection of all points in the series. Read-only.
    points []WorkbookChartPointable
}
// NewWorkbookChartSeries instantiates a new workbookChartSeries and sets the default values.
func NewWorkbookChartSeries()(*WorkbookChartSeries) {
    m := &WorkbookChartSeries{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartSeriesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookChartSeriesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartSeries(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookChartSeries) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["format"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartSeriesFormatFromDiscriminatorValue , m.SetFormat)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["points"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookChartPointFromDiscriminatorValue , m.SetPoints)
    return res
}
// GetFormat gets the format property value. Represents the formatting of a chart series, which includes fill and line formatting. Read-only.
func (m *WorkbookChartSeries) GetFormat()(WorkbookChartSeriesFormatable) {
    return m.format
}
// GetName gets the name property value. Represents the name of a series in a chart.
func (m *WorkbookChartSeries) GetName()(*string) {
    return m.name
}
// GetPoints gets the points property value. Represents a collection of all points in the series. Read-only.
func (m *WorkbookChartSeries) GetPoints()([]WorkbookChartPointable) {
    return m.points
}
// Serialize serializes information the current object
func (m *WorkbookChartSeries) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetPoints() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPoints())
        err = writer.WriteCollectionOfObjectValues("points", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormat sets the format property value. Represents the formatting of a chart series, which includes fill and line formatting. Read-only.
func (m *WorkbookChartSeries) SetFormat(value WorkbookChartSeriesFormatable)() {
    m.format = value
}
// SetName sets the name property value. Represents the name of a series in a chart.
func (m *WorkbookChartSeries) SetName(value *string)() {
    m.name = value
}
// SetPoints sets the points property value. Represents a collection of all points in the series. Read-only.
func (m *WorkbookChartSeries) SetPoints(value []WorkbookChartPointable)() {
    m.points = value
}
