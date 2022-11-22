package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AlterationResponse 
type AlterationResponse struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // Defines the original user query string.
    originalQueryString *string
    // Defines the details of the alteration information for the spelling correction.
    queryAlteration SearchAlterationable
    // Defines the type of the spelling correction. Possible values are: suggestion, modification.
    queryAlterationType *SearchAlterationType
}
// NewAlterationResponse instantiates a new alterationResponse and sets the default values.
func NewAlterationResponse()(*AlterationResponse) {
    m := &AlterationResponse{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAlterationResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAlterationResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAlterationResponse(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AlterationResponse) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AlterationResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["originalQueryString"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOriginalQueryString)
    res["queryAlteration"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSearchAlterationFromDiscriminatorValue , m.SetQueryAlteration)
    res["queryAlterationType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseSearchAlterationType , m.SetQueryAlterationType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AlterationResponse) GetOdataType()(*string) {
    return m.odataType
}
// GetOriginalQueryString gets the originalQueryString property value. Defines the original user query string.
func (m *AlterationResponse) GetOriginalQueryString()(*string) {
    return m.originalQueryString
}
// GetQueryAlteration gets the queryAlteration property value. Defines the details of the alteration information for the spelling correction.
func (m *AlterationResponse) GetQueryAlteration()(SearchAlterationable) {
    return m.queryAlteration
}
// GetQueryAlterationType gets the queryAlterationType property value. Defines the type of the spelling correction. Possible values are: suggestion, modification.
func (m *AlterationResponse) GetQueryAlterationType()(*SearchAlterationType) {
    return m.queryAlterationType
}
// Serialize serializes information the current object
func (m *AlterationResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("originalQueryString", m.GetOriginalQueryString())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("queryAlteration", m.GetQueryAlteration())
        if err != nil {
            return err
        }
    }
    if m.GetQueryAlterationType() != nil {
        cast := (*m.GetQueryAlterationType()).String()
        err := writer.WriteStringValue("queryAlterationType", &cast)
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
func (m *AlterationResponse) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AlterationResponse) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOriginalQueryString sets the originalQueryString property value. Defines the original user query string.
func (m *AlterationResponse) SetOriginalQueryString(value *string)() {
    m.originalQueryString = value
}
// SetQueryAlteration sets the queryAlteration property value. Defines the details of the alteration information for the spelling correction.
func (m *AlterationResponse) SetQueryAlteration(value SearchAlterationable)() {
    m.queryAlteration = value
}
// SetQueryAlterationType sets the queryAlterationType property value. Defines the type of the spelling correction. Possible values are: suggestion, modification.
func (m *AlterationResponse) SetQueryAlterationType(value *SearchAlterationType)() {
    m.queryAlterationType = value
}
