package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AccessPackageAssignmentReviewSettings 
type AccessPackageAssignmentReviewSettings struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The default decision to apply if the access is not reviewed. The possible values are: keepAccess, removeAccess, acceptAccessRecommendation, unknownFutureValue.
    expirationBehavior *AccessReviewExpirationBehavior
    // This collection specifies the users who will be the fallback reviewers when the primary reviewers don't respond.
    fallbackReviewers []SubjectSetable
    // If true, access reviews are required for assignments through this policy.
    isEnabled *bool
    // Specifies whether to display recommendations to the reviewer. The default value is true.
    isRecommendationEnabled *bool
    // Specifies whether the reviewer must provide justification for the approval. The default value is true.
    isReviewerJustificationRequired *bool
    // Specifies whether the principals can review their own assignments.
    isSelfReview *bool
    // The OdataType property
    odataType *string
    // This collection specifies the users or group of users who will review the access package assignments.
    primaryReviewers []SubjectSetable
    // When the first review should start and how often it should recur.
    schedule EntitlementManagementScheduleable
}
// NewAccessPackageAssignmentReviewSettings instantiates a new accessPackageAssignmentReviewSettings and sets the default values.
func NewAccessPackageAssignmentReviewSettings()(*AccessPackageAssignmentReviewSettings) {
    m := &AccessPackageAssignmentReviewSettings{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAccessPackageAssignmentReviewSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessPackageAssignmentReviewSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignmentReviewSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AccessPackageAssignmentReviewSettings) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetExpirationBehavior gets the expirationBehavior property value. The default decision to apply if the access is not reviewed. The possible values are: keepAccess, removeAccess, acceptAccessRecommendation, unknownFutureValue.
func (m *AccessPackageAssignmentReviewSettings) GetExpirationBehavior()(*AccessReviewExpirationBehavior) {
    return m.expirationBehavior
}
// GetFallbackReviewers gets the fallbackReviewers property value. This collection specifies the users who will be the fallback reviewers when the primary reviewers don't respond.
func (m *AccessPackageAssignmentReviewSettings) GetFallbackReviewers()([]SubjectSetable) {
    return m.fallbackReviewers
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessPackageAssignmentReviewSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["expirationBehavior"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseAccessReviewExpirationBehavior , m.SetExpirationBehavior)
    res["fallbackReviewers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetFallbackReviewers)
    res["isEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsEnabled)
    res["isRecommendationEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsRecommendationEnabled)
    res["isReviewerJustificationRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsReviewerJustificationRequired)
    res["isSelfReview"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsSelfReview)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["primaryReviewers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue , m.SetPrimaryReviewers)
    res["schedule"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEntitlementManagementScheduleFromDiscriminatorValue , m.SetSchedule)
    return res
}
// GetIsEnabled gets the isEnabled property value. If true, access reviews are required for assignments through this policy.
func (m *AccessPackageAssignmentReviewSettings) GetIsEnabled()(*bool) {
    return m.isEnabled
}
// GetIsRecommendationEnabled gets the isRecommendationEnabled property value. Specifies whether to display recommendations to the reviewer. The default value is true.
func (m *AccessPackageAssignmentReviewSettings) GetIsRecommendationEnabled()(*bool) {
    return m.isRecommendationEnabled
}
// GetIsReviewerJustificationRequired gets the isReviewerJustificationRequired property value. Specifies whether the reviewer must provide justification for the approval. The default value is true.
func (m *AccessPackageAssignmentReviewSettings) GetIsReviewerJustificationRequired()(*bool) {
    return m.isReviewerJustificationRequired
}
// GetIsSelfReview gets the isSelfReview property value. Specifies whether the principals can review their own assignments.
func (m *AccessPackageAssignmentReviewSettings) GetIsSelfReview()(*bool) {
    return m.isSelfReview
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AccessPackageAssignmentReviewSettings) GetOdataType()(*string) {
    return m.odataType
}
// GetPrimaryReviewers gets the primaryReviewers property value. This collection specifies the users or group of users who will review the access package assignments.
func (m *AccessPackageAssignmentReviewSettings) GetPrimaryReviewers()([]SubjectSetable) {
    return m.primaryReviewers
}
// GetSchedule gets the schedule property value. When the first review should start and how often it should recur.
func (m *AccessPackageAssignmentReviewSettings) GetSchedule()(EntitlementManagementScheduleable) {
    return m.schedule
}
// Serialize serializes information the current object
func (m *AccessPackageAssignmentReviewSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetExpirationBehavior() != nil {
        cast := (*m.GetExpirationBehavior()).String()
        err := writer.WriteStringValue("expirationBehavior", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetFallbackReviewers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetFallbackReviewers())
        err := writer.WriteCollectionOfObjectValues("fallbackReviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isRecommendationEnabled", m.GetIsRecommendationEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isReviewerJustificationRequired", m.GetIsReviewerJustificationRequired())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isSelfReview", m.GetIsSelfReview())
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
    if m.GetPrimaryReviewers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPrimaryReviewers())
        err := writer.WriteCollectionOfObjectValues("primaryReviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("schedule", m.GetSchedule())
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
func (m *AccessPackageAssignmentReviewSettings) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetExpirationBehavior sets the expirationBehavior property value. The default decision to apply if the access is not reviewed. The possible values are: keepAccess, removeAccess, acceptAccessRecommendation, unknownFutureValue.
func (m *AccessPackageAssignmentReviewSettings) SetExpirationBehavior(value *AccessReviewExpirationBehavior)() {
    m.expirationBehavior = value
}
// SetFallbackReviewers sets the fallbackReviewers property value. This collection specifies the users who will be the fallback reviewers when the primary reviewers don't respond.
func (m *AccessPackageAssignmentReviewSettings) SetFallbackReviewers(value []SubjectSetable)() {
    m.fallbackReviewers = value
}
// SetIsEnabled sets the isEnabled property value. If true, access reviews are required for assignments through this policy.
func (m *AccessPackageAssignmentReviewSettings) SetIsEnabled(value *bool)() {
    m.isEnabled = value
}
// SetIsRecommendationEnabled sets the isRecommendationEnabled property value. Specifies whether to display recommendations to the reviewer. The default value is true.
func (m *AccessPackageAssignmentReviewSettings) SetIsRecommendationEnabled(value *bool)() {
    m.isRecommendationEnabled = value
}
// SetIsReviewerJustificationRequired sets the isReviewerJustificationRequired property value. Specifies whether the reviewer must provide justification for the approval. The default value is true.
func (m *AccessPackageAssignmentReviewSettings) SetIsReviewerJustificationRequired(value *bool)() {
    m.isReviewerJustificationRequired = value
}
// SetIsSelfReview sets the isSelfReview property value. Specifies whether the principals can review their own assignments.
func (m *AccessPackageAssignmentReviewSettings) SetIsSelfReview(value *bool)() {
    m.isSelfReview = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessPackageAssignmentReviewSettings) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPrimaryReviewers sets the primaryReviewers property value. This collection specifies the users or group of users who will review the access package assignments.
func (m *AccessPackageAssignmentReviewSettings) SetPrimaryReviewers(value []SubjectSetable)() {
    m.primaryReviewers = value
}
// SetSchedule sets the schedule property value. When the first review should start and how often it should recur.
func (m *AccessPackageAssignmentReviewSettings) SetSchedule(value EntitlementManagementScheduleable)() {
    m.schedule = value
}
