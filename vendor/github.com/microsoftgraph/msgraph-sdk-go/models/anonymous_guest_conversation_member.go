package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AnonymousGuestConversationMember 
type AnonymousGuestConversationMember struct {
    ConversationMember
    // The anonymousGuestId property
    anonymousGuestId *string
}
// NewAnonymousGuestConversationMember instantiates a new AnonymousGuestConversationMember and sets the default values.
func NewAnonymousGuestConversationMember()(*AnonymousGuestConversationMember) {
    m := &AnonymousGuestConversationMember{
        ConversationMember: *NewConversationMember(),
    }
    odataTypeValue := "#microsoft.graph.anonymousGuestConversationMember";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateAnonymousGuestConversationMemberFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAnonymousGuestConversationMemberFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAnonymousGuestConversationMember(), nil
}
// GetAnonymousGuestId gets the anonymousGuestId property value. The anonymousGuestId property
func (m *AnonymousGuestConversationMember) GetAnonymousGuestId()(*string) {
    return m.anonymousGuestId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AnonymousGuestConversationMember) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConversationMember.GetFieldDeserializers()
    res["anonymousGuestId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAnonymousGuestId)
    return res
}
// Serialize serializes information the current object
func (m *AnonymousGuestConversationMember) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConversationMember.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("anonymousGuestId", m.GetAnonymousGuestId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAnonymousGuestId sets the anonymousGuestId property value. The anonymousGuestId property
func (m *AnonymousGuestConversationMember) SetAnonymousGuestId(value *string)() {
    m.anonymousGuestId = value
}
