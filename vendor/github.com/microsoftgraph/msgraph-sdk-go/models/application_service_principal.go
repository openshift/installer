package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ApplicationServicePrincipal 
type ApplicationServicePrincipal struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The application property
    application Applicationable
    // The OdataType property
    odataType *string
    // The servicePrincipal property
    servicePrincipal ServicePrincipalable
}
// NewApplicationServicePrincipal instantiates a new applicationServicePrincipal and sets the default values.
func NewApplicationServicePrincipal()(*ApplicationServicePrincipal) {
    m := &ApplicationServicePrincipal{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateApplicationServicePrincipalFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateApplicationServicePrincipalFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewApplicationServicePrincipal(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ApplicationServicePrincipal) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetApplication gets the application property value. The application property
func (m *ApplicationServicePrincipal) GetApplication()(Applicationable) {
    return m.application
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ApplicationServicePrincipal) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["application"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateApplicationFromDiscriminatorValue , m.SetApplication)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["servicePrincipal"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateServicePrincipalFromDiscriminatorValue , m.SetServicePrincipal)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ApplicationServicePrincipal) GetOdataType()(*string) {
    return m.odataType
}
// GetServicePrincipal gets the servicePrincipal property value. The servicePrincipal property
func (m *ApplicationServicePrincipal) GetServicePrincipal()(ServicePrincipalable) {
    return m.servicePrincipal
}
// Serialize serializes information the current object
func (m *ApplicationServicePrincipal) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("application", m.GetApplication())
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
        err := writer.WriteObjectValue("servicePrincipal", m.GetServicePrincipal())
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
func (m *ApplicationServicePrincipal) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetApplication sets the application property value. The application property
func (m *ApplicationServicePrincipal) SetApplication(value Applicationable)() {
    m.application = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ApplicationServicePrincipal) SetOdataType(value *string)() {
    m.odataType = value
}
// SetServicePrincipal sets the servicePrincipal property value. The servicePrincipal property
func (m *ApplicationServicePrincipal) SetServicePrincipal(value ServicePrincipalable)() {
    m.servicePrincipal = value
}
