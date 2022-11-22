package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UnifiedRoleAssignmentSchedule 
type UnifiedRoleAssignmentSchedule struct {
    UnifiedRoleScheduleBase
    // If the request is from an eligible administrator to activate a role, this parameter will show the related eligible assignment for that activation. Otherwise, it is null. Supports $expand.
    activatedUsing UnifiedRoleEligibilityScheduleable
    // Type of the assignment which can either be Assigned or Activated. Supports $filter (eq, ne).
    assignmentType *string
    // How the assignments is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
    memberType *string
    // The period of the role assignment. It can represent a single occurrence or multiple recurrences.
    scheduleInfo RequestScheduleable
}
// NewUnifiedRoleAssignmentSchedule instantiates a new UnifiedRoleAssignmentSchedule and sets the default values.
func NewUnifiedRoleAssignmentSchedule()(*UnifiedRoleAssignmentSchedule) {
    m := &UnifiedRoleAssignmentSchedule{
        UnifiedRoleScheduleBase: *NewUnifiedRoleScheduleBase(),
    }
    return m
}
// CreateUnifiedRoleAssignmentScheduleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUnifiedRoleAssignmentScheduleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleAssignmentSchedule(), nil
}
// GetActivatedUsing gets the activatedUsing property value. If the request is from an eligible administrator to activate a role, this parameter will show the related eligible assignment for that activation. Otherwise, it is null. Supports $expand.
func (m *UnifiedRoleAssignmentSchedule) GetActivatedUsing()(UnifiedRoleEligibilityScheduleable) {
    return m.activatedUsing
}
// GetAssignmentType gets the assignmentType property value. Type of the assignment which can either be Assigned or Activated. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentSchedule) GetAssignmentType()(*string) {
    return m.assignmentType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UnifiedRoleAssignmentSchedule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UnifiedRoleScheduleBase.GetFieldDeserializers()
    res["activatedUsing"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUnifiedRoleEligibilityScheduleFromDiscriminatorValue , m.SetActivatedUsing)
    res["assignmentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAssignmentType)
    res["memberType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMemberType)
    res["scheduleInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRequestScheduleFromDiscriminatorValue , m.SetScheduleInfo)
    return res
}
// GetMemberType gets the memberType property value. How the assignments is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentSchedule) GetMemberType()(*string) {
    return m.memberType
}
// GetScheduleInfo gets the scheduleInfo property value. The period of the role assignment. It can represent a single occurrence or multiple recurrences.
func (m *UnifiedRoleAssignmentSchedule) GetScheduleInfo()(RequestScheduleable) {
    return m.scheduleInfo
}
// Serialize serializes information the current object
func (m *UnifiedRoleAssignmentSchedule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UnifiedRoleScheduleBase.Serialize(writer)
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
        err = writer.WriteStringValue("memberType", m.GetMemberType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("scheduleInfo", m.GetScheduleInfo())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivatedUsing sets the activatedUsing property value. If the request is from an eligible administrator to activate a role, this parameter will show the related eligible assignment for that activation. Otherwise, it is null. Supports $expand.
func (m *UnifiedRoleAssignmentSchedule) SetActivatedUsing(value UnifiedRoleEligibilityScheduleable)() {
    m.activatedUsing = value
}
// SetAssignmentType sets the assignmentType property value. Type of the assignment which can either be Assigned or Activated. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentSchedule) SetAssignmentType(value *string)() {
    m.assignmentType = value
}
// SetMemberType sets the memberType property value. How the assignments is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentSchedule) SetMemberType(value *string)() {
    m.memberType = value
}
// SetScheduleInfo sets the scheduleInfo property value. The period of the role assignment. It can represent a single occurrence or multiple recurrences.
func (m *UnifiedRoleAssignmentSchedule) SetScheduleInfo(value RequestScheduleable)() {
    m.scheduleInfo = value
}
