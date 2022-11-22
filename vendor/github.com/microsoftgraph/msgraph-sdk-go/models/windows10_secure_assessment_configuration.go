package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10SecureAssessmentConfiguration 
type Windows10SecureAssessmentConfiguration struct {
    DeviceConfiguration
    // Indicates whether or not to allow the app from printing during the test.
    allowPrinting *bool
    // Indicates whether or not to allow screen capture capability during a test.
    allowScreenCapture *bool
    // Indicates whether or not to allow text suggestions during the test.
    allowTextSuggestion *bool
    // The account used to configure the Windows device for taking the test. The user can be a domain account (domain/user), an AAD account (username@tenant.com) or a local account (username).
    configurationAccount *string
    // Url link to an assessment that's automatically loaded when the secure assessment browser is launched. It has to be a valid Url (http[s]://msdn.microsoft.com/).
    launchUri *string
}
// NewWindows10SecureAssessmentConfiguration instantiates a new Windows10SecureAssessmentConfiguration and sets the default values.
func NewWindows10SecureAssessmentConfiguration()(*Windows10SecureAssessmentConfiguration) {
    m := &Windows10SecureAssessmentConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windows10SecureAssessmentConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWindows10SecureAssessmentConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindows10SecureAssessmentConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10SecureAssessmentConfiguration(), nil
}
// GetAllowPrinting gets the allowPrinting property value. Indicates whether or not to allow the app from printing during the test.
func (m *Windows10SecureAssessmentConfiguration) GetAllowPrinting()(*bool) {
    return m.allowPrinting
}
// GetAllowScreenCapture gets the allowScreenCapture property value. Indicates whether or not to allow screen capture capability during a test.
func (m *Windows10SecureAssessmentConfiguration) GetAllowScreenCapture()(*bool) {
    return m.allowScreenCapture
}
// GetAllowTextSuggestion gets the allowTextSuggestion property value. Indicates whether or not to allow text suggestions during the test.
func (m *Windows10SecureAssessmentConfiguration) GetAllowTextSuggestion()(*bool) {
    return m.allowTextSuggestion
}
// GetConfigurationAccount gets the configurationAccount property value. The account used to configure the Windows device for taking the test. The user can be a domain account (domain/user), an AAD account (username@tenant.com) or a local account (username).
func (m *Windows10SecureAssessmentConfiguration) GetConfigurationAccount()(*string) {
    return m.configurationAccount
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Windows10SecureAssessmentConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["allowPrinting"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowPrinting)
    res["allowScreenCapture"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowScreenCapture)
    res["allowTextSuggestion"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAllowTextSuggestion)
    res["configurationAccount"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetConfigurationAccount)
    res["launchUri"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLaunchUri)
    return res
}
// GetLaunchUri gets the launchUri property value. Url link to an assessment that's automatically loaded when the secure assessment browser is launched. It has to be a valid Url (http[s]://msdn.microsoft.com/).
func (m *Windows10SecureAssessmentConfiguration) GetLaunchUri()(*string) {
    return m.launchUri
}
// Serialize serializes information the current object
func (m *Windows10SecureAssessmentConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowPrinting", m.GetAllowPrinting())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowScreenCapture", m.GetAllowScreenCapture())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowTextSuggestion", m.GetAllowTextSuggestion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("configurationAccount", m.GetConfigurationAccount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("launchUri", m.GetLaunchUri())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowPrinting sets the allowPrinting property value. Indicates whether or not to allow the app from printing during the test.
func (m *Windows10SecureAssessmentConfiguration) SetAllowPrinting(value *bool)() {
    m.allowPrinting = value
}
// SetAllowScreenCapture sets the allowScreenCapture property value. Indicates whether or not to allow screen capture capability during a test.
func (m *Windows10SecureAssessmentConfiguration) SetAllowScreenCapture(value *bool)() {
    m.allowScreenCapture = value
}
// SetAllowTextSuggestion sets the allowTextSuggestion property value. Indicates whether or not to allow text suggestions during the test.
func (m *Windows10SecureAssessmentConfiguration) SetAllowTextSuggestion(value *bool)() {
    m.allowTextSuggestion = value
}
// SetConfigurationAccount sets the configurationAccount property value. The account used to configure the Windows device for taking the test. The user can be a domain account (domain/user), an AAD account (username@tenant.com) or a local account (username).
func (m *Windows10SecureAssessmentConfiguration) SetConfigurationAccount(value *string)() {
    m.configurationAccount = value
}
// SetLaunchUri sets the launchUri property value. Url link to an assessment that's automatically loaded when the secure assessment browser is launched. It has to be a valid Url (http[s]://msdn.microsoft.com/).
func (m *Windows10SecureAssessmentConfiguration) SetLaunchUri(value *string)() {
    m.launchUri = value
}
