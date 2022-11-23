package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BucketAggregationDefinition 
type BucketAggregationDefinition struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // True to specify the sort order as descending. The default is false, with the sort order as ascending. Optional.
    isDescending *bool
    // The minimum number of items that should be present in the aggregation to be returned in a bucket. Optional.
    minimumCount *int32
    // The OdataType property
    odataType *string
    // A filter to define a matching criteria. The key should start with the specified prefix to be returned in the response. Optional.
    prefixFilter *string
    // Specifies the manual ranges to compute the aggregations. This is only valid for non-string refiners of date or numeric type. Optional.
    ranges []BucketAggregationRangeable
    // The sortBy property
    sortBy *BucketAggregationSortProperty
}
// NewBucketAggregationDefinition instantiates a new bucketAggregationDefinition and sets the default values.
func NewBucketAggregationDefinition()(*BucketAggregationDefinition) {
    m := &BucketAggregationDefinition{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateBucketAggregationDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateBucketAggregationDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBucketAggregationDefinition(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BucketAggregationDefinition) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *BucketAggregationDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isDescending"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsDescending)
    res["minimumCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetMinimumCount)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["prefixFilter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPrefixFilter)
    res["ranges"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateBucketAggregationRangeFromDiscriminatorValue , m.SetRanges)
    res["sortBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseBucketAggregationSortProperty , m.SetSortBy)
    return res
}
// GetIsDescending gets the isDescending property value. True to specify the sort order as descending. The default is false, with the sort order as ascending. Optional.
func (m *BucketAggregationDefinition) GetIsDescending()(*bool) {
    return m.isDescending
}
// GetMinimumCount gets the minimumCount property value. The minimum number of items that should be present in the aggregation to be returned in a bucket. Optional.
func (m *BucketAggregationDefinition) GetMinimumCount()(*int32) {
    return m.minimumCount
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *BucketAggregationDefinition) GetOdataType()(*string) {
    return m.odataType
}
// GetPrefixFilter gets the prefixFilter property value. A filter to define a matching criteria. The key should start with the specified prefix to be returned in the response. Optional.
func (m *BucketAggregationDefinition) GetPrefixFilter()(*string) {
    return m.prefixFilter
}
// GetRanges gets the ranges property value. Specifies the manual ranges to compute the aggregations. This is only valid for non-string refiners of date or numeric type. Optional.
func (m *BucketAggregationDefinition) GetRanges()([]BucketAggregationRangeable) {
    return m.ranges
}
// GetSortBy gets the sortBy property value. The sortBy property
func (m *BucketAggregationDefinition) GetSortBy()(*BucketAggregationSortProperty) {
    return m.sortBy
}
// Serialize serializes information the current object
func (m *BucketAggregationDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isDescending", m.GetIsDescending())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("minimumCount", m.GetMinimumCount())
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
        err := writer.WriteStringValue("prefixFilter", m.GetPrefixFilter())
        if err != nil {
            return err
        }
    }
    if m.GetRanges() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetRanges())
        err := writer.WriteCollectionOfObjectValues("ranges", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSortBy() != nil {
        cast := (*m.GetSortBy()).String()
        err := writer.WriteStringValue("sortBy", &cast)
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
func (m *BucketAggregationDefinition) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIsDescending sets the isDescending property value. True to specify the sort order as descending. The default is false, with the sort order as ascending. Optional.
func (m *BucketAggregationDefinition) SetIsDescending(value *bool)() {
    m.isDescending = value
}
// SetMinimumCount sets the minimumCount property value. The minimum number of items that should be present in the aggregation to be returned in a bucket. Optional.
func (m *BucketAggregationDefinition) SetMinimumCount(value *int32)() {
    m.minimumCount = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BucketAggregationDefinition) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPrefixFilter sets the prefixFilter property value. A filter to define a matching criteria. The key should start with the specified prefix to be returned in the response. Optional.
func (m *BucketAggregationDefinition) SetPrefixFilter(value *string)() {
    m.prefixFilter = value
}
// SetRanges sets the ranges property value. Specifies the manual ranges to compute the aggregations. This is only valid for non-string refiners of date or numeric type. Optional.
func (m *BucketAggregationDefinition) SetRanges(value []BucketAggregationRangeable)() {
    m.ranges = value
}
// SetSortBy sets the sortBy property value. The sortBy property
func (m *BucketAggregationDefinition) SetSortBy(value *BucketAggregationSortProperty)() {
    m.sortBy = value
}
