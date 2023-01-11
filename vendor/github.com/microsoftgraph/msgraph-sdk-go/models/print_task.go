package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrintTask provides operations to manage the collection of agreement entities.
type PrintTask struct {
    Entity
    // The definition property
    definition PrintTaskDefinitionable
    // The URL for the print entity that triggered this task. For example, https://graph.microsoft.com/v1.0/print/printers/{printerId}/jobs/{jobId}. Read-only.
    parentUrl *string
    // The status property
    status PrintTaskStatusable
    // The trigger property
    trigger PrintTaskTriggerable
}
// NewPrintTask instantiates a new printTask and sets the default values.
func NewPrintTask()(*PrintTask) {
    m := &PrintTask{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrintTaskFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrintTaskFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintTask(), nil
}
// GetDefinition gets the definition property value. The definition property
func (m *PrintTask) GetDefinition()(PrintTaskDefinitionable) {
    return m.definition
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PrintTask) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["definition"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePrintTaskDefinitionFromDiscriminatorValue , m.SetDefinition)
    res["parentUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetParentUrl)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePrintTaskStatusFromDiscriminatorValue , m.SetStatus)
    res["trigger"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePrintTaskTriggerFromDiscriminatorValue , m.SetTrigger)
    return res
}
// GetParentUrl gets the parentUrl property value. The URL for the print entity that triggered this task. For example, https://graph.microsoft.com/v1.0/print/printers/{printerId}/jobs/{jobId}. Read-only.
func (m *PrintTask) GetParentUrl()(*string) {
    return m.parentUrl
}
// GetStatus gets the status property value. The status property
func (m *PrintTask) GetStatus()(PrintTaskStatusable) {
    return m.status
}
// GetTrigger gets the trigger property value. The trigger property
func (m *PrintTask) GetTrigger()(PrintTaskTriggerable) {
    return m.trigger
}
// Serialize serializes information the current object
func (m *PrintTask) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("definition", m.GetDefinition())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("parentUrl", m.GetParentUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("trigger", m.GetTrigger())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDefinition sets the definition property value. The definition property
func (m *PrintTask) SetDefinition(value PrintTaskDefinitionable)() {
    m.definition = value
}
// SetParentUrl sets the parentUrl property value. The URL for the print entity that triggered this task. For example, https://graph.microsoft.com/v1.0/print/printers/{printerId}/jobs/{jobId}. Read-only.
func (m *PrintTask) SetParentUrl(value *string)() {
    m.parentUrl = value
}
// SetStatus sets the status property value. The status property
func (m *PrintTask) SetStatus(value PrintTaskStatusable)() {
    m.status = value
}
// SetTrigger sets the trigger property value. The trigger property
func (m *PrintTask) SetTrigger(value PrintTaskTriggerable)() {
    m.trigger = value
}
