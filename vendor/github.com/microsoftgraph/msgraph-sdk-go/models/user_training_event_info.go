package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserTrainingEventInfo 
type UserTrainingEventInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Display name of the training.
    displayName *string
    // Latest status of the training assigned to the user. Possible values are: unknown, assigned, inProgress, completed, overdue, unknownFutureValue.
    latestTrainingStatus *TrainingStatus
    // The OdataType property
    odataType *string
    // Event details of the training when it was assigned to the user.
    trainingAssignedProperties UserTrainingContentEventInfoable
    // Event details of the training when it was completed by the user.
    trainingCompletedProperties UserTrainingContentEventInfoable
    // Event details of the training when it was updated/in-progress by the user.
    trainingUpdatedProperties UserTrainingContentEventInfoable
}
// NewUserTrainingEventInfo instantiates a new userTrainingEventInfo and sets the default values.
func NewUserTrainingEventInfo()(*UserTrainingEventInfo) {
    m := &UserTrainingEventInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateUserTrainingEventInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserTrainingEventInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserTrainingEventInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UserTrainingEventInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDisplayName gets the displayName property value. Display name of the training.
func (m *UserTrainingEventInfo) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserTrainingEventInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["latestTrainingStatus"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseTrainingStatus , m.SetLatestTrainingStatus)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["trainingAssignedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUserTrainingContentEventInfoFromDiscriminatorValue , m.SetTrainingAssignedProperties)
    res["trainingCompletedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUserTrainingContentEventInfoFromDiscriminatorValue , m.SetTrainingCompletedProperties)
    res["trainingUpdatedProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUserTrainingContentEventInfoFromDiscriminatorValue , m.SetTrainingUpdatedProperties)
    return res
}
// GetLatestTrainingStatus gets the latestTrainingStatus property value. Latest status of the training assigned to the user. Possible values are: unknown, assigned, inProgress, completed, overdue, unknownFutureValue.
func (m *UserTrainingEventInfo) GetLatestTrainingStatus()(*TrainingStatus) {
    return m.latestTrainingStatus
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *UserTrainingEventInfo) GetOdataType()(*string) {
    return m.odataType
}
// GetTrainingAssignedProperties gets the trainingAssignedProperties property value. Event details of the training when it was assigned to the user.
func (m *UserTrainingEventInfo) GetTrainingAssignedProperties()(UserTrainingContentEventInfoable) {
    return m.trainingAssignedProperties
}
// GetTrainingCompletedProperties gets the trainingCompletedProperties property value. Event details of the training when it was completed by the user.
func (m *UserTrainingEventInfo) GetTrainingCompletedProperties()(UserTrainingContentEventInfoable) {
    return m.trainingCompletedProperties
}
// GetTrainingUpdatedProperties gets the trainingUpdatedProperties property value. Event details of the training when it was updated/in-progress by the user.
func (m *UserTrainingEventInfo) GetTrainingUpdatedProperties()(UserTrainingContentEventInfoable) {
    return m.trainingUpdatedProperties
}
// Serialize serializes information the current object
func (m *UserTrainingEventInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetLatestTrainingStatus() != nil {
        cast := (*m.GetLatestTrainingStatus()).String()
        err := writer.WriteStringValue("latestTrainingStatus", &cast)
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
        err := writer.WriteObjectValue("trainingAssignedProperties", m.GetTrainingAssignedProperties())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("trainingCompletedProperties", m.GetTrainingCompletedProperties())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("trainingUpdatedProperties", m.GetTrainingUpdatedProperties())
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
func (m *UserTrainingEventInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDisplayName sets the displayName property value. Display name of the training.
func (m *UserTrainingEventInfo) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetLatestTrainingStatus sets the latestTrainingStatus property value. Latest status of the training assigned to the user. Possible values are: unknown, assigned, inProgress, completed, overdue, unknownFutureValue.
func (m *UserTrainingEventInfo) SetLatestTrainingStatus(value *TrainingStatus)() {
    m.latestTrainingStatus = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserTrainingEventInfo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTrainingAssignedProperties sets the trainingAssignedProperties property value. Event details of the training when it was assigned to the user.
func (m *UserTrainingEventInfo) SetTrainingAssignedProperties(value UserTrainingContentEventInfoable)() {
    m.trainingAssignedProperties = value
}
// SetTrainingCompletedProperties sets the trainingCompletedProperties property value. Event details of the training when it was completed by the user.
func (m *UserTrainingEventInfo) SetTrainingCompletedProperties(value UserTrainingContentEventInfoable)() {
    m.trainingCompletedProperties = value
}
// SetTrainingUpdatedProperties sets the trainingUpdatedProperties property value. Event details of the training when it was updated/in-progress by the user.
func (m *UserTrainingEventInfo) SetTrainingUpdatedProperties(value UserTrainingContentEventInfoable)() {
    m.trainingUpdatedProperties = value
}
