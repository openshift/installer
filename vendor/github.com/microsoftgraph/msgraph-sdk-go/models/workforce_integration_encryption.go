package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WorkforceIntegrationEncryption 
type WorkforceIntegrationEncryption struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // Possible values are: sharedSecret, unknownFutureValue.
    protocol *WorkforceIntegrationEncryptionProtocol
    // Encryption shared secret.
    secret *string
}
// NewWorkforceIntegrationEncryption instantiates a new workforceIntegrationEncryption and sets the default values.
func NewWorkforceIntegrationEncryption()(*WorkforceIntegrationEncryption) {
    m := &WorkforceIntegrationEncryption{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWorkforceIntegrationEncryptionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWorkforceIntegrationEncryptionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkforceIntegrationEncryption(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *WorkforceIntegrationEncryption) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WorkforceIntegrationEncryption) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["protocol"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWorkforceIntegrationEncryptionProtocol , m.SetProtocol)
    res["secret"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSecret)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *WorkforceIntegrationEncryption) GetOdataType()(*string) {
    return m.odataType
}
// GetProtocol gets the protocol property value. Possible values are: sharedSecret, unknownFutureValue.
func (m *WorkforceIntegrationEncryption) GetProtocol()(*WorkforceIntegrationEncryptionProtocol) {
    return m.protocol
}
// GetSecret gets the secret property value. Encryption shared secret.
func (m *WorkforceIntegrationEncryption) GetSecret()(*string) {
    return m.secret
}
// Serialize serializes information the current object
func (m *WorkforceIntegrationEncryption) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetProtocol() != nil {
        cast := (*m.GetProtocol()).String()
        err := writer.WriteStringValue("protocol", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("secret", m.GetSecret())
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
func (m *WorkforceIntegrationEncryption) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WorkforceIntegrationEncryption) SetOdataType(value *string)() {
    m.odataType = value
}
// SetProtocol sets the protocol property value. Possible values are: sharedSecret, unknownFutureValue.
func (m *WorkforceIntegrationEncryption) SetProtocol(value *WorkforceIntegrationEncryptionProtocol)() {
    m.protocol = value
}
// SetSecret sets the secret property value. Encryption shared secret.
func (m *WorkforceIntegrationEncryption) SetSecret(value *string)() {
    m.secret = value
}
