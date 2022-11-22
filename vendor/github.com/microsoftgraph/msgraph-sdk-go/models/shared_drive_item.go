package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SharedDriveItem 
type SharedDriveItem struct {
    BaseItem
    // Used to access the underlying driveItem
    driveItem DriveItemable
    // All driveItems contained in the sharing root. This collection cannot be enumerated.
    items []DriveItemable
    // Used to access the underlying list
    list Listable
    // Used to access the underlying listItem
    listItem ListItemable
    // Information about the owner of the shared item being referenced.
    owner IdentitySetable
    // Used to access the permission representing the underlying sharing link
    permission Permissionable
    // Used to access the underlying driveItem. Deprecated -- use driveItem instead.
    root DriveItemable
    // Used to access the underlying site
    site Siteable
}
// NewSharedDriveItem instantiates a new SharedDriveItem and sets the default values.
func NewSharedDriveItem()(*SharedDriveItem) {
    m := &SharedDriveItem{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.sharedDriveItem";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateSharedDriveItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSharedDriveItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharedDriveItem(), nil
}
// GetDriveItem gets the driveItem property value. Used to access the underlying driveItem
func (m *SharedDriveItem) GetDriveItem()(DriveItemable) {
    return m.driveItem
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SharedDriveItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["driveItem"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDriveItemFromDiscriminatorValue , m.SetDriveItem)
    res["items"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDriveItemFromDiscriminatorValue , m.SetItems)
    res["list"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateListFromDiscriminatorValue , m.SetList)
    res["listItem"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateListItemFromDiscriminatorValue , m.SetListItem)
    res["owner"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetOwner)
    res["permission"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePermissionFromDiscriminatorValue , m.SetPermission)
    res["root"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDriveItemFromDiscriminatorValue , m.SetRoot)
    res["site"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSiteFromDiscriminatorValue , m.SetSite)
    return res
}
// GetItems gets the items property value. All driveItems contained in the sharing root. This collection cannot be enumerated.
func (m *SharedDriveItem) GetItems()([]DriveItemable) {
    return m.items
}
// GetList gets the list property value. Used to access the underlying list
func (m *SharedDriveItem) GetList()(Listable) {
    return m.list
}
// GetListItem gets the listItem property value. Used to access the underlying listItem
func (m *SharedDriveItem) GetListItem()(ListItemable) {
    return m.listItem
}
// GetOwner gets the owner property value. Information about the owner of the shared item being referenced.
func (m *SharedDriveItem) GetOwner()(IdentitySetable) {
    return m.owner
}
// GetPermission gets the permission property value. Used to access the permission representing the underlying sharing link
func (m *SharedDriveItem) GetPermission()(Permissionable) {
    return m.permission
}
// GetRoot gets the root property value. Used to access the underlying driveItem. Deprecated -- use driveItem instead.
func (m *SharedDriveItem) GetRoot()(DriveItemable) {
    return m.root
}
// GetSite gets the site property value. Used to access the underlying site
func (m *SharedDriveItem) GetSite()(Siteable) {
    return m.site
}
// Serialize serializes information the current object
func (m *SharedDriveItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("driveItem", m.GetDriveItem())
        if err != nil {
            return err
        }
    }
    if m.GetItems() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetItems())
        err = writer.WriteCollectionOfObjectValues("items", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("list", m.GetList())
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
        err = writer.WriteObjectValue("owner", m.GetOwner())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("permission", m.GetPermission())
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
        err = writer.WriteObjectValue("site", m.GetSite())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDriveItem sets the driveItem property value. Used to access the underlying driveItem
func (m *SharedDriveItem) SetDriveItem(value DriveItemable)() {
    m.driveItem = value
}
// SetItems sets the items property value. All driveItems contained in the sharing root. This collection cannot be enumerated.
func (m *SharedDriveItem) SetItems(value []DriveItemable)() {
    m.items = value
}
// SetList sets the list property value. Used to access the underlying list
func (m *SharedDriveItem) SetList(value Listable)() {
    m.list = value
}
// SetListItem sets the listItem property value. Used to access the underlying listItem
func (m *SharedDriveItem) SetListItem(value ListItemable)() {
    m.listItem = value
}
// SetOwner sets the owner property value. Information about the owner of the shared item being referenced.
func (m *SharedDriveItem) SetOwner(value IdentitySetable)() {
    m.owner = value
}
// SetPermission sets the permission property value. Used to access the permission representing the underlying sharing link
func (m *SharedDriveItem) SetPermission(value Permissionable)() {
    m.permission = value
}
// SetRoot sets the root property value. Used to access the underlying driveItem. Deprecated -- use driveItem instead.
func (m *SharedDriveItem) SetRoot(value DriveItemable)() {
    m.root = value
}
// SetSite sets the site property value. Used to access the underlying site
func (m *SharedDriveItem) SetSite(value Siteable)() {
    m.site = value
}
