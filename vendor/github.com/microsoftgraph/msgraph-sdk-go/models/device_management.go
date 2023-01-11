package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagement 
type DeviceManagement struct {
    Entity
    // Apple push notification certificate.
    applePushNotificationCertificate ApplePushNotificationCertificateable
    // The Audit Events
    auditEvents []AuditEventable
    // The list of Compliance Management Partners configured by the tenant.
    complianceManagementPartners []ComplianceManagementPartnerable
    // The Exchange on premises conditional access settings. On premises conditional access will require devices to be both enrolled and compliant for mail access
    conditionalAccessSettings OnPremisesConditionalAccessSettingsable
    // The list of detected apps associated with a device.
    detectedApps []DetectedAppable
    // The list of device categories with the tenant.
    deviceCategories []DeviceCategoryable
    // The device compliance policies.
    deviceCompliancePolicies []DeviceCompliancePolicyable
    // The device compliance state summary for this account.
    deviceCompliancePolicyDeviceStateSummary DeviceCompliancePolicyDeviceStateSummaryable
    // The summary states of compliance policy settings for this account.
    deviceCompliancePolicySettingStateSummaries []DeviceCompliancePolicySettingStateSummaryable
    // The device configuration device state summary for this account.
    deviceConfigurationDeviceStateSummaries DeviceConfigurationDeviceStateSummaryable
    // The device configurations.
    deviceConfigurations []DeviceConfigurationable
    // The list of device enrollment configurations
    deviceEnrollmentConfigurations []DeviceEnrollmentConfigurationable
    // The list of Device Management Partners configured by the tenant.
    deviceManagementPartners []DeviceManagementPartnerable
    // The list of Exchange Connectors configured by the tenant.
    exchangeConnectors []DeviceManagementExchangeConnectorable
    // Collection of imported Windows autopilot devices.
    importedWindowsAutopilotDeviceIdentities []ImportedWindowsAutopilotDeviceIdentityable
    // Intune Account Id for given tenant
    intuneAccountId *string
    // intuneBrand contains data which is used in customizing the appearance of the Company Portal applications as well as the end user web portal.
    intuneBrand IntuneBrandable
    // The IOS software update installation statuses for this account.
    iosUpdateStatuses []IosUpdateDeviceStatusable
    // Device overview
    managedDeviceOverview ManagedDeviceOverviewable
    // The list of managed devices.
    managedDevices []ManagedDeviceable
    // The list of Mobile threat Defense connectors configured by the tenant.
    mobileThreatDefenseConnectors []MobileThreatDefenseConnectorable
    // The Notification Message Templates.
    notificationMessageTemplates []NotificationMessageTemplateable
    // The remote assist partners.
    remoteAssistancePartners []RemoteAssistancePartnerable
    // Reports singleton
    reports DeviceManagementReportsable
    // The Resource Operations.
    resourceOperations []ResourceOperationable
    // The Role Assignments.
    roleAssignments []DeviceAndAppManagementRoleAssignmentable
    // The Role Definitions.
    roleDefinitions []RoleDefinitionable
    // Account level settings.
    settings DeviceManagementSettingsable
    // The software update status summary.
    softwareUpdateStatusSummary SoftwareUpdateStatusSummaryable
    // Tenant mobile device management subscription state.
    subscriptionState *DeviceManagementSubscriptionState
    // The telecom expense management partners.
    telecomExpenseManagementPartners []TelecomExpenseManagementPartnerable
    // The terms and conditions associated with device management of the company.
    termsAndConditions []TermsAndConditionsable
    // The list of troubleshooting events for the tenant.
    troubleshootingEvents []DeviceManagementTroubleshootingEventable
    // The Windows autopilot device identities contained collection.
    windowsAutopilotDeviceIdentities []WindowsAutopilotDeviceIdentityable
    // The windows information protection app learning summaries.
    windowsInformationProtectionAppLearningSummaries []WindowsInformationProtectionAppLearningSummaryable
    // The windows information protection network learning summaries.
    windowsInformationProtectionNetworkLearningSummaries []WindowsInformationProtectionNetworkLearningSummaryable
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
    return m.applePushNotificationCertificate
}
// GetAuditEvents gets the auditEvents property value. The Audit Events
func (m *DeviceManagement) GetAuditEvents()([]AuditEventable) {
    return m.auditEvents
}
// GetComplianceManagementPartners gets the complianceManagementPartners property value. The list of Compliance Management Partners configured by the tenant.
func (m *DeviceManagement) GetComplianceManagementPartners()([]ComplianceManagementPartnerable) {
    return m.complianceManagementPartners
}
// GetConditionalAccessSettings gets the conditionalAccessSettings property value. The Exchange on premises conditional access settings. On premises conditional access will require devices to be both enrolled and compliant for mail access
func (m *DeviceManagement) GetConditionalAccessSettings()(OnPremisesConditionalAccessSettingsable) {
    return m.conditionalAccessSettings
}
// GetDetectedApps gets the detectedApps property value. The list of detected apps associated with a device.
func (m *DeviceManagement) GetDetectedApps()([]DetectedAppable) {
    return m.detectedApps
}
// GetDeviceCategories gets the deviceCategories property value. The list of device categories with the tenant.
func (m *DeviceManagement) GetDeviceCategories()([]DeviceCategoryable) {
    return m.deviceCategories
}
// GetDeviceCompliancePolicies gets the deviceCompliancePolicies property value. The device compliance policies.
func (m *DeviceManagement) GetDeviceCompliancePolicies()([]DeviceCompliancePolicyable) {
    return m.deviceCompliancePolicies
}
// GetDeviceCompliancePolicyDeviceStateSummary gets the deviceCompliancePolicyDeviceStateSummary property value. The device compliance state summary for this account.
func (m *DeviceManagement) GetDeviceCompliancePolicyDeviceStateSummary()(DeviceCompliancePolicyDeviceStateSummaryable) {
    return m.deviceCompliancePolicyDeviceStateSummary
}
// GetDeviceCompliancePolicySettingStateSummaries gets the deviceCompliancePolicySettingStateSummaries property value. The summary states of compliance policy settings for this account.
func (m *DeviceManagement) GetDeviceCompliancePolicySettingStateSummaries()([]DeviceCompliancePolicySettingStateSummaryable) {
    return m.deviceCompliancePolicySettingStateSummaries
}
// GetDeviceConfigurationDeviceStateSummaries gets the deviceConfigurationDeviceStateSummaries property value. The device configuration device state summary for this account.
func (m *DeviceManagement) GetDeviceConfigurationDeviceStateSummaries()(DeviceConfigurationDeviceStateSummaryable) {
    return m.deviceConfigurationDeviceStateSummaries
}
// GetDeviceConfigurations gets the deviceConfigurations property value. The device configurations.
func (m *DeviceManagement) GetDeviceConfigurations()([]DeviceConfigurationable) {
    return m.deviceConfigurations
}
// GetDeviceEnrollmentConfigurations gets the deviceEnrollmentConfigurations property value. The list of device enrollment configurations
func (m *DeviceManagement) GetDeviceEnrollmentConfigurations()([]DeviceEnrollmentConfigurationable) {
    return m.deviceEnrollmentConfigurations
}
// GetDeviceManagementPartners gets the deviceManagementPartners property value. The list of Device Management Partners configured by the tenant.
func (m *DeviceManagement) GetDeviceManagementPartners()([]DeviceManagementPartnerable) {
    return m.deviceManagementPartners
}
// GetExchangeConnectors gets the exchangeConnectors property value. The list of Exchange Connectors configured by the tenant.
func (m *DeviceManagement) GetExchangeConnectors()([]DeviceManagementExchangeConnectorable) {
    return m.exchangeConnectors
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceManagement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["applePushNotificationCertificate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateApplePushNotificationCertificateFromDiscriminatorValue , m.SetApplePushNotificationCertificate)
    res["auditEvents"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAuditEventFromDiscriminatorValue , m.SetAuditEvents)
    res["complianceManagementPartners"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateComplianceManagementPartnerFromDiscriminatorValue , m.SetComplianceManagementPartners)
    res["conditionalAccessSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateOnPremisesConditionalAccessSettingsFromDiscriminatorValue , m.SetConditionalAccessSettings)
    res["detectedApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDetectedAppFromDiscriminatorValue , m.SetDetectedApps)
    res["deviceCategories"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceCategoryFromDiscriminatorValue , m.SetDeviceCategories)
    res["deviceCompliancePolicies"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceCompliancePolicyFromDiscriminatorValue , m.SetDeviceCompliancePolicies)
    res["deviceCompliancePolicyDeviceStateSummary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceCompliancePolicyDeviceStateSummaryFromDiscriminatorValue , m.SetDeviceCompliancePolicyDeviceStateSummary)
    res["deviceCompliancePolicySettingStateSummaries"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceCompliancePolicySettingStateSummaryFromDiscriminatorValue , m.SetDeviceCompliancePolicySettingStateSummaries)
    res["deviceConfigurationDeviceStateSummaries"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceConfigurationDeviceStateSummaryFromDiscriminatorValue , m.SetDeviceConfigurationDeviceStateSummaries)
    res["deviceConfigurations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceConfigurationFromDiscriminatorValue , m.SetDeviceConfigurations)
    res["deviceEnrollmentConfigurations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceEnrollmentConfigurationFromDiscriminatorValue , m.SetDeviceEnrollmentConfigurations)
    res["deviceManagementPartners"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceManagementPartnerFromDiscriminatorValue , m.SetDeviceManagementPartners)
    res["exchangeConnectors"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceManagementExchangeConnectorFromDiscriminatorValue , m.SetExchangeConnectors)
    res["importedWindowsAutopilotDeviceIdentities"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateImportedWindowsAutopilotDeviceIdentityFromDiscriminatorValue , m.SetImportedWindowsAutopilotDeviceIdentities)
    res["intuneAccountId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetIntuneAccountId)
    res["intuneBrand"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIntuneBrandFromDiscriminatorValue , m.SetIntuneBrand)
    res["iosUpdateStatuses"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIosUpdateDeviceStatusFromDiscriminatorValue , m.SetIosUpdateStatuses)
    res["managedDeviceOverview"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateManagedDeviceOverviewFromDiscriminatorValue , m.SetManagedDeviceOverview)
    res["managedDevices"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateManagedDeviceFromDiscriminatorValue , m.SetManagedDevices)
    res["mobileThreatDefenseConnectors"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateMobileThreatDefenseConnectorFromDiscriminatorValue , m.SetMobileThreatDefenseConnectors)
    res["notificationMessageTemplates"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateNotificationMessageTemplateFromDiscriminatorValue , m.SetNotificationMessageTemplates)
    res["remoteAssistancePartners"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRemoteAssistancePartnerFromDiscriminatorValue , m.SetRemoteAssistancePartners)
    res["reports"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceManagementReportsFromDiscriminatorValue , m.SetReports)
    res["resourceOperations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateResourceOperationFromDiscriminatorValue , m.SetResourceOperations)
    res["roleAssignments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceAndAppManagementRoleAssignmentFromDiscriminatorValue , m.SetRoleAssignments)
    res["roleDefinitions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRoleDefinitionFromDiscriminatorValue , m.SetRoleDefinitions)
    res["settings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateDeviceManagementSettingsFromDiscriminatorValue , m.SetSettings)
    res["softwareUpdateStatusSummary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateSoftwareUpdateStatusSummaryFromDiscriminatorValue , m.SetSoftwareUpdateStatusSummary)
    res["subscriptionState"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceManagementSubscriptionState , m.SetSubscriptionState)
    res["telecomExpenseManagementPartners"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTelecomExpenseManagementPartnerFromDiscriminatorValue , m.SetTelecomExpenseManagementPartners)
    res["termsAndConditions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTermsAndConditionsFromDiscriminatorValue , m.SetTermsAndConditions)
    res["troubleshootingEvents"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceManagementTroubleshootingEventFromDiscriminatorValue , m.SetTroubleshootingEvents)
    res["windowsAutopilotDeviceIdentities"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsAutopilotDeviceIdentityFromDiscriminatorValue , m.SetWindowsAutopilotDeviceIdentities)
    res["windowsInformationProtectionAppLearningSummaries"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionAppLearningSummaryFromDiscriminatorValue , m.SetWindowsInformationProtectionAppLearningSummaries)
    res["windowsInformationProtectionNetworkLearningSummaries"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionNetworkLearningSummaryFromDiscriminatorValue , m.SetWindowsInformationProtectionNetworkLearningSummaries)
    return res
}
// GetImportedWindowsAutopilotDeviceIdentities gets the importedWindowsAutopilotDeviceIdentities property value. Collection of imported Windows autopilot devices.
func (m *DeviceManagement) GetImportedWindowsAutopilotDeviceIdentities()([]ImportedWindowsAutopilotDeviceIdentityable) {
    return m.importedWindowsAutopilotDeviceIdentities
}
// GetIntuneAccountId gets the intuneAccountId property value. Intune Account Id for given tenant
func (m *DeviceManagement) GetIntuneAccountId()(*string) {
    return m.intuneAccountId
}
// GetIntuneBrand gets the intuneBrand property value. intuneBrand contains data which is used in customizing the appearance of the Company Portal applications as well as the end user web portal.
func (m *DeviceManagement) GetIntuneBrand()(IntuneBrandable) {
    return m.intuneBrand
}
// GetIosUpdateStatuses gets the iosUpdateStatuses property value. The IOS software update installation statuses for this account.
func (m *DeviceManagement) GetIosUpdateStatuses()([]IosUpdateDeviceStatusable) {
    return m.iosUpdateStatuses
}
// GetManagedDeviceOverview gets the managedDeviceOverview property value. Device overview
func (m *DeviceManagement) GetManagedDeviceOverview()(ManagedDeviceOverviewable) {
    return m.managedDeviceOverview
}
// GetManagedDevices gets the managedDevices property value. The list of managed devices.
func (m *DeviceManagement) GetManagedDevices()([]ManagedDeviceable) {
    return m.managedDevices
}
// GetMobileThreatDefenseConnectors gets the mobileThreatDefenseConnectors property value. The list of Mobile threat Defense connectors configured by the tenant.
func (m *DeviceManagement) GetMobileThreatDefenseConnectors()([]MobileThreatDefenseConnectorable) {
    return m.mobileThreatDefenseConnectors
}
// GetNotificationMessageTemplates gets the notificationMessageTemplates property value. The Notification Message Templates.
func (m *DeviceManagement) GetNotificationMessageTemplates()([]NotificationMessageTemplateable) {
    return m.notificationMessageTemplates
}
// GetRemoteAssistancePartners gets the remoteAssistancePartners property value. The remote assist partners.
func (m *DeviceManagement) GetRemoteAssistancePartners()([]RemoteAssistancePartnerable) {
    return m.remoteAssistancePartners
}
// GetReports gets the reports property value. Reports singleton
func (m *DeviceManagement) GetReports()(DeviceManagementReportsable) {
    return m.reports
}
// GetResourceOperations gets the resourceOperations property value. The Resource Operations.
func (m *DeviceManagement) GetResourceOperations()([]ResourceOperationable) {
    return m.resourceOperations
}
// GetRoleAssignments gets the roleAssignments property value. The Role Assignments.
func (m *DeviceManagement) GetRoleAssignments()([]DeviceAndAppManagementRoleAssignmentable) {
    return m.roleAssignments
}
// GetRoleDefinitions gets the roleDefinitions property value. The Role Definitions.
func (m *DeviceManagement) GetRoleDefinitions()([]RoleDefinitionable) {
    return m.roleDefinitions
}
// GetSettings gets the settings property value. Account level settings.
func (m *DeviceManagement) GetSettings()(DeviceManagementSettingsable) {
    return m.settings
}
// GetSoftwareUpdateStatusSummary gets the softwareUpdateStatusSummary property value. The software update status summary.
func (m *DeviceManagement) GetSoftwareUpdateStatusSummary()(SoftwareUpdateStatusSummaryable) {
    return m.softwareUpdateStatusSummary
}
// GetSubscriptionState gets the subscriptionState property value. Tenant mobile device management subscription state.
func (m *DeviceManagement) GetSubscriptionState()(*DeviceManagementSubscriptionState) {
    return m.subscriptionState
}
// GetTelecomExpenseManagementPartners gets the telecomExpenseManagementPartners property value. The telecom expense management partners.
func (m *DeviceManagement) GetTelecomExpenseManagementPartners()([]TelecomExpenseManagementPartnerable) {
    return m.telecomExpenseManagementPartners
}
// GetTermsAndConditions gets the termsAndConditions property value. The terms and conditions associated with device management of the company.
func (m *DeviceManagement) GetTermsAndConditions()([]TermsAndConditionsable) {
    return m.termsAndConditions
}
// GetTroubleshootingEvents gets the troubleshootingEvents property value. The list of troubleshooting events for the tenant.
func (m *DeviceManagement) GetTroubleshootingEvents()([]DeviceManagementTroubleshootingEventable) {
    return m.troubleshootingEvents
}
// GetWindowsAutopilotDeviceIdentities gets the windowsAutopilotDeviceIdentities property value. The Windows autopilot device identities contained collection.
func (m *DeviceManagement) GetWindowsAutopilotDeviceIdentities()([]WindowsAutopilotDeviceIdentityable) {
    return m.windowsAutopilotDeviceIdentities
}
// GetWindowsInformationProtectionAppLearningSummaries gets the windowsInformationProtectionAppLearningSummaries property value. The windows information protection app learning summaries.
func (m *DeviceManagement) GetWindowsInformationProtectionAppLearningSummaries()([]WindowsInformationProtectionAppLearningSummaryable) {
    return m.windowsInformationProtectionAppLearningSummaries
}
// GetWindowsInformationProtectionNetworkLearningSummaries gets the windowsInformationProtectionNetworkLearningSummaries property value. The windows information protection network learning summaries.
func (m *DeviceManagement) GetWindowsInformationProtectionNetworkLearningSummaries()([]WindowsInformationProtectionNetworkLearningSummaryable) {
    return m.windowsInformationProtectionNetworkLearningSummaries
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAuditEvents())
        err = writer.WriteCollectionOfObjectValues("auditEvents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetComplianceManagementPartners() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetComplianceManagementPartners())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDetectedApps())
        err = writer.WriteCollectionOfObjectValues("detectedApps", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceCategories() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceCategories())
        err = writer.WriteCollectionOfObjectValues("deviceCategories", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceCompliancePolicies() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceCompliancePolicies())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceCompliancePolicySettingStateSummaries())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceConfigurations())
        err = writer.WriteCollectionOfObjectValues("deviceConfigurations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceEnrollmentConfigurations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceEnrollmentConfigurations())
        err = writer.WriteCollectionOfObjectValues("deviceEnrollmentConfigurations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceManagementPartners() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceManagementPartners())
        err = writer.WriteCollectionOfObjectValues("deviceManagementPartners", cast)
        if err != nil {
            return err
        }
    }
    if m.GetExchangeConnectors() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetExchangeConnectors())
        err = writer.WriteCollectionOfObjectValues("exchangeConnectors", cast)
        if err != nil {
            return err
        }
    }
    if m.GetImportedWindowsAutopilotDeviceIdentities() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetImportedWindowsAutopilotDeviceIdentities())
        err = writer.WriteCollectionOfObjectValues("importedWindowsAutopilotDeviceIdentities", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("intuneAccountId", m.GetIntuneAccountId())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetIosUpdateStatuses())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetManagedDevices())
        err = writer.WriteCollectionOfObjectValues("managedDevices", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMobileThreatDefenseConnectors() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMobileThreatDefenseConnectors())
        err = writer.WriteCollectionOfObjectValues("mobileThreatDefenseConnectors", cast)
        if err != nil {
            return err
        }
    }
    if m.GetNotificationMessageTemplates() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNotificationMessageTemplates())
        err = writer.WriteCollectionOfObjectValues("notificationMessageTemplates", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRemoteAssistancePartners() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetRemoteAssistancePartners())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetResourceOperations())
        err = writer.WriteCollectionOfObjectValues("resourceOperations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleAssignments() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetRoleAssignments())
        err = writer.WriteCollectionOfObjectValues("roleAssignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleDefinitions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetRoleDefinitions())
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
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTelecomExpenseManagementPartners())
        err = writer.WriteCollectionOfObjectValues("telecomExpenseManagementPartners", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTermsAndConditions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTermsAndConditions())
        err = writer.WriteCollectionOfObjectValues("termsAndConditions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTroubleshootingEvents() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTroubleshootingEvents())
        err = writer.WriteCollectionOfObjectValues("troubleshootingEvents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsAutopilotDeviceIdentities() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetWindowsAutopilotDeviceIdentities())
        err = writer.WriteCollectionOfObjectValues("windowsAutopilotDeviceIdentities", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsInformationProtectionAppLearningSummaries() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetWindowsInformationProtectionAppLearningSummaries())
        err = writer.WriteCollectionOfObjectValues("windowsInformationProtectionAppLearningSummaries", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsInformationProtectionNetworkLearningSummaries() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetWindowsInformationProtectionNetworkLearningSummaries())
        err = writer.WriteCollectionOfObjectValues("windowsInformationProtectionNetworkLearningSummaries", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplePushNotificationCertificate sets the applePushNotificationCertificate property value. Apple push notification certificate.
func (m *DeviceManagement) SetApplePushNotificationCertificate(value ApplePushNotificationCertificateable)() {
    m.applePushNotificationCertificate = value
}
// SetAuditEvents sets the auditEvents property value. The Audit Events
func (m *DeviceManagement) SetAuditEvents(value []AuditEventable)() {
    m.auditEvents = value
}
// SetComplianceManagementPartners sets the complianceManagementPartners property value. The list of Compliance Management Partners configured by the tenant.
func (m *DeviceManagement) SetComplianceManagementPartners(value []ComplianceManagementPartnerable)() {
    m.complianceManagementPartners = value
}
// SetConditionalAccessSettings sets the conditionalAccessSettings property value. The Exchange on premises conditional access settings. On premises conditional access will require devices to be both enrolled and compliant for mail access
func (m *DeviceManagement) SetConditionalAccessSettings(value OnPremisesConditionalAccessSettingsable)() {
    m.conditionalAccessSettings = value
}
// SetDetectedApps sets the detectedApps property value. The list of detected apps associated with a device.
func (m *DeviceManagement) SetDetectedApps(value []DetectedAppable)() {
    m.detectedApps = value
}
// SetDeviceCategories sets the deviceCategories property value. The list of device categories with the tenant.
func (m *DeviceManagement) SetDeviceCategories(value []DeviceCategoryable)() {
    m.deviceCategories = value
}
// SetDeviceCompliancePolicies sets the deviceCompliancePolicies property value. The device compliance policies.
func (m *DeviceManagement) SetDeviceCompliancePolicies(value []DeviceCompliancePolicyable)() {
    m.deviceCompliancePolicies = value
}
// SetDeviceCompliancePolicyDeviceStateSummary sets the deviceCompliancePolicyDeviceStateSummary property value. The device compliance state summary for this account.
func (m *DeviceManagement) SetDeviceCompliancePolicyDeviceStateSummary(value DeviceCompliancePolicyDeviceStateSummaryable)() {
    m.deviceCompliancePolicyDeviceStateSummary = value
}
// SetDeviceCompliancePolicySettingStateSummaries sets the deviceCompliancePolicySettingStateSummaries property value. The summary states of compliance policy settings for this account.
func (m *DeviceManagement) SetDeviceCompliancePolicySettingStateSummaries(value []DeviceCompliancePolicySettingStateSummaryable)() {
    m.deviceCompliancePolicySettingStateSummaries = value
}
// SetDeviceConfigurationDeviceStateSummaries sets the deviceConfigurationDeviceStateSummaries property value. The device configuration device state summary for this account.
func (m *DeviceManagement) SetDeviceConfigurationDeviceStateSummaries(value DeviceConfigurationDeviceStateSummaryable)() {
    m.deviceConfigurationDeviceStateSummaries = value
}
// SetDeviceConfigurations sets the deviceConfigurations property value. The device configurations.
func (m *DeviceManagement) SetDeviceConfigurations(value []DeviceConfigurationable)() {
    m.deviceConfigurations = value
}
// SetDeviceEnrollmentConfigurations sets the deviceEnrollmentConfigurations property value. The list of device enrollment configurations
func (m *DeviceManagement) SetDeviceEnrollmentConfigurations(value []DeviceEnrollmentConfigurationable)() {
    m.deviceEnrollmentConfigurations = value
}
// SetDeviceManagementPartners sets the deviceManagementPartners property value. The list of Device Management Partners configured by the tenant.
func (m *DeviceManagement) SetDeviceManagementPartners(value []DeviceManagementPartnerable)() {
    m.deviceManagementPartners = value
}
// SetExchangeConnectors sets the exchangeConnectors property value. The list of Exchange Connectors configured by the tenant.
func (m *DeviceManagement) SetExchangeConnectors(value []DeviceManagementExchangeConnectorable)() {
    m.exchangeConnectors = value
}
// SetImportedWindowsAutopilotDeviceIdentities sets the importedWindowsAutopilotDeviceIdentities property value. Collection of imported Windows autopilot devices.
func (m *DeviceManagement) SetImportedWindowsAutopilotDeviceIdentities(value []ImportedWindowsAutopilotDeviceIdentityable)() {
    m.importedWindowsAutopilotDeviceIdentities = value
}
// SetIntuneAccountId sets the intuneAccountId property value. Intune Account Id for given tenant
func (m *DeviceManagement) SetIntuneAccountId(value *string)() {
    m.intuneAccountId = value
}
// SetIntuneBrand sets the intuneBrand property value. intuneBrand contains data which is used in customizing the appearance of the Company Portal applications as well as the end user web portal.
func (m *DeviceManagement) SetIntuneBrand(value IntuneBrandable)() {
    m.intuneBrand = value
}
// SetIosUpdateStatuses sets the iosUpdateStatuses property value. The IOS software update installation statuses for this account.
func (m *DeviceManagement) SetIosUpdateStatuses(value []IosUpdateDeviceStatusable)() {
    m.iosUpdateStatuses = value
}
// SetManagedDeviceOverview sets the managedDeviceOverview property value. Device overview
func (m *DeviceManagement) SetManagedDeviceOverview(value ManagedDeviceOverviewable)() {
    m.managedDeviceOverview = value
}
// SetManagedDevices sets the managedDevices property value. The list of managed devices.
func (m *DeviceManagement) SetManagedDevices(value []ManagedDeviceable)() {
    m.managedDevices = value
}
// SetMobileThreatDefenseConnectors sets the mobileThreatDefenseConnectors property value. The list of Mobile threat Defense connectors configured by the tenant.
func (m *DeviceManagement) SetMobileThreatDefenseConnectors(value []MobileThreatDefenseConnectorable)() {
    m.mobileThreatDefenseConnectors = value
}
// SetNotificationMessageTemplates sets the notificationMessageTemplates property value. The Notification Message Templates.
func (m *DeviceManagement) SetNotificationMessageTemplates(value []NotificationMessageTemplateable)() {
    m.notificationMessageTemplates = value
}
// SetRemoteAssistancePartners sets the remoteAssistancePartners property value. The remote assist partners.
func (m *DeviceManagement) SetRemoteAssistancePartners(value []RemoteAssistancePartnerable)() {
    m.remoteAssistancePartners = value
}
// SetReports sets the reports property value. Reports singleton
func (m *DeviceManagement) SetReports(value DeviceManagementReportsable)() {
    m.reports = value
}
// SetResourceOperations sets the resourceOperations property value. The Resource Operations.
func (m *DeviceManagement) SetResourceOperations(value []ResourceOperationable)() {
    m.resourceOperations = value
}
// SetRoleAssignments sets the roleAssignments property value. The Role Assignments.
func (m *DeviceManagement) SetRoleAssignments(value []DeviceAndAppManagementRoleAssignmentable)() {
    m.roleAssignments = value
}
// SetRoleDefinitions sets the roleDefinitions property value. The Role Definitions.
func (m *DeviceManagement) SetRoleDefinitions(value []RoleDefinitionable)() {
    m.roleDefinitions = value
}
// SetSettings sets the settings property value. Account level settings.
func (m *DeviceManagement) SetSettings(value DeviceManagementSettingsable)() {
    m.settings = value
}
// SetSoftwareUpdateStatusSummary sets the softwareUpdateStatusSummary property value. The software update status summary.
func (m *DeviceManagement) SetSoftwareUpdateStatusSummary(value SoftwareUpdateStatusSummaryable)() {
    m.softwareUpdateStatusSummary = value
}
// SetSubscriptionState sets the subscriptionState property value. Tenant mobile device management subscription state.
func (m *DeviceManagement) SetSubscriptionState(value *DeviceManagementSubscriptionState)() {
    m.subscriptionState = value
}
// SetTelecomExpenseManagementPartners sets the telecomExpenseManagementPartners property value. The telecom expense management partners.
func (m *DeviceManagement) SetTelecomExpenseManagementPartners(value []TelecomExpenseManagementPartnerable)() {
    m.telecomExpenseManagementPartners = value
}
// SetTermsAndConditions sets the termsAndConditions property value. The terms and conditions associated with device management of the company.
func (m *DeviceManagement) SetTermsAndConditions(value []TermsAndConditionsable)() {
    m.termsAndConditions = value
}
// SetTroubleshootingEvents sets the troubleshootingEvents property value. The list of troubleshooting events for the tenant.
func (m *DeviceManagement) SetTroubleshootingEvents(value []DeviceManagementTroubleshootingEventable)() {
    m.troubleshootingEvents = value
}
// SetWindowsAutopilotDeviceIdentities sets the windowsAutopilotDeviceIdentities property value. The Windows autopilot device identities contained collection.
func (m *DeviceManagement) SetWindowsAutopilotDeviceIdentities(value []WindowsAutopilotDeviceIdentityable)() {
    m.windowsAutopilotDeviceIdentities = value
}
// SetWindowsInformationProtectionAppLearningSummaries sets the windowsInformationProtectionAppLearningSummaries property value. The windows information protection app learning summaries.
func (m *DeviceManagement) SetWindowsInformationProtectionAppLearningSummaries(value []WindowsInformationProtectionAppLearningSummaryable)() {
    m.windowsInformationProtectionAppLearningSummaries = value
}
// SetWindowsInformationProtectionNetworkLearningSummaries sets the windowsInformationProtectionNetworkLearningSummaries property value. The windows information protection network learning summaries.
func (m *DeviceManagement) SetWindowsInformationProtectionNetworkLearningSummaries(value []WindowsInformationProtectionNetworkLearningSummaryable)() {
    m.windowsInformationProtectionNetworkLearningSummaries = value
}
