package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosDeviceType contains properties of the possible iOS device types the mobile app can run on.
type IosDeviceType struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Whether the app should run on iPads.
    iPad *bool
    // Whether the app should run on iPhones and iPods.
    iPhoneAndIPod *bool
    // The OdataType property
    odataType *string
}
// NewIosDeviceType instantiates a new iosDeviceType and sets the default values.
func NewIosDeviceType()(*IosDeviceType) {
    m := &IosDeviceType{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateIosDeviceTypeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosDeviceTypeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosDeviceType(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *IosDeviceType) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosDeviceType) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["iPad"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIPad)
    res["iPhoneAndIPod"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIPhoneAndIPod)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetIPad gets the iPad property value. Whether the app should run on iPads.
func (m *IosDeviceType) GetIPad()(*bool) {
    return m.iPad
}
// GetIPhoneAndIPod gets the iPhoneAndIPod property value. Whether the app should run on iPhones and iPods.
func (m *IosDeviceType) GetIPhoneAndIPod()(*bool) {
    return m.iPhoneAndIPod
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *IosDeviceType) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *IosDeviceType) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("iPad", m.GetIPad())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("iPhoneAndIPod", m.GetIPhoneAndIPod())
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
func (m *IosDeviceType) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIPad sets the iPad property value. Whether the app should run on iPads.
func (m *IosDeviceType) SetIPad(value *bool)() {
    m.iPad = value
}
// SetIPhoneAndIPod sets the iPhoneAndIPod property value. Whether the app should run on iPhones and iPods.
func (m *IosDeviceType) SetIPhoneAndIPod(value *bool)() {
    m.iPhoneAndIPod = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *IosDeviceType) SetOdataType(value *string)() {
    m.odataType = value
}
