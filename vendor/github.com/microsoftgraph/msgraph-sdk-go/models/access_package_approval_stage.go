package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AccessPackageApprovalStage 
type AccessPackageApprovalStage struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The number of days that a request can be pending a response before it is automatically denied.
    durationBeforeAutomaticDenial *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration
    // If escalation is required, the time a request can be pending a response from a primary approver.
    durationBeforeEscalation *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration
    // If escalation is enabled and the primary approvers do not respond before the escalation time, the escalationApprovers are the users who will be asked to approve requests.
    escalationApprovers []SubjectSetable
    // The subjects, typically users, who are the fallback escalation approvers.
    fallbackEscalationApprovers []SubjectSetable
    // The subjects, typically users, who are the fallback primary approvers.
    fallbackPrimaryApprovers []SubjectSetable
    // Indicates whether the approver is required to provide a justification for approving a request.
    isApproverJustificationRequired *bool
    // If true, then one or more escalationApprovers are configured in this approval stage.
    isEscalationEnabled *bool
    // The OdataType property
    odataType *string
    // The subjects, typically users, who will be asked to approve requests. A collection of singleUser, groupMembers, requestorManager, internalSponsors or externalSponsors.
    primaryApprovers []SubjectSetable
}
// NewAccessPackageApprovalStage instantiates a new accessPackageApprovalStage and sets the default values.
func NewAccessPackageApprovalStage()(*AccessPackageApprovalStage) {
    m := &AccessPackageApprovalStage{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAccessPackageApprovalStageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessPackageApprovalStageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageApprovalStage(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AccessPackageApprovalStage) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDurationBeforeAutomaticDenial gets the durationBeforeAutomaticDenial property value. The number of days that a request can be pending a response before it is automatically denied.
func (m *AccessPackageApprovalStage) GetDurationBeforeAutomaticDenial()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    return m.durationBeforeAutomaticDenial
}
// GetDurationBeforeEscalation gets the durationBeforeEscalation property value. If escalation is required, the time a request can be pending a response from a primary approver.
func (m *AccessPackageApprovalStage) GetDurationBeforeEscalation()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    return m.durationBeforeEscalation
}
// GetEscalationApprovers gets the escalationApprovers property value. If escalation is enabled and the primary approvers do not respond before the escalation time, the escalationApprovers are the users who will be asked to approve requests.
func (m *AccessPackageApprovalStage) GetEscalationApprovers()([]SubjectSetable) {
    return m.escalationApprovers
}
// GetFallbackEscalationApprovers gets the fallbackEscalationApprovers property value. The subjects, typically users, who are the fallback escalation approvers.
func (m *AccessPackageApprovalStage) GetFallbackEscalationApprovers()([]SubjectSetable) {
    return m.fallbackEscalationApprovers
}
// GetFallbackPrimaryApprovers gets the fallbackPrimaryApprovers property value. The subjects, typically users, who are the fallback primary approvers.
func (m *AccessPackageApprovalStage) GetFallbackPrimaryApprovers()([]SubjectSetable) {
    return m.fallbackPrimaryApprovers
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessPackageApprovalStage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["durationBeforeAutomaticDenial"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetISODurationValue(m.SetDurationBeforeAutomaticDenial)
    res["durationBeforeEscalation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetISODurationValue(m.SetDurationBeforeEscalation)
    res["escalationApprovers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetEscalationApprovers)
    res["fallbackEscalationApprovers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetFallbackEscalationApprovers)
    res["fallbackPrimaryApprovers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetFallbackPrimaryApprovers)
    res["isApproverJustificationRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsApproverJustificationRequired)
    res["isEscalationEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsEscalationEnabled)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["primaryApprovers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetPrimaryApprovers)
    return res
}
// GetIsApproverJustificationRequired gets the isApproverJustificationRequired property value. Indicates whether the approver is required to provide a justification for approving a request.
func (m *AccessPackageApprovalStage) GetIsApproverJustificationRequired()(*bool) {
    return m.isApproverJustificationRequired
}
// GetIsEscalationEnabled gets the isEscalationEnabled property value. If true, then one or more escalationApprovers are configured in this approval stage.
func (m *AccessPackageApprovalStage) GetIsEscalationEnabled()(*bool) {
    return m.isEscalationEnabled
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AccessPackageApprovalStage) GetOdataType()(*string) {
    return m.odataType
}
// GetPrimaryApprovers gets the primaryApprovers property value. The subjects, typically users, who will be asked to approve requests. A collection of singleUser, groupMembers, requestorManager, internalSponsors or externalSponsors.
func (m *AccessPackageApprovalStage) GetPrimaryApprovers()([]SubjectSetable) {
    return m.primaryApprovers
}
// Serialize serializes information the current object
func (m *AccessPackageApprovalStage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteISODurationValue("durationBeforeAutomaticDenial", m.GetDurationBeforeAutomaticDenial())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("durationBeforeEscalation", m.GetDurationBeforeEscalation())
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
    if m.GetFallbackEscalationApprovers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetFallbackEscalationApprovers())
        err := writer.WriteCollectionOfObjectValues("fallbackEscalationApprovers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetFallbackPrimaryApprovers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetFallbackPrimaryApprovers())
        err := writer.WriteCollectionOfObjectValues("fallbackPrimaryApprovers", cast)
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
func (m *AccessPackageApprovalStage) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDurationBeforeAutomaticDenial sets the durationBeforeAutomaticDenial property value. The number of days that a request can be pending a response before it is automatically denied.
func (m *AccessPackageApprovalStage) SetDurationBeforeAutomaticDenial(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    m.durationBeforeAutomaticDenial = value
}
// SetDurationBeforeEscalation sets the durationBeforeEscalation property value. If escalation is required, the time a request can be pending a response from a primary approver.
func (m *AccessPackageApprovalStage) SetDurationBeforeEscalation(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    m.durationBeforeEscalation = value
}
// SetEscalationApprovers sets the escalationApprovers property value. If escalation is enabled and the primary approvers do not respond before the escalation time, the escalationApprovers are the users who will be asked to approve requests.
func (m *AccessPackageApprovalStage) SetEscalationApprovers(value []SubjectSetable)() {
    m.escalationApprovers = value
}
// SetFallbackEscalationApprovers sets the fallbackEscalationApprovers property value. The subjects, typically users, who are the fallback escalation approvers.
func (m *AccessPackageApprovalStage) SetFallbackEscalationApprovers(value []SubjectSetable)() {
    m.fallbackEscalationApprovers = value
}
// SetFallbackPrimaryApprovers sets the fallbackPrimaryApprovers property value. The subjects, typically users, who are the fallback primary approvers.
func (m *AccessPackageApprovalStage) SetFallbackPrimaryApprovers(value []SubjectSetable)() {
    m.fallbackPrimaryApprovers = value
}
// SetIsApproverJustificationRequired sets the isApproverJustificationRequired property value. Indicates whether the approver is required to provide a justification for approving a request.
func (m *AccessPackageApprovalStage) SetIsApproverJustificationRequired(value *bool)() {
    m.isApproverJustificationRequired = value
}
// SetIsEscalationEnabled sets the isEscalationEnabled property value. If true, then one or more escalationApprovers are configured in this approval stage.
func (m *AccessPackageApprovalStage) SetIsEscalationEnabled(value *bool)() {
    m.isEscalationEnabled = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessPackageApprovalStage) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPrimaryApprovers sets the primaryApprovers property value. The subjects, typically users, who will be asked to approve requests. A collection of singleUser, groupMembers, requestorManager, internalSponsors or externalSponsors.
func (m *AccessPackageApprovalStage) SetPrimaryApprovers(value []SubjectSetable)() {
    m.primaryApprovers = value
}
