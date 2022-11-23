package security

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EdiscoverySearch 
type EdiscoverySearch struct {
    Search
    // Adds an additional source to the eDiscovery search.
    additionalSources []DataSourceable
    // Adds the results of the eDiscovery search to the specified reviewSet.
    addToReviewSetOperation EdiscoveryAddToReviewSetOperationable
    // Custodian sources that are included in the eDiscovery search.
    custodianSources []DataSourceable
    // When specified, the collection will span across a service for an entire workload. Possible values are: none, allTenantMailboxes, allTenantSites, allCaseCustodians, allCaseNoncustodialDataSources.
    dataSourceScopes *DataSourceScopes
    // The last estimate operation associated with the eDiscovery search.
    lastEstimateStatisticsOperation EdiscoveryEstimateOperationable
    // noncustodialDataSource sources that are included in the eDiscovery search
    noncustodialSources []EdiscoveryNoncustodialDataSourceable
}
// NewEdiscoverySearch instantiates a new EdiscoverySearch and sets the default values.
func NewEdiscoverySearch()(*EdiscoverySearch) {
    m := &EdiscoverySearch{
        Search: *NewSearch(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoverySearch";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEdiscoverySearchFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEdiscoverySearchFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoverySearch(), nil
}
// GetAdditionalSources gets the additionalSources property value. Adds an additional source to the eDiscovery search.
func (m *EdiscoverySearch) GetAdditionalSources()([]DataSourceable) {
    return m.additionalSources
}
// GetAddToReviewSetOperation gets the addToReviewSetOperation property value. Adds the results of the eDiscovery search to the specified reviewSet.
func (m *EdiscoverySearch) GetAddToReviewSetOperation()(EdiscoveryAddToReviewSetOperationable) {
    return m.addToReviewSetOperation
}
// GetCustodianSources gets the custodianSources property value. Custodian sources that are included in the eDiscovery search.
func (m *EdiscoverySearch) GetCustodianSources()([]DataSourceable) {
    return m.custodianSources
}
// GetDataSourceScopes gets the dataSourceScopes property value. When specified, the collection will span across a service for an entire workload. Possible values are: none, allTenantMailboxes, allTenantSites, allCaseCustodians, allCaseNoncustodialDataSources.
func (m *EdiscoverySearch) GetDataSourceScopes()(*DataSourceScopes) {
    return m.dataSourceScopes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EdiscoverySearch) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Search.GetFieldDeserializers()
    res["additionalSources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDataSourceFromDiscriminatorValue , m.SetAdditionalSources)
    res["addToReviewSetOperation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEdiscoveryAddToReviewSetOperationFromDiscriminatorValue , m.SetAddToReviewSetOperation)
    res["custodianSources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDataSourceFromDiscriminatorValue , m.SetCustodianSources)
    res["dataSourceScopes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDataSourceScopes , m.SetDataSourceScopes)
    res["lastEstimateStatisticsOperation"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEdiscoveryEstimateOperationFromDiscriminatorValue , m.SetLastEstimateStatisticsOperation)
    res["noncustodialSources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEdiscoveryNoncustodialDataSourceFromDiscriminatorValue , m.SetNoncustodialSources)
    return res
}
// GetLastEstimateStatisticsOperation gets the lastEstimateStatisticsOperation property value. The last estimate operation associated with the eDiscovery search.
func (m *EdiscoverySearch) GetLastEstimateStatisticsOperation()(EdiscoveryEstimateOperationable) {
    return m.lastEstimateStatisticsOperation
}
// GetNoncustodialSources gets the noncustodialSources property value. noncustodialDataSource sources that are included in the eDiscovery search
func (m *EdiscoverySearch) GetNoncustodialSources()([]EdiscoveryNoncustodialDataSourceable) {
    return m.noncustodialSources
}
// Serialize serializes information the current object
func (m *EdiscoverySearch) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Search.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAdditionalSources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAdditionalSources())
        err = writer.WriteCollectionOfObjectValues("additionalSources", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("addToReviewSetOperation", m.GetAddToReviewSetOperation())
        if err != nil {
            return err
        }
    }
    if m.GetCustodianSources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCustodianSources())
        err = writer.WriteCollectionOfObjectValues("custodianSources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDataSourceScopes() != nil {
        cast := (*m.GetDataSourceScopes()).String()
        err = writer.WriteStringValue("dataSourceScopes", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastEstimateStatisticsOperation", m.GetLastEstimateStatisticsOperation())
        if err != nil {
            return err
        }
    }
    if m.GetNoncustodialSources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNoncustodialSources())
        err = writer.WriteCollectionOfObjectValues("noncustodialSources", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalSources sets the additionalSources property value. Adds an additional source to the eDiscovery search.
func (m *EdiscoverySearch) SetAdditionalSources(value []DataSourceable)() {
    m.additionalSources = value
}
// SetAddToReviewSetOperation sets the addToReviewSetOperation property value. Adds the results of the eDiscovery search to the specified reviewSet.
func (m *EdiscoverySearch) SetAddToReviewSetOperation(value EdiscoveryAddToReviewSetOperationable)() {
    m.addToReviewSetOperation = value
}
// SetCustodianSources sets the custodianSources property value. Custodian sources that are included in the eDiscovery search.
func (m *EdiscoverySearch) SetCustodianSources(value []DataSourceable)() {
    m.custodianSources = value
}
// SetDataSourceScopes sets the dataSourceScopes property value. When specified, the collection will span across a service for an entire workload. Possible values are: none, allTenantMailboxes, allTenantSites, allCaseCustodians, allCaseNoncustodialDataSources.
func (m *EdiscoverySearch) SetDataSourceScopes(value *DataSourceScopes)() {
    m.dataSourceScopes = value
}
// SetLastEstimateStatisticsOperation sets the lastEstimateStatisticsOperation property value. The last estimate operation associated with the eDiscovery search.
func (m *EdiscoverySearch) SetLastEstimateStatisticsOperation(value EdiscoveryEstimateOperationable)() {
    m.lastEstimateStatisticsOperation = value
}
// SetNoncustodialSources sets the noncustodialSources property value. noncustodialDataSource sources that are included in the eDiscovery search
func (m *EdiscoverySearch) SetNoncustodialSources(value []EdiscoveryNoncustodialDataSourceable)() {
    m.noncustodialSources = value
}
