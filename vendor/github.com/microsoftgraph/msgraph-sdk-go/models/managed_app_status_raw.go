package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAppStatusRaw 
type ManagedAppStatusRaw struct {
    ManagedAppStatus
}
// NewManagedAppStatusRaw instantiates a new ManagedAppStatusRaw and sets the default values.
func NewManagedAppStatusRaw()(*ManagedAppStatusRaw) {
    m := &ManagedAppStatusRaw{
        ManagedAppStatus: *NewManagedAppStatus(),
    }
    odataTypeValue := "#microsoft.graph.managedAppStatusRaw"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateManagedAppStatusRawFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateManagedAppStatusRawFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedAppStatusRaw(), nil
}
// GetContent gets the content property value. Status report content.
func (m *ManagedAppStatusRaw) GetContent()(Jsonable) {
    val, err := m.GetBackingStore().Get("content")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Jsonable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ManagedAppStatusRaw) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedAppStatus.GetFieldDeserializers()
    res["content"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateJsonFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContent(val.(Jsonable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *ManagedAppStatusRaw) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedAppStatus.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("content", m.GetContent())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContent sets the content property value. Status report content.
func (m *ManagedAppStatusRaw) SetContent(value Jsonable)() {
    err := m.GetBackingStore().Set("content", value)
    if err != nil {
        panic(err)
    }
}
// ManagedAppStatusRawable 
type ManagedAppStatusRawable interface {
    ManagedAppStatusable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContent()(Jsonable)
    SetContent(value Jsonable)()
}
