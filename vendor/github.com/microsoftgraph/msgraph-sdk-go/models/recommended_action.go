package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RecommendedAction 
type RecommendedAction struct {
    // Web URL to the recommended action.
    actionWebUrl *string
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // Potential improvement in the tenant security score from the recommended action.
    potentialScoreImpact *float64
    // Title of the recommended action.
    title *string
}
// NewRecommendedAction instantiates a new recommendedAction and sets the default values.
func NewRecommendedAction()(*RecommendedAction) {
    m := &RecommendedAction{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateRecommendedActionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRecommendedActionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRecommendedAction(), nil
}
// GetActionWebUrl gets the actionWebUrl property value. Web URL to the recommended action.
func (m *RecommendedAction) GetActionWebUrl()(*string) {
    return m.actionWebUrl
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RecommendedAction) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RecommendedAction) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["actionWebUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetActionWebUrl)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["potentialScoreImpact"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetPotentialScoreImpact)
    res["title"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTitle)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *RecommendedAction) GetOdataType()(*string) {
    return m.odataType
}
// GetPotentialScoreImpact gets the potentialScoreImpact property value. Potential improvement in the tenant security score from the recommended action.
func (m *RecommendedAction) GetPotentialScoreImpact()(*float64) {
    return m.potentialScoreImpact
}
// GetTitle gets the title property value. Title of the recommended action.
func (m *RecommendedAction) GetTitle()(*string) {
    return m.title
}
// Serialize serializes information the current object
func (m *RecommendedAction) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("actionWebUrl", m.GetActionWebUrl())
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
        err := writer.WriteFloat64Value("potentialScoreImpact", m.GetPotentialScoreImpact())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("title", m.GetTitle())
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
// SetActionWebUrl sets the actionWebUrl property value. Web URL to the recommended action.
func (m *RecommendedAction) SetActionWebUrl(value *string)() {
    m.actionWebUrl = value
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RecommendedAction) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RecommendedAction) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPotentialScoreImpact sets the potentialScoreImpact property value. Potential improvement in the tenant security score from the recommended action.
func (m *RecommendedAction) SetPotentialScoreImpact(value *float64)() {
    m.potentialScoreImpact = value
}
// SetTitle sets the title property value. Title of the recommended action.
func (m *RecommendedAction) SetTitle(value *string)() {
    m.title = value
}
