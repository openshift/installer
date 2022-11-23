package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RubricQualityFeedbackModel 
type RubricQualityFeedbackModel struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Specific feedback for one quality of this rubric.
    feedback EducationItemBodyable
    // The OdataType property
    odataType *string
    // The ID of the rubricQuality that this feedback is related to.
    qualityId *string
}
// NewRubricQualityFeedbackModel instantiates a new rubricQualityFeedbackModel and sets the default values.
func NewRubricQualityFeedbackModel()(*RubricQualityFeedbackModel) {
    m := &RubricQualityFeedbackModel{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateRubricQualityFeedbackModelFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRubricQualityFeedbackModelFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRubricQualityFeedbackModel(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RubricQualityFeedbackModel) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFeedback gets the feedback property value. Specific feedback for one quality of this rubric.
func (m *RubricQualityFeedbackModel) GetFeedback()(EducationItemBodyable) {
    return m.feedback
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RubricQualityFeedbackModel) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["feedback"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEducationItemBodyFromDiscriminatorValue , m.SetFeedback)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["qualityId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQualityId)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *RubricQualityFeedbackModel) GetOdataType()(*string) {
    return m.odataType
}
// GetQualityId gets the qualityId property value. The ID of the rubricQuality that this feedback is related to.
func (m *RubricQualityFeedbackModel) GetQualityId()(*string) {
    return m.qualityId
}
// Serialize serializes information the current object
func (m *RubricQualityFeedbackModel) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("feedback", m.GetFeedback())
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
        err := writer.WriteStringValue("qualityId", m.GetQualityId())
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
func (m *RubricQualityFeedbackModel) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetFeedback sets the feedback property value. Specific feedback for one quality of this rubric.
func (m *RubricQualityFeedbackModel) SetFeedback(value EducationItemBodyable)() {
    m.feedback = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RubricQualityFeedbackModel) SetOdataType(value *string)() {
    m.odataType = value
}
// SetQualityId sets the qualityId property value. The ID of the rubricQuality that this feedback is related to.
func (m *RubricQualityFeedbackModel) SetQualityId(value *string)() {
    m.qualityId = value
}
