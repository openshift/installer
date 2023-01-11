package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Planner 
type Planner struct {
    Entity
    // Read-only. Nullable. Returns a collection of the specified buckets
    buckets []PlannerBucketable
    // Read-only. Nullable. Returns a collection of the specified plans
    plans []PlannerPlanable
    // Read-only. Nullable. Returns a collection of the specified tasks
    tasks []PlannerTaskable
}
// NewPlanner instantiates a new Planner and sets the default values.
func NewPlanner()(*Planner) {
    m := &Planner{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlanner(), nil
}
// GetBuckets gets the buckets property value. Read-only. Nullable. Returns a collection of the specified buckets
func (m *Planner) GetBuckets()([]PlannerBucketable) {
    return m.buckets
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Planner) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["buckets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerBucketFromDiscriminatorValue , m.SetBuckets)
    res["plans"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerPlanFromDiscriminatorValue , m.SetPlans)
    res["tasks"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerTaskFromDiscriminatorValue , m.SetTasks)
    return res
}
// GetPlans gets the plans property value. Read-only. Nullable. Returns a collection of the specified plans
func (m *Planner) GetPlans()([]PlannerPlanable) {
    return m.plans
}
// GetTasks gets the tasks property value. Read-only. Nullable. Returns a collection of the specified tasks
func (m *Planner) GetTasks()([]PlannerTaskable) {
    return m.tasks
}
// Serialize serializes information the current object
func (m *Planner) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetBuckets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetBuckets())
        err = writer.WriteCollectionOfObjectValues("buckets", cast)
        if err != nil {
            return err
        }
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
// SetBuckets sets the buckets property value. Read-only. Nullable. Returns a collection of the specified buckets
func (m *Planner) SetBuckets(value []PlannerBucketable)() {
    m.buckets = value
}
// SetPlans sets the plans property value. Read-only. Nullable. Returns a collection of the specified plans
func (m *Planner) SetPlans(value []PlannerPlanable)() {
    m.plans = value
}
// SetTasks sets the tasks property value. Read-only. Nullable. Returns a collection of the specified tasks
func (m *Planner) SetTasks(value []PlannerTaskable)() {
    m.tasks = value
}
