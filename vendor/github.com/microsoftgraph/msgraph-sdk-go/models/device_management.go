package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagement 
type DeviceManagement struct {
    Entity
}
// NewDeviceManagement instantiates a new DeviceManagement and sets the default values.
func NewDeviceManagement()(*DeviceManagement) {
    m := &DeviceManagement{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceManagementFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDeviceManagementFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceManagement(), nil
}
// GetApplePushNotificationCertificate gets the applePushNotificationCertificate property value. Apple push notification certificate.
func (m *DeviceManagement) GetApplePushNotificationCertificate()(ApplePushNotificationCertificateable) {
    val, err := m.GetBackingStore().Get("applePushNotificationCertificate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ApplePushNotificationCertificateable)
    }
    return nil
}
// GetAuditEvents gets the auditEvents property value. The Audit Events
func (m *DeviceManagement) GetAuditEvents()([]AuditEventable) {
    val, err := m.GetBackingStore().Get("auditEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuditEventable)
    }
    return nil
}
// GetComplianceManagementPartners gets the complianceManagementPartners property value. The list of Compliance Management Partners configured by the tenant.
func (m *DeviceManagement) GetComplianceManagementPartners()([]ComplianceManagementPartnerable) {
    val, err := m.GetBackingStore().Get("complianceManagementPartners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ComplianceManagementPartnerable)
    }
    return nil
}
// GetConditionalAccessSettings gets the conditionalAccessSettings property value. The Exchange on premises conditional access settings. On premises conditional access will require devices to be both enrolled and compliant for mail access
func (m *DeviceManagement) GetConditionalAccessSettings()(OnPremisesConditionalAccessSettingsable) {
    val, err := m.GetBackingStore().Get("conditionalAccessSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnPremisesConditionalAccessSettingsable)
    }
    return nil
}
// GetDetectedApps gets the detectedApps property value. The list of detected apps associated with a device.
func (m *DeviceManagement) GetDetectedApps()([]DetectedAppable) {
    val, err := m.GetBackingStore().Get("detectedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DetectedAppable)
    }
    return nil
}
// GetDeviceCategories gets the deviceCategories property value. The list of device categories with the tenant.
func (m *DeviceManagement) GetDeviceCategories()([]DeviceCategoryable) {
    val, err := m.GetBackingStore().Get("deviceCategories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceCategoryable)
    }
    return nil
}
// GetDeviceCompliancePolicies gets the deviceCompliancePolicies property value. The device compliance policies.
func (m *DeviceManagement) GetDeviceCompliancePolicies()([]DeviceCompliancePolicyable) {
    val, err := m.GetBackingStore().Get("deviceCompliancePolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceCompliancePolicyable)
    }
    return nil
}
// GetDeviceCompliancePolicyDeviceStateSummary gets the deviceCompliancePolicyDeviceStateSummary property value. The device compliance state summary for this account.
func (m *DeviceManagement) GetDeviceCompliancePolicyDeviceStateSummary()(DeviceCompliancePolicyDeviceStateSummaryable) {
    val, err := m.GetBackingStore().Get("deviceCompliancePolicyDeviceStateSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceCompliancePolicyDeviceStateSummaryable)
    }
    return nil
}
// GetDeviceCompliancePolicySettingStateSummaries gets the deviceCompliancePolicySettingStateSummaries property value. The summary states of compliance policy settings for this account.
func (m *DeviceManagement) GetDeviceCompliancePolicySettingStateSummaries()([]DeviceCompliancePolicySettingStateSummaryable) {
    val, err := m.GetBackingStore().Get("deviceCompliancePolicySettingStateSummaries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceCompliancePolicySettingStateSummaryable)
    }
    return nil
}
// GetDeviceConfigurationDeviceStateSummaries gets the deviceConfigurationDeviceStateSummaries property value. The device configuration device state summary for this account.
func (m *DeviceManagement) GetDeviceConfigurationDeviceStateSummaries()(DeviceConfigurationDeviceStateSummaryable) {
    val, err := m.GetBackingStore().Get("deviceConfigurationDeviceStateSummaries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceConfigurationDeviceStateSummaryable)
    }
    return nil
}
// GetDeviceConfigurations gets the deviceConfigurations property value. The device configurations.
func (m *DeviceManagement) GetDeviceConfigurations()([]DeviceConfigurationable) {
    val, err := m.GetBackingStore().Get("deviceConfigurations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceConfigurationable)
    }
    return nil
}
// GetDeviceEnrollmentConfigurations gets the deviceEnrollmentConfigurations property value. The list of device enrollment configurations
func (m *DeviceManagement) GetDeviceEnrollmentConfigurations()([]DeviceEnrollmentConfigurationable) {
    val, err := m.GetBackingStore().Get("deviceEnrollmentConfigurations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceEnrollmentConfigurationable)
    }
    return nil
}
// GetDeviceManagementPartners gets the deviceManagementPartners property value. The list of Device Management Partners configured by the tenant.
func (m *DeviceManagement) GetDeviceManagementPartners()([]DeviceManagementPartnerable) {
    val, err := m.GetBackingStore().Get("deviceManagementPartners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceManagementPartnerable)
    }
    return nil
}
// GetExchangeConnectors gets the exchangeConnectors property value. The list of Exchange Connectors configured by the tenant.
func (m *DeviceManagement) GetExchangeConnectors()([]DeviceManagementExchangeConnectorable) {
    val, err := m.GetBackingStore().Get("exchangeConnectors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceManagementExchangeConnectorable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceManagement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["applePushNotificationCertificate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateApplePushNotificationCertificateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplePushNotificationCertificate(val.(ApplePushNotificationCertificateable))
        }
        return nil
    }
    res["auditEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuditEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuditEventable, len(val))
            for i, v := range val {
                res[i] = v.(AuditEventable)
            }
            m.SetAuditEvents(res)
        }
        return nil
    }
    res["complianceManagementPartners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateComplianceManagementPartnerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ComplianceManagementPartnerable, len(val))
            for i, v := range val {
                res[i] = v.(ComplianceManagementPartnerable)
            }
            m.SetComplianceManagementPartners(res)
        }
        return nil
    }
    res["conditionalAccessSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnPremisesConditionalAccessSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConditionalAccessSettings(val.(OnPremisesConditionalAccessSettingsable))
        }
        return nil
    }
    res["detectedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDetectedAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DetectedAppable, len(val))
            for i, v := range val {
                res[i] = v.(DetectedAppable)
            }
            m.SetDetectedApps(res)
        }
        return nil
    }
    res["deviceCategories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceCategoryable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceCategoryable)
            }
            m.SetDeviceCategories(res)
        }
        return nil
    }
    res["deviceCompliancePolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceCompliancePolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceCompliancePolicyable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceCompliancePolicyable)
            }
            m.SetDeviceCompliancePolicies(res)
        }
        return nil
    }
    res["deviceCompliancePolicyDeviceStateSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceCompliancePolicyDeviceStateSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceCompliancePolicyDeviceStateSummary(val.(DeviceCompliancePolicyDeviceStateSummaryable))
        }
        return nil
    }
    res["deviceCompliancePolicySettingStateSummaries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceCompliancePolicySettingStateSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceCompliancePolicySettingStateSummaryable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceCompliancePolicySettingStateSummaryable)
            }
            m.SetDeviceCompliancePolicySettingStateSummaries(res)
        }
        return nil
    }
    res["deviceConfigurationDeviceStateSummaries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceConfigurationDeviceStateSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceConfigurationDeviceStateSummaries(val.(DeviceConfigurationDeviceStateSummaryable))
        }
        return nil
    }
    res["deviceConfigurations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceConfigurationable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceConfigurationable)
            }
            m.SetDeviceConfigurations(res)
        }
        return nil
    }
    res["deviceEnrollmentConfigurations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceEnrollmentConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceEnrollmentConfigurationable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceEnrollmentConfigurationable)
            }
            m.SetDeviceEnrollmentConfigurations(res)
        }
        return nil
    }
    res["deviceManagementPartners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceManagementPartnerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceManagementPartnerable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceManagementPartnerable)
            }
            m.SetDeviceManagementPartners(res)
        }
        return nil
    }
    res["exchangeConnectors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceManagementExchangeConnectorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceManagementExchangeConnectorable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceManagementExchangeConnectorable)
            }
            m.SetExchangeConnectors(res)
        }
        return nil
    }
    res["importedWindowsAutopilotDeviceIdentities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateImportedWindowsAutopilotDeviceIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ImportedWindowsAutopilotDeviceIdentityable, len(val))
            for i, v := range val {
                res[i] = v.(ImportedWindowsAutopilotDeviceIdentityable)
            }
            m.SetImportedWindowsAutopilotDeviceIdentities(res)
        }
        return nil
    }
    res["intuneAccountId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIntuneAccountId(val)
        }
        return nil
    }
    res["intuneBrand"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIntuneBrandFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIntuneBrand(val.(IntuneBrandable))
        }
        return nil
    }
    res["iosUpdateStatuses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIosUpdateDeviceStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IosUpdateDeviceStatusable, len(val))
            for i, v := range val {
                res[i] = v.(IosUpdateDeviceStatusable)
            }
            m.SetIosUpdateStatuses(res)
        }
        return nil
    }
    res["managedDeviceOverview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateManagedDeviceOverviewFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagedDeviceOverview(val.(ManagedDeviceOverviewable))
        }
        return nil
    }
    res["managedDevices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedDeviceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedDeviceable, len(val))
            for i, v := range val {
                res[i] = v.(ManagedDeviceable)
            }
            m.SetManagedDevices(res)
        }
        return nil
    }
    res["mobileThreatDefenseConnectors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMobileThreatDefenseConnectorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MobileThreatDefenseConnectorable, len(val))
            for i, v := range val {
                res[i] = v.(MobileThreatDefenseConnectorable)
            }
            m.SetMobileThreatDefenseConnectors(res)
        }
        return nil
    }
    res["notificationMessageTemplates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateNotificationMessageTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]NotificationMessageTemplateable, len(val))
            for i, v := range val {
                res[i] = v.(NotificationMessageTemplateable)
            }
            m.SetNotificationMessageTemplates(res)
        }
        return nil
    }
    res["remoteAssistancePartners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRemoteAssistancePartnerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RemoteAssistancePartnerable, len(val))
            for i, v := range val {
                res[i] = v.(RemoteAssistancePartnerable)
            }
            m.SetRemoteAssistancePartners(res)
        }
        return nil
    }
    res["reports"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceManagementReportsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReports(val.(DeviceManagementReportsable))
        }
        return nil
    }
    res["resourceOperations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateResourceOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ResourceOperationable, len(val))
            for i, v := range val {
                res[i] = v.(ResourceOperationable)
            }
            m.SetResourceOperations(res)
        }
        return nil
    }
    res["roleAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceAndAppManagementRoleAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceAndAppManagementRoleAssignmentable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceAndAppManagementRoleAssignmentable)
            }
            m.SetRoleAssignments(res)
        }
        return nil
    }
    res["roleDefinitions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRoleDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RoleDefinitionable, len(val))
            for i, v := range val {
                res[i] = v.(RoleDefinitionable)
            }
            m.SetRoleDefinitions(res)
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceManagementSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettings(val.(DeviceManagementSettingsable))
        }
        return nil
    }
    res["softwareUpdateStatusSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSoftwareUpdateStatusSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSoftwareUpdateStatusSummary(val.(SoftwareUpdateStatusSummaryable))
        }
        return nil
    }
    res["subscriptionState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementSubscriptionState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubscriptionState(val.(*DeviceManagementSubscriptionState))
        }
        return nil
    }
    res["telecomExpenseManagementPartners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTelecomExpenseManagementPartnerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TelecomExpenseManagementPartnerable, len(val))
            for i, v := range val {
                res[i] = v.(TelecomExpenseManagementPartnerable)
            }
            m.SetTelecomExpenseManagementPartners(res)
        }
        return nil
    }
    res["termsAndConditions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTermsAndConditionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TermsAndConditionsable, len(val))
            for i, v := range val {
                res[i] = v.(TermsAndConditionsable)
            }
            m.SetTermsAndConditions(res)
        }
        return nil
    }
    res["troubleshootingEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceManagementTroubleshootingEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceManagementTroubleshootingEventable, len(val))
            for i, v := range val {
                res[i] = v.(DeviceManagementTroubleshootingEventable)
            }
            m.SetTroubleshootingEvents(res)
        }
        return nil
    }
    res["windowsAutopilotDeviceIdentities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsAutopilotDeviceIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsAutopilotDeviceIdentityable, len(val))
            for i, v := range val {
                res[i] = v.(WindowsAutopilotDeviceIdentityable)
            }
            m.SetWindowsAutopilotDeviceIdentities(res)
        }
        return nil
    }
    res["windowsInformationProtectionAppLearningSummaries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionAppLearningSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionAppLearningSummaryable, len(val))
            for i, v := range val {
                res[i] = v.(WindowsInformationProtectionAppLearningSummaryable)
            }
            m.SetWindowsInformationProtectionAppLearningSummaries(res)
        }
        return nil
    }
    res["windowsInformationProtectionNetworkLearningSummaries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionNetworkLearningSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionNetworkLearningSummaryable, len(val))
            for i, v := range val {
                res[i] = v.(WindowsInformationProtectionNetworkLearningSummaryable)
            }
            m.SetWindowsInformationProtectionNetworkLearningSummaries(res)
        }
        return nil
    }
    return res
}
// GetImportedWindowsAutopilotDeviceIdentities gets the importedWindowsAutopilotDeviceIdentities property value. Collection of imported Windows autopilot devices.
func (m *DeviceManagement) GetImportedWindowsAutopilotDeviceIdentities()([]ImportedWindowsAutopilotDeviceIdentityable) {
    val, err := m.GetBackingStore().Get("importedWindowsAutopilotDeviceIdentities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ImportedWindowsAutopilotDeviceIdentityable)
    }
    return nil
}
// GetIntuneAccountId gets the intuneAccountId property value. Intune Account Id for given tenant
func (m *DeviceManagement) GetIntuneAccountId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("intuneAccountId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetIntuneBrand gets the intuneBrand property value. intuneBrand contains data which is used in customizing the appearance of the Company Portal applications as well as the end user web portal.
func (m *DeviceManagement) GetIntuneBrand()(IntuneBrandable) {
    val, err := m.GetBackingStore().Get("intuneBrand")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IntuneBrandable)
    }
    return nil
}
// GetIosUpdateStatuses gets the iosUpdateStatuses property value. The IOS software update installation statuses for this account.
func (m *DeviceManagement) GetIosUpdateStatuses()([]IosUpdateDeviceStatusable) {
    val, err := m.GetBackingStore().Get("iosUpdateStatuses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IosUpdateDeviceStatusable)
    }
    return nil
}
// GetManagedDeviceOverview gets the managedDeviceOverview property value. Device overview
func (m *DeviceManagement) GetManagedDeviceOverview()(ManagedDeviceOverviewable) {
    val, err := m.GetBackingStore().Get("managedDeviceOverview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ManagedDeviceOverviewable)
    }
    return nil
}
// GetManagedDevices gets the managedDevices property value. The list of managed devices.
func (m *DeviceManagement) GetManagedDevices()([]ManagedDeviceable) {
    val, err := m.GetBackingStore().Get("managedDevices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedDeviceable)
    }
    return nil
}
// GetMobileThreatDefenseConnectors gets the mobileThreatDefenseConnectors property value. The list of Mobile threat Defense connectors configured by the tenant.
func (m *DeviceManagement) GetMobileThreatDefenseConnectors()([]MobileThreatDefenseConnectorable) {
    val, err := m.GetBackingStore().Get("mobileThreatDefenseConnectors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MobileThreatDefenseConnectorable)
    }
    return nil
}
// GetNotificationMessageTemplates gets the notificationMessageTemplates property value. The Notification Message Templates.
func (m *DeviceManagement) GetNotificationMessageTemplates()([]NotificationMessageTemplateable) {
    val, err := m.GetBackingStore().Get("notificationMessageTemplates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]NotificationMessageTemplateable)
    }
    return nil
}
// GetRemoteAssistancePartners gets the remoteAssistancePartners property value. The remote assist partners.
func (m *DeviceManagement) GetRemoteAssistancePartners()([]RemoteAssistancePartnerable) {
    val, err := m.GetBackingStore().Get("remoteAssistancePartners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RemoteAssistancePartnerable)
    }
    return nil
}
// GetReports gets the reports property value. Reports singleton
func (m *DeviceManagement) GetReports()(DeviceManagementReportsable) {
    val, err := m.GetBackingStore().Get("reports")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceManagementReportsable)
    }
    return nil
}
// GetResourceOperations gets the resourceOperations property value. The Resource Operations.
func (m *DeviceManagement) GetResourceOperations()([]ResourceOperationable) {
    val, err := m.GetBackingStore().Get("resourceOperations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ResourceOperationable)
    }
    return nil
}
// GetRoleAssignments gets the roleAssignments property value. The Role Assignments.
func (m *DeviceManagement) GetRoleAssignments()([]DeviceAndAppManagementRoleAssignmentable) {
    val, err := m.GetBackingStore().Get("roleAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceAndAppManagementRoleAssignmentable)
    }
    return nil
}
// GetRoleDefinitions gets the roleDefinitions property value. The Role Definitions.
func (m *DeviceManagement) GetRoleDefinitions()([]RoleDefinitionable) {
    val, err := m.GetBackingStore().Get("roleDefinitions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RoleDefinitionable)
    }
    return nil
}
// GetSettings gets the settings property value. Account level settings.
func (m *DeviceManagement) GetSettings()(DeviceManagementSettingsable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceManagementSettingsable)
    }
    return nil
}
// GetSoftwareUpdateStatusSummary gets the softwareUpdateStatusSummary property value. The software update status summary.
func (m *DeviceManagement) GetSoftwareUpdateStatusSummary()(SoftwareUpdateStatusSummaryable) {
    val, err := m.GetBackingStore().Get("softwareUpdateStatusSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SoftwareUpdateStatusSummaryable)
    }
    return nil
}
// GetSubscriptionState gets the subscriptionState property value. Tenant mobile device management subscription state.
func (m *DeviceManagement) GetSubscriptionState()(*DeviceManagementSubscriptionState) {
    val, err := m.GetBackingStore().Get("subscriptionState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementSubscriptionState)
    }
    return nil
}
// GetTelecomExpenseManagementPartners gets the telecomExpenseManagementPartners property value. The telecom expense management partners.
func (m *DeviceManagement) GetTelecomExpenseManagementPartners()([]TelecomExpenseManagementPartnerable) {
    val, err := m.GetBackingStore().Get("telecomExpenseManagementPartners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TelecomExpenseManagementPartnerable)
    }
    return nil
}
// GetTermsAndConditions gets the termsAndConditions property value. The terms and conditions associated with device management of the company.
func (m *DeviceManagement) GetTermsAndConditions()([]TermsAndConditionsable) {
    val, err := m.GetBackingStore().Get("termsAndConditions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TermsAndConditionsable)
    }
    return nil
}
// GetTroubleshootingEvents gets the troubleshootingEvents property value. The list of troubleshooting events for the tenant.
func (m *DeviceManagement) GetTroubleshootingEvents()([]DeviceManagementTroubleshootingEventable) {
    val, err := m.GetBackingStore().Get("troubleshootingEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceManagementTroubleshootingEventable)
    }
    return nil
}
// GetWindowsAutopilotDeviceIdentities gets the windowsAutopilotDeviceIdentities property value. The Windows autopilot device identities contained collection.
func (m *DeviceManagement) GetWindowsAutopilotDeviceIdentities()([]WindowsAutopilotDeviceIdentityable) {
    val, err := m.GetBackingStore().Get("windowsAutopilotDeviceIdentities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsAutopilotDeviceIdentityable)
    }
    return nil
}
// GetWindowsInformationProtectionAppLearningSummaries gets the windowsInformationProtectionAppLearningSummaries property value. The windows information protection app learning summaries.
func (m *DeviceManagement) GetWindowsInformationProtectionAppLearningSummaries()([]WindowsInformationProtectionAppLearningSummaryable) {
    val, err := m.GetBackingStore().Get("windowsInformationProtectionAppLearningSummaries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionAppLearningSummaryable)
    }
    return nil
}
// GetWindowsInformationProtectionNetworkLearningSummaries gets the windowsInformationProtectionNetworkLearningSummaries property value. The windows information protection network learning summaries.
func (m *DeviceManagement) GetWindowsInformationProtectionNetworkLearningSummaries()([]WindowsInformationProtectionNetworkLearningSummaryable) {
    val, err := m.GetBackingStore().Get("windowsInformationProtectionNetworkLearningSummaries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionNetworkLearningSummaryable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceManagement) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("applePushNotificationCertificate", m.GetApplePushNotificationCertificate())
        if err != nil {
            return err
        }
    }
    if m.GetAuditEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuditEvents()))
        for i, v := range m.GetAuditEvents() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("auditEvents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetComplianceManagementPartners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetComplianceManagementPartners()))
        for i, v := range m.GetComplianceManagementPartners() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("complianceManagementPartners", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("conditionalAccessSettings", m.GetConditionalAccessSettings())
        if err != nil {
            return err
        }
    }
    if m.GetDetectedApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDetectedApps()))
        for i, v := range m.GetDetectedApps() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("detectedApps", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceCategories() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceCategories()))
        for i, v := range m.GetDeviceCategories() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("deviceCategories", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceCompliancePolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceCompliancePolicies()))
        for i, v := range m.GetDeviceCompliancePolicies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("deviceCompliancePolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deviceCompliancePolicyDeviceStateSummary", m.GetDeviceCompliancePolicyDeviceStateSummary())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceCompliancePolicySettingStateSummaries() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceCompliancePolicySettingStateSummaries()))
        for i, v := range m.GetDeviceCompliancePolicySettingStateSummaries() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("deviceCompliancePolicySettingStateSummaries", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deviceConfigurationDeviceStateSummaries", m.GetDeviceConfigurationDeviceStateSummaries())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceConfigurations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceConfigurations()))
        for i, v := range m.GetDeviceConfigurations() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("deviceConfigurations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceEnrollmentConfigurations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceEnrollmentConfigurations()))
        for i, v := range m.GetDeviceEnrollmentConfigurations() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("deviceEnrollmentConfigurations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceManagementPartners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceManagementPartners()))
        for i, v := range m.GetDeviceManagementPartners() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("deviceManagementPartners", cast)
        if err != nil {
            return err
        }
    }
    if m.GetExchangeConnectors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExchangeConnectors()))
        for i, v := range m.GetExchangeConnectors() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("exchangeConnectors", cast)
        if err != nil {
            return err
        }
    }
    if m.GetImportedWindowsAutopilotDeviceIdentities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetImportedWindowsAutopilotDeviceIdentities()))
        for i, v := range m.GetImportedWindowsAutopilotDeviceIdentities() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("importedWindowsAutopilotDeviceIdentities", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("intuneAccountId", m.GetIntuneAccountId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("intuneBrand", m.GetIntuneBrand())
        if err != nil {
            return err
        }
    }
    if m.GetIosUpdateStatuses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIosUpdateStatuses()))
        for i, v := range m.GetIosUpdateStatuses() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("iosUpdateStatuses", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("managedDeviceOverview", m.GetManagedDeviceOverview())
        if err != nil {
            return err
        }
    }
    if m.GetManagedDevices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetManagedDevices()))
        for i, v := range m.GetManagedDevices() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("managedDevices", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMobileThreatDefenseConnectors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMobileThreatDefenseConnectors()))
        for i, v := range m.GetMobileThreatDefenseConnectors() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("mobileThreatDefenseConnectors", cast)
        if err != nil {
            return err
        }
    }
    if m.GetNotificationMessageTemplates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNotificationMessageTemplates()))
        for i, v := range m.GetNotificationMessageTemplates() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("notificationMessageTemplates", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRemoteAssistancePartners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRemoteAssistancePartners()))
        for i, v := range m.GetRemoteAssistancePartners() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("remoteAssistancePartners", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("reports", m.GetReports())
        if err != nil {
            return err
        }
    }
    if m.GetResourceOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResourceOperations()))
        for i, v := range m.GetResourceOperations() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("resourceOperations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleAssignments()))
        for i, v := range m.GetRoleAssignments() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("roleAssignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleDefinitions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleDefinitions()))
        for i, v := range m.GetRoleDefinitions() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("roleDefinitions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("softwareUpdateStatusSummary", m.GetSoftwareUpdateStatusSummary())
        if err != nil {
            return err
        }
    }
    if m.GetSubscriptionState() != nil {
        cast := (*m.GetSubscriptionState()).String()
        err = writer.WriteStringValue("subscriptionState", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetTelecomExpenseManagementPartners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTelecomExpenseManagementPartners()))
        for i, v := range m.GetTelecomExpenseManagementPartners() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("telecomExpenseManagementPartners", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTermsAndConditions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTermsAndConditions()))
        for i, v := range m.GetTermsAndConditions() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("termsAndConditions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTroubleshootingEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTroubleshootingEvents()))
        for i, v := range m.GetTroubleshootingEvents() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("troubleshootingEvents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsAutopilotDeviceIdentities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWindowsAutopilotDeviceIdentities()))
        for i, v := range m.GetWindowsAutopilotDeviceIdentities() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("windowsAutopilotDeviceIdentities", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsInformationProtectionAppLearningSummaries() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWindowsInformationProtectionAppLearningSummaries()))
        for i, v := range m.GetWindowsInformationProtectionAppLearningSummaries() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("windowsInformationProtectionAppLearningSummaries", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsInformationProtectionNetworkLearningSummaries() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWindowsInformationProtectionNetworkLearningSummaries()))
        for i, v := range m.GetWindowsInformationProtectionNetworkLearningSummaries() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("windowsInformationProtectionNetworkLearningSummaries", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplePushNotificationCertificate sets the applePushNotificationCertificate property value. Apple push notification certificate.
func (m *DeviceManagement) SetApplePushNotificationCertificate(value ApplePushNotificationCertificateable)() {
    err := m.GetBackingStore().Set("applePushNotificationCertificate", value)
    if err != nil {
        panic(err)
    }
}
// SetAuditEvents sets the auditEvents property value. The Audit Events
func (m *DeviceManagement) SetAuditEvents(value []AuditEventable)() {
    err := m.GetBackingStore().Set("auditEvents", value)
    if err != nil {
        panic(err)
    }
}
// SetComplianceManagementPartners sets the complianceManagementPartners property value. The list of Compliance Management Partners configured by the tenant.
func (m *DeviceManagement) SetComplianceManagementPartners(value []ComplianceManagementPartnerable)() {
    err := m.GetBackingStore().Set("complianceManagementPartners", value)
    if err != nil {
        panic(err)
    }
}
// SetConditionalAccessSettings sets the conditionalAccessSettings property value. The Exchange on premises conditional access settings. On premises conditional access will require devices to be both enrolled and compliant for mail access
func (m *DeviceManagement) SetConditionalAccessSettings(value OnPremisesConditionalAccessSettingsable)() {
    err := m.GetBackingStore().Set("conditionalAccessSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetDetectedApps sets the detectedApps property value. The list of detected apps associated with a device.
func (m *DeviceManagement) SetDetectedApps(value []DetectedAppable)() {
    err := m.GetBackingStore().Set("detectedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCategories sets the deviceCategories property value. The list of device categories with the tenant.
func (m *DeviceManagement) SetDeviceCategories(value []DeviceCategoryable)() {
    err := m.GetBackingStore().Set("deviceCategories", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCompliancePolicies sets the deviceCompliancePolicies property value. The device compliance policies.
func (m *DeviceManagement) SetDeviceCompliancePolicies(value []DeviceCompliancePolicyable)() {
    err := m.GetBackingStore().Set("deviceCompliancePolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCompliancePolicyDeviceStateSummary sets the deviceCompliancePolicyDeviceStateSummary property value. The device compliance state summary for this account.
func (m *DeviceManagement) SetDeviceCompliancePolicyDeviceStateSummary(value DeviceCompliancePolicyDeviceStateSummaryable)() {
    err := m.GetBackingStore().Set("deviceCompliancePolicyDeviceStateSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceCompliancePolicySettingStateSummaries sets the deviceCompliancePolicySettingStateSummaries property value. The summary states of compliance policy settings for this account.
func (m *DeviceManagement) SetDeviceCompliancePolicySettingStateSummaries(value []DeviceCompliancePolicySettingStateSummaryable)() {
    err := m.GetBackingStore().Set("deviceCompliancePolicySettingStateSummaries", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceConfigurationDeviceStateSummaries sets the deviceConfigurationDeviceStateSummaries property value. The device configuration device state summary for this account.
func (m *DeviceManagement) SetDeviceConfigurationDeviceStateSummaries(value DeviceConfigurationDeviceStateSummaryable)() {
    err := m.GetBackingStore().Set("deviceConfigurationDeviceStateSummaries", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceConfigurations sets the deviceConfigurations property value. The device configurations.
func (m *DeviceManagement) SetDeviceConfigurations(value []DeviceConfigurationable)() {
    err := m.GetBackingStore().Set("deviceConfigurations", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceEnrollmentConfigurations sets the deviceEnrollmentConfigurations property value. The list of device enrollment configurations
func (m *DeviceManagement) SetDeviceEnrollmentConfigurations(value []DeviceEnrollmentConfigurationable)() {
    err := m.GetBackingStore().Set("deviceEnrollmentConfigurations", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceManagementPartners sets the deviceManagementPartners property value. The list of Device Management Partners configured by the tenant.
func (m *DeviceManagement) SetDeviceManagementPartners(value []DeviceManagementPartnerable)() {
    err := m.GetBackingStore().Set("deviceManagementPartners", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeConnectors sets the exchangeConnectors property value. The list of Exchange Connectors configured by the tenant.
func (m *DeviceManagement) SetExchangeConnectors(value []DeviceManagementExchangeConnectorable)() {
    err := m.GetBackingStore().Set("exchangeConnectors", value)
    if err != nil {
        panic(err)
    }
}
// SetImportedWindowsAutopilotDeviceIdentities sets the importedWindowsAutopilotDeviceIdentities property value. Collection of imported Windows autopilot devices.
func (m *DeviceManagement) SetImportedWindowsAutopilotDeviceIdentities(value []ImportedWindowsAutopilotDeviceIdentityable)() {
    err := m.GetBackingStore().Set("importedWindowsAutopilotDeviceIdentities", value)
    if err != nil {
        panic(err)
    }
}
// SetIntuneAccountId sets the intuneAccountId property value. Intune Account Id for given tenant
func (m *DeviceManagement) SetIntuneAccountId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("intuneAccountId", value)
    if err != nil {
        panic(err)
    }
}
// SetIntuneBrand sets the intuneBrand property value. intuneBrand contains data which is used in customizing the appearance of the Company Portal applications as well as the end user web portal.
func (m *DeviceManagement) SetIntuneBrand(value IntuneBrandable)() {
    err := m.GetBackingStore().Set("intuneBrand", value)
    if err != nil {
        panic(err)
    }
}
// SetIosUpdateStatuses sets the iosUpdateStatuses property value. The IOS software update installation statuses for this account.
func (m *DeviceManagement) SetIosUpdateStatuses(value []IosUpdateDeviceStatusable)() {
    err := m.GetBackingStore().Set("iosUpdateStatuses", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedDeviceOverview sets the managedDeviceOverview property value. Device overview
func (m *DeviceManagement) SetManagedDeviceOverview(value ManagedDeviceOverviewable)() {
    err := m.GetBackingStore().Set("managedDeviceOverview", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedDevices sets the managedDevices property value. The list of managed devices.
func (m *DeviceManagement) SetManagedDevices(value []ManagedDeviceable)() {
    err := m.GetBackingStore().Set("managedDevices", value)
    if err != nil {
        panic(err)
    }
}
// SetMobileThreatDefenseConnectors sets the mobileThreatDefenseConnectors property value. The list of Mobile threat Defense connectors configured by the tenant.
func (m *DeviceManagement) SetMobileThreatDefenseConnectors(value []MobileThreatDefenseConnectorable)() {
    err := m.GetBackingStore().Set("mobileThreatDefenseConnectors", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationMessageTemplates sets the notificationMessageTemplates property value. The Notification Message Templates.
func (m *DeviceManagement) SetNotificationMessageTemplates(value []NotificationMessageTemplateable)() {
    err := m.GetBackingStore().Set("notificationMessageTemplates", value)
    if err != nil {
        panic(err)
    }
}
// SetRemoteAssistancePartners sets the remoteAssistancePartners property value. The remote assist partners.
func (m *DeviceManagement) SetRemoteAssistancePartners(value []RemoteAssistancePartnerable)() {
    err := m.GetBackingStore().Set("remoteAssistancePartners", value)
    if err != nil {
        panic(err)
    }
}
// SetReports sets the reports property value. Reports singleton
func (m *DeviceManagement) SetReports(value DeviceManagementReportsable)() {
    err := m.GetBackingStore().Set("reports", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceOperations sets the resourceOperations property value. The Resource Operations.
func (m *DeviceManagement) SetResourceOperations(value []ResourceOperationable)() {
    err := m.GetBackingStore().Set("resourceOperations", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleAssignments sets the roleAssignments property value. The Role Assignments.
func (m *DeviceManagement) SetRoleAssignments(value []DeviceAndAppManagementRoleAssignmentable)() {
    err := m.GetBackingStore().Set("roleAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleDefinitions sets the roleDefinitions property value. The Role Definitions.
func (m *DeviceManagement) SetRoleDefinitions(value []RoleDefinitionable)() {
    err := m.GetBackingStore().Set("roleDefinitions", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. Account level settings.
func (m *DeviceManagement) SetSettings(value DeviceManagementSettingsable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
// SetSoftwareUpdateStatusSummary sets the softwareUpdateStatusSummary property value. The software update status summary.
func (m *DeviceManagement) SetSoftwareUpdateStatusSummary(value SoftwareUpdateStatusSummaryable)() {
    err := m.GetBackingStore().Set("softwareUpdateStatusSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetSubscriptionState sets the subscriptionState property value. Tenant mobile device management subscription state.
func (m *DeviceManagement) SetSubscriptionState(value *DeviceManagementSubscriptionState)() {
    err := m.GetBackingStore().Set("subscriptionState", value)
    if err != nil {
        panic(err)
    }
}
// SetTelecomExpenseManagementPartners sets the telecomExpenseManagementPartners property value. The telecom expense management partners.
func (m *DeviceManagement) SetTelecomExpenseManagementPartners(value []TelecomExpenseManagementPartnerable)() {
    err := m.GetBackingStore().Set("telecomExpenseManagementPartners", value)
    if err != nil {
        panic(err)
    }
}
// SetTermsAndConditions sets the termsAndConditions property value. The terms and conditions associated with device management of the company.
func (m *DeviceManagement) SetTermsAndConditions(value []TermsAndConditionsable)() {
    err := m.GetBackingStore().Set("termsAndConditions", value)
    if err != nil {
        panic(err)
    }
}
// SetTroubleshootingEvents sets the troubleshootingEvents property value. The list of troubleshooting events for the tenant.
func (m *DeviceManagement) SetTroubleshootingEvents(value []DeviceManagementTroubleshootingEventable)() {
    err := m.GetBackingStore().Set("troubleshootingEvents", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsAutopilotDeviceIdentities sets the windowsAutopilotDeviceIdentities property value. The Windows autopilot device identities contained collection.
func (m *DeviceManagement) SetWindowsAutopilotDeviceIdentities(value []WindowsAutopilotDeviceIdentityable)() {
    err := m.GetBackingStore().Set("windowsAutopilotDeviceIdentities", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsInformationProtectionAppLearningSummaries sets the windowsInformationProtectionAppLearningSummaries property value. The windows information protection app learning summaries.
func (m *DeviceManagement) SetWindowsInformationProtectionAppLearningSummaries(value []WindowsInformationProtectionAppLearningSummaryable)() {
    err := m.GetBackingStore().Set("windowsInformationProtectionAppLearningSummaries", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsInformationProtectionNetworkLearningSummaries sets the windowsInformationProtectionNetworkLearningSummaries property value. The windows information protection network learning summaries.
func (m *DeviceManagement) SetWindowsInformationProtectionNetworkLearningSummaries(value []WindowsInformationProtectionNetworkLearningSummaryable)() {
    err := m.GetBackingStore().Set("windowsInformationProtectionNetworkLearningSummaries", value)
    if err != nil {
        panic(err)
    }
}
// DeviceManagementable 
type DeviceManagementable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplePushNotificationCertificate()(ApplePushNotificationCertificateable)
    GetAuditEvents()([]AuditEventable)
    GetComplianceManagementPartners()([]ComplianceManagementPartnerable)
    GetConditionalAccessSettings()(OnPremisesConditionalAccessSettingsable)
    GetDetectedApps()([]DetectedAppable)
    GetDeviceCategories()([]DeviceCategoryable)
    GetDeviceCompliancePolicies()([]DeviceCompliancePolicyable)
    GetDeviceCompliancePolicyDeviceStateSummary()(DeviceCompliancePolicyDeviceStateSummaryable)
    GetDeviceCompliancePolicySettingStateSummaries()([]DeviceCompliancePolicySettingStateSummaryable)
    GetDeviceConfigurationDeviceStateSummaries()(DeviceConfigurationDeviceStateSummaryable)
    GetDeviceConfigurations()([]DeviceConfigurationable)
    GetDeviceEnrollmentConfigurations()([]DeviceEnrollmentConfigurationable)
    GetDeviceManagementPartners()([]DeviceManagementPartnerable)
    GetExchangeConnectors()([]DeviceManagementExchangeConnectorable)
    GetImportedWindowsAutopilotDeviceIdentities()([]ImportedWindowsAutopilotDeviceIdentityable)
    GetIntuneAccountId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetIntuneBrand()(IntuneBrandable)
    GetIosUpdateStatuses()([]IosUpdateDeviceStatusable)
    GetManagedDeviceOverview()(ManagedDeviceOverviewable)
    GetManagedDevices()([]ManagedDeviceable)
    GetMobileThreatDefenseConnectors()([]MobileThreatDefenseConnectorable)
    GetNotificationMessageTemplates()([]NotificationMessageTemplateable)
    GetRemoteAssistancePartners()([]RemoteAssistancePartnerable)
    GetReports()(DeviceManagementReportsable)
    GetResourceOperations()([]ResourceOperationable)
    GetRoleAssignments()([]DeviceAndAppManagementRoleAssignmentable)
    GetRoleDefinitions()([]RoleDefinitionable)
    GetSettings()(DeviceManagementSettingsable)
    GetSoftwareUpdateStatusSummary()(SoftwareUpdateStatusSummaryable)
    GetSubscriptionState()(*DeviceManagementSubscriptionState)
    GetTelecomExpenseManagementPartners()([]TelecomExpenseManagementPartnerable)
    GetTermsAndConditions()([]TermsAndConditionsable)
    GetTroubleshootingEvents()([]DeviceManagementTroubleshootingEventable)
    GetWindowsAutopilotDeviceIdentities()([]WindowsAutopilotDeviceIdentityable)
    GetWindowsInformationProtectionAppLearningSummaries()([]WindowsInformationProtectionAppLearningSummaryable)
    GetWindowsInformationProtectionNetworkLearningSummaries()([]WindowsInformationProtectionNetworkLearningSummaryable)
    SetApplePushNotificationCertificate(value ApplePushNotificationCertificateable)()
    SetAuditEvents(value []AuditEventable)()
    SetComplianceManagementPartners(value []ComplianceManagementPartnerable)()
    SetConditionalAccessSettings(value OnPremisesConditionalAccessSettingsable)()
    SetDetectedApps(value []DetectedAppable)()
    SetDeviceCategories(value []DeviceCategoryable)()
    SetDeviceCompliancePolicies(value []DeviceCompliancePolicyable)()
    SetDeviceCompliancePolicyDeviceStateSummary(value DeviceCompliancePolicyDeviceStateSummaryable)()
    SetDeviceCompliancePolicySettingStateSummaries(value []DeviceCompliancePolicySettingStateSummaryable)()
    SetDeviceConfigurationDeviceStateSummaries(value DeviceConfigurationDeviceStateSummaryable)()
    SetDeviceConfigurations(value []DeviceConfigurationable)()
    SetDeviceEnrollmentConfigurations(value []DeviceEnrollmentConfigurationable)()
    SetDeviceManagementPartners(value []DeviceManagementPartnerable)()
    SetExchangeConnectors(value []DeviceManagementExchangeConnectorable)()
    SetImportedWindowsAutopilotDeviceIdentities(value []ImportedWindowsAutopilotDeviceIdentityable)()
    SetIntuneAccountId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetIntuneBrand(value IntuneBrandable)()
    SetIosUpdateStatuses(value []IosUpdateDeviceStatusable)()
    SetManagedDeviceOverview(value ManagedDeviceOverviewable)()
    SetManagedDevices(value []ManagedDeviceable)()
    SetMobileThreatDefenseConnectors(value []MobileThreatDefenseConnectorable)()
    SetNotificationMessageTemplates(value []NotificationMessageTemplateable)()
    SetRemoteAssistancePartners(value []RemoteAssistancePartnerable)()
    SetReports(value DeviceManagementReportsable)()
    SetResourceOperations(value []ResourceOperationable)()
    SetRoleAssignments(value []DeviceAndAppManagementRoleAssignmentable)()
    SetRoleDefinitions(value []RoleDefinitionable)()
    SetSettings(value DeviceManagementSettingsable)()
    SetSoftwareUpdateStatusSummary(value SoftwareUpdateStatusSummaryable)()
    SetSubscriptionState(value *DeviceManagementSubscriptionState)()
    SetTelecomExpenseManagementPartners(value []TelecomExpenseManagementPartnerable)()
    SetTermsAndConditions(value []TermsAndConditionsable)()
    SetTroubleshootingEvents(value []DeviceManagementTroubleshootingEventable)()
    SetWindowsAutopilotDeviceIdentities(value []WindowsAutopilotDeviceIdentityable)()
    SetWindowsInformationProtectionAppLearningSummaries(value []WindowsInformationProtectionAppLearningSummaryable)()
    SetWindowsInformationProtectionNetworkLearningSummaries(value []WindowsInformationProtectionNetworkLearningSummaryable)()
}
