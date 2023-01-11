package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AssignedTrainingInfo 
type AssignedTrainingInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Number of users who were assigned the training in an attack simulation and training campaign.
    assignedUserCount *int32
    // Number of users who completed the training in an attack simulation and training campaign.
    completedUserCount *int32
    // Display name of the training in an attack simulation and training campaign.
    displayName *string
    // The OdataType property
    odataType *string
}
// NewAssignedTrainingInfo instantiates a new assignedTrainingInfo and sets the default values.
func NewAssignedTrainingInfo()(*AssignedTrainingInfo) {
    m := &AssignedTrainingInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAssignedTrainingInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAssignedTrainingInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAssignedTrainingInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AssignedTrainingInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAssignedUserCount gets the assignedUserCount property value. Number of users who were assigned the training in an attack simulation and training campaign.
func (m *AssignedTrainingInfo) GetAssignedUserCount()(*int32) {
    return m.assignedUserCount
}
// GetCompletedUserCount gets the completedUserCount property value. Number of users who completed the training in an attack simulation and training campaign.
func (m *AssignedTrainingInfo) GetCompletedUserCount()(*int32) {
    return m.completedUserCount
}
// GetDisplayName gets the displayName property value. Display name of the training in an attack simulation and training campaign.
func (m *AssignedTrainingInfo) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AssignedTrainingInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["assignedUserCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetAssignedUserCount)
    res["completedUserCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetCompletedUserCount)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AssignedTrainingInfo) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *AssignedTrainingInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("assignedUserCount", m.GetAssignedUserCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("completedUserCount", m.GetCompletedUserCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
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
func (m *AssignedTrainingInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAssignedUserCount sets the assignedUserCount property value. Number of users who were assigned the training in an attack simulation and training campaign.
func (m *AssignedTrainingInfo) SetAssignedUserCount(value *int32)() {
    m.assignedUserCount = value
}
// SetCompletedUserCount sets the completedUserCount property value. Number of users who completed the training in an attack simulation and training campaign.
func (m *AssignedTrainingInfo) SetCompletedUserCount(value *int32)() {
    m.completedUserCount = value
}
// SetDisplayName sets the displayName property value. Display name of the training in an attack simulation and training campaign.
func (m *AssignedTrainingInfo) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AssignedTrainingInfo) SetOdataType(value *string)() {
    m.odataType = value
}
