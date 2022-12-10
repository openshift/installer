package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtection 
type WindowsInformationProtection struct {
    ManagedAppPolicy
    // Navigation property to list of security groups targeted for policy.
    assignments []TargetedManagedAppPolicyAssignmentable
    // Specifies whether to allow Azure RMS encryption for WIP
    azureRightsManagementServicesAllowed *bool
    // Specifies a recovery certificate that can be used for data recovery of encrypted files. This is the same as the data recovery agent(DRA) certificate for encrypting file system(EFS)
    dataRecoveryCertificate WindowsInformationProtectionDataRecoveryCertificateable
    // Possible values for WIP Protection enforcement levels
    enforcementLevel *WindowsInformationProtectionEnforcementLevel
    // Primary enterprise domain
    enterpriseDomain *string
    // This is the comma-separated list of internal proxy servers. For example, '157.54.14.28, 157.54.11.118, 10.202.14.167, 157.53.14.163, 157.69.210.59'. These proxies have been configured by the admin to connect to specific resources on the Internet. They are considered to be enterprise network locations. The proxies are only leveraged in configuring the EnterpriseProxiedDomains policy to force traffic to the matched domains through these proxies
    enterpriseInternalProxyServers []WindowsInformationProtectionResourceCollectionable
    // Sets the enterprise IP ranges that define the computers in the enterprise network. Data that comes from those computers will be considered part of the enterprise and protected. These locations will be considered a safe destination for enterprise data to be shared to
    enterpriseIPRanges []WindowsInformationProtectionIPRangeCollectionable
    // Boolean value that tells the client to accept the configured list and not to use heuristics to attempt to find other subnets. Default is false
    enterpriseIPRangesAreAuthoritative *bool
    // This is the list of domains that comprise the boundaries of the enterprise. Data from one of these domains that is sent to a device will be considered enterprise data and protected These locations will be considered a safe destination for enterprise data to be shared to
    enterpriseNetworkDomainNames []WindowsInformationProtectionResourceCollectionable
    // List of enterprise domains to be protected
    enterpriseProtectedDomainNames []WindowsInformationProtectionResourceCollectionable
    // Contains a list of Enterprise resource domains hosted in the cloud that need to be protected. Connections to these resources are considered enterprise data. If a proxy is paired with a cloud resource, traffic to the cloud resource will be routed through the enterprise network via the denoted proxy server (on Port 80). A proxy server used for this purpose must also be configured using the EnterpriseInternalProxyServers policy
    enterpriseProxiedDomains []WindowsInformationProtectionProxiedDomainCollectionable
    // This is a list of proxy servers. Any server not on this list is considered non-enterprise
    enterpriseProxyServers []WindowsInformationProtectionResourceCollectionable
    // Boolean value that tells the client to accept the configured list of proxies and not try to detect other work proxies. Default is false
    enterpriseProxyServersAreAuthoritative *bool
    // Another way to input exempt apps through xml files
    exemptAppLockerFiles []WindowsInformationProtectionAppLockerFileable
    // Exempt applications can also access enterprise data, but the data handled by those applications are not protected. This is because some critical enterprise applications may have compatibility problems with encrypted data.
    exemptApps []WindowsInformationProtectionAppable
    // Determines whether overlays are added to icons for WIP protected files in Explorer and enterprise only app tiles in the Start menu. Starting in Windows 10, version 1703 this setting also configures the visibility of the WIP icon in the title bar of a WIP-protected app
    iconsVisible *bool
    // This switch is for the Windows Search Indexer, to allow or disallow indexing of items
    indexingEncryptedStoresOrItemsBlocked *bool
    // Indicates if the policy is deployed to any inclusion groups or not.
    isAssigned *bool
    // List of domain names that can used for work or personal resource
    neutralDomainResources []WindowsInformationProtectionResourceCollectionable
    // Another way to input protected apps through xml files
    protectedAppLockerFiles []WindowsInformationProtectionAppLockerFileable
    // Protected applications can access enterprise data and the data handled by those applications are protected with encryption
    protectedApps []WindowsInformationProtectionAppable
    // Specifies whether the protection under lock feature (also known as encrypt under pin) should be configured
    protectionUnderLockConfigRequired *bool
    // This policy controls whether to revoke the WIP keys when a device unenrolls from the management service. If set to 1 (Don't revoke keys), the keys will not be revoked and the user will continue to have access to protected files after unenrollment. If the keys are not revoked, there will be no revoked file cleanup subsequently.
    revokeOnUnenrollDisabled *bool
    // TemplateID GUID to use for RMS encryption. The RMS template allows the IT admin to configure the details about who has access to RMS-protected file and how long they have access
    rightsManagementServicesTemplateId *string
    // Specifies a list of file extensions, so that files with these extensions are encrypted when copying from an SMB share within the corporate boundary
    smbAutoEncryptedFileExtensions []WindowsInformationProtectionResourceCollectionable
}
// NewWindowsInformationProtection instantiates a new WindowsInformationProtection and sets the default values.
func NewWindowsInformationProtection()(*WindowsInformationProtection) {
    m := &WindowsInformationProtection{
        ManagedAppPolicy: *NewManagedAppPolicy(),
    }
    odataTypeValue := "#microsoft.graph.windowsInformationProtection";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateWindowsInformationProtectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsInformationProtectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.mdmWindowsInformationProtectionPolicy":
                        return NewMdmWindowsInformationProtectionPolicy(), nil
                    case "#microsoft.graph.windowsInformationProtectionPolicy":
                        return NewWindowsInformationProtectionPolicy(), nil
                }
            }
        }
    }
    return NewWindowsInformationProtection(), nil
}
// GetAssignments gets the assignments property value. Navigation property to list of security groups targeted for policy.
func (m *WindowsInformationProtection) GetAssignments()([]TargetedManagedAppPolicyAssignmentable) {
    return m.assignments
}
// GetAzureRightsManagementServicesAllowed gets the azureRightsManagementServicesAllowed property value. Specifies whether to allow Azure RMS encryption for WIP
func (m *WindowsInformationProtection) GetAzureRightsManagementServicesAllowed()(*bool) {
    return m.azureRightsManagementServicesAllowed
}
// GetDataRecoveryCertificate gets the dataRecoveryCertificate property value. Specifies a recovery certificate that can be used for data recovery of encrypted files. This is the same as the data recovery agent(DRA) certificate for encrypting file system(EFS)
func (m *WindowsInformationProtection) GetDataRecoveryCertificate()(WindowsInformationProtectionDataRecoveryCertificateable) {
    return m.dataRecoveryCertificate
}
// GetEnforcementLevel gets the enforcementLevel property value. Possible values for WIP Protection enforcement levels
func (m *WindowsInformationProtection) GetEnforcementLevel()(*WindowsInformationProtectionEnforcementLevel) {
    return m.enforcementLevel
}
// GetEnterpriseDomain gets the enterpriseDomain property value. Primary enterprise domain
func (m *WindowsInformationProtection) GetEnterpriseDomain()(*string) {
    return m.enterpriseDomain
}
// GetEnterpriseInternalProxyServers gets the enterpriseInternalProxyServers property value. This is the comma-separated list of internal proxy servers. For example, '157.54.14.28, 157.54.11.118, 10.202.14.167, 157.53.14.163, 157.69.210.59'. These proxies have been configured by the admin to connect to specific resources on the Internet. They are considered to be enterprise network locations. The proxies are only leveraged in configuring the EnterpriseProxiedDomains policy to force traffic to the matched domains through these proxies
func (m *WindowsInformationProtection) GetEnterpriseInternalProxyServers()([]WindowsInformationProtectionResourceCollectionable) {
    return m.enterpriseInternalProxyServers
}
// GetEnterpriseIPRanges gets the enterpriseIPRanges property value. Sets the enterprise IP ranges that define the computers in the enterprise network. Data that comes from those computers will be considered part of the enterprise and protected. These locations will be considered a safe destination for enterprise data to be shared to
func (m *WindowsInformationProtection) GetEnterpriseIPRanges()([]WindowsInformationProtectionIPRangeCollectionable) {
    return m.enterpriseIPRanges
}
// GetEnterpriseIPRangesAreAuthoritative gets the enterpriseIPRangesAreAuthoritative property value. Boolean value that tells the client to accept the configured list and not to use heuristics to attempt to find other subnets. Default is false
func (m *WindowsInformationProtection) GetEnterpriseIPRangesAreAuthoritative()(*bool) {
    return m.enterpriseIPRangesAreAuthoritative
}
// GetEnterpriseNetworkDomainNames gets the enterpriseNetworkDomainNames property value. This is the list of domains that comprise the boundaries of the enterprise. Data from one of these domains that is sent to a device will be considered enterprise data and protected These locations will be considered a safe destination for enterprise data to be shared to
func (m *WindowsInformationProtection) GetEnterpriseNetworkDomainNames()([]WindowsInformationProtectionResourceCollectionable) {
    return m.enterpriseNetworkDomainNames
}
// GetEnterpriseProtectedDomainNames gets the enterpriseProtectedDomainNames property value. List of enterprise domains to be protected
func (m *WindowsInformationProtection) GetEnterpriseProtectedDomainNames()([]WindowsInformationProtectionResourceCollectionable) {
    return m.enterpriseProtectedDomainNames
}
// GetEnterpriseProxiedDomains gets the enterpriseProxiedDomains property value. Contains a list of Enterprise resource domains hosted in the cloud that need to be protected. Connections to these resources are considered enterprise data. If a proxy is paired with a cloud resource, traffic to the cloud resource will be routed through the enterprise network via the denoted proxy server (on Port 80). A proxy server used for this purpose must also be configured using the EnterpriseInternalProxyServers policy
func (m *WindowsInformationProtection) GetEnterpriseProxiedDomains()([]WindowsInformationProtectionProxiedDomainCollectionable) {
    return m.enterpriseProxiedDomains
}
// GetEnterpriseProxyServers gets the enterpriseProxyServers property value. This is a list of proxy servers. Any server not on this list is considered non-enterprise
func (m *WindowsInformationProtection) GetEnterpriseProxyServers()([]WindowsInformationProtectionResourceCollectionable) {
    return m.enterpriseProxyServers
}
// GetEnterpriseProxyServersAreAuthoritative gets the enterpriseProxyServersAreAuthoritative property value. Boolean value that tells the client to accept the configured list of proxies and not try to detect other work proxies. Default is false
func (m *WindowsInformationProtection) GetEnterpriseProxyServersAreAuthoritative()(*bool) {
    return m.enterpriseProxyServersAreAuthoritative
}
// GetExemptAppLockerFiles gets the exemptAppLockerFiles property value. Another way to input exempt apps through xml files
func (m *WindowsInformationProtection) GetExemptAppLockerFiles()([]WindowsInformationProtectionAppLockerFileable) {
    return m.exemptAppLockerFiles
}
// GetExemptApps gets the exemptApps property value. Exempt applications can also access enterprise data, but the data handled by those applications are not protected. This is because some critical enterprise applications may have compatibility problems with encrypted data.
func (m *WindowsInformationProtection) GetExemptApps()([]WindowsInformationProtectionAppable) {
    return m.exemptApps
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsInformationProtection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedAppPolicy.GetFieldDeserializers()
    res["assignments"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTargetedManagedAppPolicyAssignmentFromDiscriminatorValue , m.SetAssignments)
    res["azureRightsManagementServicesAllowed"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetAzureRightsManagementServicesAllowed)
    res["dataRecoveryCertificate"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateWindowsInformationProtectionDataRecoveryCertificateFromDiscriminatorValue , m.SetDataRecoveryCertificate)
    res["enforcementLevel"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseWindowsInformationProtectionEnforcementLevel , m.SetEnforcementLevel)
    res["enterpriseDomain"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEnterpriseDomain)
    res["enterpriseInternalProxyServers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue , m.SetEnterpriseInternalProxyServers)
    res["enterpriseIPRanges"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionIPRangeCollectionFromDiscriminatorValue , m.SetEnterpriseIPRanges)
    res["enterpriseIPRangesAreAuthoritative"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEnterpriseIPRangesAreAuthoritative)
    res["enterpriseNetworkDomainNames"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue , m.SetEnterpriseNetworkDomainNames)
    res["enterpriseProtectedDomainNames"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue , m.SetEnterpriseProtectedDomainNames)
    res["enterpriseProxiedDomains"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionProxiedDomainCollectionFromDiscriminatorValue , m.SetEnterpriseProxiedDomains)
    res["enterpriseProxyServers"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue , m.SetEnterpriseProxyServers)
    res["enterpriseProxyServersAreAuthoritative"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetEnterpriseProxyServersAreAuthoritative)
    res["exemptAppLockerFiles"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionAppLockerFileFromDiscriminatorValue , m.SetExemptAppLockerFiles)
    res["exemptApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionAppFromDiscriminatorValue , m.SetExemptApps)
    res["iconsVisible"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIconsVisible)
    res["indexingEncryptedStoresOrItemsBlocked"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIndexingEncryptedStoresOrItemsBlocked)
    res["isAssigned"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsAssigned)
    res["neutralDomainResources"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue , m.SetNeutralDomainResources)
    res["protectedAppLockerFiles"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionAppLockerFileFromDiscriminatorValue , m.SetProtectedAppLockerFiles)
    res["protectedApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionAppFromDiscriminatorValue , m.SetProtectedApps)
    res["protectionUnderLockConfigRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetProtectionUnderLockConfigRequired)
    res["revokeOnUnenrollDisabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetRevokeOnUnenrollDisabled)
    res["rightsManagementServicesTemplateId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRightsManagementServicesTemplateId)
    res["smbAutoEncryptedFileExtensions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateWindowsInformationProtectionResourceCollectionFromDiscriminatorValue , m.SetSmbAutoEncryptedFileExtensions)
    return res
}
// GetIconsVisible gets the iconsVisible property value. Determines whether overlays are added to icons for WIP protected files in Explorer and enterprise only app tiles in the Start menu. Starting in Windows 10, version 1703 this setting also configures the visibility of the WIP icon in the title bar of a WIP-protected app
func (m *WindowsInformationProtection) GetIconsVisible()(*bool) {
    return m.iconsVisible
}
// GetIndexingEncryptedStoresOrItemsBlocked gets the indexingEncryptedStoresOrItemsBlocked property value. This switch is for the Windows Search Indexer, to allow or disallow indexing of items
func (m *WindowsInformationProtection) GetIndexingEncryptedStoresOrItemsBlocked()(*bool) {
    return m.indexingEncryptedStoresOrItemsBlocked
}
// GetIsAssigned gets the isAssigned property value. Indicates if the policy is deployed to any inclusion groups or not.
func (m *WindowsInformationProtection) GetIsAssigned()(*bool) {
    return m.isAssigned
}
// GetNeutralDomainResources gets the neutralDomainResources property value. List of domain names that can used for work or personal resource
func (m *WindowsInformationProtection) GetNeutralDomainResources()([]WindowsInformationProtectionResourceCollectionable) {
    return m.neutralDomainResources
}
// GetProtectedAppLockerFiles gets the protectedAppLockerFiles property value. Another way to input protected apps through xml files
func (m *WindowsInformationProtection) GetProtectedAppLockerFiles()([]WindowsInformationProtectionAppLockerFileable) {
    return m.protectedAppLockerFiles
}
// GetProtectedApps gets the protectedApps property value. Protected applications can access enterprise data and the data handled by those applications are protected with encryption
func (m *WindowsInformationProtection) GetProtectedApps()([]WindowsInformationProtectionAppable) {
    return m.protectedApps
}
// GetProtectionUnderLockConfigRequired gets the protectionUnderLockConfigRequired property value. Specifies whether the protection under lock feature (also known as encrypt under pin) should be configured
func (m *WindowsInformationProtection) GetProtectionUnderLockConfigRequired()(*bool) {
    return m.protectionUnderLockConfigRequired
}
// GetRevokeOnUnenrollDisabled gets the revokeOnUnenrollDisabled property value. This policy controls whether to revoke the WIP keys when a device unenrolls from the management service. If set to 1 (Don't revoke keys), the keys will not be revoked and the user will continue to have access to protected files after unenrollment. If the keys are not revoked, there will be no revoked file cleanup subsequently.
func (m *WindowsInformationProtection) GetRevokeOnUnenrollDisabled()(*bool) {
    return m.revokeOnUnenrollDisabled
}
// GetRightsManagementServicesTemplateId gets the rightsManagementServicesTemplateId property value. TemplateID GUID to use for RMS encryption. The RMS template allows the IT admin to configure the details about who has access to RMS-protected file and how long they have access
func (m *WindowsInformationProtection) GetRightsManagementServicesTemplateId()(*string) {
    return m.rightsManagementServicesTemplateId
}
// GetSmbAutoEncryptedFileExtensions gets the smbAutoEncryptedFileExtensions property value. Specifies a list of file extensions, so that files with these extensions are encrypted when copying from an SMB share within the corporate boundary
func (m *WindowsInformationProtection) GetSmbAutoEncryptedFileExtensions()([]WindowsInformationProtectionResourceCollectionable) {
    return m.smbAutoEncryptedFileExtensions
}
// Serialize serializes information the current object
func (m *WindowsInformationProtection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedAppPolicy.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignments() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAssignments())
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("azureRightsManagementServicesAllowed", m.GetAzureRightsManagementServicesAllowed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("dataRecoveryCertificate", m.GetDataRecoveryCertificate())
        if err != nil {
            return err
        }
    }
    if m.GetEnforcementLevel() != nil {
        cast := (*m.GetEnforcementLevel()).String()
        err = writer.WriteStringValue("enforcementLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("enterpriseDomain", m.GetEnterpriseDomain())
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseInternalProxyServers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEnterpriseInternalProxyServers())
        err = writer.WriteCollectionOfObjectValues("enterpriseInternalProxyServers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseIPRanges() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEnterpriseIPRanges())
        err = writer.WriteCollectionOfObjectValues("enterpriseIPRanges", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enterpriseIPRangesAreAuthoritative", m.GetEnterpriseIPRangesAreAuthoritative())
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseNetworkDomainNames() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEnterpriseNetworkDomainNames())
        err = writer.WriteCollectionOfObjectValues("enterpriseNetworkDomainNames", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseProtectedDomainNames() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEnterpriseProtectedDomainNames())
        err = writer.WriteCollectionOfObjectValues("enterpriseProtectedDomainNames", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseProxiedDomains() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEnterpriseProxiedDomains())
        err = writer.WriteCollectionOfObjectValues("enterpriseProxiedDomains", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEnterpriseProxyServers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetEnterpriseProxyServers())
        err = writer.WriteCollectionOfObjectValues("enterpriseProxyServers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enterpriseProxyServersAreAuthoritative", m.GetEnterpriseProxyServersAreAuthoritative())
        if err != nil {
            return err
        }
    }
    if m.GetExemptAppLockerFiles() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetExemptAppLockerFiles())
        err = writer.WriteCollectionOfObjectValues("exemptAppLockerFiles", cast)
        if err != nil {
            return err
        }
    }
    if m.GetExemptApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetExemptApps())
        err = writer.WriteCollectionOfObjectValues("exemptApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iconsVisible", m.GetIconsVisible())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("indexingEncryptedStoresOrItemsBlocked", m.GetIndexingEncryptedStoresOrItemsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAssigned", m.GetIsAssigned())
        if err != nil {
            return err
        }
    }
    if m.GetNeutralDomainResources() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetNeutralDomainResources())
        err = writer.WriteCollectionOfObjectValues("neutralDomainResources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProtectedAppLockerFiles() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetProtectedAppLockerFiles())
        err = writer.WriteCollectionOfObjectValues("protectedAppLockerFiles", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProtectedApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetProtectedApps())
        err = writer.WriteCollectionOfObjectValues("protectedApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("protectionUnderLockConfigRequired", m.GetProtectionUnderLockConfigRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("revokeOnUnenrollDisabled", m.GetRevokeOnUnenrollDisabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("rightsManagementServicesTemplateId", m.GetRightsManagementServicesTemplateId())
        if err != nil {
            return err
        }
    }
    if m.GetSmbAutoEncryptedFileExtensions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetSmbAutoEncryptedFileExtensions())
        err = writer.WriteCollectionOfObjectValues("smbAutoEncryptedFileExtensions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignments sets the assignments property value. Navigation property to list of security groups targeted for policy.
func (m *WindowsInformationProtection) SetAssignments(value []TargetedManagedAppPolicyAssignmentable)() {
    m.assignments = value
}
// SetAzureRightsManagementServicesAllowed sets the azureRightsManagementServicesAllowed property value. Specifies whether to allow Azure RMS encryption for WIP
func (m *WindowsInformationProtection) SetAzureRightsManagementServicesAllowed(value *bool)() {
    m.azureRightsManagementServicesAllowed = value
}
// SetDataRecoveryCertificate sets the dataRecoveryCertificate property value. Specifies a recovery certificate that can be used for data recovery of encrypted files. This is the same as the data recovery agent(DRA) certificate for encrypting file system(EFS)
func (m *WindowsInformationProtection) SetDataRecoveryCertificate(value WindowsInformationProtectionDataRecoveryCertificateable)() {
    m.dataRecoveryCertificate = value
}
// SetEnforcementLevel sets the enforcementLevel property value. Possible values for WIP Protection enforcement levels
func (m *WindowsInformationProtection) SetEnforcementLevel(value *WindowsInformationProtectionEnforcementLevel)() {
    m.enforcementLevel = value
}
// SetEnterpriseDomain sets the enterpriseDomain property value. Primary enterprise domain
func (m *WindowsInformationProtection) SetEnterpriseDomain(value *string)() {
    m.enterpriseDomain = value
}
// SetEnterpriseInternalProxyServers sets the enterpriseInternalProxyServers property value. This is the comma-separated list of internal proxy servers. For example, '157.54.14.28, 157.54.11.118, 10.202.14.167, 157.53.14.163, 157.69.210.59'. These proxies have been configured by the admin to connect to specific resources on the Internet. They are considered to be enterprise network locations. The proxies are only leveraged in configuring the EnterpriseProxiedDomains policy to force traffic to the matched domains through these proxies
func (m *WindowsInformationProtection) SetEnterpriseInternalProxyServers(value []WindowsInformationProtectionResourceCollectionable)() {
    m.enterpriseInternalProxyServers = value
}
// SetEnterpriseIPRanges sets the enterpriseIPRanges property value. Sets the enterprise IP ranges that define the computers in the enterprise network. Data that comes from those computers will be considered part of the enterprise and protected. These locations will be considered a safe destination for enterprise data to be shared to
func (m *WindowsInformationProtection) SetEnterpriseIPRanges(value []WindowsInformationProtectionIPRangeCollectionable)() {
    m.enterpriseIPRanges = value
}
// SetEnterpriseIPRangesAreAuthoritative sets the enterpriseIPRangesAreAuthoritative property value. Boolean value that tells the client to accept the configured list and not to use heuristics to attempt to find other subnets. Default is false
func (m *WindowsInformationProtection) SetEnterpriseIPRangesAreAuthoritative(value *bool)() {
    m.enterpriseIPRangesAreAuthoritative = value
}
// SetEnterpriseNetworkDomainNames sets the enterpriseNetworkDomainNames property value. This is the list of domains that comprise the boundaries of the enterprise. Data from one of these domains that is sent to a device will be considered enterprise data and protected These locations will be considered a safe destination for enterprise data to be shared to
func (m *WindowsInformationProtection) SetEnterpriseNetworkDomainNames(value []WindowsInformationProtectionResourceCollectionable)() {
    m.enterpriseNetworkDomainNames = value
}
// SetEnterpriseProtectedDomainNames sets the enterpriseProtectedDomainNames property value. List of enterprise domains to be protected
func (m *WindowsInformationProtection) SetEnterpriseProtectedDomainNames(value []WindowsInformationProtectionResourceCollectionable)() {
    m.enterpriseProtectedDomainNames = value
}
// SetEnterpriseProxiedDomains sets the enterpriseProxiedDomains property value. Contains a list of Enterprise resource domains hosted in the cloud that need to be protected. Connections to these resources are considered enterprise data. If a proxy is paired with a cloud resource, traffic to the cloud resource will be routed through the enterprise network via the denoted proxy server (on Port 80). A proxy server used for this purpose must also be configured using the EnterpriseInternalProxyServers policy
func (m *WindowsInformationProtection) SetEnterpriseProxiedDomains(value []WindowsInformationProtectionProxiedDomainCollectionable)() {
    m.enterpriseProxiedDomains = value
}
// SetEnterpriseProxyServers sets the enterpriseProxyServers property value. This is a list of proxy servers. Any server not on this list is considered non-enterprise
func (m *WindowsInformationProtection) SetEnterpriseProxyServers(value []WindowsInformationProtectionResourceCollectionable)() {
    m.enterpriseProxyServers = value
}
// SetEnterpriseProxyServersAreAuthoritative sets the enterpriseProxyServersAreAuthoritative property value. Boolean value that tells the client to accept the configured list of proxies and not try to detect other work proxies. Default is false
func (m *WindowsInformationProtection) SetEnterpriseProxyServersAreAuthoritative(value *bool)() {
    m.enterpriseProxyServersAreAuthoritative = value
}
// SetExemptAppLockerFiles sets the exemptAppLockerFiles property value. Another way to input exempt apps through xml files
func (m *WindowsInformationProtection) SetExemptAppLockerFiles(value []WindowsInformationProtectionAppLockerFileable)() {
    m.exemptAppLockerFiles = value
}
// SetExemptApps sets the exemptApps property value. Exempt applications can also access enterprise data, but the data handled by those applications are not protected. This is because some critical enterprise applications may have compatibility problems with encrypted data.
func (m *WindowsInformationProtection) SetExemptApps(value []WindowsInformationProtectionAppable)() {
    m.exemptApps = value
}
// SetIconsVisible sets the iconsVisible property value. Determines whether overlays are added to icons for WIP protected files in Explorer and enterprise only app tiles in the Start menu. Starting in Windows 10, version 1703 this setting also configures the visibility of the WIP icon in the title bar of a WIP-protected app
func (m *WindowsInformationProtection) SetIconsVisible(value *bool)() {
    m.iconsVisible = value
}
// SetIndexingEncryptedStoresOrItemsBlocked sets the indexingEncryptedStoresOrItemsBlocked property value. This switch is for the Windows Search Indexer, to allow or disallow indexing of items
func (m *WindowsInformationProtection) SetIndexingEncryptedStoresOrItemsBlocked(value *bool)() {
    m.indexingEncryptedStoresOrItemsBlocked = value
}
// SetIsAssigned sets the isAssigned property value. Indicates if the policy is deployed to any inclusion groups or not.
func (m *WindowsInformationProtection) SetIsAssigned(value *bool)() {
    m.isAssigned = value
}
// SetNeutralDomainResources sets the neutralDomainResources property value. List of domain names that can used for work or personal resource
func (m *WindowsInformationProtection) SetNeutralDomainResources(value []WindowsInformationProtectionResourceCollectionable)() {
    m.neutralDomainResources = value
}
// SetProtectedAppLockerFiles sets the protectedAppLockerFiles property value. Another way to input protected apps through xml files
func (m *WindowsInformationProtection) SetProtectedAppLockerFiles(value []WindowsInformationProtectionAppLockerFileable)() {
    m.protectedAppLockerFiles = value
}
// SetProtectedApps sets the protectedApps property value. Protected applications can access enterprise data and the data handled by those applications are protected with encryption
func (m *WindowsInformationProtection) SetProtectedApps(value []WindowsInformationProtectionAppable)() {
    m.protectedApps = value
}
// SetProtectionUnderLockConfigRequired sets the protectionUnderLockConfigRequired property value. Specifies whether the protection under lock feature (also known as encrypt under pin) should be configured
func (m *WindowsInformationProtection) SetProtectionUnderLockConfigRequired(value *bool)() {
    m.protectionUnderLockConfigRequired = value
}
// SetRevokeOnUnenrollDisabled sets the revokeOnUnenrollDisabled property value. This policy controls whether to revoke the WIP keys when a device unenrolls from the management service. If set to 1 (Don't revoke keys), the keys will not be revoked and the user will continue to have access to protected files after unenrollment. If the keys are not revoked, there will be no revoked file cleanup subsequently.
func (m *WindowsInformationProtection) SetRevokeOnUnenrollDisabled(value *bool)() {
    m.revokeOnUnenrollDisabled = value
}
// SetRightsManagementServicesTemplateId sets the rightsManagementServicesTemplateId property value. TemplateID GUID to use for RMS encryption. The RMS template allows the IT admin to configure the details about who has access to RMS-protected file and how long they have access
func (m *WindowsInformationProtection) SetRightsManagementServicesTemplateId(value *string)() {
    m.rightsManagementServicesTemplateId = value
}
// SetSmbAutoEncryptedFileExtensions sets the smbAutoEncryptedFileExtensions property value. Specifies a list of file extensions, so that files with these extensions are encrypted when copying from an SMB share within the corporate boundary
func (m *WindowsInformationProtection) SetSmbAutoEncryptedFileExtensions(value []WindowsInformationProtectionResourceCollectionable)() {
    m.smbAutoEncryptedFileExtensions = value
}
