package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ChatMessage provides operations to manage the collection of chat entities.
type ChatMessage struct {
    Entity
    // References to attached objects like files, tabs, meetings etc.
    attachments []ChatMessageAttachmentable
    // The body property
    body ItemBodyable
    // If the message was sent in a channel, represents identity of the channel.
    channelIdentity ChannelIdentityable
    // If the message was sent in a chat, represents the identity of the chat.
    chatId *string
    // Timestamp of when the chat message was created.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Read only. Timestamp at which the chat message was deleted, or null if not deleted.
    deletedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Read-only. Version number of the chat message.
    etag *string
    // Read-only. If present, represents details of an event that happened in a chat, a channel, or a team, for example, adding new members. For event messages, the messageType property will be set to systemEventMessage.
    eventDetail EventMessageDetailable
    // Details of the sender of the chat message. Can only be set during migration.
    from ChatMessageFromIdentitySetable
    // Content in a message hosted by Microsoft Teams - for example, images or code snippets.
    hostedContents []ChatMessageHostedContentable
    // The importance property
    importance *ChatMessageImportance
    // Read only. Timestamp when edits to the chat message were made. Triggers an 'Edited' flag in the Teams UI. If no edits are made the value is null.
    lastEditedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Read only. Timestamp when the chat message is created (initial setting) or modified, including when a reaction is added or removed.
    lastModifiedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Locale of the chat message set by the client. Always set to en-us.
    locale *string
    // List of entities mentioned in the chat message. Supported entities are: user, bot, team, and channel.
    mentions []ChatMessageMentionable
    // The messageType property
    messageType *ChatMessageType
    // Defines the properties of a policy violation set by a data loss prevention (DLP) application.
    policyViolation ChatMessagePolicyViolationable
    // Reactions for this chat message (for example, Like).
    reactions []ChatMessageReactionable
    // Replies for a specified message. Supports $expand for channel messages.
    replies []ChatMessageable
    // Read-only. ID of the parent chat message or root chat message of the thread. (Only applies to chat messages in channels, not chats.)
    replyToId *string
    // The subject of the chat message, in plaintext.
    subject *string
    // Summary text of the chat message that could be used for push notifications and summary views or fall back views. Only applies to channel chat messages, not chat messages in a chat.
    summary *string
    // Read-only. Link to the message in Microsoft Teams.
    webUrl *string
}
// NewChatMessage instantiates a new chatMessage and sets the default values.
func NewChatMessage()(*ChatMessage) {
    m := &ChatMessage{
        Entity: *NewEntity(),
    }
    return m
}
// CreateChatMessageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateChatMessageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChatMessage(), nil
}
// GetAttachments gets the attachments property value. References to attached objects like files, tabs, meetings etc.
func (m *ChatMessage) GetAttachments()([]ChatMessageAttachmentable) {
    return m.attachments
}
// GetBody gets the body property value. The body property
func (m *ChatMessage) GetBody()(ItemBodyable) {
    return m.body
}
// GetChannelIdentity gets the channelIdentity property value. If the message was sent in a channel, represents identity of the channel.
func (m *ChatMessage) GetChannelIdentity()(ChannelIdentityable) {
    return m.channelIdentity
}
// GetChatId gets the chatId property value. If the message was sent in a chat, represents the identity of the chat.
func (m *ChatMessage) GetChatId()(*string) {
    return m.chatId
}
// GetCreatedDateTime gets the createdDateTime property value. Timestamp of when the chat message was created.
func (m *ChatMessage) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDeletedDateTime gets the deletedDateTime property value. Read only. Timestamp at which the chat message was deleted, or null if not deleted.
func (m *ChatMessage) GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.deletedDateTime
}
// GetEtag gets the etag property value. Read-only. Version number of the chat message.
func (m *ChatMessage) GetEtag()(*string) {
    return m.etag
}
// GetEventDetail gets the eventDetail property value. Read-only. If present, represents details of an event that happened in a chat, a channel, or a team, for example, adding new members. For event messages, the messageType property will be set to systemEventMessage.
func (m *ChatMessage) GetEventDetail()(EventMessageDetailable) {
    return m.eventDetail
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ChatMessage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attachments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChatMessageAttachmentFromDiscriminatorValue , m.SetAttachments)
    res["body"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateItemBodyFromDiscriminatorValue , m.SetBody)
    res["channelIdentity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChannelIdentityFromDiscriminatorValue , m.SetChannelIdentity)
    res["chatId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetChatId)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["deletedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetDeletedDateTime)
    res["etag"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEtag)
    res["eventDetail"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEventMessageDetailFromDiscriminatorValue , m.SetEventDetail)
    res["from"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChatMessageFromIdentitySetFromDiscriminatorValue , m.SetFrom)
    res["hostedContents"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChatMessageHostedContentFromDiscriminatorValue , m.SetHostedContents)
    res["importance"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseChatMessageImportance , m.SetImportance)
    res["lastEditedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastEditedDateTime)
    res["lastModifiedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastModifiedDateTime)
    res["locale"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLocale)
    res["mentions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChatMessageMentionFromDiscriminatorValue , m.SetMentions)
    res["messageType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseChatMessageType , m.SetMessageType)
    res["policyViolation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChatMessagePolicyViolationFromDiscriminatorValue , m.SetPolicyViolation)
    res["reactions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChatMessageReactionFromDiscriminatorValue , m.SetReactions)
    res["replies"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChatMessageFromDiscriminatorValue , m.SetReplies)
    res["replyToId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetReplyToId)
    res["subject"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSubject)
    res["summary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSummary)
    res["webUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetWebUrl)
    return res
}
// GetFrom gets the from property value. Details of the sender of the chat message. Can only be set during migration.
func (m *ChatMessage) GetFrom()(ChatMessageFromIdentitySetable) {
    return m.from
}
// GetHostedContents gets the hostedContents property value. Content in a message hosted by Microsoft Teams - for example, images or code snippets.
func (m *ChatMessage) GetHostedContents()([]ChatMessageHostedContentable) {
    return m.hostedContents
}
// GetImportance gets the importance property value. The importance property
func (m *ChatMessage) GetImportance()(*ChatMessageImportance) {
    return m.importance
}
// GetLastEditedDateTime gets the lastEditedDateTime property value. Read only. Timestamp when edits to the chat message were made. Triggers an 'Edited' flag in the Teams UI. If no edits are made the value is null.
func (m *ChatMessage) GetLastEditedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastEditedDateTime
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Read only. Timestamp when the chat message is created (initial setting) or modified, including when a reaction is added or removed.
func (m *ChatMessage) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastModifiedDateTime
}
// GetLocale gets the locale property value. Locale of the chat message set by the client. Always set to en-us.
func (m *ChatMessage) GetLocale()(*string) {
    return m.locale
}
// GetMentions gets the mentions property value. List of entities mentioned in the chat message. Supported entities are: user, bot, team, and channel.
func (m *ChatMessage) GetMentions()([]ChatMessageMentionable) {
    return m.mentions
}
// GetMessageType gets the messageType property value. The messageType property
func (m *ChatMessage) GetMessageType()(*ChatMessageType) {
    return m.messageType
}
// GetPolicyViolation gets the policyViolation property value. Defines the properties of a policy violation set by a data loss prevention (DLP) application.
func (m *ChatMessage) GetPolicyViolation()(ChatMessagePolicyViolationable) {
    return m.policyViolation
}
// GetReactions gets the reactions property value. Reactions for this chat message (for example, Like).
func (m *ChatMessage) GetReactions()([]ChatMessageReactionable) {
    return m.reactions
}
// GetReplies gets the replies property value. Replies for a specified message. Supports $expand for channel messages.
func (m *ChatMessage) GetReplies()([]ChatMessageable) {
    return m.replies
}
// GetReplyToId gets the replyToId property value. Read-only. ID of the parent chat message or root chat message of the thread. (Only applies to chat messages in channels, not chats.)
func (m *ChatMessage) GetReplyToId()(*string) {
    return m.replyToId
}
// GetSubject gets the subject property value. The subject of the chat message, in plaintext.
func (m *ChatMessage) GetSubject()(*string) {
    return m.subject
}
// GetSummary gets the summary property value. Summary text of the chat message that could be used for push notifications and summary views or fall back views. Only applies to channel chat messages, not chat messages in a chat.
func (m *ChatMessage) GetSummary()(*string) {
    return m.summary
}
// GetWebUrl gets the webUrl property value. Read-only. Link to the message in Microsoft Teams.
func (m *ChatMessage) GetWebUrl()(*string) {
    return m.webUrl
}
// Serialize serializes information the current object
func (m *ChatMessage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAttachments() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAttachments())
        err = writer.WriteCollectionOfObjectValues("attachments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("body", m.GetBody())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("channelIdentity", m.GetChannelIdentity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("chatId", m.GetChatId())
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
        err = writer.WriteTimeValue("deletedDateTime", m.GetDeletedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("etag", m.GetEtag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("eventDetail", m.GetEventDetail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("from", m.GetFrom())
        if err != nil {
            return err
        }
    }
    if m.GetHostedContents() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetHostedContents())
        err = writer.WriteCollectionOfObjectValues("hostedContents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetImportance() != nil {
        cast := (*m.GetImportance()).String()
        err = writer.WriteStringValue("importance", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastEditedDateTime", m.GetLastEditedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("locale", m.GetLocale())
        if err != nil {
            return err
        }
    }
    if m.GetMentions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMentions())
        err = writer.WriteCollectionOfObjectValues("mentions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMessageType() != nil {
        cast := (*m.GetMessageType()).String()
        err = writer.WriteStringValue("messageType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("policyViolation", m.GetPolicyViolation())
        if err != nil {
            return err
        }
    }
    if m.GetReactions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetReactions())
        err = writer.WriteCollectionOfObjectValues("reactions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetReplies() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetReplies())
        err = writer.WriteCollectionOfObjectValues("replies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("replyToId", m.GetReplyToId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("subject", m.GetSubject())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("summary", m.GetSummary())
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
// SetAttachments sets the attachments property value. References to attached objects like files, tabs, meetings etc.
func (m *ChatMessage) SetAttachments(value []ChatMessageAttachmentable)() {
    m.attachments = value
}
// SetBody sets the body property value. The body property
func (m *ChatMessage) SetBody(value ItemBodyable)() {
    m.body = value
}
// SetChannelIdentity sets the channelIdentity property value. If the message was sent in a channel, represents identity of the channel.
func (m *ChatMessage) SetChannelIdentity(value ChannelIdentityable)() {
    m.channelIdentity = value
}
// SetChatId sets the chatId property value. If the message was sent in a chat, represents the identity of the chat.
func (m *ChatMessage) SetChatId(value *string)() {
    m.chatId = value
}
// SetCreatedDateTime sets the createdDateTime property value. Timestamp of when the chat message was created.
func (m *ChatMessage) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDeletedDateTime sets the deletedDateTime property value. Read only. Timestamp at which the chat message was deleted, or null if not deleted.
func (m *ChatMessage) SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.deletedDateTime = value
}
// SetEtag sets the etag property value. Read-only. Version number of the chat message.
func (m *ChatMessage) SetEtag(value *string)() {
    m.etag = value
}
// SetEventDetail sets the eventDetail property value. Read-only. If present, represents details of an event that happened in a chat, a channel, or a team, for example, adding new members. For event messages, the messageType property will be set to systemEventMessage.
func (m *ChatMessage) SetEventDetail(value EventMessageDetailable)() {
    m.eventDetail = value
}
// SetFrom sets the from property value. Details of the sender of the chat message. Can only be set during migration.
func (m *ChatMessage) SetFrom(value ChatMessageFromIdentitySetable)() {
    m.from = value
}
// SetHostedContents sets the hostedContents property value. Content in a message hosted by Microsoft Teams - for example, images or code snippets.
func (m *ChatMessage) SetHostedContents(value []ChatMessageHostedContentable)() {
    m.hostedContents = value
}
// SetImportance sets the importance property value. The importance property
func (m *ChatMessage) SetImportance(value *ChatMessageImportance)() {
    m.importance = value
}
// SetLastEditedDateTime sets the lastEditedDateTime property value. Read only. Timestamp when edits to the chat message were made. Triggers an 'Edited' flag in the Teams UI. If no edits are made the value is null.
func (m *ChatMessage) SetLastEditedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastEditedDateTime = value
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Read only. Timestamp when the chat message is created (initial setting) or modified, including when a reaction is added or removed.
func (m *ChatMessage) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastModifiedDateTime = value
}
// SetLocale sets the locale property value. Locale of the chat message set by the client. Always set to en-us.
func (m *ChatMessage) SetLocale(value *string)() {
    m.locale = value
}
// SetMentions sets the mentions property value. List of entities mentioned in the chat message. Supported entities are: user, bot, team, and channel.
func (m *ChatMessage) SetMentions(value []ChatMessageMentionable)() {
    m.mentions = value
}
// SetMessageType sets the messageType property value. The messageType property
func (m *ChatMessage) SetMessageType(value *ChatMessageType)() {
    m.messageType = value
}
// SetPolicyViolation sets the policyViolation property value. Defines the properties of a policy violation set by a data loss prevention (DLP) application.
func (m *ChatMessage) SetPolicyViolation(value ChatMessagePolicyViolationable)() {
    m.policyViolation = value
}
// SetReactions sets the reactions property value. Reactions for this chat message (for example, Like).
func (m *ChatMessage) SetReactions(value []ChatMessageReactionable)() {
    m.reactions = value
}
// SetReplies sets the replies property value. Replies for a specified message. Supports $expand for channel messages.
func (m *ChatMessage) SetReplies(value []ChatMessageable)() {
    m.replies = value
}
// SetReplyToId sets the replyToId property value. Read-only. ID of the parent chat message or root chat message of the thread. (Only applies to chat messages in channels, not chats.)
func (m *ChatMessage) SetReplyToId(value *string)() {
    m.replyToId = value
}
// SetSubject sets the subject property value. The subject of the chat message, in plaintext.
func (m *ChatMessage) SetSubject(value *string)() {
    m.subject = value
}
// SetSummary sets the summary property value. Summary text of the chat message that could be used for push notifications and summary views or fall back views. Only applies to channel chat messages, not chat messages in a chat.
func (m *ChatMessage) SetSummary(value *string)() {
    m.summary = value
}
// SetWebUrl sets the webUrl property value. Read-only. Link to the message in Microsoft Teams.
func (m *ChatMessage) SetWebUrl(value *string)() {
    m.webUrl = value
}
