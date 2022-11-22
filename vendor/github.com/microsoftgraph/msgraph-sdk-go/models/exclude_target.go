package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ExcludeTarget 
type ExcludeTarget struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The object identifier of an Azure Active Directory user or group.
    id *string
    // The OdataType property
    odataType *string
    // The targetType property
    targetType *AuthenticationMethodTargetType
}
// NewExcludeTarget instantiates a new excludeTarget and sets the default values.
func NewExcludeTarget()(*ExcludeTarget) {
    m := &ExcludeTarget{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateExcludeTargetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateExcludeTargetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewExcludeTarget(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ExcludeTarget) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ExcludeTarget) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["id"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["targetType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseAuthenticationMethodTargetType , m.SetTargetType)
    return res
}
// GetId gets the id property value. The object identifier of an Azure Active Directory user or group.
func (m *ExcludeTarget) GetId()(*string) {
    return m.id
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ExcludeTarget) GetOdataType()(*string) {
    return m.odataType
}
// GetTargetType gets the targetType property value. The targetType property
func (m *ExcludeTarget) GetTargetType()(*AuthenticationMethodTargetType) {
    return m.targetType
}
// Serialize serializes information the current object
func (m *ExcludeTarget) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("id", m.GetId())
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
    if m.GetTargetType() != nil {
        cast := (*m.GetTargetType()).String()
        err := writer.WriteStringValue("targetType", &cast)
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
func (m *ExcludeTarget) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetId sets the id property value. The object identifier of an Azure Active Directory user or group.
func (m *ExcludeTarget) SetId(value *string)() {
    m.id = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ExcludeTarget) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTargetType sets the targetType property value. The targetType property
func (m *ExcludeTarget) SetTargetType(value *AuthenticationMethodTargetType)() {
    m.targetType = value
}
