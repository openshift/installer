package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EdgeSearchEngineCustom 
type EdgeSearchEngineCustom struct {
    EdgeSearchEngineBase
    // Points to a https link containing the OpenSearch xml file that contains, at minimum, the short name and the URL to the search Engine.
    edgeSearchEngineOpenSearchXmlUrl *string
}
// NewEdgeSearchEngineCustom instantiates a new EdgeSearchEngineCustom and sets the default values.
func NewEdgeSearchEngineCustom()(*EdgeSearchEngineCustom) {
    m := &EdgeSearchEngineCustom{
        EdgeSearchEngineBase: *NewEdgeSearchEngineBase(),
    }
    odataTypeValue := "#microsoft.graph.edgeSearchEngineCustom";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateEdgeSearchEngineCustomFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEdgeSearchEngineCustomFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdgeSearchEngineCustom(), nil
}
// GetEdgeSearchEngineOpenSearchXmlUrl gets the edgeSearchEngineOpenSearchXmlUrl property value. Points to a https link containing the OpenSearch xml file that contains, at minimum, the short name and the URL to the search Engine.
func (m *EdgeSearchEngineCustom) GetEdgeSearchEngineOpenSearchXmlUrl()(*string) {
    return m.edgeSearchEngineOpenSearchXmlUrl
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EdgeSearchEngineCustom) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EdgeSearchEngineBase.GetFieldDeserializers()
    res["edgeSearchEngineOpenSearchXmlUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEdgeSearchEngineOpenSearchXmlUrl)
    return res
}
// Serialize serializes information the current object
func (m *EdgeSearchEngineCustom) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EdgeSearchEngineBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("edgeSearchEngineOpenSearchXmlUrl", m.GetEdgeSearchEngineOpenSearchXmlUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEdgeSearchEngineOpenSearchXmlUrl sets the edgeSearchEngineOpenSearchXmlUrl property value. Points to a https link containing the OpenSearch xml file that contains, at minimum, the short name and the URL to the search Engine.
func (m *EdgeSearchEngineCustom) SetEdgeSearchEngineOpenSearchXmlUrl(value *string)() {
    m.edgeSearchEngineOpenSearchXmlUrl = value
}
