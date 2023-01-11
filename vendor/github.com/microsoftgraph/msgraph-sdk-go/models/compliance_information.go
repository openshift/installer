package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ComplianceInformation 
type ComplianceInformation struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Collection of the certification controls associated with certification
    certificationControls []CertificationControlable
    // Compliance certification name (for example, ISO 27018:2014, GDPR, FedRAMP, NIST 800-171)
    certificationName *string
    // The OdataType property
    odataType *string
}
// NewComplianceInformation instantiates a new complianceInformation and sets the default values.
func NewComplianceInformation()(*ComplianceInformation) {
    m := &ComplianceInformation{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateComplianceInformationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateComplianceInformationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewComplianceInformation(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ComplianceInformation) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCertificationControls gets the certificationControls property value. Collection of the certification controls associated with certification
func (m *ComplianceInformation) GetCertificationControls()([]CertificationControlable) {
    return m.certificationControls
}
// GetCertificationName gets the certificationName property value. Compliance certification name (for example, ISO 27018:2014, GDPR, FedRAMP, NIST 800-171)
func (m *ComplianceInformation) GetCertificationName()(*string) {
    return m.certificationName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ComplianceInformation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["certificationControls"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCertificationControlFromDiscriminatorValue , m.SetCertificationControls)
    res["certificationName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCertificationName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ComplianceInformation) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ComplianceInformation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetCertificationControls() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCertificationControls())
        err := writer.WriteCollectionOfObjectValues("certificationControls", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("certificationName", m.GetCertificationName())
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
func (m *ComplianceInformation) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCertificationControls sets the certificationControls property value. Collection of the certification controls associated with certification
func (m *ComplianceInformation) SetCertificationControls(value []CertificationControlable)() {
    m.certificationControls = value
}
// SetCertificationName sets the certificationName property value. Compliance certification name (for example, ISO 27018:2014, GDPR, FedRAMP, NIST 800-171)
func (m *ComplianceInformation) SetCertificationName(value *string)() {
    m.certificationName = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ComplianceInformation) SetOdataType(value *string)() {
    m.odataType = value
}
