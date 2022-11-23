package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ConversationThread provides operations to manage the collection of agreement entities.
type ConversationThread struct {
    Entity
    // The Cc: recipients for the thread. Returned only on $select.
    ccRecipients []Recipientable
    // Indicates whether any of the posts within this thread has at least one attachment. Returned by default.
    hasAttachments *bool
    // Indicates if the thread is locked. Returned by default.
    isLocked *bool
    // The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.Returned by default.
    lastDeliveredDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The posts property
    posts []Postable
    // A short summary from the body of the latest post in this conversation. Returned by default.
    preview *string
    // The topic of the conversation. This property can be set when the conversation is created, but it cannot be updated. Returned by default.
    topic *string
    // The To: recipients for the thread. Returned only on $select.
    toRecipients []Recipientable
    // All the users that sent a message to this thread. Returned by default.
    uniqueSenders []string
}
// NewConversationThread instantiates a new conversationThread and sets the default values.
func NewConversationThread()(*ConversationThread) {
    m := &ConversationThread{
        Entity: *NewEntity(),
    }
    return m
}
// CreateConversationThreadFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateConversationThreadFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConversationThread(), nil
}
// GetCcRecipients gets the ccRecipients property value. The Cc: recipients for the thread. Returned only on $select.
func (m *ConversationThread) GetCcRecipients()([]Recipientable) {
    return m.ccRecipients
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ConversationThread) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["ccRecipients"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRecipientFromDiscriminatorValue , m.SetCcRecipients)
    res["hasAttachments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHasAttachments)
    res["isLocked"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsLocked)
    res["lastDeliveredDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastDeliveredDateTime)
    res["posts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePostFromDiscriminatorValue , m.SetPosts)
    res["preview"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPreview)
    res["topic"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTopic)
    res["toRecipients"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRecipientFromDiscriminatorValue , m.SetToRecipients)
    res["uniqueSenders"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetUniqueSenders)
    return res
}
// GetHasAttachments gets the hasAttachments property value. Indicates whether any of the posts within this thread has at least one attachment. Returned by default.
func (m *ConversationThread) GetHasAttachments()(*bool) {
    return m.hasAttachments
}
// GetIsLocked gets the isLocked property value. Indicates if the thread is locked. Returned by default.
func (m *ConversationThread) GetIsLocked()(*bool) {
    return m.isLocked
}
// GetLastDeliveredDateTime gets the lastDeliveredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.Returned by default.
func (m *ConversationThread) GetLastDeliveredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastDeliveredDateTime
}
// GetPosts gets the posts property value. The posts property
func (m *ConversationThread) GetPosts()([]Postable) {
    return m.posts
}
// GetPreview gets the preview property value. A short summary from the body of the latest post in this conversation. Returned by default.
func (m *ConversationThread) GetPreview()(*string) {
    return m.preview
}
// GetTopic gets the topic property value. The topic of the conversation. This property can be set when the conversation is created, but it cannot be updated. Returned by default.
func (m *ConversationThread) GetTopic()(*string) {
    return m.topic
}
// GetToRecipients gets the toRecipients property value. The To: recipients for the thread. Returned only on $select.
func (m *ConversationThread) GetToRecipients()([]Recipientable) {
    return m.toRecipients
}
// GetUniqueSenders gets the uniqueSenders property value. All the users that sent a message to this thread. Returned by default.
func (m *ConversationThread) GetUniqueSenders()([]string) {
    return m.uniqueSenders
}
// Serialize serializes information the current object
func (m *ConversationThread) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCcRecipients() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCcRecipients())
        err = writer.WriteCollectionOfObjectValues("ccRecipients", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasAttachments", m.GetHasAttachments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isLocked", m.GetIsLocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastDeliveredDateTime", m.GetLastDeliveredDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetPosts() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPosts())
        err = writer.WriteCollectionOfObjectValues("posts", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("preview", m.GetPreview())
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
    if m.GetToRecipients() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetToRecipients())
        err = writer.WriteCollectionOfObjectValues("toRecipients", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUniqueSenders() != nil {
        err = writer.WriteCollectionOfStringValues("uniqueSenders", m.GetUniqueSenders())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCcRecipients sets the ccRecipients property value. The Cc: recipients for the thread. Returned only on $select.
func (m *ConversationThread) SetCcRecipients(value []Recipientable)() {
    m.ccRecipients = value
}
// SetHasAttachments sets the hasAttachments property value. Indicates whether any of the posts within this thread has at least one attachment. Returned by default.
func (m *ConversationThread) SetHasAttachments(value *bool)() {
    m.hasAttachments = value
}
// SetIsLocked sets the isLocked property value. Indicates if the thread is locked. Returned by default.
func (m *ConversationThread) SetIsLocked(value *bool)() {
    m.isLocked = value
}
// SetLastDeliveredDateTime sets the lastDeliveredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.Returned by default.
func (m *ConversationThread) SetLastDeliveredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastDeliveredDateTime = value
}
// SetPosts sets the posts property value. The posts property
func (m *ConversationThread) SetPosts(value []Postable)() {
    m.posts = value
}
// SetPreview sets the preview property value. A short summary from the body of the latest post in this conversation. Returned by default.
func (m *ConversationThread) SetPreview(value *string)() {
    m.preview = value
}
// SetTopic sets the topic property value. The topic of the conversation. This property can be set when the conversation is created, but it cannot be updated. Returned by default.
func (m *ConversationThread) SetTopic(value *string)() {
    m.topic = value
}
// SetToRecipients sets the toRecipients property value. The To: recipients for the thread. Returned only on $select.
func (m *ConversationThread) SetToRecipients(value []Recipientable)() {
    m.toRecipients = value
}
// SetUniqueSenders sets the uniqueSenders property value. All the users that sent a message to this thread. Returned by default.
func (m *ConversationThread) SetUniqueSenders(value []string)() {
    m.uniqueSenders = value
}
