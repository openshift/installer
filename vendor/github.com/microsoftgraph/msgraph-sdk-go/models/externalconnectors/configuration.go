package externalconnectors

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Configuration 
type Configuration struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // A collection of application IDs for registered Azure Active Directory apps that are allowed to manage the externalConnection and to index content in the externalConnection.
    authorizedAppIds []string
    // The OdataType property
    odataType *string
}
// NewConfiguration instantiates a new configuration and sets the default values.
func NewConfiguration()(*Configuration) {
    m := &Configuration{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConfiguration(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Configuration) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAuthorizedAppIds gets the authorizedAppIds property value. A collection of application IDs for registered Azure Active Directory apps that are allowed to manage the externalConnection and to index content in the externalConnection.
func (m *Configuration) GetAuthorizedAppIds()([]string) {
    return m.authorizedAppIds
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Configuration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["authorizedAppIds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetAuthorizedAppIds)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Configuration) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *Configuration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAuthorizedAppIds() != nil {
        err := writer.WriteCollectionOfStringValues("authorizedAppIds", m.GetAuthorizedAppIds())
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
func (m *Configuration) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAuthorizedAppIds sets the authorizedAppIds property value. A collection of application IDs for registered Azure Active Directory apps that are allowed to manage the externalConnection and to index content in the externalConnection.
func (m *Configuration) SetAuthorizedAppIds(value []string)() {
    m.authorizedAppIds = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Configuration) SetOdataType(value *string)() {
    m.odataType = value
}
