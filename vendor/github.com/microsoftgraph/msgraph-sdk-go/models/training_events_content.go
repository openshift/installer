package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TrainingEventsContent 
type TrainingEventsContent struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // List of assigned trainings and their information in an attack simulation and training campaign.
    assignedTrainingsInfos []AssignedTrainingInfoable
    // The OdataType property
    odataType *string
    // Number of users who were assigned trainings in an attack simulation and training campaign.
    trainingsAssignedUserCount *int32
}
// NewTrainingEventsContent instantiates a new trainingEventsContent and sets the default values.
func NewTrainingEventsContent()(*TrainingEventsContent) {
    m := &TrainingEventsContent{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateTrainingEventsContentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTrainingEventsContentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTrainingEventsContent(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TrainingEventsContent) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAssignedTrainingsInfos gets the assignedTrainingsInfos property value. List of assigned trainings and their information in an attack simulation and training campaign.
func (m *TrainingEventsContent) GetAssignedTrainingsInfos()([]AssignedTrainingInfoable) {
    return m.assignedTrainingsInfos
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TrainingEventsContent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["assignedTrainingsInfos"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAssignedTrainingInfoFromDiscriminatorValue , m.SetAssignedTrainingsInfos)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["trainingsAssignedUserCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetTrainingsAssignedUserCount)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *TrainingEventsContent) GetOdataType()(*string) {
    return m.odataType
}
// GetTrainingsAssignedUserCount gets the trainingsAssignedUserCount property value. Number of users who were assigned trainings in an attack simulation and training campaign.
func (m *TrainingEventsContent) GetTrainingsAssignedUserCount()(*int32) {
    return m.trainingsAssignedUserCount
}
// Serialize serializes information the current object
func (m *TrainingEventsContent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAssignedTrainingsInfos() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAssignedTrainingsInfos())
        err := writer.WriteCollectionOfObjectValues("assignedTrainingsInfos", cast)
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
        err := writer.WriteInt32Value("trainingsAssignedUserCount", m.GetTrainingsAssignedUserCount())
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
func (m *TrainingEventsContent) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAssignedTrainingsInfos sets the assignedTrainingsInfos property value. List of assigned trainings and their information in an attack simulation and training campaign.
func (m *TrainingEventsContent) SetAssignedTrainingsInfos(value []AssignedTrainingInfoable)() {
    m.assignedTrainingsInfos = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TrainingEventsContent) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTrainingsAssignedUserCount sets the trainingsAssignedUserCount property value. Number of users who were assigned trainings in an attack simulation and training campaign.
func (m *TrainingEventsContent) SetTrainingsAssignedUserCount(value *int32)() {
    m.trainingsAssignedUserCount = value
}
