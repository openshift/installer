package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserFlowApiConnectorConfiguration 
type UserFlowApiConnectorConfiguration struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // The postAttributeCollection property
    postAttributeCollection IdentityApiConnectorable
    // The postFederationSignup property
    postFederationSignup IdentityApiConnectorable
}
// NewUserFlowApiConnectorConfiguration instantiates a new userFlowApiConnectorConfiguration and sets the default values.
func NewUserFlowApiConnectorConfiguration()(*UserFlowApiConnectorConfiguration) {
    m := &UserFlowApiConnectorConfiguration{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateUserFlowApiConnectorConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserFlowApiConnectorConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserFlowApiConnectorConfiguration(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UserFlowApiConnectorConfiguration) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserFlowApiConnectorConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["postAttributeCollection"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentityApiConnectorFromDiscriminatorValue , m.SetPostAttributeCollection)
    res["postFederationSignup"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentityApiConnectorFromDiscriminatorValue , m.SetPostFederationSignup)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *UserFlowApiConnectorConfiguration) GetOdataType()(*string) {
    return m.odataType
}
// GetPostAttributeCollection gets the postAttributeCollection property value. The postAttributeCollection property
func (m *UserFlowApiConnectorConfiguration) GetPostAttributeCollection()(IdentityApiConnectorable) {
    return m.postAttributeCollection
}
// GetPostFederationSignup gets the postFederationSignup property value. The postFederationSignup property
func (m *UserFlowApiConnectorConfiguration) GetPostFederationSignup()(IdentityApiConnectorable) {
    return m.postFederationSignup
}
// Serialize serializes information the current object
func (m *UserFlowApiConnectorConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("postAttributeCollection", m.GetPostAttributeCollection())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("postFederationSignup", m.GetPostFederationSignup())
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
func (m *UserFlowApiConnectorConfiguration) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserFlowApiConnectorConfiguration) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPostAttributeCollection sets the postAttributeCollection property value. The postAttributeCollection property
func (m *UserFlowApiConnectorConfiguration) SetPostAttributeCollection(value IdentityApiConnectorable)() {
    m.postAttributeCollection = value
}
// SetPostFederationSignup sets the postFederationSignup property value. The postFederationSignup property
func (m *UserFlowApiConnectorConfiguration) SetPostFederationSignup(value IdentityApiConnectorable)() {
    m.postFederationSignup = value
}
