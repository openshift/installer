package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CloudCommunications provides operations to manage the cloudCommunications singleton.
type CloudCommunications struct {
    Entity
    // The calls property
    calls []Callable
    // The onlineMeetings property
    onlineMeetings []OnlineMeetingable
    // The presences property
    presences []Presenceable
}
// NewCloudCommunications instantiates a new cloudCommunications and sets the default values.
func NewCloudCommunications()(*CloudCommunications) {
    m := &CloudCommunications{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCloudCommunicationsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCloudCommunicationsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudCommunications(), nil
}
// GetCalls gets the calls property value. The calls property
func (m *CloudCommunications) GetCalls()([]Callable) {
    return m.calls
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CloudCommunications) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["calls"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCallFromDiscriminatorValue , m.SetCalls)
    res["onlineMeetings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateOnlineMeetingFromDiscriminatorValue , m.SetOnlineMeetings)
    res["presences"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePresenceFromDiscriminatorValue , m.SetPresences)
    return res
}
// GetOnlineMeetings gets the onlineMeetings property value. The onlineMeetings property
func (m *CloudCommunications) GetOnlineMeetings()([]OnlineMeetingable) {
    return m.onlineMeetings
}
// GetPresences gets the presences property value. The presences property
func (m *CloudCommunications) GetPresences()([]Presenceable) {
    return m.presences
}
// Serialize serializes information the current object
func (m *CloudCommunications) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCalls() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCalls())
        err = writer.WriteCollectionOfObjectValues("calls", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOnlineMeetings() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOnlineMeetings())
        err = writer.WriteCollectionOfObjectValues("onlineMeetings", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPresences() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPresences())
        err = writer.WriteCollectionOfObjectValues("presences", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCalls sets the calls property value. The calls property
func (m *CloudCommunications) SetCalls(value []Callable)() {
    m.calls = value
}
// SetOnlineMeetings sets the onlineMeetings property value. The onlineMeetings property
func (m *CloudCommunications) SetOnlineMeetings(value []OnlineMeetingable)() {
    m.onlineMeetings = value
}
// SetPresences sets the presences property value. The presences property
func (m *CloudCommunications) SetPresences(value []Presenceable)() {
    m.presences = value
}
