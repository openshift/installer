package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10NetworkProxyServer network Proxy Server Policy.
type Windows10NetworkProxyServer struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Address to the proxy server. Specify an address in the format [':']
    address *string
    // Addresses that should not use the proxy server. The system will not use the proxy server for addresses beginning with what is specified in this node.
    exceptions []string
    // The OdataType property
    odataType *string
    // Specifies whether the proxy server should be used for local (intranet) addresses.
    useForLocalAddresses *bool
}
// NewWindows10NetworkProxyServer instantiates a new windows10NetworkProxyServer and sets the default values.
func NewWindows10NetworkProxyServer()(*Windows10NetworkProxyServer) {
    m := &Windows10NetworkProxyServer{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWindows10NetworkProxyServerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindows10NetworkProxyServerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10NetworkProxyServer(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Windows10NetworkProxyServer) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAddress gets the address property value. Address to the proxy server. Specify an address in the format [':']
func (m *Windows10NetworkProxyServer) GetAddress()(*string) {
    return m.address
}
// GetExceptions gets the exceptions property value. Addresses that should not use the proxy server. The system will not use the proxy server for addresses beginning with what is specified in this node.
func (m *Windows10NetworkProxyServer) GetExceptions()([]string) {
    return m.exceptions
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Windows10NetworkProxyServer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["address"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAddress)
    res["exceptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetExceptions)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["useForLocalAddresses"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetUseForLocalAddresses)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Windows10NetworkProxyServer) GetOdataType()(*string) {
    return m.odataType
}
// GetUseForLocalAddresses gets the useForLocalAddresses property value. Specifies whether the proxy server should be used for local (intranet) addresses.
func (m *Windows10NetworkProxyServer) GetUseForLocalAddresses()(*bool) {
    return m.useForLocalAddresses
}
// Serialize serializes information the current object
func (m *Windows10NetworkProxyServer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("address", m.GetAddress())
        if err != nil {
            return err
        }
    }
    if m.GetExceptions() != nil {
        err := writer.WriteCollectionOfStringValues("exceptions", m.GetExceptions())
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
        err := writer.WriteBoolValue("useForLocalAddresses", m.GetUseForLocalAddresses())
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
func (m *Windows10NetworkProxyServer) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAddress sets the address property value. Address to the proxy server. Specify an address in the format [':']
func (m *Windows10NetworkProxyServer) SetAddress(value *string)() {
    m.address = value
}
// SetExceptions sets the exceptions property value. Addresses that should not use the proxy server. The system will not use the proxy server for addresses beginning with what is specified in this node.
func (m *Windows10NetworkProxyServer) SetExceptions(value []string)() {
    m.exceptions = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Windows10NetworkProxyServer) SetOdataType(value *string)() {
    m.odataType = value
}
// SetUseForLocalAddresses sets the useForLocalAddresses property value. Specifies whether the proxy server should be used for local (intranet) addresses.
func (m *Windows10NetworkProxyServer) SetUseForLocalAddresses(value *bool)() {
    m.useForLocalAddresses = value
}
