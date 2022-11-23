package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosNetworkUsageRule network Usage Rules allow enterprises to specify how managed apps use networks, such as cellular data networks.
type IosNetworkUsageRule struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // If set to true, corresponding managed apps will not be allowed to use cellular data at any time.
    cellularDataBlocked *bool
    // If set to true, corresponding managed apps will not be allowed to use cellular data when roaming.
    cellularDataBlockWhenRoaming *bool
    // Information about the managed apps that this rule is going to apply to. This collection can contain a maximum of 500 elements.
    managedApps []AppListItemable
    // The OdataType property
    odataType *string
}
// NewIosNetworkUsageRule instantiates a new iosNetworkUsageRule and sets the default values.
func NewIosNetworkUsageRule()(*IosNetworkUsageRule) {
    m := &IosNetworkUsageRule{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateIosNetworkUsageRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosNetworkUsageRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosNetworkUsageRule(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *IosNetworkUsageRule) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCellularDataBlocked gets the cellularDataBlocked property value. If set to true, corresponding managed apps will not be allowed to use cellular data at any time.
func (m *IosNetworkUsageRule) GetCellularDataBlocked()(*bool) {
    return m.cellularDataBlocked
}
// GetCellularDataBlockWhenRoaming gets the cellularDataBlockWhenRoaming property value. If set to true, corresponding managed apps will not be allowed to use cellular data when roaming.
func (m *IosNetworkUsageRule) GetCellularDataBlockWhenRoaming()(*bool) {
    return m.cellularDataBlockWhenRoaming
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosNetworkUsageRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["cellularDataBlocked"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetCellularDataBlocked)
    res["cellularDataBlockWhenRoaming"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetCellularDataBlockWhenRoaming)
    res["managedApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAppListItemFromDiscriminatorValue , m.SetManagedApps)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetManagedApps gets the managedApps property value. Information about the managed apps that this rule is going to apply to. This collection can contain a maximum of 500 elements.
func (m *IosNetworkUsageRule) GetManagedApps()([]AppListItemable) {
    return m.managedApps
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *IosNetworkUsageRule) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *IosNetworkUsageRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("cellularDataBlocked", m.GetCellularDataBlocked())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("cellularDataBlockWhenRoaming", m.GetCellularDataBlockWhenRoaming())
        if err != nil {
            return err
        }
    }
    if m.GetManagedApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetManagedApps())
        err := writer.WriteCollectionOfObjectValues("managedApps", cast)
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
func (m *IosNetworkUsageRule) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCellularDataBlocked sets the cellularDataBlocked property value. If set to true, corresponding managed apps will not be allowed to use cellular data at any time.
func (m *IosNetworkUsageRule) SetCellularDataBlocked(value *bool)() {
    m.cellularDataBlocked = value
}
// SetCellularDataBlockWhenRoaming sets the cellularDataBlockWhenRoaming property value. If set to true, corresponding managed apps will not be allowed to use cellular data when roaming.
func (m *IosNetworkUsageRule) SetCellularDataBlockWhenRoaming(value *bool)() {
    m.cellularDataBlockWhenRoaming = value
}
// SetManagedApps sets the managedApps property value. Information about the managed apps that this rule is going to apply to. This collection can contain a maximum of 500 elements.
func (m *IosNetworkUsageRule) SetManagedApps(value []AppListItemable)() {
    m.managedApps = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *IosNetworkUsageRule) SetOdataType(value *string)() {
    m.odataType = value
}
