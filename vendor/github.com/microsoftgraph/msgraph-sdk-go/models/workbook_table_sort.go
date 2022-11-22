package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookTableSort 
type WorkbookTableSort struct {
    Entity
    // Represents the current conditions used to last sort the table. Read-only.
    fields []WorkbookSortFieldable
    // Represents whether the casing impacted the last sort of the table. Read-only.
    matchCase *bool
    // Represents Chinese character ordering method last used to sort the table. The possible values are: PinYin, StrokeCount. Read-only.
    method *string
}
// NewWorkbookTableSort instantiates a new workbookTableSort and sets the default values.
func NewWorkbookTableSort()(*WorkbookTableSort) {
    m := &WorkbookTableSort{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookTableSortFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookTableSortFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookTableSort(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookTableSort) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["fields"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookSortFieldFromDiscriminatorValue , m.SetFields)
    res["matchCase"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetMatchCase)
    res["method"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMethod)
    return res
}
// GetFields gets the fields property value. Represents the current conditions used to last sort the table. Read-only.
func (m *WorkbookTableSort) GetFields()([]WorkbookSortFieldable) {
    return m.fields
}
// GetMatchCase gets the matchCase property value. Represents whether the casing impacted the last sort of the table. Read-only.
func (m *WorkbookTableSort) GetMatchCase()(*bool) {
    return m.matchCase
}
// GetMethod gets the method property value. Represents Chinese character ordering method last used to sort the table. The possible values are: PinYin, StrokeCount. Read-only.
func (m *WorkbookTableSort) GetMethod()(*string) {
    return m.method
}
// Serialize serializes information the current object
func (m *WorkbookTableSort) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetFields() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetFields())
        err = writer.WriteCollectionOfObjectValues("fields", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("matchCase", m.GetMatchCase())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("method", m.GetMethod())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFields sets the fields property value. Represents the current conditions used to last sort the table. Read-only.
func (m *WorkbookTableSort) SetFields(value []WorkbookSortFieldable)() {
    m.fields = value
}
// SetMatchCase sets the matchCase property value. Represents whether the casing impacted the last sort of the table. Read-only.
func (m *WorkbookTableSort) SetMatchCase(value *bool)() {
    m.matchCase = value
}
// SetMethod sets the method property value. Represents Chinese character ordering method last used to sort the table. The possible values are: PinYin, StrokeCount. Read-only.
func (m *WorkbookTableSort) SetMethod(value *string)() {
    m.method = value
}
