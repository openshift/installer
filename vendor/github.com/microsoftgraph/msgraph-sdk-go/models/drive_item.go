package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DriveItem provides operations to manage the collection of agreement entities.
type DriveItem struct {
    BaseItem
    // Analytics about the view activities that took place on this item.
    analytics ItemAnalyticsable
    // Audio metadata, if the item is an audio file. Read-only. Read-only. Only on OneDrive Personal.
    audio Audioable
    // Bundle metadata, if the item is a bundle. Read-only.
    bundle Bundleable
    // Collection containing Item objects for the immediate children of Item. Only items representing folders have children. Read-only. Nullable.
    children []DriveItemable
    // The content stream, if the item represents a file.
    content []byte
    // An eTag for the content of the item. This eTag is not changed if only the metadata is changed. Note This property is not returned if the item is a folder. Read-only.
    cTag *string
    // Information about the deleted state of the item. Read-only.
    deleted Deletedable
    // File metadata, if the item is a file. Read-only.
    file Fileable
    // File system information on client. Read-write.
    fileSystemInfo FileSystemInfoable
    // Folder metadata, if the item is a folder. Read-only.
    folder Folderable
    // Image metadata, if the item is an image. Read-only.
    image Imageable
    // For drives in SharePoint, the associated document library list item. Read-only. Nullable.
    listItem ListItemable
    // Location metadata, if the item has location data. Read-only.
    location GeoCoordinatesable
    // Malware metadata, if the item was detected to contain malware. Read-only.
    malware Malwareable
    // If present, indicates that this item is a package instead of a folder or file. Packages are treated like files in some contexts and folders in others. Read-only.
    package_escaped Package_escapedable
    // If present, indicates that one or more operations that might affect the state of the driveItem are pending completion. Read-only.
    pendingOperations PendingOperationsable
    // The set of permissions for the item. Read-only. Nullable.
    permissions []Permissionable
    // Photo metadata, if the item is a photo. Read-only.
    photo Photoable
    // Provides information about the published or checked-out state of an item, in locations that support such actions. This property is not returned by default. Read-only.
    publication PublicationFacetable
    // Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.
    remoteItem RemoteItemable
    // If this property is non-null, it indicates that the driveItem is the top-most driveItem in the drive.
    root Rootable
    // Search metadata, if the item is from a search result. Read-only.
    searchResult SearchResultable
    // Indicates that the item has been shared with others and provides information about the shared state of the item. Read-only.
    shared Sharedable
    // Returns identifiers useful for SharePoint REST compatibility. Read-only.
    sharepointIds SharepointIdsable
    // Size of the item in bytes. Read-only.
    size *int64
    // If the current item is also available as a special folder, this facet is returned. Read-only.
    specialFolder SpecialFolderable
    // The set of subscriptions on the item. Only supported on the root of a drive.
    subscriptions []Subscriptionable
    // Collection containing [ThumbnailSet][] objects associated with the item. For more info, see [getting thumbnails][]. Read-only. Nullable.
    thumbnails []ThumbnailSetable
    // The list of previous versions of the item. For more info, see [getting previous versions][]. Read-only. Nullable.
    versions []DriveItemVersionable
    // Video metadata, if the item is a video. Read-only.
    video Videoable
    // WebDAV compatible URL for the item.
    webDavUrl *string
    // For files that are Excel spreadsheets, accesses the workbook API to work with the spreadsheet's contents. Nullable.
    workbook Workbookable
}
// NewDriveItem instantiates a new driveItem and sets the default values.
func NewDriveItem()(*DriveItem) {
    m := &DriveItem{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.driveItem";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateDriveItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDriveItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDriveItem(), nil
}
// GetAnalytics gets the analytics property value. Analytics about the view activities that took place on this item.
func (m *DriveItem) GetAnalytics()(ItemAnalyticsable) {
    return m.analytics
}
// GetAudio gets the audio property value. Audio metadata, if the item is an audio file. Read-only. Read-only. Only on OneDrive Personal.
func (m *DriveItem) GetAudio()(Audioable) {
    return m.audio
}
// GetBundle gets the bundle property value. Bundle metadata, if the item is a bundle. Read-only.
func (m *DriveItem) GetBundle()(Bundleable) {
    return m.bundle
}
// GetChildren gets the children property value. Collection containing Item objects for the immediate children of Item. Only items representing folders have children. Read-only. Nullable.
func (m *DriveItem) GetChildren()([]DriveItemable) {
    return m.children
}
// GetContent gets the content property value. The content stream, if the item represents a file.
func (m *DriveItem) GetContent()([]byte) {
    return m.content
}
// GetCTag gets the cTag property value. An eTag for the content of the item. This eTag is not changed if only the metadata is changed. Note This property is not returned if the item is a folder. Read-only.
func (m *DriveItem) GetCTag()(*string) {
    return m.cTag
}
// GetDeleted gets the deleted property value. Information about the deleted state of the item. Read-only.
func (m *DriveItem) GetDeleted()(Deletedable) {
    return m.deleted
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DriveItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["analytics"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateItemAnalyticsFromDiscriminatorValue , m.SetAnalytics)
    res["audio"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateAudioFromDiscriminatorValue , m.SetAudio)
    res["bundle"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateBundleFromDiscriminatorValue , m.SetBundle)
    res["children"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDriveItemFromDiscriminatorValue , m.SetChildren)
    res["content"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetByteArrayValue(m.SetContent)
    res["cTag"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCTag)
    res["deleted"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeletedFromDiscriminatorValue , m.SetDeleted)
    res["file"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateFileFromDiscriminatorValue , m.SetFile)
    res["fileSystemInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateFileSystemInfoFromDiscriminatorValue , m.SetFileSystemInfo)
    res["folder"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateFolderFromDiscriminatorValue , m.SetFolder)
    res["image"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateImageFromDiscriminatorValue , m.SetImage)
    res["listItem"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateListItemFromDiscriminatorValue , m.SetListItem)
    res["location"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateGeoCoordinatesFromDiscriminatorValue , m.SetLocation)
    res["malware"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateMalwareFromDiscriminatorValue , m.SetMalware)
    res["package"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePackage_escapedFromDiscriminatorValue , m.SetPackage)
    res["pendingOperations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePendingOperationsFromDiscriminatorValue , m.SetPendingOperations)
    res["permissions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePermissionFromDiscriminatorValue , m.SetPermissions)
    res["photo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePhotoFromDiscriminatorValue , m.SetPhoto)
    res["publication"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePublicationFacetFromDiscriminatorValue , m.SetPublication)
    res["remoteItem"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRemoteItemFromDiscriminatorValue , m.SetRemoteItem)
    res["root"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRootFromDiscriminatorValue , m.SetRoot)
    res["searchResult"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSearchResultFromDiscriminatorValue , m.SetSearchResult)
    res["shared"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSharedFromDiscriminatorValue , m.SetShared)
    res["sharepointIds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSharepointIdsFromDiscriminatorValue , m.SetSharepointIds)
    res["size"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt64Value(m.SetSize)
    res["specialFolder"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSpecialFolderFromDiscriminatorValue , m.SetSpecialFolder)
    res["subscriptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubscriptionFromDiscriminatorValue , m.SetSubscriptions)
    res["thumbnails"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateThumbnailSetFromDiscriminatorValue , m.SetThumbnails)
    res["versions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDriveItemVersionFromDiscriminatorValue , m.SetVersions)
    res["video"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateVideoFromDiscriminatorValue , m.SetVideo)
    res["webDavUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetWebDavUrl)
    res["workbook"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWorkbookFromDiscriminatorValue , m.SetWorkbook)
    return res
}
// GetFile gets the file property value. File metadata, if the item is a file. Read-only.
func (m *DriveItem) GetFile()(Fileable) {
    return m.file
}
// GetFileSystemInfo gets the fileSystemInfo property value. File system information on client. Read-write.
func (m *DriveItem) GetFileSystemInfo()(FileSystemInfoable) {
    return m.fileSystemInfo
}
// GetFolder gets the folder property value. Folder metadata, if the item is a folder. Read-only.
func (m *DriveItem) GetFolder()(Folderable) {
    return m.folder
}
// GetImage gets the image property value. Image metadata, if the item is an image. Read-only.
func (m *DriveItem) GetImage()(Imageable) {
    return m.image
}
// GetListItem gets the listItem property value. For drives in SharePoint, the associated document library list item. Read-only. Nullable.
func (m *DriveItem) GetListItem()(ListItemable) {
    return m.listItem
}
// GetLocation gets the location property value. Location metadata, if the item has location data. Read-only.
func (m *DriveItem) GetLocation()(GeoCoordinatesable) {
    return m.location
}
// GetMalware gets the malware property value. Malware metadata, if the item was detected to contain malware. Read-only.
func (m *DriveItem) GetMalware()(Malwareable) {
    return m.malware
}
// GetPackage gets the package property value. If present, indicates that this item is a package instead of a folder or file. Packages are treated like files in some contexts and folders in others. Read-only.
func (m *DriveItem) GetPackage()(Package_escapedable) {
    return m.package_escaped
}
// GetPendingOperations gets the pendingOperations property value. If present, indicates that one or more operations that might affect the state of the driveItem are pending completion. Read-only.
func (m *DriveItem) GetPendingOperations()(PendingOperationsable) {
    return m.pendingOperations
}
// GetPermissions gets the permissions property value. The set of permissions for the item. Read-only. Nullable.
func (m *DriveItem) GetPermissions()([]Permissionable) {
    return m.permissions
}
// GetPhoto gets the photo property value. Photo metadata, if the item is a photo. Read-only.
func (m *DriveItem) GetPhoto()(Photoable) {
    return m.photo
}
// GetPublication gets the publication property value. Provides information about the published or checked-out state of an item, in locations that support such actions. This property is not returned by default. Read-only.
func (m *DriveItem) GetPublication()(PublicationFacetable) {
    return m.publication
}
// GetRemoteItem gets the remoteItem property value. Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.
func (m *DriveItem) GetRemoteItem()(RemoteItemable) {
    return m.remoteItem
}
// GetRoot gets the root property value. If this property is non-null, it indicates that the driveItem is the top-most driveItem in the drive.
func (m *DriveItem) GetRoot()(Rootable) {
    return m.root
}
// GetSearchResult gets the searchResult property value. Search metadata, if the item is from a search result. Read-only.
func (m *DriveItem) GetSearchResult()(SearchResultable) {
    return m.searchResult
}
// GetShared gets the shared property value. Indicates that the item has been shared with others and provides information about the shared state of the item. Read-only.
func (m *DriveItem) GetShared()(Sharedable) {
    return m.shared
}
// GetSharepointIds gets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *DriveItem) GetSharepointIds()(SharepointIdsable) {
    return m.sharepointIds
}
// GetSize gets the size property value. Size of the item in bytes. Read-only.
func (m *DriveItem) GetSize()(*int64) {
    return m.size
}
// GetSpecialFolder gets the specialFolder property value. If the current item is also available as a special folder, this facet is returned. Read-only.
func (m *DriveItem) GetSpecialFolder()(SpecialFolderable) {
    return m.specialFolder
}
// GetSubscriptions gets the subscriptions property value. The set of subscriptions on the item. Only supported on the root of a drive.
func (m *DriveItem) GetSubscriptions()([]Subscriptionable) {
    return m.subscriptions
}
// GetThumbnails gets the thumbnails property value. Collection containing [ThumbnailSet][] objects associated with the item. For more info, see [getting thumbnails][]. Read-only. Nullable.
func (m *DriveItem) GetThumbnails()([]ThumbnailSetable) {
    return m.thumbnails
}
// GetVersions gets the versions property value. The list of previous versions of the item. For more info, see [getting previous versions][]. Read-only. Nullable.
func (m *DriveItem) GetVersions()([]DriveItemVersionable) {
    return m.versions
}
// GetVideo gets the video property value. Video metadata, if the item is a video. Read-only.
func (m *DriveItem) GetVideo()(Videoable) {
    return m.video
}
// GetWebDavUrl gets the webDavUrl property value. WebDAV compatible URL for the item.
func (m *DriveItem) GetWebDavUrl()(*string) {
    return m.webDavUrl
}
// GetWorkbook gets the workbook property value. For files that are Excel spreadsheets, accesses the workbook API to work with the spreadsheet's contents. Nullable.
func (m *DriveItem) GetWorkbook()(Workbookable) {
    return m.workbook
}
// Serialize serializes information the current object
func (m *DriveItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteObjectValue("audio", m.GetAudio())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("bundle", m.GetBundle())
        if err != nil {
            return err
        }
    }
    if m.GetChildren() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetChildren())
        err = writer.WriteCollectionOfObjectValues("children", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("content", m.GetContent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("cTag", m.GetCTag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deleted", m.GetDeleted())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("file", m.GetFile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("fileSystemInfo", m.GetFileSystemInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("folder", m.GetFolder())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("image", m.GetImage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("listItem", m.GetListItem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("location", m.GetLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("malware", m.GetMalware())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("package", m.GetPackage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("pendingOperations", m.GetPendingOperations())
        if err != nil {
            return err
        }
    }
    if m.GetPermissions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPermissions())
        err = writer.WriteCollectionOfObjectValues("permissions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("photo", m.GetPhoto())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("publication", m.GetPublication())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("remoteItem", m.GetRemoteItem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("root", m.GetRoot())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("searchResult", m.GetSearchResult())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("shared", m.GetShared())
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
    {
        err = writer.WriteInt64Value("size", m.GetSize())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("specialFolder", m.GetSpecialFolder())
        if err != nil {
            return err
        }
    }
    if m.GetSubscriptions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSubscriptions())
        err = writer.WriteCollectionOfObjectValues("subscriptions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetThumbnails() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetThumbnails())
        err = writer.WriteCollectionOfObjectValues("thumbnails", cast)
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
    {
        err = writer.WriteObjectValue("video", m.GetVideo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webDavUrl", m.GetWebDavUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("workbook", m.GetWorkbook())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAnalytics sets the analytics property value. Analytics about the view activities that took place on this item.
func (m *DriveItem) SetAnalytics(value ItemAnalyticsable)() {
    m.analytics = value
}
// SetAudio sets the audio property value. Audio metadata, if the item is an audio file. Read-only. Read-only. Only on OneDrive Personal.
func (m *DriveItem) SetAudio(value Audioable)() {
    m.audio = value
}
// SetBundle sets the bundle property value. Bundle metadata, if the item is a bundle. Read-only.
func (m *DriveItem) SetBundle(value Bundleable)() {
    m.bundle = value
}
// SetChildren sets the children property value. Collection containing Item objects for the immediate children of Item. Only items representing folders have children. Read-only. Nullable.
func (m *DriveItem) SetChildren(value []DriveItemable)() {
    m.children = value
}
// SetContent sets the content property value. The content stream, if the item represents a file.
func (m *DriveItem) SetContent(value []byte)() {
    m.content = value
}
// SetCTag sets the cTag property value. An eTag for the content of the item. This eTag is not changed if only the metadata is changed. Note This property is not returned if the item is a folder. Read-only.
func (m *DriveItem) SetCTag(value *string)() {
    m.cTag = value
}
// SetDeleted sets the deleted property value. Information about the deleted state of the item. Read-only.
func (m *DriveItem) SetDeleted(value Deletedable)() {
    m.deleted = value
}
// SetFile sets the file property value. File metadata, if the item is a file. Read-only.
func (m *DriveItem) SetFile(value Fileable)() {
    m.file = value
}
// SetFileSystemInfo sets the fileSystemInfo property value. File system information on client. Read-write.
func (m *DriveItem) SetFileSystemInfo(value FileSystemInfoable)() {
    m.fileSystemInfo = value
}
// SetFolder sets the folder property value. Folder metadata, if the item is a folder. Read-only.
func (m *DriveItem) SetFolder(value Folderable)() {
    m.folder = value
}
// SetImage sets the image property value. Image metadata, if the item is an image. Read-only.
func (m *DriveItem) SetImage(value Imageable)() {
    m.image = value
}
// SetListItem sets the listItem property value. For drives in SharePoint, the associated document library list item. Read-only. Nullable.
func (m *DriveItem) SetListItem(value ListItemable)() {
    m.listItem = value
}
// SetLocation sets the location property value. Location metadata, if the item has location data. Read-only.
func (m *DriveItem) SetLocation(value GeoCoordinatesable)() {
    m.location = value
}
// SetMalware sets the malware property value. Malware metadata, if the item was detected to contain malware. Read-only.
func (m *DriveItem) SetMalware(value Malwareable)() {
    m.malware = value
}
// SetPackage sets the package property value. If present, indicates that this item is a package instead of a folder or file. Packages are treated like files in some contexts and folders in others. Read-only.
func (m *DriveItem) SetPackage(value Package_escapedable)() {
    m.package_escaped = value
}
// SetPendingOperations sets the pendingOperations property value. If present, indicates that one or more operations that might affect the state of the driveItem are pending completion. Read-only.
func (m *DriveItem) SetPendingOperations(value PendingOperationsable)() {
    m.pendingOperations = value
}
// SetPermissions sets the permissions property value. The set of permissions for the item. Read-only. Nullable.
func (m *DriveItem) SetPermissions(value []Permissionable)() {
    m.permissions = value
}
// SetPhoto sets the photo property value. Photo metadata, if the item is a photo. Read-only.
func (m *DriveItem) SetPhoto(value Photoable)() {
    m.photo = value
}
// SetPublication sets the publication property value. Provides information about the published or checked-out state of an item, in locations that support such actions. This property is not returned by default. Read-only.
func (m *DriveItem) SetPublication(value PublicationFacetable)() {
    m.publication = value
}
// SetRemoteItem sets the remoteItem property value. Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.
func (m *DriveItem) SetRemoteItem(value RemoteItemable)() {
    m.remoteItem = value
}
// SetRoot sets the root property value. If this property is non-null, it indicates that the driveItem is the top-most driveItem in the drive.
func (m *DriveItem) SetRoot(value Rootable)() {
    m.root = value
}
// SetSearchResult sets the searchResult property value. Search metadata, if the item is from a search result. Read-only.
func (m *DriveItem) SetSearchResult(value SearchResultable)() {
    m.searchResult = value
}
// SetShared sets the shared property value. Indicates that the item has been shared with others and provides information about the shared state of the item. Read-only.
func (m *DriveItem) SetShared(value Sharedable)() {
    m.shared = value
}
// SetSharepointIds sets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *DriveItem) SetSharepointIds(value SharepointIdsable)() {
    m.sharepointIds = value
}
// SetSize sets the size property value. Size of the item in bytes. Read-only.
func (m *DriveItem) SetSize(value *int64)() {
    m.size = value
}
// SetSpecialFolder sets the specialFolder property value. If the current item is also available as a special folder, this facet is returned. Read-only.
func (m *DriveItem) SetSpecialFolder(value SpecialFolderable)() {
    m.specialFolder = value
}
// SetSubscriptions sets the subscriptions property value. The set of subscriptions on the item. Only supported on the root of a drive.
func (m *DriveItem) SetSubscriptions(value []Subscriptionable)() {
    m.subscriptions = value
}
// SetThumbnails sets the thumbnails property value. Collection containing [ThumbnailSet][] objects associated with the item. For more info, see [getting thumbnails][]. Read-only. Nullable.
func (m *DriveItem) SetThumbnails(value []ThumbnailSetable)() {
    m.thumbnails = value
}
// SetVersions sets the versions property value. The list of previous versions of the item. For more info, see [getting previous versions][]. Read-only. Nullable.
func (m *DriveItem) SetVersions(value []DriveItemVersionable)() {
    m.versions = value
}
// SetVideo sets the video property value. Video metadata, if the item is a video. Read-only.
func (m *DriveItem) SetVideo(value Videoable)() {
    m.video = value
}
// SetWebDavUrl sets the webDavUrl property value. WebDAV compatible URL for the item.
func (m *DriveItem) SetWebDavUrl(value *string)() {
    m.webDavUrl = value
}
// SetWorkbook sets the workbook property value. For files that are Excel spreadsheets, accesses the workbook API to work with the spreadsheet's contents. Nullable.
func (m *DriveItem) SetWorkbook(value Workbookable)() {
    m.workbook = value
}
