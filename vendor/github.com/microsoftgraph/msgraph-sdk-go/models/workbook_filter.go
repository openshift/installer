package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookFilter 
type WorkbookFilter struct {
    Entity
    // The currently applied filter on the given column. Read-only.
    criteria WorkbookFilterCriteriaable
}
// NewWorkbookFilter instantiates a new workbookFilter and sets the default values.
func NewWorkbookFilter()(*WorkbookFilter) {
    m := &WorkbookFilter{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookFilterFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookFilterFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookFilter(), nil
}
// GetCriteria gets the criteria property value. The currently applied filter on the given column. Read-only.
func (m *WorkbookFilter) GetCriteria()(WorkbookFilterCriteriaable) {
    return m.criteria
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookFilter) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["criteria"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookFilterCriteriaFromDiscriminatorValue , m.SetCriteria)
    return res
}
// Serialize serializes information the current object
func (m *WorkbookFilter) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("criteria", m.GetCriteria())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCriteria sets the criteria property value. The currently applied filter on the given column. Read-only.
func (m *WorkbookFilter) SetCriteria(value WorkbookFilterCriteriaable)() {
    m.criteria = value
}
