package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ItemReference 
type ItemReference struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Unique identifier of the drive instance that contains the item. Read-only.
    driveId *string
    // Identifies the type of drive. See [drive][] resource for values.
    driveType *string
    // Unique identifier of the item in the drive. Read-only.
    id *string
    // The name of the item being referenced. Read-only.
    name *string
    // The OdataType property
    odataType *string
    // Path that can be used to navigate to the item. Read-only.
    path *string
    // A unique identifier for a shared resource that can be accessed via the [Shares][] API.
    shareId *string
    // Returns identifiers useful for SharePoint REST compatibility. Read-only.
    sharepointIds SharepointIdsable
    // For OneDrive for Business and SharePoint, this property represents the ID of the site that contains the parent document library of the driveItem resource. The value is the same as the id property of that [site][] resource. It is an opaque string that consists of three identifiers of the site. For OneDrive, this property is not populated.
    siteId *string
}
// NewItemReference instantiates a new itemReference and sets the default values.
func NewItemReference()(*ItemReference) {
    m := &ItemReference{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateItemReferenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateItemReferenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemReference(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ItemReference) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDriveId gets the driveId property value. Unique identifier of the drive instance that contains the item. Read-only.
func (m *ItemReference) GetDriveId()(*string) {
    return m.driveId
}
// GetDriveType gets the driveType property value. Identifies the type of drive. See [drive][] resource for values.
func (m *ItemReference) GetDriveType()(*string) {
    return m.driveType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ItemReference) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["driveId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDriveId)
    res["driveType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDriveType)
    res["id"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetId)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["path"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPath)
    res["shareId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetShareId)
    res["sharepointIds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSharepointIdsFromDiscriminatorValue , m.SetSharepointIds)
    res["siteId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSiteId)
    return res
}
// GetId gets the id property value. Unique identifier of the item in the drive. Read-only.
func (m *ItemReference) GetId()(*string) {
    return m.id
}
// GetName gets the name property value. The name of the item being referenced. Read-only.
func (m *ItemReference) GetName()(*string) {
    return m.name
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ItemReference) GetOdataType()(*string) {
    return m.odataType
}
// GetPath gets the path property value. Path that can be used to navigate to the item. Read-only.
func (m *ItemReference) GetPath()(*string) {
    return m.path
}
// GetShareId gets the shareId property value. A unique identifier for a shared resource that can be accessed via the [Shares][] API.
func (m *ItemReference) GetShareId()(*string) {
    return m.shareId
}
// GetSharepointIds gets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *ItemReference) GetSharepointIds()(SharepointIdsable) {
    return m.sharepointIds
}
// GetSiteId gets the siteId property value. For OneDrive for Business and SharePoint, this property represents the ID of the site that contains the parent document library of the driveItem resource. The value is the same as the id property of that [site][] resource. It is an opaque string that consists of three identifiers of the site. For OneDrive, this property is not populated.
func (m *ItemReference) GetSiteId()(*string) {
    return m.siteId
}
// Serialize serializes information the current object
func (m *ItemReference) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("driveId", m.GetDriveId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("driveType", m.GetDriveType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("id", m.GetId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("path", m.GetPath())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("shareId", m.GetShareId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("sharepointIds", m.GetSharepointIds())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("siteId", m.GetSiteId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ItemReference) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDriveId sets the driveId property value. Unique identifier of the drive instance that contains the item. Read-only.
func (m *ItemReference) SetDriveId(value *string)() {
    m.driveId = value
}
// SetDriveType sets the driveType property value. Identifies the type of drive. See [drive][] resource for values.
func (m *ItemReference) SetDriveType(value *string)() {
    m.driveType = value
}
// SetId sets the id property value. Unique identifier of the item in the drive. Read-only.
func (m *ItemReference) SetId(value *string)() {
    m.id = value
}
// SetName sets the name property value. The name of the item being referenced. Read-only.
func (m *ItemReference) SetName(value *string)() {
    m.name = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ItemReference) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPath sets the path property value. Path that can be used to navigate to the item. Read-only.
func (m *ItemReference) SetPath(value *string)() {
    m.path = value
}
// SetShareId sets the shareId property value. A unique identifier for a shared resource that can be accessed via the [Shares][] API.
func (m *ItemReference) SetShareId(value *string)() {
    m.shareId = value
}
// SetSharepointIds sets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *ItemReference) SetSharepointIds(value SharepointIdsable)() {
    m.sharepointIds = value
}
// SetSiteId sets the siteId property value. For OneDrive for Business and SharePoint, this property represents the ID of the site that contains the parent document library of the driveItem resource. The value is the same as the id property of that [site][] resource. It is an opaque string that consists of three identifiers of the site. For OneDrive, this property is not populated.
func (m *ItemReference) SetSiteId(value *string)() {
    m.siteId = value
}
