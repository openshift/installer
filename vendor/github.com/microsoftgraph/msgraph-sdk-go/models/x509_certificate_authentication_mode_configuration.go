package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// X509CertificateAuthenticationModeConfiguration 
type X509CertificateAuthenticationModeConfiguration struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // Rules are configured in addition to the authentication mode to bind a specific x509CertificateRuleType to an x509CertificateAuthenticationMode. For example, bind the policyOID with identifier 1.32.132.343 to x509CertificateMultiFactor authentication mode.
    rules []X509CertificateRuleable
    // The type of strong authentication mode. The possible values are: x509CertificateSingleFactor, x509CertificateMultiFactor, unknownFutureValue.
    x509CertificateAuthenticationDefaultMode *X509CertificateAuthenticationMode
}
// NewX509CertificateAuthenticationModeConfiguration instantiates a new x509CertificateAuthenticationModeConfiguration and sets the default values.
func NewX509CertificateAuthenticationModeConfiguration()(*X509CertificateAuthenticationModeConfiguration) {
    m := &X509CertificateAuthenticationModeConfiguration{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateX509CertificateAuthenticationModeConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateX509CertificateAuthenticationModeConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewX509CertificateAuthenticationModeConfiguration(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *X509CertificateAuthenticationModeConfiguration) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *X509CertificateAuthenticationModeConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["rules"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateX509CertificateRuleFromDiscriminatorValue , m.SetRules)
    res["x509CertificateAuthenticationDefaultMode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseX509CertificateAuthenticationMode , m.SetX509CertificateAuthenticationDefaultMode)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *X509CertificateAuthenticationModeConfiguration) GetOdataType()(*string) {
    return m.odataType
}
// GetRules gets the rules property value. Rules are configured in addition to the authentication mode to bind a specific x509CertificateRuleType to an x509CertificateAuthenticationMode. For example, bind the policyOID with identifier 1.32.132.343 to x509CertificateMultiFactor authentication mode.
func (m *X509CertificateAuthenticationModeConfiguration) GetRules()([]X509CertificateRuleable) {
    return m.rules
}
// GetX509CertificateAuthenticationDefaultMode gets the x509CertificateAuthenticationDefaultMode property value. The type of strong authentication mode. The possible values are: x509CertificateSingleFactor, x509CertificateMultiFactor, unknownFutureValue.
func (m *X509CertificateAuthenticationModeConfiguration) GetX509CertificateAuthenticationDefaultMode()(*X509CertificateAuthenticationMode) {
    return m.x509CertificateAuthenticationDefaultMode
}
// Serialize serializes information the current object
func (m *X509CertificateAuthenticationModeConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetRules() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetRules())
        err := writer.WriteCollectionOfObjectValues("rules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetX509CertificateAuthenticationDefaultMode() != nil {
        cast := (*m.GetX509CertificateAuthenticationDefaultMode()).String()
        err := writer.WriteStringValue("x509CertificateAuthenticationDefaultMode", &cast)
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
func (m *X509CertificateAuthenticationModeConfiguration) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *X509CertificateAuthenticationModeConfiguration) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRules sets the rules property value. Rules are configured in addition to the authentication mode to bind a specific x509CertificateRuleType to an x509CertificateAuthenticationMode. For example, bind the policyOID with identifier 1.32.132.343 to x509CertificateMultiFactor authentication mode.
func (m *X509CertificateAuthenticationModeConfiguration) SetRules(value []X509CertificateRuleable)() {
    m.rules = value
}
// SetX509CertificateAuthenticationDefaultMode sets the x509CertificateAuthenticationDefaultMode property value. The type of strong authentication mode. The possible values are: x509CertificateSingleFactor, x509CertificateMultiFactor, unknownFutureValue.
func (m *X509CertificateAuthenticationModeConfiguration) SetX509CertificateAuthenticationDefaultMode(value *X509CertificateAuthenticationMode)() {
    m.x509CertificateAuthenticationDefaultMode = value
}
