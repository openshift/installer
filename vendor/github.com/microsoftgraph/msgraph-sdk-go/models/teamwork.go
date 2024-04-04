package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Teamwork 
type Teamwork struct {
    Entity
}
// NewTeamwork instantiates a new Teamwork and sets the default values.
func NewTeamwork()(*Teamwork) {
    m := &Teamwork{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamworkFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamworkFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamwork(), nil
}
// GetDeletedTeams gets the deletedTeams property value. The deleted team.
func (m *Teamwork) GetDeletedTeams()([]DeletedTeamable) {
    val, err := m.GetBackingStore().Get("deletedTeams")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeletedTeamable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Teamwork) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deletedTeams"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeletedTeamFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeletedTeamable, len(val))
            for i, v := range val {
                res[i] = v.(DeletedTeamable)
            }
            m.SetDeletedTeams(res)
        }
        return nil
    }
    res["workforceIntegrations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkforceIntegrationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkforceIntegrationable, len(val))
            for i, v := range val {
                res[i] = v.(WorkforceIntegrationable)
            }
            m.SetWorkforceIntegrations(res)
        }
        return nil
    }
    return res
}
// GetWorkforceIntegrations gets the workforceIntegrations property value. The workforceIntegrations property
func (m *Teamwork) GetWorkforceIntegrations()([]WorkforceIntegrationable) {
    val, err := m.GetBackingStore().Get("workforceIntegrations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkforceIntegrationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Teamwork) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDeletedTeams() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeletedTeams()))
        for i, v := range m.GetDeletedTeams() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("deletedTeams", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWorkforceIntegrations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWorkforceIntegrations()))
        for i, v := range m.GetWorkforceIntegrations() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("workforceIntegrations", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeletedTeams sets the deletedTeams property value. The deleted team.
func (m *Teamwork) SetDeletedTeams(value []DeletedTeamable)() {
    err := m.GetBackingStore().Set("deletedTeams", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkforceIntegrations sets the workforceIntegrations property value. The workforceIntegrations property
func (m *Teamwork) SetWorkforceIntegrations(value []WorkforceIntegrationable)() {
    err := m.GetBackingStore().Set("workforceIntegrations", value)
    if err != nil {
        panic(err)
    }
}
// Teamworkable 
type Teamworkable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeletedTeams()([]DeletedTeamable)
    GetWorkforceIntegrations()([]WorkforceIntegrationable)
    SetDeletedTeams(value []DeletedTeamable)()
    SetWorkforceIntegrations(value []WorkforceIntegrationable)()
}
