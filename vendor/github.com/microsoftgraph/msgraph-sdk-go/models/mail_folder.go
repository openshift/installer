package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MailFolder provides operations to manage the collection of agreement entities.
type MailFolder struct {
    Entity
    // The number of immediate child mailFolders in the current mailFolder.
    childFolderCount *int32
    // The collection of child folders in the mailFolder.
    childFolders []MailFolderable
    // The mailFolder's display name.
    displayName *string
    // Indicates whether the mailFolder is hidden. This property can be set only when creating the folder. Find more information in Hidden mail folders.
    isHidden *bool
    // The collection of rules that apply to the user's Inbox folder.
    messageRules []MessageRuleable
    // The collection of messages in the mailFolder.
    messages []Messageable
    // The collection of multi-value extended properties defined for the mailFolder. Read-only. Nullable.
    multiValueExtendedProperties []MultiValueLegacyExtendedPropertyable
    // The unique identifier for the mailFolder's parent mailFolder.
    parentFolderId *string
    // The collection of single-value extended properties defined for the mailFolder. Read-only. Nullable.
    singleValueExtendedProperties []SingleValueLegacyExtendedPropertyable
    // The number of items in the mailFolder.
    totalItemCount *int32
    // The number of items in the mailFolder marked as unread.
    unreadItemCount *int32
}
// NewMailFolder instantiates a new mailFolder and sets the default values.
func NewMailFolder()(*MailFolder) {
    m := &MailFolder{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMailFolderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMailFolderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.mailSearchFolder":
                        return NewMailSearchFolder(), nil
                }
            }
        }
    }
    return NewMailFolder(), nil
}
// GetChildFolderCount gets the childFolderCount property value. The number of immediate child mailFolders in the current mailFolder.
func (m *MailFolder) GetChildFolderCount()(*int32) {
    return m.childFolderCount
}
// GetChildFolders gets the childFolders property value. The collection of child folders in the mailFolder.
func (m *MailFolder) GetChildFolders()([]MailFolderable) {
    return m.childFolders
}
// GetDisplayName gets the displayName property value. The mailFolder's display name.
func (m *MailFolder) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *MailFolder) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["childFolderCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetChildFolderCount)
    res["childFolders"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMailFolderFromDiscriminatorValue , m.SetChildFolders)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["isHidden"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsHidden)
    res["messageRules"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMessageRuleFromDiscriminatorValue , m.SetMessageRules)
    res["messages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMessageFromDiscriminatorValue , m.SetMessages)
    res["multiValueExtendedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMultiValueLegacyExtendedPropertyFromDiscriminatorValue , m.SetMultiValueExtendedProperties)
    res["parentFolderId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetParentFolderId)
    res["singleValueExtendedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue , m.SetSingleValueExtendedProperties)
    res["totalItemCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetTotalItemCount)
    res["unreadItemCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetUnreadItemCount)
    return res
}
// GetIsHidden gets the isHidden property value. Indicates whether the mailFolder is hidden. This property can be set only when creating the folder. Find more information in Hidden mail folders.
func (m *MailFolder) GetIsHidden()(*bool) {
    return m.isHidden
}
// GetMessageRules gets the messageRules property value. The collection of rules that apply to the user's Inbox folder.
func (m *MailFolder) GetMessageRules()([]MessageRuleable) {
    return m.messageRules
}
// GetMessages gets the messages property value. The collection of messages in the mailFolder.
func (m *MailFolder) GetMessages()([]Messageable) {
    return m.messages
}
// GetMultiValueExtendedProperties gets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the mailFolder. Read-only. Nullable.
func (m *MailFolder) GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable) {
    return m.multiValueExtendedProperties
}
// GetParentFolderId gets the parentFolderId property value. The unique identifier for the mailFolder's parent mailFolder.
func (m *MailFolder) GetParentFolderId()(*string) {
    return m.parentFolderId
}
// GetSingleValueExtendedProperties gets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the mailFolder. Read-only. Nullable.
func (m *MailFolder) GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable) {
    return m.singleValueExtendedProperties
}
// GetTotalItemCount gets the totalItemCount property value. The number of items in the mailFolder.
func (m *MailFolder) GetTotalItemCount()(*int32) {
    return m.totalItemCount
}
// GetUnreadItemCount gets the unreadItemCount property value. The number of items in the mailFolder marked as unread.
func (m *MailFolder) GetUnreadItemCount()(*int32) {
    return m.unreadItemCount
}
// Serialize serializes information the current object
func (m *MailFolder) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("childFolderCount", m.GetChildFolderCount())
        if err != nil {
            return err
        }
    }
    if m.GetChildFolders() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetChildFolders())
        err = writer.WriteCollectionOfObjectValues("childFolders", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isHidden", m.GetIsHidden())
        if err != nil {
            return err
        }
    }
    if m.GetMessageRules() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMessageRules())
        err = writer.WriteCollectionOfObjectValues("messageRules", cast)
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
    if m.GetMultiValueExtendedProperties() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMultiValueExtendedProperties())
        err = writer.WriteCollectionOfObjectValues("multiValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("parentFolderId", m.GetParentFolderId())
        if err != nil {
            return err
        }
    }
    if m.GetSingleValueExtendedProperties() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSingleValueExtendedProperties())
        err = writer.WriteCollectionOfObjectValues("singleValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalItemCount", m.GetTotalItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("unreadItemCount", m.GetUnreadItemCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChildFolderCount sets the childFolderCount property value. The number of immediate child mailFolders in the current mailFolder.
func (m *MailFolder) SetChildFolderCount(value *int32)() {
    m.childFolderCount = value
}
// SetChildFolders sets the childFolders property value. The collection of child folders in the mailFolder.
func (m *MailFolder) SetChildFolders(value []MailFolderable)() {
    m.childFolders = value
}
// SetDisplayName sets the displayName property value. The mailFolder's display name.
func (m *MailFolder) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetIsHidden sets the isHidden property value. Indicates whether the mailFolder is hidden. This property can be set only when creating the folder. Find more information in Hidden mail folders.
func (m *MailFolder) SetIsHidden(value *bool)() {
    m.isHidden = value
}
// SetMessageRules sets the messageRules property value. The collection of rules that apply to the user's Inbox folder.
func (m *MailFolder) SetMessageRules(value []MessageRuleable)() {
    m.messageRules = value
}
// SetMessages sets the messages property value. The collection of messages in the mailFolder.
func (m *MailFolder) SetMessages(value []Messageable)() {
    m.messages = value
}
// SetMultiValueExtendedProperties sets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the mailFolder. Read-only. Nullable.
func (m *MailFolder) SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)() {
    m.multiValueExtendedProperties = value
}
// SetParentFolderId sets the parentFolderId property value. The unique identifier for the mailFolder's parent mailFolder.
func (m *MailFolder) SetParentFolderId(value *string)() {
    m.parentFolderId = value
}
// SetSingleValueExtendedProperties sets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the mailFolder. Read-only. Nullable.
func (m *MailFolder) SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)() {
    m.singleValueExtendedProperties = value
}
// SetTotalItemCount sets the totalItemCount property value. The number of items in the mailFolder.
func (m *MailFolder) SetTotalItemCount(value *int32)() {
    m.totalItemCount = value
}
// SetUnreadItemCount sets the unreadItemCount property value. The number of items in the mailFolder marked as unread.
func (m *MailFolder) SetUnreadItemCount(value *int32)() {
    m.unreadItemCount = value
}
