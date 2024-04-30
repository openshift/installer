package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// RequestSignatureVerification 
type RequestSignatureVerification struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewRequestSignatureVerification instantiates a new requestSignatureVerification and sets the default values.
func NewRequestSignatureVerification()(*RequestSignatureVerification) {
    m := &RequestSignatureVerification{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateRequestSignatureVerificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRequestSignatureVerificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRequestSignatureVerification(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RequestSignatureVerification) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAllowedWeakAlgorithms gets the allowedWeakAlgorithms property value. The allowedWeakAlgorithms property
func (m *RequestSignatureVerification) GetAllowedWeakAlgorithms()(*WeakAlgorithms) {
    val, err := m.GetBackingStore().Get("allowedWeakAlgorithms")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WeakAlgorithms)
    }
    return nil
}
// GetBackingStore gets the backingStore property value. Stores model information.
func (m *RequestSignatureVerification) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RequestSignatureVerification) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowedWeakAlgorithms"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWeakAlgorithms)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedWeakAlgorithms(val.(*WeakAlgorithms))
        }
        return nil
    }
    res["isSignedRequestRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSignedRequestRequired(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    return res
}
// GetIsSignedRequestRequired gets the isSignedRequestRequired property value. The isSignedRequestRequired property
func (m *RequestSignatureVerification) GetIsSignedRequestRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isSignedRequestRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *RequestSignatureVerification) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RequestSignatureVerification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAllowedWeakAlgorithms() != nil {
        cast := (*m.GetAllowedWeakAlgorithms()).String()
        err := writer.WriteStringValue("allowedWeakAlgorithms", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isSignedRequestRequired", m.GetIsSignedRequestRequired())
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
func (m *RequestSignatureVerification) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedWeakAlgorithms sets the allowedWeakAlgorithms property value. The allowedWeakAlgorithms property
func (m *RequestSignatureVerification) SetAllowedWeakAlgorithms(value *WeakAlgorithms)() {
    err := m.GetBackingStore().Set("allowedWeakAlgorithms", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the backingStore property value. Stores model information.
func (m *RequestSignatureVerification) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIsSignedRequestRequired sets the isSignedRequestRequired property value. The isSignedRequestRequired property
func (m *RequestSignatureVerification) SetIsSignedRequestRequired(value *bool)() {
    err := m.GetBackingStore().Set("isSignedRequestRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RequestSignatureVerification) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// RequestSignatureVerificationable 
type RequestSignatureVerificationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedWeakAlgorithms()(*WeakAlgorithms)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIsSignedRequestRequired()(*bool)
    GetOdataType()(*string)
    SetAllowedWeakAlgorithms(value *WeakAlgorithms)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIsSignedRequestRequired(value *bool)()
    SetOdataType(value *string)()
}
