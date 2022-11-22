package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OptionalClaim 
type OptionalClaim struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
    additionalProperties []string
    // If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
    essential *bool
    // The name of the optional claim.
    name *string
    // The OdataType property
    odataType *string
    // The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.
    source *string
}
// NewOptionalClaim instantiates a new optionalClaim and sets the default values.
func NewOptionalClaim()(*OptionalClaim) {
    m := &OptionalClaim{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateOptionalClaimFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOptionalClaimFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOptionalClaim(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *OptionalClaim) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAdditionalProperties gets the additionalProperties property value. Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
func (m *OptionalClaim) GetAdditionalProperties()([]string) {
    return m.additionalProperties
}
// GetEssential gets the essential property value. If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
func (m *OptionalClaim) GetEssential()(*bool) {
    return m.essential
}
// GetFieldDeserializers the deserialization information for the current model
func (m *OptionalClaim) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["additionalProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetAdditionalProperties)
    res["essential"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEssential)
    res["name"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetName)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["source"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSource)
    return res
}
// GetName gets the name property value. The name of the optional claim.
func (m *OptionalClaim) GetName()(*string) {
    return m.name
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *OptionalClaim) GetOdataType()(*string) {
    return m.odataType
}
// GetSource gets the source property value. The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.
func (m *OptionalClaim) GetSource()(*string) {
    return m.source
}
// Serialize serializes information the current object
func (m *OptionalClaim) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAdditionalProperties() != nil {
        err := writer.WriteCollectionOfStringValues("additionalProperties", m.GetAdditionalProperties())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("essential", m.GetEssential())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
        err := writer.WriteStringValue("source", m.GetSource())
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
func (m *OptionalClaim) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAdditionalProperties sets the additionalProperties property value. Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
func (m *OptionalClaim) SetAdditionalProperties(value []string)() {
    m.additionalProperties = value
}
// SetEssential sets the essential property value. If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
func (m *OptionalClaim) SetEssential(value *bool)() {
    m.essential = value
}
// SetName sets the name property value. The name of the optional claim.
func (m *OptionalClaim) SetName(value *string)() {
    m.name = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *OptionalClaim) SetOdataType(value *string)() {
    m.odataType = value
}
// SetSource sets the source property value. The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.
func (m *OptionalClaim) SetSource(value *string)() {
    m.source = value
}
