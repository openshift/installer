package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SearchAlteration 
type SearchAlteration struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Defines the altered highlighted query string with spelling correction. The annotation around the corrected segment is: /ue000, /ue001.
    alteredHighlightedQueryString *string
    // Defines the altered query string with spelling correction.
    alteredQueryString *string
    // Represents changed segments related to an original user query.
    alteredQueryTokens []AlteredQueryTokenable
    // The OdataType property
    odataType *string
}
// NewSearchAlteration instantiates a new searchAlteration and sets the default values.
func NewSearchAlteration()(*SearchAlteration) {
    m := &SearchAlteration{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSearchAlterationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSearchAlterationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSearchAlteration(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SearchAlteration) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAlteredHighlightedQueryString gets the alteredHighlightedQueryString property value. Defines the altered highlighted query string with spelling correction. The annotation around the corrected segment is: /ue000, /ue001.
func (m *SearchAlteration) GetAlteredHighlightedQueryString()(*string) {
    return m.alteredHighlightedQueryString
}
// GetAlteredQueryString gets the alteredQueryString property value. Defines the altered query string with spelling correction.
func (m *SearchAlteration) GetAlteredQueryString()(*string) {
    return m.alteredQueryString
}
// GetAlteredQueryTokens gets the alteredQueryTokens property value. Represents changed segments related to an original user query.
func (m *SearchAlteration) GetAlteredQueryTokens()([]AlteredQueryTokenable) {
    return m.alteredQueryTokens
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SearchAlteration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["alteredHighlightedQueryString"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAlteredHighlightedQueryString)
    res["alteredQueryString"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAlteredQueryString)
    res["alteredQueryTokens"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAlteredQueryTokenFromDiscriminatorValue , m.SetAlteredQueryTokens)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SearchAlteration) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *SearchAlteration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("alteredHighlightedQueryString", m.GetAlteredHighlightedQueryString())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("alteredQueryString", m.GetAlteredQueryString())
        if err != nil {
            return err
        }
    }
    if m.GetAlteredQueryTokens() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAlteredQueryTokens())
        err := writer.WriteCollectionOfObjectValues("alteredQueryTokens", cast)
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
func (m *SearchAlteration) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAlteredHighlightedQueryString sets the alteredHighlightedQueryString property value. Defines the altered highlighted query string with spelling correction. The annotation around the corrected segment is: /ue000, /ue001.
func (m *SearchAlteration) SetAlteredHighlightedQueryString(value *string)() {
    m.alteredHighlightedQueryString = value
}
// SetAlteredQueryString sets the alteredQueryString property value. Defines the altered query string with spelling correction.
func (m *SearchAlteration) SetAlteredQueryString(value *string)() {
    m.alteredQueryString = value
}
// SetAlteredQueryTokens sets the alteredQueryTokens property value. Represents changed segments related to an original user query.
func (m *SearchAlteration) SetAlteredQueryTokens(value []AlteredQueryTokenable)() {
    m.alteredQueryTokens = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SearchAlteration) SetOdataType(value *string)() {
    m.odataType = value
}
