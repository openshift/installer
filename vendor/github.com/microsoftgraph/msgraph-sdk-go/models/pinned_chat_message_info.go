package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PinnedChatMessageInfo provides operations to manage the collection of chat entities.
type PinnedChatMessageInfo struct {
    Entity
    // Represents details about the chat message that is pinned.
    message ChatMessageable
}
// NewPinnedChatMessageInfo instantiates a new pinnedChatMessageInfo and sets the default values.
func NewPinnedChatMessageInfo()(*PinnedChatMessageInfo) {
    m := &PinnedChatMessageInfo{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePinnedChatMessageInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePinnedChatMessageInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPinnedChatMessageInfo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PinnedChatMessageInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["message"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChatMessageFromDiscriminatorValue , m.SetMessage)
    return res
}
// GetMessage gets the message property value. Represents details about the chat message that is pinned.
func (m *PinnedChatMessageInfo) GetMessage()(ChatMessageable) {
    return m.message
}
// Serialize serializes information the current object
func (m *PinnedChatMessageInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("message", m.GetMessage())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMessage sets the message property value. Represents details about the chat message that is pinned.
func (m *PinnedChatMessageInfo) SetMessage(value ChatMessageable)() {
    m.message = value
}
