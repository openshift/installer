package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DriveRecipient 
type DriveRecipient struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The alias of the domain object, for cases where an email address is unavailable (e.g. security groups).
    alias *string
    // The email address for the recipient, if the recipient has an associated email address.
    email *string
    // The unique identifier for the recipient in the directory.
    objectId *string
    // The OdataType property
    odataType *string
}
// NewDriveRecipient instantiates a new driveRecipient and sets the default values.
func NewDriveRecipient()(*DriveRecipient) {
    m := &DriveRecipient{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateDriveRecipientFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDriveRecipientFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDriveRecipient(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *DriveRecipient) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAlias gets the alias property value. The alias of the domain object, for cases where an email address is unavailable (e.g. security groups).
func (m *DriveRecipient) GetAlias()(*string) {
    return m.alias
}
// GetEmail gets the email property value. The email address for the recipient, if the recipient has an associated email address.
func (m *DriveRecipient) GetEmail()(*string) {
    return m.email
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DriveRecipient) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["alias"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAlias)
    res["email"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEmail)
    res["objectId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetObjectId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetObjectId gets the objectId property value. The unique identifier for the recipient in the directory.
func (m *DriveRecipient) GetObjectId()(*string) {
    return m.objectId
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *DriveRecipient) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *DriveRecipient) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("alias", m.GetAlias())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("objectId", m.GetObjectId())
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
func (m *DriveRecipient) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAlias sets the alias property value. The alias of the domain object, for cases where an email address is unavailable (e.g. security groups).
func (m *DriveRecipient) SetAlias(value *string)() {
    m.alias = value
}
// SetEmail sets the email property value. The email address for the recipient, if the recipient has an associated email address.
func (m *DriveRecipient) SetEmail(value *string)() {
    m.email = value
}
// SetObjectId sets the objectId property value. The unique identifier for the recipient in the directory.
func (m *DriveRecipient) SetObjectId(value *string)() {
    m.objectId = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DriveRecipient) SetOdataType(value *string)() {
    m.odataType = value
}
