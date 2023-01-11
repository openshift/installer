package callrecords

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Media 
type Media struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Device information associated with the callee endpoint of this media.
    calleeDevice DeviceInfoable
    // Network information associated with the callee endpoint of this media.
    calleeNetwork NetworkInfoable
    // Device information associated with the caller endpoint of this media.
    callerDevice DeviceInfoable
    // Network information associated with the caller endpoint of this media.
    callerNetwork NetworkInfoable
    // How the media was identified during media negotiation stage.
    label *string
    // The OdataType property
    odataType *string
    // Network streams associated with this media.
    streams []MediaStreamable
}
// NewMedia instantiates a new media and sets the default values.
func NewMedia()(*Media) {
    m := &Media{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateMediaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateMediaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMedia(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Media) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCalleeDevice gets the calleeDevice property value. Device information associated with the callee endpoint of this media.
func (m *Media) GetCalleeDevice()(DeviceInfoable) {
    return m.calleeDevice
}
// GetCalleeNetwork gets the calleeNetwork property value. Network information associated with the callee endpoint of this media.
func (m *Media) GetCalleeNetwork()(NetworkInfoable) {
    return m.calleeNetwork
}
// GetCallerDevice gets the callerDevice property value. Device information associated with the caller endpoint of this media.
func (m *Media) GetCallerDevice()(DeviceInfoable) {
    return m.callerDevice
}
// GetCallerNetwork gets the callerNetwork property value. Network information associated with the caller endpoint of this media.
func (m *Media) GetCallerNetwork()(NetworkInfoable) {
    return m.callerNetwork
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Media) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["calleeDevice"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceInfoFromDiscriminatorValue , m.SetCalleeDevice)
    res["calleeNetwork"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateNetworkInfoFromDiscriminatorValue , m.SetCalleeNetwork)
    res["callerDevice"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceInfoFromDiscriminatorValue , m.SetCallerDevice)
    res["callerNetwork"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateNetworkInfoFromDiscriminatorValue , m.SetCallerNetwork)
    res["label"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLabel)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["streams"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMediaStreamFromDiscriminatorValue , m.SetStreams)
    return res
}
// GetLabel gets the label property value. How the media was identified during media negotiation stage.
func (m *Media) GetLabel()(*string) {
    return m.label
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Media) GetOdataType()(*string) {
    return m.odataType
}
// GetStreams gets the streams property value. Network streams associated with this media.
func (m *Media) GetStreams()([]MediaStreamable) {
    return m.streams
}
// Serialize serializes information the current object
func (m *Media) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("calleeDevice", m.GetCalleeDevice())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("calleeNetwork", m.GetCalleeNetwork())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("callerDevice", m.GetCallerDevice())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("callerNetwork", m.GetCallerNetwork())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("label", m.GetLabel())
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
    if m.GetStreams() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetStreams())
        err := writer.WriteCollectionOfObjectValues("streams", cast)
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
func (m *Media) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCalleeDevice sets the calleeDevice property value. Device information associated with the callee endpoint of this media.
func (m *Media) SetCalleeDevice(value DeviceInfoable)() {
    m.calleeDevice = value
}
// SetCalleeNetwork sets the calleeNetwork property value. Network information associated with the callee endpoint of this media.
func (m *Media) SetCalleeNetwork(value NetworkInfoable)() {
    m.calleeNetwork = value
}
// SetCallerDevice sets the callerDevice property value. Device information associated with the caller endpoint of this media.
func (m *Media) SetCallerDevice(value DeviceInfoable)() {
    m.callerDevice = value
}
// SetCallerNetwork sets the callerNetwork property value. Network information associated with the caller endpoint of this media.
func (m *Media) SetCallerNetwork(value NetworkInfoable)() {
    m.callerNetwork = value
}
// SetLabel sets the label property value. How the media was identified during media negotiation stage.
func (m *Media) SetLabel(value *string)() {
    m.label = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Media) SetOdataType(value *string)() {
    m.odataType = value
}
// SetStreams sets the streams property value. Network streams associated with this media.
func (m *Media) SetStreams(value []MediaStreamable)() {
    m.streams = value
}
