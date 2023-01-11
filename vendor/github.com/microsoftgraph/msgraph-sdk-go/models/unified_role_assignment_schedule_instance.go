package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UnifiedRoleAssignmentScheduleInstance 
type UnifiedRoleAssignmentScheduleInstance struct {
    UnifiedRoleScheduleInstanceBase
    // If the request is from an eligible administrator to activate a role, this parameter will show the related eligible assignment for that activation. Otherwise, it is null. Supports $expand.
    activatedUsing UnifiedRoleEligibilityScheduleInstanceable
    // Type of the assignment which can either be Assigned or Activated. Supports $filter (eq, ne).
    assignmentType *string
    // The end date of the schedule instance.
    endDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // How the assignments is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
    memberType *string
    // The identifier of the role assignment in Azure AD. Supports $filter (eq, ne).
    roleAssignmentOriginId *string
    // The identifier of the unifiedRoleAssignmentSchedule object from which this instance was created. Supports $filter (eq, ne).
    roleAssignmentScheduleId *string
    // When this instance starts.
    startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
}
// NewUnifiedRoleAssignmentScheduleInstance instantiates a new UnifiedRoleAssignmentScheduleInstance and sets the default values.
func NewUnifiedRoleAssignmentScheduleInstance()(*UnifiedRoleAssignmentScheduleInstance) {
    m := &UnifiedRoleAssignmentScheduleInstance{
        UnifiedRoleScheduleInstanceBase: *NewUnifiedRoleScheduleInstanceBase(),
    }
    return m
}
// CreateUnifiedRoleAssignmentScheduleInstanceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUnifiedRoleAssignmentScheduleInstanceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleAssignmentScheduleInstance(), nil
}
// GetActivatedUsing gets the activatedUsing property value. If the request is from an eligible administrator to activate a role, this parameter will show the related eligible assignment for that activation. Otherwise, it is null. Supports $expand.
func (m *UnifiedRoleAssignmentScheduleInstance) GetActivatedUsing()(UnifiedRoleEligibilityScheduleInstanceable) {
    return m.activatedUsing
}
// GetAssignmentType gets the assignmentType property value. Type of the assignment which can either be Assigned or Activated. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) GetAssignmentType()(*string) {
    return m.assignmentType
}
// GetEndDateTime gets the endDateTime property value. The end date of the schedule instance.
func (m *UnifiedRoleAssignmentScheduleInstance) GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.endDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UnifiedRoleAssignmentScheduleInstance) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UnifiedRoleScheduleInstanceBase.GetFieldDeserializers()
    res["activatedUsing"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUnifiedRoleEligibilityScheduleInstanceFromDiscriminatorValue , m.SetActivatedUsing)
    res["assignmentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAssignmentType)
    res["endDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetEndDateTime)
    res["memberType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMemberType)
    res["roleAssignmentOriginId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRoleAssignmentOriginId)
    res["roleAssignmentScheduleId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRoleAssignmentScheduleId)
    res["startDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetStartDateTime)
    return res
}
// GetMemberType gets the memberType property value. How the assignments is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) GetMemberType()(*string) {
    return m.memberType
}
// GetRoleAssignmentOriginId gets the roleAssignmentOriginId property value. The identifier of the role assignment in Azure AD. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) GetRoleAssignmentOriginId()(*string) {
    return m.roleAssignmentOriginId
}
// GetRoleAssignmentScheduleId gets the roleAssignmentScheduleId property value. The identifier of the unifiedRoleAssignmentSchedule object from which this instance was created. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) GetRoleAssignmentScheduleId()(*string) {
    return m.roleAssignmentScheduleId
}
// GetStartDateTime gets the startDateTime property value. When this instance starts.
func (m *UnifiedRoleAssignmentScheduleInstance) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.startDateTime
}
// Serialize serializes information the current object
func (m *UnifiedRoleAssignmentScheduleInstance) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UnifiedRoleScheduleInstanceBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("activatedUsing", m.GetActivatedUsing())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("assignmentType", m.GetAssignmentType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("endDateTime", m.GetEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("memberType", m.GetMemberType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleAssignmentOriginId", m.GetRoleAssignmentOriginId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleAssignmentScheduleId", m.GetRoleAssignmentScheduleId())
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
    return nil
}
// SetActivatedUsing sets the activatedUsing property value. If the request is from an eligible administrator to activate a role, this parameter will show the related eligible assignment for that activation. Otherwise, it is null. Supports $expand.
func (m *UnifiedRoleAssignmentScheduleInstance) SetActivatedUsing(value UnifiedRoleEligibilityScheduleInstanceable)() {
    m.activatedUsing = value
}
// SetAssignmentType sets the assignmentType property value. Type of the assignment which can either be Assigned or Activated. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetAssignmentType(value *string)() {
    m.assignmentType = value
}
// SetEndDateTime sets the endDateTime property value. The end date of the schedule instance.
func (m *UnifiedRoleAssignmentScheduleInstance) SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.endDateTime = value
}
// SetMemberType sets the memberType property value. How the assignments is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetMemberType(value *string)() {
    m.memberType = value
}
// SetRoleAssignmentOriginId sets the roleAssignmentOriginId property value. The identifier of the role assignment in Azure AD. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetRoleAssignmentOriginId(value *string)() {
    m.roleAssignmentOriginId = value
}
// SetRoleAssignmentScheduleId sets the roleAssignmentScheduleId property value. The identifier of the unifiedRoleAssignmentSchedule object from which this instance was created. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetRoleAssignmentScheduleId(value *string)() {
    m.roleAssignmentScheduleId = value
}
// SetStartDateTime sets the startDateTime property value. When this instance starts.
func (m *UnifiedRoleAssignmentScheduleInstance) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.startDateTime = value
}
