package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AccessReviewQueryScope 
type AccessReviewQueryScope struct {
    AccessReviewScope
    // The query representing what will be reviewed in an access review.
    query *string
    // In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query is specified. For example, ./manager.
    queryRoot *string
    // Indicates the type of query. Types include MicrosoftGraph and ARM.
    queryType *string
}
// NewAccessReviewQueryScope instantiates a new AccessReviewQueryScope and sets the default values.
func NewAccessReviewQueryScope()(*AccessReviewQueryScope) {
    m := &AccessReviewQueryScope{
        AccessReviewScope: *NewAccessReviewScope(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewQueryScope";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateAccessReviewQueryScopeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessReviewQueryScopeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.accessReviewInactiveUsersQueryScope":
                        return NewAccessReviewInactiveUsersQueryScope(), nil
                }
            }
        }
    }
    return NewAccessReviewQueryScope(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessReviewQueryScope) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewScope.GetFieldDeserializers()
    res["query"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQuery)
    res["queryRoot"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQueryRoot)
    res["queryType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQueryType)
    return res
}
// GetQuery gets the query property value. The query representing what will be reviewed in an access review.
func (m *AccessReviewQueryScope) GetQuery()(*string) {
    return m.query
}
// GetQueryRoot gets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query is specified. For example, ./manager.
func (m *AccessReviewQueryScope) GetQueryRoot()(*string) {
    return m.queryRoot
}
// GetQueryType gets the queryType property value. Indicates the type of query. Types include MicrosoftGraph and ARM.
func (m *AccessReviewQueryScope) GetQueryType()(*string) {
    return m.queryType
}
// Serialize serializes information the current object
func (m *AccessReviewQueryScope) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewScope.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("query", m.GetQuery())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("queryRoot", m.GetQueryRoot())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("queryType", m.GetQueryType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetQuery sets the query property value. The query representing what will be reviewed in an access review.
func (m *AccessReviewQueryScope) SetQuery(value *string)() {
    m.query = value
}
// SetQueryRoot sets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query is specified. For example, ./manager.
func (m *AccessReviewQueryScope) SetQueryRoot(value *string)() {
    m.queryRoot = value
}
// SetQueryType sets the queryType property value. Indicates the type of query. Types include MicrosoftGraph and ARM.
func (m *AccessReviewQueryScope) SetQueryType(value *string)() {
    m.queryType = value
}
