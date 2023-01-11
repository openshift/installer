package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SkypeForBusinessUserConversationMember 
type SkypeForBusinessUserConversationMember struct {
    ConversationMember
    // The tenantId property
    tenantId *string
    // The userId property
    userId *string
}
// NewSkypeForBusinessUserConversationMember instantiates a new SkypeForBusinessUserConversationMember and sets the default values.
func NewSkypeForBusinessUserConversationMember()(*SkypeForBusinessUserConversationMember) {
    m := &SkypeForBusinessUserConversationMember{
        ConversationMember: *NewConversationMember(),
    }
    odataTypeValue := "#microsoft.graph.skypeForBusinessUserConversationMember";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateSkypeForBusinessUserConversationMemberFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSkypeForBusinessUserConversationMemberFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSkypeForBusinessUserConversationMember(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SkypeForBusinessUserConversationMember) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConversationMember.GetFieldDeserializers()
    res["tenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTenantId)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    return res
}
// GetTenantId gets the tenantId property value. The tenantId property
func (m *SkypeForBusinessUserConversationMember) GetTenantId()(*string) {
    return m.tenantId
}
// GetUserId gets the userId property value. The userId property
func (m *SkypeForBusinessUserConversationMember) GetUserId()(*string) {
    return m.userId
}
// Serialize serializes information the current object
func (m *SkypeForBusinessUserConversationMember) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConversationMember.Serialize(writer)
    if err != nil {
        return err
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
// SetTenantId sets the tenantId property value. The tenantId property
func (m *SkypeForBusinessUserConversationMember) SetTenantId(value *string)() {
    m.tenantId = value
}
// SetUserId sets the userId property value. The userId property
func (m *SkypeForBusinessUserConversationMember) SetUserId(value *string)() {
    m.userId = value
}
