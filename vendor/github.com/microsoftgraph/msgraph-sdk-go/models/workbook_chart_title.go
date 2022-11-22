package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookChartTitle 
type WorkbookChartTitle struct {
    Entity
    // Represents the formatting of a chart title, which includes fill and font formatting. Read-only.
    format WorkbookChartTitleFormatable
    // Boolean value representing if the chart title will overlay the chart or not.
    overlay *bool
    // Represents the title text of a chart.
    text *string
    // A boolean value the represents the visibility of a chart title object.
    visible *bool
}
// NewWorkbookChartTitle instantiates a new workbookChartTitle and sets the default values.
func NewWorkbookChartTitle()(*WorkbookChartTitle) {
    m := &WorkbookChartTitle{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartTitleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookChartTitleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartTitle(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookChartTitle) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["format"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartTitleFormatFromDiscriminatorValue , m.SetFormat)
    res["overlay"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetOverlay)
    res["text"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetText)
    res["visible"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetVisible)
    return res
}
// GetFormat gets the format property value. Represents the formatting of a chart title, which includes fill and font formatting. Read-only.
func (m *WorkbookChartTitle) GetFormat()(WorkbookChartTitleFormatable) {
    return m.format
}
// GetOverlay gets the overlay property value. Boolean value representing if the chart title will overlay the chart or not.
func (m *WorkbookChartTitle) GetOverlay()(*bool) {
    return m.overlay
}
// GetText gets the text property value. Represents the title text of a chart.
func (m *WorkbookChartTitle) GetText()(*string) {
    return m.text
}
// GetVisible gets the visible property value. A boolean value the represents the visibility of a chart title object.
func (m *WorkbookChartTitle) GetVisible()(*bool) {
    return m.visible
}
// Serialize serializes information the current object
func (m *WorkbookChartTitle) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteBoolValue("overlay", m.GetOverlay())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("text", m.GetText())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("visible", m.GetVisible())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormat sets the format property value. Represents the formatting of a chart title, which includes fill and font formatting. Read-only.
func (m *WorkbookChartTitle) SetFormat(value WorkbookChartTitleFormatable)() {
    m.format = value
}
// SetOverlay sets the overlay property value. Boolean value representing if the chart title will overlay the chart or not.
func (m *WorkbookChartTitle) SetOverlay(value *bool)() {
    m.overlay = value
}
// SetText sets the text property value. Represents the title text of a chart.
func (m *WorkbookChartTitle) SetText(value *string)() {
    m.text = value
}
// SetVisible sets the visible property value. A boolean value the represents the visibility of a chart title object.
func (m *WorkbookChartTitle) SetVisible(value *bool)() {
    m.visible = value
}
