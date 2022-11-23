package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SkypeUserConversationMember 
type SkypeUserConversationMember struct {
    ConversationMember
    // The skypeId property
    skypeId *string
}
// NewSkypeUserConversationMember instantiates a new SkypeUserConversationMember and sets the default values.
func NewSkypeUserConversationMember()(*SkypeUserConversationMember) {
    m := &SkypeUserConversationMember{
        ConversationMember: *NewConversationMember(),
    }
    odataTypeValue := "#microsoft.graph.skypeUserConversationMember";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateSkypeUserConversationMemberFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSkypeUserConversationMemberFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSkypeUserConversationMember(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SkypeUserConversationMember) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConversationMember.GetFieldDeserializers()
    res["skypeId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSkypeId)
    return res
}
// GetSkypeId gets the skypeId property value. The skypeId property
func (m *SkypeUserConversationMember) GetSkypeId()(*string) {
    return m.skypeId
}
// Serialize serializes information the current object
func (m *SkypeUserConversationMember) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConversationMember.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("skypeId", m.GetSkypeId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSkypeId sets the skypeId property value. The skypeId property
func (m *SkypeUserConversationMember) SetSkypeId(value *string)() {
    m.skypeId = value
}
