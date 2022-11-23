package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookChartPointFormat 
type WorkbookChartPointFormat struct {
    Entity
    // Represents the fill format of a chart, which includes background formating information. Read-only.
    fill WorkbookChartFillable
}
// NewWorkbookChartPointFormat instantiates a new workbookChartPointFormat and sets the default values.
func NewWorkbookChartPointFormat()(*WorkbookChartPointFormat) {
    m := &WorkbookChartPointFormat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartPointFormatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookChartPointFormatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartPointFormat(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookChartPointFormat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["fill"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartFillFromDiscriminatorValue , m.SetFill)
    return res
}
// GetFill gets the fill property value. Represents the fill format of a chart, which includes background formating information. Read-only.
func (m *WorkbookChartPointFormat) GetFill()(WorkbookChartFillable) {
    return m.fill
}
// Serialize serializes information the current object
func (m *WorkbookChartPointFormat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("fill", m.GetFill())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFill sets the fill property value. Represents the fill format of a chart, which includes background formating information. Read-only.
func (m *WorkbookChartPointFormat) SetFill(value WorkbookChartFillable)() {
    m.fill = value
}
