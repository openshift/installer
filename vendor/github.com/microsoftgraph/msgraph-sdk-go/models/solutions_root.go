package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SolutionsRoot 
type SolutionsRoot struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The bookingBusinesses property
    bookingBusinesses []BookingBusinessable
    // The bookingCurrencies property
    bookingCurrencies []BookingCurrencyable
    // The OdataType property
    odataType *string
}
// NewSolutionsRoot instantiates a new SolutionsRoot and sets the default values.
func NewSolutionsRoot()(*SolutionsRoot) {
    m := &SolutionsRoot{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSolutionsRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSolutionsRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSolutionsRoot(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SolutionsRoot) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetBookingBusinesses gets the bookingBusinesses property value. The bookingBusinesses property
func (m *SolutionsRoot) GetBookingBusinesses()([]BookingBusinessable) {
    return m.bookingBusinesses
}
// GetBookingCurrencies gets the bookingCurrencies property value. The bookingCurrencies property
func (m *SolutionsRoot) GetBookingCurrencies()([]BookingCurrencyable) {
    return m.bookingCurrencies
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SolutionsRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["bookingBusinesses"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateBookingBusinessFromDiscriminatorValue , m.SetBookingBusinesses)
    res["bookingCurrencies"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateBookingCurrencyFromDiscriminatorValue , m.SetBookingCurrencies)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SolutionsRoot) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *SolutionsRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetBookingBusinesses() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetBookingBusinesses())
        err := writer.WriteCollectionOfObjectValues("bookingBusinesses", cast)
        if err != nil {
            return err
        }
    }
    if m.GetBookingCurrencies() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetBookingCurrencies())
        err := writer.WriteCollectionOfObjectValues("bookingCurrencies", cast)
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
func (m *SolutionsRoot) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetBookingBusinesses sets the bookingBusinesses property value. The bookingBusinesses property
func (m *SolutionsRoot) SetBookingBusinesses(value []BookingBusinessable)() {
    m.bookingBusinesses = value
}
// SetBookingCurrencies sets the bookingCurrencies property value. The bookingCurrencies property
func (m *SolutionsRoot) SetBookingCurrencies(value []BookingCurrencyable)() {
    m.bookingCurrencies = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SolutionsRoot) SetOdataType(value *string)() {
    m.odataType = value
}
