package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TeamInfo provides operations to manage the collection of agreement entities.
type TeamInfo struct {
    Entity
    // The name of the team.
    displayName *string
    // The team property
    team Teamable
    // The ID of the Azure Active Directory tenant.
    tenantId *string
}
// NewTeamInfo instantiates a new teamInfo and sets the default values.
func NewTeamInfo()(*TeamInfo) {
    m := &TeamInfo{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.associatedTeamInfo":
                        return NewAssociatedTeamInfo(), nil
                    case "#microsoft.graph.sharedWithChannelTeamInfo":
                        return NewSharedWithChannelTeamInfo(), nil
                }
            }
        }
    }
    return NewTeamInfo(), nil
}
// GetDisplayName gets the displayName property value. The name of the team.
func (m *TeamInfo) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TeamInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["team"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamFromDiscriminatorValue , m.SetTeam)
    res["tenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTenantId)
    return res
}
// GetTeam gets the team property value. The team property
func (m *TeamInfo) GetTeam()(Teamable) {
    return m.team
}
// GetTenantId gets the tenantId property value. The ID of the Azure Active Directory tenant.
func (m *TeamInfo) GetTenantId()(*string) {
    return m.tenantId
}
// Serialize serializes information the current object
func (m *TeamInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteObjectValue("team", m.GetTeam())
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
    return nil
}
// SetDisplayName sets the displayName property value. The name of the team.
func (m *TeamInfo) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetTeam sets the team property value. The team property
func (m *TeamInfo) SetTeam(value Teamable)() {
    m.team = value
}
// SetTenantId sets the tenantId property value. The ID of the Azure Active Directory tenant.
func (m *TeamInfo) SetTenantId(value *string)() {
    m.tenantId = value
}
