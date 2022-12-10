package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SiteCollection 
type SiteCollection struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The geographic region code for where this site collection resides. Read-only.
    dataLocationCode *string
    // The hostname for the site collection. Read-only.
    hostname *string
    // The OdataType property
    odataType *string
    // If present, indicates that this is a root site collection in SharePoint. Read-only.
    root Rootable
}
// NewSiteCollection instantiates a new siteCollection and sets the default values.
func NewSiteCollection()(*SiteCollection) {
    m := &SiteCollection{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSiteCollectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSiteCollectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSiteCollection(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SiteCollection) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDataLocationCode gets the dataLocationCode property value. The geographic region code for where this site collection resides. Read-only.
func (m *SiteCollection) GetDataLocationCode()(*string) {
    return m.dataLocationCode
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SiteCollection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["dataLocationCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDataLocationCode)
    res["hostname"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetHostname)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["root"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateRootFromDiscriminatorValue , m.SetRoot)
    return res
}
// GetHostname gets the hostname property value. The hostname for the site collection. Read-only.
func (m *SiteCollection) GetHostname()(*string) {
    return m.hostname
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SiteCollection) GetOdataType()(*string) {
    return m.odataType
}
// GetRoot gets the root property value. If present, indicates that this is a root site collection in SharePoint. Read-only.
func (m *SiteCollection) GetRoot()(Rootable) {
    return m.root
}
// Serialize serializes information the current object
func (m *SiteCollection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("dataLocationCode", m.GetDataLocationCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("hostname", m.GetHostname())
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
        err := writer.WriteObjectValue("root", m.GetRoot())
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
func (m *SiteCollection) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDataLocationCode sets the dataLocationCode property value. The geographic region code for where this site collection resides. Read-only.
func (m *SiteCollection) SetDataLocationCode(value *string)() {
    m.dataLocationCode = value
}
// SetHostname sets the hostname property value. The hostname for the site collection. Read-only.
func (m *SiteCollection) SetHostname(value *string)() {
    m.hostname = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SiteCollection) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRoot sets the root property value. If present, indicates that this is a root site collection in SharePoint. Read-only.
func (m *SiteCollection) SetRoot(value Rootable)() {
    m.root = value
}
