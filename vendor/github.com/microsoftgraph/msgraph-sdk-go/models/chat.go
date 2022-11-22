package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Chat provides operations to manage the collection of chat entities.
type Chat struct {
    Entity
    // The chatType property
    chatType *ChatType
    // Date and time at which the chat was created. Read-only.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // A collection of all the apps in the chat. Nullable.
    installedApps []TeamsAppInstallationable
    // Preview of the last message sent in the chat. Null if no messages have been sent in the chat. Currently, only the list chats operation supports this property.
    lastMessagePreview ChatMessageInfoable
    // Date and time at which the chat was renamed or list of members were last changed. Read-only.
    lastUpdatedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // A collection of all the members in the chat. Nullable.
    members []ConversationMemberable
    // A collection of all the messages in the chat. Nullable.
    messages []ChatMessageable
    // Represents details about an online meeting. If the chat isn't associated with an online meeting, the property is empty. Read-only.
    onlineMeetingInfo TeamworkOnlineMeetingInfoable
    // A collection of all the pinned messages in the chat. Nullable.
    pinnedMessages []PinnedChatMessageInfoable
    // A collection of all the tabs in the chat. Nullable.
    tabs []TeamsTabable
    // The identifier of the tenant in which the chat was created. Read-only.
    tenantId *string
    // (Optional) Subject or topic for the chat. Only available for group chats.
    topic *string
    // Represents caller-specific information about the chat, such as last message read date and time. This property is populated only when the request is made in a delegated context.
    viewpoint ChatViewpointable
    // The URL for the chat in Microsoft Teams. The URL should be treated as an opaque blob, and not parsed. Read-only.
    webUrl *string
}
// NewChat instantiates a new chat and sets the default values.
func NewChat()(*Chat) {
    m := &Chat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateChatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateChatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChat(), nil
}
// GetChatType gets the chatType property value. The chatType property
func (m *Chat) GetChatType()(*ChatType) {
    return m.chatType
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time at which the chat was created. Read-only.
func (m *Chat) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Chat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["chatType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseChatType , m.SetChatType)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["installedApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamsAppInstallationFromDiscriminatorValue , m.SetInstalledApps)
    res["lastMessagePreview"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChatMessageInfoFromDiscriminatorValue , m.SetLastMessagePreview)
    res["lastUpdatedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastUpdatedDateTime)
    res["members"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateConversationMemberFromDiscriminatorValue , m.SetMembers)
    res["messages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChatMessageFromDiscriminatorValue , m.SetMessages)
    res["onlineMeetingInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamworkOnlineMeetingInfoFromDiscriminatorValue , m.SetOnlineMeetingInfo)
    res["pinnedMessages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePinnedChatMessageInfoFromDiscriminatorValue , m.SetPinnedMessages)
    res["tabs"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamsTabFromDiscriminatorValue , m.SetTabs)
    res["tenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTenantId)
    res["topic"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTopic)
    res["viewpoint"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChatViewpointFromDiscriminatorValue , m.SetViewpoint)
    res["webUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetWebUrl)
    return res
}
// GetInstalledApps gets the installedApps property value. A collection of all the apps in the chat. Nullable.
func (m *Chat) GetInstalledApps()([]TeamsAppInstallationable) {
    return m.installedApps
}
// GetLastMessagePreview gets the lastMessagePreview property value. Preview of the last message sent in the chat. Null if no messages have been sent in the chat. Currently, only the list chats operation supports this property.
func (m *Chat) GetLastMessagePreview()(ChatMessageInfoable) {
    return m.lastMessagePreview
}
// GetLastUpdatedDateTime gets the lastUpdatedDateTime property value. Date and time at which the chat was renamed or list of members were last changed. Read-only.
func (m *Chat) GetLastUpdatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastUpdatedDateTime
}
// GetMembers gets the members property value. A collection of all the members in the chat. Nullable.
func (m *Chat) GetMembers()([]ConversationMemberable) {
    return m.members
}
// GetMessages gets the messages property value. A collection of all the messages in the chat. Nullable.
func (m *Chat) GetMessages()([]ChatMessageable) {
    return m.messages
}
// GetOnlineMeetingInfo gets the onlineMeetingInfo property value. Represents details about an online meeting. If the chat isn't associated with an online meeting, the property is empty. Read-only.
func (m *Chat) GetOnlineMeetingInfo()(TeamworkOnlineMeetingInfoable) {
    return m.onlineMeetingInfo
}
// GetPinnedMessages gets the pinnedMessages property value. A collection of all the pinned messages in the chat. Nullable.
func (m *Chat) GetPinnedMessages()([]PinnedChatMessageInfoable) {
    return m.pinnedMessages
}
// GetTabs gets the tabs property value. A collection of all the tabs in the chat. Nullable.
func (m *Chat) GetTabs()([]TeamsTabable) {
    return m.tabs
}
// GetTenantId gets the tenantId property value. The identifier of the tenant in which the chat was created. Read-only.
func (m *Chat) GetTenantId()(*string) {
    return m.tenantId
}
// GetTopic gets the topic property value. (Optional) Subject or topic for the chat. Only available for group chats.
func (m *Chat) GetTopic()(*string) {
    return m.topic
}
// GetViewpoint gets the viewpoint property value. Represents caller-specific information about the chat, such as last message read date and time. This property is populated only when the request is made in a delegated context.
func (m *Chat) GetViewpoint()(ChatViewpointable) {
    return m.viewpoint
}
// GetWebUrl gets the webUrl property value. The URL for the chat in Microsoft Teams. The URL should be treated as an opaque blob, and not parsed. Read-only.
func (m *Chat) GetWebUrl()(*string) {
    return m.webUrl
}
// Serialize serializes information the current object
func (m *Chat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChatType() != nil {
        cast := (*m.GetChatType()).String()
        err = writer.WriteStringValue("chatType", &cast)
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
    if m.GetInstalledApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetInstalledApps())
        err = writer.WriteCollectionOfObjectValues("installedApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastMessagePreview", m.GetLastMessagePreview())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastUpdatedDateTime", m.GetLastUpdatedDateTime())
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
    if m.GetMessages() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMessages())
        err = writer.WriteCollectionOfObjectValues("messages", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("onlineMeetingInfo", m.GetOnlineMeetingInfo())
        if err != nil {
            return err
        }
    }
    if m.GetPinnedMessages() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPinnedMessages())
        err = writer.WriteCollectionOfObjectValues("pinnedMessages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTabs() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTabs())
        err = writer.WriteCollectionOfObjectValues("tabs", cast)
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
        err = writer.WriteStringValue("topic", m.GetTopic())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("viewpoint", m.GetViewpoint())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webUrl", m.GetWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChatType sets the chatType property value. The chatType property
func (m *Chat) SetChatType(value *ChatType)() {
    m.chatType = value
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time at which the chat was created. Read-only.
func (m *Chat) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetInstalledApps sets the installedApps property value. A collection of all the apps in the chat. Nullable.
func (m *Chat) SetInstalledApps(value []TeamsAppInstallationable)() {
    m.installedApps = value
}
// SetLastMessagePreview sets the lastMessagePreview property value. Preview of the last message sent in the chat. Null if no messages have been sent in the chat. Currently, only the list chats operation supports this property.
func (m *Chat) SetLastMessagePreview(value ChatMessageInfoable)() {
    m.lastMessagePreview = value
}
// SetLastUpdatedDateTime sets the lastUpdatedDateTime property value. Date and time at which the chat was renamed or list of members were last changed. Read-only.
func (m *Chat) SetLastUpdatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastUpdatedDateTime = value
}
// SetMembers sets the members property value. A collection of all the members in the chat. Nullable.
func (m *Chat) SetMembers(value []ConversationMemberable)() {
    m.members = value
}
// SetMessages sets the messages property value. A collection of all the messages in the chat. Nullable.
func (m *Chat) SetMessages(value []ChatMessageable)() {
    m.messages = value
}
// SetOnlineMeetingInfo sets the onlineMeetingInfo property value. Represents details about an online meeting. If the chat isn't associated with an online meeting, the property is empty. Read-only.
func (m *Chat) SetOnlineMeetingInfo(value TeamworkOnlineMeetingInfoable)() {
    m.onlineMeetingInfo = value
}
// SetPinnedMessages sets the pinnedMessages property value. A collection of all the pinned messages in the chat. Nullable.
func (m *Chat) SetPinnedMessages(value []PinnedChatMessageInfoable)() {
    m.pinnedMessages = value
}
// SetTabs sets the tabs property value. A collection of all the tabs in the chat. Nullable.
func (m *Chat) SetTabs(value []TeamsTabable)() {
    m.tabs = value
}
// SetTenantId sets the tenantId property value. The identifier of the tenant in which the chat was created. Read-only.
func (m *Chat) SetTenantId(value *string)() {
    m.tenantId = value
}
// SetTopic sets the topic property value. (Optional) Subject or topic for the chat. Only available for group chats.
func (m *Chat) SetTopic(value *string)() {
    m.topic = value
}
// SetViewpoint sets the viewpoint property value. Represents caller-specific information about the chat, such as last message read date and time. This property is populated only when the request is made in a delegated context.
func (m *Chat) SetViewpoint(value ChatViewpointable)() {
    m.viewpoint = value
}
// SetWebUrl sets the webUrl property value. The URL for the chat in Microsoft Teams. The URL should be treated as an opaque blob, and not parsed. Read-only.
func (m *Chat) SetWebUrl(value *string)() {
    m.webUrl = value
}
