package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Security 
type Security struct {
    Entity
    // The alerts property
    alerts []Alertable
    // The attackSimulation property
    attackSimulation AttackSimulationRootable
    // The secureScoreControlProfiles property
    secureScoreControlProfiles []SecureScoreControlProfileable
    // The secureScores property
    secureScores []SecureScoreable
}
// NewSecurity instantiates a new Security and sets the default values.
func NewSecurity()(*Security) {
    m := &Security{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSecurityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSecurityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSecurity(), nil
}
// GetAlerts gets the alerts property value. The alerts property
func (m *Security) GetAlerts()([]Alertable) {
    return m.alerts
}
// GetAttackSimulation gets the attackSimulation property value. The attackSimulation property
func (m *Security) GetAttackSimulation()(AttackSimulationRootable) {
    return m.attackSimulation
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Security) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["alerts"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAlertFromDiscriminatorValue , m.SetAlerts)
    res["attackSimulation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateAttackSimulationRootFromDiscriminatorValue , m.SetAttackSimulation)
    res["secureScoreControlProfiles"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSecureScoreControlProfileFromDiscriminatorValue , m.SetSecureScoreControlProfiles)
    res["secureScores"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSecureScoreFromDiscriminatorValue , m.SetSecureScores)
    return res
}
// GetSecureScoreControlProfiles gets the secureScoreControlProfiles property value. The secureScoreControlProfiles property
func (m *Security) GetSecureScoreControlProfiles()([]SecureScoreControlProfileable) {
    return m.secureScoreControlProfiles
}
// GetSecureScores gets the secureScores property value. The secureScores property
func (m *Security) GetSecureScores()([]SecureScoreable) {
    return m.secureScores
}
// Serialize serializes information the current object
func (m *Security) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAlerts() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAlerts())
        err = writer.WriteCollectionOfObjectValues("alerts", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("attackSimulation", m.GetAttackSimulation())
        if err != nil {
            return err
        }
    }
    if m.GetSecureScoreControlProfiles() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSecureScoreControlProfiles())
        err = writer.WriteCollectionOfObjectValues("secureScoreControlProfiles", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSecureScores() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSecureScores())
        err = writer.WriteCollectionOfObjectValues("secureScores", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAlerts sets the alerts property value. The alerts property
func (m *Security) SetAlerts(value []Alertable)() {
    m.alerts = value
}
// SetAttackSimulation sets the attackSimulation property value. The attackSimulation property
func (m *Security) SetAttackSimulation(value AttackSimulationRootable)() {
    m.attackSimulation = value
}
// SetSecureScoreControlProfiles sets the secureScoreControlProfiles property value. The secureScoreControlProfiles property
func (m *Security) SetSecureScoreControlProfiles(value []SecureScoreControlProfileable)() {
    m.secureScoreControlProfiles = value
}
// SetSecureScores sets the secureScores property value. The secureScores property
func (m *Security) SetSecureScores(value []SecureScoreable)() {
    m.secureScores = value
}
