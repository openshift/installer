package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SharedInsight provides operations to manage the collection of agreement entities.
type SharedInsight struct {
    Entity
    // Details about the shared item. Read only.
    lastShared SharingDetailable
    // The lastSharedMethod property
    lastSharedMethod Entityable
    // Used for navigating to the item that was shared. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
    resource Entityable
    // Reference properties of the shared document, such as the url and type of the document. Read-only
    resourceReference ResourceReferenceable
    // Properties that you can use to visualize the document in your experience. Read-only
    resourceVisualization ResourceVisualizationable
    // The sharingHistory property
    sharingHistory []SharingDetailable
}
// NewSharedInsight instantiates a new sharedInsight and sets the default values.
func NewSharedInsight()(*SharedInsight) {
    m := &SharedInsight{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSharedInsightFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSharedInsightFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharedInsight(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SharedInsight) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["lastShared"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSharingDetailFromDiscriminatorValue , m.SetLastShared)
    res["lastSharedMethod"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEntityFromDiscriminatorValue , m.SetLastSharedMethod)
    res["resource"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEntityFromDiscriminatorValue , m.SetResource)
    res["resourceReference"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateResourceReferenceFromDiscriminatorValue , m.SetResourceReference)
    res["resourceVisualization"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateResourceVisualizationFromDiscriminatorValue , m.SetResourceVisualization)
    res["sharingHistory"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSharingDetailFromDiscriminatorValue , m.SetSharingHistory)
    return res
}
// GetLastShared gets the lastShared property value. Details about the shared item. Read only.
func (m *SharedInsight) GetLastShared()(SharingDetailable) {
    return m.lastShared
}
// GetLastSharedMethod gets the lastSharedMethod property value. The lastSharedMethod property
func (m *SharedInsight) GetLastSharedMethod()(Entityable) {
    return m.lastSharedMethod
}
// GetResource gets the resource property value. Used for navigating to the item that was shared. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
func (m *SharedInsight) GetResource()(Entityable) {
    return m.resource
}
// GetResourceReference gets the resourceReference property value. Reference properties of the shared document, such as the url and type of the document. Read-only
func (m *SharedInsight) GetResourceReference()(ResourceReferenceable) {
    return m.resourceReference
}
// GetResourceVisualization gets the resourceVisualization property value. Properties that you can use to visualize the document in your experience. Read-only
func (m *SharedInsight) GetResourceVisualization()(ResourceVisualizationable) {
    return m.resourceVisualization
}
// GetSharingHistory gets the sharingHistory property value. The sharingHistory property
func (m *SharedInsight) GetSharingHistory()([]SharingDetailable) {
    return m.sharingHistory
}
// Serialize serializes information the current object
func (m *SharedInsight) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("lastShared", m.GetLastShared())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastSharedMethod", m.GetLastSharedMethod())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("resource", m.GetResource())
        if err != nil {
            return err
        }
    }
    if m.GetSharingHistory() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSharingHistory())
        err = writer.WriteCollectionOfObjectValues("sharingHistory", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLastShared sets the lastShared property value. Details about the shared item. Read only.
func (m *SharedInsight) SetLastShared(value SharingDetailable)() {
    m.lastShared = value
}
// SetLastSharedMethod sets the lastSharedMethod property value. The lastSharedMethod property
func (m *SharedInsight) SetLastSharedMethod(value Entityable)() {
    m.lastSharedMethod = value
}
// SetResource sets the resource property value. Used for navigating to the item that was shared. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
func (m *SharedInsight) SetResource(value Entityable)() {
    m.resource = value
}
// SetResourceReference sets the resourceReference property value. Reference properties of the shared document, such as the url and type of the document. Read-only
func (m *SharedInsight) SetResourceReference(value ResourceReferenceable)() {
    m.resourceReference = value
}
// SetResourceVisualization sets the resourceVisualization property value. Properties that you can use to visualize the document in your experience. Read-only
func (m *SharedInsight) SetResourceVisualization(value ResourceVisualizationable)() {
    m.resourceVisualization = value
}
// SetSharingHistory sets the sharingHistory property value. The sharingHistory property
func (m *SharedInsight) SetSharingHistory(value []SharingDetailable)() {
    m.sharingHistory = value
}
