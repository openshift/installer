package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SimulationReportOverview 
type SimulationReportOverview struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // List of recommended actions for a tenant to improve its security posture based on the attack simulation and training campaign attack type.
    recommendedActions []RecommendedActionable
    // Number of valid users in the attack simulation and training campaign.
    resolvedTargetsCount *int32
    // Summary of simulation events in the attack simulation and training campaign.
    simulationEventsContent SimulationEventsContentable
    // Summary of assigned trainings in the attack simulation and training campaign.
    trainingEventsContent TrainingEventsContentable
}
// NewSimulationReportOverview instantiates a new simulationReportOverview and sets the default values.
func NewSimulationReportOverview()(*SimulationReportOverview) {
    m := &SimulationReportOverview{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSimulationReportOverviewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSimulationReportOverviewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSimulationReportOverview(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SimulationReportOverview) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SimulationReportOverview) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["recommendedActions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRecommendedActionFromDiscriminatorValue , m.SetRecommendedActions)
    res["resolvedTargetsCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetResolvedTargetsCount)
    res["simulationEventsContent"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSimulationEventsContentFromDiscriminatorValue , m.SetSimulationEventsContent)
    res["trainingEventsContent"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTrainingEventsContentFromDiscriminatorValue , m.SetTrainingEventsContent)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SimulationReportOverview) GetOdataType()(*string) {
    return m.odataType
}
// GetRecommendedActions gets the recommendedActions property value. List of recommended actions for a tenant to improve its security posture based on the attack simulation and training campaign attack type.
func (m *SimulationReportOverview) GetRecommendedActions()([]RecommendedActionable) {
    return m.recommendedActions
}
// GetResolvedTargetsCount gets the resolvedTargetsCount property value. Number of valid users in the attack simulation and training campaign.
func (m *SimulationReportOverview) GetResolvedTargetsCount()(*int32) {
    return m.resolvedTargetsCount
}
// GetSimulationEventsContent gets the simulationEventsContent property value. Summary of simulation events in the attack simulation and training campaign.
func (m *SimulationReportOverview) GetSimulationEventsContent()(SimulationEventsContentable) {
    return m.simulationEventsContent
}
// GetTrainingEventsContent gets the trainingEventsContent property value. Summary of assigned trainings in the attack simulation and training campaign.
func (m *SimulationReportOverview) GetTrainingEventsContent()(TrainingEventsContentable) {
    return m.trainingEventsContent
}
// Serialize serializes information the current object
func (m *SimulationReportOverview) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetRecommendedActions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetRecommendedActions())
        err := writer.WriteCollectionOfObjectValues("recommendedActions", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("resolvedTargetsCount", m.GetResolvedTargetsCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("simulationEventsContent", m.GetSimulationEventsContent())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("trainingEventsContent", m.GetTrainingEventsContent())
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
func (m *SimulationReportOverview) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SimulationReportOverview) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRecommendedActions sets the recommendedActions property value. List of recommended actions for a tenant to improve its security posture based on the attack simulation and training campaign attack type.
func (m *SimulationReportOverview) SetRecommendedActions(value []RecommendedActionable)() {
    m.recommendedActions = value
}
// SetResolvedTargetsCount sets the resolvedTargetsCount property value. Number of valid users in the attack simulation and training campaign.
func (m *SimulationReportOverview) SetResolvedTargetsCount(value *int32)() {
    m.resolvedTargetsCount = value
}
// SetSimulationEventsContent sets the simulationEventsContent property value. Summary of simulation events in the attack simulation and training campaign.
func (m *SimulationReportOverview) SetSimulationEventsContent(value SimulationEventsContentable)() {
    m.simulationEventsContent = value
}
// SetTrainingEventsContent sets the trainingEventsContent property value. Summary of assigned trainings in the attack simulation and training campaign.
func (m *SimulationReportOverview) SetTrainingEventsContent(value TrainingEventsContentable)() {
    m.trainingEventsContent = value
}
