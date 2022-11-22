package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BaseDeltaFunctionResponse 
type BaseDeltaFunctionResponse struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataDeltaLink property
    odataDeltaLink *string
    // The OdataNextLink property
    odataNextLink *string
}
// NewBaseDeltaFunctionResponse instantiates a new BaseDeltaFunctionResponse and sets the default values.
func NewBaseDeltaFunctionResponse()(*BaseDeltaFunctionResponse) {
    m := &BaseDeltaFunctionResponse{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateBaseDeltaFunctionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateBaseDeltaFunctionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBaseDeltaFunctionResponse(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BaseDeltaFunctionResponse) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *BaseDeltaFunctionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.deltaLink"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataDeltaLink)
    res["@odata.nextLink"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataNextLink)
    return res
}
// GetOdataDeltaLink gets the @odata.deltaLink property value. The OdataDeltaLink property
func (m *BaseDeltaFunctionResponse) GetOdataDeltaLink()(*string) {
    return m.odataDeltaLink
}
// GetOdataNextLink gets the @odata.nextLink property value. The OdataNextLink property
func (m *BaseDeltaFunctionResponse) GetOdataNextLink()(*string) {
    return m.odataNextLink
}
// Serialize serializes information the current object
func (m *BaseDeltaFunctionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.deltaLink", m.GetOdataDeltaLink())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.nextLink", m.GetOdataNextLink())
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
func (m *BaseDeltaFunctionResponse) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataDeltaLink sets the @odata.deltaLink property value. The OdataDeltaLink property
func (m *BaseDeltaFunctionResponse) SetOdataDeltaLink(value *string)() {
    m.odataDeltaLink = value
}
// SetOdataNextLink sets the @odata.nextLink property value. The OdataNextLink property
func (m *BaseDeltaFunctionResponse) SetOdataNextLink(value *string)() {
    m.odataNextLink = value
}
