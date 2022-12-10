package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PlannerPlan provides operations to manage the collection of agreement entities.
type PlannerPlan struct {
    Entity
    // Read-only. Nullable. Collection of buckets in the plan.
    buckets []PlannerBucketable
    // Identifies the container of the plan. After it is set, this property can’t be updated. Required.
    container PlannerPlanContainerable
    // Read-only. The user who created the plan.
    createdBy IdentitySetable
    // Read-only. Date and time at which the plan is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Read-only. Nullable. Additional details about the plan.
    details PlannerPlanDetailsable
    // The owner property
    owner *string
    // Read-only. Nullable. Collection of tasks in the plan.
    tasks []PlannerTaskable
    // Required. Title of the plan.
    title *string
}
// NewPlannerPlan instantiates a new plannerPlan and sets the default values.
func NewPlannerPlan()(*PlannerPlan) {
    m := &PlannerPlan{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerPlanFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerPlanFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerPlan(), nil
}
// GetBuckets gets the buckets property value. Read-only. Nullable. Collection of buckets in the plan.
func (m *PlannerPlan) GetBuckets()([]PlannerBucketable) {
    return m.buckets
}
// GetContainer gets the container property value. Identifies the container of the plan. After it is set, this property can’t be updated. Required.
func (m *PlannerPlan) GetContainer()(PlannerPlanContainerable) {
    return m.container
}
// GetCreatedBy gets the createdBy property value. Read-only. The user who created the plan.
func (m *PlannerPlan) GetCreatedBy()(IdentitySetable) {
    return m.createdBy
}
// GetCreatedDateTime gets the createdDateTime property value. Read-only. Date and time at which the plan is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerPlan) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDetails gets the details property value. Read-only. Nullable. Additional details about the plan.
func (m *PlannerPlan) GetDetails()(PlannerPlanDetailsable) {
    return m.details
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PlannerPlan) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["buckets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerBucketFromDiscriminatorValue , m.SetBuckets)
    res["container"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerPlanContainerFromDiscriminatorValue , m.SetContainer)
    res["createdBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetCreatedBy)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["details"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerPlanDetailsFromDiscriminatorValue , m.SetDetails)
    res["owner"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOwner)
    res["tasks"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePlannerTaskFromDiscriminatorValue , m.SetTasks)
    res["title"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTitle)
    return res
}
// GetOwner gets the owner property value. The owner property
func (m *PlannerPlan) GetOwner()(*string) {
    return m.owner
}
// GetTasks gets the tasks property value. Read-only. Nullable. Collection of tasks in the plan.
func (m *PlannerPlan) GetTasks()([]PlannerTaskable) {
    return m.tasks
}
// GetTitle gets the title property value. Required. Title of the plan.
func (m *PlannerPlan) GetTitle()(*string) {
    return m.title
}
// Serialize serializes information the current object
func (m *PlannerPlan) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    {
        err = writer.WriteObjectValue("container", m.GetContainer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("details", m.GetDetails())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("owner", m.GetOwner())
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
    {
        err = writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBuckets sets the buckets property value. Read-only. Nullable. Collection of buckets in the plan.
func (m *PlannerPlan) SetBuckets(value []PlannerBucketable)() {
    m.buckets = value
}
// SetContainer sets the container property value. Identifies the container of the plan. After it is set, this property can’t be updated. Required.
func (m *PlannerPlan) SetContainer(value PlannerPlanContainerable)() {
    m.container = value
}
// SetCreatedBy sets the createdBy property value. Read-only. The user who created the plan.
func (m *PlannerPlan) SetCreatedBy(value IdentitySetable)() {
    m.createdBy = value
}
// SetCreatedDateTime sets the createdDateTime property value. Read-only. Date and time at which the plan is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerPlan) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDetails sets the details property value. Read-only. Nullable. Additional details about the plan.
func (m *PlannerPlan) SetDetails(value PlannerPlanDetailsable)() {
    m.details = value
}
// SetOwner sets the owner property value. The owner property
func (m *PlannerPlan) SetOwner(value *string)() {
    m.owner = value
}
// SetTasks sets the tasks property value. Read-only. Nullable. Collection of tasks in the plan.
func (m *PlannerPlan) SetTasks(value []PlannerTaskable)() {
    m.tasks = value
}
// SetTitle sets the title property value. Required. Title of the plan.
func (m *PlannerPlan) SetTitle(value *string)() {
    m.title = value
}
