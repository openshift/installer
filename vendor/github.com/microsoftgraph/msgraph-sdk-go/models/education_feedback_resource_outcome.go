package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EducationFeedbackResourceOutcome 
type EducationFeedbackResourceOutcome struct {
    EducationOutcome
    // The actual feedback resource.
    feedbackResource EducationResourceable
    // The status of the feedback resource. The possible values are: notPublished, pendingPublish, published, failedPublish, unknownFutureValue.
    resourceStatus *EducationFeedbackResourceOutcomeStatus
}
// NewEducationFeedbackResourceOutcome instantiates a new EducationFeedbackResourceOutcome and sets the default values.
func NewEducationFeedbackResourceOutcome()(*EducationFeedbackResourceOutcome) {
    m := &EducationFeedbackResourceOutcome{
        EducationOutcome: *NewEducationOutcome(),
    }
    odataTypeValue := "#microsoft.graph.educationFeedbackResourceOutcome";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEducationFeedbackResourceOutcomeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEducationFeedbackResourceOutcomeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationFeedbackResourceOutcome(), nil
}
// GetFeedbackResource gets the feedbackResource property value. The actual feedback resource.
func (m *EducationFeedbackResourceOutcome) GetFeedbackResource()(EducationResourceable) {
    return m.feedbackResource
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EducationFeedbackResourceOutcome) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationOutcome.GetFieldDeserializers()
    res["feedbackResource"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEducationResourceFromDiscriminatorValue , m.SetFeedbackResource)
    res["resourceStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseEducationFeedbackResourceOutcomeStatus , m.SetResourceStatus)
    return res
}
// GetResourceStatus gets the resourceStatus property value. The status of the feedback resource. The possible values are: notPublished, pendingPublish, published, failedPublish, unknownFutureValue.
func (m *EducationFeedbackResourceOutcome) GetResourceStatus()(*EducationFeedbackResourceOutcomeStatus) {
    return m.resourceStatus
}
// Serialize serializes information the current object
func (m *EducationFeedbackResourceOutcome) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationOutcome.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("feedbackResource", m.GetFeedbackResource())
        if err != nil {
            return err
        }
    }
    if m.GetResourceStatus() != nil {
        cast := (*m.GetResourceStatus()).String()
        err = writer.WriteStringValue("resourceStatus", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFeedbackResource sets the feedbackResource property value. The actual feedback resource.
func (m *EducationFeedbackResourceOutcome) SetFeedbackResource(value EducationResourceable)() {
    m.feedbackResource = value
}
// SetResourceStatus sets the resourceStatus property value. The status of the feedback resource. The possible values are: notPublished, pendingPublish, published, failedPublish, unknownFutureValue.
func (m *EducationFeedbackResourceOutcome) SetResourceStatus(value *EducationFeedbackResourceOutcomeStatus)() {
    m.resourceStatus = value
}
