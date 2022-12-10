package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserTrainingContentEventInfo 
type UserTrainingContentEventInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Browser of the user from where the training event was generated.
    browser *string
    // Date and time of the training content playback by the user.
    contentDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // IP address of the user for the training event.
    ipAddress *string
    // The OdataType property
    odataType *string
    // The operating system, platform, and device details of the user for the training event.
    osPlatformDeviceDetails *string
    // Potential improvement in the tenant security posture after completion of the training by the user.
    potentialScoreImpact *float64
}
// NewUserTrainingContentEventInfo instantiates a new userTrainingContentEventInfo and sets the default values.
func NewUserTrainingContentEventInfo()(*UserTrainingContentEventInfo) {
    m := &UserTrainingContentEventInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateUserTrainingContentEventInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserTrainingContentEventInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserTrainingContentEventInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *UserTrainingContentEventInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetBrowser gets the browser property value. Browser of the user from where the training event was generated.
func (m *UserTrainingContentEventInfo) GetBrowser()(*string) {
    return m.browser
}
// GetContentDateTime gets the contentDateTime property value. Date and time of the training content playback by the user.
func (m *UserTrainingContentEventInfo) GetContentDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.contentDateTime
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserTrainingContentEventInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["browser"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetBrowser)
    res["contentDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetContentDateTime)
    res["ipAddress"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetIpAddress)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["osPlatformDeviceDetails"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOsPlatformDeviceDetails)
    res["potentialScoreImpact"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetFloat64Value(m.SetPotentialScoreImpact)
    return res
}
// GetIpAddress gets the ipAddress property value. IP address of the user for the training event.
func (m *UserTrainingContentEventInfo) GetIpAddress()(*string) {
    return m.ipAddress
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *UserTrainingContentEventInfo) GetOdataType()(*string) {
    return m.odataType
}
// GetOsPlatformDeviceDetails gets the osPlatformDeviceDetails property value. The operating system, platform, and device details of the user for the training event.
func (m *UserTrainingContentEventInfo) GetOsPlatformDeviceDetails()(*string) {
    return m.osPlatformDeviceDetails
}
// GetPotentialScoreImpact gets the potentialScoreImpact property value. Potential improvement in the tenant security posture after completion of the training by the user.
func (m *UserTrainingContentEventInfo) GetPotentialScoreImpact()(*float64) {
    return m.potentialScoreImpact
}
// Serialize serializes information the current object
func (m *UserTrainingContentEventInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("browser", m.GetBrowser())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("contentDateTime", m.GetContentDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("ipAddress", m.GetIpAddress())
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
        err := writer.WriteStringValue("osPlatformDeviceDetails", m.GetOsPlatformDeviceDetails())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("potentialScoreImpact", m.GetPotentialScoreImpact())
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
func (m *UserTrainingContentEventInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetBrowser sets the browser property value. Browser of the user from where the training event was generated.
func (m *UserTrainingContentEventInfo) SetBrowser(value *string)() {
    m.browser = value
}
// SetContentDateTime sets the contentDateTime property value. Date and time of the training content playback by the user.
func (m *UserTrainingContentEventInfo) SetContentDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.contentDateTime = value
}
// SetIpAddress sets the ipAddress property value. IP address of the user for the training event.
func (m *UserTrainingContentEventInfo) SetIpAddress(value *string)() {
    m.ipAddress = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserTrainingContentEventInfo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetOsPlatformDeviceDetails sets the osPlatformDeviceDetails property value. The operating system, platform, and device details of the user for the training event.
func (m *UserTrainingContentEventInfo) SetOsPlatformDeviceDetails(value *string)() {
    m.osPlatformDeviceDetails = value
}
// SetPotentialScoreImpact sets the potentialScoreImpact property value. Potential improvement in the tenant security posture after completion of the training by the user.
func (m *UserTrainingContentEventInfo) SetPotentialScoreImpact(value *float64)() {
    m.potentialScoreImpact = value
}
