package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DataPolicyOperation provides operations to manage the collection of dataPolicyOperation entities.
type DataPolicyOperation struct {
    Entity
    // Represents when the request for this data policy operation was completed, in UTC time, using the ISO 8601 format. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Null until the operation completes.
    completedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Specifies the progress of an operation.
    progress *float64
    // Possible values are: notStarted, running, complete, failed, unknownFutureValue.
    status *DataPolicyOperationStatus
    // The URL location to where data is being exported for export requests.
    storageLocation *string
    // Represents when the request for this data operation was submitted, in UTC time, using the ISO 8601 format. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    submittedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The id for the user on whom the operation is performed.
    userId *string
}
// NewDataPolicyOperation instantiates a new dataPolicyOperation and sets the default values.
func NewDataPolicyOperation()(*DataPolicyOperation) {
    m := &DataPolicyOperation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDataPolicyOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDataPolicyOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDataPolicyOperation(), nil
}
// GetCompletedDateTime gets the completedDateTime property value. Represents when the request for this data policy operation was completed, in UTC time, using the ISO 8601 format. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Null until the operation completes.
func (m *DataPolicyOperation) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.completedDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DataPolicyOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["completedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCompletedDateTime)
    res["progress"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetProgress)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDataPolicyOperationStatus , m.SetStatus)
    res["storageLocation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetStorageLocation)
    res["submittedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetSubmittedDateTime)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    return res
}
// GetProgress gets the progress property value. Specifies the progress of an operation.
func (m *DataPolicyOperation) GetProgress()(*float64) {
    return m.progress
}
// GetStatus gets the status property value. Possible values are: notStarted, running, complete, failed, unknownFutureValue.
func (m *DataPolicyOperation) GetStatus()(*DataPolicyOperationStatus) {
    return m.status
}
// GetStorageLocation gets the storageLocation property value. The URL location to where data is being exported for export requests.
func (m *DataPolicyOperation) GetStorageLocation()(*string) {
    return m.storageLocation
}
// GetSubmittedDateTime gets the submittedDateTime property value. Represents when the request for this data operation was submitted, in UTC time, using the ISO 8601 format. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *DataPolicyOperation) GetSubmittedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.submittedDateTime
}
// GetUserId gets the userId property value. The id for the user on whom the operation is performed.
func (m *DataPolicyOperation) GetUserId()(*string) {
    return m.userId
}
// Serialize serializes information the current object
func (m *DataPolicyOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("completedDateTime", m.GetCompletedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("progress", m.GetProgress())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("storageLocation", m.GetStorageLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("submittedDateTime", m.GetSubmittedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompletedDateTime sets the completedDateTime property value. Represents when the request for this data policy operation was completed, in UTC time, using the ISO 8601 format. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Null until the operation completes.
func (m *DataPolicyOperation) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.completedDateTime = value
}
// SetProgress sets the progress property value. Specifies the progress of an operation.
func (m *DataPolicyOperation) SetProgress(value *float64)() {
    m.progress = value
}
// SetStatus sets the status property value. Possible values are: notStarted, running, complete, failed, unknownFutureValue.
func (m *DataPolicyOperation) SetStatus(value *DataPolicyOperationStatus)() {
    m.status = value
}
// SetStorageLocation sets the storageLocation property value. The URL location to where data is being exported for export requests.
func (m *DataPolicyOperation) SetStorageLocation(value *string)() {
    m.storageLocation = value
}
// SetSubmittedDateTime sets the submittedDateTime property value. Represents when the request for this data operation was submitted, in UTC time, using the ISO 8601 format. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *DataPolicyOperation) SetSubmittedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.submittedDateTime = value
}
// SetUserId sets the userId property value. The id for the user on whom the operation is performed.
func (m *DataPolicyOperation) SetUserId(value *string)() {
    m.userId = value
}
