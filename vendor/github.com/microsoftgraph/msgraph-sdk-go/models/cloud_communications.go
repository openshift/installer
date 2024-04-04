package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CloudCommunications 
type CloudCommunications struct {
    Entity
}
// NewCloudCommunications instantiates a new CloudCommunications and sets the default values.
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
    val, err := m.GetBackingStore().Get("calls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Callable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CloudCommunications) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["calls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCallFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Callable, len(val))
            for i, v := range val {
                res[i] = v.(Callable)
            }
            m.SetCalls(res)
        }
        return nil
    }
    res["onlineMeetings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnlineMeetingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnlineMeetingable, len(val))
            for i, v := range val {
                res[i] = v.(OnlineMeetingable)
            }
            m.SetOnlineMeetings(res)
        }
        return nil
    }
    res["presences"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePresenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Presenceable, len(val))
            for i, v := range val {
                res[i] = v.(Presenceable)
            }
            m.SetPresences(res)
        }
        return nil
    }
    return res
}
// GetOnlineMeetings gets the onlineMeetings property value. The onlineMeetings property
func (m *CloudCommunications) GetOnlineMeetings()([]OnlineMeetingable) {
    val, err := m.GetBackingStore().Get("onlineMeetings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnlineMeetingable)
    }
    return nil
}
// GetPresences gets the presences property value. The presences property
func (m *CloudCommunications) GetPresences()([]Presenceable) {
    val, err := m.GetBackingStore().Get("presences")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Presenceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CloudCommunications) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCalls() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCalls()))
        for i, v := range m.GetCalls() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("calls", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOnlineMeetings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOnlineMeetings()))
        for i, v := range m.GetOnlineMeetings() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("onlineMeetings", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPresences() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPresences()))
        for i, v := range m.GetPresences() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("presences", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCalls sets the calls property value. The calls property
func (m *CloudCommunications) SetCalls(value []Callable)() {
    err := m.GetBackingStore().Set("calls", value)
    if err != nil {
        panic(err)
    }
}
// SetOnlineMeetings sets the onlineMeetings property value. The onlineMeetings property
func (m *CloudCommunications) SetOnlineMeetings(value []OnlineMeetingable)() {
    err := m.GetBackingStore().Set("onlineMeetings", value)
    if err != nil {
        panic(err)
    }
}
// SetPresences sets the presences property value. The presences property
func (m *CloudCommunications) SetPresences(value []Presenceable)() {
    err := m.GetBackingStore().Set("presences", value)
    if err != nil {
        panic(err)
    }
}
// CloudCommunicationsable 
type CloudCommunicationsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCalls()([]Callable)
    GetOnlineMeetings()([]OnlineMeetingable)
    GetPresences()([]Presenceable)
    SetCalls(value []Callable)()
    SetOnlineMeetings(value []OnlineMeetingable)()
    SetPresences(value []Presenceable)()
}
