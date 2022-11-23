package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TeamworkTagMember provides operations to manage the collection of agreement entities.
type TeamworkTagMember struct {
    Entity
    // The member's display name.
    displayName *string
    // The ID of the tenant that the tag member is a part of.
    tenantId *string
    // The user ID of the member.
    userId *string
}
// NewTeamworkTagMember instantiates a new teamworkTagMember and sets the default values.
func NewTeamworkTagMember()(*TeamworkTagMember) {
    m := &TeamworkTagMember{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamworkTagMemberFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamworkTagMemberFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamworkTagMember(), nil
}
// GetDisplayName gets the displayName property value. The member's display name.
func (m *TeamworkTagMember) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TeamworkTagMember) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["tenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTenantId)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    return res
}
// GetTenantId gets the tenantId property value. The ID of the tenant that the tag member is a part of.
func (m *TeamworkTagMember) GetTenantId()(*string) {
    return m.tenantId
}
// GetUserId gets the userId property value. The user ID of the member.
func (m *TeamworkTagMember) GetUserId()(*string) {
    return m.userId
}
// Serialize serializes information the current object
func (m *TeamworkTagMember) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The member's display name.
func (m *TeamworkTagMember) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetTenantId sets the tenantId property value. The ID of the tenant that the tag member is a part of.
func (m *TeamworkTagMember) SetTenantId(value *string)() {
    m.tenantId = value
}
// SetUserId sets the userId property value. The user ID of the member.
func (m *TeamworkTagMember) SetUserId(value *string)() {
    m.userId = value
}
