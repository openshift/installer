package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppRoleAssignment provides operations to manage the collection of agreement entities.
type AppRoleAssignment struct {
    DirectoryObject
    // The identifier (id) for the app role which is assigned to the principal. This app role must be exposed in the appRoles property on the resource application's service principal (resourceId). If the resource application has not declared any app roles, a default app role ID of 00000000-0000-0000-0000-000000000000 can be specified to signal that the principal is assigned to the resource app without any specific app roles. Required on create.
    appRoleId *string
    // The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The display name of the user, group, or service principal that was granted the app role assignment. Read-only. Supports $filter (eq and startswith).
    principalDisplayName *string
    // The unique identifier (id) for the user, security group, or service principal being granted the app role. Security groups with dynamic memberships are supported. Required on create.
    principalId *string
    // The type of the assigned principal. This can either be User, Group, or ServicePrincipal. Read-only.
    principalType *string
    // The display name of the resource app's service principal to which the assignment is made.
    resourceDisplayName *string
    // The unique identifier (id) for the resource service principal for which the assignment is made. Required on create. Supports $filter (eq only).
    resourceId *string
}
// NewAppRoleAssignment instantiates a new appRoleAssignment and sets the default values.
func NewAppRoleAssignment()(*AppRoleAssignment) {
    m := &AppRoleAssignment{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.appRoleAssignment";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateAppRoleAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAppRoleAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppRoleAssignment(), nil
}
// GetAppRoleId gets the appRoleId property value. The identifier (id) for the app role which is assigned to the principal. This app role must be exposed in the appRoles property on the resource application's service principal (resourceId). If the resource application has not declared any app roles, a default app role ID of 00000000-0000-0000-0000-000000000000 can be specified to signal that the principal is assigned to the resource app without any specific app roles. Required on create.
func (m *AppRoleAssignment) GetAppRoleId()(*string) {
    return m.appRoleId
}
// GetCreatedDateTime gets the createdDateTime property value. The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *AppRoleAssignment) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AppRoleAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
    res["appRoleId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppRoleId)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["principalDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPrincipalDisplayName)
    res["principalId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPrincipalId)
    res["principalType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPrincipalType)
    res["resourceDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetResourceDisplayName)
    res["resourceId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetResourceId)
    return res
}
// GetPrincipalDisplayName gets the principalDisplayName property value. The display name of the user, group, or service principal that was granted the app role assignment. Read-only. Supports $filter (eq and startswith).
func (m *AppRoleAssignment) GetPrincipalDisplayName()(*string) {
    return m.principalDisplayName
}
// GetPrincipalId gets the principalId property value. The unique identifier (id) for the user, security group, or service principal being granted the app role. Security groups with dynamic memberships are supported. Required on create.
func (m *AppRoleAssignment) GetPrincipalId()(*string) {
    return m.principalId
}
// GetPrincipalType gets the principalType property value. The type of the assigned principal. This can either be User, Group, or ServicePrincipal. Read-only.
func (m *AppRoleAssignment) GetPrincipalType()(*string) {
    return m.principalType
}
// GetResourceDisplayName gets the resourceDisplayName property value. The display name of the resource app's service principal to which the assignment is made.
func (m *AppRoleAssignment) GetResourceDisplayName()(*string) {
    return m.resourceDisplayName
}
// GetResourceId gets the resourceId property value. The unique identifier (id) for the resource service principal for which the assignment is made. Required on create. Supports $filter (eq only).
func (m *AppRoleAssignment) GetResourceId()(*string) {
    return m.resourceId
}
// Serialize serializes information the current object
func (m *AppRoleAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appRoleId", m.GetAppRoleId())
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
        err = writer.WriteStringValue("principalDisplayName", m.GetPrincipalDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("principalId", m.GetPrincipalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("principalType", m.GetPrincipalType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceDisplayName", m.GetResourceDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceId", m.GetResourceId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppRoleId sets the appRoleId property value. The identifier (id) for the app role which is assigned to the principal. This app role must be exposed in the appRoles property on the resource application's service principal (resourceId). If the resource application has not declared any app roles, a default app role ID of 00000000-0000-0000-0000-000000000000 can be specified to signal that the principal is assigned to the resource app without any specific app roles. Required on create.
func (m *AppRoleAssignment) SetAppRoleId(value *string)() {
    m.appRoleId = value
}
// SetCreatedDateTime sets the createdDateTime property value. The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *AppRoleAssignment) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetPrincipalDisplayName sets the principalDisplayName property value. The display name of the user, group, or service principal that was granted the app role assignment. Read-only. Supports $filter (eq and startswith).
func (m *AppRoleAssignment) SetPrincipalDisplayName(value *string)() {
    m.principalDisplayName = value
}
// SetPrincipalId sets the principalId property value. The unique identifier (id) for the user, security group, or service principal being granted the app role. Security groups with dynamic memberships are supported. Required on create.
func (m *AppRoleAssignment) SetPrincipalId(value *string)() {
    m.principalId = value
}
// SetPrincipalType sets the principalType property value. The type of the assigned principal. This can either be User, Group, or ServicePrincipal. Read-only.
func (m *AppRoleAssignment) SetPrincipalType(value *string)() {
    m.principalType = value
}
// SetResourceDisplayName sets the resourceDisplayName property value. The display name of the resource app's service principal to which the assignment is made.
func (m *AppRoleAssignment) SetResourceDisplayName(value *string)() {
    m.resourceDisplayName = value
}
// SetResourceId sets the resourceId property value. The unique identifier (id) for the resource service principal for which the assignment is made. Required on create. Supports $filter (eq only).
func (m *AppRoleAssignment) SetResourceId(value *string)() {
    m.resourceId = value
}
