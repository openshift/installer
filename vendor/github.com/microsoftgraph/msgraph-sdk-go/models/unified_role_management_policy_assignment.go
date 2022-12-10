package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UnifiedRoleManagementPolicyAssignment provides operations to manage the collection of agreement entities.
type UnifiedRoleManagementPolicyAssignment struct {
    Entity
    // The policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
    policy UnifiedRoleManagementPolicyable
    // The id of the policy. Inherited from entity.
    policyId *string
    // The identifier of the role definition object where the policy applies. If not specified, the policy applies to all roles. Supports $filter (eq).
    roleDefinitionId *string
    // The identifier of the scope where the policy is assigned.  Can be / for the tenant or a group ID. Required.
    scopeId *string
    // The type of the scope where the policy is assigned. One of Directory, DirectoryRole. Required.
    scopeType *string
}
// NewUnifiedRoleManagementPolicyAssignment instantiates a new unifiedRoleManagementPolicyAssignment and sets the default values.
func NewUnifiedRoleManagementPolicyAssignment()(*UnifiedRoleManagementPolicyAssignment) {
    m := &UnifiedRoleManagementPolicyAssignment{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUnifiedRoleManagementPolicyAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUnifiedRoleManagementPolicyAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleManagementPolicyAssignment(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UnifiedRoleManagementPolicyAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["policy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUnifiedRoleManagementPolicyFromDiscriminatorValue , m.SetPolicy)
    res["policyId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPolicyId)
    res["roleDefinitionId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRoleDefinitionId)
    res["scopeId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetScopeId)
    res["scopeType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetScopeType)
    return res
}
// GetPolicy gets the policy property value. The policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
func (m *UnifiedRoleManagementPolicyAssignment) GetPolicy()(UnifiedRoleManagementPolicyable) {
    return m.policy
}
// GetPolicyId gets the policyId property value. The id of the policy. Inherited from entity.
func (m *UnifiedRoleManagementPolicyAssignment) GetPolicyId()(*string) {
    return m.policyId
}
// GetRoleDefinitionId gets the roleDefinitionId property value. The identifier of the role definition object where the policy applies. If not specified, the policy applies to all roles. Supports $filter (eq).
func (m *UnifiedRoleManagementPolicyAssignment) GetRoleDefinitionId()(*string) {
    return m.roleDefinitionId
}
// GetScopeId gets the scopeId property value. The identifier of the scope where the policy is assigned.  Can be / for the tenant or a group ID. Required.
func (m *UnifiedRoleManagementPolicyAssignment) GetScopeId()(*string) {
    return m.scopeId
}
// GetScopeType gets the scopeType property value. The type of the scope where the policy is assigned. One of Directory, DirectoryRole. Required.
func (m *UnifiedRoleManagementPolicyAssignment) GetScopeType()(*string) {
    return m.scopeType
}
// Serialize serializes information the current object
func (m *UnifiedRoleManagementPolicyAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("policy", m.GetPolicy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("policyId", m.GetPolicyId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleDefinitionId", m.GetRoleDefinitionId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("scopeId", m.GetScopeId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("scopeType", m.GetScopeType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPolicy sets the policy property value. The policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
func (m *UnifiedRoleManagementPolicyAssignment) SetPolicy(value UnifiedRoleManagementPolicyable)() {
    m.policy = value
}
// SetPolicyId sets the policyId property value. The id of the policy. Inherited from entity.
func (m *UnifiedRoleManagementPolicyAssignment) SetPolicyId(value *string)() {
    m.policyId = value
}
// SetRoleDefinitionId sets the roleDefinitionId property value. The identifier of the role definition object where the policy applies. If not specified, the policy applies to all roles. Supports $filter (eq).
func (m *UnifiedRoleManagementPolicyAssignment) SetRoleDefinitionId(value *string)() {
    m.roleDefinitionId = value
}
// SetScopeId sets the scopeId property value. The identifier of the scope where the policy is assigned.  Can be / for the tenant or a group ID. Required.
func (m *UnifiedRoleManagementPolicyAssignment) SetScopeId(value *string)() {
    m.scopeId = value
}
// SetScopeType sets the scopeType property value. The type of the scope where the policy is assigned. One of Directory, DirectoryRole. Required.
func (m *UnifiedRoleManagementPolicyAssignment) SetScopeType(value *string)() {
    m.scopeType = value
}
