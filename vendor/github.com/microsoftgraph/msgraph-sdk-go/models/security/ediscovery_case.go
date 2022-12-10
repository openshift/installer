package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// EdiscoveryCase 
type EdiscoveryCase struct {
    Case_escaped
    // The user who closed the case.
    closedBy iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable
    // The date and time when the case was closed. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
    closedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Returns a list of case ediscoveryCustodian objects for this case.
    custodians []EdiscoveryCustodianable
    // The external case number for customer reference.
    externalId *string
    // Returns a list of case ediscoveryNoncustodialDataSource objects for this case.
    noncustodialDataSources []EdiscoveryNoncustodialDataSourceable
    // Returns a list of case caseOperation objects for this case.
    operations []CaseOperationable
    // Returns a list of eDiscoveryReviewSet objects in the case.
    reviewSets []EdiscoveryReviewSetable
    // Returns a list of eDiscoverySearch objects associated with this case.
    searches []EdiscoverySearchable
    // Returns a list of eDIscoverySettings objects in the case.
    settings EdiscoveryCaseSettingsable
    // Returns a list of ediscoveryReviewTag objects associated to this case.
    tags []EdiscoveryReviewTagable
}
// NewEdiscoveryCase instantiates a new EdiscoveryCase and sets the default values.
func NewEdiscoveryCase()(*EdiscoveryCase) {
    m := &EdiscoveryCase{
        Case_escaped: *NewCase_escaped(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoveryCase";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEdiscoveryCaseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEdiscoveryCaseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryCase(), nil
}
// GetClosedBy gets the closedBy property value. The user who closed the case.
func (m *EdiscoveryCase) GetClosedBy()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable) {
    return m.closedBy
}
// GetClosedDateTime gets the closedDateTime property value. The date and time when the case was closed. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EdiscoveryCase) GetClosedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.closedDateTime
}
// GetCustodians gets the custodians property value. Returns a list of case ediscoveryCustodian objects for this case.
func (m *EdiscoveryCase) GetCustodians()([]EdiscoveryCustodianable) {
    return m.custodians
}
// GetExternalId gets the externalId property value. The external case number for customer reference.
func (m *EdiscoveryCase) GetExternalId()(*string) {
    return m.externalId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EdiscoveryCase) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Case_escaped.GetFieldDeserializers()
    res["closedBy"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentitySetFromDiscriminatorValue , m.SetClosedBy)
    res["closedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetClosedDateTime)
    res["custodians"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEdiscoveryCustodianFromDiscriminatorValue , m.SetCustodians)
    res["externalId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetExternalId)
    res["noncustodialDataSources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEdiscoveryNoncustodialDataSourceFromDiscriminatorValue , m.SetNoncustodialDataSources)
    res["operations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCaseOperationFromDiscriminatorValue , m.SetOperations)
    res["reviewSets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEdiscoveryReviewSetFromDiscriminatorValue , m.SetReviewSets)
    res["searches"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEdiscoverySearchFromDiscriminatorValue , m.SetSearches)
    res["settings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateEdiscoveryCaseSettingsFromDiscriminatorValue , m.SetSettings)
    res["tags"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateEdiscoveryReviewTagFromDiscriminatorValue , m.SetTags)
    return res
}
// GetNoncustodialDataSources gets the noncustodialDataSources property value. Returns a list of case ediscoveryNoncustodialDataSource objects for this case.
func (m *EdiscoveryCase) GetNoncustodialDataSources()([]EdiscoveryNoncustodialDataSourceable) {
    return m.noncustodialDataSources
}
// GetOperations gets the operations property value. Returns a list of case caseOperation objects for this case.
func (m *EdiscoveryCase) GetOperations()([]CaseOperationable) {
    return m.operations
}
// GetReviewSets gets the reviewSets property value. Returns a list of eDiscoveryReviewSet objects in the case.
func (m *EdiscoveryCase) GetReviewSets()([]EdiscoveryReviewSetable) {
    return m.reviewSets
}
// GetSearches gets the searches property value. Returns a list of eDiscoverySearch objects associated with this case.
func (m *EdiscoveryCase) GetSearches()([]EdiscoverySearchable) {
    return m.searches
}
// GetSettings gets the settings property value. Returns a list of eDIscoverySettings objects in the case.
func (m *EdiscoveryCase) GetSettings()(EdiscoveryCaseSettingsable) {
    return m.settings
}
// GetTags gets the tags property value. Returns a list of ediscoveryReviewTag objects associated to this case.
func (m *EdiscoveryCase) GetTags()([]EdiscoveryReviewTagable) {
    return m.tags
}
// Serialize serializes information the current object
func (m *EdiscoveryCase) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Case_escaped.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("closedBy", m.GetClosedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("closedDateTime", m.GetClosedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetCustodians() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCustodians())
        err = writer.WriteCollectionOfObjectValues("custodians", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    if m.GetNoncustodialDataSources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNoncustodialDataSources())
        err = writer.WriteCollectionOfObjectValues("noncustodialDataSources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOperations())
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetReviewSets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetReviewSets())
        err = writer.WriteCollectionOfObjectValues("reviewSets", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSearches() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSearches())
        err = writer.WriteCollectionOfObjectValues("searches", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTags())
        err = writer.WriteCollectionOfObjectValues("tags", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClosedBy sets the closedBy property value. The user who closed the case.
func (m *EdiscoveryCase) SetClosedBy(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)() {
    m.closedBy = value
}
// SetClosedDateTime sets the closedDateTime property value. The date and time when the case was closed. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EdiscoveryCase) SetClosedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.closedDateTime = value
}
// SetCustodians sets the custodians property value. Returns a list of case ediscoveryCustodian objects for this case.
func (m *EdiscoveryCase) SetCustodians(value []EdiscoveryCustodianable)() {
    m.custodians = value
}
// SetExternalId sets the externalId property value. The external case number for customer reference.
func (m *EdiscoveryCase) SetExternalId(value *string)() {
    m.externalId = value
}
// SetNoncustodialDataSources sets the noncustodialDataSources property value. Returns a list of case ediscoveryNoncustodialDataSource objects for this case.
func (m *EdiscoveryCase) SetNoncustodialDataSources(value []EdiscoveryNoncustodialDataSourceable)() {
    m.noncustodialDataSources = value
}
// SetOperations sets the operations property value. Returns a list of case caseOperation objects for this case.
func (m *EdiscoveryCase) SetOperations(value []CaseOperationable)() {
    m.operations = value
}
// SetReviewSets sets the reviewSets property value. Returns a list of eDiscoveryReviewSet objects in the case.
func (m *EdiscoveryCase) SetReviewSets(value []EdiscoveryReviewSetable)() {
    m.reviewSets = value
}
// SetSearches sets the searches property value. Returns a list of eDiscoverySearch objects associated with this case.
func (m *EdiscoveryCase) SetSearches(value []EdiscoverySearchable)() {
    m.searches = value
}
// SetSettings sets the settings property value. Returns a list of eDIscoverySettings objects in the case.
func (m *EdiscoveryCase) SetSettings(value EdiscoveryCaseSettingsable)() {
    m.settings = value
}
// SetTags sets the tags property value. Returns a list of ediscoveryReviewTag objects associated to this case.
func (m *EdiscoveryCase) SetTags(value []EdiscoveryReviewTagable)() {
    m.tags = value
}
