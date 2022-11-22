package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AttackSimulationRoot 
type AttackSimulationRoot struct {
    Entity
    // Represents simulation automation created to run on a tenant.
    simulationAutomations []SimulationAutomationable
    // Represents an attack simulation training campaign in a tenant.
    simulations []Simulationable
}
// NewAttackSimulationRoot instantiates a new attackSimulationRoot and sets the default values.
func NewAttackSimulationRoot()(*AttackSimulationRoot) {
    m := &AttackSimulationRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAttackSimulationRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAttackSimulationRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttackSimulationRoot(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AttackSimulationRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["simulationAutomations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSimulationAutomationFromDiscriminatorValue , m.SetSimulationAutomations)
    res["simulations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSimulationFromDiscriminatorValue , m.SetSimulations)
    return res
}
// GetSimulationAutomations gets the simulationAutomations property value. Represents simulation automation created to run on a tenant.
func (m *AttackSimulationRoot) GetSimulationAutomations()([]SimulationAutomationable) {
    return m.simulationAutomations
}
// GetSimulations gets the simulations property value. Represents an attack simulation training campaign in a tenant.
func (m *AttackSimulationRoot) GetSimulations()([]Simulationable) {
    return m.simulations
}
// Serialize serializes information the current object
func (m *AttackSimulationRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetSimulationAutomations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSimulationAutomations())
        err = writer.WriteCollectionOfObjectValues("simulationAutomations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSimulations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSimulations())
        err = writer.WriteCollectionOfObjectValues("simulations", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSimulationAutomations sets the simulationAutomations property value. Represents simulation automation created to run on a tenant.
func (m *AttackSimulationRoot) SetSimulationAutomations(value []SimulationAutomationable)() {
    m.simulationAutomations = value
}
// SetSimulations sets the simulations property value. Represents an attack simulation training campaign in a tenant.
func (m *AttackSimulationRoot) SetSimulations(value []Simulationable)() {
    m.simulations = value
}
