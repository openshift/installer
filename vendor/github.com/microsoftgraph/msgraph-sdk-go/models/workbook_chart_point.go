package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookChartPoint provides operations to manage the collection of agreement entities.
type WorkbookChartPoint struct {
    Entity
    // Encapsulates the format properties chart point. Read-only.
    format WorkbookChartPointFormatable
    // Returns the value of a chart point. Read-only.
    value Jsonable
}
// NewWorkbookChartPoint instantiates a new workbookChartPoint and sets the default values.
func NewWorkbookChartPoint()(*WorkbookChartPoint) {
    m := &WorkbookChartPoint{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartPointFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookChartPointFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartPoint(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookChartPoint) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["format"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartPointFormatFromDiscriminatorValue , m.SetFormat)
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateJsonFromDiscriminatorValue , m.SetValue)
    return res
}
// GetFormat gets the format property value. Encapsulates the format properties chart point. Read-only.
func (m *WorkbookChartPoint) GetFormat()(WorkbookChartPointFormatable) {
    return m.format
}
// GetValue gets the value property value. Returns the value of a chart point. Read-only.
func (m *WorkbookChartPoint) GetValue()(Jsonable) {
    return m.value
}
// Serialize serializes information the current object
func (m *WorkbookChartPoint) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteObjectValue("value", m.GetValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormat sets the format property value. Encapsulates the format properties chart point. Read-only.
func (m *WorkbookChartPoint) SetFormat(value WorkbookChartPointFormatable)() {
    m.format = value
}
// SetValue sets the value property value. Returns the value of a chart point. Read-only.
func (m *WorkbookChartPoint) SetValue(value Jsonable)() {
    m.value = value
}
