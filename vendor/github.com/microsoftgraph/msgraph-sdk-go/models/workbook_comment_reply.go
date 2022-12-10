package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookCommentReply provides operations to manage the collection of agreement entities.
type WorkbookCommentReply struct {
    Entity
    // The content of a comment reply.
    content *string
    // Indicates the type for the comment reply.
    contentType *string
}
// NewWorkbookCommentReply instantiates a new workbookCommentReply and sets the default values.
func NewWorkbookCommentReply()(*WorkbookCommentReply) {
    m := &WorkbookCommentReply{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookCommentReplyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkbookCommentReplyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookCommentReply(), nil
}
// GetContent gets the content property value. The content of a comment reply.
func (m *WorkbookCommentReply) GetContent()(*string) {
    return m.content
}
// GetContentType gets the contentType property value. Indicates the type for the comment reply.
func (m *WorkbookCommentReply) GetContentType()(*string) {
    return m.contentType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookCommentReply) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["content"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetContent)
    res["contentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetContentType)
    return res
}
// Serialize serializes information the current object
func (m *WorkbookCommentReply) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    return nil
}
// SetContent sets the content property value. The content of a comment reply.
func (m *WorkbookCommentReply) SetContent(value *string)() {
    m.content = value
}
// SetContentType sets the contentType property value. Indicates the type for the comment reply.
func (m *WorkbookCommentReply) SetContentType(value *string)() {
    m.contentType = value
}
