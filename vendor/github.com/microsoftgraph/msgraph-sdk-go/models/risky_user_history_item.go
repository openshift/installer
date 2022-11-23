package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RiskyUserHistoryItem provides operations to manage the collection of agreement entities.
type RiskyUserHistoryItem struct {
    RiskyUser
    // The activity related to user risk level change.
    activity RiskUserActivityable
    // The ID of actor that does the operation.
    initiatedBy *string
    // The ID of the user.
    userId *string
}
// NewRiskyUserHistoryItem instantiates a new riskyUserHistoryItem and sets the default values.
func NewRiskyUserHistoryItem()(*RiskyUserHistoryItem) {
    m := &RiskyUserHistoryItem{
        RiskyUser: *NewRiskyUser(),
    }
    return m
}
// CreateRiskyUserHistoryItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRiskyUserHistoryItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRiskyUserHistoryItem(), nil
}
// GetActivity gets the activity property value. The activity related to user risk level change.
func (m *RiskyUserHistoryItem) GetActivity()(RiskUserActivityable) {
    return m.activity
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RiskyUserHistoryItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RiskyUser.GetFieldDeserializers()
    res["activity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRiskUserActivityFromDiscriminatorValue , m.SetActivity)
    res["initiatedBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInitiatedBy)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    return res
}
// GetInitiatedBy gets the initiatedBy property value. The ID of actor that does the operation.
func (m *RiskyUserHistoryItem) GetInitiatedBy()(*string) {
    return m.initiatedBy
}
// GetUserId gets the userId property value. The ID of the user.
func (m *RiskyUserHistoryItem) GetUserId()(*string) {
    return m.userId
}
// Serialize serializes information the current object
func (m *RiskyUserHistoryItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RiskyUser.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("activity", m.GetActivity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("initiatedBy", m.GetInitiatedBy())
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
// SetActivity sets the activity property value. The activity related to user risk level change.
func (m *RiskyUserHistoryItem) SetActivity(value RiskUserActivityable)() {
    m.activity = value
}
// SetInitiatedBy sets the initiatedBy property value. The ID of actor that does the operation.
func (m *RiskyUserHistoryItem) SetInitiatedBy(value *string)() {
    m.initiatedBy = value
}
// SetUserId sets the userId property value. The ID of the user.
func (m *RiskyUserHistoryItem) SetUserId(value *string)() {
    m.userId = value
}
