package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Pkcs12CertificateInformation 
type Pkcs12CertificateInformation struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Represents whether the certificate is the active certificate to be used for calling the API connector. The active certificate is the most recently uploaded certificate which is not yet expired but whose notBefore time is in the past.
    isActive *bool
    // The certificate's expiry. This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)
    notAfter *int64
    // The certificate's issue time (not before). This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)
    notBefore *int64
    // The OdataType property
    odataType *string
    // The certificate thumbprint.
    thumbprint *string
}
// NewPkcs12CertificateInformation instantiates a new pkcs12CertificateInformation and sets the default values.
func NewPkcs12CertificateInformation()(*Pkcs12CertificateInformation) {
    m := &Pkcs12CertificateInformation{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePkcs12CertificateInformationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePkcs12CertificateInformationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPkcs12CertificateInformation(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Pkcs12CertificateInformation) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Pkcs12CertificateInformation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isActive"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsActive)
    res["notAfter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt64Value(m.SetNotAfter)
    res["notBefore"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt64Value(m.SetNotBefore)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["thumbprint"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetThumbprint)
    return res
}
// GetIsActive gets the isActive property value. Represents whether the certificate is the active certificate to be used for calling the API connector. The active certificate is the most recently uploaded certificate which is not yet expired but whose notBefore time is in the past.
func (m *Pkcs12CertificateInformation) GetIsActive()(*bool) {
    return m.isActive
}
// GetNotAfter gets the notAfter property value. The certificate's expiry. This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)
func (m *Pkcs12CertificateInformation) GetNotAfter()(*int64) {
    return m.notAfter
}
// GetNotBefore gets the notBefore property value. The certificate's issue time (not before). This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)
func (m *Pkcs12CertificateInformation) GetNotBefore()(*int64) {
    return m.notBefore
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Pkcs12CertificateInformation) GetOdataType()(*string) {
    return m.odataType
}
// GetThumbprint gets the thumbprint property value. The certificate thumbprint.
func (m *Pkcs12CertificateInformation) GetThumbprint()(*string) {
    return m.thumbprint
}
// Serialize serializes information the current object
func (m *Pkcs12CertificateInformation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isActive", m.GetIsActive())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("notAfter", m.GetNotAfter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("notBefore", m.GetNotBefore())
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
        err := writer.WriteStringValue("thumbprint", m.GetThumbprint())
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
func (m *Pkcs12CertificateInformation) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIsActive sets the isActive property value. Represents whether the certificate is the active certificate to be used for calling the API connector. The active certificate is the most recently uploaded certificate which is not yet expired but whose notBefore time is in the past.
func (m *Pkcs12CertificateInformation) SetIsActive(value *bool)() {
    m.isActive = value
}
// SetNotAfter sets the notAfter property value. The certificate's expiry. This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)
func (m *Pkcs12CertificateInformation) SetNotAfter(value *int64)() {
    m.notAfter = value
}
// SetNotBefore sets the notBefore property value. The certificate's issue time (not before). This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)
func (m *Pkcs12CertificateInformation) SetNotBefore(value *int64)() {
    m.notBefore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Pkcs12CertificateInformation) SetOdataType(value *string)() {
    m.odataType = value
}
// SetThumbprint sets the thumbprint property value. The certificate thumbprint.
func (m *Pkcs12CertificateInformation) SetThumbprint(value *string)() {
    m.thumbprint = value
}
