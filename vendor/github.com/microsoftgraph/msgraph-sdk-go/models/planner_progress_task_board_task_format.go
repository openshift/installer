package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PlannerProgressTaskBoardTaskFormat 
type PlannerProgressTaskBoardTaskFormat struct {
    Entity
    // Hint value used to order the task on the Progress view of the Task Board. The format is defined as outlined here.
    orderHint *string
}
// NewPlannerProgressTaskBoardTaskFormat instantiates a new plannerProgressTaskBoardTaskFormat and sets the default values.
func NewPlannerProgressTaskBoardTaskFormat()(*PlannerProgressTaskBoardTaskFormat) {
    m := &PlannerProgressTaskBoardTaskFormat{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerProgressTaskBoardTaskFormatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerProgressTaskBoardTaskFormatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerProgressTaskBoardTaskFormat(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PlannerProgressTaskBoardTaskFormat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["orderHint"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOrderHint)
    return res
}
// GetOrderHint gets the orderHint property value. Hint value used to order the task on the Progress view of the Task Board. The format is defined as outlined here.
func (m *PlannerProgressTaskBoardTaskFormat) GetOrderHint()(*string) {
    return m.orderHint
}
// Serialize serializes information the current object
func (m *PlannerProgressTaskBoardTaskFormat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("orderHint", m.GetOrderHint())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOrderHint sets the orderHint property value. Hint value used to order the task on the Progress view of the Task Board. The format is defined as outlined here.
func (m *PlannerProgressTaskBoardTaskFormat) SetOrderHint(value *string)() {
    m.orderHint = value
}
