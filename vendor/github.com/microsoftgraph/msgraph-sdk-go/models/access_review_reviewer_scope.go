package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AccessReviewReviewerScope 
type AccessReviewReviewerScope struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The OdataType property
    odataType *string
    // The query specifying who will be the reviewer.
    query *string
    // In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query, for example, ./manager, is specified. Possible value: decisions.
    queryRoot *string
    // The type of query. Examples include MicrosoftGraph and ARM.
    queryType *string
}
// NewAccessReviewReviewerScope instantiates a new accessReviewReviewerScope and sets the default values.
func NewAccessReviewReviewerScope()(*AccessReviewReviewerScope) {
    m := &AccessReviewReviewerScope{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateAccessReviewReviewerScopeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessReviewReviewerScopeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewReviewerScope(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AccessReviewReviewerScope) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessReviewReviewerScope) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["query"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQuery)
    res["queryRoot"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQueryRoot)
    res["queryType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQueryType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *AccessReviewReviewerScope) GetOdataType()(*string) {
    return m.odataType
}
// GetQuery gets the query property value. The query specifying who will be the reviewer.
func (m *AccessReviewReviewerScope) GetQuery()(*string) {
    return m.query
}
// GetQueryRoot gets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query, for example, ./manager, is specified. Possible value: decisions.
func (m *AccessReviewReviewerScope) GetQueryRoot()(*string) {
    return m.queryRoot
}
// GetQueryType gets the queryType property value. The type of query. Examples include MicrosoftGraph and ARM.
func (m *AccessReviewReviewerScope) GetQueryType()(*string) {
    return m.queryType
}
// Serialize serializes information the current object
func (m *AccessReviewReviewerScope) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("query", m.GetQuery())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("queryRoot", m.GetQueryRoot())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("queryType", m.GetQueryType())
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
func (m *AccessReviewReviewerScope) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessReviewReviewerScope) SetOdataType(value *string)() {
    m.odataType = value
}
// SetQuery sets the query property value. The query specifying who will be the reviewer.
func (m *AccessReviewReviewerScope) SetQuery(value *string)() {
    m.query = value
}
// SetQueryRoot sets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query, for example, ./manager, is specified. Possible value: decisions.
func (m *AccessReviewReviewerScope) SetQueryRoot(value *string)() {
    m.queryRoot = value
}
// SetQueryType sets the queryType property value. The type of query. Examples include MicrosoftGraph and ARM.
func (m *AccessReviewReviewerScope) SetQueryType(value *string)() {
    m.queryType = value
}
