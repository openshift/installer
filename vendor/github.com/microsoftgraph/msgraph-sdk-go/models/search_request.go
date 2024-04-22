package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// SearchRequest 
type SearchRequest struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSearchRequest instantiates a new searchRequest and sets the default values.
func NewSearchRequest()(*SearchRequest) {
    m := &SearchRequest{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSearchRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSearchRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSearchRequest(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SearchRequest) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAggregationFilters gets the aggregationFilters property value. The aggregationFilters property
func (m *SearchRequest) GetAggregationFilters()([]string) {
    val, err := m.GetBackingStore().Get("aggregationFilters")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetAggregations gets the aggregations property value. The aggregations property
func (m *SearchRequest) GetAggregations()([]AggregationOptionable) {
    val, err := m.GetBackingStore().Get("aggregations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AggregationOptionable)
    }
    return nil
}
// GetBackingStore gets the backingStore property value. Stores model information.
func (m *SearchRequest) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetContentSources gets the contentSources property value. The contentSources property
func (m *SearchRequest) GetContentSources()([]string) {
    val, err := m.GetBackingStore().Get("contentSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetEnableTopResults gets the enableTopResults property value. The enableTopResults property
func (m *SearchRequest) GetEnableTopResults()(*bool) {
    val, err := m.GetBackingStore().Get("enableTopResults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEntityTypes gets the entityTypes property value. The entityTypes property
func (m *SearchRequest) GetEntityTypes()([]EntityType) {
    val, err := m.GetBackingStore().Get("entityTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EntityType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SearchRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["aggregationFilters"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetAggregationFilters(res)
        }
        return nil
    }
    res["aggregations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAggregationOptionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AggregationOptionable, len(val))
            for i, v := range val {
                res[i] = v.(AggregationOptionable)
            }
            m.SetAggregations(res)
        }
        return nil
    }
    res["contentSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetContentSources(res)
        }
        return nil
    }
    res["enableTopResults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableTopResults(val)
        }
        return nil
    }
    res["entityTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseEntityType)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EntityType, len(val))
            for i, v := range val {
                res[i] = *(v.(*EntityType))
            }
            m.SetEntityTypes(res)
        }
        return nil
    }
    res["fields"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetFields(res)
        }
        return nil
    }
    res["from"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFrom(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["query"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSearchQueryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuery(val.(SearchQueryable))
        }
        return nil
    }
    res["queryAlterationOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSearchAlterationOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQueryAlterationOptions(val.(SearchAlterationOptionsable))
        }
        return nil
    }
    res["region"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegion(val)
        }
        return nil
    }
    res["resultTemplateOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResultTemplateOptionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResultTemplateOptions(val.(ResultTemplateOptionable))
        }
        return nil
    }
    res["sharePointOneDriveOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharePointOneDriveOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharePointOneDriveOptions(val.(SharePointOneDriveOptionsable))
        }
        return nil
    }
    res["size"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSize(val)
        }
        return nil
    }
    res["sortProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSortPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SortPropertyable, len(val))
            for i, v := range val {
                res[i] = v.(SortPropertyable)
            }
            m.SetSortProperties(res)
        }
        return nil
    }
    return res
}
// GetFields gets the fields property value. The fields property
func (m *SearchRequest) GetFields()([]string) {
    val, err := m.GetBackingStore().Get("fields")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFrom gets the from property value. The from property
func (m *SearchRequest) GetFrom()(*int32) {
    val, err := m.GetBackingStore().Get("from")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SearchRequest) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQuery gets the query property value. The query property
func (m *SearchRequest) GetQuery()(SearchQueryable) {
    val, err := m.GetBackingStore().Get("query")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SearchQueryable)
    }
    return nil
}
// GetQueryAlterationOptions gets the queryAlterationOptions property value. The queryAlterationOptions property
func (m *SearchRequest) GetQueryAlterationOptions()(SearchAlterationOptionsable) {
    val, err := m.GetBackingStore().Get("queryAlterationOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SearchAlterationOptionsable)
    }
    return nil
}
// GetRegion gets the region property value. The region property
func (m *SearchRequest) GetRegion()(*string) {
    val, err := m.GetBackingStore().Get("region")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResultTemplateOptions gets the resultTemplateOptions property value. The resultTemplateOptions property
func (m *SearchRequest) GetResultTemplateOptions()(ResultTemplateOptionable) {
    val, err := m.GetBackingStore().Get("resultTemplateOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResultTemplateOptionable)
    }
    return nil
}
// GetSharePointOneDriveOptions gets the sharePointOneDriveOptions property value. The sharePointOneDriveOptions property
func (m *SearchRequest) GetSharePointOneDriveOptions()(SharePointOneDriveOptionsable) {
    val, err := m.GetBackingStore().Get("sharePointOneDriveOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharePointOneDriveOptionsable)
    }
    return nil
}
// GetSize gets the size property value. The size property
func (m *SearchRequest) GetSize()(*int32) {
    val, err := m.GetBackingStore().Get("size")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSortProperties gets the sortProperties property value. The sortProperties property
func (m *SearchRequest) GetSortProperties()([]SortPropertyable) {
    val, err := m.GetBackingStore().Get("sortProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SortPropertyable)
    }
    return nil
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
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAggregations()))
        for i, v := range m.GetAggregations() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
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
        err := writer.WriteStringValue("region", m.GetRegion())
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
        err := writer.WriteObjectValue("sharePointOneDriveOptions", m.GetSharePointOneDriveOptions())
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
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSortProperties()))
        for i, v := range m.GetSortProperties() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
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
func (m *SearchRequest) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAggregationFilters sets the aggregationFilters property value. The aggregationFilters property
func (m *SearchRequest) SetAggregationFilters(value []string)() {
    err := m.GetBackingStore().Set("aggregationFilters", value)
    if err != nil {
        panic(err)
    }
}
// SetAggregations sets the aggregations property value. The aggregations property
func (m *SearchRequest) SetAggregations(value []AggregationOptionable)() {
    err := m.GetBackingStore().Set("aggregations", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the backingStore property value. Stores model information.
func (m *SearchRequest) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetContentSources sets the contentSources property value. The contentSources property
func (m *SearchRequest) SetContentSources(value []string)() {
    err := m.GetBackingStore().Set("contentSources", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableTopResults sets the enableTopResults property value. The enableTopResults property
func (m *SearchRequest) SetEnableTopResults(value *bool)() {
    err := m.GetBackingStore().Set("enableTopResults", value)
    if err != nil {
        panic(err)
    }
}
// SetEntityTypes sets the entityTypes property value. The entityTypes property
func (m *SearchRequest) SetEntityTypes(value []EntityType)() {
    err := m.GetBackingStore().Set("entityTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetFields sets the fields property value. The fields property
func (m *SearchRequest) SetFields(value []string)() {
    err := m.GetBackingStore().Set("fields", value)
    if err != nil {
        panic(err)
    }
}
// SetFrom sets the from property value. The from property
func (m *SearchRequest) SetFrom(value *int32)() {
    err := m.GetBackingStore().Set("from", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SearchRequest) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetQuery sets the query property value. The query property
func (m *SearchRequest) SetQuery(value SearchQueryable)() {
    err := m.GetBackingStore().Set("query", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryAlterationOptions sets the queryAlterationOptions property value. The queryAlterationOptions property
func (m *SearchRequest) SetQueryAlterationOptions(value SearchAlterationOptionsable)() {
    err := m.GetBackingStore().Set("queryAlterationOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetRegion sets the region property value. The region property
func (m *SearchRequest) SetRegion(value *string)() {
    err := m.GetBackingStore().Set("region", value)
    if err != nil {
        panic(err)
    }
}
// SetResultTemplateOptions sets the resultTemplateOptions property value. The resultTemplateOptions property
func (m *SearchRequest) SetResultTemplateOptions(value ResultTemplateOptionable)() {
    err := m.GetBackingStore().Set("resultTemplateOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetSharePointOneDriveOptions sets the sharePointOneDriveOptions property value. The sharePointOneDriveOptions property
func (m *SearchRequest) SetSharePointOneDriveOptions(value SharePointOneDriveOptionsable)() {
    err := m.GetBackingStore().Set("sharePointOneDriveOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetSize sets the size property value. The size property
func (m *SearchRequest) SetSize(value *int32)() {
    err := m.GetBackingStore().Set("size", value)
    if err != nil {
        panic(err)
    }
}
// SetSortProperties sets the sortProperties property value. The sortProperties property
func (m *SearchRequest) SetSortProperties(value []SortPropertyable)() {
    err := m.GetBackingStore().Set("sortProperties", value)
    if err != nil {
        panic(err)
    }
}
// SearchRequestable 
type SearchRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAggregationFilters()([]string)
    GetAggregations()([]AggregationOptionable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetContentSources()([]string)
    GetEnableTopResults()(*bool)
    GetEntityTypes()([]EntityType)
    GetFields()([]string)
    GetFrom()(*int32)
    GetOdataType()(*string)
    GetQuery()(SearchQueryable)
    GetQueryAlterationOptions()(SearchAlterationOptionsable)
    GetRegion()(*string)
    GetResultTemplateOptions()(ResultTemplateOptionable)
    GetSharePointOneDriveOptions()(SharePointOneDriveOptionsable)
    GetSize()(*int32)
    GetSortProperties()([]SortPropertyable)
    SetAggregationFilters(value []string)()
    SetAggregations(value []AggregationOptionable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetContentSources(value []string)()
    SetEnableTopResults(value *bool)()
    SetEntityTypes(value []EntityType)()
    SetFields(value []string)()
    SetFrom(value *int32)()
    SetOdataType(value *string)()
    SetQuery(value SearchQueryable)()
    SetQueryAlterationOptions(value SearchAlterationOptionsable)()
    SetRegion(value *string)()
    SetResultTemplateOptions(value ResultTemplateOptionable)()
    SetSharePointOneDriveOptions(value SharePointOneDriveOptionsable)()
    SetSize(value *int32)()
    SetSortProperties(value []SortPropertyable)()
}
