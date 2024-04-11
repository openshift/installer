package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkbookComment 
type WorkbookComment struct {
    Entity
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
    val, err := m.GetBackingStore().Get("content")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContentType gets the contentType property value. Indicates the type for the comment.
func (m *WorkbookComment) GetContentType()(*string) {
    val, err := m.GetBackingStore().Get("contentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkbookComment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["content"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContent(val)
        }
        return nil
    }
    res["contentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentType(val)
        }
        return nil
    }
    res["replies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookCommentReplyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookCommentReplyable, len(val))
            for i, v := range val {
                res[i] = v.(WorkbookCommentReplyable)
            }
            m.SetReplies(res)
        }
        return nil
    }
    return res
}
// GetReplies gets the replies property value. The replies property
func (m *WorkbookComment) GetReplies()([]WorkbookCommentReplyable) {
    val, err := m.GetBackingStore().Get("replies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookCommentReplyable)
    }
    return nil
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
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReplies()))
        for i, v := range m.GetReplies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("replies", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContent sets the content property value. The content of comment.
func (m *WorkbookComment) SetContent(value *string)() {
    err := m.GetBackingStore().Set("content", value)
    if err != nil {
        panic(err)
    }
}
// SetContentType sets the contentType property value. Indicates the type for the comment.
func (m *WorkbookComment) SetContentType(value *string)() {
    err := m.GetBackingStore().Set("contentType", value)
    if err != nil {
        panic(err)
    }
}
// SetReplies sets the replies property value. The replies property
func (m *WorkbookComment) SetReplies(value []WorkbookCommentReplyable)() {
    err := m.GetBackingStore().Set("replies", value)
    if err != nil {
        panic(err)
    }
}
// WorkbookCommentable 
type WorkbookCommentable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContent()(*string)
    GetContentType()(*string)
    GetReplies()([]WorkbookCommentReplyable)
    SetContent(value *string)()
    SetContentType(value *string)()
    SetReplies(value []WorkbookCommentReplyable)()
}
