package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SigningCertificateUpdateStatus 
type SigningCertificateUpdateStatus struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Status of the last certificate update. Read-only. For a list of statuses, see certificateUpdateResult status.
    certificateUpdateResult *string
    // Date and time in ISO 8601 format and in UTC time when the certificate was last updated. Read-only.
    lastRunDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The OdataType property
    odataType *string
}
// NewSigningCertificateUpdateStatus instantiates a new signingCertificateUpdateStatus and sets the default values.
func NewSigningCertificateUpdateStatus()(*SigningCertificateUpdateStatus) {
    m := &SigningCertificateUpdateStatus{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSigningCertificateUpdateStatusFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSigningCertificateUpdateStatusFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSigningCertificateUpdateStatus(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SigningCertificateUpdateStatus) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCertificateUpdateResult gets the certificateUpdateResult property value. Status of the last certificate update. Read-only. For a list of statuses, see certificateUpdateResult status.
func (m *SigningCertificateUpdateStatus) GetCertificateUpdateResult()(*string) {
    return m.certificateUpdateResult
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SigningCertificateUpdateStatus) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["certificateUpdateResult"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCertificateUpdateResult)
    res["lastRunDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastRunDateTime)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetLastRunDateTime gets the lastRunDateTime property value. Date and time in ISO 8601 format and in UTC time when the certificate was last updated. Read-only.
func (m *SigningCertificateUpdateStatus) GetLastRunDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastRunDateTime
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SigningCertificateUpdateStatus) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *SigningCertificateUpdateStatus) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("certificateUpdateResult", m.GetCertificateUpdateResult())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("lastRunDateTime", m.GetLastRunDateTime())
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
func (m *SigningCertificateUpdateStatus) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCertificateUpdateResult sets the certificateUpdateResult property value. Status of the last certificate update. Read-only. For a list of statuses, see certificateUpdateResult status.
func (m *SigningCertificateUpdateStatus) SetCertificateUpdateResult(value *string)() {
    m.certificateUpdateResult = value
}
// SetLastRunDateTime sets the lastRunDateTime property value. Date and time in ISO 8601 format and in UTC time when the certificate was last updated. Read-only.
func (m *SigningCertificateUpdateStatus) SetLastRunDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastRunDateTime = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SigningCertificateUpdateStatus) SetOdataType(value *string)() {
    m.odataType = value
}
