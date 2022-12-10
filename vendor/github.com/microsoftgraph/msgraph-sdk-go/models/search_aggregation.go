package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SearchAggregation 
type SearchAggregation struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The buckets property
    buckets []SearchBucketable
    // The field property
    field *string
    // The OdataType property
    odataType *string
}
// NewSearchAggregation instantiates a new searchAggregation and sets the default values.
func NewSearchAggregation()(*SearchAggregation) {
    m := &SearchAggregation{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSearchAggregationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSearchAggregationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSearchAggregation(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SearchAggregation) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetBuckets gets the buckets property value. The buckets property
func (m *SearchAggregation) GetBuckets()([]SearchBucketable) {
    return m.buckets
}
// GetField gets the field property value. The field property
func (m *SearchAggregation) GetField()(*string) {
    return m.field
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SearchAggregation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["buckets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSearchBucketFromDiscriminatorValue , m.SetBuckets)
    res["field"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetField)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SearchAggregation) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *SearchAggregation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetBuckets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetBuckets())
        err := writer.WriteCollectionOfObjectValues("buckets", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("field", m.GetField())
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
func (m *SearchAggregation) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetBuckets sets the buckets property value. The buckets property
func (m *SearchAggregation) SetBuckets(value []SearchBucketable)() {
    m.buckets = value
}
// SetField sets the field property value. The field property
func (m *SearchAggregation) SetField(value *string)() {
    m.field = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SearchAggregation) SetOdataType(value *string)() {
    m.odataType = value
}
