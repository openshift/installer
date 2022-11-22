package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosVppEBook 
type IosVppEBook struct {
    ManagedEBook
    // The Apple ID associated with Vpp token.
    appleId *string
    // Genres.
    genres []string
    // Language.
    language *string
    // Seller.
    seller *string
    // Total license count.
    totalLicenseCount *int32
    // Used license count.
    usedLicenseCount *int32
    // The Vpp token's organization name.
    vppOrganizationName *string
    // The Vpp token ID.
    vppTokenId *string
}
// NewIosVppEBook instantiates a new IosVppEBook and sets the default values.
func NewIosVppEBook()(*IosVppEBook) {
    m := &IosVppEBook{
        ManagedEBook: *NewManagedEBook(),
    }
    odataTypeValue := "#microsoft.graph.iosVppEBook";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateIosVppEBookFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosVppEBookFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosVppEBook(), nil
}
// GetAppleId gets the appleId property value. The Apple ID associated with Vpp token.
func (m *IosVppEBook) GetAppleId()(*string) {
    return m.appleId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosVppEBook) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedEBook.GetFieldDeserializers()
    res["appleId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppleId)
    res["genres"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetGenres)
    res["language"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLanguage)
    res["seller"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSeller)
    res["totalLicenseCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetTotalLicenseCount)
    res["usedLicenseCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetUsedLicenseCount)
    res["vppOrganizationName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVppOrganizationName)
    res["vppTokenId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVppTokenId)
    return res
}
// GetGenres gets the genres property value. Genres.
func (m *IosVppEBook) GetGenres()([]string) {
    return m.genres
}
// GetLanguage gets the language property value. Language.
func (m *IosVppEBook) GetLanguage()(*string) {
    return m.language
}
// GetSeller gets the seller property value. Seller.
func (m *IosVppEBook) GetSeller()(*string) {
    return m.seller
}
// GetTotalLicenseCount gets the totalLicenseCount property value. Total license count.
func (m *IosVppEBook) GetTotalLicenseCount()(*int32) {
    return m.totalLicenseCount
}
// GetUsedLicenseCount gets the usedLicenseCount property value. Used license count.
func (m *IosVppEBook) GetUsedLicenseCount()(*int32) {
    return m.usedLicenseCount
}
// GetVppOrganizationName gets the vppOrganizationName property value. The Vpp token's organization name.
func (m *IosVppEBook) GetVppOrganizationName()(*string) {
    return m.vppOrganizationName
}
// GetVppTokenId gets the vppTokenId property value. The Vpp token ID.
func (m *IosVppEBook) GetVppTokenId()(*string) {
    return m.vppTokenId
}
// Serialize serializes information the current object
func (m *IosVppEBook) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedEBook.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appleId", m.GetAppleId())
        if err != nil {
            return err
        }
    }
    if m.GetGenres() != nil {
        err = writer.WriteCollectionOfStringValues("genres", m.GetGenres())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("language", m.GetLanguage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("seller", m.GetSeller())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalLicenseCount", m.GetTotalLicenseCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("usedLicenseCount", m.GetUsedLicenseCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vppOrganizationName", m.GetVppOrganizationName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vppTokenId", m.GetVppTokenId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppleId sets the appleId property value. The Apple ID associated with Vpp token.
func (m *IosVppEBook) SetAppleId(value *string)() {
    m.appleId = value
}
// SetGenres sets the genres property value. Genres.
func (m *IosVppEBook) SetGenres(value []string)() {
    m.genres = value
}
// SetLanguage sets the language property value. Language.
func (m *IosVppEBook) SetLanguage(value *string)() {
    m.language = value
}
// SetSeller sets the seller property value. Seller.
func (m *IosVppEBook) SetSeller(value *string)() {
    m.seller = value
}
// SetTotalLicenseCount sets the totalLicenseCount property value. Total license count.
func (m *IosVppEBook) SetTotalLicenseCount(value *int32)() {
    m.totalLicenseCount = value
}
// SetUsedLicenseCount sets the usedLicenseCount property value. Used license count.
func (m *IosVppEBook) SetUsedLicenseCount(value *int32)() {
    m.usedLicenseCount = value
}
// SetVppOrganizationName sets the vppOrganizationName property value. The Vpp token's organization name.
func (m *IosVppEBook) SetVppOrganizationName(value *string)() {
    m.vppOrganizationName = value
}
// SetVppTokenId sets the vppTokenId property value. The Vpp token ID.
func (m *IosVppEBook) SetVppTokenId(value *string)() {
    m.vppTokenId = value
}
