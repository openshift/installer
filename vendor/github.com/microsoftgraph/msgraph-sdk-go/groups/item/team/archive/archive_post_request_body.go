package archive

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ArchivePostRequestBody provides operations to call the archive method.
type ArchivePostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The shouldSetSpoSiteReadOnlyForMembers property
    shouldSetSpoSiteReadOnlyForMembers *bool
}
// NewArchivePostRequestBody instantiates a new archivePostRequestBody and sets the default values.
func NewArchivePostRequestBody()(*ArchivePostRequestBody) {
    m := &ArchivePostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateArchivePostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateArchivePostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewArchivePostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ArchivePostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ArchivePostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["shouldSetSpoSiteReadOnlyForMembers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetShouldSetSpoSiteReadOnlyForMembers)
    return res
}
// GetShouldSetSpoSiteReadOnlyForMembers gets the shouldSetSpoSiteReadOnlyForMembers property value. The shouldSetSpoSiteReadOnlyForMembers property
func (m *ArchivePostRequestBody) GetShouldSetSpoSiteReadOnlyForMembers()(*bool) {
    return m.shouldSetSpoSiteReadOnlyForMembers
}
// Serialize serializes information the current object
func (m *ArchivePostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("shouldSetSpoSiteReadOnlyForMembers", m.GetShouldSetSpoSiteReadOnlyForMembers())
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
func (m *ArchivePostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetShouldSetSpoSiteReadOnlyForMembers sets the shouldSetSpoSiteReadOnlyForMembers property value. The shouldSetSpoSiteReadOnlyForMembers property
func (m *ArchivePostRequestBody) SetShouldSetSpoSiteReadOnlyForMembers(value *bool)() {
    m.shouldSetSpoSiteReadOnlyForMembers = value
}
