package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PlannerGroup 
type PlannerGroup struct {
    Entity
    // Read-only. Nullable. Returns the plannerPlans owned by the group.
    plans []PlannerPlanable
}
// NewPlannerGroup instantiates a new plannerGroup and sets the default values.
func NewPlannerGroup()(*PlannerGroup) {
    m := &PlannerGroup{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerGroupFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerGroupFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerGroup(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PlannerGroup) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["plans"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerPlanFromDiscriminatorValue , m.SetPlans)
    return res
}
// GetPlans gets the plans property value. Read-only. Nullable. Returns the plannerPlans owned by the group.
func (m *PlannerGroup) GetPlans()([]PlannerPlanable) {
    return m.plans
}
// Serialize serializes information the current object
func (m *PlannerGroup) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    return nil
}
// SetPlans sets the plans property value. Read-only. Nullable. Returns the plannerPlans owned by the group.
func (m *PlannerGroup) SetPlans(value []PlannerPlanable)() {
    m.plans = value
}
