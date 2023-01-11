package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ApprovalSettings 
type ApprovalSettings struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // One of SingleStage, Serial, Parallel, NoApproval (default). NoApproval is used when isApprovalRequired is false.
    approvalMode *string
    // If approval is required, the one or two elements of this collection define each of the stages of approval. An empty array if no approval is required.
    approvalStages []UnifiedApprovalStageable
    // Indicates whether approval is required for requests in this policy.
    isApprovalRequired *bool
    // Indicates whether approval is required for a user to extend their assignment.
    isApprovalRequiredForExtension *bool
    // Indicates whether the requestor is required to supply a justification in their request.
    isRequestorJustificationRequired *bool
    // The OdataType property
    odataType *string
}
// NewApprovalSettings instantiates a new approvalSettings and sets the default values.
func NewApprovalSettings()(*ApprovalSettings) {
    m := &ApprovalSettings{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateApprovalSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateApprovalSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewApprovalSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ApprovalSettings) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetApprovalMode gets the approvalMode property value. One of SingleStage, Serial, Parallel, NoApproval (default). NoApproval is used when isApprovalRequired is false.
func (m *ApprovalSettings) GetApprovalMode()(*string) {
    return m.approvalMode
}
// GetApprovalStages gets the approvalStages property value. If approval is required, the one or two elements of this collection define each of the stages of approval. An empty array if no approval is required.
func (m *ApprovalSettings) GetApprovalStages()([]UnifiedApprovalStageable) {
    return m.approvalStages
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ApprovalSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["approvalMode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetApprovalMode)
    res["approvalStages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateUnifiedApprovalStageFromDiscriminatorValue , m.SetApprovalStages)
    res["isApprovalRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsApprovalRequired)
    res["isApprovalRequiredForExtension"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsApprovalRequiredForExtension)
    res["isRequestorJustificationRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsRequestorJustificationRequired)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetIsApprovalRequired gets the isApprovalRequired property value. Indicates whether approval is required for requests in this policy.
func (m *ApprovalSettings) GetIsApprovalRequired()(*bool) {
    return m.isApprovalRequired
}
// GetIsApprovalRequiredForExtension gets the isApprovalRequiredForExtension property value. Indicates whether approval is required for a user to extend their assignment.
func (m *ApprovalSettings) GetIsApprovalRequiredForExtension()(*bool) {
    return m.isApprovalRequiredForExtension
}
// GetIsRequestorJustificationRequired gets the isRequestorJustificationRequired property value. Indicates whether the requestor is required to supply a justification in their request.
func (m *ApprovalSettings) GetIsRequestorJustificationRequired()(*bool) {
    return m.isRequestorJustificationRequired
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ApprovalSettings) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ApprovalSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("approvalMode", m.GetApprovalMode())
        if err != nil {
            return err
        }
    }
    if m.GetApprovalStages() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetApprovalStages())
        err := writer.WriteCollectionOfObjectValues("approvalStages", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isApprovalRequired", m.GetIsApprovalRequired())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isApprovalRequiredForExtension", m.GetIsApprovalRequiredForExtension())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isRequestorJustificationRequired", m.GetIsRequestorJustificationRequired())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ApprovalSettings) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetApprovalMode sets the approvalMode property value. One of SingleStage, Serial, Parallel, NoApproval (default). NoApproval is used when isApprovalRequired is false.
func (m *ApprovalSettings) SetApprovalMode(value *string)() {
    m.approvalMode = value
}
// SetApprovalStages sets the approvalStages property value. If approval is required, the one or two elements of this collection define each of the stages of approval. An empty array if no approval is required.
func (m *ApprovalSettings) SetApprovalStages(value []UnifiedApprovalStageable)() {
    m.approvalStages = value
}
// SetIsApprovalRequired sets the isApprovalRequired property value. Indicates whether approval is required for requests in this policy.
func (m *ApprovalSettings) SetIsApprovalRequired(value *bool)() {
    m.isApprovalRequired = value
}
// SetIsApprovalRequiredForExtension sets the isApprovalRequiredForExtension property value. Indicates whether approval is required for a user to extend their assignment.
func (m *ApprovalSettings) SetIsApprovalRequiredForExtension(value *bool)() {
    m.isApprovalRequiredForExtension = value
}
// SetIsRequestorJustificationRequired sets the isRequestorJustificationRequired property value. Indicates whether the requestor is required to supply a justification in their request.
func (m *ApprovalSettings) SetIsRequestorJustificationRequired(value *bool)() {
    m.isRequestorJustificationRequired = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ApprovalSettings) SetOdataType(value *string)() {
    m.odataType = value
}
