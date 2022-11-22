package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UnifiedApprovalStage 
type UnifiedApprovalStage struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The number of days that a request can be pending a response before it is automatically denied.
    approvalStageTimeOutInDays *int32
    // The escalation approvers for this stage when the primary approvers don't respond.
    escalationApprovers []SubjectSetable
    // The time a request can be pending a response from a primary approver before it can be escalated to the escalation approvers.
    escalationTimeInMinutes *int32
    // Indicates whether the approver must provide justification for their reponse.
    isApproverJustificationRequired *bool
    // Indicates whether escalation if enabled.
    isEscalationEnabled *bool
    // The OdataType property
    odataType *string
    // The primary approvers of this stage.
    primaryApprovers []SubjectSetable
}
// NewUnifiedApprovalStage instantiates a new unifiedApprovalStage and sets the default values.
func NewUnifiedApprovalStage()(*UnifiedApprovalStage) {
    m := &UnifiedApprovalStage{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateUnifiedApprovalStageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUnifiedApprovalStageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedApprovalStage(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UnifiedApprovalStage) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetApprovalStageTimeOutInDays gets the approvalStageTimeOutInDays property value. The number of days that a request can be pending a response before it is automatically denied.
func (m *UnifiedApprovalStage) GetApprovalStageTimeOutInDays()(*int32) {
    return m.approvalStageTimeOutInDays
}
// GetEscalationApprovers gets the escalationApprovers property value. The escalation approvers for this stage when the primary approvers don't respond.
func (m *UnifiedApprovalStage) GetEscalationApprovers()([]SubjectSetable) {
    return m.escalationApprovers
}
// GetEscalationTimeInMinutes gets the escalationTimeInMinutes property value. The time a request can be pending a response from a primary approver before it can be escalated to the escalation approvers.
func (m *UnifiedApprovalStage) GetEscalationTimeInMinutes()(*int32) {
    return m.escalationTimeInMinutes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UnifiedApprovalStage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["approvalStageTimeOutInDays"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetApprovalStageTimeOutInDays)
    res["escalationApprovers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetEscalationApprovers)
    res["escalationTimeInMinutes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetEscalationTimeInMinutes)
    res["isApproverJustificationRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsApproverJustificationRequired)
    res["isEscalationEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsEscalationEnabled)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["primaryApprovers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetPrimaryApprovers)
    return res
}
// GetIsApproverJustificationRequired gets the isApproverJustificationRequired property value. Indicates whether the approver must provide justification for their reponse.
func (m *UnifiedApprovalStage) GetIsApproverJustificationRequired()(*bool) {
    return m.isApproverJustificationRequired
}
// GetIsEscalationEnabled gets the isEscalationEnabled property value. Indicates whether escalation if enabled.
func (m *UnifiedApprovalStage) GetIsEscalationEnabled()(*bool) {
    return m.isEscalationEnabled
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *UnifiedApprovalStage) GetOdataType()(*string) {
    return m.odataType
}
// GetPrimaryApprovers gets the primaryApprovers property value. The primary approvers of this stage.
func (m *UnifiedApprovalStage) GetPrimaryApprovers()([]SubjectSetable) {
    return m.primaryApprovers
}
// Serialize serializes information the current object
func (m *UnifiedApprovalStage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("approvalStageTimeOutInDays", m.GetApprovalStageTimeOutInDays())
        if err != nil {
            return err
        }
    }
    if m.GetEscalationApprovers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEscalationApprovers())
        err := writer.WriteCollectionOfObjectValues("escalationApprovers", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("escalationTimeInMinutes", m.GetEscalationTimeInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isApproverJustificationRequired", m.GetIsApproverJustificationRequired())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isEscalationEnabled", m.GetIsEscalationEnabled())
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
    if m.GetPrimaryApprovers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPrimaryApprovers())
        err := writer.WriteCollectionOfObjectValues("primaryApprovers", cast)
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
func (m *UnifiedApprovalStage) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetApprovalStageTimeOutInDays sets the approvalStageTimeOutInDays property value. The number of days that a request can be pending a response before it is automatically denied.
func (m *UnifiedApprovalStage) SetApprovalStageTimeOutInDays(value *int32)() {
    m.approvalStageTimeOutInDays = value
}
// SetEscalationApprovers sets the escalationApprovers property value. The escalation approvers for this stage when the primary approvers don't respond.
func (m *UnifiedApprovalStage) SetEscalationApprovers(value []SubjectSetable)() {
    m.escalationApprovers = value
}
// SetEscalationTimeInMinutes sets the escalationTimeInMinutes property value. The time a request can be pending a response from a primary approver before it can be escalated to the escalation approvers.
func (m *UnifiedApprovalStage) SetEscalationTimeInMinutes(value *int32)() {
    m.escalationTimeInMinutes = value
}
// SetIsApproverJustificationRequired sets the isApproverJustificationRequired property value. Indicates whether the approver must provide justification for their reponse.
func (m *UnifiedApprovalStage) SetIsApproverJustificationRequired(value *bool)() {
    m.isApproverJustificationRequired = value
}
// SetIsEscalationEnabled sets the isEscalationEnabled property value. Indicates whether escalation if enabled.
func (m *UnifiedApprovalStage) SetIsEscalationEnabled(value *bool)() {
    m.isEscalationEnabled = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UnifiedApprovalStage) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPrimaryApprovers sets the primaryApprovers property value. The primary approvers of this stage.
func (m *UnifiedApprovalStage) SetPrimaryApprovers(value []SubjectSetable)() {
    m.primaryApprovers = value
}
