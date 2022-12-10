package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AuthorizationInfo 
type AuthorizationInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The certificateUserIds property
    certificateUserIds []string
    // The OdataType property
    odataType *string
}
// NewAuthorizationInfo instantiates a new authorizationInfo and sets the default values.
func NewAuthorizationInfo()(*AuthorizationInfo) {
    m := &AuthorizationInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAuthorizationInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAuthorizationInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuthorizationInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AuthorizationInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCertificateUserIds gets the certificateUserIds property value. The certificateUserIds property
func (m *AuthorizationInfo) GetCertificateUserIds()([]string) {
    return m.certificateUserIds
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AuthorizationInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["certificateUserIds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetCertificateUserIds)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AuthorizationInfo) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *AuthorizationInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetCertificateUserIds() != nil {
        err := writer.WriteCollectionOfStringValues("certificateUserIds", m.GetCertificateUserIds())
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
func (m *AuthorizationInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCertificateUserIds sets the certificateUserIds property value. The certificateUserIds property
func (m *AuthorizationInfo) SetCertificateUserIds(value []string)() {
    m.certificateUserIds = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AuthorizationInfo) SetOdataType(value *string)() {
    m.odataType = value
}
