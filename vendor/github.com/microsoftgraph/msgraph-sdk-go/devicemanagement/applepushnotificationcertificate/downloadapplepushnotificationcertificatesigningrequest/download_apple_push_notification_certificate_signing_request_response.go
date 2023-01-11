package downloadapplepushnotificationcertificatesigningrequest

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DownloadApplePushNotificationCertificateSigningRequestResponse provides operations to call the downloadApplePushNotificationCertificateSigningRequest method.
type DownloadApplePushNotificationCertificateSigningRequestResponse struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The value property
    value *string
}
// NewDownloadApplePushNotificationCertificateSigningRequestResponse instantiates a new downloadApplePushNotificationCertificateSigningRequestResponse and sets the default values.
func NewDownloadApplePushNotificationCertificateSigningRequestResponse()(*DownloadApplePushNotificationCertificateSigningRequestResponse) {
    m := &DownloadApplePushNotificationCertificateSigningRequestResponse{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateDownloadApplePushNotificationCertificateSigningRequestResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDownloadApplePushNotificationCertificateSigningRequestResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDownloadApplePushNotificationCertificateSigningRequestResponse(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *DownloadApplePushNotificationCertificateSigningRequestResponse) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DownloadApplePushNotificationCertificateSigningRequestResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetValue)
    return res
}
// GetValue gets the value property value. The value property
func (m *DownloadApplePushNotificationCertificateSigningRequestResponse) GetValue()(*string) {
    return m.value
}
// Serialize serializes information the current object
func (m *DownloadApplePushNotificationCertificateSigningRequestResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("value", m.GetValue())
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
func (m *DownloadApplePushNotificationCertificateSigningRequestResponse) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetValue sets the value property value. The value property
func (m *DownloadApplePushNotificationCertificateSigningRequestResponse) SetValue(value *string)() {
    m.value = value
}
