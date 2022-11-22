package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// LocationConstraintItem 
type LocationConstraintItem struct {
    Location
    // If set to true and the specified resource is busy, findMeetingTimes looks for another resource that is free. If set to false and the specified resource is busy, findMeetingTimes returns the resource best ranked in the user's cache without checking if it's free. Default is true.
    resolveAvailability *bool
}
// NewLocationConstraintItem instantiates a new LocationConstraintItem and sets the default values.
func NewLocationConstraintItem()(*LocationConstraintItem) {
    m := &LocationConstraintItem{
        Location: *NewLocation(),
    }
    odataTypeValue := "#microsoft.graph.locationConstraintItem";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateLocationConstraintItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateLocationConstraintItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLocationConstraintItem(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *LocationConstraintItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Location.GetFieldDeserializers()
    res["resolveAvailability"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetResolveAvailability)
    return res
}
// GetResolveAvailability gets the resolveAvailability property value. If set to true and the specified resource is busy, findMeetingTimes looks for another resource that is free. If set to false and the specified resource is busy, findMeetingTimes returns the resource best ranked in the user's cache without checking if it's free. Default is true.
func (m *LocationConstraintItem) GetResolveAvailability()(*bool) {
    return m.resolveAvailability
}
// Serialize serializes information the current object
func (m *LocationConstraintItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Location.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("resolveAvailability", m.GetResolveAvailability())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetResolveAvailability sets the resolveAvailability property value. If set to true and the specified resource is busy, findMeetingTimes looks for another resource that is free. If set to false and the specified resource is busy, findMeetingTimes returns the resource best ranked in the user's cache without checking if it's free. Default is true.
func (m *LocationConstraintItem) SetResolveAvailability(value *bool)() {
    m.resolveAvailability = value
}
