package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DocumentSetVersion 
type DocumentSetVersion struct {
    ListItemVersion
    // Comment about the captured version.
    comment *string
    // User who captured the version.
    createdBy IdentitySetable
    // Date and time when this version was created.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Items within the document set that are captured as part of this version.
    items []DocumentSetVersionItemable
    // If true, minor versions of items are also captured; otherwise, only major versions will be captured. Default value is false.
    shouldCaptureMinorVersion *bool
}
// NewDocumentSetVersion instantiates a new DocumentSetVersion and sets the default values.
func NewDocumentSetVersion()(*DocumentSetVersion) {
    m := &DocumentSetVersion{
        ListItemVersion: *NewListItemVersion(),
    }
    odataTypeValue := "#microsoft.graph.documentSetVersion";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateDocumentSetVersionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDocumentSetVersionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDocumentSetVersion(), nil
}
// GetComment gets the comment property value. Comment about the captured version.
func (m *DocumentSetVersion) GetComment()(*string) {
    return m.comment
}
// GetCreatedBy gets the createdBy property value. User who captured the version.
func (m *DocumentSetVersion) GetCreatedBy()(IdentitySetable) {
    return m.createdBy
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time when this version was created.
func (m *DocumentSetVersion) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DocumentSetVersion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ListItemVersion.GetFieldDeserializers()
    res["comment"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetComment)
    res["createdBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetCreatedBy)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["items"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDocumentSetVersionItemFromDiscriminatorValue , m.SetItems)
    res["shouldCaptureMinorVersion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetShouldCaptureMinorVersion)
    return res
}
// GetItems gets the items property value. Items within the document set that are captured as part of this version.
func (m *DocumentSetVersion) GetItems()([]DocumentSetVersionItemable) {
    return m.items
}
// GetShouldCaptureMinorVersion gets the shouldCaptureMinorVersion property value. If true, minor versions of items are also captured; otherwise, only major versions will be captured. Default value is false.
func (m *DocumentSetVersion) GetShouldCaptureMinorVersion()(*bool) {
    return m.shouldCaptureMinorVersion
}
// Serialize serializes information the current object
func (m *DocumentSetVersion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ListItemVersion.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("comment", m.GetComment())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
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
        err = writer.WriteBoolValue("shouldCaptureMinorVersion", m.GetShouldCaptureMinorVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetComment sets the comment property value. Comment about the captured version.
func (m *DocumentSetVersion) SetComment(value *string)() {
    m.comment = value
}
// SetCreatedBy sets the createdBy property value. User who captured the version.
func (m *DocumentSetVersion) SetCreatedBy(value IdentitySetable)() {
    m.createdBy = value
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time when this version was created.
func (m *DocumentSetVersion) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetItems sets the items property value. Items within the document set that are captured as part of this version.
func (m *DocumentSetVersion) SetItems(value []DocumentSetVersionItemable)() {
    m.items = value
}
// SetShouldCaptureMinorVersion sets the shouldCaptureMinorVersion property value. If true, minor versions of items are also captured; otherwise, only major versions will be captured. Default value is false.
func (m *DocumentSetVersion) SetShouldCaptureMinorVersion(value *bool)() {
    m.shouldCaptureMinorVersion = value
}
