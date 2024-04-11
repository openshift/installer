package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ServicePrincipal 
type ServicePrincipal struct {
    DirectoryObject
}
// NewServicePrincipal instantiates a new ServicePrincipal and sets the default values.
func NewServicePrincipal()(*ServicePrincipal) {
    m := &ServicePrincipal{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.servicePrincipal"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateServicePrincipalFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateServicePrincipalFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServicePrincipal(), nil
}
// GetAccountEnabled gets the accountEnabled property value. true if the service principal account is enabled; otherwise, false. If set to false, then no users will be able to sign in to this app, even if they are assigned to it. Supports $filter (eq, ne, not, in).
func (m *ServicePrincipal) GetAccountEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("accountEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAddIns gets the addIns property value. Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams may set the addIns property for its 'FileHandler' functionality. This will let services like Microsoft 365 call the application in the context of a document the user is working on.
func (m *ServicePrincipal) GetAddIns()([]AddInable) {
    val, err := m.GetBackingStore().Get("addIns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AddInable)
    }
    return nil
}
// GetAlternativeNames gets the alternativeNames property value. Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities. Supports $filter (eq, not, ge, le, startsWith).
func (m *ServicePrincipal) GetAlternativeNames()([]string) {
    val, err := m.GetBackingStore().Get("alternativeNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetAppDescription gets the appDescription property value. The description exposed by the associated application.
func (m *ServicePrincipal) GetAppDescription()(*string) {
    val, err := m.GetBackingStore().Get("appDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppDisplayName gets the appDisplayName property value. The display name exposed by the associated application.
func (m *ServicePrincipal) GetAppDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("appDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppId gets the appId property value. The unique identifier for the associated application (its appId property). Supports $filter (eq, ne, not, in, startsWith).
func (m *ServicePrincipal) GetAppId()(*string) {
    val, err := m.GetBackingStore().Get("appId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetApplicationTemplateId gets the applicationTemplateId property value. Unique identifier of the applicationTemplate that the servicePrincipal was created from. Read-only. Supports $filter (eq, ne, NOT, startsWith).
func (m *ServicePrincipal) GetApplicationTemplateId()(*string) {
    val, err := m.GetBackingStore().Get("applicationTemplateId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppManagementPolicies gets the appManagementPolicies property value. The appManagementPolicies property
func (m *ServicePrincipal) GetAppManagementPolicies()([]AppManagementPolicyable) {
    val, err := m.GetBackingStore().Get("appManagementPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppManagementPolicyable)
    }
    return nil
}
// GetAppOwnerOrganizationId gets the appOwnerOrganizationId property value. Contains the tenant id where the application is registered. This is applicable only to service principals backed by applications. Supports $filter (eq, ne, NOT, ge, le).
func (m *ServicePrincipal) GetAppOwnerOrganizationId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("appOwnerOrganizationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetAppRoleAssignedTo gets the appRoleAssignedTo property value. App role assignments for this app or service, granted to users, groups, and other service principals. Supports $expand.
func (m *ServicePrincipal) GetAppRoleAssignedTo()([]AppRoleAssignmentable) {
    val, err := m.GetBackingStore().Get("appRoleAssignedTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppRoleAssignmentable)
    }
    return nil
}
// GetAppRoleAssignmentRequired gets the appRoleAssignmentRequired property value. Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports $filter (eq, ne, NOT).
func (m *ServicePrincipal) GetAppRoleAssignmentRequired()(*bool) {
    val, err := m.GetBackingStore().Get("appRoleAssignmentRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppRoleAssignments gets the appRoleAssignments property value. App role assignment for another app or service, granted to this service principal. Supports $expand.
func (m *ServicePrincipal) GetAppRoleAssignments()([]AppRoleAssignmentable) {
    val, err := m.GetBackingStore().Get("appRoleAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppRoleAssignmentable)
    }
    return nil
}
// GetAppRoles gets the appRoles property value. The roles exposed by the application which this service principal represents. For more information see the appRoles property definition on the application entity. Not nullable.
func (m *ServicePrincipal) GetAppRoles()([]AppRoleable) {
    val, err := m.GetBackingStore().Get("appRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppRoleable)
    }
    return nil
}
// GetClaimsMappingPolicies gets the claimsMappingPolicies property value. The claimsMappingPolicies assigned to this service principal. Supports $expand.
func (m *ServicePrincipal) GetClaimsMappingPolicies()([]ClaimsMappingPolicyable) {
    val, err := m.GetBackingStore().Get("claimsMappingPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ClaimsMappingPolicyable)
    }
    return nil
}
// GetCreatedObjects gets the createdObjects property value. Directory objects created by this service principal. Read-only. Nullable.
func (m *ServicePrincipal) GetCreatedObjects()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("createdObjects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetDelegatedPermissionClassifications gets the delegatedPermissionClassifications property value. The delegatedPermissionClassifications property
func (m *ServicePrincipal) GetDelegatedPermissionClassifications()([]DelegatedPermissionClassificationable) {
    val, err := m.GetBackingStore().Get("delegatedPermissionClassifications")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DelegatedPermissionClassificationable)
    }
    return nil
}
// GetDescription gets the description property value. Free text field to provide an internal end-user facing description of the service principal. End-user portals such MyApps will display the application description in this field. The maximum allowed size is 1024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.
func (m *ServicePrincipal) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisabledByMicrosoftStatus gets the disabledByMicrosoftStatus property value. Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).
func (m *ServicePrincipal) GetDisabledByMicrosoftStatus()(*string) {
    val, err := m.GetBackingStore().Get("disabledByMicrosoftStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name for the service principal. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderBy.
func (m *ServicePrincipal) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEndpoints gets the endpoints property value. The endpoints property
func (m *ServicePrincipal) GetEndpoints()([]Endpointable) {
    val, err := m.GetBackingStore().Get("endpoints")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Endpointable)
    }
    return nil
}
// GetFederatedIdentityCredentials gets the federatedIdentityCredentials property value. Federated identities for a specific type of service principal - managed identity. Supports $expand and $filter (/$count eq 0, /$count ne 0).
func (m *ServicePrincipal) GetFederatedIdentityCredentials()([]FederatedIdentityCredentialable) {
    val, err := m.GetBackingStore().Get("federatedIdentityCredentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]FederatedIdentityCredentialable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ServicePrincipal) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
    res["accountEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountEnabled(val)
        }
        return nil
    }
    res["addIns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAddInFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AddInable, len(val))
            for i, v := range val {
                res[i] = v.(AddInable)
            }
            m.SetAddIns(res)
        }
        return nil
    }
    res["alternativeNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetAlternativeNames(res)
        }
        return nil
    }
    res["appDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppDescription(val)
        }
        return nil
    }
    res["appDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppDisplayName(val)
        }
        return nil
    }
    res["appId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppId(val)
        }
        return nil
    }
    res["applicationTemplateId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationTemplateId(val)
        }
        return nil
    }
    res["appManagementPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppManagementPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppManagementPolicyable, len(val))
            for i, v := range val {
                res[i] = v.(AppManagementPolicyable)
            }
            m.SetAppManagementPolicies(res)
        }
        return nil
    }
    res["appOwnerOrganizationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppOwnerOrganizationId(val)
        }
        return nil
    }
    res["appRoleAssignedTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppRoleAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppRoleAssignmentable, len(val))
            for i, v := range val {
                res[i] = v.(AppRoleAssignmentable)
            }
            m.SetAppRoleAssignedTo(res)
        }
        return nil
    }
    res["appRoleAssignmentRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppRoleAssignmentRequired(val)
        }
        return nil
    }
    res["appRoleAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppRoleAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppRoleAssignmentable, len(val))
            for i, v := range val {
                res[i] = v.(AppRoleAssignmentable)
            }
            m.SetAppRoleAssignments(res)
        }
        return nil
    }
    res["appRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppRoleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppRoleable, len(val))
            for i, v := range val {
                res[i] = v.(AppRoleable)
            }
            m.SetAppRoles(res)
        }
        return nil
    }
    res["claimsMappingPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateClaimsMappingPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ClaimsMappingPolicyable, len(val))
            for i, v := range val {
                res[i] = v.(ClaimsMappingPolicyable)
            }
            m.SetClaimsMappingPolicies(res)
        }
        return nil
    }
    res["createdObjects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                res[i] = v.(DirectoryObjectable)
            }
            m.SetCreatedObjects(res)
        }
        return nil
    }
    res["delegatedPermissionClassifications"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDelegatedPermissionClassificationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DelegatedPermissionClassificationable, len(val))
            for i, v := range val {
                res[i] = v.(DelegatedPermissionClassificationable)
            }
            m.SetDelegatedPermissionClassifications(res)
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["disabledByMicrosoftStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisabledByMicrosoftStatus(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["endpoints"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEndpointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Endpointable, len(val))
            for i, v := range val {
                res[i] = v.(Endpointable)
            }
            m.SetEndpoints(res)
        }
        return nil
    }
    res["federatedIdentityCredentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateFederatedIdentityCredentialFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]FederatedIdentityCredentialable, len(val))
            for i, v := range val {
                res[i] = v.(FederatedIdentityCredentialable)
            }
            m.SetFederatedIdentityCredentials(res)
        }
        return nil
    }
    res["homepage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHomepage(val)
        }
        return nil
    }
    res["homeRealmDiscoveryPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateHomeRealmDiscoveryPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]HomeRealmDiscoveryPolicyable, len(val))
            for i, v := range val {
                res[i] = v.(HomeRealmDiscoveryPolicyable)
            }
            m.SetHomeRealmDiscoveryPolicies(res)
        }
        return nil
    }
    res["info"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateInformationalUrlFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInfo(val.(InformationalUrlable))
        }
        return nil
    }
    res["keyCredentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyCredentialFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyCredentialable, len(val))
            for i, v := range val {
                res[i] = v.(KeyCredentialable)
            }
            m.SetKeyCredentials(res)
        }
        return nil
    }
    res["loginUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLoginUrl(val)
        }
        return nil
    }
    res["logoutUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogoutUrl(val)
        }
        return nil
    }
    res["memberOf"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                res[i] = v.(DirectoryObjectable)
            }
            m.SetMemberOf(res)
        }
        return nil
    }
    res["notes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotes(val)
        }
        return nil
    }
    res["notificationEmailAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetNotificationEmailAddresses(res)
        }
        return nil
    }
    res["oauth2PermissionGrants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOAuth2PermissionGrantFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OAuth2PermissionGrantable, len(val))
            for i, v := range val {
                res[i] = v.(OAuth2PermissionGrantable)
            }
            m.SetOauth2PermissionGrants(res)
        }
        return nil
    }
    res["oauth2PermissionScopes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePermissionScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PermissionScopeable, len(val))
            for i, v := range val {
                res[i] = v.(PermissionScopeable)
            }
            m.SetOauth2PermissionScopes(res)
        }
        return nil
    }
    res["ownedObjects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                res[i] = v.(DirectoryObjectable)
            }
            m.SetOwnedObjects(res)
        }
        return nil
    }
    res["owners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                res[i] = v.(DirectoryObjectable)
            }
            m.SetOwners(res)
        }
        return nil
    }
    res["passwordCredentials"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePasswordCredentialFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PasswordCredentialable, len(val))
            for i, v := range val {
                res[i] = v.(PasswordCredentialable)
            }
            m.SetPasswordCredentials(res)
        }
        return nil
    }
    res["preferredSingleSignOnMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreferredSingleSignOnMode(val)
        }
        return nil
    }
    res["preferredTokenSigningKeyThumbprint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreferredTokenSigningKeyThumbprint(val)
        }
        return nil
    }
    res["replyUrls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetReplyUrls(res)
        }
        return nil
    }
    res["resourceSpecificApplicationPermissions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateResourceSpecificPermissionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ResourceSpecificPermissionable, len(val))
            for i, v := range val {
                res[i] = v.(ResourceSpecificPermissionable)
            }
            m.SetResourceSpecificApplicationPermissions(res)
        }
        return nil
    }
    res["samlSingleSignOnSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSamlSingleSignOnSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSamlSingleSignOnSettings(val.(SamlSingleSignOnSettingsable))
        }
        return nil
    }
    res["servicePrincipalNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetServicePrincipalNames(res)
        }
        return nil
    }
    res["servicePrincipalType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalType(val)
        }
        return nil
    }
    res["signInAudience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignInAudience(val)
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                res[i] = *(v.(*string))
            }
            m.SetTags(res)
        }
        return nil
    }
    res["tokenEncryptionKeyId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTokenEncryptionKeyId(val)
        }
        return nil
    }
    res["tokenIssuancePolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTokenIssuancePolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TokenIssuancePolicyable, len(val))
            for i, v := range val {
                res[i] = v.(TokenIssuancePolicyable)
            }
            m.SetTokenIssuancePolicies(res)
        }
        return nil
    }
    res["tokenLifetimePolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTokenLifetimePolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TokenLifetimePolicyable, len(val))
            for i, v := range val {
                res[i] = v.(TokenLifetimePolicyable)
            }
            m.SetTokenLifetimePolicies(res)
        }
        return nil
    }
    res["transitiveMemberOf"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                res[i] = v.(DirectoryObjectable)
            }
            m.SetTransitiveMemberOf(res)
        }
        return nil
    }
    res["verifiedPublisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVerifiedPublisherFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVerifiedPublisher(val.(VerifiedPublisherable))
        }
        return nil
    }
    return res
}
// GetHomepage gets the homepage property value. Home page or landing page of the application.
func (m *ServicePrincipal) GetHomepage()(*string) {
    val, err := m.GetBackingStore().Get("homepage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHomeRealmDiscoveryPolicies gets the homeRealmDiscoveryPolicies property value. The homeRealmDiscoveryPolicies assigned to this service principal. Supports $expand.
func (m *ServicePrincipal) GetHomeRealmDiscoveryPolicies()([]HomeRealmDiscoveryPolicyable) {
    val, err := m.GetBackingStore().Get("homeRealmDiscoveryPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]HomeRealmDiscoveryPolicyable)
    }
    return nil
}
// GetInfo gets the info property value. Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Azure AD apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).
func (m *ServicePrincipal) GetInfo()(InformationalUrlable) {
    val, err := m.GetBackingStore().Get("info")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(InformationalUrlable)
    }
    return nil
}
// GetKeyCredentials gets the keyCredentials property value. The collection of key credentials associated with the service principal. Not nullable. Supports $filter (eq, not, ge, le).
func (m *ServicePrincipal) GetKeyCredentials()([]KeyCredentialable) {
    val, err := m.GetBackingStore().Get("keyCredentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyCredentialable)
    }
    return nil
}
// GetLoginUrl gets the loginUrl property value. Specifies the URL where the service provider redirects the user to Azure AD to authenticate. Azure AD uses the URL to launch the application from Microsoft 365 or the Azure AD My Apps. When blank, Azure AD performs IdP-initiated sign-on for applications configured with SAML-based single sign-on. The user launches the application from Microsoft 365, the Azure AD My Apps, or the Azure AD SSO URL.
func (m *ServicePrincipal) GetLoginUrl()(*string) {
    val, err := m.GetBackingStore().Get("loginUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLogoutUrl gets the logoutUrl property value. Specifies the URL that will be used by Microsoft's authorization service to logout an user using OpenId Connect front-channel, back-channel or SAML logout protocols.
func (m *ServicePrincipal) GetLogoutUrl()(*string) {
    val, err := m.GetBackingStore().Get("logoutUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMemberOf gets the memberOf property value. Roles that this service principal is a member of. HTTP Methods: GET Read-only. Nullable. Supports $expand.
func (m *ServicePrincipal) GetMemberOf()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("memberOf")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetNotes gets the notes property value. Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1024 characters.
func (m *ServicePrincipal) GetNotes()(*string) {
    val, err := m.GetBackingStore().Get("notes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotificationEmailAddresses gets the notificationEmailAddresses property value. Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications.
func (m *ServicePrincipal) GetNotificationEmailAddresses()([]string) {
    val, err := m.GetBackingStore().Get("notificationEmailAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetOauth2PermissionGrants gets the oauth2PermissionGrants property value. Delegated permission grants authorizing this service principal to access an API on behalf of a signed-in user. Read-only. Nullable.
func (m *ServicePrincipal) GetOauth2PermissionGrants()([]OAuth2PermissionGrantable) {
    val, err := m.GetBackingStore().Get("oauth2PermissionGrants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OAuth2PermissionGrantable)
    }
    return nil
}
// GetOauth2PermissionScopes gets the oauth2PermissionScopes property value. The delegated permissions exposed by the application. For more information see the oauth2PermissionScopes property on the application entity's api property. Not nullable.
func (m *ServicePrincipal) GetOauth2PermissionScopes()([]PermissionScopeable) {
    val, err := m.GetBackingStore().Get("oauth2PermissionScopes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PermissionScopeable)
    }
    return nil
}
// GetOwnedObjects gets the ownedObjects property value. Directory objects that are owned by this service principal. Read-only. Nullable. Supports $expand and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
func (m *ServicePrincipal) GetOwnedObjects()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("ownedObjects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetOwners gets the owners property value. Directory objects that are owners of this servicePrincipal. The owners are a set of non-admin users or servicePrincipals who are allowed to modify this object. Read-only. Nullable.  Supports $expand and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
func (m *ServicePrincipal) GetOwners()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("owners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetPasswordCredentials gets the passwordCredentials property value. The collection of password credentials associated with the application. Not nullable.
func (m *ServicePrincipal) GetPasswordCredentials()([]PasswordCredentialable) {
    val, err := m.GetBackingStore().Get("passwordCredentials")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PasswordCredentialable)
    }
    return nil
}
// GetPreferredSingleSignOnMode gets the preferredSingleSignOnMode property value. Specifies the single sign-on mode configured for this application. Azure AD uses the preferred single sign-on mode to launch the application from Microsoft 365 or the Azure AD My Apps. The supported values are password, saml, notSupported, and oidc.
func (m *ServicePrincipal) GetPreferredSingleSignOnMode()(*string) {
    val, err := m.GetBackingStore().Get("preferredSingleSignOnMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPreferredTokenSigningKeyThumbprint gets the preferredTokenSigningKeyThumbprint property value. Reserved for internal use only. Do not write or otherwise rely on this property. May be removed in future versions.
func (m *ServicePrincipal) GetPreferredTokenSigningKeyThumbprint()(*string) {
    val, err := m.GetBackingStore().Get("preferredTokenSigningKeyThumbprint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReplyUrls gets the replyUrls property value. The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application. Not nullable.
func (m *ServicePrincipal) GetReplyUrls()([]string) {
    val, err := m.GetBackingStore().Get("replyUrls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetResourceSpecificApplicationPermissions gets the resourceSpecificApplicationPermissions property value. The resource-specific application permissions exposed by this application. Currently, resource-specific permissions are only supported for Teams apps accessing to specific chats and teams using Microsoft Graph. Read-only.
func (m *ServicePrincipal) GetResourceSpecificApplicationPermissions()([]ResourceSpecificPermissionable) {
    val, err := m.GetBackingStore().Get("resourceSpecificApplicationPermissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ResourceSpecificPermissionable)
    }
    return nil
}
// GetSamlSingleSignOnSettings gets the samlSingleSignOnSettings property value. The collection for settings related to saml single sign-on.
func (m *ServicePrincipal) GetSamlSingleSignOnSettings()(SamlSingleSignOnSettingsable) {
    val, err := m.GetBackingStore().Get("samlSingleSignOnSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SamlSingleSignOnSettingsable)
    }
    return nil
}
// GetServicePrincipalNames gets the servicePrincipalNames property value. Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD. For example,Client apps can specify a resource URI which is based on the values of this property to acquire an access token, which is the URI returned in the 'aud' claim.The any operator is required for filter expressions on multi-valued properties. Not nullable.  Supports $filter (eq, not, ge, le, startsWith).
func (m *ServicePrincipal) GetServicePrincipalNames()([]string) {
    val, err := m.GetBackingStore().Get("servicePrincipalNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetServicePrincipalType gets the servicePrincipalType property value. Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Azure AD internally. The servicePrincipalType property can be set to three different values: __Application - A service principal that represents an application or service. The appId property identifies the associated app registration, and matches the appId of an application, possibly from a different tenant. If the associated app registration is missing, tokens are not issued for the service principal.__ManagedIdentity - A service principal that represents a managed identity. Service principals representing managed identities can be granted access and permissions, but cannot be updated or modified directly.__Legacy - A service principal that represents an app created before app registrations, or through legacy experiences. Legacy service principal can have credentials, service principal names, reply URLs, and other properties which are editable by an authorized user, but does not have an associated app registration. The appId value does not associate the service principal with an app registration. The service principal can only be used in the tenant where it was created.__SocialIdp - For internal use.
func (m *ServicePrincipal) GetServicePrincipalType()(*string) {
    val, err := m.GetBackingStore().Get("servicePrincipalType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSignInAudience gets the signInAudience property value. Specifies the Microsoft accounts that are supported for the current application. Read-only. Supported values are:AzureADMyOrg: Users with a Microsoft work or school account in my organization’s Azure AD tenant (single-tenant).AzureADMultipleOrgs: Users with a Microsoft work or school account in any organization’s Azure AD tenant (multi-tenant).AzureADandPersonalMicrosoftAccount: Users with a personal Microsoft account, or a work or school account in any organization’s Azure AD tenant.PersonalMicrosoftAccount: Users with a personal Microsoft account only.
func (m *ServicePrincipal) GetSignInAudience()(*string) {
    val, err := m.GetBackingStore().Get("signInAudience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTags gets the tags property value. Custom strings that can be used to categorize and identify the service principal. Not nullable. Supports $filter (eq, not, ge, le, startsWith).
func (m *ServicePrincipal) GetTags()([]string) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetTokenEncryptionKeyId gets the tokenEncryptionKeyId property value. Specifies the keyId of a public key from the keyCredentials collection. When configured, Azure AD issues tokens for this application encrypted using the key specified by this property. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.
func (m *ServicePrincipal) GetTokenEncryptionKeyId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("tokenEncryptionKeyId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetTokenIssuancePolicies gets the tokenIssuancePolicies property value. The tokenIssuancePolicies assigned to this service principal.
func (m *ServicePrincipal) GetTokenIssuancePolicies()([]TokenIssuancePolicyable) {
    val, err := m.GetBackingStore().Get("tokenIssuancePolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TokenIssuancePolicyable)
    }
    return nil
}
// GetTokenLifetimePolicies gets the tokenLifetimePolicies property value. The tokenLifetimePolicies assigned to this service principal.
func (m *ServicePrincipal) GetTokenLifetimePolicies()([]TokenLifetimePolicyable) {
    val, err := m.GetBackingStore().Get("tokenLifetimePolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TokenLifetimePolicyable)
    }
    return nil
}
// GetTransitiveMemberOf gets the transitiveMemberOf property value. The transitiveMemberOf property
func (m *ServicePrincipal) GetTransitiveMemberOf()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("transitiveMemberOf")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetVerifiedPublisher gets the verifiedPublisher property value. Specifies the verified publisher of the application which this service principal represents.
func (m *ServicePrincipal) GetVerifiedPublisher()(VerifiedPublisherable) {
    val, err := m.GetBackingStore().Get("verifiedPublisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(VerifiedPublisherable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ServicePrincipal) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("accountEnabled", m.GetAccountEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetAddIns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAddIns()))
        for i, v := range m.GetAddIns() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("addIns", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAlternativeNames() != nil {
        err = writer.WriteCollectionOfStringValues("alternativeNames", m.GetAlternativeNames())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appDescription", m.GetAppDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appDisplayName", m.GetAppDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("applicationTemplateId", m.GetApplicationTemplateId())
        if err != nil {
            return err
        }
    }
    if m.GetAppManagementPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppManagementPolicies()))
        for i, v := range m.GetAppManagementPolicies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("appManagementPolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("appOwnerOrganizationId", m.GetAppOwnerOrganizationId())
        if err != nil {
            return err
        }
    }
    if m.GetAppRoleAssignedTo() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppRoleAssignedTo()))
        for i, v := range m.GetAppRoleAssignedTo() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("appRoleAssignedTo", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appRoleAssignmentRequired", m.GetAppRoleAssignmentRequired())
        if err != nil {
            return err
        }
    }
    if m.GetAppRoleAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppRoleAssignments()))
        for i, v := range m.GetAppRoleAssignments() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("appRoleAssignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAppRoles() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppRoles()))
        for i, v := range m.GetAppRoles() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("appRoles", cast)
        if err != nil {
            return err
        }
    }
    if m.GetClaimsMappingPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetClaimsMappingPolicies()))
        for i, v := range m.GetClaimsMappingPolicies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("claimsMappingPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCreatedObjects() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCreatedObjects()))
        for i, v := range m.GetCreatedObjects() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("createdObjects", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDelegatedPermissionClassifications() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDelegatedPermissionClassifications()))
        for i, v := range m.GetDelegatedPermissionClassifications() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("delegatedPermissionClassifications", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("disabledByMicrosoftStatus", m.GetDisabledByMicrosoftStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetEndpoints() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEndpoints()))
        for i, v := range m.GetEndpoints() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("endpoints", cast)
        if err != nil {
            return err
        }
    }
    if m.GetFederatedIdentityCredentials() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFederatedIdentityCredentials()))
        for i, v := range m.GetFederatedIdentityCredentials() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("federatedIdentityCredentials", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("homepage", m.GetHomepage())
        if err != nil {
            return err
        }
    }
    if m.GetHomeRealmDiscoveryPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHomeRealmDiscoveryPolicies()))
        for i, v := range m.GetHomeRealmDiscoveryPolicies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("homeRealmDiscoveryPolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("info", m.GetInfo())
        if err != nil {
            return err
        }
    }
    if m.GetKeyCredentials() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetKeyCredentials()))
        for i, v := range m.GetKeyCredentials() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("keyCredentials", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("loginUrl", m.GetLoginUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("logoutUrl", m.GetLogoutUrl())
        if err != nil {
            return err
        }
    }
    if m.GetMemberOf() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMemberOf()))
        for i, v := range m.GetMemberOf() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("memberOf", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notes", m.GetNotes())
        if err != nil {
            return err
        }
    }
    if m.GetNotificationEmailAddresses() != nil {
        err = writer.WriteCollectionOfStringValues("notificationEmailAddresses", m.GetNotificationEmailAddresses())
        if err != nil {
            return err
        }
    }
    if m.GetOauth2PermissionGrants() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOauth2PermissionGrants()))
        for i, v := range m.GetOauth2PermissionGrants() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("oauth2PermissionGrants", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOauth2PermissionScopes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOauth2PermissionScopes()))
        for i, v := range m.GetOauth2PermissionScopes() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("oauth2PermissionScopes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOwnedObjects() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOwnedObjects()))
        for i, v := range m.GetOwnedObjects() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("ownedObjects", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOwners() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOwners()))
        for i, v := range m.GetOwners() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("owners", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPasswordCredentials() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPasswordCredentials()))
        for i, v := range m.GetPasswordCredentials() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("passwordCredentials", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("preferredSingleSignOnMode", m.GetPreferredSingleSignOnMode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("preferredTokenSigningKeyThumbprint", m.GetPreferredTokenSigningKeyThumbprint())
        if err != nil {
            return err
        }
    }
    if m.GetReplyUrls() != nil {
        err = writer.WriteCollectionOfStringValues("replyUrls", m.GetReplyUrls())
        if err != nil {
            return err
        }
    }
    if m.GetResourceSpecificApplicationPermissions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResourceSpecificApplicationPermissions()))
        for i, v := range m.GetResourceSpecificApplicationPermissions() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("resourceSpecificApplicationPermissions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("samlSingleSignOnSettings", m.GetSamlSingleSignOnSettings())
        if err != nil {
            return err
        }
    }
    if m.GetServicePrincipalNames() != nil {
        err = writer.WriteCollectionOfStringValues("servicePrincipalNames", m.GetServicePrincipalNames())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("servicePrincipalType", m.GetServicePrincipalType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("signInAudience", m.GetSignInAudience())
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        err = writer.WriteCollectionOfStringValues("tags", m.GetTags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("tokenEncryptionKeyId", m.GetTokenEncryptionKeyId())
        if err != nil {
            return err
        }
    }
    if m.GetTokenIssuancePolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTokenIssuancePolicies()))
        for i, v := range m.GetTokenIssuancePolicies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("tokenIssuancePolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTokenLifetimePolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTokenLifetimePolicies()))
        for i, v := range m.GetTokenLifetimePolicies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("tokenLifetimePolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTransitiveMemberOf() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTransitiveMemberOf()))
        for i, v := range m.GetTransitiveMemberOf() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("transitiveMemberOf", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("verifiedPublisher", m.GetVerifiedPublisher())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccountEnabled sets the accountEnabled property value. true if the service principal account is enabled; otherwise, false. If set to false, then no users will be able to sign in to this app, even if they are assigned to it. Supports $filter (eq, ne, not, in).
func (m *ServicePrincipal) SetAccountEnabled(value *bool)() {
    err := m.GetBackingStore().Set("accountEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetAddIns sets the addIns property value. Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams may set the addIns property for its 'FileHandler' functionality. This will let services like Microsoft 365 call the application in the context of a document the user is working on.
func (m *ServicePrincipal) SetAddIns(value []AddInable)() {
    err := m.GetBackingStore().Set("addIns", value)
    if err != nil {
        panic(err)
    }
}
// SetAlternativeNames sets the alternativeNames property value. Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities. Supports $filter (eq, not, ge, le, startsWith).
func (m *ServicePrincipal) SetAlternativeNames(value []string)() {
    err := m.GetBackingStore().Set("alternativeNames", value)
    if err != nil {
        panic(err)
    }
}
// SetAppDescription sets the appDescription property value. The description exposed by the associated application.
func (m *ServicePrincipal) SetAppDescription(value *string)() {
    err := m.GetBackingStore().Set("appDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetAppDisplayName sets the appDisplayName property value. The display name exposed by the associated application.
func (m *ServicePrincipal) SetAppDisplayName(value *string)() {
    err := m.GetBackingStore().Set("appDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetAppId sets the appId property value. The unique identifier for the associated application (its appId property). Supports $filter (eq, ne, not, in, startsWith).
func (m *ServicePrincipal) SetAppId(value *string)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationTemplateId sets the applicationTemplateId property value. Unique identifier of the applicationTemplate that the servicePrincipal was created from. Read-only. Supports $filter (eq, ne, NOT, startsWith).
func (m *ServicePrincipal) SetApplicationTemplateId(value *string)() {
    err := m.GetBackingStore().Set("applicationTemplateId", value)
    if err != nil {
        panic(err)
    }
}
// SetAppManagementPolicies sets the appManagementPolicies property value. The appManagementPolicies property
func (m *ServicePrincipal) SetAppManagementPolicies(value []AppManagementPolicyable)() {
    err := m.GetBackingStore().Set("appManagementPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetAppOwnerOrganizationId sets the appOwnerOrganizationId property value. Contains the tenant id where the application is registered. This is applicable only to service principals backed by applications. Supports $filter (eq, ne, NOT, ge, le).
func (m *ServicePrincipal) SetAppOwnerOrganizationId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("appOwnerOrganizationId", value)
    if err != nil {
        panic(err)
    }
}
// SetAppRoleAssignedTo sets the appRoleAssignedTo property value. App role assignments for this app or service, granted to users, groups, and other service principals. Supports $expand.
func (m *ServicePrincipal) SetAppRoleAssignedTo(value []AppRoleAssignmentable)() {
    err := m.GetBackingStore().Set("appRoleAssignedTo", value)
    if err != nil {
        panic(err)
    }
}
// SetAppRoleAssignmentRequired sets the appRoleAssignmentRequired property value. Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports $filter (eq, ne, NOT).
func (m *ServicePrincipal) SetAppRoleAssignmentRequired(value *bool)() {
    err := m.GetBackingStore().Set("appRoleAssignmentRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetAppRoleAssignments sets the appRoleAssignments property value. App role assignment for another app or service, granted to this service principal. Supports $expand.
func (m *ServicePrincipal) SetAppRoleAssignments(value []AppRoleAssignmentable)() {
    err := m.GetBackingStore().Set("appRoleAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetAppRoles sets the appRoles property value. The roles exposed by the application which this service principal represents. For more information see the appRoles property definition on the application entity. Not nullable.
func (m *ServicePrincipal) SetAppRoles(value []AppRoleable)() {
    err := m.GetBackingStore().Set("appRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetClaimsMappingPolicies sets the claimsMappingPolicies property value. The claimsMappingPolicies assigned to this service principal. Supports $expand.
func (m *ServicePrincipal) SetClaimsMappingPolicies(value []ClaimsMappingPolicyable)() {
    err := m.GetBackingStore().Set("claimsMappingPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedObjects sets the createdObjects property value. Directory objects created by this service principal. Read-only. Nullable.
func (m *ServicePrincipal) SetCreatedObjects(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("createdObjects", value)
    if err != nil {
        panic(err)
    }
}
// SetDelegatedPermissionClassifications sets the delegatedPermissionClassifications property value. The delegatedPermissionClassifications property
func (m *ServicePrincipal) SetDelegatedPermissionClassifications(value []DelegatedPermissionClassificationable)() {
    err := m.GetBackingStore().Set("delegatedPermissionClassifications", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Free text field to provide an internal end-user facing description of the service principal. End-user portals such MyApps will display the application description in this field. The maximum allowed size is 1024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.
func (m *ServicePrincipal) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisabledByMicrosoftStatus sets the disabledByMicrosoftStatus property value. Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).
func (m *ServicePrincipal) SetDisabledByMicrosoftStatus(value *string)() {
    err := m.GetBackingStore().Set("disabledByMicrosoftStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name for the service principal. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderBy.
func (m *ServicePrincipal) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEndpoints sets the endpoints property value. The endpoints property
func (m *ServicePrincipal) SetEndpoints(value []Endpointable)() {
    err := m.GetBackingStore().Set("endpoints", value)
    if err != nil {
        panic(err)
    }
}
// SetFederatedIdentityCredentials sets the federatedIdentityCredentials property value. Federated identities for a specific type of service principal - managed identity. Supports $expand and $filter (/$count eq 0, /$count ne 0).
func (m *ServicePrincipal) SetFederatedIdentityCredentials(value []FederatedIdentityCredentialable)() {
    err := m.GetBackingStore().Set("federatedIdentityCredentials", value)
    if err != nil {
        panic(err)
    }
}
// SetHomepage sets the homepage property value. Home page or landing page of the application.
func (m *ServicePrincipal) SetHomepage(value *string)() {
    err := m.GetBackingStore().Set("homepage", value)
    if err != nil {
        panic(err)
    }
}
// SetHomeRealmDiscoveryPolicies sets the homeRealmDiscoveryPolicies property value. The homeRealmDiscoveryPolicies assigned to this service principal. Supports $expand.
func (m *ServicePrincipal) SetHomeRealmDiscoveryPolicies(value []HomeRealmDiscoveryPolicyable)() {
    err := m.GetBackingStore().Set("homeRealmDiscoveryPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetInfo sets the info property value. Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Azure AD apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).
func (m *ServicePrincipal) SetInfo(value InformationalUrlable)() {
    err := m.GetBackingStore().Set("info", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyCredentials sets the keyCredentials property value. The collection of key credentials associated with the service principal. Not nullable. Supports $filter (eq, not, ge, le).
func (m *ServicePrincipal) SetKeyCredentials(value []KeyCredentialable)() {
    err := m.GetBackingStore().Set("keyCredentials", value)
    if err != nil {
        panic(err)
    }
}
// SetLoginUrl sets the loginUrl property value. Specifies the URL where the service provider redirects the user to Azure AD to authenticate. Azure AD uses the URL to launch the application from Microsoft 365 or the Azure AD My Apps. When blank, Azure AD performs IdP-initiated sign-on for applications configured with SAML-based single sign-on. The user launches the application from Microsoft 365, the Azure AD My Apps, or the Azure AD SSO URL.
func (m *ServicePrincipal) SetLoginUrl(value *string)() {
    err := m.GetBackingStore().Set("loginUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetLogoutUrl sets the logoutUrl property value. Specifies the URL that will be used by Microsoft's authorization service to logout an user using OpenId Connect front-channel, back-channel or SAML logout protocols.
func (m *ServicePrincipal) SetLogoutUrl(value *string)() {
    err := m.GetBackingStore().Set("logoutUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberOf sets the memberOf property value. Roles that this service principal is a member of. HTTP Methods: GET Read-only. Nullable. Supports $expand.
func (m *ServicePrincipal) SetMemberOf(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("memberOf", value)
    if err != nil {
        panic(err)
    }
}
// SetNotes sets the notes property value. Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1024 characters.
func (m *ServicePrincipal) SetNotes(value *string)() {
    err := m.GetBackingStore().Set("notes", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationEmailAddresses sets the notificationEmailAddresses property value. Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications.
func (m *ServicePrincipal) SetNotificationEmailAddresses(value []string)() {
    err := m.GetBackingStore().Set("notificationEmailAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetOauth2PermissionGrants sets the oauth2PermissionGrants property value. Delegated permission grants authorizing this service principal to access an API on behalf of a signed-in user. Read-only. Nullable.
func (m *ServicePrincipal) SetOauth2PermissionGrants(value []OAuth2PermissionGrantable)() {
    err := m.GetBackingStore().Set("oauth2PermissionGrants", value)
    if err != nil {
        panic(err)
    }
}
// SetOauth2PermissionScopes sets the oauth2PermissionScopes property value. The delegated permissions exposed by the application. For more information see the oauth2PermissionScopes property on the application entity's api property. Not nullable.
func (m *ServicePrincipal) SetOauth2PermissionScopes(value []PermissionScopeable)() {
    err := m.GetBackingStore().Set("oauth2PermissionScopes", value)
    if err != nil {
        panic(err)
    }
}
// SetOwnedObjects sets the ownedObjects property value. Directory objects that are owned by this service principal. Read-only. Nullable. Supports $expand and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
func (m *ServicePrincipal) SetOwnedObjects(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("ownedObjects", value)
    if err != nil {
        panic(err)
    }
}
// SetOwners sets the owners property value. Directory objects that are owners of this servicePrincipal. The owners are a set of non-admin users or servicePrincipals who are allowed to modify this object. Read-only. Nullable.  Supports $expand and $filter (/$count eq 0, /$count ne 0, /$count eq 1, /$count ne 1).
func (m *ServicePrincipal) SetOwners(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("owners", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordCredentials sets the passwordCredentials property value. The collection of password credentials associated with the application. Not nullable.
func (m *ServicePrincipal) SetPasswordCredentials(value []PasswordCredentialable)() {
    err := m.GetBackingStore().Set("passwordCredentials", value)
    if err != nil {
        panic(err)
    }
}
// SetPreferredSingleSignOnMode sets the preferredSingleSignOnMode property value. Specifies the single sign-on mode configured for this application. Azure AD uses the preferred single sign-on mode to launch the application from Microsoft 365 or the Azure AD My Apps. The supported values are password, saml, notSupported, and oidc.
func (m *ServicePrincipal) SetPreferredSingleSignOnMode(value *string)() {
    err := m.GetBackingStore().Set("preferredSingleSignOnMode", value)
    if err != nil {
        panic(err)
    }
}
// SetPreferredTokenSigningKeyThumbprint sets the preferredTokenSigningKeyThumbprint property value. Reserved for internal use only. Do not write or otherwise rely on this property. May be removed in future versions.
func (m *ServicePrincipal) SetPreferredTokenSigningKeyThumbprint(value *string)() {
    err := m.GetBackingStore().Set("preferredTokenSigningKeyThumbprint", value)
    if err != nil {
        panic(err)
    }
}
// SetReplyUrls sets the replyUrls property value. The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application. Not nullable.
func (m *ServicePrincipal) SetReplyUrls(value []string)() {
    err := m.GetBackingStore().Set("replyUrls", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceSpecificApplicationPermissions sets the resourceSpecificApplicationPermissions property value. The resource-specific application permissions exposed by this application. Currently, resource-specific permissions are only supported for Teams apps accessing to specific chats and teams using Microsoft Graph. Read-only.
func (m *ServicePrincipal) SetResourceSpecificApplicationPermissions(value []ResourceSpecificPermissionable)() {
    err := m.GetBackingStore().Set("resourceSpecificApplicationPermissions", value)
    if err != nil {
        panic(err)
    }
}
// SetSamlSingleSignOnSettings sets the samlSingleSignOnSettings property value. The collection for settings related to saml single sign-on.
func (m *ServicePrincipal) SetSamlSingleSignOnSettings(value SamlSingleSignOnSettingsable)() {
    err := m.GetBackingStore().Set("samlSingleSignOnSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalNames sets the servicePrincipalNames property value. Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD. For example,Client apps can specify a resource URI which is based on the values of this property to acquire an access token, which is the URI returned in the 'aud' claim.The any operator is required for filter expressions on multi-valued properties. Not nullable.  Supports $filter (eq, not, ge, le, startsWith).
func (m *ServicePrincipal) SetServicePrincipalNames(value []string)() {
    err := m.GetBackingStore().Set("servicePrincipalNames", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalType sets the servicePrincipalType property value. Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Azure AD internally. The servicePrincipalType property can be set to three different values: __Application - A service principal that represents an application or service. The appId property identifies the associated app registration, and matches the appId of an application, possibly from a different tenant. If the associated app registration is missing, tokens are not issued for the service principal.__ManagedIdentity - A service principal that represents a managed identity. Service principals representing managed identities can be granted access and permissions, but cannot be updated or modified directly.__Legacy - A service principal that represents an app created before app registrations, or through legacy experiences. Legacy service principal can have credentials, service principal names, reply URLs, and other properties which are editable by an authorized user, but does not have an associated app registration. The appId value does not associate the service principal with an app registration. The service principal can only be used in the tenant where it was created.__SocialIdp - For internal use.
func (m *ServicePrincipal) SetServicePrincipalType(value *string)() {
    err := m.GetBackingStore().Set("servicePrincipalType", value)
    if err != nil {
        panic(err)
    }
}
// SetSignInAudience sets the signInAudience property value. Specifies the Microsoft accounts that are supported for the current application. Read-only. Supported values are:AzureADMyOrg: Users with a Microsoft work or school account in my organization’s Azure AD tenant (single-tenant).AzureADMultipleOrgs: Users with a Microsoft work or school account in any organization’s Azure AD tenant (multi-tenant).AzureADandPersonalMicrosoftAccount: Users with a personal Microsoft account, or a work or school account in any organization’s Azure AD tenant.PersonalMicrosoftAccount: Users with a personal Microsoft account only.
func (m *ServicePrincipal) SetSignInAudience(value *string)() {
    err := m.GetBackingStore().Set("signInAudience", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. Custom strings that can be used to categorize and identify the service principal. Not nullable. Supports $filter (eq, not, ge, le, startsWith).
func (m *ServicePrincipal) SetTags(value []string)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenEncryptionKeyId sets the tokenEncryptionKeyId property value. Specifies the keyId of a public key from the keyCredentials collection. When configured, Azure AD issues tokens for this application encrypted using the key specified by this property. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.
func (m *ServicePrincipal) SetTokenEncryptionKeyId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("tokenEncryptionKeyId", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenIssuancePolicies sets the tokenIssuancePolicies property value. The tokenIssuancePolicies assigned to this service principal.
func (m *ServicePrincipal) SetTokenIssuancePolicies(value []TokenIssuancePolicyable)() {
    err := m.GetBackingStore().Set("tokenIssuancePolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetTokenLifetimePolicies sets the tokenLifetimePolicies property value. The tokenLifetimePolicies assigned to this service principal.
func (m *ServicePrincipal) SetTokenLifetimePolicies(value []TokenLifetimePolicyable)() {
    err := m.GetBackingStore().Set("tokenLifetimePolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetTransitiveMemberOf sets the transitiveMemberOf property value. The transitiveMemberOf property
func (m *ServicePrincipal) SetTransitiveMemberOf(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("transitiveMemberOf", value)
    if err != nil {
        panic(err)
    }
}
// SetVerifiedPublisher sets the verifiedPublisher property value. Specifies the verified publisher of the application which this service principal represents.
func (m *ServicePrincipal) SetVerifiedPublisher(value VerifiedPublisherable)() {
    err := m.GetBackingStore().Set("verifiedPublisher", value)
    if err != nil {
        panic(err)
    }
}
// ServicePrincipalable 
type ServicePrincipalable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountEnabled()(*bool)
    GetAddIns()([]AddInable)
    GetAlternativeNames()([]string)
    GetAppDescription()(*string)
    GetAppDisplayName()(*string)
    GetAppId()(*string)
    GetApplicationTemplateId()(*string)
    GetAppManagementPolicies()([]AppManagementPolicyable)
    GetAppOwnerOrganizationId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetAppRoleAssignedTo()([]AppRoleAssignmentable)
    GetAppRoleAssignmentRequired()(*bool)
    GetAppRoleAssignments()([]AppRoleAssignmentable)
    GetAppRoles()([]AppRoleable)
    GetClaimsMappingPolicies()([]ClaimsMappingPolicyable)
    GetCreatedObjects()([]DirectoryObjectable)
    GetDelegatedPermissionClassifications()([]DelegatedPermissionClassificationable)
    GetDescription()(*string)
    GetDisabledByMicrosoftStatus()(*string)
    GetDisplayName()(*string)
    GetEndpoints()([]Endpointable)
    GetFederatedIdentityCredentials()([]FederatedIdentityCredentialable)
    GetHomepage()(*string)
    GetHomeRealmDiscoveryPolicies()([]HomeRealmDiscoveryPolicyable)
    GetInfo()(InformationalUrlable)
    GetKeyCredentials()([]KeyCredentialable)
    GetLoginUrl()(*string)
    GetLogoutUrl()(*string)
    GetMemberOf()([]DirectoryObjectable)
    GetNotes()(*string)
    GetNotificationEmailAddresses()([]string)
    GetOauth2PermissionGrants()([]OAuth2PermissionGrantable)
    GetOauth2PermissionScopes()([]PermissionScopeable)
    GetOwnedObjects()([]DirectoryObjectable)
    GetOwners()([]DirectoryObjectable)
    GetPasswordCredentials()([]PasswordCredentialable)
    GetPreferredSingleSignOnMode()(*string)
    GetPreferredTokenSigningKeyThumbprint()(*string)
    GetReplyUrls()([]string)
    GetResourceSpecificApplicationPermissions()([]ResourceSpecificPermissionable)
    GetSamlSingleSignOnSettings()(SamlSingleSignOnSettingsable)
    GetServicePrincipalNames()([]string)
    GetServicePrincipalType()(*string)
    GetSignInAudience()(*string)
    GetTags()([]string)
    GetTokenEncryptionKeyId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetTokenIssuancePolicies()([]TokenIssuancePolicyable)
    GetTokenLifetimePolicies()([]TokenLifetimePolicyable)
    GetTransitiveMemberOf()([]DirectoryObjectable)
    GetVerifiedPublisher()(VerifiedPublisherable)
    SetAccountEnabled(value *bool)()
    SetAddIns(value []AddInable)()
    SetAlternativeNames(value []string)()
    SetAppDescription(value *string)()
    SetAppDisplayName(value *string)()
    SetAppId(value *string)()
    SetApplicationTemplateId(value *string)()
    SetAppManagementPolicies(value []AppManagementPolicyable)()
    SetAppOwnerOrganizationId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetAppRoleAssignedTo(value []AppRoleAssignmentable)()
    SetAppRoleAssignmentRequired(value *bool)()
    SetAppRoleAssignments(value []AppRoleAssignmentable)()
    SetAppRoles(value []AppRoleable)()
    SetClaimsMappingPolicies(value []ClaimsMappingPolicyable)()
    SetCreatedObjects(value []DirectoryObjectable)()
    SetDelegatedPermissionClassifications(value []DelegatedPermissionClassificationable)()
    SetDescription(value *string)()
    SetDisabledByMicrosoftStatus(value *string)()
    SetDisplayName(value *string)()
    SetEndpoints(value []Endpointable)()
    SetFederatedIdentityCredentials(value []FederatedIdentityCredentialable)()
    SetHomepage(value *string)()
    SetHomeRealmDiscoveryPolicies(value []HomeRealmDiscoveryPolicyable)()
    SetInfo(value InformationalUrlable)()
    SetKeyCredentials(value []KeyCredentialable)()
    SetLoginUrl(value *string)()
    SetLogoutUrl(value *string)()
    SetMemberOf(value []DirectoryObjectable)()
    SetNotes(value *string)()
    SetNotificationEmailAddresses(value []string)()
    SetOauth2PermissionGrants(value []OAuth2PermissionGrantable)()
    SetOauth2PermissionScopes(value []PermissionScopeable)()
    SetOwnedObjects(value []DirectoryObjectable)()
    SetOwners(value []DirectoryObjectable)()
    SetPasswordCredentials(value []PasswordCredentialable)()
    SetPreferredSingleSignOnMode(value *string)()
    SetPreferredTokenSigningKeyThumbprint(value *string)()
    SetReplyUrls(value []string)()
    SetResourceSpecificApplicationPermissions(value []ResourceSpecificPermissionable)()
    SetSamlSingleSignOnSettings(value SamlSingleSignOnSettingsable)()
    SetServicePrincipalNames(value []string)()
    SetServicePrincipalType(value *string)()
    SetSignInAudience(value *string)()
    SetTags(value []string)()
    SetTokenEncryptionKeyId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetTokenIssuancePolicies(value []TokenIssuancePolicyable)()
    SetTokenLifetimePolicies(value []TokenLifetimePolicyable)()
    SetTransitiveMemberOf(value []DirectoryObjectable)()
    SetVerifiedPublisher(value VerifiedPublisherable)()
}
