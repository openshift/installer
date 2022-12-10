package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DocumentSet 
type DocumentSet struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Content types allowed in document set.
    allowedContentTypes []ContentTypeInfoable
    // Default contents of document set.
    defaultContents []DocumentSetContentable
    // The OdataType property
    odataType *string
    // Specifies whether to push welcome page changes to inherited content types.
    propagateWelcomePageChanges *bool
    // The sharedColumns property
    sharedColumns []ColumnDefinitionable
    // Indicates whether to add the name of the document set to each file name.
    shouldPrefixNameToFile *bool
    // The welcomePageColumns property
    welcomePageColumns []ColumnDefinitionable
    // Welcome page absolute URL.
    welcomePageUrl *string
}
// NewDocumentSet instantiates a new documentSet and sets the default values.
func NewDocumentSet()(*DocumentSet) {
    m := &DocumentSet{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateDocumentSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDocumentSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDocumentSet(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *DocumentSet) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAllowedContentTypes gets the allowedContentTypes property value. Content types allowed in document set.
func (m *DocumentSet) GetAllowedContentTypes()([]ContentTypeInfoable) {
    return m.allowedContentTypes
}
// GetDefaultContents gets the defaultContents property value. Default contents of document set.
func (m *DocumentSet) GetDefaultContents()([]DocumentSetContentable) {
    return m.defaultContents
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DocumentSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowedContentTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateContentTypeInfoFromDiscriminatorValue , m.SetAllowedContentTypes)
    res["defaultContents"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDocumentSetContentFromDiscriminatorValue , m.SetDefaultContents)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["propagateWelcomePageChanges"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetPropagateWelcomePageChanges)
    res["sharedColumns"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateColumnDefinitionFromDiscriminatorValue , m.SetSharedColumns)
    res["shouldPrefixNameToFile"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetShouldPrefixNameToFile)
    res["welcomePageColumns"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateColumnDefinitionFromDiscriminatorValue , m.SetWelcomePageColumns)
    res["welcomePageUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetWelcomePageUrl)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *DocumentSet) GetOdataType()(*string) {
    return m.odataType
}
// GetPropagateWelcomePageChanges gets the propagateWelcomePageChanges property value. Specifies whether to push welcome page changes to inherited content types.
func (m *DocumentSet) GetPropagateWelcomePageChanges()(*bool) {
    return m.propagateWelcomePageChanges
}
// GetSharedColumns gets the sharedColumns property value. The sharedColumns property
func (m *DocumentSet) GetSharedColumns()([]ColumnDefinitionable) {
    return m.sharedColumns
}
// GetShouldPrefixNameToFile gets the shouldPrefixNameToFile property value. Indicates whether to add the name of the document set to each file name.
func (m *DocumentSet) GetShouldPrefixNameToFile()(*bool) {
    return m.shouldPrefixNameToFile
}
// GetWelcomePageColumns gets the welcomePageColumns property value. The welcomePageColumns property
func (m *DocumentSet) GetWelcomePageColumns()([]ColumnDefinitionable) {
    return m.welcomePageColumns
}
// GetWelcomePageUrl gets the welcomePageUrl property value. Welcome page absolute URL.
func (m *DocumentSet) GetWelcomePageUrl()(*string) {
    return m.welcomePageUrl
}
// Serialize serializes information the current object
func (m *DocumentSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAllowedContentTypes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAllowedContentTypes())
        err := writer.WriteCollectionOfObjectValues("allowedContentTypes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDefaultContents() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDefaultContents())
        err := writer.WriteCollectionOfObjectValues("defaultContents", cast)
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
        err := writer.WriteBoolValue("propagateWelcomePageChanges", m.GetPropagateWelcomePageChanges())
        if err != nil {
            return err
        }
    }
    if m.GetSharedColumns() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSharedColumns())
        err := writer.WriteCollectionOfObjectValues("sharedColumns", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("shouldPrefixNameToFile", m.GetShouldPrefixNameToFile())
        if err != nil {
            return err
        }
    }
    if m.GetWelcomePageColumns() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetWelcomePageColumns())
        err := writer.WriteCollectionOfObjectValues("welcomePageColumns", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("welcomePageUrl", m.GetWelcomePageUrl())
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
func (m *DocumentSet) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAllowedContentTypes sets the allowedContentTypes property value. Content types allowed in document set.
func (m *DocumentSet) SetAllowedContentTypes(value []ContentTypeInfoable)() {
    m.allowedContentTypes = value
}
// SetDefaultContents sets the defaultContents property value. Default contents of document set.
func (m *DocumentSet) SetDefaultContents(value []DocumentSetContentable)() {
    m.defaultContents = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DocumentSet) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPropagateWelcomePageChanges sets the propagateWelcomePageChanges property value. Specifies whether to push welcome page changes to inherited content types.
func (m *DocumentSet) SetPropagateWelcomePageChanges(value *bool)() {
    m.propagateWelcomePageChanges = value
}
// SetSharedColumns sets the sharedColumns property value. The sharedColumns property
func (m *DocumentSet) SetSharedColumns(value []ColumnDefinitionable)() {
    m.sharedColumns = value
}
// SetShouldPrefixNameToFile sets the shouldPrefixNameToFile property value. Indicates whether to add the name of the document set to each file name.
func (m *DocumentSet) SetShouldPrefixNameToFile(value *bool)() {
    m.shouldPrefixNameToFile = value
}
// SetWelcomePageColumns sets the welcomePageColumns property value. The welcomePageColumns property
func (m *DocumentSet) SetWelcomePageColumns(value []ColumnDefinitionable)() {
    m.welcomePageColumns = value
}
// SetWelcomePageUrl sets the welcomePageUrl property value. Welcome page absolute URL.
func (m *DocumentSet) SetWelcomePageUrl(value *string)() {
    m.welcomePageUrl = value
}
