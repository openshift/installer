package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DirectoryAudit provides operations to manage the collection of agreement entities.
type DirectoryAudit struct {
    Entity
    // Indicates the date and time the activity was performed. The Timestamp type is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
    activityDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Indicates the activity name or the operation name (examples: 'Create User' and 'Add member to group'). For full list, see Azure AD activity list.
    activityDisplayName *string
    // Indicates additional details on the activity.
    additionalDetails []KeyValueable
    // Indicates which resource category that's targeted by the activity. For example: UserManagement, GroupManagement, ApplicationManagement, RoleManagement.
    category *string
    // Indicates a unique ID that helps correlate activities that span across various services. Can be used to trace logs across services.
    correlationId *string
    // The initiatedBy property
    initiatedBy AuditActivityInitiatorable
    // Indicates information on which service initiated the activity (For example: Self-service Password Management, Core Directory, B2C, Invited Users, Microsoft Identity Manager, Privileged Identity Management.
    loggedByService *string
    // Indicates the type of operation that was performed. The possible values include but are not limited to the following: Add, Assign, Update, Unassign, and Delete.
    operationType *string
    // Indicates the result of the activity. Possible values are: success, failure, timeout, unknownFutureValue.
    result *OperationResult
    // Indicates the reason for failure if the result is failure or timeout.
    resultReason *string
    // Indicates information on which resource was changed due to the activity. Target Resource Type can be User, Device, Directory, App, Role, Group, Policy or Other.
    targetResources []TargetResourceable
}
// NewDirectoryAudit instantiates a new directoryAudit and sets the default values.
func NewDirectoryAudit()(*DirectoryAudit) {
    m := &DirectoryAudit{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDirectoryAuditFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDirectoryAuditFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectoryAudit(), nil
}
// GetActivityDateTime gets the activityDateTime property value. Indicates the date and time the activity was performed. The Timestamp type is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *DirectoryAudit) GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.activityDateTime
}
// GetActivityDisplayName gets the activityDisplayName property value. Indicates the activity name or the operation name (examples: 'Create User' and 'Add member to group'). For full list, see Azure AD activity list.
func (m *DirectoryAudit) GetActivityDisplayName()(*string) {
    return m.activityDisplayName
}
// GetAdditionalDetails gets the additionalDetails property value. Indicates additional details on the activity.
func (m *DirectoryAudit) GetAdditionalDetails()([]KeyValueable) {
    return m.additionalDetails
}
// GetCategory gets the category property value. Indicates which resource category that's targeted by the activity. For example: UserManagement, GroupManagement, ApplicationManagement, RoleManagement.
func (m *DirectoryAudit) GetCategory()(*string) {
    return m.category
}
// GetCorrelationId gets the correlationId property value. Indicates a unique ID that helps correlate activities that span across various services. Can be used to trace logs across services.
func (m *DirectoryAudit) GetCorrelationId()(*string) {
    return m.correlationId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DirectoryAudit) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activityDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetActivityDateTime)
    res["activityDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetActivityDisplayName)
    res["additionalDetails"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateKeyValueFromDiscriminatorValue , m.SetAdditionalDetails)
    res["category"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCategory)
    res["correlationId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCorrelationId)
    res["initiatedBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateAuditActivityInitiatorFromDiscriminatorValue , m.SetInitiatedBy)
    res["loggedByService"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLoggedByService)
    res["operationType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOperationType)
    res["result"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseOperationResult , m.SetResult)
    res["resultReason"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetResultReason)
    res["targetResources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTargetResourceFromDiscriminatorValue , m.SetTargetResources)
    return res
}
// GetInitiatedBy gets the initiatedBy property value. The initiatedBy property
func (m *DirectoryAudit) GetInitiatedBy()(AuditActivityInitiatorable) {
    return m.initiatedBy
}
// GetLoggedByService gets the loggedByService property value. Indicates information on which service initiated the activity (For example: Self-service Password Management, Core Directory, B2C, Invited Users, Microsoft Identity Manager, Privileged Identity Management.
func (m *DirectoryAudit) GetLoggedByService()(*string) {
    return m.loggedByService
}
// GetOperationType gets the operationType property value. Indicates the type of operation that was performed. The possible values include but are not limited to the following: Add, Assign, Update, Unassign, and Delete.
func (m *DirectoryAudit) GetOperationType()(*string) {
    return m.operationType
}
// GetResult gets the result property value. Indicates the result of the activity. Possible values are: success, failure, timeout, unknownFutureValue.
func (m *DirectoryAudit) GetResult()(*OperationResult) {
    return m.result
}
// GetResultReason gets the resultReason property value. Indicates the reason for failure if the result is failure or timeout.
func (m *DirectoryAudit) GetResultReason()(*string) {
    return m.resultReason
}
// GetTargetResources gets the targetResources property value. Indicates information on which resource was changed due to the activity. Target Resource Type can be User, Device, Directory, App, Role, Group, Policy or Other.
func (m *DirectoryAudit) GetTargetResources()([]TargetResourceable) {
    return m.targetResources
}
// Serialize serializes information the current object
func (m *DirectoryAudit) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("activityDateTime", m.GetActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityDisplayName", m.GetActivityDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetAdditionalDetails() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAdditionalDetails())
        err = writer.WriteCollectionOfObjectValues("additionalDetails", cast)
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
        err = writer.WriteStringValue("correlationId", m.GetCorrelationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("initiatedBy", m.GetInitiatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("loggedByService", m.GetLoggedByService())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("operationType", m.GetOperationType())
        if err != nil {
            return err
        }
    }
    if m.GetResult() != nil {
        cast := (*m.GetResult()).String()
        err = writer.WriteStringValue("result", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resultReason", m.GetResultReason())
        if err != nil {
            return err
        }
    }
    if m.GetTargetResources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTargetResources())
        err = writer.WriteCollectionOfObjectValues("targetResources", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivityDateTime sets the activityDateTime property value. Indicates the date and time the activity was performed. The Timestamp type is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *DirectoryAudit) SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.activityDateTime = value
}
// SetActivityDisplayName sets the activityDisplayName property value. Indicates the activity name or the operation name (examples: 'Create User' and 'Add member to group'). For full list, see Azure AD activity list.
func (m *DirectoryAudit) SetActivityDisplayName(value *string)() {
    m.activityDisplayName = value
}
// SetAdditionalDetails sets the additionalDetails property value. Indicates additional details on the activity.
func (m *DirectoryAudit) SetAdditionalDetails(value []KeyValueable)() {
    m.additionalDetails = value
}
// SetCategory sets the category property value. Indicates which resource category that's targeted by the activity. For example: UserManagement, GroupManagement, ApplicationManagement, RoleManagement.
func (m *DirectoryAudit) SetCategory(value *string)() {
    m.category = value
}
// SetCorrelationId sets the correlationId property value. Indicates a unique ID that helps correlate activities that span across various services. Can be used to trace logs across services.
func (m *DirectoryAudit) SetCorrelationId(value *string)() {
    m.correlationId = value
}
// SetInitiatedBy sets the initiatedBy property value. The initiatedBy property
func (m *DirectoryAudit) SetInitiatedBy(value AuditActivityInitiatorable)() {
    m.initiatedBy = value
}
// SetLoggedByService sets the loggedByService property value. Indicates information on which service initiated the activity (For example: Self-service Password Management, Core Directory, B2C, Invited Users, Microsoft Identity Manager, Privileged Identity Management.
func (m *DirectoryAudit) SetLoggedByService(value *string)() {
    m.loggedByService = value
}
// SetOperationType sets the operationType property value. Indicates the type of operation that was performed. The possible values include but are not limited to the following: Add, Assign, Update, Unassign, and Delete.
func (m *DirectoryAudit) SetOperationType(value *string)() {
    m.operationType = value
}
// SetResult sets the result property value. Indicates the result of the activity. Possible values are: success, failure, timeout, unknownFutureValue.
func (m *DirectoryAudit) SetResult(value *OperationResult)() {
    m.result = value
}
// SetResultReason sets the resultReason property value. Indicates the reason for failure if the result is failure or timeout.
func (m *DirectoryAudit) SetResultReason(value *string)() {
    m.resultReason = value
}
// SetTargetResources sets the targetResources property value. Indicates information on which resource was changed due to the activity. Target Resource Type can be User, Device, Directory, App, Role, Group, Policy or Other.
func (m *DirectoryAudit) SetTargetResources(value []TargetResourceable)() {
    m.targetResources = value
}
