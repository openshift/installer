package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SimulationReport 
type SimulationReport struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // Overview of an attack simulation and training campaign.
    overview SimulationReportOverviewable
    // The tenant users and their online actions in an attack simulation and training campaign.
    simulationUsers []UserSimulationDetailsable
}
// NewSimulationReport instantiates a new simulationReport and sets the default values.
func NewSimulationReport()(*SimulationReport) {
    m := &SimulationReport{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSimulationReportFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSimulationReportFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSimulationReport(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SimulationReport) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SimulationReport) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["overview"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSimulationReportOverviewFromDiscriminatorValue , m.SetOverview)
    res["simulationUsers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateUserSimulationDetailsFromDiscriminatorValue , m.SetSimulationUsers)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SimulationReport) GetOdataType()(*string) {
    return m.odataType
}
// GetOverview gets the overview property value. Overview of an attack simulation and training campaign.
func (m *SimulationReport) GetOverview()(SimulationReportOverviewable) {
    return m.overview
}
// GetSimulationUsers gets the simulationUsers property value. The tenant users and their online actions in an attack simulation and training campaign.
func (m *SimulationReport) GetSimulationUsers()([]UserSimulationDetailsable) {
    return m.simulationUsers
}
// Serialize serializes information the current object
func (m *SimulationReport) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("overview", m.GetOverview())
        if err != nil {
            return err
        }
    }
    if m.GetSimulationUsers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSimulationUsers())
        err := writer.WriteCollectionOfObjectValues("simulationUsers", cast)
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
func (m *SimulationReport) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SimulationReport) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOverview sets the overview property value. Overview of an attack simulation and training campaign.
func (m *SimulationReport) SetOverview(value SimulationReportOverviewable)() {
    m.overview = value
}
// SetSimulationUsers sets the simulationUsers property value. The tenant users and their online actions in an attack simulation and training campaign.
func (m *SimulationReport) SetSimulationUsers(value []UserSimulationDetailsable)() {
    m.simulationUsers = value
}
