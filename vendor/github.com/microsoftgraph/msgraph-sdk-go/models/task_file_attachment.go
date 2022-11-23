package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TaskFileAttachment 
type TaskFileAttachment struct {
    AttachmentBase
    // The contentBytes property
    contentBytes []byte
}
// NewTaskFileAttachment instantiates a new TaskFileAttachment and sets the default values.
func NewTaskFileAttachment()(*TaskFileAttachment) {
    m := &TaskFileAttachment{
        AttachmentBase: *NewAttachmentBase(),
    }
    odataTypeValue := "#microsoft.graph.taskFileAttachment";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateTaskFileAttachmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTaskFileAttachmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTaskFileAttachment(), nil
}
// GetContentBytes gets the contentBytes property value. The contentBytes property
func (m *TaskFileAttachment) GetContentBytes()([]byte) {
    return m.contentBytes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TaskFileAttachment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AttachmentBase.GetFieldDeserializers()
    res["contentBytes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteArrayValue(m.SetContentBytes)
    return res
}
// Serialize serializes information the current object
func (m *TaskFileAttachment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AttachmentBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteByteArrayValue("contentBytes", m.GetContentBytes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContentBytes sets the contentBytes property value. The contentBytes property
func (m *TaskFileAttachment) SetContentBytes(value []byte)() {
    m.contentBytes = value
}
