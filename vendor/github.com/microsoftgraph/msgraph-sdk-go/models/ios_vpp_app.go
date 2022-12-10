package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosVppApp 
type IosVppApp struct {
    MobileApp
    // The applicable iOS Device Type.
    applicableDeviceType IosDeviceTypeable
    // The store URL.
    appStoreUrl *string
    // The Identity Name.
    bundleId *string
    // The supported License Type.
    licensingType VppLicensingTypeable
    // The VPP application release date and time.
    releaseDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The total number of VPP licenses.
    totalLicenseCount *int32
    // The number of VPP licenses in use.
    usedLicenseCount *int32
    // Possible types of an Apple Volume Purchase Program token.
    vppTokenAccountType *VppTokenAccountType
    // The Apple Id associated with the given Apple Volume Purchase Program Token.
    vppTokenAppleId *string
    // The organization associated with the Apple Volume Purchase Program Token
    vppTokenOrganizationName *string
}
// NewIosVppApp instantiates a new IosVppApp and sets the default values.
func NewIosVppApp()(*IosVppApp) {
    m := &IosVppApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.iosVppApp";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateIosVppAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosVppAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosVppApp(), nil
}
// GetApplicableDeviceType gets the applicableDeviceType property value. The applicable iOS Device Type.
func (m *IosVppApp) GetApplicableDeviceType()(IosDeviceTypeable) {
    return m.applicableDeviceType
}
// GetAppStoreUrl gets the appStoreUrl property value. The store URL.
func (m *IosVppApp) GetAppStoreUrl()(*string) {
    return m.appStoreUrl
}
// GetBundleId gets the bundleId property value. The Identity Name.
func (m *IosVppApp) GetBundleId()(*string) {
    return m.bundleId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosVppApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    res["applicableDeviceType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIosDeviceTypeFromDiscriminatorValue , m.SetApplicableDeviceType)
    res["appStoreUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppStoreUrl)
    res["bundleId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetBundleId)
    res["licensingType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateVppLicensingTypeFromDiscriminatorValue , m.SetLicensingType)
    res["releaseDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetReleaseDateTime)
    res["totalLicenseCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetTotalLicenseCount)
    res["usedLicenseCount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetUsedLicenseCount)
    res["vppTokenAccountType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseVppTokenAccountType , m.SetVppTokenAccountType)
    res["vppTokenAppleId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVppTokenAppleId)
    res["vppTokenOrganizationName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVppTokenOrganizationName)
    return res
}
// GetLicensingType gets the licensingType property value. The supported License Type.
func (m *IosVppApp) GetLicensingType()(VppLicensingTypeable) {
    return m.licensingType
}
// GetReleaseDateTime gets the releaseDateTime property value. The VPP application release date and time.
func (m *IosVppApp) GetReleaseDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.releaseDateTime
}
// GetTotalLicenseCount gets the totalLicenseCount property value. The total number of VPP licenses.
func (m *IosVppApp) GetTotalLicenseCount()(*int32) {
    return m.totalLicenseCount
}
// GetUsedLicenseCount gets the usedLicenseCount property value. The number of VPP licenses in use.
func (m *IosVppApp) GetUsedLicenseCount()(*int32) {
    return m.usedLicenseCount
}
// GetVppTokenAccountType gets the vppTokenAccountType property value. Possible types of an Apple Volume Purchase Program token.
func (m *IosVppApp) GetVppTokenAccountType()(*VppTokenAccountType) {
    return m.vppTokenAccountType
}
// GetVppTokenAppleId gets the vppTokenAppleId property value. The Apple Id associated with the given Apple Volume Purchase Program Token.
func (m *IosVppApp) GetVppTokenAppleId()(*string) {
    return m.vppTokenAppleId
}
// GetVppTokenOrganizationName gets the vppTokenOrganizationName property value. The organization associated with the Apple Volume Purchase Program Token
func (m *IosVppApp) GetVppTokenOrganizationName()(*string) {
    return m.vppTokenOrganizationName
}
// Serialize serializes information the current object
func (m *IosVppApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("applicableDeviceType", m.GetApplicableDeviceType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appStoreUrl", m.GetAppStoreUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("bundleId", m.GetBundleId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("licensingType", m.GetLicensingType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("releaseDateTime", m.GetReleaseDateTime())
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
    if m.GetVppTokenAccountType() != nil {
        cast := (*m.GetVppTokenAccountType()).String()
        err = writer.WriteStringValue("vppTokenAccountType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vppTokenAppleId", m.GetVppTokenAppleId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vppTokenOrganizationName", m.GetVppTokenOrganizationName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicableDeviceType sets the applicableDeviceType property value. The applicable iOS Device Type.
func (m *IosVppApp) SetApplicableDeviceType(value IosDeviceTypeable)() {
    m.applicableDeviceType = value
}
// SetAppStoreUrl sets the appStoreUrl property value. The store URL.
func (m *IosVppApp) SetAppStoreUrl(value *string)() {
    m.appStoreUrl = value
}
// SetBundleId sets the bundleId property value. The Identity Name.
func (m *IosVppApp) SetBundleId(value *string)() {
    m.bundleId = value
}
// SetLicensingType sets the licensingType property value. The supported License Type.
func (m *IosVppApp) SetLicensingType(value VppLicensingTypeable)() {
    m.licensingType = value
}
// SetReleaseDateTime sets the releaseDateTime property value. The VPP application release date and time.
func (m *IosVppApp) SetReleaseDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.releaseDateTime = value
}
// SetTotalLicenseCount sets the totalLicenseCount property value. The total number of VPP licenses.
func (m *IosVppApp) SetTotalLicenseCount(value *int32)() {
    m.totalLicenseCount = value
}
// SetUsedLicenseCount sets the usedLicenseCount property value. The number of VPP licenses in use.
func (m *IosVppApp) SetUsedLicenseCount(value *int32)() {
    m.usedLicenseCount = value
}
// SetVppTokenAccountType sets the vppTokenAccountType property value. Possible types of an Apple Volume Purchase Program token.
func (m *IosVppApp) SetVppTokenAccountType(value *VppTokenAccountType)() {
    m.vppTokenAccountType = value
}
// SetVppTokenAppleId sets the vppTokenAppleId property value. The Apple Id associated with the given Apple Volume Purchase Program Token.
func (m *IosVppApp) SetVppTokenAppleId(value *string)() {
    m.vppTokenAppleId = value
}
// SetVppTokenOrganizationName sets the vppTokenOrganizationName property value. The organization associated with the Apple Volume Purchase Program Token
func (m *IosVppApp) SetVppTokenOrganizationName(value *string)() {
    m.vppTokenOrganizationName = value
}
