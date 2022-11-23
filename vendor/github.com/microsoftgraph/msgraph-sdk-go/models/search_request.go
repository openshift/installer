package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SearchRequest 
type SearchRequest struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The aggregationFilters property
    aggregationFilters []string
    // The aggregations property
    aggregations []AggregationOptionable
    // The contentSources property
    contentSources []string
    // The enableTopResults property
    enableTopResults *bool
    // The entityTypes property
    entityTypes []EntityType
    // The fields property
    fields []string
    // The from property
    from *int32
    // The OdataType property
    odataType *string
    // The query property
    query SearchQueryable
    // The queryAlterationOptions property
    queryAlterationOptions SearchAlterationOptionsable
    // The resultTemplateOptions property
    resultTemplateOptions ResultTemplateOptionable
    // The size property
    size *int32
    // The sortProperties property
    sortProperties []SortPropertyable
}
// NewSearchRequest instantiates a new searchRequest and sets the default values.
func NewSearchRequest()(*SearchRequest) {
    m := &SearchRequest{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSearchRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSearchRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSearchRequest(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SearchRequest) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAggregationFilters gets the aggregationFilters property value. The aggregationFilters property
func (m *SearchRequest) GetAggregationFilters()([]string) {
    return m.aggregationFilters
}
// GetAggregations gets the aggregations property value. The aggregations property
func (m *SearchRequest) GetAggregations()([]AggregationOptionable) {
    return m.aggregations
}
// GetContentSources gets the contentSources property value. The contentSources property
func (m *SearchRequest) GetContentSources()([]string) {
    return m.contentSources
}
// GetEnableTopResults gets the enableTopResults property value. The enableTopResults property
func (m *SearchRequest) GetEnableTopResults()(*bool) {
    return m.enableTopResults
}
// GetEntityTypes gets the entityTypes property value. The entityTypes property
func (m *SearchRequest) GetEntityTypes()([]EntityType) {
    return m.entityTypes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SearchRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["aggregationFilters"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetAggregationFilters)
    res["aggregations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAggregationOptionFromDiscriminatorValue , m.SetAggregations)
    res["contentSources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetContentSources)
    res["enableTopResults"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEnableTopResults)
    res["entityTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParseEntityType , m.SetEntityTypes)
    res["fields"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetFields)
    res["from"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetFrom)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["query"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSearchQueryFromDiscriminatorValue , m.SetQuery)
    res["queryAlterationOptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSearchAlterationOptionsFromDiscriminatorValue , m.SetQueryAlterationOptions)
    res["resultTemplateOptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateResultTemplateOptionFromDiscriminatorValue , m.SetResultTemplateOptions)
    res["size"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetSize)
    res["sortProperties"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateSortPropertyFromDiscriminatorValue , m.SetSortProperties)
    return res
}
// GetFields gets the fields property value. The fields property
func (m *SearchRequest) GetFields()([]string) {
    return m.fields
}
// GetFrom gets the from property value. The from property
func (m *SearchRequest) GetFrom()(*int32) {
    return m.from
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SearchRequest) GetOdataType()(*string) {
    return m.odataType
}
// GetQuery gets the query property value. The query property
func (m *SearchRequest) GetQuery()(SearchQueryable) {
    return m.query
}
// GetQueryAlterationOptions gets the queryAlterationOptions property value. The queryAlterationOptions property
func (m *SearchRequest) GetQueryAlterationOptions()(SearchAlterationOptionsable) {
    return m.queryAlterationOptions
}
// GetResultTemplateOptions gets the resultTemplateOptions property value. The resultTemplateOptions property
func (m *SearchRequest) GetResultTemplateOptions()(ResultTemplateOptionable) {
    return m.resultTemplateOptions
}
// GetSize gets the size property value. The size property
func (m *SearchRequest) GetSize()(*int32) {
    return m.size
}
// GetSortProperties gets the sortProperties property value. The sortProperties property
func (m *SearchRequest) GetSortProperties()([]SortPropertyable) {
    return m.sortProperties
}
// Serialize serializes information the current object
func (m *SearchRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAggregationFilters() != nil {
        err := writer.WriteCollectionOfStringValues("aggregationFilters", m.GetAggregationFilters())
        if err != nil {
            return err
        }
    }
    if m.GetAggregations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAggregations())
        err := writer.WriteCollectionOfObjectValues("aggregations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetContentSources() != nil {
        err := writer.WriteCollectionOfStringValues("contentSources", m.GetContentSources())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableTopResults", m.GetEnableTopResults())
        if err != nil {
            return err
        }
    }
    if m.GetEntityTypes() != nil {
        err := writer.WriteCollectionOfStringValues("entityTypes", SerializeEntityType(m.GetEntityTypes()))
        if err != nil {
            return err
        }
    }
    if m.GetFields() != nil {
        err := writer.WriteCollectionOfStringValues("fields", m.GetFields())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("from", m.GetFrom())
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
        err := writer.WriteObjectValue("query", m.GetQuery())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("queryAlterationOptions", m.GetQueryAlterationOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("resultTemplateOptions", m.GetResultTemplateOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("size", m.GetSize())
        if err != nil {
            return err
        }
    }
    if m.GetSortProperties() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSortProperties())
        err := writer.WriteCollectionOfObjectValues("sortProperties", cast)
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
func (m *SearchRequest) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAggregationFilters sets the aggregationFilters property value. The aggregationFilters property
func (m *SearchRequest) SetAggregationFilters(value []string)() {
    m.aggregationFilters = value
}
// SetAggregations sets the aggregations property value. The aggregations property
func (m *SearchRequest) SetAggregations(value []AggregationOptionable)() {
    m.aggregations = value
}
// SetContentSources sets the contentSources property value. The contentSources property
func (m *SearchRequest) SetContentSources(value []string)() {
    m.contentSources = value
}
// SetEnableTopResults sets the enableTopResults property value. The enableTopResults property
func (m *SearchRequest) SetEnableTopResults(value *bool)() {
    m.enableTopResults = value
}
// SetEntityTypes sets the entityTypes property value. The entityTypes property
func (m *SearchRequest) SetEntityTypes(value []EntityType)() {
    m.entityTypes = value
}
// SetFields sets the fields property value. The fields property
func (m *SearchRequest) SetFields(value []string)() {
    m.fields = value
}
// SetFrom sets the from property value. The from property
func (m *SearchRequest) SetFrom(value *int32)() {
    m.from = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SearchRequest) SetOdataType(value *string)() {
    m.odataType = value
}
// SetQuery sets the query property value. The query property
func (m *SearchRequest) SetQuery(value SearchQueryable)() {
    m.query = value
}
// SetQueryAlterationOptions sets the queryAlterationOptions property value. The queryAlterationOptions property
func (m *SearchRequest) SetQueryAlterationOptions(value SearchAlterationOptionsable)() {
    m.queryAlterationOptions = value
}
// SetResultTemplateOptions sets the resultTemplateOptions property value. The resultTemplateOptions property
func (m *SearchRequest) SetResultTemplateOptions(value ResultTemplateOptionable)() {
    m.resultTemplateOptions = value
}
// SetSize sets the size property value. The size property
func (m *SearchRequest) SetSize(value *int32)() {
    m.size = value
}
// SetSortProperties sets the sortProperties property value. The sortProperties property
func (m *SearchRequest) SetSortProperties(value []SortPropertyable)() {
    m.sortProperties = value
}
