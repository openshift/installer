package drives

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody 
type ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody instantiates a new ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody and sets the default values.
func NewItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody()(*ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) {
    m := &ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) GetAdditionalData()(map[string]any) {
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
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCriteria1 gets the criteria1 property value. The criteria1 property
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) GetCriteria1()(*string) {
    val, err := m.GetBackingStore().Get("criteria1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCriteria2 gets the criteria2 property value. The criteria2 property
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) GetCriteria2()(*string) {
    val, err := m.GetBackingStore().Get("criteria2")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["criteria1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCriteria1(val)
        }
        return nil
    }
    res["criteria2"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCriteria2(val)
        }
        return nil
    }
    res["oper"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOper(val)
        }
        return nil
    }
    return res
}
// GetOper gets the oper property value. The oper property
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) GetOper()(*string) {
    val, err := m.GetBackingStore().Get("oper")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("criteria1", m.GetCriteria1())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("criteria2", m.GetCriteria2())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("oper", m.GetOper())
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
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the backingStore property value. Stores model information.
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCriteria1 sets the criteria1 property value. The criteria1 property
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) SetCriteria1(value *string)() {
    err := m.GetBackingStore().Set("criteria1", value)
    if err != nil {
        panic(err)
    }
}
// SetCriteria2 sets the criteria2 property value. The criteria2 property
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) SetCriteria2(value *string)() {
    err := m.GetBackingStore().Set("criteria2", value)
    if err != nil {
        panic(err)
    }
}
// SetOper sets the oper property value. The oper property
func (m *ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBody) SetOper(value *string)() {
    err := m.GetBackingStore().Set("oper", value)
    if err != nil {
        panic(err)
    }
}
// ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBodyable 
type ItemItemsItemWorkbookTablesItemColumnsItemFilterApplyCustomFilterPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCriteria1()(*string)
    GetCriteria2()(*string)
    GetOper()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCriteria1(value *string)()
    SetCriteria2(value *string)()
    SetOper(value *string)()
}
