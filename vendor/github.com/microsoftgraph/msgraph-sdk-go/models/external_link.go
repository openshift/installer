package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ExternalLink 
type ExternalLink struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The url of the link.
    href *string
    // The OdataType property
    odataType *string
}
// NewExternalLink instantiates a new externalLink and sets the default values.
func NewExternalLink()(*ExternalLink) {
    m := &ExternalLink{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateExternalLinkFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateExternalLinkFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewExternalLink(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ExternalLink) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ExternalLink) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["href"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetHref)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetHref gets the href property value. The url of the link.
func (m *ExternalLink) GetHref()(*string) {
    return m.href
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ExternalLink) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ExternalLink) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("href", m.GetHref())
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
func (m *ExternalLink) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetHref sets the href property value. The url of the link.
func (m *ExternalLink) SetHref(value *string)() {
    m.href = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ExternalLink) SetOdataType(value *string)() {
    m.odataType = value
}
