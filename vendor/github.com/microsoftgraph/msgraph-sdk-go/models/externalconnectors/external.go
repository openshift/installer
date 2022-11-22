package externalconnectors

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// External 
type External struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The connections property
    connections []ExternalConnectionable
    // The OdataType property
    odataType *string
}
// NewExternal instantiates a new External and sets the default values.
func NewExternal()(*External) {
    m := &External{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateExternalFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateExternalFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewExternal(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *External) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetConnections gets the connections property value. The connections property
func (m *External) GetConnections()([]ExternalConnectionable) {
    return m.connections
}
// GetFieldDeserializers the deserialization information for the current model
func (m *External) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["connections"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateExternalConnectionFromDiscriminatorValue , m.SetConnections)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *External) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *External) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetConnections() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetConnections())
        err := writer.WriteCollectionOfObjectValues("connections", cast)
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
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *External) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetConnections sets the connections property value. The connections property
func (m *External) SetConnections(value []ExternalConnectionable)() {
    m.connections = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *External) SetOdataType(value *string)() {
    m.odataType = value
}
