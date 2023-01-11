package translateexchangeids

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// TranslateExchangeIdsPostRequestBody provides operations to call the translateExchangeIds method.
type TranslateExchangeIdsPostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The InputIds property
    inputIds []string
    // The SourceIdType property
    sourceIdType *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat
    // The TargetIdType property
    targetIdType *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat
}
// NewTranslateExchangeIdsPostRequestBody instantiates a new translateExchangeIdsPostRequestBody and sets the default values.
func NewTranslateExchangeIdsPostRequestBody()(*TranslateExchangeIdsPostRequestBody) {
    m := &TranslateExchangeIdsPostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateTranslateExchangeIdsPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTranslateExchangeIdsPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTranslateExchangeIdsPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *TranslateExchangeIdsPostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TranslateExchangeIdsPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["inputIds"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetInputIds)
    res["sourceIdType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseExchangeIdFormat , m.SetSourceIdType)
    res["targetIdType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseExchangeIdFormat , m.SetTargetIdType)
    return res
}
// GetInputIds gets the inputIds property value. The InputIds property
func (m *TranslateExchangeIdsPostRequestBody) GetInputIds()([]string) {
    return m.inputIds
}
// GetSourceIdType gets the sourceIdType property value. The SourceIdType property
func (m *TranslateExchangeIdsPostRequestBody) GetSourceIdType()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat) {
    return m.sourceIdType
}
// GetTargetIdType gets the targetIdType property value. The TargetIdType property
func (m *TranslateExchangeIdsPostRequestBody) GetTargetIdType()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat) {
    return m.targetIdType
}
// Serialize serializes information the current object
func (m *TranslateExchangeIdsPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetInputIds() != nil {
        err := writer.WriteCollectionOfStringValues("inputIds", m.GetInputIds())
        if err != nil {
            return err
        }
    }
    if m.GetSourceIdType() != nil {
        cast := (*m.GetSourceIdType()).String()
        err := writer.WriteStringValue("sourceIdType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetTargetIdType() != nil {
        cast := (*m.GetTargetIdType()).String()
        err := writer.WriteStringValue("targetIdType", &cast)
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
func (m *TranslateExchangeIdsPostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetInputIds sets the inputIds property value. The InputIds property
func (m *TranslateExchangeIdsPostRequestBody) SetInputIds(value []string)() {
    m.inputIds = value
}
// SetSourceIdType sets the sourceIdType property value. The SourceIdType property
func (m *TranslateExchangeIdsPostRequestBody) SetSourceIdType(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)() {
    m.sourceIdType = value
}
// SetTargetIdType sets the targetIdType property value. The TargetIdType property
func (m *TranslateExchangeIdsPostRequestBody) SetTargetIdType(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ExchangeIdFormat)() {
    m.targetIdType = value
}
