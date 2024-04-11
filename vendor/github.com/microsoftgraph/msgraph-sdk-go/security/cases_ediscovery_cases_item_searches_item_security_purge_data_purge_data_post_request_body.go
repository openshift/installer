package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody 
type CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody instantiates a new CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody and sets the default values.
func NewCasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody()(*CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) {
    m := &CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the backingStore property value. Stores model information.
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["purgeAreas"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ParsePurgeAreas)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPurgeAreas(val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas))
        }
        return nil
    }
    res["purgeType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ParsePurgeType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPurgeType(val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType))
        }
        return nil
    }
    return res
}
// GetPurgeAreas gets the purgeAreas property value. The purgeAreas property
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) GetPurgeAreas()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas) {
    val, err := m.GetBackingStore().Get("purgeAreas")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas)
    }
    return nil
}
// GetPurgeType gets the purgeType property value. The purgeType property
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) GetPurgeType()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType) {
    val, err := m.GetBackingStore().Get("purgeType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetPurgeAreas() != nil {
        cast := (*m.GetPurgeAreas()).String()
        err := writer.WriteStringValue("purgeAreas", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetPurgeType() != nil {
        cast := (*m.GetPurgeType()).String()
        err := writer.WriteStringValue("purgeType", &cast)
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
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the backingStore property value. Stores model information.
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetPurgeAreas sets the purgeAreas property value. The purgeAreas property
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) SetPurgeAreas(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas)() {
    err := m.GetBackingStore().Set("purgeAreas", value)
    if err != nil {
        panic(err)
    }
}
// SetPurgeType sets the purgeType property value. The purgeType property
func (m *CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBody) SetPurgeType(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType)() {
    err := m.GetBackingStore().Set("purgeType", value)
    if err != nil {
        panic(err)
    }
}
// CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBodyable 
type CasesEdiscoveryCasesItemSearchesItemSecurityPurgeDataPurgeDataPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetPurgeAreas()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas)
    GetPurgeType()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetPurgeAreas(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeAreas)()
    SetPurgeType(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.PurgeType)()
}
