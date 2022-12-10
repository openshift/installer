package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PlannerPlanContainer 
type PlannerPlanContainer struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The identifier of the resource that contains the plan.
    containerId *string
    // The OdataType property
    odataType *string
    // The type of the resource that contains the plan. For supported types, see the previous table. Possible values are: group, unknownFutureValue, roster. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value in this evolvable enum: roster.
    type_escaped *PlannerContainerType
    // The full canonical URL of the container.
    url *string
}
// NewPlannerPlanContainer instantiates a new plannerPlanContainer and sets the default values.
func NewPlannerPlanContainer()(*PlannerPlanContainer) {
    m := &PlannerPlanContainer{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePlannerPlanContainerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePlannerPlanContainerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerPlanContainer(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PlannerPlanContainer) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetContainerId gets the containerId property value. The identifier of the resource that contains the plan.
func (m *PlannerPlanContainer) GetContainerId()(*string) {
    return m.containerId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PlannerPlanContainer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["containerId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetContainerId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParsePlannerContainerType , m.SetType)
    res["url"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUrl)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *PlannerPlanContainer) GetOdataType()(*string) {
    return m.odataType
}
// GetType gets the type property value. The type of the resource that contains the plan. For supported types, see the previous table. Possible values are: group, unknownFutureValue, roster. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value in this evolvable enum: roster.
func (m *PlannerPlanContainer) GetType()(*PlannerContainerType) {
    return m.type_escaped
}
// GetUrl gets the url property value. The full canonical URL of the container.
func (m *PlannerPlanContainer) GetUrl()(*string) {
    return m.url
}
// Serialize serializes information the current object
func (m *PlannerPlanContainer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("containerId", m.GetContainerId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetType() != nil {
        cast := (*m.GetType()).String()
        err := writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("url", m.GetUrl())
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
func (m *PlannerPlanContainer) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetContainerId sets the containerId property value. The identifier of the resource that contains the plan.
func (m *PlannerPlanContainer) SetContainerId(value *string)() {
    m.containerId = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PlannerPlanContainer) SetOdataType(value *string)() {
    m.odataType = value
}
// SetType sets the type property value. The type of the resource that contains the plan. For supported types, see the previous table. Possible values are: group, unknownFutureValue, roster. Note that you must use the Prefer: include-unknown-enum-members request header to get the following value in this evolvable enum: roster.
func (m *PlannerPlanContainer) SetType(value *PlannerContainerType)() {
    m.type_escaped = value
}
// SetUrl sets the url property value. The full canonical URL of the container.
func (m *PlannerPlanContainer) SetUrl(value *string)() {
    m.url = value
}
