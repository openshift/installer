package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AttackSimulationUser 
type AttackSimulationUser struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Display name of the user.
    displayName *string
    // Email address of the user.
    email *string
    // The OdataType property
    odataType *string
    // This is the id property value of the user resource that represents the user in the Azure Active Directory tenant.
    userId *string
}
// NewAttackSimulationUser instantiates a new attackSimulationUser and sets the default values.
func NewAttackSimulationUser()(*AttackSimulationUser) {
    m := &AttackSimulationUser{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAttackSimulationUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAttackSimulationUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttackSimulationUser(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AttackSimulationUser) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDisplayName gets the displayName property value. Display name of the user.
func (m *AttackSimulationUser) GetDisplayName()(*string) {
    return m.displayName
}
// GetEmail gets the email property value. Email address of the user.
func (m *AttackSimulationUser) GetEmail()(*string) {
    return m.email
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AttackSimulationUser) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["email"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEmail)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AttackSimulationUser) GetOdataType()(*string) {
    return m.odataType
}
// GetUserId gets the userId property value. This is the id property value of the user resource that represents the user in the Azure Active Directory tenant.
func (m *AttackSimulationUser) GetUserId()(*string) {
    return m.userId
}
// Serialize serializes information the current object
func (m *AttackSimulationUser) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
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
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("userId", m.GetUserId())
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
func (m *AttackSimulationUser) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDisplayName sets the displayName property value. Display name of the user.
func (m *AttackSimulationUser) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetEmail sets the email property value. Email address of the user.
func (m *AttackSimulationUser) SetEmail(value *string)() {
    m.email = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AttackSimulationUser) SetOdataType(value *string)() {
    m.odataType = value
}
// SetUserId sets the userId property value. This is the id property value of the user resource that represents the user in the Azure Active Directory tenant.
func (m *AttackSimulationUser) SetUserId(value *string)() {
    m.userId = value
}
