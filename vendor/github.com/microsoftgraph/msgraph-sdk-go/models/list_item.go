package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ListItem 
type ListItem struct {
    BaseItem
    // Analytics about the view activities that took place on this item.
    analytics ItemAnalyticsable
    // The content type of this list item
    contentType ContentTypeInfoable
    // Version information for a document set version created by a user.
    documentSetVersions []DocumentSetVersionable
    // For document libraries, the driveItem relationship exposes the listItem as a [driveItem][]
    driveItem DriveItemable
    // The values of the columns set on this list item.
    fields FieldValueSetable
    // Returns identifiers useful for SharePoint REST compatibility. Read-only.
    sharepointIds SharepointIdsable
    // The list of previous versions of the list item.
    versions []ListItemVersionable
}
// NewListItem instantiates a new listItem and sets the default values.
func NewListItem()(*ListItem) {
    m := &ListItem{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.listItem";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateListItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateListItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewListItem(), nil
}
// GetAnalytics gets the analytics property value. Analytics about the view activities that took place on this item.
func (m *ListItem) GetAnalytics()(ItemAnalyticsable) {
    return m.analytics
}
// GetContentType gets the contentType property value. The content type of this list item
func (m *ListItem) GetContentType()(ContentTypeInfoable) {
    return m.contentType
}
// GetDocumentSetVersions gets the documentSetVersions property value. Version information for a document set version created by a user.
func (m *ListItem) GetDocumentSetVersions()([]DocumentSetVersionable) {
    return m.documentSetVersions
}
// GetDriveItem gets the driveItem property value. For document libraries, the driveItem relationship exposes the listItem as a [driveItem][]
func (m *ListItem) GetDriveItem()(DriveItemable) {
    return m.driveItem
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ListItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["analytics"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateItemAnalyticsFromDiscriminatorValue , m.SetAnalytics)
    res["contentType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateContentTypeInfoFromDiscriminatorValue , m.SetContentType)
    res["documentSetVersions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDocumentSetVersionFromDiscriminatorValue , m.SetDocumentSetVersions)
    res["driveItem"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDriveItemFromDiscriminatorValue , m.SetDriveItem)
    res["fields"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateFieldValueSetFromDiscriminatorValue , m.SetFields)
    res["sharepointIds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSharepointIdsFromDiscriminatorValue , m.SetSharepointIds)
    res["versions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateListItemVersionFromDiscriminatorValue , m.SetVersions)
    return res
}
// GetFields gets the fields property value. The values of the columns set on this list item.
func (m *ListItem) GetFields()(FieldValueSetable) {
    return m.fields
}
// GetSharepointIds gets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *ListItem) GetSharepointIds()(SharepointIdsable) {
    return m.sharepointIds
}
// GetVersions gets the versions property value. The list of previous versions of the list item.
func (m *ListItem) GetVersions()([]ListItemVersionable) {
    return m.versions
}
// Serialize serializes information the current object
func (m *ListItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("analytics", m.GetAnalytics())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("contentType", m.GetContentType())
        if err != nil {
            return err
        }
    }
    if m.GetDocumentSetVersions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDocumentSetVersions())
        err = writer.WriteCollectionOfObjectValues("documentSetVersions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("driveItem", m.GetDriveItem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("fields", m.GetFields())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharepointIds", m.GetSharepointIds())
        if err != nil {
            return err
        }
    }
    if m.GetVersions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetVersions())
        err = writer.WriteCollectionOfObjectValues("versions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAnalytics sets the analytics property value. Analytics about the view activities that took place on this item.
func (m *ListItem) SetAnalytics(value ItemAnalyticsable)() {
    m.analytics = value
}
// SetContentType sets the contentType property value. The content type of this list item
func (m *ListItem) SetContentType(value ContentTypeInfoable)() {
    m.contentType = value
}
// SetDocumentSetVersions sets the documentSetVersions property value. Version information for a document set version created by a user.
func (m *ListItem) SetDocumentSetVersions(value []DocumentSetVersionable)() {
    m.documentSetVersions = value
}
// SetDriveItem sets the driveItem property value. For document libraries, the driveItem relationship exposes the listItem as a [driveItem][]
func (m *ListItem) SetDriveItem(value DriveItemable)() {
    m.driveItem = value
}
// SetFields sets the fields property value. The values of the columns set on this list item.
func (m *ListItem) SetFields(value FieldValueSetable)() {
    m.fields = value
}
// SetSharepointIds sets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *ListItem) SetSharepointIds(value SharepointIdsable)() {
    m.sharepointIds = value
}
// SetVersions sets the versions property value. The list of previous versions of the list item.
func (m *ListItem) SetVersions(value []ListItemVersionable)() {
    m.versions = value
}
