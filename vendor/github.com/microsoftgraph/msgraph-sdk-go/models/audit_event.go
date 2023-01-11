package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AuditEvent 
type AuditEvent struct {
    Entity
    // Friendly name of the activity.
    activity *string
    // The date time in UTC when the activity was performed.
    activityDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The HTTP operation type of the activity.
    activityOperationType *string
    // The result of the activity.
    activityResult *string
    // The type of activity that was being performed.
    activityType *string
    // AAD user and application that are associated with the audit event.
    actor AuditActorable
    // Audit category.
    category *string
    // Component name.
    componentName *string
    // The client request Id that is used to correlate activity within the system.
    correlationId *string
    // Event display name.
    displayName *string
    // Resources being modified.
    resources []AuditResourceable
}
// NewAuditEvent instantiates a new AuditEvent and sets the default values.
func NewAuditEvent()(*AuditEvent) {
    m := &AuditEvent{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuditEventFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAuditEventFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuditEvent(), nil
}
// GetActivity gets the activity property value. Friendly name of the activity.
func (m *AuditEvent) GetActivity()(*string) {
    return m.activity
}
// GetActivityDateTime gets the activityDateTime property value. The date time in UTC when the activity was performed.
func (m *AuditEvent) GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.activityDateTime
}
// GetActivityOperationType gets the activityOperationType property value. The HTTP operation type of the activity.
func (m *AuditEvent) GetActivityOperationType()(*string) {
    return m.activityOperationType
}
// GetActivityResult gets the activityResult property value. The result of the activity.
func (m *AuditEvent) GetActivityResult()(*string) {
    return m.activityResult
}
// GetActivityType gets the activityType property value. The type of activity that was being performed.
func (m *AuditEvent) GetActivityType()(*string) {
    return m.activityType
}
// GetActor gets the actor property value. AAD user and application that are associated with the audit event.
func (m *AuditEvent) GetActor()(AuditActorable) {
    return m.actor
}
// GetCategory gets the category property value. Audit category.
func (m *AuditEvent) GetCategory()(*string) {
    return m.category
}
// GetComponentName gets the componentName property value. Component name.
func (m *AuditEvent) GetComponentName()(*string) {
    return m.componentName
}
// GetCorrelationId gets the correlationId property value. The client request Id that is used to correlate activity within the system.
func (m *AuditEvent) GetCorrelationId()(*string) {
    return m.correlationId
}
// GetDisplayName gets the displayName property value. Event display name.
func (m *AuditEvent) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AuditEvent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetActivity)
    res["activityDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetActivityDateTime)
    res["activityOperationType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetActivityOperationType)
    res["activityResult"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetActivityResult)
    res["activityType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetActivityType)
    res["actor"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateAuditActorFromDiscriminatorValue , m.SetActor)
    res["category"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCategory)
    res["componentName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetComponentName)
    res["correlationId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCorrelationId)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["resources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAuditResourceFromDiscriminatorValue , m.SetResources)
    return res
}
// GetResources gets the resources property value. Resources being modified.
func (m *AuditEvent) GetResources()([]AuditResourceable) {
    return m.resources
}
// Serialize serializes information the current object
func (m *AuditEvent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("activity", m.GetActivity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("activityDateTime", m.GetActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityOperationType", m.GetActivityOperationType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityResult", m.GetActivityResult())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityType", m.GetActivityType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("actor", m.GetActor())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("category", m.GetCategory())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("componentName", m.GetComponentName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("correlationId", m.GetCorrelationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetResources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetResources())
        err = writer.WriteCollectionOfObjectValues("resources", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivity sets the activity property value. Friendly name of the activity.
func (m *AuditEvent) SetActivity(value *string)() {
    m.activity = value
}
// SetActivityDateTime sets the activityDateTime property value. The date time in UTC when the activity was performed.
func (m *AuditEvent) SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.activityDateTime = value
}
// SetActivityOperationType sets the activityOperationType property value. The HTTP operation type of the activity.
func (m *AuditEvent) SetActivityOperationType(value *string)() {
    m.activityOperationType = value
}
// SetActivityResult sets the activityResult property value. The result of the activity.
func (m *AuditEvent) SetActivityResult(value *string)() {
    m.activityResult = value
}
// SetActivityType sets the activityType property value. The type of activity that was being performed.
func (m *AuditEvent) SetActivityType(value *string)() {
    m.activityType = value
}
// SetActor sets the actor property value. AAD user and application that are associated with the audit event.
func (m *AuditEvent) SetActor(value AuditActorable)() {
    m.actor = value
}
// SetCategory sets the category property value. Audit category.
func (m *AuditEvent) SetCategory(value *string)() {
    m.category = value
}
// SetComponentName sets the componentName property value. Component name.
func (m *AuditEvent) SetComponentName(value *string)() {
    m.componentName = value
}
// SetCorrelationId sets the correlationId property value. The client request Id that is used to correlate activity within the system.
func (m *AuditEvent) SetCorrelationId(value *string)() {
    m.correlationId = value
}
// SetDisplayName sets the displayName property value. Event display name.
func (m *AuditEvent) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetResources sets the resources property value. Resources being modified.
func (m *AuditEvent) SetResources(value []AuditResourceable)() {
    m.resources = value
}
