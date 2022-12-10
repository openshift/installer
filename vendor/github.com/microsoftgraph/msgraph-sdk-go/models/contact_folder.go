package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ContactFolder provides operations to manage the collection of agreement entities.
type ContactFolder struct {
    Entity
    // The collection of child folders in the folder. Navigation property. Read-only. Nullable.
    childFolders []ContactFolderable
    // The contacts in the folder. Navigation property. Read-only. Nullable.
    contacts []Contactable
    // The folder's display name.
    displayName *string
    // The collection of multi-value extended properties defined for the contactFolder. Read-only. Nullable.
    multiValueExtendedProperties []MultiValueLegacyExtendedPropertyable
    // The ID of the folder's parent folder.
    parentFolderId *string
    // The collection of single-value extended properties defined for the contactFolder. Read-only. Nullable.
    singleValueExtendedProperties []SingleValueLegacyExtendedPropertyable
}
// NewContactFolder instantiates a new contactFolder and sets the default values.
func NewContactFolder()(*ContactFolder) {
    m := &ContactFolder{
        Entity: *NewEntity(),
    }
    return m
}
// CreateContactFolderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateContactFolderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewContactFolder(), nil
}
// GetChildFolders gets the childFolders property value. The collection of child folders in the folder. Navigation property. Read-only. Nullable.
func (m *ContactFolder) GetChildFolders()([]ContactFolderable) {
    return m.childFolders
}
// GetContacts gets the contacts property value. The contacts in the folder. Navigation property. Read-only. Nullable.
func (m *ContactFolder) GetContacts()([]Contactable) {
    return m.contacts
}
// GetDisplayName gets the displayName property value. The folder's display name.
func (m *ContactFolder) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ContactFolder) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["childFolders"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateContactFolderFromDiscriminatorValue , m.SetChildFolders)
    res["contacts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateContactFromDiscriminatorValue , m.SetContacts)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["multiValueExtendedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMultiValueLegacyExtendedPropertyFromDiscriminatorValue , m.SetMultiValueExtendedProperties)
    res["parentFolderId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetParentFolderId)
    res["singleValueExtendedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue , m.SetSingleValueExtendedProperties)
    return res
}
// GetMultiValueExtendedProperties gets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the contactFolder. Read-only. Nullable.
func (m *ContactFolder) GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable) {
    return m.multiValueExtendedProperties
}
// GetParentFolderId gets the parentFolderId property value. The ID of the folder's parent folder.
func (m *ContactFolder) GetParentFolderId()(*string) {
    return m.parentFolderId
}
// GetSingleValueExtendedProperties gets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the contactFolder. Read-only. Nullable.
func (m *ContactFolder) GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable) {
    return m.singleValueExtendedProperties
}
// Serialize serializes information the current object
func (m *ContactFolder) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChildFolders() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetChildFolders())
        err = writer.WriteCollectionOfObjectValues("childFolders", cast)
        if err != nil {
            return err
        }
    }
    if m.GetContacts() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetContacts())
        err = writer.WriteCollectionOfObjectValues("contacts", cast)
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
    return nil
}
// SetChildFolders sets the childFolders property value. The collection of child folders in the folder. Navigation property. Read-only. Nullable.
func (m *ContactFolder) SetChildFolders(value []ContactFolderable)() {
    m.childFolders = value
}
// SetContacts sets the contacts property value. The contacts in the folder. Navigation property. Read-only. Nullable.
func (m *ContactFolder) SetContacts(value []Contactable)() {
    m.contacts = value
}
// SetDisplayName sets the displayName property value. The folder's display name.
func (m *ContactFolder) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetMultiValueExtendedProperties sets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the contactFolder. Read-only. Nullable.
func (m *ContactFolder) SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)() {
    m.multiValueExtendedProperties = value
}
// SetParentFolderId sets the parentFolderId property value. The ID of the folder's parent folder.
func (m *ContactFolder) SetParentFolderId(value *string)() {
    m.parentFolderId = value
}
// SetSingleValueExtendedProperties sets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the contactFolder. Read-only. Nullable.
func (m *ContactFolder) SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)() {
    m.singleValueExtendedProperties = value
}
