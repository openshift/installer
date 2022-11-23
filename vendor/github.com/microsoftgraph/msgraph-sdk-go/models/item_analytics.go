package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ItemAnalytics 
type ItemAnalytics struct {
    Entity
    // The allTime property
    allTime ItemActivityStatable
    // The itemActivityStats property
    itemActivityStats []ItemActivityStatable
    // The lastSevenDays property
    lastSevenDays ItemActivityStatable
}
// NewItemAnalytics instantiates a new itemAnalytics and sets the default values.
func NewItemAnalytics()(*ItemAnalytics) {
    m := &ItemAnalytics{
        Entity: *NewEntity(),
    }
    return m
}
// CreateItemAnalyticsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateItemAnalyticsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemAnalytics(), nil
}
// GetAllTime gets the allTime property value. The allTime property
func (m *ItemAnalytics) GetAllTime()(ItemActivityStatable) {
    return m.allTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ItemAnalytics) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateItemActivityStatFromDiscriminatorValue , m.SetAllTime)
    res["itemActivityStats"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateItemActivityStatFromDiscriminatorValue , m.SetItemActivityStats)
    res["lastSevenDays"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateItemActivityStatFromDiscriminatorValue , m.SetLastSevenDays)
    return res
}
// GetItemActivityStats gets the itemActivityStats property value. The itemActivityStats property
func (m *ItemAnalytics) GetItemActivityStats()([]ItemActivityStatable) {
    return m.itemActivityStats
}
// GetLastSevenDays gets the lastSevenDays property value. The lastSevenDays property
func (m *ItemAnalytics) GetLastSevenDays()(ItemActivityStatable) {
    return m.lastSevenDays
}
// Serialize serializes information the current object
func (m *ItemAnalytics) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("allTime", m.GetAllTime())
        if err != nil {
            return err
        }
    }
    if m.GetItemActivityStats() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetItemActivityStats())
        err = writer.WriteCollectionOfObjectValues("itemActivityStats", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastSevenDays", m.GetLastSevenDays())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllTime sets the allTime property value. The allTime property
func (m *ItemAnalytics) SetAllTime(value ItemActivityStatable)() {
    m.allTime = value
}
// SetItemActivityStats sets the itemActivityStats property value. The itemActivityStats property
func (m *ItemAnalytics) SetItemActivityStats(value []ItemActivityStatable)() {
    m.itemActivityStats = value
}
// SetLastSevenDays sets the lastSevenDays property value. The lastSevenDays property
func (m *ItemAnalytics) SetLastSevenDays(value ItemActivityStatable)() {
    m.lastSevenDays = value
}
