package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookComment provides operations to manage the collection of agreement entities.
type WorkbookComment struct {
    Entity
    // The content of comment.
    content *string
    // Indicates the type for the comment.
    contentType *string
    // The replies property
    replies []WorkbookCommentReplyable
}
// NewWorkbookComment instantiates a new workbookComment and sets the default values.
func NewWorkbookComment()(*WorkbookComment) {
    m := &WorkbookComment{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookCommentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookCommentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookComment(), nil
}
// GetContent gets the content property value. The content of comment.
func (m *WorkbookComment) GetContent()(*string) {
    return m.content
}
// GetContentType gets the contentType property value. Indicates the type for the comment.
func (m *WorkbookComment) GetContentType()(*string) {
    return m.contentType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookComment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["content"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetContent)
    res["contentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetContentType)
    res["replies"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWorkbookCommentReplyFromDiscriminatorValue , m.SetReplies)
    return res
}
// GetReplies gets the replies property value. The replies property
func (m *WorkbookComment) GetReplies()([]WorkbookCommentReplyable) {
    return m.replies
}
// Serialize serializes information the current object
func (m *WorkbookComment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("content", m.GetContent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("contentType", m.GetContentType())
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
    return nil
}
// SetContent sets the content property value. The content of comment.
func (m *WorkbookComment) SetContent(value *string)() {
    m.content = value
}
// SetContentType sets the contentType property value. Indicates the type for the comment.
func (m *WorkbookComment) SetContentType(value *string)() {
    m.contentType = value
}
// SetReplies sets the replies property value. The replies property
func (m *WorkbookComment) SetReplies(value []WorkbookCommentReplyable)() {
    m.replies = value
}
