package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PlannerUser 
type PlannerUser struct {
    Entity
    // Read-only. Nullable. Returns the plannerTasks assigned to the user.
    plans []PlannerPlanable
    // Read-only. Nullable. Returns the plannerPlans shared with the user.
    tasks []PlannerTaskable
}
// NewPlannerUser instantiates a new plannerUser and sets the default values.
func NewPlannerUser()(*PlannerUser) {
    m := &PlannerUser{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerUser(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PlannerUser) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["plans"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerPlanFromDiscriminatorValue , m.SetPlans)
    res["tasks"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerTaskFromDiscriminatorValue , m.SetTasks)
    return res
}
// GetPlans gets the plans property value. Read-only. Nullable. Returns the plannerTasks assigned to the user.
func (m *PlannerUser) GetPlans()([]PlannerPlanable) {
    return m.plans
}
// GetTasks gets the tasks property value. Read-only. Nullable. Returns the plannerPlans shared with the user.
func (m *PlannerUser) GetTasks()([]PlannerTaskable) {
    return m.tasks
}
// Serialize serializes information the current object
func (m *PlannerUser) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetPlans() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPlans())
        err = writer.WriteCollectionOfObjectValues("plans", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTasks() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTasks())
        err = writer.WriteCollectionOfObjectValues("tasks", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPlans sets the plans property value. Read-only. Nullable. Returns the plannerTasks assigned to the user.
func (m *PlannerUser) SetPlans(value []PlannerPlanable)() {
    m.plans = value
}
// SetTasks sets the tasks property value. Read-only. Nullable. Returns the plannerPlans shared with the user.
func (m *PlannerUser) SetTasks(value []PlannerTaskable)() {
    m.tasks = value
}
