package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Onenote 
type Onenote struct {
    Entity
    // The collection of OneNote notebooks that are owned by the user or group. Read-only. Nullable.
    notebooks []Notebookable
    // The status of OneNote operations. Getting an operations collection is not supported, but you can get the status of long-running operations if the Operation-Location header is returned in the response. Read-only. Nullable.
    operations []OnenoteOperationable
    // The pages in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
    pages []OnenotePageable
    // The image and other file resources in OneNote pages. Getting a resources collection is not supported, but you can get the binary content of a specific resource. Read-only. Nullable.
    resources []OnenoteResourceable
    // The section groups in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
    sectionGroups []SectionGroupable
    // The sections in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
    sections []OnenoteSectionable
}
// NewOnenote instantiates a new onenote and sets the default values.
func NewOnenote()(*Onenote) {
    m := &Onenote{
        Entity: *NewEntity(),
    }
    return m
}
// CreateOnenoteFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOnenoteFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnenote(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Onenote) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["notebooks"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateNotebookFromDiscriminatorValue , m.SetNotebooks)
    res["operations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOnenoteOperationFromDiscriminatorValue , m.SetOperations)
    res["pages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOnenotePageFromDiscriminatorValue , m.SetPages)
    res["resources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOnenoteResourceFromDiscriminatorValue , m.SetResources)
    res["sectionGroups"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSectionGroupFromDiscriminatorValue , m.SetSectionGroups)
    res["sections"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOnenoteSectionFromDiscriminatorValue , m.SetSections)
    return res
}
// GetNotebooks gets the notebooks property value. The collection of OneNote notebooks that are owned by the user or group. Read-only. Nullable.
func (m *Onenote) GetNotebooks()([]Notebookable) {
    return m.notebooks
}
// GetOperations gets the operations property value. The status of OneNote operations. Getting an operations collection is not supported, but you can get the status of long-running operations if the Operation-Location header is returned in the response. Read-only. Nullable.
func (m *Onenote) GetOperations()([]OnenoteOperationable) {
    return m.operations
}
// GetPages gets the pages property value. The pages in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) GetPages()([]OnenotePageable) {
    return m.pages
}
// GetResources gets the resources property value. The image and other file resources in OneNote pages. Getting a resources collection is not supported, but you can get the binary content of a specific resource. Read-only. Nullable.
func (m *Onenote) GetResources()([]OnenoteResourceable) {
    return m.resources
}
// GetSectionGroups gets the sectionGroups property value. The section groups in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) GetSectionGroups()([]SectionGroupable) {
    return m.sectionGroups
}
// GetSections gets the sections property value. The sections in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) GetSections()([]OnenoteSectionable) {
    return m.sections
}
// Serialize serializes information the current object
func (m *Onenote) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetNotebooks() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNotebooks())
        err = writer.WriteCollectionOfObjectValues("notebooks", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOperations())
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPages() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPages())
        err = writer.WriteCollectionOfObjectValues("pages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetResources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetResources())
        err = writer.WriteCollectionOfObjectValues("resources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSectionGroups() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSectionGroups())
        err = writer.WriteCollectionOfObjectValues("sectionGroups", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSections() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSections())
        err = writer.WriteCollectionOfObjectValues("sections", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetNotebooks sets the notebooks property value. The collection of OneNote notebooks that are owned by the user or group. Read-only. Nullable.
func (m *Onenote) SetNotebooks(value []Notebookable)() {
    m.notebooks = value
}
// SetOperations sets the operations property value. The status of OneNote operations. Getting an operations collection is not supported, but you can get the status of long-running operations if the Operation-Location header is returned in the response. Read-only. Nullable.
func (m *Onenote) SetOperations(value []OnenoteOperationable)() {
    m.operations = value
}
// SetPages sets the pages property value. The pages in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) SetPages(value []OnenotePageable)() {
    m.pages = value
}
// SetResources sets the resources property value. The image and other file resources in OneNote pages. Getting a resources collection is not supported, but you can get the binary content of a specific resource. Read-only. Nullable.
func (m *Onenote) SetResources(value []OnenoteResourceable)() {
    m.resources = value
}
// SetSectionGroups sets the sectionGroups property value. The section groups in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) SetSectionGroups(value []SectionGroupable)() {
    m.sectionGroups = value
}
// SetSections sets the sections property value. The sections in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) SetSections(value []OnenoteSectionable)() {
    m.sections = value
}
