package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MembersJoinedEventMessageDetail 
type MembersJoinedEventMessageDetail struct {
    EventMessageDetail
    // Initiator of the event.
    initiator IdentitySetable
    // List of members who joined the chat.
    members []TeamworkUserIdentityable
}
// NewMembersJoinedEventMessageDetail instantiates a new MembersJoinedEventMessageDetail and sets the default values.
func NewMembersJoinedEventMessageDetail()(*MembersJoinedEventMessageDetail) {
    m := &MembersJoinedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.membersJoinedEventMessageDetail";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateMembersJoinedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMembersJoinedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMembersJoinedEventMessageDetail(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MembersJoinedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["initiator"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetInitiator)
    res["members"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamworkUserIdentityFromDiscriminatorValue , m.SetMembers)
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
func (m *MembersJoinedEventMessageDetail) GetInitiator()(IdentitySetable) {
    return m.initiator
}
// GetMembers gets the members property value. List of members who joined the chat.
func (m *MembersJoinedEventMessageDetail) GetMembers()([]TeamworkUserIdentityable) {
    return m.members
}
// Serialize serializes information the current object
func (m *MembersJoinedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
        if err != nil {
            return err
        }
    }
    if m.GetMembers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMembers())
        err = writer.WriteCollectionOfObjectValues("members", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *MembersJoinedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    m.initiator = value
}
// SetMembers sets the members property value. List of members who joined the chat.
func (m *MembersJoinedEventMessageDetail) SetMembers(value []TeamworkUserIdentityable)() {
    m.members = value
}
