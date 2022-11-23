package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SharedWithChannelTeamInfo 
type SharedWithChannelTeamInfo struct {
    TeamInfo
    // A collection of team members who have access to the shared channel.
    allowedMembers []ConversationMemberable
    // Indicates whether the team is the host of the channel.
    isHostTeam *bool
}
// NewSharedWithChannelTeamInfo instantiates a new SharedWithChannelTeamInfo and sets the default values.
func NewSharedWithChannelTeamInfo()(*SharedWithChannelTeamInfo) {
    m := &SharedWithChannelTeamInfo{
        TeamInfo: *NewTeamInfo(),
    }
    return m
}
// CreateSharedWithChannelTeamInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSharedWithChannelTeamInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharedWithChannelTeamInfo(), nil
}
// GetAllowedMembers gets the allowedMembers property value. A collection of team members who have access to the shared channel.
func (m *SharedWithChannelTeamInfo) GetAllowedMembers()([]ConversationMemberable) {
    return m.allowedMembers
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SharedWithChannelTeamInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TeamInfo.GetFieldDeserializers()
    res["allowedMembers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateConversationMemberFromDiscriminatorValue , m.SetAllowedMembers)
    res["isHostTeam"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsHostTeam)
    return res
}
// GetIsHostTeam gets the isHostTeam property value. Indicates whether the team is the host of the channel.
func (m *SharedWithChannelTeamInfo) GetIsHostTeam()(*bool) {
    return m.isHostTeam
}
// Serialize serializes information the current object
func (m *SharedWithChannelTeamInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TeamInfo.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedMembers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAllowedMembers())
        err = writer.WriteCollectionOfObjectValues("allowedMembers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isHostTeam", m.GetIsHostTeam())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedMembers sets the allowedMembers property value. A collection of team members who have access to the shared channel.
func (m *SharedWithChannelTeamInfo) SetAllowedMembers(value []ConversationMemberable)() {
    m.allowedMembers = value
}
// SetIsHostTeam sets the isHostTeam property value. Indicates whether the team is the host of the channel.
func (m *SharedWithChannelTeamInfo) SetIsHostTeam(value *bool)() {
    m.isHostTeam = value
}
