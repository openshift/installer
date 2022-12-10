package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookChartAreaFormat 
type WorkbookChartAreaFormat struct {
    Entity
    // Represents the fill format of an object, which includes background formatting information. Read-only.
    fill WorkbookChartFillable
    // Represents the font attributes (font name, font size, color, etc.) for the current object. Read-only.
    font WorkbookChartFontable
}
// NewWorkbookChartAreaFormat instantiates a new workbookChartAreaFormat and sets the default values.
func NewWorkbookChartAreaFormat()(*WorkbookChartAreaFormat) {
    m := &WorkbookChartAreaFormat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartAreaFormatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookChartAreaFormatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartAreaFormat(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookChartAreaFormat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["fill"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartFillFromDiscriminatorValue , m.SetFill)
    res["font"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookChartFontFromDiscriminatorValue , m.SetFont)
    return res
}
// GetFill gets the fill property value. Represents the fill format of an object, which includes background formatting information. Read-only.
func (m *WorkbookChartAreaFormat) GetFill()(WorkbookChartFillable) {
    return m.fill
}
// GetFont gets the font property value. Represents the font attributes (font name, font size, color, etc.) for the current object. Read-only.
func (m *WorkbookChartAreaFormat) GetFont()(WorkbookChartFontable) {
    return m.font
}
// Serialize serializes information the current object
func (m *WorkbookChartAreaFormat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    {
        err = writer.WriteObjectValue("font", m.GetFont())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFill sets the fill property value. Represents the fill format of an object, which includes background formatting information. Read-only.
func (m *WorkbookChartAreaFormat) SetFill(value WorkbookChartFillable)() {
    m.fill = value
}
// SetFont sets the font property value. Represents the font attributes (font name, font size, color, etc.) for the current object. Read-only.
func (m *WorkbookChartAreaFormat) SetFont(value WorkbookChartFontable)() {
    m.font = value
}
