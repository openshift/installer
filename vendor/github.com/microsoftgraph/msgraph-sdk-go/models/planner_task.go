package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PlannerTask provides operations to manage the collection of agreement entities.
type PlannerTask struct {
    Entity
    // Number of checklist items with value set to false, representing incomplete items.
    activeChecklistItemCount *int32
    // The categories to which the task has been applied. See applied Categories for possible values.
    appliedCategories PlannerAppliedCategoriesable
    // Read-only. Nullable. Used to render the task correctly in the task board view when grouped by assignedTo.
    assignedToTaskBoardFormat PlannerAssignedToTaskBoardTaskFormatable
    // Hint used to order items of this type in a list view. The format is defined as outlined here.
    assigneePriority *string
    // The set of assignees the task is assigned to.
    assignments PlannerAssignmentsable
    // Bucket ID to which the task belongs. The bucket needs to be in the plan that the task is in. It is 28 characters long and case-sensitive. Format validation is done on the service.
    bucketId *string
    // Read-only. Nullable. Used to render the task correctly in the task board view when grouped by bucket.
    bucketTaskBoardFormat PlannerBucketTaskBoardTaskFormatable
    // Number of checklist items that are present on the task.
    checklistItemCount *int32
    // Identity of the user that completed the task.
    completedBy IdentitySetable
    // Read-only. Date and time at which the 'percentComplete' of the task is set to '100'. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    completedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Thread ID of the conversation on the task. This is the ID of the conversation thread object created in the group.
    conversationThreadId *string
    // Identity of the user that created the task.
    createdBy IdentitySetable
    // Read-only. Date and time at which the task is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Read-only. Nullable. Additional details about the task.
    details PlannerTaskDetailsable
    // Date and time at which the task is due. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    dueDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Read-only. Value is true if the details object of the task has a non-empty description and false otherwise.
    hasDescription *bool
    // Hint used to order items of this type in a list view. The format is defined as outlined here.
    orderHint *string
    // Percentage of task completion. When set to 100, the task is considered completed.
    percentComplete *int32
    // Plan ID to which the task belongs.
    planId *string
    // This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference.
    previewType *PlannerPreviewType
    // Priority of the task. The valid range of values is between 0 and 10, with the increasing value being lower priority (0 has the highest priority and 10 has the lowest priority).  Currently, Planner interprets values 0 and 1 as 'urgent', 2, 3 and 4 as 'important', 5, 6, and 7 as 'medium', and 8, 9, and 10 as 'low'.  Additionally, Planner sets the value 1 for 'urgent', 3 for 'important', 5 for 'medium', and 9 for 'low'.
    priority *int32
    // Read-only. Nullable. Used to render the task correctly in the task board view when grouped by progress.
    progressTaskBoardFormat PlannerProgressTaskBoardTaskFormatable
    // Number of external references that exist on the task.
    referenceCount *int32
    // Date and time at which the task starts. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Title of the task.
    title *string
}
// NewPlannerTask instantiates a new plannerTask and sets the default values.
func NewPlannerTask()(*PlannerTask) {
    m := &PlannerTask{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerTaskFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerTaskFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerTask(), nil
}
// GetActiveChecklistItemCount gets the activeChecklistItemCount property value. Number of checklist items with value set to false, representing incomplete items.
func (m *PlannerTask) GetActiveChecklistItemCount()(*int32) {
    return m.activeChecklistItemCount
}
// GetAppliedCategories gets the appliedCategories property value. The categories to which the task has been applied. See applied Categories for possible values.
func (m *PlannerTask) GetAppliedCategories()(PlannerAppliedCategoriesable) {
    return m.appliedCategories
}
// GetAssignedToTaskBoardFormat gets the assignedToTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by assignedTo.
func (m *PlannerTask) GetAssignedToTaskBoardFormat()(PlannerAssignedToTaskBoardTaskFormatable) {
    return m.assignedToTaskBoardFormat
}
// GetAssigneePriority gets the assigneePriority property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
func (m *PlannerTask) GetAssigneePriority()(*string) {
    return m.assigneePriority
}
// GetAssignments gets the assignments property value. The set of assignees the task is assigned to.
func (m *PlannerTask) GetAssignments()(PlannerAssignmentsable) {
    return m.assignments
}
// GetBucketId gets the bucketId property value. Bucket ID to which the task belongs. The bucket needs to be in the plan that the task is in. It is 28 characters long and case-sensitive. Format validation is done on the service.
func (m *PlannerTask) GetBucketId()(*string) {
    return m.bucketId
}
// GetBucketTaskBoardFormat gets the bucketTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by bucket.
func (m *PlannerTask) GetBucketTaskBoardFormat()(PlannerBucketTaskBoardTaskFormatable) {
    return m.bucketTaskBoardFormat
}
// GetChecklistItemCount gets the checklistItemCount property value. Number of checklist items that are present on the task.
func (m *PlannerTask) GetChecklistItemCount()(*int32) {
    return m.checklistItemCount
}
// GetCompletedBy gets the completedBy property value. Identity of the user that completed the task.
func (m *PlannerTask) GetCompletedBy()(IdentitySetable) {
    return m.completedBy
}
// GetCompletedDateTime gets the completedDateTime property value. Read-only. Date and time at which the 'percentComplete' of the task is set to '100'. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.completedDateTime
}
// GetConversationThreadId gets the conversationThreadId property value. Thread ID of the conversation on the task. This is the ID of the conversation thread object created in the group.
func (m *PlannerTask) GetConversationThreadId()(*string) {
    return m.conversationThreadId
}
// GetCreatedBy gets the createdBy property value. Identity of the user that created the task.
func (m *PlannerTask) GetCreatedBy()(IdentitySetable) {
    return m.createdBy
}
// GetCreatedDateTime gets the createdDateTime property value. Read-only. Date and time at which the task is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDetails gets the details property value. Read-only. Nullable. Additional details about the task.
func (m *PlannerTask) GetDetails()(PlannerTaskDetailsable) {
    return m.details
}
// GetDueDateTime gets the dueDateTime property value. Date and time at which the task is due. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) GetDueDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.dueDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PlannerTask) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activeChecklistItemCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetActiveChecklistItemCount)
    res["appliedCategories"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerAppliedCategoriesFromDiscriminatorValue , m.SetAppliedCategories)
    res["assignedToTaskBoardFormat"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerAssignedToTaskBoardTaskFormatFromDiscriminatorValue , m.SetAssignedToTaskBoardFormat)
    res["assigneePriority"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAssigneePriority)
    res["assignments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerAssignmentsFromDiscriminatorValue , m.SetAssignments)
    res["bucketId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetBucketId)
    res["bucketTaskBoardFormat"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerBucketTaskBoardTaskFormatFromDiscriminatorValue , m.SetBucketTaskBoardFormat)
    res["checklistItemCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetChecklistItemCount)
    res["completedBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetCompletedBy)
    res["completedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCompletedDateTime)
    res["conversationThreadId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetConversationThreadId)
    res["createdBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetCreatedBy)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["details"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerTaskDetailsFromDiscriminatorValue , m.SetDetails)
    res["dueDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetDueDateTime)
    res["hasDescription"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHasDescription)
    res["orderHint"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOrderHint)
    res["percentComplete"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPercentComplete)
    res["planId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPlanId)
    res["previewType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePlannerPreviewType , m.SetPreviewType)
    res["priority"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetPriority)
    res["progressTaskBoardFormat"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePlannerProgressTaskBoardTaskFormatFromDiscriminatorValue , m.SetProgressTaskBoardFormat)
    res["referenceCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetReferenceCount)
    res["startDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetStartDateTime)
    res["title"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTitle)
    return res
}
// GetHasDescription gets the hasDescription property value. Read-only. Value is true if the details object of the task has a non-empty description and false otherwise.
func (m *PlannerTask) GetHasDescription()(*bool) {
    return m.hasDescription
}
// GetOrderHint gets the orderHint property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
func (m *PlannerTask) GetOrderHint()(*string) {
    return m.orderHint
}
// GetPercentComplete gets the percentComplete property value. Percentage of task completion. When set to 100, the task is considered completed.
func (m *PlannerTask) GetPercentComplete()(*int32) {
    return m.percentComplete
}
// GetPlanId gets the planId property value. Plan ID to which the task belongs.
func (m *PlannerTask) GetPlanId()(*string) {
    return m.planId
}
// GetPreviewType gets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference.
func (m *PlannerTask) GetPreviewType()(*PlannerPreviewType) {
    return m.previewType
}
// GetPriority gets the priority property value. Priority of the task. The valid range of values is between 0 and 10, with the increasing value being lower priority (0 has the highest priority and 10 has the lowest priority).  Currently, Planner interprets values 0 and 1 as 'urgent', 2, 3 and 4 as 'important', 5, 6, and 7 as 'medium', and 8, 9, and 10 as 'low'.  Additionally, Planner sets the value 1 for 'urgent', 3 for 'important', 5 for 'medium', and 9 for 'low'.
func (m *PlannerTask) GetPriority()(*int32) {
    return m.priority
}
// GetProgressTaskBoardFormat gets the progressTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by progress.
func (m *PlannerTask) GetProgressTaskBoardFormat()(PlannerProgressTaskBoardTaskFormatable) {
    return m.progressTaskBoardFormat
}
// GetReferenceCount gets the referenceCount property value. Number of external references that exist on the task.
func (m *PlannerTask) GetReferenceCount()(*int32) {
    return m.referenceCount
}
// GetStartDateTime gets the startDateTime property value. Date and time at which the task starts. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.startDateTime
}
// GetTitle gets the title property value. Title of the task.
func (m *PlannerTask) GetTitle()(*string) {
    return m.title
}
// Serialize serializes information the current object
func (m *PlannerTask) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("activeChecklistItemCount", m.GetActiveChecklistItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("appliedCategories", m.GetAppliedCategories())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignedToTaskBoardFormat", m.GetAssignedToTaskBoardFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("assigneePriority", m.GetAssigneePriority())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignments", m.GetAssignments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("bucketId", m.GetBucketId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("bucketTaskBoardFormat", m.GetBucketTaskBoardFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("checklistItemCount", m.GetChecklistItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("completedBy", m.GetCompletedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("completedDateTime", m.GetCompletedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("conversationThreadId", m.GetConversationThreadId())
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
        err = writer.WriteTimeValue("dueDateTime", m.GetDueDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasDescription", m.GetHasDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("orderHint", m.GetOrderHint())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("percentComplete", m.GetPercentComplete())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("planId", m.GetPlanId())
        if err != nil {
            return err
        }
    }
    if m.GetPreviewType() != nil {
        cast := (*m.GetPreviewType()).String()
        err = writer.WriteStringValue("previewType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("priority", m.GetPriority())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("progressTaskBoardFormat", m.GetProgressTaskBoardFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("referenceCount", m.GetReferenceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
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
// SetActiveChecklistItemCount sets the activeChecklistItemCount property value. Number of checklist items with value set to false, representing incomplete items.
func (m *PlannerTask) SetActiveChecklistItemCount(value *int32)() {
    m.activeChecklistItemCount = value
}
// SetAppliedCategories sets the appliedCategories property value. The categories to which the task has been applied. See applied Categories for possible values.
func (m *PlannerTask) SetAppliedCategories(value PlannerAppliedCategoriesable)() {
    m.appliedCategories = value
}
// SetAssignedToTaskBoardFormat sets the assignedToTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by assignedTo.
func (m *PlannerTask) SetAssignedToTaskBoardFormat(value PlannerAssignedToTaskBoardTaskFormatable)() {
    m.assignedToTaskBoardFormat = value
}
// SetAssigneePriority sets the assigneePriority property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
func (m *PlannerTask) SetAssigneePriority(value *string)() {
    m.assigneePriority = value
}
// SetAssignments sets the assignments property value. The set of assignees the task is assigned to.
func (m *PlannerTask) SetAssignments(value PlannerAssignmentsable)() {
    m.assignments = value
}
// SetBucketId sets the bucketId property value. Bucket ID to which the task belongs. The bucket needs to be in the plan that the task is in. It is 28 characters long and case-sensitive. Format validation is done on the service.
func (m *PlannerTask) SetBucketId(value *string)() {
    m.bucketId = value
}
// SetBucketTaskBoardFormat sets the bucketTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by bucket.
func (m *PlannerTask) SetBucketTaskBoardFormat(value PlannerBucketTaskBoardTaskFormatable)() {
    m.bucketTaskBoardFormat = value
}
// SetChecklistItemCount sets the checklistItemCount property value. Number of checklist items that are present on the task.
func (m *PlannerTask) SetChecklistItemCount(value *int32)() {
    m.checklistItemCount = value
}
// SetCompletedBy sets the completedBy property value. Identity of the user that completed the task.
func (m *PlannerTask) SetCompletedBy(value IdentitySetable)() {
    m.completedBy = value
}
// SetCompletedDateTime sets the completedDateTime property value. Read-only. Date and time at which the 'percentComplete' of the task is set to '100'. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.completedDateTime = value
}
// SetConversationThreadId sets the conversationThreadId property value. Thread ID of the conversation on the task. This is the ID of the conversation thread object created in the group.
func (m *PlannerTask) SetConversationThreadId(value *string)() {
    m.conversationThreadId = value
}
// SetCreatedBy sets the createdBy property value. Identity of the user that created the task.
func (m *PlannerTask) SetCreatedBy(value IdentitySetable)() {
    m.createdBy = value
}
// SetCreatedDateTime sets the createdDateTime property value. Read-only. Date and time at which the task is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDetails sets the details property value. Read-only. Nullable. Additional details about the task.
func (m *PlannerTask) SetDetails(value PlannerTaskDetailsable)() {
    m.details = value
}
// SetDueDateTime sets the dueDateTime property value. Date and time at which the task is due. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetDueDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.dueDateTime = value
}
// SetHasDescription sets the hasDescription property value. Read-only. Value is true if the details object of the task has a non-empty description and false otherwise.
func (m *PlannerTask) SetHasDescription(value *bool)() {
    m.hasDescription = value
}
// SetOrderHint sets the orderHint property value. Hint used to order items of this type in a list view. The format is defined as outlined here.
func (m *PlannerTask) SetOrderHint(value *string)() {
    m.orderHint = value
}
// SetPercentComplete sets the percentComplete property value. Percentage of task completion. When set to 100, the task is considered completed.
func (m *PlannerTask) SetPercentComplete(value *int32)() {
    m.percentComplete = value
}
// SetPlanId sets the planId property value. Plan ID to which the task belongs.
func (m *PlannerTask) SetPlanId(value *string)() {
    m.planId = value
}
// SetPreviewType sets the previewType property value. This sets the type of preview that shows up on the task. The possible values are: automatic, noPreview, checklist, description, reference.
func (m *PlannerTask) SetPreviewType(value *PlannerPreviewType)() {
    m.previewType = value
}
// SetPriority sets the priority property value. Priority of the task. The valid range of values is between 0 and 10, with the increasing value being lower priority (0 has the highest priority and 10 has the lowest priority).  Currently, Planner interprets values 0 and 1 as 'urgent', 2, 3 and 4 as 'important', 5, 6, and 7 as 'medium', and 8, 9, and 10 as 'low'.  Additionally, Planner sets the value 1 for 'urgent', 3 for 'important', 5 for 'medium', and 9 for 'low'.
func (m *PlannerTask) SetPriority(value *int32)() {
    m.priority = value
}
// SetProgressTaskBoardFormat sets the progressTaskBoardFormat property value. Read-only. Nullable. Used to render the task correctly in the task board view when grouped by progress.
func (m *PlannerTask) SetProgressTaskBoardFormat(value PlannerProgressTaskBoardTaskFormatable)() {
    m.progressTaskBoardFormat = value
}
// SetReferenceCount sets the referenceCount property value. Number of external references that exist on the task.
func (m *PlannerTask) SetReferenceCount(value *int32)() {
    m.referenceCount = value
}
// SetStartDateTime sets the startDateTime property value. Date and time at which the task starts. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *PlannerTask) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.startDateTime = value
}
// SetTitle sets the title property value. Title of the task.
func (m *PlannerTask) SetTitle(value *string)() {
    m.title = value
}
