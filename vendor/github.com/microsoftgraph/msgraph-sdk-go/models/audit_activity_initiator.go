package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AuditActivityInitiator 
type AuditActivityInitiator struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // If the resource initiating the activity is an app, this property indicates all the app related information like appId, Name, servicePrincipalId, Name.
    app AppIdentityable
    // The OdataType property
    odataType *string
    // If the resource initiating the activity is a user, this property Indicates all the user related information like userId, Name, UserPrinicpalName.
    user UserIdentityable
}
// NewAuditActivityInitiator instantiates a new auditActivityInitiator and sets the default values.
func NewAuditActivityInitiator()(*AuditActivityInitiator) {
    m := &AuditActivityInitiator{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAuditActivityInitiatorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAuditActivityInitiatorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuditActivityInitiator(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AuditActivityInitiator) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetApp gets the app property value. If the resource initiating the activity is an app, this property indicates all the app related information like appId, Name, servicePrincipalId, Name.
func (m *AuditActivityInitiator) GetApp()(AppIdentityable) {
    return m.app
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AuditActivityInitiator) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["app"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateAppIdentityFromDiscriminatorValue , m.SetApp)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["user"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUserIdentityFromDiscriminatorValue , m.SetUser)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AuditActivityInitiator) GetOdataType()(*string) {
    return m.odataType
}
// GetUser gets the user property value. If the resource initiating the activity is a user, this property Indicates all the user related information like userId, Name, UserPrinicpalName.
func (m *AuditActivityInitiator) GetUser()(UserIdentityable) {
    return m.user
}
// Serialize serializes information the current object
func (m *AuditActivityInitiator) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("app", m.GetApp())
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
        err := writer.WriteObjectValue("user", m.GetUser())
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
func (m *AuditActivityInitiator) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetApp sets the app property value. If the resource initiating the activity is an app, this property indicates all the app related information like appId, Name, servicePrincipalId, Name.
func (m *AuditActivityInitiator) SetApp(value AppIdentityable)() {
    m.app = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AuditActivityInitiator) SetOdataType(value *string)() {
    m.odataType = value
}
// SetUser sets the user property value. If the resource initiating the activity is a user, this property Indicates all the user related information like userId, Name, UserPrinicpalName.
func (m *AuditActivityInitiator) SetUser(value UserIdentityable)() {
    m.user = value
}
